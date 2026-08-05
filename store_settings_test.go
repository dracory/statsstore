package statsstore

import (
	"context"
	"testing"
)

func initStoreWithSettings() (StoreInterface, error) {
	db, err := initDB()
	if err != nil {
		return nil, err
	}

	return NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		SettingsTableName:  "settings_table",
		AutomigrateEnabled: true,
	})
}

func TestSettingGet_NonExistentKey(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	val, err := store.SettingGet(context.Background(), "does-not-exist")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if val != "" {
		t.Fatalf("expected empty string, got %q", val)
	}
}

func TestSettingSetAndGet(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.SettingSet(ctx, "test-key", "test-value"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	val, err := store.SettingGet(ctx, "test-key")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if val != "test-value" {
		t.Fatalf("expected %q, got %q", "test-value", val)
	}
}

func TestSettingSet_Upsert(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Insert.
	if err := store.SettingSet(ctx, "upsert-key", "v1"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	val, _ := store.SettingGet(ctx, "upsert-key")
	if val != "v1" {
		t.Fatalf("expected %q, got %q", "v1", val)
	}

	// Update.
	if err := store.SettingSet(ctx, "upsert-key", "v2"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	val, _ = store.SettingGet(ctx, "upsert-key")
	if val != "v2" {
		t.Fatalf("expected %q, got %q", "v2", val)
	}
}

func TestSettingHas(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	exists, err := store.SettingHas(ctx, "has-key")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if exists {
		t.Fatal("expected false for non-existent key")
	}

	if err := store.SettingSet(ctx, "has-key", "yes"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	exists, err = store.SettingHas(ctx, "has-key")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if !exists {
		t.Fatal("expected true for existing key")
	}
}

func TestSettingDelete(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create.
	if err := store.SettingSet(ctx, "delete-me", "value"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete.
	if err := store.SettingDelete(ctx, "delete-me"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Verify gone.
	val, err := store.SettingGet(ctx, "delete-me")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if val != "" {
		t.Fatalf("expected empty string after delete, got %q", val)
	}

	// Delete non-existent — should not error.
	if err := store.SettingDelete(ctx, "never-existed"); err != nil {
		t.Fatal("unexpected error for non-existent key:", err)
	}
}

func TestSettingGet_DefaultSettingsTable(t *testing.T) {
	// initStore does not set SettingsTableName, but NewStore defaults to
	// DEFAULT_SETTINGS_TABLE. The settings table should still work.
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Should be able to set and get.
	if err := store.SettingSet(ctx, "default-table-key", "ok"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	val, err := store.SettingGet(ctx, "default-table-key")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if val != "ok" {
		t.Fatalf("expected %q, got %q", "ok", val)
	}
}

func TestSettingGet_EmptyKey(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	_, err = store.SettingGet(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSettingSet_EmptyKey(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.SettingSet(context.Background(), "", "val"); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSettingDelete_EmptyKey(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.SettingDelete(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSettingHas_EmptyKey(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	_, err = store.SettingHas(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSettingList_Empty(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settings, err := store.SettingList(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(settings) != 0 {
		t.Fatalf("expected empty map, got %d items", len(settings))
	}
}

func TestSettingList_WithEntries(t *testing.T) {
	store, err := initStoreWithSettings()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	if err := store.SettingSet(ctx, "k1", "v1"); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := store.SettingSet(ctx, "k2", "v2"); err != nil {
		t.Fatal("unexpected error:", err)
	}

	settings, err := store.SettingList(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(settings) != 2 {
		t.Fatalf("expected 2 items, got %d", len(settings))
	}
	if settings["k1"] != "v1" {
		t.Fatalf("expected %q, got %q", "v1", settings["k1"])
	}
	if settings["k2"] != "v2" {
		t.Fatalf("expected %q, got %q", "v2", settings["k2"])
	}
}
