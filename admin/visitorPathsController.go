package admin

import (
	"net/http"
	"strings"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
	"github.com/samber/lo"
)

// == CONSTRUCTOR ==============================================================

func visitorPaths(ui ui) PageInterface {
	return &visitorPathsController{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

type visitorPathsController struct {
	ui ui
}

// == PUBLIC METHODS ===========================================================

func (c *visitorPathsController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

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

func (c *visitorPathsController) ToHTML() string {
	return c.ToTag(c.ui.response, c.ui.request).ToHTML()
}

// == PRIVATE METHODS ==========================================================

func (c *visitorPathsController) page(data visitorPathControllerData) *hb.Tag {
	breadcrumbs := breadcrumbs(c.ui.request, []Breadcrumb{
		{
			Name: "Visitor Paths",
			URL:  url(c.ui.request, pathVisitorPaths, nil),
		},
	})

	title := hb.Heading1().
		HTML("Kalleidoscope. Visitor Paths")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader(c.ui)).
		Child(hb.HR()).
		Child(title).
		Child(c.cardVisitorPaths(data))
}

// == PRIVATE METHODS ==========================================================

func (c *visitorPathsController) cardVisitorPaths(data visitorPathControllerData) hb.TagInterface {
	visitors := data.vistorActivities

	return hb.Div().
		Class("card").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading2().
				Class("card-title").
				Style("margin-bottom: 0px;").
				HTML("Visitor Paths"))).
		Child(hb.Div().
			Class("card-body").
			Child(c.tableVisitors(visitors)).
			Child(hb.BR())).
		Child(hb.Div().
			Class("card-footer"))
}

func (c *visitorPathsController) tableVisitors(visitors []vistorActivity) hb.TagInterface {
	cards := hb.Wrap().
		Children(lo.Map(visitors, func(visitor vistorActivity, index int) hb.TagInterface {
			ip := visitor.VisitorIP
			country := visitor.VisitorCountry
			device := visitor.VisitorDevice
			deviceType := visitor.VisitorDeviceType
			browser := visitor.VisitorBrowser
			browserVersion := visitor.VisitorBrowserVersion
			// userAgent := visitor.UserAgent()
			// userLanguage := visitor.UserAcceptLanguage()
			// encoding := visitor.UserAcceptEncoding()
			os := visitor.VisitorOS
			osVersion := visitor.VisitorOSVersion
			referrer := visitor.VisitorPaths[0].Referrer
			if referrer == "" {
				referrer = "(No referring link)"
			}

			countryFlagSrc := "https://flagicons.lipis.dev/flags/4x3/" + strings.ToLower(country) + ".svg"
			countryFlagSrc = strings.ReplaceAll(countryFlagSrc, "un", "xx")
			countryFlag := hb.Img(countryFlagSrc).
				Style("width: 20px;").
				Title("Country code: " + country)

			tableHeader := hb.Table().
				Style("width: 100%;").
				Style("font-size: 14px;").
				Child(hb.TR().
					Child(hb.TD().
						Style("width: 140px;").
						Child(countryFlag).
						Text(" ").
						Text(ip)).
					Child(hb.TD().
						Style("width: 200px;").
						Text(os).
						Text(" ").
						Text(osVersion).
						Text(", ").
						Text(deviceType).
						Text(" ").
						Text(device),
					).
					Child(hb.TD().
						Text(browser).
						Text(" ").
						Text(browserVersion),
					))

			tableBody := hb.Table().
				Style("width: 100%;").
				Child(hb.TR().
					Child(hb.TD().
						Style("width: 100px;").
						Text("")).
					Child(hb.TD().
						Style("width: 100px;").
						Text("")).
					Child(hb.TD().
						Child(hb.Div().Style("color: limegreen;").Text("Referrer: ").Text(referrer)))).
				Children(lo.Map(visitor.VisitorPaths, func(path visitorPath, index int) hb.TagInterface {
					date := path.Date.ToDateString()
					time := path.Date.ToTimeString()
					link := c.ui.websiteUrl + path.Path
					link = strings.ReplaceAll(link, "[GET]", "")

					hyperlink := hb.Hyperlink().
						Target("_blank").
						Href(link).
						Text(path.Path)

					return hb.TR().
						Child(hb.TD().
							Style("width: 100px;").
							Text(date)).
						Child(hb.TD().
							Style("width: 100px;").
							Text(time)).
						Child(hb.TD().
							Child(hb.Div().Style("color: crimson;").Child(hyperlink)))
				}))

			return hb.Div().
				Class("card mb-3").
				Child(hb.Div().
					Class("card-header").
					Child(tableHeader)).
				Child(hb.Div().
					Class("card-body").
					Child(tableBody))
		}))

	return cards
}

func (c *visitorPathsController) prepareData(r *http.Request) (data visitorPathControllerData, errorMessage string) {
	startDate := carbon.Now().SubDays(31).ToDateString()
	endDate := carbon.Now().ToDateString()

	visits, err := c.ui.store.VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		CreatedAtGte: startDate + " 00:00:00",
		CreatedAtLte: endDate + " 23:59:59",
	})

	if err != nil {
		return data, err.Error()
	}

	vistorActivities := []vistorActivity{}

	for _, visit := range visits {
		fingerprint := visit.Fingerprint()

		activity := vistorActivity{
			VisitorCountry:        visit.Country(),
			VisitorDevice:         visit.UserDevice(),
			VisitorDeviceType:     visit.UserDeviceType(),
			VisitorFingerprint:    visit.Fingerprint(),
			VisitorIP:             visit.IpAddress(),
			VisitorOS:             visit.UserOs(),
			VisitorOSVersion:      visit.UserOsVersion(),
			VisitorBrowser:        visit.UserBrowser(),
			VisitorBrowserVersion: visit.UserBrowserVersion(),
			VisitorPaths: []visitorPath{
				{
					Date:     visit.CreatedAtCarbon(),
					Path:     visit.Path(),
					Referrer: visit.UserReferrer(),
				},
			},
		}

		_, index, isFound := lo.FindIndexOf(vistorActivities, func(v vistorActivity) bool {
			return v.VisitorFingerprint == fingerprint
		})

		if isFound {
			vistorActivities[index].VisitorPaths = append(vistorActivities[index].VisitorPaths, visitorPath{
				Date:     visit.CreatedAtCarbon(),
				Path:     visit.Path(),
				Referrer: visit.UserReferrer(),
			})
		} else {
			vistorActivities = append(vistorActivities, activity)
		}
	}

	data.vistorActivities = vistorActivities

	return data, ""
}

type visitorPathControllerData struct {
	vistorActivities []vistorActivity
}

type vistorActivity struct {
	VisitorCountry        string
	VisitorDevice         string
	VisitorDeviceType     string
	VisitorFingerprint    string
	VisitorIP             string
	VisitorOS             string
	VisitorOSVersion      string
	VisitorBrowser        string
	VisitorBrowserVersion string
	VisitorPaths          []visitorPath
}

type visitorPath struct {
	Date     *carbon.Carbon
	Path     string
	Referrer string
}
