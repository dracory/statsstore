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
	tableHeaders := []string{"Device", "Location", "ISP / IP", "Referrer", "Pages", "Exit Time", "Actions"}

	headers := lo.Map(tableHeaders, func(title string, _ int) hb.TagInterface {
		return hb.TH().Text(title)
	})

	headRow := hb.TR().Children(headers)
	tableHead := hb.Thead().Child(headRow)

	bodyRows := lo.Map(visitors, func(visitor statsstore.VisitorInterface, _ int) hb.TagInterface {
		createdAt := visitor.CreatedAt()
		formattedDate := ""
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			formattedDate = t.Format("2 Jan 2006 15:04:05")
		}

		deviceSummary := hb.Div().
			Child(hb.Div().Class("d-flex align-items-center gap-2").
				Child(deviceIcon(visitor)).
				Child(osIcon(visitor))).
			Child(hb.Small().Class("text-muted d-block").Text(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())).
			Child(hb.Small().Class("text-muted d-block").Text(visitor.UserOs()))

		location := visitor.Country()
		if location == "" {
			location = "Unknown"
		}
		locationCell := hb.TD().Text(location)

		ipValue := visitor.IpAddress()
		ipCell := hb.TD().
			Child(hb.Div().Text(shared.StrTruncate(ipValue, 20))).
			Child(hb.Small().Class("text-muted d-block").Text("(" + ipValue + ")"))

		referrer := visitor.UserReferrer()
		referrerDisplay := "(No referring link)"
		if referrer != "" {
			referrerDisplay = shared.StrTruncate(referrer, 40)
		}
		referrerCell := hb.TD().Child(hb.Div().Class("text-success").Text(referrerDisplay))

		pageViewsCell := hb.TD().Text("1")
		timeCell := hb.TD().Text(formattedDate)

		viewDetailsButton := hb.Button().
			Class("btn btn-sm btn-outline-primary").
			Attr("type", "button").
			Attr("data-bs-toggle", "modal").
			Attr("data-bs-target", "#visitorDetailModal").
			Attr("hx-get", shared.UrlVisitorActivity(data.Request, map[string]string{"visitor_id": cast.ToString(visitor.ID()), "modal": "1"})).
			Attr("hx-target", "#visitorDetailModalContent").
			Attr("hx-swap", "innerHTML").
			Attr("title", "View details").
			Child(hb.I().Class("bi bi-eye"))

		actionCell := hb.TD().Child(viewDetailsButton)

		return hb.TR().Children([]hb.TagInterface{
			hb.TD().Child(deviceSummary),
			locationCell,
			ipCell,
			referrerCell,
			pageViewsCell,
			timeCell,
			actionCell,
		})
	})

	body := hb.Tbody().Children(bodyRows)

	table := hb.Table().
		ID("visitor-activity-table").
		Class("table table-striped table-hover").
		Children([]hb.TagInterface{tableHead, body})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}
