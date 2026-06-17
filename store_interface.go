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

	VisitorCount(ctx context.Context, query VisitorQueryInterface) (int64, error)
	VisitorCreate(ctx context.Context, user VisitorInterface) error
	VisitorDelete(ctx context.Context, user VisitorInterface) error
	VisitorDeleteByID(ctx context.Context, id string) error
	VisitorFindByID(ctx context.Context, userID string) (VisitorInterface, error)
	VisitorList(ctx context.Context, query VisitorQueryInterface) ([]VisitorInterface, error)
	VisitorRegister(ctx context.Context, r *http.Request) error
	VisitorSoftDelete(ctx context.Context, user VisitorInterface) error
	VisitorSoftDeleteByID(ctx context.Context, id string) error
	VisitorUpdate(ctx context.Context, user VisitorInterface) error
}
