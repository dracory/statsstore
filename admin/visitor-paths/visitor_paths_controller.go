package visitorpaths

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

//go:embed visitor_paths.html
var visitorPathsHTML string

//go:embed visitor_paths.js
var visitorPathsJS string

// == CONSTRUCTOR ==============================================================

// New creates a new visitor paths controller
func New(ui shared.ControllerOptions) http.Handler {
	return &visitorPathsController{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the visitor paths page
type visitorPathsController struct {
	ui shared.ControllerOptions
}

// ServeHTTP implements the http.Handler interface
func (c *visitorPathsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Handler(w, r)))
}

// Handler renders the controller output using the shared layout
func (c *visitorPathsController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetString(r, "action")

	// AJAX endpoints for Vue.js
	switch action {
	case "list-ajax":
		return c.handleListAjax(w, r)
	case "export":
		return c.handleExport(w, r)
	}

	c.ui.Layout.SetTitle("Visitor Paths | Visitor Analytics")

	scriptURLs := []string{
		cdn.VueJs_3_5_32(),
	}

	scripts := []string{
		visitorPathsJS,
	}

	c.ui.Layout.SetBody(c.pageShell(r).ToHTML())
	c.ui.Layout.SetScriptURLs(scriptURLs)
	c.ui.Layout.SetScripts(scripts)

	return c.ui.Layout.Render(w, r)
}

// handleExport exports visitor path data as CSV
func (c *visitorPathsController) handleExport(w http.ResponseWriter, r *http.Request) string {
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
	if filters.PathContains != "" {
		options = options.SetPathContains(filters.PathContains)
	}
	if filters.PathExact != "" {
		options = options.SetPathExact(filters.PathExact)
	}
	if filters.Device != "" {
		options = options.SetDeviceType(filters.Device)
	}

	visitors, err := c.ui.Store.VisitorList(r.Context(), options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err.Error()
	}

	headers := []string{
		"Visit Time",
		"Path",
		"Absolute URL",
		"Country",
		"IP Address",
		"Referrer",
		"Session",
		"Device",
		"Browser",
	}

	rows := make([][]string, 0, len(visitors))
	for _, visitor := range visitors {
		browser := strings.TrimSpace(visitor.GetUserBrowser() + " " + visitor.GetUserBrowserVersion())
		absoluteURL := shared.FullPathURL(c.ui, visitor.GetPath())
		rows = append(rows, []string{
			shared.FormatTimestamp(visitor.GetCreatedAt()),
			visitor.GetPath(),
			absoluteURL,
			shared.ResolvedCountryName(c.ui, visitor.GetCountry()),
			visitor.GetIpAddress(),
			visitor.GetUserReferrer(),
			fmt.Sprintf("Sessions: %d", sessionCount(visitors, visitor)),
			visitor.GetUserDevice(),
			browser,
		})
	}

	return shared.ExportCSV(w, shared.ExportFilename("visitor-paths"), headers, rows)
}

// pageShell builds the page shell (breadcrumbs, header, nav) and embeds
// the Vue.js visitor paths template. No DB queries are made here —
// all data is loaded via AJAX from the per-section endpoints.
func (c *visitorPathsController) pageShell(r *http.Request) hb.TagInterface {
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
			Name: "Visitor Paths",
			URL:  shared.UrlVisitorPaths(r),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Paths")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(hb.Raw(visitorPathsHTML))
}
