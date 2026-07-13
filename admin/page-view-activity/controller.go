package pageviewactivity

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

// New creates a new page view activity controller.
func New(ui ControllerOptions) http.Handler {
	return &Controller{ui: ui}
}

// Controller handles rendering the page view activity screen.
type Controller struct {
	ui ControllerOptions
}

// ServeHTTP implements the http.Handler interface.
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Handler(w, r)))
}

// Handler prepares the layout and returns the rendered HTML.
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := buildControllerData(r, c.ui.Store)

	if action := r.URL.Query().Get("action"); action == "export" {
		if errorMessage != "" {
			w.WriteHeader(http.StatusInternalServerError)
			return errorMessage
		}
		return c.exportCSV(w, data)
	}

	c.ui.Layout.SetTitle("Page View Activity | Visitor Analytics")

	if errorMessage != "" {
		c.ui.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return c.ui.Layout.Render(w, r)
	}

	scripts := []string{
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
		`
		document.addEventListener('DOMContentLoaded', function() {
			var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
			tooltipTriggerList.map(function (el) {
				return new bootstrap.Tooltip(el);
			});
		});
		`,
	}

	c.ui.Layout.SetBody(c.page(data).ToHTML())
	c.ui.Layout.SetScripts(scripts)

	return c.ui.Layout.Render(w, r)
}

// exportCSV generates and writes a CSV export of the current page view data.
func (c *Controller) exportCSV(w http.ResponseWriter, data ControllerData) string {
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

	rows := make([][]string, 0, len(data.Visitors))
	for _, visitor := range data.Visitors {
		date, timeStr := splitTimestamp(visitor.GetCreatedAt())
		rows = append(rows, []string{
			date,
			timeStr,
			visitor.GetPath(),
			shared.FullPathURL(c.ui, visitor.GetPath()),
			shared.ResolvedCountryName(c.ui, visitor.GetCountry()),
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

// ToTag renders the controller to an HTML tag (useful for embedding).
func (c *Controller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	return hb.Raw(c.Handler(w, r))
}

// page constructs the main page container.
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
			Name: "Page View Activity",
			URL:  shared.UrlPageViewActivity(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Page View Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(CardPageViewActivity(data, c.ui))
}
