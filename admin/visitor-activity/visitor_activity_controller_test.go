package visitoractivity

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	_ "modernc.org/sqlite"
)

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

func TestVisitorActivityControllerHandlerSuccess(t *testing.T) {
	store := newTestStore(t, true)

	now := carbon.Now()
	visitor := statsstore.NewVisitor().
		SetID("visitor-1").
		SetCountry("us").
		SetCreatedAt(now.ToDateTimeString(carbon.UTC)).
		SetIpAddress("127.0.0.1").
		SetUserReferrer("https://example.com").
		SetUserDevice("Desktop").
		SetUserBrowser("Firefox").
		SetUserBrowserVersion("118")
	visitor = visitor.SetFingerprint("fingerprint-same")

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

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-activity", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if body := rr.Body.String(); body != "rendered" {
		t.Fatalf("unexpected response body: %s", body)
	}

	if layout.title != "Visitor Activity | Visitor Analytics" {
		t.Fatalf("unexpected title: %s", layout.title)
	}

	if len(layout.scripts) != 1 {
		t.Fatalf("expected 1 script, got %d", len(layout.scripts))
	}

	if !strings.Contains(layout.body, "visitor-activity-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}

	if !layout.renderCalled {
		t.Fatalf("expected render to be called")
	}
}

func TestVisitorActivityControllerHandlerError(t *testing.T) {
	store := newTestStore(t, false)
	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-activity", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if len(layout.scripts) != 1 {
		t.Fatalf("expected 1 script even on error (page shell still renders), got %d", len(layout.scripts))
	}

	if !strings.Contains(layout.body, "visitor-activity-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}
}

func TestVisitorActivityControllerEmptyState(t *testing.T) {
	store := newTestStore(t, true)

	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:      store,
		Layout:     layout,
		HomeURL:    "https://admin.local",
		WebsiteUrl: "https://example.com",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-activity", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if !strings.Contains(layout.body, "visitor-activity-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}
}

func TestVisitorActivityControllerWithFilters(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/admin/visitor-activity?country=US&device=desktop", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if !strings.Contains(layout.body, "visitor-activity-app") {
		t.Fatalf("expected body to contain Vue app div, got: %s", layout.body)
	}
}

func TestVisitorActivityListAjax(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodPost, "/admin/visitor-activity", strings.NewReader(body))
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
