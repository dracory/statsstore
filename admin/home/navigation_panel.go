package home

import (
	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

// navigationPanel creates the navigation options panel
func navigationPanel(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card mb-4 shadow-sm").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Navigation"))).
		Child(hb.Div().
			Class("card-body").
			Child(hb.Div().
				Class("row").
				Child(hb.Div().
					Class("col-md-6").
					Child(shared.NavCardUI("Visitor Activity", shared.UrlVisitorActivity(data.Request), "bi bi-activity", "Track visitor interactions"))).
				Child(hb.Div().
					Class("col-md-6").
					Child(shared.NavCardUI("Visitor Paths", shared.UrlVisitorPaths(data.Request), "bi bi-signpost-split", "Analyze visitor navigation paths")))))
}
