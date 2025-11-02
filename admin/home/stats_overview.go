package home

import (
	"fmt"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// statsOverview creates a summary of key statistics
func statsOverview(data ControllerData) hb.TagInterface {
	totalUniqueVisitors := lo.Sum(data.uniqueVisits)
	totalVisitors := lo.Sum(data.totalVisits)
	avgUniqueVisits := float64(totalUniqueVisitors) / float64(len(data.dates))
	avgTotalVisits := float64(totalVisitors) / float64(len(data.dates))

	return hb.Div().
		Class("row g-4 text-center").
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Total Unique Visitors", fmt.Sprintf("%d", totalUniqueVisitors), "bi bi-person", "primary"))).
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Total Visitors", fmt.Sprintf("%d", totalVisitors), "bi bi-people", "success"))).
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Avg. Unique Visitors", fmt.Sprintf("%.2f", avgUniqueVisits), "bi bi-graph-up", "info"))).
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Avg. Total Visitors", fmt.Sprintf("%.2f", avgTotalVisits), "bi bi-bar-chart", "warning")))
}
