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
	visitorTableName   string
	db                 *neat.Database
	automigrateEnabled bool
	debugEnabled       bool
	logger             *slog.Logger
}

// == MIGRATE ==================================================================

// MigrateUp creates the visitor table if it does not already exist.
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.visitorTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: table already exists", "table", st.visitorTableName)
		}
		return nil
	}

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
	})

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateUp failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the visitor table.
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
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

// == VISITOR OPERATIONS =======================================================

// VisitorRegister creates a visitor from an HTTP request.
func (st *storeImplementation) VisitorRegister(ctx context.Context, r *http.Request) error {
	path := r.URL.Path
	ip := req.GetIP(r)
	userAgent := r.UserAgent()

	visitor := NewVisitor().
		SetPath(path).
		SetIpAddress(ip).
		SetUserAgent(userAgent)

	return st.VisitorCreate(ctx, visitor)
}

// VisitorCount counts visitors based on a query.
func (st *storeImplementation) VisitorCount(ctx context.Context, query VisitorQueryInterface) (int64, error) {
	if query.HasDistinct() && query.Distinct() != "" {
		q := st.buildQuery(query)
		var results []map[string]any
		err := q.Table(st.visitorTableName).Select("DISTINCT " + query.Distinct()).Get(&results)
		if err != nil {
			return 0, err
		}
		return int64(len(results)), nil
	}

	q := st.buildQuery(query)
	var count int64
	err := q.Table(st.visitorTableName).Count(&count)
	return count, err
}

// VisitorCreate creates a new visitor.
func (st *storeImplementation) VisitorCreate(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}

	visitor.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
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
	if err := q.Table(st.visitorTableName).Get(&rows); err != nil {
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
		v.CreatedAtField.CreatedAt = r.CreatedAt
		v.UpdatedAtField.UpdatedAt = r.UpdatedAt
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

// == QUERY BUILDER ============================================================

func (st *storeImplementation) buildQuery(query VisitorQueryInterface) contractsorm.Query {
	q := st.db.Query()

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
		q = q.Where(COLUMN_CREATED_AT+" >= ?", query.CreatedAtGte())
	}
	if query.HasCreatedAtLte() && query.CreatedAtLte() != "" {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", query.CreatedAtLte())
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

	if query.HasSoftDeletedIncluded() && query.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())
	}

	return q
}
