package home

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samber/lo"
)

// handleComparisonAjax returns the period comparison table data as JSON.
// This loads visitors for both current and previous periods to compute extended stats.
func (c *Controller) handleComparisonAjax(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "application/json")

	data, errorMessage := c.prepareData(r)
	if errorMessage != "" {
		return fmt.Sprintf(`{"error":%q}`, errorMessage)
	}

	totalUniqueVisitors := lo.Sum(data.uniqueVisits)
	totalVisitors := lo.Sum(data.totalVisits)
	totalFirstVisits := lo.Sum(data.firstVisits)
	totalReturningVisits := lo.Sum(data.returnVisits)

	ext := data.currentStats

	comparisons := []comparisonRowJSON{
		{"Total Unique Visitors", formatCount(totalUniqueVisitors), formatCount(data.previousPeriodUnique), changePercentInt(totalUniqueVisitors, data.previousPeriodUnique), false},
		{"Total Visitors", formatCount(totalVisitors), formatCount(data.previousPeriodTotal), changePercentInt(totalVisitors, data.previousPeriodTotal), false},
		{"First Time Visits", formatCount(totalFirstVisits), formatCount(data.previousPeriodFirst), changePercentInt(totalFirstVisits, data.previousPeriodFirst), false},
		{"Returning Visits", formatCount(totalReturningVisits), formatCount(data.previousPeriodReturning), changePercentInt(totalReturningVisits, data.previousPeriodReturning), false},
		{"Bounce Rate", formatFloat2(ext.BounceRateValue) + "%", formatFloat2(data.previousStats.BounceRateValue) + "%", changePercentFloat(ext.BounceRateValue, data.previousStats.BounceRateValue), true},
		{"Avg. Visit Duration", formatDuration(ext.SessionDurationSeconds), formatDuration(data.previousStats.SessionDurationSeconds), changePercentFloat(ext.SessionDurationSeconds, data.previousStats.SessionDurationSeconds), false},
	}

	statCards := []statCardJSON{
		{"Total Unique Visitors", formatCount(totalUniqueVisitors), "bi bi-person", "primary"},
		{"Total Visitors", formatCount(totalVisitors), "bi bi-people", "success"},
		{"Avg. Daily First Time Visits", formatFloat(float64(totalFirstVisits) / float64(maxInt(len(data.dates), 1))), "bi bi-person-plus", "secondary"},
		{"Avg. Daily Returning Visits", formatFloat(float64(totalReturningVisits) / float64(maxInt(len(data.dates), 1))), "bi bi-person-check", "dark"},
		{"Sessions", ext.Sessions, "bi bi-activity", "primary"},
		{"Pageviews", ext.Pageviews, "bi bi-collection", "success"},
		{"Pages per Session", ext.PagesPerSession, "bi bi-diagram-3", "info"},
		{"Bounce Rate", ext.BounceRate, "bi bi-arrow-repeat", "warning"},
		{"Avg. Visit Duration", ext.SessionDuration, "bi bi-clock-history", "secondary"},
	}

	result := map[string]any{
		"statCards":           statCards,
		"comparisonRows":      comparisons,
		"previousPeriodLabel": data.previousPeriodLabel,
	}

	b, _ := json.Marshal(result)
	return string(b)
}
