package home

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/statsstore/admin/shared"
)

func TestLiveVisitorCount(t *testing.T) {
	store := newTestStore(t, true)
	layout := &fakeLayout{renderReturn: "rendered"}
	controller := New(shared.ControllerOptions{
		Store:   store,
		Layout:  layout,
		HomeURL: "https://admin.local",
	}).(*Controller)

	req := httptest.NewRequest(http.MethodGet, "/admin/home", nil)
	count, err := controller.liveVisitorCount(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count < 0 {
		t.Fatalf("expected non-negative count, got %d", count)
	}
}
