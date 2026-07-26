package home

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/dromara/carbon/v2"
)

// handleOverviewAjax returns just the stat cards + live visitor count as JSON.
// This is a lightweight endpoint that only needs aggregate counts, not full visitor rows.
func (c *Controller) handleOverviewAjax(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "application/json")

	periodBounds, err := c.getPeriodBounds(r)
	if err != "" {
		return fmt.Sprintf(`{"error":%q}`, err)
	}

	currentCount, dbErr := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(periodBounds.createdAtGte).
		SetCreatedAtLte(periodBounds.createdAtLte))
	if dbErr != nil {
		return fmt.Sprintf(`{"error":%q}`, dbErr.Error())
	}

	liveGte := carbon.Now(carbon.UTC).SubMinutes(15).ToDateTimeString(carbon.UTC)
	liveCount, _ := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().SetCreatedAtGte(liveGte))

	prevCount, _ := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().
		SetCreatedAtGte(periodBounds.prevCreatedAtGte).
		SetCreatedAtLte(periodBounds.prevCreatedAtLte))

	result := map[string]any{
		"liveVisitorCount": liveCount,
		"totalVisits":      currentCount,
		"previousVisits":   prevCount,
		"selectedPeriod":   periodBounds.selectedPeriod,
		"periodOptions":    periodBounds.periodOptions,
	}

	b, _ := json.Marshal(result)
	return string(b)
}
