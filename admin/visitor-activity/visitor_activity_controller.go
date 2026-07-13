package visitoractivity

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

// New creates a new visitor activity controller
func New(ui ControllerOptions) http.Handler {
	return &Controller{UI: ui}
}

// Controller handles the visitor activity page
type Controller struct {
	UI ControllerOptions
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Handler(w, r)))
}

// Handler renders the controller output using the shared layout
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := buildControllerData(r, c.UI.Store)

	if action := r.URL.Query().Get("action"); action == "export" {
		if errorMessage != "" {
			w.WriteHeader(http.StatusInternalServerError)
			return errorMessage
		}
		return c.exportCSV(w, data)
	}

	c.UI.Layout.SetTitle("Visitor Activity | Visitor Analytics")

	if errorMessage != "" {
		c.UI.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return c.UI.Layout.Render(w, r)
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
	}

	c.UI.Layout.SetBody(c.page(data).ToHTML())
	c.UI.Layout.SetScripts(scripts)

	return c.UI.Layout.Render(w, r)
}

func (c *Controller) exportCSV(w http.ResponseWriter, data ControllerData) string {
	headers := []string{
		"Visit Time",
		"Path",
		"Country",
		"IP Address",
		"Referrer",
		"Browser",
		"OS",
		"User Agent",
	}

	rows := make([][]string, 0, len(data.Visitors))
	for _, visitor := range data.Visitors {
		rows = append(rows, []string{
			formatVisitorTimestamp(visitor.GetCreatedAt()),
			visitor.GetPath(),
			strings.ToUpper(visitor.GetCountry()),
			visitor.GetIpAddress(),
			visitor.GetUserReferrer(),
			strings.TrimSpace(visitor.GetUserBrowser() + " " + visitor.GetUserBrowserVersion()),
			strings.TrimSpace(visitor.GetUserOs() + " " + visitor.GetUserOsVersion()),
			visitor.GetUserAgent(),
		})
	}

	return shared.ExportCSV(w, shared.ExportFilename("visitor-activity"), headers, rows)
}

// ToTag renders the controller to an HTML tag
func (c *Controller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	return hb.Raw(c.Handler(w, r))
}

// == PRIVATE METHODS ==========================================================

// page builds the main page layout
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(data.Request, []shared.Breadcrumb{
		{
			Name: "Home",
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(data.Request),
		},
		{
			Name: "Visitor Activity",
			URL:  shared.UrlVisitorActivity(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.UI.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(CardVisitorActivity(data, c.UI))
}
