package home

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/statsstore"
)

// handleTrafficAjax returns the traffic source breakdown cards as JSON.
// This loads visitors for the current period only.
func (c *Controller) handleTrafficAjax(w http.ResponseWriter, r *http.Request) string {
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

	data := ControllerData{
		visitors: visitors,
		ui:       c.ui,
	}

	tsd := computeTrafficSources(data)
	trafficCards := buildTrafficCardsJSON(tsd)

	result := map[string]any{
		"trafficCards": trafficCards,
	}

	b, _ := json.Marshal(result)
	return string(b)
}
