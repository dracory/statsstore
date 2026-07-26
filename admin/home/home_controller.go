package home

import (
	_ "embed"
	"net/http"

	"github.com/dromara/carbon/v2"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

//go:embed home.html
var homeHTML string

//go:embed home.js
var homeJS string

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

	// Per-section AJAX endpoints for Vue.js
	switch action {
	case "dashboard-ajax":
		return c.handleJSONAjax(w, r)
	case "overview-ajax":
		return c.handleOverviewAjax(w, r)
	case "comparison-ajax":
		return c.handleComparisonAjax(w, r)
	case "daily-ajax":
		return c.handleDailyAjax(w, r)
	case "traffic-ajax":
		return c.handleTrafficAjax(w, r)
	case "heatmap-ajax":
		return c.handleHeatmapAjax(w, r)
	case "live-ajax":
		return c.handleLiveAjax(w, r)
	case "export":
		return c.handleExport(w, r)
	}

	c.ui.Layout.SetTitle("Dashboard | Visitor Analytics")

	scriptURLs := []string{
		cdn.VueJs_3_5_32(),
	}

	scripts := []string{
		homeJS,
	}

	c.ui.Layout.SetBody(c.pageShell(r).ToHTML())
	c.ui.Layout.SetScriptURLs(scriptURLs)
	c.ui.Layout.SetScripts(scripts)

	return c.ui.Layout.Render(w, r)
}

// == PRIVATE METHODS ==========================================================

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

// pageShell builds the page shell (breadcrumbs, header, nav) and embeds
// the Vue.js dashboard template from home.html. No DB queries are made here —
// all data is loaded via AJAX from the per-section endpoints.
func (c *Controller) pageShell(r *http.Request) hb.TagInterface {
	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{
			Name: "Home",
			URL:  shared.UrlHome(r),
		},
		{
			Name: "Visitor Analytics",
			URL:  shared.UrlHome(r),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Analytics Dashboard")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(shared.AdminHeaderUI(r, c.ui.HomeURL)).
		Child(hb.HR()).
		Child(title).
		Child(hb.Raw(homeHTML))
}
