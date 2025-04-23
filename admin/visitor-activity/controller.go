package visitoractivity

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

// New creates a new visitor activity controller
func New(ui shared.UIContext) http.Handler {
	return &Controller{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the visitor activity page
type Controller struct {
	ui shared.UIContext
}

// ControllerData contains the data needed for the visitor activity page
type ControllerData struct {
	visitors   []statsstore.VisitorInterface
	page       int
	totalPages int
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.ToHTML()))
}

// ToTag renders the controller to an HTML tag
func (c *Controller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.ui.GetLayout().SetTitle("Visitor Activity | Visitor Analytics")

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

// prepareData prepares the data for the visitor activity page
func (c *Controller) prepareData(r *http.Request) (data ControllerData, errorMessage string) {
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPage := 10
	offset := (pageInt - 1) * perPage

	// Get visitors with pagination
	visitors, err := c.ui.GetStore().VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		Limit:     perPage,
		Offset:    offset,
		OrderBy:   statsstore.COLUMN_CREATED_AT,
		SortOrder: sb.DESC,
	})

	if err != nil {
		return data, err.Error()
	}

	visitorCount, err := c.ui.GetStore().VisitorCount(r.Context(), statsstore.VisitorQueryOptions{})
	if err != nil {
		return data, err.Error()
	}

	totalPages := (int(visitorCount) + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	return ControllerData{
		visitors:   visitors,
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
			Name: "Visitor Activity",
			URL:  c.ui.URL(c.ui.GetPathVisitorActivity(), nil),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(c.ui.AdminHeader()).
		Child(hb.HR()).
		Child(title).
		Child(c.cardVisitorActivity(data))
}

// cardVisitorActivity creates the visitor activity card
func (c *Controller) cardVisitorActivity(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header bg-light d-flex justify-content-between align-items-center").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Visitor Activity")).
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
							Attr("onclick", "exportTableToCSV('visitor-activity-table', 'visitor_activity.csv')").
							Text("Export to CSV")))))).
		Child(hb.Div().
			Class("card-body").
			Child(c.tableVisitorActivity(data.visitors)).
			Child(c.pagination(data.page, data.totalPages)))
}

// tableVisitorActivity creates the visitor activity table
func (c *Controller) tableVisitorActivity(visitors []statsstore.VisitorInterface) hb.TagInterface {
	table := hb.Table().
		ID("visitor-activity-table").
		Class("table table-striped table-hover").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-light").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("ID"),
						hb.TH().Text("IP Address"),
						hb.TH().Text("Path"),
						hb.TH().Text("Referrer"),
						hb.TH().Text("User Agent"),
						hb.TH().Text("Created At"),
						hb.TH().Text("Actions"),
					}),
				}),
			hb.Tbody().Children(lo.Map(visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(cast.ToString(visitor.ID())),
					hb.TD().Text(visitor.IpAddress()),
					hb.TD().Text(shared.StrTruncate(visitor.Path(), 30)),
					hb.TD().Text(shared.StrTruncate(visitor.UserReferrer(), 30)),
					hb.TD().Text(shared.StrTruncate(visitor.UserAgent(), 30)),
					hb.TD().Text(visitor.CreatedAt()),
					hb.TD().Child(hb.A().
						Class("btn btn-sm btn-outline-primary").
						Attr("data-bs-toggle", "tooltip").
						Attr("title", "View details").
						Href(c.ui.URL("/admin/visitor-activity/"+cast.ToString(visitor.ID()), nil)).
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
		return c.ui.URL(c.ui.GetPathVisitorActivity(), map[string]string{
			"page": cast.ToString(p),
		})
	}

	return shared.PaginationUI(page, totalPages, urlFunc)
}
