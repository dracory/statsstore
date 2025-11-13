package visitorpaths

import (
	"context"
	"database/sql"
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	_ "modernc.org/sqlite"
)

func TestVisitorPathsControllerExportCSV(t *testing.T) {
	store := newTestStore(t, true)

	visitorOne := statsstore.NewVisitor().
		SetID("visitor-1").
		SetPath("/hello").
		SetCountry("US").
		SetCreatedAt("2023-01-02T15:04:05Z").
		SetIpAddress("127.0.0.1").
		SetUserReferrer("https://referrer.one").
		SetUserDevice("Desktop").
		SetUserBrowser("Firefox").
		SetUserBrowserVersion("118").
		SetFingerprint("fingerprint-same")
	seededVisitor(t, store, visitorOne)

	visitorTwo := statsstore.NewVisitor().
		SetID("visitor-2").
		SetPath("/world").
		SetCountry("US").
		SetCreatedAt("2023-01-02T14:00:00Z").
		SetIpAddress("192.168.0.2").
		SetUserReferrer("https://referrer.two").
		SetUserDevice("Mobile").
		SetUserBrowser("Firefox").
		SetUserBrowserVersion("118").
		SetFingerprint("fingerprint-same")
	seededVisitor(t, store, visitorTwo)

	handler := New(shared.ControllerOptions{
		Store:      store,
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-paths?action=export", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusOK {
		t.Fatalf("unexpected status: %d", status)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "text/csv; charset=utf-8" {
		t.Fatalf("unexpected content type: %s", contentType)
	}

	disposition := rr.Header().Get("Content-Disposition")
	if !strings.HasPrefix(disposition, "attachment; filename=\"visitor-paths-") || !strings.HasSuffix(disposition, ".csv\"") {
		t.Fatalf("unexpected content disposition: %s", disposition)
	}

	records, err := csv.NewReader(strings.NewReader(rr.Body.String())).ReadAll()
	if err != nil {
		t.Fatalf("failed to parse csv: %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("expected 3 csv rows, got %d", len(records))
	}

	expectedHeader := []string{
		"Visit Time",
		"Path",
		"Absolute URL",
		"Country",
		"IP Address",
		"Referrer",
		"Session",
		"Device",
		"Browser",
	}

	if !reflect.DeepEqual(records[0], expectedHeader) {
		t.Fatalf("unexpected header row: %+v", records[0])
	}

	firstDataRow := records[1]
	if firstDataRow[0] == "" {
		t.Fatalf("expected visit time value")
	}
	if firstDataRow[2] != "https://example.com/hello" {
		t.Fatalf("unexpected absolute url: %s", firstDataRow[2])
	}
	if firstDataRow[6] != "Sessions: 2" {
		t.Fatalf("unexpected session label: %s", firstDataRow[6])
	}
	if firstDataRow[7] != "Desktop" {
		t.Fatalf("unexpected device: %s", firstDataRow[7])
	}
	if firstDataRow[8] != "Firefox 118" {
		t.Fatalf("unexpected browser: %s", firstDataRow[8])
	}

	secondDataRow := records[2]
	if secondDataRow[6] != "Sessions: 2" {
		t.Fatalf("unexpected session label for second row: %s", secondDataRow[6])
	}
	if secondDataRow[1] != "/world" {
		t.Fatalf("unexpected path for second row: %s", secondDataRow[1])
	}
}

func TestVisitorPathsControllerExportCSVStoreError(t *testing.T) {
	handler := New(shared.ControllerOptions{
		Store: newTestStore(t, false),
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-paths?action=export", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Result().StatusCode; status != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", status)
	}

	if body := strings.TrimSpace(rr.Body.String()); !strings.Contains(strings.ToLower(body), "no such table") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func newTestStore(t testing.TB, automigrate bool) statsstore.StoreInterface {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?parseTime=true")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	store, err := statsstore.NewStore(statsstore.NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: automigrate,
	})
	if err != nil {
		_ = db.Close()
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return store
}

func seededVisitor(t testing.TB, store statsstore.StoreInterface, visitor statsstore.VisitorInterface) statsstore.VisitorInterface {
	t.Helper()
	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatalf("failed to seed visitor: %v", err)
	}
	return visitor
}
