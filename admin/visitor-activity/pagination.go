package visitoractivity

import (
	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/spf13/cast"
)

// pagination creates the pagination component
func pagination(data ControllerData, page int, totalPages int) hb.TagInterface {
	if totalPages <= 1 {
		return hb.Div()
	}

	urlFunc := func(p int) string {
		return shared.UrlVisitorActivity(data.Request, map[string]string{
			"page": cast.ToString(p),
		})
	}

	return shared.PaginationUI(page, totalPages, urlFunc)
}
