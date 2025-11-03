package visitorpaths

import (
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// tableVisitorPaths builds the hidden export table used for CSV downloads.
func tableVisitorPaths(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	head := hb.Thead().
		Child(hb.TR().Children([]hb.TagInterface{
			hb.TH().Text("Visit Time"),
			hb.TH().Text("Path"),
			hb.TH().Text("Absolute URL"),
			hb.TH().Text("Country"),
			hb.TH().Text("IP Address"),
			hb.TH().Text("Referrer"),
			hb.TH().Text("Session"),
			hb.TH().Text("Device"),
			hb.TH().Text("Browser"),
		}))

	body := hb.Tbody().
		Children(lo.Map(data.Paths, func(visitor statsstore.VisitorInterface, _ int) hb.TagInterface {
			absolute := fullPathURL(ui, visitor.Path())
			browser := strings.TrimSpace(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())
			countryName := resolvedCountryName(ui, visitor.Country())

			return hb.TR().Children([]hb.TagInterface{
				hb.TD().Text(formatTimestamp(visitor.CreatedAt())),
				hb.TD().Text(visitor.Path()),
				hb.TD().Text(absolute),
				hb.TD().Text(countryName),
				hb.TD().Text(visitor.IpAddress()),
				hb.TD().Text(visitor.UserReferrer()),
				hb.TD().Text(sessionLabel(data.Paths, visitor)),
				hb.TD().Text(visitor.UserDevice()),
				hb.TD().Text(browser),
			})
		}))

	return hb.Table().
		Class("table table-sm d-none").
		ID("visitor-paths-table").
		Child(head).
		Child(body)
}
