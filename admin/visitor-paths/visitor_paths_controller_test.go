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

// stripBOM removes the UTF-8 BOM prefix if present so csv.NewReader can parse cleanly.
func stripBOM(s string) string {
	return strings.TrimPrefix(s, "\xEF\xBB\xBF")
}

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

	records, err := csv.NewReader(strings.NewReader(stripBOM(rr.Body.String()))).ReadAll()
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

	if body := strings.TrimSpace(rr.Body.String()); !strings.Contains(strings.ToLower(body), "database operation failed") {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestVisitorPathsControllerHandlerSuccess(t *testing.T) {
	store := newTestStore(t, true)

	visitor := statsstore.NewVisitor().
		SetID("visitor-1").
		SetPath("/hello").
		SetCountry("us").
		SetCreatedAt("2023-01-02T15:04:05Z").
		SetIpAddress("127.0.0.1").
		SetUserReferrer("https://example.com").
		SetUserDevice("Desktop").
		SetUserBrowser("Firefox").
		SetUserBrowserVersion("118").
		SetFingerprint("fingerprint-same")

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatalf("failed to seed visitor: %v", err)
	}

	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-paths", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if body := rr.Body.String(); body != "rendered" {
		t.Fatalf("unexpected response body: %s", body)
	}

	if layout.title != "Visitor Paths | Visitor Analytics" {
		t.Fatalf("unexpected title: %s", layout.title)
	}

	if len(layout.scripts) != 1 {
		t.Fatalf("expected 1 script, got %d", len(layout.scripts))
	}

	if !strings.Contains(layout.body, "visitor-paths-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}

	if !layout.renderCalled {
		t.Fatalf("expected render to be called")
	}
}

func TestVisitorPathsControllerHandlerError(t *testing.T) {
	store := newTestStore(t, false)
	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-paths", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if len(layout.scripts) != 1 {
		t.Fatalf("expected 1 script even on error (page shell still renders), got %d", len(layout.scripts))
	}

	if !strings.Contains(layout.body, "visitor-paths-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}
}

func TestVisitorPathsControllerEmptyState(t *testing.T) {
	store := newTestStore(t, true)

	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-paths", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if !strings.Contains(layout.body, "visitor-paths-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}
}

func TestVisitorPathsControllerWithFilters(t *testing.T) {
	store := newTestStore(t, true)

	visitor := statsstore.NewVisitor().
		SetID("visitor-1").
		SetPath("/hello").
		SetCountry("US").
		SetCreatedAt("2023-01-02T15:04:05Z").
		SetIpAddress("127.0.0.1").
		SetUserDevice("Desktop").
		SetFingerprint("fingerprint-same")

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatalf("failed to seed visitor: %v", err)
	}

	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-paths?country=US&device=desktop", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if !strings.Contains(layout.body, "visitor-paths-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}
}

func TestVisitorPathsListAjax(t *testing.T) {
	store := newTestStore(t, true)

	visitor := statsstore.NewVisitor().
		SetID("visitor-1").
		SetPath("/hello").
		SetCountry("US").
		SetCreatedAt("2023-01-02T15:04:05Z").
		SetIpAddress("127.0.0.1").
		SetUserDevice("Desktop").
		SetFingerprint("fingerprint-same")

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatalf("failed to seed visitor: %v", err)
	}

	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	body := "action=list-ajax&page=1&per_page=10"
	req := httptest.NewRequest(http.MethodPost, "/admin/visitor-paths", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	respBody := rr.Body.String()
	if !strings.Contains(respBody, `"status":"success"`) {
		t.Fatalf("expected success status in JSON, got: %s", respBody)
	}

	if !strings.Contains(respBody, `"visitor-1"`) {
		t.Fatalf("expected visitor-1 in JSON, got: %s", respBody)
	}

	if !strings.Contains(respBody, `"totalCount":1`) {
		t.Fatalf("expected totalCount 1 in JSON, got: %s", respBody)
	}
}

// == TEST HELPERS ============================================================

type fakeLayout struct {
	title            string
	scripts          []string
	scriptURLs       []string
	styles           []string
	styleURLs        []string
	body             string
	renderReturn     string
	renderCalled     bool
	lastRenderReq    *http.Request
	lastRenderWriter http.ResponseWriter
	countryLookup    func(string) (string, error)
}

func (l *fakeLayout) SetTitle(title string) {
	l.title = title
}

func (l *fakeLayout) SetScriptURLs(scripts []string) {
	l.scriptURLs = append([]string{}, scripts...)
}

func (l *fakeLayout) SetScripts(scripts []string) {
	l.scripts = append([]string{}, scripts...)
}

func (l *fakeLayout) SetStyleURLs(styles []string) {
	l.styleURLs = append([]string{}, styles...)
}

func (l *fakeLayout) SetStyles(styles []string) {
	l.styles = append([]string{}, styles...)
}

func (l *fakeLayout) SetBody(body string) {
	l.body = body
}

func (l *fakeLayout) SetCountryNameByIso2(fn func(string) (string, error)) {
	l.countryLookup = fn
}

func (l *fakeLayout) Render(w http.ResponseWriter, r *http.Request) string {
	l.renderCalled = true
	l.lastRenderReq = r
	l.lastRenderWriter = w
	if l.renderReturn == "" {
		return "render"
	}
	return l.renderReturn
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
