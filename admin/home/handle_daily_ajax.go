package home

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/samber/lo"
)

// handleDailyAjax returns the daily stats table + chart data as JSON.
// This loads visitors for the current period only.
func (c *Controller) handleDailyAjax(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "application/json")

	periodBounds, err := c.getPeriodBounds(r)
	if err != "" {
		return fmt.Sprintf(`{"error":%q}`, err)
	}

	visitors, dbErr := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(periodBounds.createdAtGte).
		SetCreatedAtLte(periodBounds.createdAtLte))
	if dbErr != nil {
		return fmt.Sprintf(`{"error":%q}`, dbErr.Error())
	}

	currentStats := computePeriodStats(visitors, periodBounds.dateRange)

	daily := make([]dailyStatJSON, 0, len(currentStats.dates))
	for i, date := range currentStats.dates {
		daily = append(daily, dailyStatJSON{
			Date:         formatSummaryDate(date),
			TotalVisits:  currentStats.totalVisits[i],
			UniqueVisits: currentStats.uniqueVisits[i],
			FirstVisits:  currentStats.firstVisits[i],
			ReturnVisits: currentStats.returnVisits[i],
		})
	}

	result := map[string]any{
		"dailyStats":        daily,
		"chartLabels":       currentStats.dates,
		"chartUniqueVisits": currentStats.uniqueVisits,
		"chartTotalVisits":  currentStats.totalVisits,
		"totals": totalsJSON{
			TotalVisits:  lo.Sum(currentStats.totalVisits),
			UniqueVisits: lo.Sum(currentStats.uniqueVisits),
			FirstVisits:  lo.Sum(currentStats.firstVisits),
			ReturnVisits: lo.Sum(currentStats.returnVisits),
		},
	}

	b, _ := json.Marshal(result)
	return string(b)
}
