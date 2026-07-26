package home

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

	if len(layout.scripts) != 1 {
		t.Fatalf("expected 1 script, got %d", len(layout.scripts))
	}

	if len(layout.scriptURLs) != 1 {
		t.Fatalf("expected 1 script URL (Vue.js), got %d", len(layout.scriptURLs))
	}

	if !strings.Contains(layout.body, "Visitor Analytics Dashboard") {
		t.Fatalf("expected dashboard content, got: %s", layout.body)
	}

	if !strings.Contains(layout.body, "dashboard-app") {
		t.Fatalf("expected Vue.js mount point in body, got: %s", layout.body)
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

	// Page shell should always render successfully (no DB queries in Handle now)
	req := httptest.NewRequest(http.MethodGet, "/admin/home", nil)
	rr := httptest.NewRecorder()

	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	if body := rr.Body.String(); body != "rendered" {
		t.Fatalf("unexpected body: %s", body)
	}

	// Page shell always includes scripts (Vue.js + dashboard app)
	if len(layout.scripts) != 1 {
		t.Fatalf("expected 1 script (dashboard app), got %d", len(layout.scripts))
	}

	if len(layout.scriptURLs) != 1 {
		t.Fatalf("expected 1 script URL (Vue.js), got %d", len(layout.scriptURLs))
	}

	if !strings.Contains(layout.body, "dashboard-app") {
		t.Fatalf("expected Vue.js mount point in body, got: %s", layout.body)
	}

	// DB errors now appear in AJAX endpoint responses, not the page shell
	ajaxReq := httptest.NewRequest(http.MethodGet, "/admin/home?action=comparison-ajax", nil)
	ajaxRR := httptest.NewRecorder()
	controller.ServeHTTP(ajaxRR, ajaxReq)

	ajaxBody := ajaxRR.Body.String()
	if !strings.Contains(ajaxBody, "error") {
		t.Fatalf("expected error in AJAX response when DB fails, got: %s", ajaxBody)
	}
}

func TestHomeControllerDashboardMetrics(t *testing.T) {
	store := newTestStore(t, true)
	now := carbon.Now(carbon.UTC)
	base := now.StdTime()

	visitors := []struct {
		fp   string
		ip   string
		t    time.Time
		path string
	}{
		{"fp1", "1.1.1.1", base.Add(-2 * time.Hour), "/"},
		{"fp1", "1.1.1.1", base.Add(-1 * time.Hour), "/docs"},
		{"fp2", "2.2.2.2", base.Add(-10 * time.Minute), "/"},
	}

	for _, v := range visitors {
		visitor := statsstore.NewVisitor().
			SetFingerprint(v.fp).
			SetIpAddress(v.ip).
			SetPath(v.path).
			SetCreatedAt(carbon.CreateFromStdTime(v.t).ToDateTimeString(carbon.UTC))
		if err := store.VisitorCreate(context.Background(), visitor); err != nil {
			t.Fatalf("failed to create visitor: %v", err)
		}
	}

	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	// Test the comparison AJAX endpoint — this provides statCards and comparisonRows
	req := httptest.NewRequest(http.MethodGet, "/admin/home?period=last-7-days&action=comparison-ajax", nil)
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	body := rr.Body.String()

	// The JSON response should contain the computed metrics
	for _, label := range []string{"Total Unique Visitors", "Bounce Rate", "Avg. Visit Duration", "comparisonRows"} {
		if !strings.Contains(body, label) {
			t.Errorf("expected JSON to contain %q", label)
		}
	}

	// 2 sessions, 1 bounce (fp2) -> 50.0%
	if !strings.Contains(body, "50.0%") {
		t.Errorf("expected bounce rate 50.0%% in JSON, got: %s", body)
	}
	// fp1 had a 1 hour interval between page views
	if !strings.Contains(body, "1h 0m") {
		t.Errorf("expected avg visit duration 1h 0m in JSON, got: %s", body)
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
