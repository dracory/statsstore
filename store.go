package statsstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/req"
	"github.com/dromara/carbon/v2"
)

// == INTERFACE ================================================================

var _ StoreInterface = (*storeImplementation)(nil)

// storeImplementation implements StoreInterface for visitor operations.
type storeImplementation struct {
	visitorTableName     string
	settingsTableName    string
	db                   *neat.Database
	automigrateEnabled   bool
	debugEnabled         bool
	botFilterEnabled     bool
	excludedPathPrefixes []string
	excludedIPs          []string
	geoIPResolver        GeoIPResolver
	enhanceBatchSize     int
	logger               *slog.Logger
}

// == MIGRATE ==================================================================

// MigrateUp creates the visitor table and excluded IPs table if they do not already exist.
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.visitorTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: table already exists", "table", st.visitorTableName)
		}
	} else {
		err := st.db.Schema().Create(st.visitorTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_PATH, 510)
			table.String(COLUMN_FINGERPRINT, 40)
			table.String(COLUMN_IP_ADDRESS, 40)
			table.String(COLUMN_COUNTRY, 2)
			table.String(COLUMN_USER_ACCEPT_LANGUAGE, 100)
			table.String(COLUMN_USER_ACCEPT_ENCODING, 40)
			table.String(COLUMN_USER_AGENT, 510)
			table.String(COLUMN_USER_OS, 12)
			table.String(COLUMN_USER_OS_VERSION, 12)
			table.String(COLUMN_USER_DEVICE, 40)
			table.String(COLUMN_USER_DEVICE_TYPE, 12)
			table.String(COLUMN_USER_BROWSER, 40)
			table.String(COLUMN_USER_BROWSER_VERSION, 24)
			table.String(COLUMN_USER_REFERRER, 510)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)

			table.Index(COLUMN_CREATED_AT)
			table.Index(COLUMN_IP_ADDRESS)
			table.Index(COLUMN_FINGERPRINT)
		})

		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateUp failed", "error", err)
			}
			return err
		}
	}

	if st.settingsTableName != "" && !st.db.Schema().HasTable(st.settingsTableName) {
		err := st.db.Schema().Create(st.settingsTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_KEY, 100)
			table.Primary(COLUMN_KEY)
			table.Text(COLUMN_VALUE)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
		})

		if err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateUp: settings table creation failed", "error", err)
			}
			return err
		}
	}

	return nil
}

// MigrateDown drops the visitor table and settings table.
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if st.settingsTableName != "" && st.db.Schema().HasTable(st.settingsTableName) {
		if err := st.db.Schema().Drop(st.settingsTableName); err != nil {
			if st.debugEnabled {
				st.logger.Error("MigrateDown: settings table drop failed", "error", err)
			}
			return err
		}
	}

	if !st.db.Schema().HasTable(st.visitorTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateDown: table does not exist", "table", st.visitorTableName)
		}
		return nil
	}

	err := st.db.Schema().Drop(st.visitorTableName)
	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateDown failed", "error", err)
		}
		return err
	}
	return nil
}

// == DEBUG ====================================================================

// EnableDebug enables or disables debug mode.
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// == DB =======================================================================

// GetDB returns the underlying *sql.DB.
func (st *storeImplementation) GetDB() *sql.DB {
	db, _ := st.db.DB()
	return db
}

// == BOT FILTERING ===========================================================

// SetBotFilterEnabled enables or disables bot/referrer spam/data center filtering.
func (st *storeImplementation) SetBotFilterEnabled(enabled bool) {
	st.botFilterEnabled = enabled
}

// IsBotFilterEnabled returns whether bot filtering is currently enabled.
func (st *storeImplementation) IsBotFilterEnabled() bool {
	return st.botFilterEnabled
}

// SetExcludedPathPrefixes sets path prefixes that should be excluded from visitor tracking.
// Requests whose URL path starts with any of these prefixes will be silently skipped
// by VisitorRegister. This is useful for excluding admin panel traffic (e.g. "/admin/").
func (st *storeImplementation) SetExcludedPathPrefixes(prefixes []string) {
	st.excludedPathPrefixes = prefixes
}

// GetExcludedPathPrefixes returns the currently configured excluded path prefixes.
func (st *storeImplementation) GetExcludedPathPrefixes() []string {
	return st.excludedPathPrefixes
}

// SetExcludedIPs sets IP addresses that should be excluded from visitor tracking.
// Requests from these IPs will be silently skipped by VisitorRegister.
func (st *storeImplementation) SetExcludedIPs(ips []string) {
	st.excludedIPs = ips
}

// GetExcludedIPs returns the currently configured excluded IP addresses.
func (st *storeImplementation) GetExcludedIPs() []string {
	return st.excludedIPs
}

// == VISITOR OPERATIONS =======================================================

// VisitorRegister creates a visitor from an HTTP request.
// If bot filtering is enabled, requests from known bots/crawlers, referrer spam
// domains, or data center IP ranges are silently skipped (returns nil).
func (st *storeImplementation) VisitorRegister(ctx context.Context, r *http.Request) error {
	path := r.URL.Path

	for _, prefix := range st.excludedPathPrefixes {
		if strings.HasPrefix(path, prefix) {
			if st.debugEnabled {
				st.logger.Info("path-filter: skipping excluded path", "path", path, "prefix", prefix)
			}
			return nil
		}
	}

	ip := req.GetIP(r)
	userAgent := r.UserAgent()
	referrer := r.Header.Get("Referer")

	for _, excludedIP := range st.excludedIPs {
		if ip == excludedIP {
			if st.debugEnabled {
				st.logger.Info("ip-filter: skipping excluded IP", "ip", ip)
			}
			return nil
		}
	}

	if st.botFilterEnabled {
		if IsBot(userAgent) {
			if st.debugEnabled {
				st.logger.Info("bot-filter: skipping bot visit", "user_agent", userAgent)
			}
			return nil
		}

		if IsReferrerSpam(referrer) {
			if st.debugEnabled {
				st.logger.Info("bot-filter: skipping referrer spam visit", "referrer", referrer)
			}
			return nil
		}

		if IsDataCenterIP(ip) {
			if st.debugEnabled {
				st.logger.Info("bot-filter: skipping data center IP visit", "ip", ip)
			}
			return nil
		}
	}

	uaInfo := ParseUserAgent(userAgent)

	visitor := NewVisitor().
		SetPath(path).
		SetIpAddress(ip).
		SetUserAgent(userAgent).
		SetUserBrowser(uaInfo.Browser).
		SetUserBrowserVersion(uaInfo.BrowserVersion).
		SetUserOs(uaInfo.Os).
		SetUserOsVersion(uaInfo.OsVersion).
		SetUserDevice(uaInfo.Device).
		SetUserDeviceType(uaInfo.DeviceType).
		SetUserReferrer(referrer)

	return st.VisitorCreate(ctx, visitor)
}

// VisitorCount counts visitors based on a query.
func (st *storeImplementation) VisitorCount(ctx context.Context, query VisitorQueryInterface) (int64, error) {
	if query.HasDistinct() && query.Distinct() != "" {
		q := st.buildQuery(query)
		var results []map[string]any
		err := q.Select("DISTINCT " + query.Distinct()).Get(&results)
		if err != nil {
			return 0, err
		}
		return int64(len(results)), nil
	}

	q := st.buildQuery(query)
	var count int64
	err := q.Count(&count)
	return count, err
}

// VisitorCreate creates a new visitor.
func (st *storeImplementation) VisitorCreate(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}

	if visitor.GetCreatedAt() == "" {
		visitor.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	}
	visitor.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_ID:                   visitor.GetID(),
		COLUMN_PATH:                 visitor.GetPath(),
		COLUMN_FINGERPRINT:          visitor.GetFingerprint(),
		COLUMN_IP_ADDRESS:           visitor.GetIpAddress(),
		COLUMN_COUNTRY:              visitor.GetCountry(),
		COLUMN_USER_ACCEPT_LANGUAGE: visitor.GetUserAcceptLanguage(),
		COLUMN_USER_ACCEPT_ENCODING: visitor.GetUserAcceptEncoding(),
		COLUMN_USER_AGENT:           visitor.GetUserAgent(),
		COLUMN_USER_OS:              visitor.GetUserOs(),
		COLUMN_USER_OS_VERSION:      visitor.GetUserOsVersion(),
		COLUMN_USER_DEVICE:          visitor.GetUserDevice(),
		COLUMN_USER_DEVICE_TYPE:     visitor.GetUserDeviceType(),
		COLUMN_USER_BROWSER:         visitor.GetUserBrowser(),
		COLUMN_USER_BROWSER_VERSION: visitor.GetUserBrowserVersion(),
		COLUMN_USER_REFERRER:        visitor.GetUserReferrer(),
		COLUMN_CREATED_AT:           visitor.GetCreatedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:           visitor.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:      visitor.GetSoftDeletedAtCarbon().StdTime(),
	}

	return st.db.Query().Table(st.visitorTableName).Create(row)
}

// VisitorDelete permanently deletes a visitor.
func (st *storeImplementation) VisitorDelete(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}
	return st.VisitorDeleteByID(ctx, visitor.GetID())
}

// VisitorDeleteByID permanently deletes a visitor by ID.
func (st *storeImplementation) VisitorDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("visitor id is empty")
	}

	_, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

// VisitorDeleteByIP permanently deletes all visitor records matching the given IP address.
// Returns the number of deleted rows.
func (st *storeImplementation) VisitorDeleteByIP(ctx context.Context, ip string) (int64, error) {
	if ip == "" {
		return 0, errors.New("visitor ip is empty")
	}

	rowsAffected, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_IP_ADDRESS+" = ?", ip).
		Delete()
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected, nil
}

// VisitorFindByID finds a visitor by ID.
func (st *storeImplementation) VisitorFindByID(ctx context.Context, id string) (VisitorInterface, error) {
	if id == "" {
		return nil, errors.New("visitor id is empty")
	}

	list, err := st.VisitorList(ctx, VisitorQuery().SetID(id).SetLimit(1))
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// VisitorList lists visitors based on a query.
func (st *storeImplementation) VisitorList(ctx context.Context, query VisitorQueryInterface) ([]VisitorInterface, error) {
	q := st.buildQuery(query)

	type visitorRow struct {
		ID                 string    `db:"id"`
		Path               string    `db:"path"`
		Fingerprint        string    `db:"fingerprint"`
		IPAddress          string    `db:"ip_address"`
		Country            string    `db:"country"`
		UserAcceptLanguage string    `db:"user_accept_language"`
		UserAcceptEncoding string    `db:"user_accept_encoding"`
		UserAgent          string    `db:"user_agent"`
		UserOs             string    `db:"user_os"`
		UserOsVersion      string    `db:"user_os_version"`
		UserDevice         string    `db:"user_device"`
		UserDeviceType     string    `db:"user_device_type"`
		UserBrowser        string    `db:"user_browser"`
		UserBrowserVersion string    `db:"user_browser_version"`
		UserReferrer       string    `db:"user_referrer"`
		CreatedAt          time.Time `db:"created_at"`
		UpdatedAt          time.Time `db:"updated_at"`
		SoftDeletedAt      time.Time `db:"soft_deleted_at"`
	}

	var rows []visitorRow
	if err := q.Get(&rows); err != nil {
		return []VisitorInterface{}, err
	}

	list := make([]VisitorInterface, 0, len(rows))
	for _, r := range rows {
		v := &visitorImplementation{}
		v.SetID(r.ID)
		v.SetPath(r.Path)
		v.SetFingerprint(r.Fingerprint)
		v.SetIpAddress(r.IPAddress)
		v.SetCountry(r.Country)
		v.SetUserAcceptLanguage(r.UserAcceptLanguage)
		v.SetUserAcceptEncoding(r.UserAcceptEncoding)
		v.SetUserAgent(r.UserAgent)
		v.SetUserOs(r.UserOs)
		v.SetUserOsVersion(r.UserOsVersion)
		v.SetUserDevice(r.UserDevice)
		v.SetUserDeviceType(r.UserDeviceType)
		v.SetUserBrowser(r.UserBrowser)
		v.SetUserBrowserVersion(r.UserBrowserVersion)
		v.SetUserReferrer(r.UserReferrer)
		v.CreatedAt.CreatedAt = r.CreatedAt
		v.UpdatedAt.UpdatedAt = r.UpdatedAt
		v.SoftDeletesMaxDate.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, v)
	}

	return list, nil
}

// VisitorSoftDelete soft deletes a visitor.
func (st *storeImplementation) VisitorSoftDelete(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}

	visitor.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: visitor.GetSoftDeletedAtCarbon().StdTime(),
		COLUMN_UPDATED_AT:      carbon.Now(carbon.UTC).StdTime(),
	}

	_, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_ID+" = ?", visitor.GetID()).
		Update(row)

	return err
}

// VisitorSoftDeleteByID soft deletes a visitor by ID.
func (st *storeImplementation) VisitorSoftDeleteByID(ctx context.Context, id string) error {
	visitor, err := st.VisitorFindByID(ctx, id)
	if err != nil {
		return err
	}
	if visitor == nil {
		return nil
	}
	return st.VisitorSoftDelete(ctx, visitor)
}

// VisitorUpdate updates a visitor.
func (st *storeImplementation) VisitorUpdate(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}

	visitor.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	row := map[string]any{
		COLUMN_PATH:                 visitor.GetPath(),
		COLUMN_FINGERPRINT:          visitor.GetFingerprint(),
		COLUMN_IP_ADDRESS:           visitor.GetIpAddress(),
		COLUMN_COUNTRY:              visitor.GetCountry(),
		COLUMN_USER_ACCEPT_LANGUAGE: visitor.GetUserAcceptLanguage(),
		COLUMN_USER_ACCEPT_ENCODING: visitor.GetUserAcceptEncoding(),
		COLUMN_USER_AGENT:           visitor.GetUserAgent(),
		COLUMN_USER_OS:              visitor.GetUserOs(),
		COLUMN_USER_OS_VERSION:      visitor.GetUserOsVersion(),
		COLUMN_USER_DEVICE:          visitor.GetUserDevice(),
		COLUMN_USER_DEVICE_TYPE:     visitor.GetUserDeviceType(),
		COLUMN_USER_BROWSER:         visitor.GetUserBrowser(),
		COLUMN_USER_BROWSER_VERSION: visitor.GetUserBrowserVersion(),
		COLUMN_USER_REFERRER:        visitor.GetUserReferrer(),
		COLUMN_UPDATED_AT:           visitor.GetUpdatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT:      visitor.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_ID+" = ?", visitor.GetID()).
		Update(row)

	return err
}

// == ENHANCE ==================================================================

// VisitorEnhance enriches visitor records that have an empty country field.
// For each record it:
//  1. Parses the user agent to fill in browser, OS, device, and device type
//     (if those fields are empty — e.g. when the record was created via
//     VisitorCreate instead of VisitorRegister)
//  2. Looks up the country via the configured GeoIPResolver
//
// Records are grouped by IP so that each unique IP is resolved only once.
// The country is then bulk-updated for ALL records sharing that IP (not just
// the current batch), significantly reducing DB writes when many visitors
// share the same IP.
//
// Records whose geo-IP lookup fails are still updated with UA data, but
// their country field is left empty so they get retried on the next call.
// Returns the number of records that were fully processed (country + UA).
func (st *storeImplementation) VisitorEnhance(ctx context.Context) (int, error) {
	if st.geoIPResolver == nil {
		return 0, errors.New("stats store: GeoIPResolver is not configured")
	}

	batchSize := st.enhanceBatchSize
	if batchSize <= 0 {
		batchSize = 10
	}

	visitors, err := st.VisitorList(ctx, VisitorQuery().
		SetCountry("empty").
		SetLimit(batchSize))
	if err != nil {
		return 0, err
	}

	if len(visitors) == 0 {
		return 0, nil
	}

	// Collect unique IPs to resolve each only once
	seenIPs := make(map[string]bool)
	ipOrder := []string{}
	for _, v := range visitors {
		ip := v.GetIpAddress()
		if !seenIPs[ip] {
			seenIPs[ip] = true
			ipOrder = append(ipOrder, ip)
		}
	}

	// Resolve each unique IP once, then bulk-update country for ALL records
	// with that IP (not just the current batch) to minimize DB writes.
	// A 2-second delay is added between API calls to avoid overwhelming the
	// geo-IP service (e.g. ip2c.org rate limits).
	resolvedCountries := make(map[string]string)
	for i, ip := range ipOrder {
		if i > 0 {
			select {
			case <-time.After(2 * time.Second):
			case <-ctx.Done():
				return 0, ctx.Err()
			}
		}

		country, err := st.geoIPResolver.Resolve(ctx, ip)
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("VisitorEnhance: geo-IP lookup failed",
					"ip", ip, "error", err)
			}
			continue
		}

		_, err = st.db.Query().
			Table(st.visitorTableName).
			Where(COLUMN_IP_ADDRESS+" = ?", ip).
			Where(COLUMN_COUNTRY+" = ?", "").
			Update(map[string]any{
				COLUMN_COUNTRY:    country,
				COLUMN_UPDATED_AT: carbon.Now(carbon.UTC).StdTime(),
			})
		if err != nil {
			if st.debugEnabled {
				st.logger.Error("VisitorEnhance: bulk country update failed",
					"ip", ip, "error", err)
			}
			continue
		}

		resolvedCountries[ip] = country
	}

	// Update UA fields per record (UA differs per visitor even for the same IP)
	processed := 0
	for _, visitor := range visitors {
		uaInfo := ParseUserAgent(visitor.GetUserAgent())
		if visitor.GetUserBrowser() == "" {
			visitor.SetUserBrowser(uaInfo.Browser)
		}
		if visitor.GetUserBrowserVersion() == "" {
			visitor.SetUserBrowserVersion(uaInfo.BrowserVersion)
		}
		if visitor.GetUserOs() == "" {
			visitor.SetUserOs(uaInfo.Os)
		}
		if visitor.GetUserOsVersion() == "" {
			visitor.SetUserOsVersion(uaInfo.OsVersion)
		}
		if visitor.GetUserDevice() == "" {
			visitor.SetUserDevice(uaInfo.Device)
		}
		if visitor.GetUserDeviceType() == "" {
			visitor.SetUserDeviceType(uaInfo.DeviceType)
		}

		// Set country from the resolved map so VisitorUpdate persists it
		ipResolved := false
		if country, ok := resolvedCountries[visitor.GetIpAddress()]; ok {
			visitor.SetCountry(country)
			ipResolved = true
		}

		if err := st.VisitorUpdate(ctx, visitor); err != nil {
			if st.debugEnabled {
				st.logger.Error("VisitorEnhance: update failed",
					"id", visitor.GetID(), "error", err)
			}
			continue
		}

		if ipResolved {
			processed++
		}
	}

	return processed, nil
}

// == QUERY BUILDER ============================================================

func (st *storeImplementation) buildQuery(query VisitorQueryInterface) contractsorm.Query {
	// Use Model() to enable neat's automatic soft delete handling via SoftDeletesMaxDate
	q := st.db.Query().Model(&visitorImplementation{}).Table(st.visitorTableName)

	if query.HasID() && query.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.ID())
	}

	if query.HasIDIn() && len(query.IDIn()) > 0 {
		args := make([]any, len(query.IDIn()))
		for i, id := range query.IDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}

	if query.HasIPIn() && len(query.IPIn()) > 0 {
		args := make([]any, len(query.IPIn()))
		for i, ip := range query.IPIn() {
			args[i] = ip
		}
		q = q.WhereIn(COLUMN_IP_ADDRESS, args)
	}

	if query.HasIPNotIn() && len(query.IPNotIn()) > 0 {
		args := make([]any, len(query.IPNotIn()))
		for i, ip := range query.IPNotIn() {
			args[i] = ip
		}
		q = q.WhereNotIn(COLUMN_IP_ADDRESS, args)
	}

	if query.HasCountry() && query.Country() != "" {
		if strings.EqualFold(query.Country(), "empty") {
			q = q.Where(COLUMN_COUNTRY+" = ?", "")
		} else {
			q = q.Where(COLUMN_COUNTRY+" = ?", query.Country())
		}
	}

	if query.HasPathExact() && query.PathExact() != "" {
		q = q.Where(COLUMN_PATH+" = ?", query.PathExact())
	} else if query.HasPathContains() && query.PathContains() != "" {
		q = q.Where(COLUMN_PATH+" LIKE ?", "%"+query.PathContains()+"%")
	}

	if query.HasDeviceType() && query.DeviceType() != "" {
		if strings.EqualFold(query.DeviceType(), "empty") {
			q = q.Where(COLUMN_USER_DEVICE_TYPE+" = ?", "")
		} else {
			q = q.Where(COLUMN_USER_DEVICE_TYPE+" = ?", query.DeviceType())
		}
	}

	if query.HasCreatedAtGte() && query.CreatedAtGte() != "" {
		if createdAt, ok := parseCreatedAt(query.CreatedAtGte()); ok {
			q = q.Where(COLUMN_CREATED_AT+" >= ?", createdAt)
		} else {
			return q.Where("1 = 0")
		}
	}
	if query.HasCreatedAtLte() && query.CreatedAtLte() != "" {
		if createdAt, ok := parseCreatedAt(query.CreatedAtLte()); ok {
			q = q.Where(COLUMN_CREATED_AT+" <= ?", createdAt)
		} else {
			return q.Where("1 = 0")
		}
	}

	if query.HasLimit() && query.Limit() > 0 {
		q = q.Limit(query.Limit())
	}

	if query.HasOffset() && query.Offset() > 0 {
		q = q.Offset(query.Offset())
	}

	if query.HasOrderBy() && query.OrderBy() != "" {
		sortOrder := "desc"
		if query.HasSortOrder() && query.SortOrder() != "" {
			sortOrder = query.SortOrder()
		}
		q = q.OrderBy(query.OrderBy(), sortOrder)
	}

	// Handle soft delete filtering via neat's automatic handling (SoftDeletesMaxDate)
	if query.HasSoftDeletedIncluded() && query.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	}

	return q
}

// parseCreatedAt parses a created_at bound string into a time.Time value. It
// accepts RFC3339 and common SQL datetime formats so callers can pass strings
// to the visitor query interface while the ORM receives a typed DateTime
// comparison.
func parseCreatedAt(value string) (time.Time, bool) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, value); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
