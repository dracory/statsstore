package settings

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// handleRemoveIpAjax removes an IP from the exclusion list
func (c *Controller) handleRemoveIpAjax(w http.ResponseWriter, r *http.Request) string {
	ip := strings.TrimSpace(req.GetString(r, "ip_address"))
	if ip == "" {
		api.Respond(w, r, api.Error("IP address cannot be empty"))
		return ""
	}

	if err := c.UI.Store.ExcludedIPRemove(r.Context(), ip); err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	api.Respond(w, r, api.Success("IP removed from exclusion list"))

	return ""
}
