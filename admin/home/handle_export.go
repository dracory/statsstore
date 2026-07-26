package home

import (
	"net/http"
	"strconv"

	"github.com/dracory/statsstore/admin/shared"
)

// handleExport exports the dashboard data as CSV.
func (c *Controller) handleExport(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)
	if errorMessage != "" {
		w.WriteHeader(http.StatusInternalServerError)
		return errorMessage
	}
	return c.exportCSV(w, data)
}

func (c *Controller) exportCSV(w http.ResponseWriter, data ControllerData) string {
	headers := []string{
		"Date",
		"Page Views",
		"Unique Visits",
		"First Time Visits",
		"Returning Visits",
	}

	rows := make([][]string, 0, len(data.dates))
	for i, date := range data.dates {
		rows = append(rows, []string{
			formatSummaryDate(date),
			strconv.FormatInt(data.totalVisits[i], 10),
			strconv.FormatInt(data.uniqueVisits[i], 10),
			strconv.FormatInt(data.firstVisits[i], 10),
			strconv.FormatInt(data.returnVisits[i], 10),
		})
	}

	return shared.ExportCSV(w, shared.ExportFilename("visitor-stats"), headers, rows)
}
