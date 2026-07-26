package home

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/statsstore/admin/shared"
)

func TestHandleLiveAjaxSuccess(t *testing.T) {
	store := newTestStore(t, true)
	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	})

	form := strings.NewReader("action=live-ajax")
	req := httptest.NewRequest(http.MethodPost, "/admin/home", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	controller.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "liveVisitorCount") {
		t.Fatalf("expected liveVisitorCount in response, got: %s", body)
	}
}
