package statsstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"net/netip"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dracory/req"
	"github.com/dracory/str"
	"github.com/dromara/carbon/v2"
)

// == INTERFACE ================================================================

var _ StoreInterface = (*storeImplementation)(nil)

// storeImplementation implements StoreInterface for visitor operations.
type storeImplementation struct {
	visitorTableName   string
	settingsTableName  string
	db                 *neat.Database
	automigrateEnabled bool
	debugEnabled       bool
	botFilterEnabled   bool
	botAutoTagEnabled  bool
	// filtersMu guards the in-memory filter lists below; they are read on
	// every VisitorRegister/buildQuery call while admin endpoints may mutate
	// them concurrently.
	filtersMu            sync.RWMutex
	excludedPathPrefixes []string
	includedPathPrefixes []string
	excludedIPs          []string
	// The *Opts fields preserve the values supplied via NewStoreOptions so
	// that DB writes never persist option-configured entries and option
	// entries remain active even if "removed" via the admin UI.
	excludedPathPrefixesOpts []string
	includedPathPrefixesOpts []string
	excludedIPsOpts          []string
	geoIPResolver            GeoIPResolver
	enhanceBatchSize         int
	trustedProxies           []string
	fingerprintSalt          string
	logger                   *slog.Logger
}

// == MIGRATE ==================================================================

// MigrateUp creates the visitor table and excluded IPs table if they do not already exist.
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.visitorTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: table already exists", "table", st.visitorTableName)
		}

		// Add bot/threat columns to existing tables that predate them.
		// neat's HasColumn guards against re-adding.
		if !st.db.Schema().HasColumn(st.visitorTableName, COLUMN_BOT) {
			if err := st.db.Schema().Table(st.visitorTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_BOT, MAX_LEN_BOT).Default(VALUE_NO)
			}); err != nil {
				if st.debugEnabled {
					st.logger.Error("MigrateUp: add bot column failed", "error", err)
				}
				return err
			}
		}
		if !st.db.Schema().HasColumn(st.visitorTableName, COLUMN_THREAT) {
			if err := st.db.Schema().Table(st.visitorTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_THREAT, MAX_LEN_THREAT).Default(VALUE_NO)
			}); err != nil {
				if st.debugEnabled {
					st.logger.Error("MigrateUp: add threat column failed", "error", err)
				}
				return err
			}
		}

		// Add indexes on bot/threat columns for existing tables that predate them.
		// neat generates index names as "{table}_{column}_index".
		botIndexName := strings.ToLower(st.visitorTableName + "_" + COLUMN_BOT + "_index")
		if !st.db.Schema().HasIndex(st.visitorTableName, botIndexName) {
			if err := st.db.Schema().Table(st.visitorTableName, func(table contractsschema.Blueprint) {
				table.Index(COLUMN_BOT)
			}); err != nil {
				if st.debugEnabled {
					st.logger.Error("MigrateUp: add bot index failed", "error", err)
				}
				return err
			}
		}
		threatIndexName := strings.ToLower(st.visitorTableName + "_" + COLUMN_THREAT + "_index")
		if !st.db.Schema().HasIndex(st.visitorTableName, threatIndexName) {
			if err := st.db.Schema().Table(st.visitorTableName, func(table contractsschema.Blueprint) {
				table.Index(COLUMN_THREAT)
			}); err != nil {
				if st.debugEnabled {
					st.logger.Error("MigrateUp: add threat index failed", "error", err)
				}
				return err
			}
		}
	} else {
		err := st.db.Schema().Create(st.visitorTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, MAX_LEN_ID)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_PATH, MAX_LEN_PATH)
			table.String(COLUMN_FINGERPRINT, MAX_LEN_FINGERPRINT)
			table.String(COLUMN_IP_ADDRESS, MAX_LEN_IP_ADDRESS)
			table.String(COLUMN_COUNTRY, MAX_LEN_COUNTRY)
			table.String(COLUMN_USER_ACCEPT_LANGUAGE, MAX_LEN_USER_ACCEPT_LANGUAGE)
			table.String(COLUMN_USER_ACCEPT_ENCODING, MAX_LEN_USER_ACCEPT_ENCODING)
			table.String(COLUMN_USER_AGENT, MAX_LEN_USER_AGENT)
			table.String(COLUMN_USER_OS, MAX_LEN_USER_OS)
			table.String(COLUMN_USER_OS_VERSION, MAX_LEN_USER_OS_VERSION)
			table.String(COLUMN_USER_DEVICE, MAX_LEN_USER_DEVICE)
			table.String(COLUMN_USER_DEVICE_TYPE, MAX_LEN_USER_DEVICE_TYPE)
			table.String(COLUMN_USER_BROWSER, MAX_LEN_USER_BROWSER)
			table.String(COLUMN_USER_BROWSER_VERSION, MAX_LEN_USER_BROWSER_VERSION)
			table.String(COLUMN_USER_REFERRER, MAX_LEN_USER_REFERRER)
			table.String(COLUMN_BOT, MAX_LEN_BOT).Default(VALUE_NO)
			table.String(COLUMN_THREAT, MAX_LEN_THREAT).Default(VALUE_NO)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)

			table.Index(COLUMN_CREATED_AT)
			table.Index(COLUMN_IP_ADDRESS)
			table.Index(COLUMN_FINGERPRINT)
			table.Index(COLUMN_BOT)
			table.Index(COLUMN_THREAT)
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
			table.LongText(COLUMN_VALUE)
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

// SetBotFilterEnabled enables or disables skipping bot/threat traffic at
// ingestion. When enabled, requests matching bot/referrer-spam/data-center/
// bot-path/malicious-path heuristics are silently skipped (not inserted).
// This is independent of BotAutoTagEnabled, which controls flag computation.
func (st *storeImplementation) SetBotFilterEnabled(enabled bool) {
	st.botFilterEnabled = enabled
}

// IsBotFilterEnabled returns whether bot filtering is currently enabled.
func (st *storeImplementation) IsBotFilterEnabled() bool {
	return st.botFilterEnabled
}

// SetBotAutoTagEnabled controls whether bot/threat traffic is tagged (bot='yes'/
// threat='yes') and kept, instead of being dropped at ingestion. When true,
// VisitorRegister tags matched rows instead of dropping them, and
// VisitorCreate/VisitorUpdate auto-compute the flags from the visitor's own
// fields. Defaults to false (drop behaviour).
func (st *storeImplementation) SetBotAutoTagEnabled(enabled bool) {
	st.botAutoTagEnabled = enabled
}

// IsBotAutoTagEnabled returns whether bot auto-tagging is currently enabled.
func (st *storeImplementation) IsBotAutoTagEnabled() bool {
	return st.botAutoTagEnabled
}

// SetExcludedPathPrefixes sets path patterns that should be excluded from
// visitor tracking. Entries without a "*" act as prefix matches ("/admin/"
// covers "/admin/users"); "*" matches any character sequence including "/"
// (e.g. "/api/*", "/*.php"). Matching requests are silently skipped by
// VisitorRegister.
func (st *storeImplementation) SetExcludedPathPrefixes(prefixes []string) {
	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()
	st.excludedPathPrefixes = prefixes
	st.excludedPathPrefixesOpts = slices.Clone(prefixes)
}

// GetExcludedPathPrefixes returns the currently configured excluded path prefixes.
func (st *storeImplementation) GetExcludedPathPrefixes() []string {
	st.filtersMu.RLock()
	defer st.filtersMu.RUnlock()
	return slices.Clone(st.excludedPathPrefixes)
}

// SetIncludedPathPrefixes sets path patterns that should be included in visitor
// tracking. Entries use the same pattern syntax as SetExcludedPathPrefixes
// (prefix match, or glob-style "*"). When the list is non-empty it acts as a
// whitelist: requests whose URL path does not match are silently skipped by
// VisitorRegister. An empty list means all paths are tracked.
func (st *storeImplementation) SetIncludedPathPrefixes(prefixes []string) {
	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()
	st.includedPathPrefixes = prefixes
	st.includedPathPrefixesOpts = slices.Clone(prefixes)
}

// GetIncludedPathPrefixes returns the currently configured included path prefixes.
func (st *storeImplementation) GetIncludedPathPrefixes() []string {
	st.filtersMu.RLock()
	defer st.filtersMu.RUnlock()
	return slices.Clone(st.includedPathPrefixes)
}

// SetExcludedIPs sets IP addresses that should be excluded from visitor tracking.
// Requests from these IPs will be silently skipped by VisitorRegister.
func (st *storeImplementation) SetExcludedIPs(ips []string) {
	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()
	st.excludedIPs = ips
	st.excludedIPsOpts = slices.Clone(ips)
}

// GetExcludedIPs returns the currently configured excluded IP addresses.
func (st *storeImplementation) GetExcludedIPs() []string {
	st.filtersMu.RLock()
	defer st.filtersMu.RUnlock()
	return slices.Clone(st.excludedIPs)
}

// == VISITOR OPERATIONS =======================================================

// VisitorRegister creates a visitor from an HTTP request.
// If BotFilterEnabled is true, bot/threat traffic is skipped (not inserted).
// If BotAutoTagEnabled is true, bot/threat flags are computed and set on the
// row. The two flags are independent.
func (st *storeImplementation) VisitorRegister(ctx context.Context, r *http.Request) error {
	path := r.URL.Path

	st.filtersMu.RLock()
	excludedPaths := st.excludedPathPrefixes
	includedPaths := st.includedPathPrefixes
	excludedIPs := st.excludedIPs
	st.filtersMu.RUnlock()

	if pattern, matched := pathPatternMatchAny(excludedPaths, path); matched {
		if st.debugEnabled {
			st.logger.Info("path-filter: skipping excluded path", "path", path, "pattern", pattern)
		}
		return nil
	}

	// Whitelist semantics: when includedPathPrefixes is non-empty, only paths
	// matching a prefix are tracked; everything else is skipped.
	if len(includedPaths) > 0 {
		if _, matched := pathPatternMatchAny(includedPaths, path); !matched {
			if st.debugEnabled {
				st.logger.Info("path-filter: skipping non-included path", "path", path)
			}
			return nil
		}
	}

	ip := st.clientIP(r)
	if ip == "" {
		if st.debugEnabled {
			st.logger.Info("ip-filter: skipping unparseable client IP", "remote_addr", r.RemoteAddr)
		}
		return nil
	}
	userAgent := r.UserAgent()
	referrer := r.Header.Get("Referer")

	if slices.Contains(excludedIPs, ip) {
		if st.debugEnabled {
			st.logger.Info("ip-filter: skipping excluded IP", "ip", redactIP(ip))
		}
		return nil
	}

	uaInfo := ParseUserAgent(userAgent)

	// BotFilterEnabled: detect and skip bot/threat traffic at ingestion.
	if st.botFilterEnabled {
		isBot := IsBot(userAgent) || IsReferrerSpam(referrer) || IsDataCenterIP(ip) || IsBotPath(path) ||
			strings.EqualFold(uaInfo.DeviceType, "bot") || strings.EqualFold(uaInfo.Device, "bot")
		isThreat := IsMaliciousPath(path)

		if isBot || isThreat {
			if st.debugEnabled {
				st.logger.Info("bot-filter: skipping bot/threat visit",
					"user_agent", userAgent, "ip", redactIP(ip), "path", path)
			}
			return nil
		}
	}

	// BotAutoTagEnabled: compute and set bot/threat flags on the row.
	botVal := VALUE_NO
	threatVal := VALUE_NO

	if st.botAutoTagEnabled {
		isBot := IsBot(userAgent) || IsReferrerSpam(referrer) || IsDataCenterIP(ip) || IsBotPath(path) ||
			strings.EqualFold(uaInfo.DeviceType, "bot") || strings.EqualFold(uaInfo.Device, "bot")
		isThreat := IsMaliciousPath(path)

		if isBot {
			botVal = VALUE_YES
		}
		if isThreat {
			threatVal = VALUE_YES
		}

		if st.debugEnabled && (isBot || isThreat) {
			st.logger.Info("bot-tag: tagging visit",
				"bot", botVal, "threat", threatVal,
				"user_agent", userAgent, "ip", redactIP(ip), "path", path)
		}
	}

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
		SetUserReferrer(referrer).
		SetBot(botVal).
		SetThreat(threatVal)

	visitor.SetFingerprint(st.visitorFingerprint(visitor))

	return st.VisitorCreate(ctx, visitor)
}

// clientIP resolves the client IP for a request. When TrustedProxies is
// configured, forwarding headers are only honoured from trusted peers;
// otherwise req.GetIP's default precedence (X-Real-IP, X-Forwarded-For,
// RemoteAddr) applies. The result is validated and canonicalized; ""
// means no parseable IP could be determined.
func (st *storeImplementation) clientIP(r *http.Request) string {
	var ip string
	if len(st.trustedProxies) > 0 {
		ip = req.GetIPWithOptions(r, req.IPOptions{
			TrustedProxies:            st.trustedProxies,
			Validate:                  true,
			ReturnPrivateIfAllPrivate: true,
		})
	} else {
		ip = req.GetIP(r)
	}
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return ""
	}
	return addr.String()
}

// visitorFingerprint computes the stored fingerprint for a visitor. With
// FingerprintSalt configured it is SHA-256(salt + ip + ua); otherwise
// plain SHA-256(ip + ua), truncated to MAX_LEN_FINGERPRINT hex chars.
func (st *storeImplementation) visitorFingerprint(v VisitorInterface) string {
	if st.fingerprintSalt == "" {
		return v.FingerprintCalculate()
	}
	sum := str.SHA256(st.fingerprintSalt + v.GetIpAddress() + v.GetUserAgent())
	if len(sum) > MAX_LEN_FINGERPRINT {
		sum = sum[:MAX_LEN_FINGERPRINT]
	}
	return sum
}

// redactIP masks the last octet (IPv4) or last 16 bits (IPv6) of an
// address for debug logging so full IPs are not written to stdout.
func redactIP(ip string) string {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return "<invalid>"
	}
	if addr.Is4() {
		a4 := addr.As4()
		return netip.AddrFrom4([4]byte{a4[0], a4[1], a4[2], 0}).String()
	}
	a16 := addr.As16()
	a16[14], a16[15] = 0, 0
	return netip.AddrFrom16(a16).String()
}

// VisitorCount counts visitors based on a query.
func (st *storeImplementation) VisitorCount(ctx context.Context, query VisitorQueryInterface) (int64, error) {
	if query.HasDistinct() && query.Distinct() != "" {
		if !visitorQueryColumns[query.Distinct()] {
			return 0, errors.New("visitor query: distinct is not an allowed column")
		}
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

// ensureBotThreatFlags computes and sets the bot/threat flags on the visitor
// when they have not been explicitly set (empty string) and botAutoTagEnabled
// is true. The flags are derived from the visitor's own user-agent, IP,
// referrer, and path fields using the same heuristics as VisitorRegister.
// When botAutoTagEnabled is false, flags are left as-is (empty or whatever
// the caller set).
func (st *storeImplementation) ensureBotThreatFlags(visitor VisitorInterface) {
	if !st.botAutoTagEnabled {
		return
	}

	isBot := IsBot(visitor.GetUserAgent()) ||
		strings.EqualFold(visitor.GetUserDeviceType(), "bot") ||
		strings.EqualFold(visitor.GetUserDevice(), "bot") ||
		IsReferrerSpam(visitor.GetUserReferrer()) ||
		IsDataCenterIP(visitor.GetIpAddress()) ||
		IsBotPath(visitor.GetPath())

	if isBot {
		visitor.SetBot(VALUE_YES)
	} else if visitor.GetBot() == "" {
		visitor.SetBot(VALUE_NO)
	}

	if visitor.GetThreat() == "" {
		if IsMaliciousPath(visitor.GetPath()) {
			visitor.SetThreat(VALUE_YES)
		} else {
			visitor.SetThreat(VALUE_NO)
		}
	}
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

	st.ensureBotThreatFlags(visitor)
	if visitor.GetBot() == "" {
		visitor.SetBot(VALUE_NO)
	}
	if visitor.GetThreat() == "" {
		visitor.SetThreat(VALUE_NO)
	}

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
		COLUMN_BOT:                  visitor.GetBot(),
		COLUMN_THREAT:               visitor.GetThreat(),
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

// VisitorDeleteByIDs permanently deletes all visitor records with the given
// IDs in a single DELETE. Returns the number of deleted rows.
func (st *storeImplementation) VisitorDeleteByIDs(ctx context.Context, ids []string) (int64, error) {
	if len(ids) == 0 {
		return 0, errors.New("visitor ids are empty")
	}

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rowsAffected, err := st.db.Query().
		Table(st.visitorTableName).
		WhereIn(COLUMN_ID, args).
		Delete()
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected, nil
}

// VisitorSetBotByIP sets the bot flag on all visitor records for the given
// IP in a single UPDATE. Returns the number of updated rows.
func (st *storeImplementation) VisitorSetBotByIP(ctx context.Context, ip, bot string) (int64, error) {
	if ip == "" {
		return 0, errors.New("visitor ip is empty")
	}
	if bot != VALUE_YES && bot != VALUE_NO {
		return 0, errors.New("invalid bot flag value")
	}

	rowsAffected, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_IP_ADDRESS+" = ?", ip).
		Update(map[string]any{
			COLUMN_BOT:        bot,
			COLUMN_UPDATED_AT: carbon.Now(carbon.UTC).StdTime(),
		})
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected, nil
}

// VisitorSetThreatByIP sets the threat flag on all visitor records for the
// given IP in a single UPDATE. Returns the number of updated rows.
func (st *storeImplementation) VisitorSetThreatByIP(ctx context.Context, ip, threat string) (int64, error) {
	if ip == "" {
		return 0, errors.New("visitor ip is empty")
	}
	if threat != VALUE_YES && threat != VALUE_NO {
		return 0, errors.New("invalid threat flag value")
	}

	rowsAffected, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_IP_ADDRESS+" = ?", ip).
		Update(map[string]any{
			COLUMN_THREAT:     threat,
			COLUMN_UPDATED_AT: carbon.Now(carbon.UTC).StdTime(),
		})
	if err != nil {
		return 0, err
	}

	return rowsAffected.RowsAffected, nil
}

// VisitorDeleteOlderThan permanently deletes all visitor records created before the given timestamp.
// Returns the number of deleted rows.
func (st *storeImplementation) VisitorDeleteOlderThan(ctx context.Context, timestamp string) (int64, error) {
	if timestamp == "" {
		return 0, errors.New("timestamp is empty")
	}

	ts, ok := parseCreatedAt(timestamp)
	if !ok {
		return 0, errors.New("invalid timestamp format")
	}

	rowsAffected, err := st.db.Query().
		Table(st.visitorTableName).
		Where(COLUMN_CREATED_AT+" < ?", ts).
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
		Bot                string    `db:"bot"`
		Threat             string    `db:"threat"`
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
		v.SetBot(r.Bot)
		v.SetThreat(r.Threat)
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

	st.ensureBotThreatFlags(visitor)
	if visitor.GetBot() == "" {
		visitor.SetBot(VALUE_NO)
	}
	if visitor.GetThreat() == "" {
		visitor.SetThreat(VALUE_NO)
	}

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
		COLUMN_BOT:                  visitor.GetBot(),
		COLUMN_THREAT:               visitor.GetThreat(),
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

		if strings.EqualFold(visitor.GetUserDeviceType(), "bot") || strings.EqualFold(visitor.GetUserDevice(), "bot") || IsBot(visitor.GetUserAgent()) {
			visitor.SetBot(VALUE_YES)
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

// visitorQueryColumns is the whitelist of columns accepted by the raw-string
// Distinct and OrderBy query options. Values not on this list are rejected or
// ignored so callers may forward request-derived input without risking SQL
// injection through the query API.
var visitorQueryColumns = map[string]bool{
	COLUMN_ID:                   true,
	COLUMN_PATH:                 true,
	COLUMN_FINGERPRINT:          true,
	COLUMN_IP_ADDRESS:           true,
	COLUMN_COUNTRY:              true,
	COLUMN_USER_ACCEPT_LANGUAGE: true,
	COLUMN_USER_ACCEPT_ENCODING: true,
	COLUMN_USER_AGENT:           true,
	COLUMN_USER_OS:              true,
	COLUMN_USER_OS_VERSION:      true,
	COLUMN_USER_DEVICE:          true,
	COLUMN_USER_DEVICE_TYPE:     true,
	COLUMN_USER_BROWSER:         true,
	COLUMN_USER_BROWSER_VERSION: true,
	COLUMN_USER_REFERRER:        true,
	COLUMN_BOT:                  true,
	COLUMN_THREAT:               true,
	COLUMN_CREATED_AT:           true,
	COLUMN_UPDATED_AT:           true,
	COLUMN_SOFT_DELETED_AT:      true,
}

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

	// Combine query.IPNotIn and st.excludedIPs (unless query.HasIPIn is set)
	var excluded []string
	if query.HasIPNotIn() && len(query.IPNotIn()) > 0 {
		excluded = append(excluded, query.IPNotIn()...)
	}
	if !query.HasIPIn() {
		for _, ip := range st.GetExcludedIPs() {
			if !slices.Contains(excluded, ip) {
				excluded = append(excluded, ip)
			}
		}
	}
	if len(excluded) > 0 {
		args := make([]any, len(excluded))
		for i, ip := range excluded {
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
		q = q.Where("LOWER("+COLUMN_PATH+") LIKE LOWER(?)", "%"+query.PathContains()+"%")
	}

	if query.HasDeviceType() && query.DeviceType() != "" {
		if strings.EqualFold(query.DeviceType(), "empty") {
			q = q.Where(COLUMN_USER_DEVICE_TYPE+" = ?", "")
		} else {
			q = q.Where(COLUMN_USER_DEVICE_TYPE+" = ?", query.DeviceType())
		}
	}

	if query.HasBot() && query.Bot() != "" {
		if query.Bot() == VALUE_NO {
			q = q.Where("("+COLUMN_BOT+" = ? OR "+COLUMN_BOT+" = ? OR "+COLUMN_BOT+" IS NULL)", VALUE_NO, "").
				Where("("+COLUMN_USER_DEVICE_TYPE+" != ? OR "+COLUMN_USER_DEVICE_TYPE+" IS NULL)", "bot").
				Where("("+COLUMN_USER_DEVICE+" != ? OR "+COLUMN_USER_DEVICE+" IS NULL)", "bot")
		} else {
			q = q.Where("("+COLUMN_BOT+" = ? OR "+COLUMN_USER_DEVICE_TYPE+" = ? OR "+COLUMN_USER_DEVICE+" = ?)", VALUE_YES, "bot", "bot")
		}
	}

	if query.HasThreat() && query.Threat() != "" {
		q = q.Where(COLUMN_THREAT+" = ?", query.Threat())
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

	if query.HasOrderBy() && visitorQueryColumns[query.OrderBy()] {
		sortOrder := "desc"
		if strings.EqualFold(query.SortOrder(), "asc") {
			sortOrder = "asc"
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
