package statsstore

import (
	"context"
	"testing"
	"time"
)

func initStoreWith(t *testing.T, extra NewStoreOptions) StoreInterface {
	t.Helper()
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	extra.DB = db
	extra.VisitorTableName = "visitor_table"
	extra.SettingsTableName = "settings_table"
	extra.AutomigrateEnabled = true
	store, err := NewStore(extra)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return store
}

// TestVisitorCreateAsync verifies async inserts are queued, flushed by
// DrainVisitorWrites, and visible to VisitorCount.
func TestVisitorCreateAsync(t *testing.T) {
	store := initStoreWith(t, NewStoreOptions{
		VisitorCreateAsyncEnabled: true,
	})

	for i := 0; i < 5; i++ {
		if err := store.VisitorCreate(context.Background(), NewVisitor()); err != nil {
			t.Fatalf("VisitorCreate: %v", err)
		}
	}

	store.DrainVisitorWrites()

	count, err := store.VisitorCount(context.Background(), NewVisitorQuery())
	if err != nil {
		t.Fatalf("VisitorCount: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected 5 visitors after drain, got %d", count)
	}
	if store.VisitorDroppedCount() != 0 {
		t.Fatalf("expected 0 dropped, got %d", store.VisitorDroppedCount())
	}
}

// TestVisitorCreateAsyncDoesNotBlock verifies VisitorCreate returns
// promptly — the INSERT happens on the worker.
func TestVisitorCreateAsyncDoesNotBlock(t *testing.T) {
	store := initStoreWith(t, NewStoreOptions{
		VisitorCreateAsyncEnabled: true,
	})
	defer store.DrainVisitorWrites()

	start := time.Now()
	for i := 0; i < 20; i++ {
		if err := store.VisitorCreate(context.Background(), NewVisitor()); err != nil {
			t.Fatalf("VisitorCreate: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("VisitorCreate blocked for %v, expected near-instant enqueue", elapsed)
	}
}

// TestVisitorCreateAsyncAfterDrain verifies writes after Drain drop
// silently instead of panicking on a closed channel.
func TestVisitorCreateAsyncAfterDrain(t *testing.T) {
	store := initStoreWith(t, NewStoreOptions{
		VisitorCreateAsyncEnabled: true,
	})

	if err := store.VisitorCreate(context.Background(), NewVisitor()); err != nil {
		t.Fatalf("VisitorCreate: %v", err)
	}
	store.DrainVisitorWrites()

	if err := store.VisitorCreate(context.Background(), NewVisitor()); err != nil {
		t.Fatalf("VisitorCreate after drain: %v", err)
	}
	if store.VisitorDroppedCount() != 1 {
		t.Fatalf("expected 1 dropped, got %d", store.VisitorDroppedCount())
	}
}
