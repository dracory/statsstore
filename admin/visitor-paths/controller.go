package visitorpaths

import (
	"net/http"

	"github.com/dracory/statsstore/admin/shared"
	"github.com/gouniverse/hb"
)

// == CONSTRUCTOR ==============================================================

// New creates a new visitor paths controller
func New(ui shared.ControllerOptions) http.Handler {
	return &Controller{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the visitor paths page
type Controller struct {
	ui shared.ControllerOptions
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Handler(w, r)))
}

// ToTag renders the controller to an HTML tag
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := buildControllerData(r, c.ui.Store)

	c.ui.Layout.SetTitle("Visitor Paths | Visitor Analytics")

	if errorMessage != "" {
		c.ui.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return c.ui.Layout.Render(w, r)
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

	c.ui.Layout.SetBody(c.page(data).ToHTML())
	c.ui.Layout.SetScripts(scripts)

	return c.ui.Layout.Render(w, r)
}

// == PRIVATE METHODS ==========================================================

// page builds the main page layout
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(data.Request, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  c.ui.HomeURL,
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(data.Request),
		},
		{
			Name: "Visitor Paths",
			URL:  shared.UrlVisitorPaths(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Paths")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(CardVisitorPaths(data, c.ui))
}
