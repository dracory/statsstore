package statsstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

// == TYPE ====================================================================

type Store struct {
	visitorTableName   string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sqlStr := store.sqlVisitorTableCreate()

	if sqlStr == "" {
		return errors.New("visitor table create sql is empty")
	}

	if store.db == nil {
		return errors.New("visitorstore: database is nil")
	}

	_, err := store.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

// DB returns the database
func (store *Store) DB() *sql.DB {
	return store.db
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (store *Store) VisitorRegister(ctx context.Context, r *http.Request) error {
	path := r.URL.Path
	ip := utils.IP(r)
	userAgent := r.UserAgent()

	visitor := NewVisitor().
		SetPath(path).
		SetIpAddress(ip).
		SetUserAgent(userAgent)

	return store.VisitorCreate(ctx, visitor)
}

func (store *Store) VisitorCount(ctx context.Context, options VisitorQueryOptions) (int64, error) {
	options.CountOnly = true
	q := store.visitorQuery(options)

	if options.Distinct != "" {
		innerq := q.Select(options.Distinct).Distinct()

		q = goqu.Select(goqu.COUNT(goqu.Star()).As("count")).From(innerq)
	} else {
		q = q.Select(goqu.COUNT(goqu.Star()).As("count"))
	}

	q = q.Prepared(true).
		Limit(1)

	sqlStr, params, errSql := q.ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if options.Distinct != "" {
		sqlStr = strings.Replace(sqlStr, `AS "t1"`, "AS `t1`", 1)
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *Store) VisitorCreate(ctx context.Context, visitor VisitorInterface) error {
	visitor.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	visitor.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := visitor.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.visitorTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	if store.db == nil {
		return errors.New("visitorstore: database is nil")
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	visitor.MarkAsNotDirty()

	return nil
}

func (store *Store) VisitorDelete(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}

	return store.VisitorDeleteByID(ctx, visitor.ID())
}

func (store *Store) VisitorDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("visitor id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.visitorTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) VisitorFindByID(ctx context.Context, id string) (VisitorInterface, error) {
	if id == "" {
		return nil, errors.New("visitor id is empty")
	}

	list, err := store.VisitorList(ctx, VisitorQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) VisitorList(ctx context.Context, options VisitorQueryOptions) ([]VisitorInterface, error) {
	if store.db == nil {
		return []VisitorInterface{}, errors.New("visitorstore: database is nil")
	}

	q := store.visitorQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []VisitorInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

	if err != nil {
		return []VisitorInterface{}, err
	}

	list := []VisitorInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewVisitorFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) VisitorSoftDelete(ctx context.Context, visitor VisitorInterface) error {
	if visitor == nil {
		return errors.New("visitor is nil")
	}

	visitor.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.VisitorUpdate(ctx, visitor)
}

func (store *Store) VisitorSoftDeleteByID(ctx context.Context, id string) error {
	visitor, err := store.VisitorFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.VisitorSoftDelete(ctx, visitor)
}

func (store *Store) VisitorUpdate(ctx context.Context, visitor VisitorInterface) error {
	if store.db == nil {
		return errors.New("visitorstore: database is nil")
	}

	if visitor == nil {
		return errors.New("visitor is nil")
	}

	visitor.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := visitor.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updatable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.visitorTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(visitor.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	visitor.MarkAsNotDirty()

	return err
}

func (store *Store) toQuerableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}

	return database.Context(ctx, store.db)
}
