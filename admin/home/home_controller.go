package home

import (
	"net/http"
	"strconv"

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
	action := r.URL.Query().Get("action")

	if action == "live" {
		return c.handleLive(w, r)
	}

	data, errorMessage := c.prepareData(r)

	if action == "export" {
		if errorMessage != "" {
			w.WriteHeader(http.StatusInternalServerError)
			return errorMessage
		}
		return c.exportCSV(w, data)
	}

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

func (c *Controller) exportCSV(w http.ResponseWriter, data ControllerData) string {
	headers := []string{
		"Date",
		"Page Views",
		"Unique Visits",
		"First Time Visits",
		"Returning Visits",
	}

	rows := make([][]string, 0, len(data.dates))
	for i, date := range data.dates {
		rows = append(rows, []string{
			formatSummaryDate(date),
			strconv.FormatInt(data.totalVisits[i], 10),
			strconv.FormatInt(data.uniqueVisits[i], 10),
			strconv.FormatInt(data.firstVisits[i], 10),
			strconv.FormatInt(data.returnVisits[i], 10),
		})
	}

	return shared.ExportCSV(w, shared.ExportFilename("visitor-stats"), headers, rows)
}

func (c *Controller) handleLive(w http.ResponseWriter, r *http.Request) string {
	count, err := c.liveVisitorCount(r)
	if err != nil {
		count = 0
	}
	return liveVisitorCard(count, r).ToHTML()
}

func (c *Controller) liveVisitorCount(r *http.Request) (int64, error) {
	liveGte := carbon.Now(carbon.UTC).SubMinutes(15).ToDateTimeString(carbon.UTC)
	return c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().SetCreatedAtGte(liveGte))
}

func previousPeriodBounds(selectedPeriod string, start, end *carbon.Carbon) (*carbon.Carbon, *carbon.Carbon, string) {
	switch selectedPeriod {
	case "today":
		return start.Copy().SubDays(1), end.Copy().SubDays(1), "Yesterday"
	case "yesterday":
		return start.Copy().SubDays(1), end.Copy().SubDays(1), "Day Before Yesterday"
	case "last-7-days", "previous-7-days":
		return start.Copy().SubDays(7), end.Copy().SubDays(7), "Previous 7 Days"
	case "this-week", "last-week":
		return start.Copy().SubDays(7), end.Copy().SubDays(7), "Previous Week"
	case "this-month":
		return start.Copy().SubMonths(1).StartOfMonth(), start.Copy().SubDays(1).EndOfDay(), "Previous Month"
	case "last-month":
		return start.Copy().SubMonths(1).StartOfMonth(), start.Copy().SubDays(1).EndOfDay(), "Previous Month"
	default: // this-week
		return start.Copy().SubDays(7), end.Copy().SubDays(7), "Previous Week"
	}
}

// prepareData prepares the data for the home page

func (c *Controller) prepareData(r *http.Request) (data ControllerData, errorMessage string) {
	data.Request = r

	periodOptions := []periodOption{
		{Value: "today", Label: "Today"},
		{Value: "yesterday", Label: "Yesterday"},
		{Value: "last-7-days", Label: "Last 7 Days"},
		{Value: "previous-7-days", Label: "Previous 7 Days"},
		{Value: "this-week", Label: "This Week"},
		{Value: "last-week", Label: "Last Week"},
		{Value: "this-month", Label: "This Month"},
		{Value: "last-month", Label: "Last Month"},
	}

	selectedPeriod := r.URL.Query().Get("period")
	if selectedPeriod == "" {
		selectedPeriod = "this-week"
	}

	now := carbon.Now(carbon.UTC)
	start := now.Copy()
	end := now.Copy()

	switch selectedPeriod {
	case "today":
		start = now.Copy().StartOfDay()
		end = now.Copy().EndOfDay()
	case "yesterday":
		start = now.Copy().SubDays(1).StartOfDay()
		end = start.Copy().EndOfDay()
	case "last-7-days":
		start = now.Copy().SubDays(6).StartOfDay()
		end = now.Copy().EndOfDay()
	case "previous-7-days":
		end = now.Copy().SubDays(7).EndOfDay()
		start = end.Copy().SubDays(6).StartOfDay()
	case "last-week":
		start = now.SubWeeks(1).StartOfWeek()
		end = start.Copy().EndOfWeek()
	case "this-month":
		start = now.StartOfMonth()
		end = now.EndOfMonth()
	case "last-month":
		start = now.SubMonths(1).StartOfMonth()
		end = start.Copy().EndOfMonth()
	default: // this-week
		start = now.StartOfWeek()
		end = now.EndOfWeek()
	}

	dateRange := datesInRange(start.Copy(), end.Copy())
	createdAtGte := start.ToDateTimeString(carbon.UTC)
	createdAtLte := end.ToDateTimeString(carbon.UTC)

	visitors, err := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(createdAtGte).
		SetCreatedAtLte(createdAtLte))

	if err != nil {
		return data, err.Error()
	}

	currentStats := computePeriodStats(visitors, dateRange)

	prevStart, prevEnd, prevLabel := previousPeriodBounds(selectedPeriod, start, end)
	prevDateRange := datesInRange(prevStart.Copy(), prevEnd.Copy())
	prevCreatedAtGte := prevStart.ToDateTimeString(carbon.UTC)
	prevCreatedAtLte := prevEnd.ToDateTimeString(carbon.UTC)

	previousVisitors, err := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(prevCreatedAtGte).
		SetCreatedAtLte(prevCreatedAtLte))

	if err != nil {
		return data, err.Error()
	}

	previousStats := computePeriodStats(previousVisitors, prevDateRange)

	liveGte := carbon.Now(carbon.UTC).SubMinutes(15).ToDateTimeString(carbon.UTC)
	liveCount, err := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().SetCreatedAtGte(liveGte))
	if err != nil {
		return data, err.Error()
	}

	data.visitors = visitors
	data.ui = c.ui
	data.dates = currentStats.dates
	data.uniqueVisits = currentStats.uniqueVisits
	data.totalVisits = currentStats.totalVisits
	data.firstVisits = currentStats.firstVisits
	data.returnVisits = currentStats.returnVisits
	data.selectedPeriod = selectedPeriod
	data.periodOptions = periodOptions
	data.liveVisitorCount = liveCount
	data.previousPeriodLabel = prevLabel
	data.previousPeriodUnique = previousStats.totalUnique
	data.previousPeriodTotal = previousStats.totalTotal
	data.previousPeriodFirst = previousStats.totalFirst
	data.previousPeriodReturning = previousStats.totalReturning
	data.previousPeriodVisitors = previousVisitors
	data.currentStats = computeStatsOverview(visitors)
	data.previousStats = computeStatsOverview(previousVisitors)

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
		Child(cardStatsSummary(data)).
		Child(trafficSourcesCards(data))
}
