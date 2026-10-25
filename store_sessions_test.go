package statsstore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// createVisitorAt inserts a visitor row with a controlled fingerprint and
// created_at so session grouping can be tested deterministically.
func createVisitorAt(t *testing.T, store StoreInterface, fingerprint, path, createdAt string) VisitorInterface {
	t.Helper()

	visitor := NewVisitor().
		SetFingerprint(fingerprint).
		SetIpAddress("1.2.3.4").
		SetPath(path).
		SetCreatedAt(createdAt)

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatal("unexpected error:", err)
	}

	return visitor
}

func TestSessionIDAssignedOnCreate(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	visitor := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")

	if visitor.GetSessionID() == "" {
		t.Fatal("expected session id to be assigned")
	}

	found, err := store.VisitorFindByID(context.Background(), visitor.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found.GetSessionID() != visitor.GetSessionID() {
		t.Fatal("expected persisted session id to match")
	}
}

func TestSessionGroupingGapWindow(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	v1 := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	v2 := createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:10:00") // within gap
	v3 := createVisitorAt(t, store, "fp-1", "/c", "2026-01-01 10:50:00") // 40 min after v2: new session
	v4 := createVisitorAt(t, store, "fp-2", "/d", "2026-01-01 10:20:00") // different fingerprint

	if v1.GetSessionID() != v2.GetSessionID() {
		t.Fatal("expected v1 and v2 to share a session")
	}
	if v2.GetSessionID() == v3.GetSessionID() {
		t.Fatal("expected v3 to start a new session after the gap")
	}
	if v4.GetSessionID() == v1.GetSessionID() || v4.GetSessionID() == v3.GetSessionID() {
		t.Fatal("expected v4 to have its own session")
	}
}

func TestSessionListAndCount(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// session A: two page views
	createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:10:00")
	// session B: one page view (bounce)
	createVisitorAt(t, store, "fp-1", "/c", "2026-01-01 11:00:00")
	// session C: different visitor
	createVisitorAt(t, store, "fp-2", "/d", "2026-01-01 09:00:00")
	createVisitorAt(t, store, "fp-2", "/e", "2026-01-01 09:05:00")

	count, err := store.SessionCount(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 sessions, got %d", count)
	}

	sessions, err := store.SessionList(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(sessions) != 3 {
		t.Fatalf("expected 3 sessions, got %d", len(sessions))
	}

	// sessions are ordered by last activity desc: B (11:00), A (10:10), C (09:05)
	b, a, c := sessions[0], sessions[1], sessions[2]

	if !b.IsBounce() {
		t.Fatal("expected session B to be a bounce (single page view)")
	}
	if b.GetEntryPath() != "/c" || b.GetExitPath() != "/c" {
		t.Fatalf("expected session B entry/exit /c, got %s/%s", b.GetEntryPath(), b.GetExitPath())
	}

	if a.GetPageViews() != 2 {
		t.Fatalf("expected session A to have 2 page views, got %d", a.GetPageViews())
	}
	if a.GetEntryPath() != "/a" || a.GetExitPath() != "/b" {
		t.Fatalf("expected session A entry/exit /a//b, got %s/%s", a.GetEntryPath(), a.GetExitPath())
	}
	if a.GetDurationSeconds() != 600 {
		t.Fatalf("expected session A duration 600s, got %d", a.GetDurationSeconds())
	}
	if a.GetStartedAt() != "2026-01-01 10:00:00" || a.GetEndedAt() != "2026-01-01 10:10:00" {
		t.Fatalf("unexpected session A bounds: %s - %s", a.GetStartedAt(), a.GetEndedAt())
	}

	if c.GetPageViews() != 2 || c.GetIpAddress() != "1.2.3.4" {
		t.Fatalf("unexpected session C: pageViews=%d ip=%s", c.GetPageViews(), c.GetIpAddress())
	}

	// limit paginates sessions
	limited, err := store.SessionList(context.Background(), VisitorQuery().SetLimit(2))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(limited) != 2 {
		t.Fatalf("expected 2 sessions with limit, got %d", len(limited))
	}
}

func TestSessionListFiltersApplyToRows(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:10:00")

	sessions, err := store.SessionList(context.Background(), VisitorQuery().SetPathExact("/a"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session for path /a, got %d", len(sessions))
	}
	if sessions[0].GetPageViews() != 1 {
		t.Fatalf("expected filtered session to count 1 page view, got %d", sessions[0].GetPageViews())
	}
}

func TestSessionGapBoundaryExactly30Minutes(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	v1 := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	v2 := createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:30:00") // exactly at the boundary: same session
	v3 := createVisitorAt(t, store, "fp-1", "/c", "2026-01-01 11:00:01") // just past the boundary: new session

	if v1.GetSessionID() != v2.GetSessionID() {
		t.Fatal("expected a view exactly 30 min later to share the session")
	}
	if v2.GetSessionID() == v3.GetSessionID() {
		t.Fatal("expected a view past the gap to start a new session")
	}
}

func TestSessionIDChainsAcrossViews(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Each view is within the gap of the previous, so all join one session
	// even though the first and last are more than 30 min apart.
	v1 := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	v2 := createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:25:00")
	v3 := createVisitorAt(t, store, "fp-1", "/c", "2026-01-01 10:50:00")

	if v1.GetSessionID() != v3.GetSessionID() {
		t.Fatal("expected chained views to keep the same session")
	}
	_ = v2
}

func TestSessionIDEmptyFingerprint(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Fingerprintless visitors still get a session id (fresh each time)
	visitor := NewVisitor().
		SetIpAddress("1.2.3.4").
		SetPath("/a").
		SetCreatedAt("2026-01-01 10:00:00")

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if visitor.GetSessionID() == "" {
		t.Fatal("expected session id even without fingerprint")
	}
}

func TestSessionIDPresetNotOverwritten(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	visitor := NewVisitor().
		SetFingerprint("fp-1").
		SetSessionID("custom-session").
		SetPath("/a").
		SetCreatedAt("2026-01-01 10:00:00")

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if visitor.GetSessionID() != "custom-session" {
		t.Fatalf("expected preset session id kept, got %s", visitor.GetSessionID())
	}
}

func TestSessionListEmpty(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sessions, err := store.SessionList(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(sessions) != 0 {
		t.Fatalf("expected no sessions, got %d", len(sessions))
	}

	count, err := store.SessionCount(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 sessions, got %d", count)
	}
}

func TestSessionListNilQuery(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if _, err := store.SessionList(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil query")
	}
	if _, err := store.SessionCount(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil query")
	}
	if _, err := store.SessionIDBackfill(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil query")
	}
}

func TestSessionListOffset(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 11:00:00") // gap -> second session
	createVisitorAt(t, store, "fp-1", "/c", "2026-01-01 12:00:00") // gap -> third session

	all, err := store.SessionList(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	page, err := store.SessionList(context.Background(), VisitorQuery().SetLimit(1).SetOffset(1))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(page) != 1 {
		t.Fatalf("expected 1 session on page, got %d", len(page))
	}
	if page[0].GetSessionID() != all[1].GetSessionID() {
		t.Fatal("expected offset pagination to return the second-newest session")
	}
}

func TestSessionCountDateRangeFilter(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-1", "/b", "2026-01-02 10:00:00")

	count, err := store.SessionCount(context.Background(), VisitorQuery().
		SetCreatedAtGte("2026-01-02 00:00:00"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 session in range, got %d", count)
	}
}

func TestSessionSoftDeletedRowExcluded(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	v1 := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:10:00")

	if err := store.VisitorSoftDeleteByID(context.Background(), v1.GetID()); err != nil {
		t.Fatal("unexpected error:", err)
	}

	sessions, err := store.SessionList(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].GetPageViews() != 1 {
		t.Fatalf("expected soft-deleted view excluded from session, got %d page views", sessions[0].GetPageViews())
	}
}

func TestSessionIDBackfill(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Simulate legacy rows by clearing session ids directly
	v1 := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	v2 := createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:10:00")
	v3 := createVisitorAt(t, store, "fp-1", "/c", "2026-01-01 11:00:00")

	_, err = store.GetDB().Exec("UPDATE visitor_table SET session_id = ''")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	updated, err := store.SessionIDBackfill(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if updated != 3 {
		t.Fatalf("expected 3 backfilled rows, got %d", updated)
	}

	f1, _ := store.VisitorFindByID(context.Background(), v1.GetID())
	f2, _ := store.VisitorFindByID(context.Background(), v2.GetID())
	f3, _ := store.VisitorFindByID(context.Background(), v3.GetID())

	if f1.GetSessionID() == "" || f2.GetSessionID() == "" || f3.GetSessionID() == "" {
		t.Fatal("expected all rows to have session ids after backfill")
	}
	if f1.GetSessionID() != f2.GetSessionID() {
		t.Fatal("expected v1/v2 to share a session after backfill")
	}
	if f2.GetSessionID() == f3.GetSessionID() {
		t.Fatal("expected v3 in a new session after backfill")
	}
}

func TestSessionListNewVisitorFromExistingData(t *testing.T) {
	data := NewVisitorFromExistingData(map[string]string{
		COLUMN_SESSION_ID: "sess-1",
		COLUMN_PATH:       "/x",
	})
	if data.GetSessionID() != "sess-1" {
		t.Fatalf("expected session id loaded from data, got %s", data.GetSessionID())
	}
}

func TestVisitorListSessionIDFilter(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	v1 := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-1", "/b", "2026-01-01 10:10:00")
	createVisitorAt(t, store, "fp-2", "/c", "2026-01-01 10:00:00")

	list, err := store.VisitorList(context.Background(), VisitorQuery().SetSessionID(v1.GetSessionID()))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 visitors in session, got %d", len(list))
	}
}

func TestMigrateUpAddsSessionIDToExistingTable(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Simulate a table created before the session_id column existed
	_, err = db.Exec("CREATE TABLE visitor_table (" +
		"id VARCHAR(40) PRIMARY KEY, " +
		"path VARCHAR(510), " +
		"fingerprint VARCHAR(40), " +
		"ip_address VARCHAR(40), " +
		"country VARCHAR(2), " +
		"user_accept_language VARCHAR(100), " +
		"user_accept_encoding VARCHAR(40), " +
		"user_agent VARCHAR(510), " +
		"user_os VARCHAR(50), " +
		"user_os_version VARCHAR(50), " +
		"user_device VARCHAR(100), " +
		"user_device_type VARCHAR(50), " +
		"user_browser VARCHAR(100), " +
		"user_browser_version VARCHAR(100), " +
		"user_referrer VARCHAR(510), " +
		"bot VARCHAR(3), " +
		"threat VARCHAR(3), " +
		"created_at DATETIME, " +
		"updated_at DATETIME, " +
		"soft_deleted_at DATETIME)")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Insert a legacy row without session_id
	_, err = db.Exec("INSERT INTO visitor_table " +
		"(id, path, fingerprint, ip_address, country, bot, threat, created_at, updated_at, soft_deleted_at) " +
		"VALUES ('legacy-1', '/old', 'fp-legacy', '9.9.9.9', '', 'no', 'no', " +
		"'2020-01-01 00:00:00', '2020-01-01 00:00:00', '9999-12-31 23:59:59')")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Column must exist and be readable
	found, err := store.VisitorFindByID(context.Background(), "legacy-1")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("expected legacy visitor")
	}
	if found.GetSessionID() != "" {
		t.Fatalf("expected empty session id on legacy row, got %s", found.GetSessionID())
	}

	// New writes after migration must have session ids
	v := createVisitorAt(t, store, "fp-new", "/new", "2026-01-01 10:00:00")
	if v.GetSessionID() == "" {
		t.Fatal("expected session id on new row after migration")
	}

	// Backfill handles the legacy row
	updated, err := store.SessionIDBackfill(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if updated != 1 {
		t.Fatalf("expected 1 backfilled row, got %d", updated)
	}
}

func TestVisitorRegisterAssignsSessionID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/page", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	if err := store.VisitorRegister(context.Background(), req); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.VisitorList(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 visitor, got %d", len(list))
	}
	if list[0].GetSessionID() == "" {
		t.Fatal("expected VisitorRegister to assign a session id")
	}
	if list[0].GetFingerprint() == "" {
		t.Fatal("expected fingerprint to be set")
	}
}

func TestVisitorUpdatePersistsSessionID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	visitor := createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	sessionID := visitor.GetSessionID()

	visitor.SetPath("/changed")
	if err := store.VisitorUpdate(context.Background(), visitor); err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.VisitorFindByID(context.Background(), visitor.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found.GetSessionID() != sessionID {
		t.Fatalf("expected session id %s kept on update, got %s", sessionID, found.GetSessionID())
	}
	if found.GetPath() != "/changed" {
		t.Fatalf("expected updated path, got %s", found.GetPath())
	}
}

func TestVisitorQueryOrderBySessionID(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	createVisitorAt(t, store, "fp-1", "/a", "2026-01-01 10:00:00")
	createVisitorAt(t, store, "fp-2", "/b", "2026-01-01 10:00:00")

	// session_id is whitelisted for OrderBy and Distinct
	list, err := store.VisitorList(context.Background(), VisitorQuery().SetOrderBy(COLUMN_SESSION_ID))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 visitors, got %d", len(list))
	}

	count, err := store.VisitorCount(context.Background(), VisitorQuery().SetDistinct(COLUMN_SESSION_ID))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 distinct sessions, got %d", count)
	}
}
