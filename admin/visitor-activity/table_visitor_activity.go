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
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Device"),
						hb.TH().Text("Location"),
						hb.TH().Text("ISP / IP"),
						hb.TH().Text("Referrer"),
						hb.TH().Text("Pages"),
						hb.TH().Text("Exit Time"),
						hb.TH().Text("Actions"),
					}),
				}),
			hb.Tbody().Children(lo.Map(visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
				// Format the created at date
				createdAt := visitor.CreatedAt()
				formattedDate := ""
				if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
					formattedDate = t.Format("2 Jan 2006 15:04:05")
				}

				// Build device info
				deviceInfo := hb.Div().
					Child(hb.Div().Class("d-flex align-items-center gap-2").
						Child(deviceIcon(visitor)).
						Child(osIcon(visitor))).
					Child(hb.Small().Class("text-muted d-block").Text(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())).
					Child(hb.Small().Class("text-muted d-block").Text(visitor.UserOs()))

				// Location with flag
				location := visitor.Country()
				if location == "" {
					location = "Unknown"
				}

				// Referrer display
				referrer := visitor.UserReferrer()
				referrerDisplay := "(No referring link)"
				if referrer != "" {
					referrerDisplay = shared.StrTruncate(referrer, 40)
				}

				// Page views (for now showing path, could be enhanced to count)
				pageViews := "1"

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Child(deviceInfo),
					hb.TD().Text(location),
					hb.TD().
						Child(hb.Div().Text(shared.StrTruncate(visitor.IpAddress(), 20))).
						Child(hb.Small().Class("text-muted d-block").Text("(" + visitor.IpAddress() + ")")),
					hb.TD().
						Child(hb.Div().Class("text-success").Text(referrerDisplay)),
					hb.TD().Text(pageViews),
					hb.TD().Text(formattedDate),
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
