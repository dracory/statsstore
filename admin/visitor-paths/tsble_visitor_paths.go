package visitorpaths

import (
	"net/http"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

// tableVisitorPaths creates the visitor paths table
func tableVisitorPaths(r *http.Request, paths []statsstore.VisitorInterface) hb.TagInterface {
	table := hb.Table().
		ID("visitor-paths-table").
		Class("table table-striped table-hover").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-light").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("URL"),
						hb.TH().Class("text-end").Text("Visit Count"),
						hb.TH().Text("Last Visit"),
						hb.TH().Text("Actions"),
					}),
				}),
			hb.Tbody().Children(lo.Map(paths, func(path statsstore.VisitorInterface, index int) hb.TagInterface {
				// For now, we'll just show the path and created date
				// In a real implementation, we would need to add count functionality to the statsstore
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(shared.StrTruncate(path.Path(), 50)),
					hb.TD().Class("text-end").Text("1"), // Placeholder for count
					hb.TD().Text(path.CreatedAt()),
					hb.TD().Child(hb.A().
						Class("btn btn-sm btn-outline-primary").
						Attr("data-bs-toggle", "tooltip").
						Attr("title", "View visitors for this path").
						Href(shared.UrlVisitorActivity(r, map[string]string{
							"path": path.Path(),
						})).
						Child(hb.I().Class("bi bi-eye"))),
				})
			})),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}
