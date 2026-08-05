package settings

import (
	_ "embed"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/statsstore/admin/shared"
)

//go:embed settings.html
var settingsHTML string

//go:embed settings.js
var settingsJS string

// New creates a new settings controller
func New(ui ControllerOptions) http.Handler {
	return &Controller{UI: ui}
}

// Controller handles the settings page
type Controller struct {
	UI ControllerOptions
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(c.Handler(w, r)))
}

// Handler renders the controller output using the shared layout
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetString(r, "action")

	// AJAX endpoints for Vue.js
	switch action {
	case "list-ajax":
		return c.handleListAjax(w, r)
	case "add-ip-ajax":
		return c.handleAddIpAjax(w, r)
	case "remove-ip-ajax":
		return c.handleRemoveIpAjax(w, r)
	case "delete-visitors-ajax":
		return c.handleDeleteVisitorsAjax(w, r)
	}

	c.UI.Layout.SetTitle("Settings | Visitor Analytics")

	scriptURLs := []string{
		cdn.VueJs_3_5_32(),
	}

	scripts := []string{
		settingsJS,
	}

	c.UI.Layout.SetBody(c.pageShell(r).ToHTML())
	c.UI.Layout.SetScriptURLs(scriptURLs)
	c.UI.Layout.SetScripts(scripts)

	return c.UI.Layout.Render(w, r)
}

// pageShell builds the page shell (breadcrumbs, header, nav) and embeds
// the Vue.js settings template from settings.html. No DB queries are made here —
// all data is loaded via AJAX from the per-section endpoints.
func (c *Controller) pageShell(r *http.Request) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  shared.UrlHome(r),
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(r),
		},
		{
			Name: "Settings",
			URL:  shared.UrlSettings(r),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Settings")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, c.UI.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(hb.Raw(settingsHTML))
}
