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

func TestHandleComparisonAjaxSuccess(t *testing.T) {
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

	form := strings.NewReader("action=comparison-ajax&period=this-week")
	req := httptest.NewRequest(http.MethodPost, "/admin/home", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "statCards") {
		t.Fatalf("expected statCards in response, got: %s", body)
	}
	if !strings.Contains(body, "comparisonRows") {
		t.Fatalf("expected comparisonRows in response, got: %s", body)
	}
}

func TestHandleComparisonAjaxError(t *testing.T) {
	store := newTestStore(t, false)
	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	form := strings.NewReader("action=comparison-ajax")
	req := httptest.NewRequest(http.MethodPost, "/admin/home", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "error") {
		t.Fatalf("expected error in response, got: %s", body)
	}
}
