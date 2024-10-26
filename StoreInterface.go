package statsstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
	VisitorCreate(user VisitorInterface) error
	VisitorDelete(user VisitorInterface) error
	VisitorDeleteByID(id string) error
	VisitorFindByID(userID string) (VisitorInterface, error)
	VisitorList(options VisitorQueryOptions) ([]VisitorInterface, error)
	VisitorSoftDelete(user VisitorInterface) error
	VisitorSoftDeleteByID(id string) error
	VisitorUpdate(user VisitorInterface) error
}
