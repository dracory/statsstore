package home

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/statsstore"
	"github.com/samber/lo"
)

// handleDashboardDataAjax returns daily stats, traffic cards, and heatmap
// in a single response. This eliminates 2 duplicate DB round-trips that
// occurred when daily, traffic, and heatmap endpoints each independently
// queried VisitorList with the same period bounds.
func (c *Controller) handleDashboardDataAjax(w http.ResponseWriter, r *http.Request) string {
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

	// Daily stats
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

	// Traffic cards
	data := ControllerData{
		visitors: visitors,
		ui:       c.ui,
	}
	tsd := computeTrafficSources(data)
	trafficCards := buildTrafficCardsJSON(tsd)

	// Heatmap
	hm := computeHeatmap(visitors)

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
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
		"trafficCards": trafficCards,
		"heatmap": heatmapJSON{
			Days:        hm.Days,
			Slots:       hm.Slots,
			Intensities: hm.Intensities,
		},
	}))

	return ""
}
