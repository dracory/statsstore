package home

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

func TestHandleDashboardAjaxSuccess(t *testing.T) {
	store := newTestStore(t, true)
	now := carbon.Now(carbon.UTC)
	visitor := statsstore.NewVisitor().
		SetID("v1").
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

	req := httptest.NewRequest(http.MethodGet, "/admin/home?action=dashboard-ajax&period=this-week", nil)
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "statCards") {
		t.Fatalf("expected statCards in response, got: %s", body)
	}
}

func TestHandleDashboardAjaxError(t *testing.T) {
	store := newTestStore(t, false)
	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/home?action=dashboard-ajax", nil)
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "error") {
		t.Fatalf("expected error in response, got: %s", body)
	}
}
