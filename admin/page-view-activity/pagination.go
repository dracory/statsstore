package pageviewactivity

import (
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/spf13/cast"
)

// pagination creates the pagination component for the page view activity screen.
func pagination(r *http.Request, page int, totalPages int) hb.TagInterface {
	if totalPages <= 1 {
		return hb.Div()
	}

	urlFunc := func(p int) string {
		return shared.UrlPageViewActivity(r, map[string]string{
			"page": cast.ToString(p),
		})
	}

	return shared.PaginationUI(page, totalPages, urlFunc)
}
