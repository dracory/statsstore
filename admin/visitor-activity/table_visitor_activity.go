package visitoractivity

import (
	"time"

	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// tableVisitorActivity creates the visitor activity table
func tableVisitorActivity(data ControllerData, visitors []statsstore.VisitorInterface) hb.TagInterface {
	table := hb.Table().
		ID("visitor-activity-table").
		Class("table table-striped table-hover").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-light").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("ID"),
						hb.TH().Text("IP Address"),
						hb.TH().Text("Path"),
						hb.TH().Text("Referrer"),
						hb.TH().Text("User Agent"),
						hb.TH().Text("Created At"),
						hb.TH().Text("Actions"),
					}),
				}),
			hb.Tbody().Children(lo.Map(visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
				// Format the created at date to match the UI
				createdAt := visitor.CreatedAt()
				if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
					createdAt = t.Format("2006-01-02 15:04:05 -0700 UTC")
				}

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(cast.ToString(visitor.ID())),
					hb.TD().Text(visitor.IpAddress()),
					hb.TD().Text(shared.StrTruncate(visitor.Path(), 30)),
					hb.TD().Text(shared.StrTruncate(visitor.UserReferrer(), 30)),
					hb.TD().Text(shared.StrTruncate(visitor.UserAgent(), 30)),
					hb.TD().Text(createdAt),
					hb.TD().Child(hb.A().
						Class("btn btn-sm btn-outline-primary").
						Attr("data-bs-toggle", "tooltip").
						Attr("title", "View details").
						Href(shared.UrlVisitorActivity(data.Request, map[string]string{"path": "/admin/visitor-activity/" + cast.ToString(visitor.ID())})).
						Child(hb.I().Class("bi bi-eye"))),
				})
			})),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}
