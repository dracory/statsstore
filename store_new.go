package statsstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/dracory/neat"
)

// NewStoreOptions defines the options for creating a new stats store.
type NewStoreOptions struct {
	VisitorTableName     string
	SettingsTableName    string
	DB                   *sql.DB
	AutomigrateEnabled   bool
	DebugEnabled         bool
	BotFilterEnabled     bool // when true, skip bot/threat traffic at ingestion (don't insert). Independent of BotAutoTagEnabled.
	BotAutoTagEnabled    bool // when true, compute and set bot/threat flags on inserted rows. Also auto-computes flags on VisitorCreate/VisitorUpdate.
	ExcludedPathPrefixes []string
	IncludedPathPrefixes []string // when non-empty, acts as a whitelist: only matching paths are tracked
	ExcludedIPs          []string
	GeoIPResolver        GeoIPResolver // optional; enables VisitorEnhance for batch country enrichment
	EnhanceBatchSize     int           // number of records per VisitorEnhance call; default 10
}

// NewStore creates a new stats store.
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.VisitorTableName == "" {
		return nil, errors.New("stats store: VisitorTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("stats store: DB is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	settingsTable := opts.SettingsTableName
	if settingsTable == "" {
		settingsTable = DEFAULT_SETTINGS_TABLE
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		visitorTableName:         opts.VisitorTableName,
		settingsTableName:        settingsTable,
		db:                       neatDB,
		automigrateEnabled:       opts.AutomigrateEnabled,
		debugEnabled:             opts.DebugEnabled,
		botFilterEnabled:         opts.BotFilterEnabled,
		botAutoTagEnabled:        opts.BotAutoTagEnabled,
		excludedPathPrefixes:     opts.ExcludedPathPrefixes,
		includedPathPrefixes:     opts.IncludedPathPrefixes,
		excludedIPs:              opts.ExcludedIPs,
		excludedPathPrefixesOpts: opts.ExcludedPathPrefixes,
		includedPathPrefixesOpts: opts.IncludedPathPrefixes,
		excludedIPsOpts:          opts.ExcludedIPs,
		geoIPResolver:            opts.GeoIPResolver,
		enhanceBatchSize:         opts.EnhanceBatchSize,
		logger:                   logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	// Load excluded IPs from the database and merge with any provided via options
	store.excludedIPs = mergeWithDBList(store.excludedIPsLoadFromDB, opts.ExcludedIPs)

	// Load excluded/included path prefixes from the database and merge with
	// any provided via options
	store.excludedPathPrefixes = mergeWithDBList(func(ctx context.Context) ([]string, error) {
		return store.pathListLoadFromDB(ctx, SETTING_EXCLUDED_PATHS)
	}, opts.ExcludedPathPrefixes)
	store.includedPathPrefixes = mergeWithDBList(func(ctx context.Context) ([]string, error) {
		return store.pathListLoadFromDB(ctx, SETTING_INCLUDED_PATHS)
	}, opts.IncludedPathPrefixes)

	return store, nil
}

// mergeWithDBList merges a string list loaded from the settings table with a
// list provided via options. Option values win ordering; DB-only values are
// appended. On load error the option list is returned unchanged.
func mergeWithDBList(load func(ctx context.Context) ([]string, error), fromOpts []string) []string {
	merged := append([]string{}, fromOpts...)

	dbValues, err := load(context.Background())
	if err != nil {
		return merged
	}

	for _, dbValue := range dbValues {
		found := false
		for _, existing := range merged {
			if existing == dbValue {
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, dbValue)
		}
	}

	return merged
}
