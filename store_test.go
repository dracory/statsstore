package statsstore

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	dsn := ":memory:?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore() (StoreInterface, error) {
	db, err := initDB()

	if err != nil {
		return nil, err
	}

	return NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
	})
}

func TestStoreVisitorCreate(t *testing.T) {
	store, err := initStore()

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

func TestStoreVisitorFindByID(t *testing.T) {
	store, err := initStore()

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

	visitorFound, errFind := store.VisitorFindByID(ctx, visitor.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if visitorFound == nil {
		t.Fatal("visitor MUST NOT be nil")
	}

	if visitorFound.GetID() != visitor.GetID() {
		t.Fatal("IDs do not match")
	}
}

func TestStoreVisitorSoftDelete(t *testing.T) {
	store, err := initStore()

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

	err = store.VisitorSoftDeleteByID(ctx, visitor.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if visitor.GetSoftDeletedAt() != MAX_DATETIME {
		t.Fatal("visitor MUST NOT be soft deleted")
	}

	visitorFound, errFind := store.VisitorFindByID(ctx, visitor.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if visitorFound != nil {
		t.Fatal("visitor MUST be nil after soft delete")
	}

	visitorFindWithDeleted, err := store.VisitorList(ctx, VisitorQuery().
		SetID(visitor.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(visitorFindWithDeleted) == 0 {
		t.Fatal("visitor MUST be found with soft deleted included")
	}

	if strings.Contains(visitorFindWithDeleted[0].GetSoftDeletedAt(), MAX_DATETIME) {
		t.Fatal("visitor MUST be soft deleted", visitor.GetSoftDeletedAt())
	}
}

func TestStoreVisitorListCreatedAtRange(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	now := time.Now().UTC()

	visitors := []VisitorInterface{
		NewVisitor().SetCreatedAt(now.Add(-2 * time.Hour).Format(time.RFC3339)),
		NewVisitor().SetCreatedAt(now.Add(-10 * time.Minute).Format(time.RFC3339)),
		NewVisitor().SetCreatedAt(now.Add(-48 * time.Hour).Format(time.RFC3339)),
	}

	for _, v := range visitors {
		if err := store.VisitorCreate(ctx, v); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	from := now.Add(-24 * time.Hour).Format(time.RFC3339)
	to := now.Format(time.RFC3339)
	result, err := store.VisitorList(ctx, VisitorQuery().SetCreatedAtGte(from).SetCreatedAtLte(to))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 visitors in range, got %d", len(result))
	}
}

func TestExcludedIPAdd(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPAdd(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ExcludedIPAdd(ctx, "10.0.0.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	ips, err := store.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 2 {
		t.Fatalf("expected 2 excluded IPs, got %d", len(ips))
	}
}

func TestExcludedIPAddDuplicate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPAdd(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ExcludedIPAdd(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	ips, err := store.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 1 {
		t.Fatalf("expected 1 excluded IP (duplicate ignored), got %d", len(ips))
	}
}

func TestExcludedIPAddEmpty(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPAdd(ctx, ""); err == nil {
		t.Fatal("expected error for empty IP, got nil")
	}

	if err := store.ExcludedIPAdd(ctx, "   "); err == nil {
		t.Fatal("expected error for whitespace-only IP, got nil")
	}
}

func TestExcludedIPAddTrimsWhitespace(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPAdd(ctx, "  192.168.1.1  "); err != nil {
		t.Fatal("unexpected error:", err)
	}

	ips, err := store.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 1 {
		t.Fatalf("expected 1 excluded IP, got %d", len(ips))
	}

	if ips[0] != "192.168.1.1" {
		t.Fatalf("expected trimmed IP '192.168.1.1', got '%s'", ips[0])
	}
}

func TestExcludedIPRemove(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPAdd(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ExcludedIPAdd(ctx, "10.0.0.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ExcludedIPRemove(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	ips, err := store.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 1 {
		t.Fatalf("expected 1 excluded IP after removal, got %d", len(ips))
	}

	if ips[0] != "10.0.0.1" {
		t.Fatalf("expected remaining IP '10.0.0.1', got '%s'", ips[0])
	}
}

func TestExcludedIPRemoveNonExistent(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPAdd(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.ExcludedIPRemove(ctx, "10.0.0.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	ips, err := store.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 1 {
		t.Fatalf("expected 1 excluded IP (unchanged), got %d", len(ips))
	}
}

func TestExcludedIPRemoveEmpty(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.ExcludedIPRemove(ctx, ""); err == nil {
		t.Fatal("expected error for empty IP, got nil")
	}
}

func TestExcludedIPListEmpty(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	ips, err := store.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 0 {
		t.Fatalf("expected 0 excluded IPs, got %d", len(ips))
	}
}

func TestExcludedIPListPersistsAcrossInstances(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	store1, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store1.ExcludedIPAdd(ctx, "192.168.1.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store1.ExcludedIPAdd(ctx, "10.0.0.1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	store2, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ips, err := store2.ExcludedIPList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(ips) != 2 {
		t.Fatalf("expected 2 excluded IPs persisted across store instances, got %d", len(ips))
	}
}

func TestVisitorDeleteByIP(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	v1 := NewVisitor().SetIpAddress("192.168.1.1")
	v2 := NewVisitor().SetIpAddress("192.168.1.1")
	v3 := NewVisitor().SetIpAddress("10.0.0.1")

	for _, v := range []VisitorInterface{v1, v2, v3} {
		if err := store.VisitorCreate(ctx, v); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	deleted, err := store.VisitorDeleteByIP(ctx, "192.168.1.1")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if deleted != 2 {
		t.Fatalf("expected 2 deleted rows, got %d", deleted)
	}

	count, err := store.VisitorCount(ctx, VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 remaining visitor, got %d", count)
	}
}

func TestVisitorDeleteByIPNoMatch(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	v := NewVisitor().SetIpAddress("192.168.1.1")
	if err := store.VisitorCreate(ctx, v); err != nil {
		t.Fatal("unexpected error:", err)
	}

	deleted, err := store.VisitorDeleteByIP(ctx, "10.0.0.1")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if deleted != 0 {
		t.Fatalf("expected 0 deleted rows, got %d", deleted)
	}

	count, err := store.VisitorCount(ctx, VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 remaining visitor, got %d", count)
	}
}

func TestVisitorDeleteByIPEmpty(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	_, err = store.VisitorDeleteByIP(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty IP, got nil")
	}
}

func TestStoreVisitorRegisterExcludedPath(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		VisitorTableName:     "visitor_table",
		AutomigrateEnabled:   true,
		ExcludedPathPrefixes: []string{"/admin/"},
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Request to /admin/home should be skipped
	req1 := httptest.NewRequest(http.MethodGet, "/admin/home", nil)
	if err := store.VisitorRegister(ctx, req1); err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Request to /about should be tracked
	req2 := httptest.NewRequest(http.MethodGet, "/about", nil)
	if err := store.VisitorRegister(ctx, req2); err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, err := store.VisitorCount(ctx, VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 visitor (admin path excluded), got %d", count)
	}
}
