package visitorpaths

import (
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// tableVisitorPaths builds the hidden export table used for CSV downloads.
func tableVisitorPaths(data visitorPathsControllerData, ui shared.ControllerOptions) hb.TagInterface {
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
			absolute := fullPathURL(ui, visitor.GetPath())
			browser := strings.TrimSpace(visitor.GetUserBrowser() + " " + visitor.GetUserBrowserVersion())
			countryName := resolvedCountryName(ui, visitor.GetCountry())

			return hb.TR().Children([]hb.TagInterface{
				hb.TD().Text(formatTimestamp(visitor.GetCreatedAt())),
				hb.TD().Text(visitor.GetPath()),
				hb.TD().Text(absolute),
				hb.TD().Text(countryName),
				hb.TD().Text(visitor.GetIpAddress()),
				hb.TD().Text(visitor.GetUserReferrer()),
				hb.TD().Text(sessionLabel(data.Paths, visitor)),
				hb.TD().Text(visitor.GetUserDevice()),
				hb.TD().Text(browser),
			})
		}))

	return hb.Table().
		Class("table table-sm d-none").
		ID("visitor-paths-table").
		Child(head).
		Child(body)
}
