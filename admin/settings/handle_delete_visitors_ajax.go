package settings

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// handleDeleteVisitorsAjax deletes all visitor records for a given IP
func (c *Controller) handleDeleteVisitorsAjax(w http.ResponseWriter, r *http.Request) string {
	ip := strings.TrimSpace(req.GetString(r, "ip_address"))
	if ip == "" {
		api.Respond(w, r, api.Error("IP address cannot be empty"))
		return ""
	}

	count, err := c.UI.Store.VisitorDeleteByIP(r.Context(), ip)
	if err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData(
		fmt.Sprintf("Deleted %d visitor record(s) for IP %s", count, ip),
		map[string]any{
			"deletedCount": count,
			"ip":           ip,
		},
	))

	return ""
}
