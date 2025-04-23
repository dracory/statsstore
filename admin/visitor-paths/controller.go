package visitorpaths

import (
	"net/http"
	"strconv"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/statsstore/admin/shared"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONSTRUCTOR ==============================================================

// New creates a new visitor paths controller
func New(ui shared.UIContext) http.Handler {
	return &Controller{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the visitor paths page
type Controller struct {
	ui shared.UIContext
}

// ControllerData contains the data needed for the visitor paths page
type ControllerData struct {
	paths      []statsstore.VisitorInterface
	page       int
	totalPages int
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.ToTag(w, r)
}

// ToTag renders the controller to an HTML tag
func (c *Controller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.ui.GetLayout().SetTitle("Visitor Paths | Visitor Analytics")

	if errorMessage != "" {
		c.ui.GetLayout().SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.ui.GetLayout().Render(w, r))
	}

	// Load required scripts asynchronously
	scripts := []string{
		// Load HTMX
		`
		if (!window.htmx) {
			const loadHtmx = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = 'https://unpkg.com/htmx.org@1.9.6';
				await new Promise(resolve => script.onload = resolve);
				console.log('HTMX loaded');
			};
			loadHtmx();
		}
		`,
		// Load SweetAlert2
		`
		if (!window.Swal) {
			const loadSwal = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = 'https://cdn.jsdelivr.net/npm/sweetalert2@11';
				await new Promise(resolve => script.onload = resolve);
				console.log('SweetAlert2 loaded');
			};
			loadSwal();
		}
		`,
		// Add export functionality
		`
		function exportTableToCSV(tableId, filename) {
			const table = document.getElementById(tableId);
			if (!table) return;
			
			let csv = [];
			const rows = table.querySelectorAll('tr');
			
			for (let i = 0; i < rows.length; i++) {
				const row = [], cols = rows[i].querySelectorAll('td, th');
				
				for (let j = 0; j < cols.length; j++) {
					row.push('"' + cols[j].innerText.replace(/"/g, '""') + '"');
				}
				
				csv.push(row.join(','));
			}
			
			const csvContent = csv.join('\n');
			const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
			const link = document.createElement('a');
			
			link.href = URL.createObjectURL(blob);
			link.setAttribute('download', filename);
			link.click();
		}
		`,
	}

	c.ui.GetLayout().SetBody(c.page(data).ToHTML())
	c.ui.GetLayout().SetScripts(scripts)

	return hb.Raw(c.ui.GetLayout().Render(w, r))
}

// ToHTML renders the controller to HTML string
func (c *Controller) ToHTML() string {
	return c.ToTag(c.ui.GetResponse(), c.ui.GetRequest()).ToHTML()
}

// == PRIVATE METHODS ==========================================================

// prepareData prepares the data for the visitor paths page
func (c *Controller) prepareData(r *http.Request) (data ControllerData, errorMessage string) {
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPage := 10
	offset := (pageInt - 1) * perPage

	// Get the most visited paths
	paths, err := c.ui.GetStore().VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		Limit:     perPage,
		Offset:    offset,
		OrderBy:   statsstore.COLUMN_CREATED_AT,
		SortOrder: sb.DESC,
		// Note: We need to implement proper path grouping in the statsstore package
	})

	if err != nil {
		return data, err.Error()
	}

	// Get total unique paths
	pathCount, err := c.ui.GetStore().VisitorCount(r.Context(), statsstore.VisitorQueryOptions{
		Distinct: "path",
	})

	if err != nil {
		return data, err.Error()
	}

	totalPages := (int(pathCount) + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	return ControllerData{
		paths:      paths,
		page:       pageInt,
		totalPages: totalPages,
	}, ""
}

// page builds the main page layout
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := c.ui.Breadcrumbs([]shared.Breadcrumb{
		{
			Name: "Home",
			URL:  c.ui.URL(c.ui.GetHomeURL(), nil),
		},
		{
			Name: "Visitor Analytics",
			URL:  c.ui.URL(c.ui.GetPathHome(), nil),
		},
		{
			Name: "Visitor Paths",
			URL:  c.ui.URL(c.ui.GetPathVisitorPaths(), nil),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Paths")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(c.ui.AdminHeader()).
		Child(hb.HR()).
		Child(title).
		Child(c.cardVisitorPaths(data))
}

// cardVisitorPaths creates the visitor paths card
func (c *Controller) cardVisitorPaths(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header bg-light d-flex justify-content-between align-items-center").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Most Visited Paths")).
			Child(hb.Div().
				Class("dropdown").
				Child(hb.Button().
					Class("btn btn-sm btn-outline-secondary dropdown-toggle").
					Attr("type", "button").
					Attr("data-bs-toggle", "dropdown").
					Attr("aria-expanded", "false").
					Text("Export")).
				Child(hb.UL().
					Class("dropdown-menu").
					Child(hb.LI().
						Child(hb.A().
							Class("dropdown-item").
							Href("#").
							Attr("onclick", "exportTableToCSV('visitor-paths-table', 'visitor_paths.csv')").
							Text("Export to CSV")))))).
		Child(hb.Div().
			Class("card-body").
			Child(c.tableVisitorPaths(data.paths)).
			Child(c.pagination(data.page, data.totalPages)))
}

// tableVisitorPaths creates the visitor paths table
func (c *Controller) tableVisitorPaths(paths []statsstore.VisitorInterface) hb.TagInterface {
	table := hb.Table().
		ID("visitor-paths-table").
		Class("table table-striped table-hover").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-light").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("URL"),
						hb.TH().Class("text-end").Text("Visit Count"),
						hb.TH().Text("Last Visit"),
						hb.TH().Text("Actions"),
					}),
				}),
			hb.Tbody().Children(lo.Map(paths, func(path statsstore.VisitorInterface, index int) hb.TagInterface {
				// For now, we'll just show the path and created date
				// In a real implementation, we would need to add count functionality to the statsstore
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(shared.StrTruncate(path.Path(), 50)),
					hb.TD().Class("text-end").Text("1"), // Placeholder for count
					hb.TD().Text(path.CreatedAt()),
					hb.TD().Child(hb.A().
						Class("btn btn-sm btn-outline-primary").
						Attr("data-bs-toggle", "tooltip").
						Attr("title", "View visitors for this path").
						Href(c.ui.URL(c.ui.GetPathVisitorActivity(), map[string]string{
							"path": path.Path(),
						})).
						Child(hb.I().Class("bi bi-eye"))),
				})
			})),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}

// pagination creates the pagination component
func (c *Controller) pagination(page int, totalPages int) hb.TagInterface {
	if totalPages <= 1 {
		return hb.Div()
	}

	urlFunc := func(p int) string {
		return c.ui.URL(c.ui.GetPathVisitorPaths(), map[string]string{
			"page": cast.ToString(p),
		})
	}

	return shared.PaginationUI(page, totalPages, urlFunc)
}
