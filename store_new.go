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
	VisitorTableName   string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		visitorTableName:   opts.VisitorTableName,
		db:                 neatDB,
		automigrateEnabled: opts.AutomigrateEnabled,
		debugEnabled:       opts.DebugEnabled,
		logger:             logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}
