package settings

import (
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/req"
)

// handleAddIpAjax adds an IP to the exclusion list
func (c *Controller) handleAddIpAjax(w http.ResponseWriter, r *http.Request) string {
	ip := strings.TrimSpace(req.GetString(r, "ip_address"))
	if ip == "" {
		api.Respond(w, r, api.Error("IP address cannot be empty"))
		return ""
	}

	if err := c.UI.Store.ExcludedIPAdd(r.Context(), ip); err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	api.Respond(w, r, api.Success("IP added to exclusion list"))

	return ""
}
