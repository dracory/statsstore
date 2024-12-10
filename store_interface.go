package statsstore

import (
	"context"
	"database/sql"
	"net/http"
)

type StoreInterface interface {
	AutoMigrate() error
	DB() *sql.DB
	EnableDebug(debug bool)
	VisitorCount(ctx context.Context, options VisitorQueryOptions) (int64, error)
	VisitorCreate(ctx context.Context, user VisitorInterface) error
	VisitorDelete(ctx context.Context, user VisitorInterface) error
	VisitorDeleteByID(ctx context.Context, id string) error
	VisitorFindByID(ctx context.Context, userID string) (VisitorInterface, error)
	VisitorList(ctx context.Context, options VisitorQueryOptions) ([]VisitorInterface, error)
	VisitorRegister(ctx context.Context, r *http.Request) error
	VisitorSoftDelete(ctx context.Context, user VisitorInterface) error
	VisitorSoftDeleteByID(ctx context.Context, id string) error
	VisitorUpdate(ctx context.Context, user VisitorInterface) error
}
