package visitoractivity

import (
	"github.com/dracory/statsstore"
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

// CardVisitorActivity builds the visitor activity card with detail modal
func CardVisitorActivity(data ControllerData) hb.TagInterface {
	card := hb.Div().
		Class("card shadow-sm mb-4").
		Child(cardHeader("Visitor Activity")).
		Child(cardBody(data))

	return hb.Div().
		Child(card).
		Child(visitorDetailModal())
}

func cardHeader(title string) hb.TagInterface {
	return hb.Div().
		Class("card-header d-flex justify-content-between align-items-center").
		Child(hb.Heading4().
			Class("card-title mb-0").
			HTML(title)).
		Child(exportDropdown())
}

func cardBody(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card-body").
		Child(hb.Table().
			Class("table table-dark table-striped").
			ID("visitor-activity-table").
			Children([]hb.TagInterface{
				tableHead(),
				tableBody(data.Visitors),
			})).
		Child(pagination(data, data.Page, data.TotalPages))
}

func exportDropdown() hb.TagInterface {
	return hb.Div().
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
					Text("Export to CSV"))))
}

func tableHead() hb.TagInterface {
	return hb.Thead().
		Class("table-dark").
		Children([]hb.TagInterface{
			hb.TR().Children([]hb.TagInterface{
				hb.TH().Text("ID"),
				hb.TH().Text("Path"),
				hb.TH().Text("Timestamp"),
				hb.TH().Text("Duration"),
			}),
		})
}

func tableBody(visitors []statsstore.VisitorInterface) hb.TagInterface {
	return hb.Tbody().Children(lo.Map(visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
		timestamp := formatVisitorTimestamp(visitor.CreatedAt())
		duration := formatVisitDuration(visitor, visitors, index)

		return hb.TR().Children([]hb.TagInterface{
			hb.TD().Text(visitor.ID()),
			hb.TD().HTML(getVisitPageLink(visitor.Path())),
			hb.TD().Text(timestamp),
			hb.TD().Text(duration),
		})
	}))
}
