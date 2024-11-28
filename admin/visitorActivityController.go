package admin

import (
	"net/http"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
	"github.com/samber/lo"
)

// == CONSTRUCTOR ==============================================================

func visitorActivity(ui ui) PageInterface {
	return &visitorActivityController{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

type visitorActivityController struct {
	ui ui
}

type visitorActivityControllerData struct {
	visitors []statsstore.VisitorInterface
}

// == PUBLIC METHODS ===========================================================

func (c *visitorActivityController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData()

	c.ui.layout.SetTitle("Visitor Activity | Kalleidoscope")

	if errorMessage != "" {
		c.ui.layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.ui.layout.Render(w, r))
	}

	htmxScript := `setTimeout(() => async function() {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
			await script.onload
		}
	}, 1000);`

	swalScript := `setTimeout(() => async function() {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
			await script.onload
		}
	}, 1000);`

	// cdn.Jquery_3_7_1(),
	// // `https://cdnjs.cloudflare.com/ajax/libs/Chart.js/1.0.2/Chart.min.js`,
	// `https://cdn.jsdelivr.net/npm/chart.js`,

	c.ui.layout.SetBody(c.page(data).ToHTML())
	c.ui.layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.ui.layout.Render(w, r))
}

func (c *visitorActivityController) ToHTML() string {
	return c.ToTag(nil, nil).ToHTML()
}

// == PRIVATE METHODS ==========================================================

func (c *visitorActivityController) page(data visitorActivityControllerData) *hb.Tag {
	breadcrumbs := breadcrumbs(c.ui.request, []Breadcrumb{
		{
			Name: "Visitor Activity",
			URL:  url(c.ui.request, pathVisitorActivity, nil),
		},
	})

	title := hb.Heading1().
		HTML("Kalleidoscope. Visitor Activity")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader(c.ui)).
		Child(hb.HR()).
		Child(title).
		Child(c.cardVisitorActivity(data))
}

// == PRIVATE METHODS ==========================================================

func (c *visitorActivityController) cardVisitorActivity(data visitorActivityControllerData) hb.TagInterface {
	return hb.Div().
		Class("card").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading2().
				Class("card-title").
				Style("margin-bottom: 0px;").
				HTML("Visitor Activity"))).
		Child(hb.Div().
			Class("card-body").
			Child(c.tableVisitors(data.visitors)).
			Child(hb.BR())).
		Child(hb.Div().
			Class("card-footer"))
}

func (c *visitorActivityController) tableVisitors(visitors []statsstore.VisitorInterface) hb.TagInterface {
	// Page Views:
	// 1
	// Latest Page View:
	// 27 Oct 2024 00:28:02
	// Resolution:
	// 414x736
	// System:
	// Chrome for Android/Android
	// Vivo Y5s

	// Total Sessions:
	// 1
	// Location:
	// [China] China
	// ISP / IP Address:
	// China Mobile (39.173.105.141)
	// Referring URL:
	// (No referring link)
	// Visit Page:
	//  https://lesichkov.co.uk/

	table := hb.Table().
		Class("table table-striped table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					// hb.TH().Text("Date"),
					// hb.TH().Text("Unique Visitors"),
					// hb.TH().Text("Total Visitors"),
				}),
			}),
			hb.Tbody().Children(lo.Map(visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
				date := visitor.CreatedAtCarbon().ToDateString()
				ip := visitor.IpAddress()
				country := visitor.Country()
				device := visitor.UserDevice()
				deviceType := visitor.UserDeviceType()
				browser := visitor.UserBrowser()
				browserVersion := visitor.UserBrowserVersion()
				userAgent := visitor.UserAgent()
				userLanguage := visitor.UserAcceptLanguage()
				encoding := visitor.UserAcceptEncoding()
				os := visitor.UserOs()
				referrer := visitor.UserReferrer()
				path := visitor.Path()
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Child(hb.NewSpan().Text("Path: ")).Text(path),
					hb.TD().Child(hb.NewSpan().Text("Date: ")).Text(date),
					hb.TD().Child(hb.NewSpan().Text("IP: ")).Text(ip),
					hb.TD().Child(hb.NewSpan().Text("Country: ")).Text(country),
					hb.TD().Child(hb.NewSpan().Text("Device: ")).Text(device),
					hb.TD().Child(hb.NewSpan().Text("Device Type: ")).Text(deviceType),
					hb.TD().Child(hb.NewSpan().Text("Browser: ")).Text(browser),
					hb.TD().Child(hb.NewSpan().Text("Browser Version: ")).Text(browserVersion),
					hb.TD().Child(hb.NewSpan().Text("OS: ")).Text(os),
					hb.TD().Child(hb.NewSpan().Text("Referrer: ")).Text(referrer),
					hb.TD().Child(hb.NewSpan().Text("User Agent: ")).Text(userAgent),
					hb.TD().Child(hb.NewSpan().Text("User Language: ")).Text(userLanguage),
					hb.TD().Child(hb.NewSpan().Text("User Encoding: ")).Text(encoding),
				})
			})),
			hb.Tfoot().Children([]hb.TagInterface{
				// hb.TR().Children([]hb.TagInterface{
				// 	hb.TH().Text("Total"),
				// 	hb.TH().Text(cast.ToString(lo.Sum(uniqueVisits))),
				// 	hb.TH().Text(cast.ToString(lo.Sum(totalVisits))),
				// }),
			}),
		})

	return hb.Wrap().
		Child(table)
}

// func (c *visitorActivityController) datesInRange(timeStart, timeEnd carbon.Carbon) []string {
// 	rangeDates := []string{}

// 	if timeStart.Lte(timeEnd) {
// 		rangeDates = append(rangeDates, timeStart.ToDateString())
// 		for timeStart.Lt(timeEnd) {
// 			timeStart = timeStart.AddDays(1) // += 86400 // add 24 hours
// 			rangeDates = append(rangeDates, timeStart.ToDateString())
// 		}
// 	}

// 	return rangeDates
// }

func (c *visitorActivityController) prepareData() (data visitorActivityControllerData, errorMessage string) {
	startDate := carbon.Now().SubDays(31).ToDateString()
	endDate := carbon.Now().ToDateString()

	visitors, err := c.ui.store.VisitorList(statsstore.VisitorQueryOptions{
		CreatedAtGte: startDate + " 00:00:00",
		CreatedAtLte: endDate + " 23:59:59",
	})

	if err != nil {
		return data, err.Error()
	}

	data.visitors = visitors

	return data, ""
}
