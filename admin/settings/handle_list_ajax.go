package settings

import (
	"net/http"

	"github.com/dracory/api"
)

// handleListAjax returns the current excluded IPs list as JSON
func (c *Controller) handleListAjax(w http.ResponseWriter, r *http.Request) string {
	ips, err := c.UI.Store.ExcludedIPList(r.Context())
	if err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"excludedIps": ips,
	}))

	return ""
}
