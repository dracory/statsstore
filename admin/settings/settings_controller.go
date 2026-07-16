package settings

import (
	"net/http"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

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
	w.Write([]byte(c.Handler(w, r)))
}

// Handler renders the controller output using the shared layout
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")

	// Handle POST actions (add/remove IP)
	if r.Method == http.MethodPost {
		return c.handlePost(w, r, action)
	}

	c.UI.Layout.SetTitle("Settings | Visitor Analytics")

	data, err := c.prepareData(r)

	if err != nil {
		c.UI.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(err.Error()).ToHTML())
		return c.UI.Layout.Render(w, r)
	}

	c.UI.Layout.SetBody(c.page(data).ToHTML())
	return c.UI.Layout.Render(w, r)
}

// handlePost processes form submissions for adding/removing excluded IPs
func (c *Controller) handlePost(w http.ResponseWriter, r *http.Request, action string) string {
	if err := r.ParseForm(); err != nil {
		c.UI.Layout.SetTitle("Settings | Visitor Analytics")
		c.UI.Layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text("Failed to parse form: " + err.Error()).ToHTML())
		return c.UI.Layout.Render(w, r)
	}

	var actionError string

	switch action {
	case "add_ip":
		ip := strings.TrimSpace(r.FormValue("ip_address"))
		if ip == "" {
			actionError = "IP address cannot be empty"
		} else if err := c.UI.Store.ExcludedIPAdd(r.Context(), ip); err != nil {
			actionError = err.Error()
		}

	case "remove_ip":
		ip := strings.TrimSpace(r.FormValue("ip_address"))
		if ip == "" {
			actionError = "IP address cannot be empty"
		} else if err := c.UI.Store.ExcludedIPRemove(r.Context(), ip); err != nil {
			actionError = err.Error()
		}

	case "delete_visitors_by_ip":
		ip := strings.TrimSpace(r.FormValue("ip_address"))
		if ip == "" {
			actionError = "IP address cannot be empty"
		} else {
			count, err := c.UI.Store.VisitorDeleteByIP(r.Context(), ip)
			if err != nil {
				actionError = err.Error()
			} else {
				actionError = ""
				_ = count
			}
		}
	}

	// Redirect back to settings page (PRG pattern)
	redirectURL := shared.UrlSettings(r)
	if actionError != "" {
		redirectURL = shared.UrlSettings(r, map[string]string{"error": actionError})
	}

	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusSeeOther)
	return ""
}

// prepareData loads the current excluded IPs list
func (c *Controller) prepareData(r *http.Request) (ControllerData, error) {
	data := ControllerData{
		Request: r,
	}

	ips, err := c.UI.Store.ExcludedIPList(r.Context())
	if err != nil {
		return data, err
	}

	data.ExcludedIPs = ips
	data.ErrorMessage = r.URL.Query().Get("error")

	return data, nil
}

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
			Name: "Settings",
			URL:  shared.UrlSettings(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Settings")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.UI.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(c.excludedIPsCard(data))
}

// excludedIPsCard renders the excluded IPs management card
func (c *Controller) excludedIPsCard(data ControllerData) hb.TagInterface {
	// Error alert
	var errorAlert hb.TagInterface
	if data.ErrorMessage != "" {
		errorAlert = hb.Div().
			Class("alert alert-danger alert-dismissible fade show").
			Attr("role", "alert").
			Child(hb.Text(data.ErrorMessage)).
			Child(hb.Button().
				Class("btn-close").
				Attr("type", "button").
				Attr("data-bs-dismiss", "alert").
				Attr("aria-label", "Close"))
	}

	// Add IP form
	addForm := hb.Form().
		Class("d-flex gap-2").
		Attr("method", "POST").
		Attr("action", shared.UrlSettings(data.Request, map[string]string{"action": "add_ip"})).
		Child(hb.Input().
			Class("form-control").
			Attr("type", "text").
			Attr("name", "ip_address").
			Attr("placeholder", "e.g. 192.168.1.1").
			Attr("required", "required")).
		Child(hb.Button().
			Class("btn btn-primary").
			Attr("type", "submit").
			HTML("<i class=\"bi bi-plus-circle\"></i> Add IP"))

	// Current excluded IPs table
	var tableBody hb.TagInterface
	if len(data.ExcludedIPs) == 0 {
		tableBody = hb.TBody().
			Child(hb.TR().
				Child(hb.TD().
					Class("text-center text-muted py-3").
					Attr("colspan", "3").
					Text("No excluded IPs. Add one above.")))
	} else {
		tbody := hb.TBody()
		for _, ip := range data.ExcludedIPs {
			// Delete visitors by IP form
			deleteVisitorsForm := hb.Form().
				Class("d-inline").
				Attr("method", "POST").
				Attr("action", shared.UrlSettings(data.Request, map[string]string{"action": "delete_visitors_by_ip"})).
				Attr("onsubmit", "return confirm('Permanently delete ALL visitor records from IP "+ip+"? This cannot be undone.')").
				Child(hb.Input().
					Attr("type", "hidden").
					Attr("name", "ip_address").
					Attr("value", ip)).
				Child(hb.Button().
					Class("btn btn-sm btn-outline-danger").
					Attr("type", "submit").
					Attr("title", "Delete all visitor records from this IP").
					HTML("<i class=\"bi bi-trash\"></i> Delete Stats"))

			// Remove from excluded list form
			removeForm := hb.Form().
				Class("d-inline").
				Attr("method", "POST").
				Attr("action", shared.UrlSettings(data.Request, map[string]string{"action": "remove_ip"})).
				Child(hb.Input().
					Attr("type", "hidden").
					Attr("name", "ip_address").
					Attr("value", ip)).
				Child(hb.Button().
					Class("btn btn-sm btn-outline-secondary").
					Attr("type", "submit").
					Attr("title", "Remove from exclusion list").
					HTML("<i class=\"bi bi-x-circle\"></i> Stop Excluding"))

			tbody = tbody.Child(hb.TR().
				Child(hb.TD().
					Class("align-middle font-monospace").
					Text(ip)).
				Child(hb.TD().
					Class("align-middle text-nowrap").
					Child(deleteVisitorsForm)).
				Child(hb.TD().
					Class("align-middle text-nowrap").
					Child(removeForm)))
		}
		tableBody = tbody
	}

	table := hb.Table().
		Class("table table-striped table-hover mb-0").
		Child(hb.Thead().
			Child(hb.TR().
				Child(hb.TH().Class("w-50").Text("IP Address")).
				Child(hb.TH().Class("text-center").Text("Delete Stats")).
				Child(hb.TH().Class("text-center").Text("Stop Excluding")))).
		Child(tableBody)

	cardBody := hb.Div().Class("card-body")

	if errorAlert != nil {
		cardBody = cardBody.Child(errorAlert)
	}

	cardBody = cardBody.
		Child(hb.P().
			Class("text-muted small mb-3").
			Text("IP addresses in this list are excluded from visitor tracking. New visits from these IPs will not be recorded. You can also permanently delete existing visitor records for a given IP.")).
		Child(addForm).
		Child(hb.HR().Class("my-3")).
		Child(hb.Div().
			Class("table-responsive").
			Child(table))

	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("<i class=\"bi bi-shield-exclamation\"></i> Excluded IP Addresses"))).
		Child(cardBody)
}
