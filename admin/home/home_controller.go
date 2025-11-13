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

	now := carbon.Now()
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
	createdAtGte := start.ToDateString() + " 00:00:00"
	createdAtLte := end.ToDateString() + " 23:59:59"

	visitors, err := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		CreatedAtGte: createdAtGte,
		CreatedAtLte: createdAtLte,
	})

	if err != nil {
		return data, err.Error()
	}

	dailyPageViews := map[string]int64{}
	dailyUniqueIPs := map[string]map[string]struct{}{}
	firstVisitByIP := map[string]string{}

	for _, visitor := range visitors {
		createdAt := visitor.CreatedAtCarbon()
		if createdAt == nil {
			continue
		}

		visitDate := createdAt.ToDateString()
		identifier := visitor.IpAddress()
		if identifier == "" {
			identifier = "unknown-ip"
		}

		dailyPageViews[visitDate]++

		if _, ok := dailyUniqueIPs[visitDate]; !ok {
			dailyUniqueIPs[visitDate] = map[string]struct{}{}
		}

		dailyUniqueIPs[visitDate][identifier] = struct{}{}

		if existingDate, ok := firstVisitByIP[identifier]; !ok || visitDate < existingDate {
			firstVisitByIP[identifier] = visitDate
		}
	}

	dates := make([]string, 0, len(dateRange))
	uniqueVisits := make([]int64, 0, len(dateRange))
	totalVisits := make([]int64, 0, len(dateRange))
	firstVisits := make([]int64, 0, len(dateRange))
	returnVisits := make([]int64, 0, len(dateRange))

	for _, date := range dateRange {
		dates = append(dates, date)

		uniqueSet := dailyUniqueIPs[date]
		uniqueCount := int64(len(uniqueSet))

		var firstCount int64
		for ip := range uniqueSet {
			if firstVisitByIP[ip] == date {
				firstCount++
			}
		}

		returnCount := uniqueCount - firstCount
		if returnCount < 0 {
			returnCount = 0
		}

		uniqueVisits = append(uniqueVisits, uniqueCount)
		totalVisits = append(totalVisits, dailyPageViews[date])
		firstVisits = append(firstVisits, firstCount)
		returnVisits = append(returnVisits, returnCount)
	}

	data.dates = dates
	data.uniqueVisits = uniqueVisits
	data.totalVisits = totalVisits
	data.firstVisits = firstVisits
	data.returnVisits = returnVisits
	data.selectedPeriod = selectedPeriod
	data.periodOptions = periodOptions

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
