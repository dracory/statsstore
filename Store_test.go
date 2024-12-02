package statsstore

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) (*sql.DB, error) {
	if filepath != ":memory:" && utils.FileExists(filepath) {
		err := os.Remove(filepath) // remove database

		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			return nil, err
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore(filepath string) (StoreInterface, error) {
	db, err := initDB(filepath)

	if err != nil {
		return nil, err
	}

	return NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
	})
}

func TestStorevisitorCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	visitor := NewVisitor()

	err = store.VisitorCreate(context.Background(), visitor)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStorevisitorFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	visitor := NewVisitor()

	ctx := context.Background()

	err = store.VisitorCreate(ctx, visitor)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	visitorFound, errFind := store.VisitorFindByID(ctx, visitor.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if visitorFound == nil {
		t.Fatal("visitor MUST NOT be nil")
	}

	if visitorFound.ID() != visitor.ID() {
		t.Fatal("IDs do not match")
	}
}

func TestStorevisitorSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	ctx := context.Background()

	visitor := NewVisitor()

	err = store.VisitorCreate(ctx, visitor)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.VisitorSoftDeleteByID(ctx, visitor.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if visitor.DeletedAt() != sb.MAX_DATETIME {
		t.Fatal("visitor MUST NOT be soft deleted")
	}

	visitorFound, errFind := store.VisitorFindByID(ctx, visitor.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if visitorFound != nil {
		t.Fatal("visitor MUST be nil")
	}

	visitorFindWithDeleted, err := store.VisitorList(ctx, VisitorQueryOptions{
		ID:          visitor.ID(),
		Limit:       1,
		WithDeleted: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(visitorFindWithDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(visitorFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("visitor MUST be soft deleted", visitor.DeletedAt())
	}

}
