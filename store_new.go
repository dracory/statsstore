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
	BotFilterEnabled     bool
	ExcludedPathPrefixes []string
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
		visitorTableName:     opts.VisitorTableName,
		settingsTableName:    settingsTable,
		db:                   neatDB,
		automigrateEnabled:   opts.AutomigrateEnabled,
		debugEnabled:         opts.DebugEnabled,
		botFilterEnabled:     opts.BotFilterEnabled,
		excludedPathPrefixes: opts.ExcludedPathPrefixes,
		excludedIPs:          opts.ExcludedIPs,
		geoIPResolver:        opts.GeoIPResolver,
		enhanceBatchSize:     opts.EnhanceBatchSize,
		logger:               logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	// Load excluded IPs from the database and merge with any provided via options
	if dbIPs, err := store.excludedIPsLoadFromDB(context.Background()); err == nil {
		merged := append([]string{}, opts.ExcludedIPs...)
		for _, dbIP := range dbIPs {
			found := false
			for _, existing := range merged {
				if existing == dbIP {
					found = true
					break
				}
			}
			if !found {
				merged = append(merged, dbIP)
			}
		}
		store.excludedIPs = merged
	}

	return store, nil
}
