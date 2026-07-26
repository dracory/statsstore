package home

import (
	"net/http"
	"strconv"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

// handleExport exports the dashboard data as CSV.
func (c *Controller) handleExport(w http.ResponseWriter, r *http.Request) string {
	periodBounds, err := c.getPeriodBounds(r)
	if err != "" {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	visitors, dbErr := c.ui.Store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(periodBounds.createdAtGte).
		SetCreatedAtLte(periodBounds.createdAtLte))
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return dbErr.Error()
	}

	stats := computePeriodStats(visitors, periodBounds.dateRange)

	headers := []string{
		"Date",
		"Page Views",
		"Unique Visits",
		"First Time Visits",
		"Returning Visits",
	}

	rows := make([][]string, 0, len(stats.dates))
	for i, date := range stats.dates {
		rows = append(rows, []string{
			formatSummaryDate(date),
			strconv.FormatInt(stats.totalVisits[i], 10),
			strconv.FormatInt(stats.uniqueVisits[i], 10),
			strconv.FormatInt(stats.firstVisits[i], 10),
			strconv.FormatInt(stats.returnVisits[i], 10),
		})
	}

	return shared.ExportCSV(w, shared.ExportFilename("visitor-stats"), headers, rows)
}
