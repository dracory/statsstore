package home

import (
	"net/http"

	"github.com/dromara/carbon/v2"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

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
	w.Write([]byte(c.Handle(w, r)))
}

// Handle renders the controller to an HTML tag
func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)

	c.ui.Layout.SetTitle("Dashboard | Visitor Analytics")

	if errorMessage != "" {
		c.ui.Layout.
			SetBody(hb.Div().Class("alert alert-danger").Text(errorMessage).ToHTML())
		return c.ui.Layout.Render(w, r)
	}

	// Load required scripts asynchronously
	scripts := []string{
		// Load Chart.js
		`
		if (!window.Chart) {
			const loadChartJS = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
				await new Promise(resolve => script.onload = resolve);
				console.log('Chart.js loaded');
			};
			loadChartJS();
		}
		`,
		// Load HTMX
		`
		if (!window.htmx) {
			const loadHtmx = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = '` + cdn.Htmx_2_0_0() + `';
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
				script.src = '` + cdn.Sweetalert2_11() + `';
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

// == PRIVATE METHODS ==========================================================

// prepareData prepares the data for the home page
func (c *Controller) prepareData(r *http.Request) (data ControllerData, errorMessage string) {
	data.Request = r

	datesInRange := datesInRange(carbon.Now().SubDays(31), carbon.Now())

	dates := []string{}
	uniqueVisits := []int64{}
	totalVisits := []int64{}

	for _, date := range datesInRange {
		uniqueVisitorCount, err := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQueryOptions{
			CreatedAtGte: date + " 00:00:00",
			CreatedAtLte: date + " 23:59:59",
			Distinct:     statsstore.COLUMN_IP_ADDRESS,
		})

		if err != nil {
			return data, err.Error()
		}

		totalVisitorCount, err := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQueryOptions{
			CreatedAtGte: date + " 00:00:00",
			CreatedAtLte: date + " 23:59:59",
		})

		if err != nil {
			return data, err.Error()
		}

		dates = append(dates, date)
		uniqueVisits = append(uniqueVisits, uniqueVisitorCount)
		totalVisits = append(totalVisits, totalVisitorCount)
	}

	data.dates = dates
	data.uniqueVisits = uniqueVisits
	data.totalVisits = totalVisits

	return data, ""
}

// page builds the main page layout
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(data.Request, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  shared.UrlHome(data.Request),
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(data.Request),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Analytics Dashboard")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(data.Request, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(navigationPanel(data)).
		Child(cardStatsSummary(data))
}
