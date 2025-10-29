package visitoractivity

import "github.com/gouniverse/hb"

// cardVisitorActivity creates the visitor activity card
func cardVisitorActivity(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header bg-light d-flex justify-content-between align-items-center").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Visitor Activity")).
			Child(hb.Div().
				Class("dropdown").
				Child(hb.Button().
					Class("btn btn-sm btn-outline-secondary dropdown-toggle").
					Attr("type", "button").
					Attr("data-bs-toggle", "dropdown").
					Attr("aria-expanded", "false").
					Text("Export")).
				Child(hb.UL().
					Class("dropdown-menu").
					Child(hb.LI().
						Child(hb.A().
							Class("dropdown-item").
							Href("#").
							Attr("onclick", "exportTableToCSV('visitor-activity-table', 'visitor_activity.csv')").
							Text("Export to CSV")))))).
		Child(hb.Div().
			Class("card-body").
			Child(tableVisitorActivity(data, data.visitors)).
			Child(pagination(data, data.page, data.totalPages)))
}
