package home

import (
	"fmt"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

// liveVisitorCard renders a live visitor count card that auto-refreshes via
// HTMX every 30 seconds. The card itself is the HTMX target so the whole
// element is replaced with each poll.
func liveVisitorCard(count int64, r *http.Request) hb.TagInterface {
	refreshURL := shared.UrlHome(r, map[string]string{"action": "live"})

	return hb.Div().
		ID("live-visitor-card").
		Class("col").
		Attr("hx-get", refreshURL).
		Attr("hx-trigger", "every 30s").
		Attr("hx-target", "#live-visitor-card").
		Attr("hx-swap", "outerHTML").
		Child(shared.StatCardUI("Live Visitors", fmt.Sprintf("%d", count), "bi bi-broadcast", "danger"))
}
