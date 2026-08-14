package admin

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/statsstore/admin/dashboard"
	"github.com/dracory/statsstore/admin/ipdetails"
	"github.com/dracory/statsstore/admin/sessions"
	"github.com/dracory/statsstore/admin/settings"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/dracory/statsstore/admin/visitors"
)

// Handler dispatches the request to the matching controller based on the
// "controller" query parameter. It writes the response directly to the
// ResponseWriter.
func Handler(opts shared.ControllerOptions, w http.ResponseWriter, r *http.Request) {
	controller := req.GetStringTrimmed(r, shared.ParamController)

	var body string
	switch controller {
	case shared.CONTROLLER_VISITORS:
		body = visitors.NewVisitorsController(opts).Handler(w, r)
	case shared.CONTROLLER_SESSIONS:
		body = sessions.NewSessionsController(opts).Handler(w, r)
	case shared.CONTROLLER_SETTINGS:
		body = settings.NewSettingsController(opts).Handler(w, r)
	case shared.CONTROLLER_IP_DETAILS:
		body = ipdetails.NewIPDetailsController(opts).Handler(w, r)
	default:
		body = dashboard.NewDashboardController(opts).Handler(w, r)
	}

	// Controllers return an HTML string for page renders; AJAX endpoints
	// write directly to the ResponseWriter and return "".
	if body != "" {
		_, _ = w.Write([]byte(body))
	}
}
