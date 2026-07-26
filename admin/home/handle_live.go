package home

import (
	"net/http"

	"github.com/dromara/carbon/v2"

	"github.com/dracory/statsstore"
)

func (c *Controller) liveVisitorCount(r *http.Request) (int64, error) {
	liveGte := carbon.Now(carbon.UTC).SubMinutes(15).ToDateTimeString(carbon.UTC)
	return c.ui.Store.VisitorCount(r.Context(), statsstore.VisitorQuery().SetCreatedAtGte(liveGte))
}
