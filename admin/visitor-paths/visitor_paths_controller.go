package visitorpaths

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

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

// ToTag renders the controller to an HTML tag
func (c *visitorPathsController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := buildControllerData(r, c.ui.Store)

	if action := r.URL.Query().Get("action"); action == "export" {
		if errorMessage != "" {
			w.WriteHeader(http.StatusInternalServerError)
			return errorMessage
		}
		return c.exportCSV(w, data)
	}

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
	}

	c.ui.Layout.SetBody(c.page(data).ToHTML())
	c.ui.Layout.SetScripts(scripts)

	return c.ui.Layout.Render(w, r)
}

func (c *visitorPathsController) exportCSV(w http.ResponseWriter, data visitorPathsControllerData) string {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)

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

	if err := writer.Write(headers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Failed to generate export"
	}

	for _, visitor := range data.Paths {
		browser := strings.TrimSpace(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())
		absoluteURL := fullPathURL(c.ui, visitor.Path())
		row := []string{
			formatTimestamp(visitor.CreatedAt()),
			visitor.Path(),
			absoluteURL,
			resolvedCountryName(c.ui, visitor.Country()),
			visitor.IpAddress(),
			visitor.UserReferrer(),
			sessionLabel(data.Paths, visitor),
			visitor.UserDevice(),
			browser,
		}

		if err := writer.Write(row); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return "Failed to generate export"
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "Failed to generate export"
	}

	filename := fmt.Sprintf("visitor-paths-%s.csv", time.Now().UTC().Format("2006-01-02"))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return buffer.String()
}

// == PRIVATE METHODS ==========================================================

// page builds the main page layout
func (c *visitorPathsController) page(data visitorPathsControllerData) hb.TagInterface {
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
