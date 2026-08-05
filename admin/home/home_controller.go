package home

import (
	_ "embed"
	"net/http"

	"github.com/dromara/carbon/v2"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/statsstore/admin/shared"
)

//go:embed home.html
var homeHTML string

//go:embed home.js
var homeJS string

// == CONSTRUCTOR ==============================================================

// New creates a new home page controller instance
func New(ui shared.ControllerOptions) http.Handler {
	return &Controller{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the dashboard home page
type Controller struct {
	ui shared.ControllerOptions
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(c.Handle(w, r)))
}

// Handle renders the controller to an HTML tag
func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) string {
	action := req.GetString(r, "action")

	// Per-section AJAX endpoints for Vue.js
	switch action {
	case "overview-ajax":
		return c.handleOverviewAjax(w, r)
	case "comparison-ajax":
		return c.handleComparisonAjax(w, r)
	case "dashboard-data-ajax":
		return c.handleDashboardDataAjax(w, r)
	case "live-ajax":
		return c.handleLiveAjax(w, r)
	case "export":
		return c.handleExport(w, r)
	}

	c.ui.Layout.SetTitle("Dashboard | Visitor Analytics")

	scriptURLs := []string{
		cdn.VueJs_3_5_32(),
	}

	scripts := []string{
		homeJS,
	}

	c.ui.Layout.SetBody(c.pageShell(r).ToHTML())
	c.ui.Layout.SetScriptURLs(scriptURLs)
	c.ui.Layout.SetScripts(scripts)

	return c.ui.Layout.Render(w, r)
}

// == PRIVATE METHODS ==========================================================

func previousPeriodBounds(selectedPeriod string, start, end *carbon.Carbon) (*carbon.Carbon, *carbon.Carbon, string) {
	switch selectedPeriod {
	case "today":
		return start.Copy().SubDays(1), end.Copy().SubDays(1), "Yesterday"
	case "yesterday":
		return start.Copy().SubDays(1), end.Copy().SubDays(1), "Day Before Yesterday"
	case "last-7-days", "previous-7-days":
		return start.Copy().SubDays(7), end.Copy().SubDays(7), "Previous 7 Days"
	case "this-week", "last-week":
		return start.Copy().SubDays(7), end.Copy().SubDays(7), "Previous Week"
	case "this-month":
		return start.Copy().SubMonths(1).StartOfMonth(), start.Copy().SubDays(1).EndOfDay(), "Previous Month"
	case "last-month":
		return start.Copy().SubMonths(1).StartOfMonth(), start.Copy().SubDays(1).EndOfDay(), "Previous Month"
	default: // this-week
		return start.Copy().SubDays(7), end.Copy().SubDays(7), "Previous Week"
	}
}

// pageShell builds the page shell (breadcrumbs, header, nav) and embeds
// the Vue.js dashboard template from home.html. No DB queries are made here —
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
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Analytics Dashboard")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(hb.Raw(homeHTML))
}
