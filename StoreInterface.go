package statsstore

import "net/http"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
	VisitorCount(options VisitorQueryOptions) (int64, error)
	VisitorCreate(user VisitorInterface) error
	VisitorDelete(user VisitorInterface) error
	VisitorDeleteByID(id string) error
	VisitorFindByID(userID string) (VisitorInterface, error)
	VisitorList(options VisitorQueryOptions) ([]VisitorInterface, error)
	VisitorRegister(r *http.Request) error
	VisitorSoftDelete(user VisitorInterface) error
	VisitorSoftDeleteByID(id string) error
	VisitorUpdate(user VisitorInterface) error
}
