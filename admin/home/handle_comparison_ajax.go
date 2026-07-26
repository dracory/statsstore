package home

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/statsstore"
	"github.com/samber/lo"
)

// handleComparisonAjax returns the period comparison table data as JSON.
// This loads visitors for both current and previous periods to compute extended stats.
func (c *Controller) handleComparisonAjax(w http.ResponseWriter, r *http.Request) string {
	periodBounds, err := c.getPeriodBounds(r)
	if err != "" {
		api.Respond(w, r, api.Error(err))
		return ""
	}

	visitors, dbErr := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(periodBounds.createdAtGte).
		SetCreatedAtLte(periodBounds.createdAtLte))
	if dbErr != nil {
		api.Respond(w, r, api.Error(dbErr.Error()))
		return ""
	}

	prevVisitors, dbErr := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(periodBounds.prevCreatedAtGte).
		SetCreatedAtLte(periodBounds.prevCreatedAtLte))
	if dbErr != nil {
		api.Respond(w, r, api.Error(dbErr.Error()))
		return ""
	}

	currentStats := computePeriodStats(visitors, periodBounds.dateRange)
	prevStats := computePeriodStats(prevVisitors, periodBounds.prevDateRange)

	totalUniqueVisitors := lo.Sum(currentStats.uniqueVisits)
	totalVisitors := lo.Sum(currentStats.totalVisits)
	totalFirstVisits := lo.Sum(currentStats.firstVisits)
	totalReturningVisits := lo.Sum(currentStats.returnVisits)

	ext := computeStatsOverview(visitors)
	prevExt := computeStatsOverview(prevVisitors)

	comparisons := []comparisonRowJSON{
		{"Total Unique Visitors", formatCount(totalUniqueVisitors), formatCount(prevStats.totalUnique), changePercentInt(totalUniqueVisitors, prevStats.totalUnique), false},
		{"Total Visitors", formatCount(totalVisitors), formatCount(prevStats.totalTotal), changePercentInt(totalVisitors, prevStats.totalTotal), false},
		{"First Time Visits", formatCount(totalFirstVisits), formatCount(prevStats.totalFirst), changePercentInt(totalFirstVisits, prevStats.totalFirst), false},
		{"Returning Visits", formatCount(totalReturningVisits), formatCount(prevStats.totalReturning), changePercentInt(totalReturningVisits, prevStats.totalReturning), false},
		{"Bounce Rate", formatFloat2(ext.BounceRateValue) + "%", formatFloat2(prevExt.BounceRateValue) + "%", changePercentFloat(ext.BounceRateValue, prevExt.BounceRateValue), true},
		{"Avg. Visit Duration", formatDuration(ext.SessionDurationSeconds), formatDuration(prevExt.SessionDurationSeconds), changePercentFloat(ext.SessionDurationSeconds, prevExt.SessionDurationSeconds), false},
	}

	statCards := []statCardJSON{
		{"Total Unique Visitors", formatCount(totalUniqueVisitors), "bi bi-person", "primary"},
		{"Total Visitors", formatCount(totalVisitors), "bi bi-people", "success"},
		{"Avg. Daily First Time Visits", formatFloat(float64(totalFirstVisits) / float64(maxInt(len(currentStats.dates), 1))), "bi bi-person-plus", "secondary"},
		{"Avg. Daily Returning Visits", formatFloat(float64(totalReturningVisits) / float64(maxInt(len(currentStats.dates), 1))), "bi bi-person-check", "dark"},
		{"Sessions", ext.Sessions, "bi bi-activity", "primary"},
		{"Pageviews", ext.Pageviews, "bi bi-collection", "success"},
		{"Pages per Session", ext.PagesPerSession, "bi bi-diagram-3", "info"},
		{"Bounce Rate", ext.BounceRate, "bi bi-arrow-repeat", "warning"},
		{"Avg. Visit Duration", ext.SessionDuration, "bi bi-clock-history", "secondary"},
	}

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"statCards":           statCards,
		"comparisonRows":      comparisons,
		"previousPeriodLabel": periodBounds.prevLabel,
	}))

	return ""
}
