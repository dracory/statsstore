package home

import (
	"fmt"
	"net/http"
)

// handleLiveAjax returns just the live visitor count as JSON.
func (c *Controller) handleLiveAjax(w http.ResponseWriter, r *http.Request) string {
	count, err := c.liveVisitorCount(r)
	if err != nil {
		count = 0
	}
	w.Header().Set("Content-Type", "application/json")
	return fmt.Sprintf(`{"liveVisitorCount":%d}`, count)
}
