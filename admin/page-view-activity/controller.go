package pageviewactivity

import (
	_ "embed"
	"net/http"
	"strings"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

//go:embed page_view_activity.html
var pageViewActivityHTML string

//go:embed page_view_activity.js
var pageViewActivityJS string

// New creates a new page view activity controller.
func New(ui ControllerOptions) http.Handler {
	return &Controller{UI: ui}
}

// Controller handles rendering the page view activity screen.
type Controller struct {
	UI ControllerOptions
}

// ServeHTTP implements the http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(c.Handler(w, r)))
}

// Handler prepares the layout and returns the rendered HTML.
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetString(r, "action")

	// AJAX endpoints for Vue.js
	switch action {
	case "list-ajax":
		return c.handleListAjax(w, r)
	case "export":
		return c.handleExport(w, r)
	}

	c.UI.Layout.SetTitle("Page View Activity | Visitor Analytics")

	scriptURLs := []string{
		cdn.VueJs_3_5_32(),
	}

	scripts := []string{
		pageViewActivityJS,
	}

	c.UI.Layout.SetBody(c.pageShell(r).ToHTML())
	c.UI.Layout.SetScriptURLs(scriptURLs)
	c.UI.Layout.SetScripts(scripts)

	return c.UI.Layout.Render(w, r)
}

// handleExport exports page view data as CSV
func (c *Controller) handleExport(w http.ResponseWriter, r *http.Request) string {
	filters := parseFiltersFromReq(r)
	page := shared.ParseIntWithDefault(req.GetString(r, "page"), 1)
	perPage := shared.ClampPerPage(shared.ParseIntWithDefault(req.GetString(r, "per_page"), 10))
	offset := (page - 1) * perPage

	options := statsstore.VisitorQuery().
		SetLimit(perPage).
		SetOffset(offset).
		SetOrderBy(statsstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC")

	if filters.Country != "" {
		options = options.SetCountry(filters.Country)
	}
	if filters.From != "" {
		options = options.SetCreatedAtGte(filters.From)
	}
	if filters.To != "" {
		options = options.SetCreatedAtLte(filters.To)
	}
	if filters.Device != "" {
		options = options.SetDeviceType(filters.Device)
	}

	visitors, err := c.UI.Store.VisitorList(r.Context(), options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err.Error()
	}

	headers := []string{
		"Date",
		"Time",
		"Path",
		"Absolute URL",
		"Country",
		"IP Address",
		"Referrer",
		"Device",
		"Browser",
		"OS",
		"User Agent",
	}

	rows := make([][]string, 0, len(visitors))
	for _, visitor := range visitors {
		date, timeStr := splitTimestamp(visitor.GetCreatedAt())
		rows = append(rows, []string{
			date,
			timeStr,
			shared.StripMethodPrefix(visitor.GetPath()),
			shared.FullPathURL(c.UI, visitor.GetPath()),
			shared.ResolvedCountryName(c.UI, visitor.GetCountry()),
			visitor.GetIpAddress(),
			visitor.GetUserReferrer(),
			visitor.GetUserDevice(),
			strings.TrimSpace(visitor.GetUserBrowser() + " " + visitor.GetUserBrowserVersion()),
			strings.TrimSpace(visitor.GetUserOs() + " " + visitor.GetUserOsVersion()),
			visitor.GetUserAgent(),
		})
	}

	return shared.ExportCSV(w, shared.ExportFilename("page-view-activity"), headers, rows)
}

// pageShell builds the page shell (breadcrumbs, header, nav) and embeds
// the Vue.js page view activity template. No DB queries are made here —
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
			Name: "Page View Activity",
			URL:  shared.UrlPageViewActivity(r),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Page View Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, c.UI.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(hb.Raw(pageViewActivityHTML))
}
