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
	totalFirstVisits := lo.Sum(data.firstVisits)
	totalReturningVisits := lo.Sum(data.returnVisits)

	days := len(data.dates)
	if days == 0 {
		days = 1
	}

	avgUniqueVisits := float64(totalUniqueVisitors) / float64(days)
	avgTotalVisits := float64(totalVisitors) / float64(days)
	avgFirstVisits := float64(totalFirstVisits) / float64(days)
	avgReturningVisits := float64(totalReturningVisits) / float64(days)

	return hb.Div().
		Class("row row-cols-1 row-cols-sm-2 row-cols-lg-3 row-cols-xl-6 g-4 text-center").
		Child(hb.Div().
			Class("col").
			Child(shared.StatCardUI("Total Unique Visitors", fmt.Sprintf("%d", totalUniqueVisitors), "bi bi-person", "primary"))).
		Child(hb.Div().
			Class("col").
			Child(shared.StatCardUI("Total Visitors", fmt.Sprintf("%d", totalVisitors), "bi bi-people", "success"))).
		Child(hb.Div().
			Class("col").
			Child(shared.StatCardUI("Avg. Unique Visitors", fmt.Sprintf("%.2f", avgUniqueVisits), "bi bi-graph-up", "info"))).
		Child(hb.Div().
			Class("col").
			Child(shared.StatCardUI("Avg. Total Visitors", fmt.Sprintf("%.2f", avgTotalVisits), "bi bi-bar-chart", "warning"))).
		Child(hb.Div().
			Class("col").
			Child(shared.StatCardUI("Avg. Daily First Time Visits", fmt.Sprintf("%.2f", avgFirstVisits), "bi bi-person-plus", "secondary"))).
		Child(hb.Div().
			Class("col").
			Child(shared.StatCardUI("Avg. Daily Returning Visits", fmt.Sprintf("%.2f", avgReturningVisits), "bi bi-person-check", "dark")))
}
