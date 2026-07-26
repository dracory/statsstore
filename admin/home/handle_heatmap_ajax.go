package home

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/statsstore"
)

// handleHeatmapAjax returns the weekly heatmap data as JSON.
// This loads visitors for the current period only.
func (c *Controller) handleHeatmapAjax(w http.ResponseWriter, r *http.Request) string {
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

	hm := computeHeatmap(visitors)

	result := heatmapJSON{
		Days:        hm.Days,
		Slots:       hm.Slots,
		Intensities: hm.Intensities,
	}

	b, _ := json.Marshal(result)
	return string(b)
}
