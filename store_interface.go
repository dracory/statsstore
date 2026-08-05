package statsstore

import (
	"context"
	"database/sql"
	"net/http"
)

// StoreInterface defines the interface for a stats store.
type StoreInterface interface {
	// MigrateDown drops the stats store tables
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error

	// MigrateUp creates the stats store tables
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	EnableDebug(debug bool)
	GetDB() *sql.DB

	SetBotFilterEnabled(enabled bool)
	IsBotFilterEnabled() bool
	SetBotAutoTagEnabled(enabled bool)
	IsBotAutoTagEnabled() bool

	SetExcludedPathPrefixes(prefixes []string)
	GetExcludedPathPrefixes() []string

	SetExcludedIPs(ips []string)
	GetExcludedIPs() []string

	ExcludedIPList(ctx context.Context) ([]string, error)
	ExcludedIPAdd(ctx context.Context, ip string) error
	ExcludedIPRemove(ctx context.Context, ip string) error

	// SettingGet retrieves a setting value by key. Returns empty string and
	// nil error if the key does not exist.
	SettingGet(ctx context.Context, key string) (string, error)
	// SettingSet stores a setting value by key, using upsert semantics.
	SettingSet(ctx context.Context, key, value string) error
	// SettingDelete removes a setting by key. No error if the key is absent.
	SettingDelete(ctx context.Context, key string) error
	// SettingHas reports whether a setting key exists.
	SettingHas(ctx context.Context, key string) (bool, error)
	// SettingList returns all settings as a map of key to value.
	// Returns an empty map if the settings table does not exist or is empty.
	SettingList(ctx context.Context) (map[string]string, error)

	VisitorCount(ctx context.Context, query VisitorQueryInterface) (int64, error)
	VisitorCreate(ctx context.Context, user VisitorInterface) error
	VisitorDelete(ctx context.Context, user VisitorInterface) error
	VisitorDeleteByID(ctx context.Context, id string) error
	VisitorDeleteByIP(ctx context.Context, ip string) (int64, error)
	VisitorFindByID(ctx context.Context, userID string) (VisitorInterface, error)
	VisitorList(ctx context.Context, query VisitorQueryInterface) ([]VisitorInterface, error)
	VisitorRegister(ctx context.Context, r *http.Request) error
	VisitorSoftDelete(ctx context.Context, user VisitorInterface) error
	VisitorSoftDeleteByID(ctx context.Context, id string) error
	VisitorUpdate(ctx context.Context, user VisitorInterface) error

	// VisitorEnhance enriches visitor records that have an empty country field.
	// It parses the user agent (browser, OS, device type) and looks up the
	// country via the configured GeoIPResolver. UA fields are updated even if
	// the geo-IP lookup fails; the country stays empty for retry.
	// Returns the number of records fully processed (country + UA).
	// Call this from a background task/cron on whatever schedule suits your
	// traffic (e.g. every 5 minutes).
	//
	// If no GeoIPResolver was configured, it returns 0 and an error.
	VisitorEnhance(ctx context.Context) (int, error)
}
