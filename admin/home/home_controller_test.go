package home

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
	title        string
	scripts      []string
	scriptURLs   []string
	body         string
	renderReturn string
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

func (l *fakeLayout) SetStyleURLs(styles []string) {}

func (l *fakeLayout) SetStyles(styles []string) {}

func (l *fakeLayout) SetBody(body string) {
	l.body = body
}

func (l *fakeLayout) SetCountryNameByIso2(func(string) (string, error)) {}

func (l *fakeLayout) Render(http.ResponseWriter, *http.Request) string {
	if l.renderReturn == "" {
		return "render"
	}
	return l.renderReturn
}

func TestHomeControllerHandleSuccess(t *testing.T) {
	store := newTestStore(t, true)
	now := carbon.Now()
	visitor := statsstore.NewVisitor().
		SetID("visitor-1").
		SetCountry("US").
		SetCreatedAt(now.ToDateTimeString(carbon.UTC)).
		SetIpAddress("127.0.0.1")

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatalf("failed to seed visitor: %v", err)
	}
	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/home", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if body := rr.Body.String(); body != "rendered" {
		t.Fatalf("unexpected body: %s", body)
	}

	if layout.title != "Dashboard | Visitor Analytics" {
		t.Fatalf("unexpected title: %s", layout.title)
	}

	if len(layout.scripts) != 3 {
		t.Fatalf("expected 3 scripts, got %d", len(layout.scripts))
	}

	if !strings.Contains(layout.body, "Visitor Analytics Dashboard") {
		t.Fatalf("expected dashboard content, got: %s", layout.body)
	}
}

func TestHomeControllerHandleError(t *testing.T) {
	store := newTestStore(t, false)
	layout := &fakeLayout{renderReturn: "rendered"}

	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/home", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if body := rr.Body.String(); body != "rendered" {
		t.Fatalf("unexpected body: %s", body)
	}

	if len(layout.scripts) != 0 {
		t.Fatalf("expected no scripts when error occurs, got %d", len(layout.scripts))
	}

	if !strings.Contains(strings.ToLower(layout.body), "no such table") {
		t.Fatalf("expected missing table error in body, got: %s", layout.body)
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
