package visitorpaths

import (
	"net/http"

	"github.com/dracory/statsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/spf13/cast"
)

// pagination creates the pagination component
func pagination(r *http.Request, page int, totalPages int) hb.TagInterface {
	if totalPages <= 1 {
		return hb.Div()
	}

	urlFunc := func(p int) string {
		return shared.UrlVisitorPaths(r, map[string]string{
			"page": cast.ToString(p),
		})
	}

	return shared.PaginationUI(page, totalPages, urlFunc)
}
