package home

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/dromara/carbon/v2"
)

func TestHandleDashboardDataAjaxSuccess(t *testing.T) {
	store := newTestStore(t, true)
	now := carbon.Now(carbon.UTC)
	base := now.StdTime()

	for _, v := range []struct {
		fp, ip, path string
		t            interface{}
	}{
		{"fp1", "1.1.1.1", "/", base},
		{"fp2", "2.2.2.2", "/docs", base},
	} {
		visitor := statsstore.NewVisitor().
			SetFingerprint(v.fp).
			SetIpAddress(v.ip).
			SetPath(v.path).
			SetCreatedAt(carbon.CreateFromStdTime(base).ToDateTimeString(carbon.UTC))
		if err := store.VisitorCreate(nil, visitor); err != nil {
			t.Fatalf("failed to create visitor: %v", err)
		}
	}

	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  &fakeLayout{renderReturn: "rendered"},
		HomeURL: "https://admin.local",
	})

	form := strings.NewReader("action=dashboard-data-ajax&period=this-week")
	req := httptest.NewRequest(http.MethodPost, "/admin/home", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	body := rr.Body.String()
	for _, key := range []string{"dailyStats", "trafficCards", "heatmap", "totals"} {
		if !strings.Contains(body, key) {
			t.Errorf("expected JSON to contain %q, got: %s", key, body)
		}
	}
}

func TestHandleDashboardDataAjaxError(t *testing.T) {
	controller := New(shared.ControllerOptions{
		Store:   newTestStore(t, true),
		Layout:  &fakeLayout{renderReturn: "rendered"},
		HomeURL: "https://admin.local",
	})

	form := strings.NewReader("action=dashboard-data-ajax")
	req := httptest.NewRequest(http.MethodPost, "/admin/home", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "dailyStats") {
		t.Errorf("expected dailyStats in response even with no visitors, got: %s", body)
	}
}
