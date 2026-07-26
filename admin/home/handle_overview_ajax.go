package home

import (
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/statsstore"
	"github.com/dromara/carbon/v2"
)

// handleOverviewAjax returns just the stat cards + live visitor count as JSON.
// This is a lightweight endpoint that only needs aggregate counts, not full visitor rows.
func (c *Controller) handleOverviewAjax(w http.ResponseWriter, r *http.Request) string {
	periodBounds, err := c.getPeriodBounds(r)
	if err != "" {
		api.Respond(w, r, api.Error(err))
		return ""
	}

	liveGte := carbon.Now(carbon.UTC).SubMinutes(15).ToDateTimeString(carbon.UTC)
	liveCount, _ := c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().SetCreatedAtGte(liveGte))

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"liveVisitorCount": liveCount,
		"selectedPeriod":   periodBounds.selectedPeriod,
		"periodOptions":    periodBounds.periodOptions,
	}))

	return ""
}
