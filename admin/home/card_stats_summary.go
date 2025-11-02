package home

import "github.com/dracory/hb"

// cardStatsSummary creates the stats summary card
func cardStatsSummary(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header d-flex flex-column flex-md-row gap-2 justify-content-between align-items-md-center").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Stats Summary")).
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
							Attr("onclick", "exportTableToCSV('stats-table', 'visitor_stats.csv')").
							Text("Export to CSV"))).
					Child(hb.LI().
						Child(hb.A().
							Class("dropdown-item").
							Href("#").
							Attr("onclick", "exportTableToPDF('stats-table', 'visitor_stats.pdf')").
							Text("Export to PDF")))))).
		Child(hb.Div().
			Class("card-body").
			Child(statsOverview(data)).
			Child(hb.HR().Class("my-4")).
			Child(chartStatsSummary(data)).
			Child(hb.HR().Class("my-4")).
			Child(tableStatsSummary(data)))
}
