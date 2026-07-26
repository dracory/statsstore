package home

import (
	"net/http"

	"github.com/dracory/api"
)

// handleLiveAjax returns just the live visitor count as JSON.
func (c *Controller) handleLiveAjax(w http.ResponseWriter, r *http.Request) string {
	count, err := c.liveVisitorCount(r)
	if err != nil {
		count = 0
	}

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"liveVisitorCount": count,
	}))

	return ""
}
