package home

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// dashboardJSON is the full JSON response sent to the Vue.js frontend.
type dashboardJSON struct {
	SelectedPeriod      string              `json:"selectedPeriod"`
	PeriodOptions       []periodOption      `json:"periodOptions"`
	LiveVisitorCount    int64               `json:"liveVisitorCount"`
	StatCards           []statCardJSON      `json:"statCards"`
	ComparisonRows      []comparisonRowJSON `json:"comparisonRows"`
	PreviousPeriodLabel string              `json:"previousPeriodLabel"`
	DailyStats          []dailyStatJSON     `json:"dailyStats"`
	Totals              totalsJSON          `json:"totals"`
	TrafficCards        []trafficCardJSON   `json:"trafficCards"`
	Heatmap             heatmapJSON         `json:"heatmap"`
	ChartLabels         []string            `json:"chartLabels"`
	ChartUniqueVisits   []int64             `json:"chartUniqueVisits"`
	ChartTotalVisits    []int64             `json:"chartTotalVisits"`
}

type statCardJSON struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

type comparisonRowJSON struct {
	Label    string  `json:"label"`
	Current  string  `json:"current"`
	Previous string  `json:"previous"`
	Change   float64 `json:"change"`
	Inverted bool    `json:"inverted"`
}

type dailyStatJSON struct {
	Date         string `json:"date"`
	TotalVisits  int64  `json:"totalVisits"`
	UniqueVisits int64  `json:"uniqueVisits"`
	FirstVisits  int64  `json:"firstVisits"`
	ReturnVisits int64  `json:"returnVisits"`
}

type totalsJSON struct {
	TotalVisits  int64 `json:"totalVisits"`
	UniqueVisits int64 `json:"uniqueVisits"`
	FirstVisits  int64 `json:"firstVisits"`
	ReturnVisits int64 `json:"returnVisits"`
}

type trafficCardJSON struct {
	Title      string           `json:"title"`
	ValueLabel string           `json:"valueLabel"`
	Tabs       []trafficTabJSON `json:"tabs"`
}

type trafficTabJSON struct {
	Label   string               `json:"label"`
	Entries []trafficSourceEntry `json:"entries"`
}

type heatmapJSON struct {
	Days        []string `json:"days"`
	Slots       []string `json:"slots"`
	Intensities [][]int  `json:"intensities"`
}

// handleJSON returns all dashboard data as a JSON response for the Vue.js frontend.
func (c *Controller) handleJSON(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)
	if errorMessage != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Sprintf(`{"error":%q}`, errorMessage)
	}

	return c.buildDashboardJSON(w, data)
}

// handleLiveJSON returns just the live visitor count as JSON.
func (c *Controller) handleLiveJSON(w http.ResponseWriter, r *http.Request) string {
	count, err := c.liveVisitorCount(r)
	if err != nil {
		count = 0
	}
	w.Header().Set("Content-Type", "application/json")
	return fmt.Sprintf(`{"liveVisitorCount":%d}`, count)
}

func (c *Controller) buildDashboardJSON(w http.ResponseWriter, data ControllerData) string {
	w.Header().Set("Content-Type", "application/json")

	totalUniqueVisitors := lo.Sum(data.uniqueVisits)
	totalVisitors := lo.Sum(data.totalVisits)
	totalFirstVisits := lo.Sum(data.firstVisits)
	totalReturningVisits := lo.Sum(data.returnVisits)

	days := len(data.dates)
	if days == 0 {
		days = 1
	}

	ext := data.currentStats

	cards := []statCardJSON{
		{"Total Unique Visitors", formatCount(totalUniqueVisitors), "bi bi-person", "primary"},
		{"Total Visitors", formatCount(totalVisitors), "bi bi-people", "success"},
		{"Avg. Unique Visitors", formatFloat(float64(totalUniqueVisitors) / float64(days)), "bi bi-graph-up", "info"},
		{"Avg. Total Visitors", formatFloat(float64(totalVisitors) / float64(days)), "bi bi-bar-chart", "warning"},
		{"Avg. Daily First Time Visits", formatFloat(float64(totalFirstVisits) / float64(days)), "bi bi-person-plus", "secondary"},
		{"Avg. Daily Returning Visits", formatFloat(float64(totalReturningVisits) / float64(days)), "bi bi-person-check", "dark"},
		{"Sessions", ext.Sessions, "bi bi-activity", "primary"},
		{"Pageviews", ext.Pageviews, "bi bi-collection", "success"},
		{"Pages per Session", ext.PagesPerSession, "bi bi-diagram-3", "info"},
		{"Bounce Rate", ext.BounceRate, "bi bi-arrow-repeat", "warning"},
		{"Avg. Visit Duration", ext.SessionDuration, "bi bi-clock-history", "secondary"},
	}

	// Build comparison rows
	prevUnique := data.previousPeriodUnique
	prevTotal := data.previousPeriodTotal
	prevFirst := data.previousPeriodFirst
	prevReturning := data.previousPeriodReturning

	comparisons := []comparisonRowJSON{
		{"Total Unique Visitors", formatCount(totalUniqueVisitors), formatCount(prevUnique), changePercentInt(totalUniqueVisitors, prevUnique), false},
		{"Total Visitors", formatCount(totalVisitors), formatCount(prevTotal), changePercentInt(totalVisitors, prevTotal), false},
		{"First Time Visits", formatCount(totalFirstVisits), formatCount(prevFirst), changePercentInt(totalFirstVisits, prevFirst), false},
		{"Returning Visits", formatCount(totalReturningVisits), formatCount(prevReturning), changePercentInt(totalReturningVisits, prevReturning), false},
		{"Bounce Rate", formatFloat2(ext.BounceRateValue) + "%", formatFloat2(data.previousStats.BounceRateValue) + "%", changePercentFloat(ext.BounceRateValue, data.previousStats.BounceRateValue), true},
		{"Avg. Visit Duration", formatDuration(ext.SessionDurationSeconds), formatDuration(data.previousStats.SessionDurationSeconds), changePercentFloat(ext.SessionDurationSeconds, data.previousStats.SessionDurationSeconds), false},
	}

	// Build daily stats
	daily := make([]dailyStatJSON, 0, len(data.dates))
	for i, date := range data.dates {
		daily = append(daily, dailyStatJSON{
			Date:         formatSummaryDate(date),
			TotalVisits:  data.totalVisits[i],
			UniqueVisits: data.uniqueVisits[i],
			FirstVisits:  data.firstVisits[i],
			ReturnVisits: data.returnVisits[i],
		})
	}

	// Build traffic cards
	tsd := computeTrafficSources(data)
	trafficCards := buildTrafficCardsJSON(tsd)

	// Build heatmap
	hm := heatmapJSON{
		Days:        tsd.Heatmap.Days,
		Slots:       tsd.Heatmap.Slots,
		Intensities: tsd.Heatmap.Intensities,
	}

	result := dashboardJSON{
		SelectedPeriod:      data.selectedPeriod,
		PeriodOptions:       data.periodOptions,
		LiveVisitorCount:    data.liveVisitorCount,
		StatCards:           cards,
		ComparisonRows:      comparisons,
		PreviousPeriodLabel: data.previousPeriodLabel,
		DailyStats:          daily,
		Totals: totalsJSON{
			TotalVisits:  totalVisitors,
			UniqueVisits: totalUniqueVisitors,
			FirstVisits:  totalFirstVisits,
			ReturnVisits: totalReturningVisits,
		},
		TrafficCards:      trafficCards,
		Heatmap:           hm,
		ChartLabels:       data.dates,
		ChartUniqueVisits: data.uniqueVisits,
		ChartTotalVisits:  data.totalVisits,
	}

	return jsonMarshal(result)
}

func buildTrafficCardsJSON(tsd trafficSourcesData) []trafficCardJSON {
	ensure := func(entries []trafficSourceEntry, label string) []trafficSourceEntry {
		if len(entries) == 0 {
			return []trafficSourceEntry{{Label: label, Sessions: "0"}}
		}
		return entries
	}

	return []trafficCardJSON{
		{
			Title: "Referrers", ValueLabel: "Sessions",
			Tabs: []trafficTabJSON{
				{"Referrers", ensure(tsd.Referrers, "(No data)")},
				{"Channels", ensure(tsd.Channels, "(No data)")},
				{"Source", ensure(tsd.Sources, "(No data)")},
				{"Medium", ensure(tsd.Mediums, "(No data)")},
				{"Campaign", ensure(tsd.Campaigns, "(No data)")},
				{"Term", ensure(tsd.Terms, "(No data)")},
			},
		},
		{
			Title: "Pages", ValueLabel: "Sessions",
			Tabs: []trafficTabJSON{
				{"Pages", ensure(tsd.Pages, "(No data)")},
				{"Page Titles", ensure(nil, "(No data)")},
				{"Entry Pages", ensure(tsd.EntryPages, "(No data)")},
				{"Exit Pages", ensure(tsd.ExitPages, "(No data)")},
				{"Hostnames", ensure(nil, "(No data)")},
			},
		},
		{
			Title: "Browsers", ValueLabel: "Sessions",
			Tabs: []trafficTabJSON{
				{"Browsers", ensure(tsd.Browsers, "(No data)")},
				{"Devices", ensure(tsd.Devices, "(No data)")},
				{"Operating Systems", ensure(tsd.OperatingSystems, "(No data)")},
				{"Screen Dimensions", ensure(nil, "(No data)")},
			},
		},
		{
			Title: "Countries", ValueLabel: "Sessions",
			Tabs: []trafficTabJSON{
				{"Countries", ensure(tsd.Countries, "(No data)")},
				{"Regions", ensure(nil, "(No data)")},
				{"Cities", ensure(nil, "(No data)")},
				{"Languages", ensure(tsd.Languages, "(No data)")},
				{"Map", ensure(nil, "(No map data)")},
				{"Timezones", ensure(nil, "(No data)")},
			},
		},
		{
			Title: "Custom Events", ValueLabel: "Count",
			Tabs: []trafficTabJSON{
				{"Custom Events", ensure(tsd.Events, "(No events)")},
				{"Outbound Links", ensure(tsd.OutboundLinks, "(No outbound links)")},
			},
		},
	}
}

func jsonMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `{"error":"failed to marshal JSON"}`
	}
	return string(b)
}

// formatSummaryDate is already defined in table_stats_summary.go
// formatCount, formatFloat, formatFloat2, formatDuration are in traffic_sources_data.go
// changePercentInt, changePercentFloat are in period_comparison.go

var _ = strings.TrimSpace
var _ = shared.PathHome
