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

func TestHandleExportSuccess(t *testing.T) {
	store := newTestStore(t, true)
	now := carbon.Now(carbon.UTC)
	visitor := statsstore.NewVisitor().
		SetIpAddress("1.1.1.1").
		SetCreatedAt(now.ToDateTimeString(carbon.UTC))

	if err := store.VisitorCreate(context.Background(), visitor); err != nil {
		t.Fatalf("failed to seed visitor: %v", err)
	}

	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/home?action=export&period=this-week", nil)
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "Date") || !strings.Contains(body, "Page Views") {
		t.Fatalf("expected CSV headers in response, got: %s", body)
	}
}

func TestHandleExportError(t *testing.T) {
	store := newTestStore(t, false)
	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/home?action=export", nil)
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 on DB error, got: %d", rr.Code)
	}
}
