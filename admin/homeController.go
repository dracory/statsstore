package admin

import (
	"fmt"
	"net/http"

	// "project/config"
	// "project/internal/layouts"
	// "project/internal/links"

	"github.com/golang-module/carbon/v2"
	// "github.com/gouniverse/bs"

	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONSTRUCTOR ==============================================================

func home(ui ui) PageInterface {
	return &homeController{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	ui ui
}

type homeControllerData struct {
	dates        []string
	uniqueVisits []int64
	totalVisits  []int64
}

func (c *homeController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData()

	c.ui.layout.SetTitle("Dashboard | Kalleidoscope")

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

func (c *homeController) ToHTML() string {
	return c.ToTag(nil, nil).ToHTML()
}

// == PRIVATE METHODS ==========================================================

func (c *homeController) prepareData() (data homeControllerData, errorMessage string) {
	datesInRange := c.datesInRange(carbon.Now().SubDays(31), carbon.Now())

	dates := []string{}
	uniqueVisits := []int64{}
	totalVisits := []int64{}

	for _, date := range datesInRange {
		uniqueVisitorCount, err := c.ui.store.VisitorCount(statsstore.VisitorQueryOptions{
			CreatedAtGte: date + " 00:00:00",
			CreatedAtLte: date + " 23:59:59",
			Distinct:     statsstore.COLUMN_IP_ADDRESS,
		})

		if err != nil {
			return data, err.Error()
		}

		totalVisitorCount, err := c.ui.store.VisitorCount(statsstore.VisitorQueryOptions{
			CreatedAtGte: date + " 00:00:00",
			CreatedAtLte: date + " 23:59:59",
		})

		if err != nil {
			return data, err.Error()
		}

		dates = append(dates, date)
		uniqueVisits = append(uniqueVisits, uniqueVisitorCount)
		totalVisits = append(totalVisits, totalVisitorCount)
	}

	return homeControllerData{
		dates:        dates,
		uniqueVisits: uniqueVisits,
		totalVisits:  totalVisits,
	}, ""
}

func (c *homeController) page(data homeControllerData) hb.TagInterface {
	breadcrumbs := breadcrumbs(c.ui.request, []Breadcrumb{
		{
			Name: "Home",
			URL:  url(c.ui.request, c.ui.homeURL, nil),
		},
		{
			Name: "Kalleidoscope",
			URL:  url(c.ui.request, pathHome, nil),
		},
	})

	title := hb.Heading1().
		HTML("Kalleidoscope. Home")

	options :=
		hb.Section().
			Class("mb-3 mt-3").
			Style("background-color: #f8f9fa;").
			Child(
				hb.UL().
					Class("list-group").
					Child(hb.LI().
						Class("list-group-item").
						Child(hb.A().
							Href(url(c.ui.request, pathVisitorActivity, nil)).
							Text("Visitor Activity")).
						Child(hb.LI().
							Class("list-group-item").
							Child(hb.A().
								Href(url(c.ui.request, pathVisitorPaths, nil)).
								Text("Visitor Paths")))))

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader(c.ui)).
		Child(hb.HR()).
		Child(title).
		Child(options).
		Child(c.cardStatsSummary(data))
}

// == PRIVATE METHODS ==========================================================

func (c *homeController) cardStatsSummary(data homeControllerData) hb.TagInterface {
	return hb.Div().
		Class("card").
		Child(hb.Div().
			Class("card-header").
			Child(hb.Heading2().
				Class("card-title").
				Style("margin-bottom: 0px;").
				HTML("Stats Summary"))).
		Child(hb.Div().
			Class("card-body").
			Child(c.chartStatsSummary(data)).
			Child(hb.BR()).
			Child(c.tableStatsSummary(data)).
			Child(hb.BR())).
		Child(hb.Div().
			Class("card-footer"))
}

func (c *homeController) chartStatsSummary(data homeControllerData) hb.TagInterface {
	labels := data.dates
	uniqueVisitValues := data.uniqueVisits
	totalVisitValues := data.totalVisits

	labelsJSON, err := utils.ToJSON(labels)

	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	uniqueVisitvaluesJSON, err := utils.ToJSON(uniqueVisitValues)

	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	totalVisitValuesJSON, err := utils.ToJSON(totalVisitValues)

	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	script := hb.Script(`
			setTimeout(function () {
				generateVisitorsChart();
			}, 1000);
			function generateVisitorsChart() {
				var visitorData = {
					labels: ` + labelsJSON + `,
					datasets:
							[
								{
									label: "Unique Visitors",
									fillColor: "rgba(172,194,132,0.4)",
									strokeColor: "#ACC26D",
									pointColor: "#fff",
									pointStrokeColor: "#9DB86D",
									data: ` + uniqueVisitvaluesJSON + `
								},
								{
									label: "Total Visitors",
									fillColor: "rgba(91,192,222,0.4)",
									strokeColor: "#5BC0DE",
									pointColor: "#fff",
									pointStrokeColor: "#39B7CD",
									data: ` + totalVisitValuesJSON + `
								}
							]
				};

				var visitorContext = document.getElementById('StatsSummary').getContext('2d');
				
				new Chart(visitorContext, {
					type: 'bar',
					data: visitorData
				});
			}
		`)

	canvas := hb.Canvas().ID("StatsSummary").Style("width:100%;height:300px;")
	return hb.Wrap().
		Child(canvas).
		Child(script)
}

func (c *homeController) tableStatsSummary(data homeControllerData) hb.TagInterface {
	avgUniqueVisits := float64(lo.Sum(data.uniqueVisits)) / float64(len(data.dates))
	avgTotalVisits := float64(lo.Sum(data.totalVisits)) / float64(len(data.dates))

	cardAvgUniqueVisits := hb.Heading3().Text(fmt.Sprintf("Average Unique Visitors: %.2f", avgUniqueVisits))
	cardAvgTotalVisits := hb.Heading3().Text(fmt.Sprintf("Average Total Visitors: %.2f", avgTotalVisits))

	table := hb.Table().
		Class("table table-striped table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().Text("Date"),
					hb.TH().Text("Unique Visitors"),
					hb.TH().Text("Total Visitors"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.dates, func(date string, index int) hb.TagInterface {
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(date),
					hb.TD().Text(cast.ToString(data.uniqueVisits[index])),
					hb.TD().Text(cast.ToString(data.totalVisits[index])),
				})
			})),
			hb.Tfoot().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().Text("Total"),
					hb.TH().Text(cast.ToString(lo.Sum(data.uniqueVisits))),
					hb.TH().Text(cast.ToString(lo.Sum(data.totalVisits))),
				}),
			}),
		})

	return hb.Wrap().
		Child(hb.Div().
			Class("row g-4").
			Child(hb.Div().
				Class("col-6").
				Child(cardAvgUniqueVisits)).
			Child(hb.Div().
				Class("col-6").
				Child(cardAvgTotalVisits))).
		Child(table)
}

func (c *homeController) datesInRange(timeStart, timeEnd carbon.Carbon) []string {
	rangeDates := []string{}

	if timeStart.Lte(timeEnd) {
		rangeDates = append(rangeDates, timeStart.ToDateString())
		for timeStart.Lt(timeEnd) {
			timeStart = timeStart.AddDays(1) // += 86400 // add 24 hours
			rangeDates = append(rangeDates, timeStart.ToDateString())
		}
	}

	return rangeDates
}

// func (c *homeController) visitorsData() (dates []string, uniqueVisits []int64, totalVisits []int64, err error) {
// 	datesInRange := c.datesInRange(carbon.Now().SubDays(31), carbon.Now())

// 	for _, date := range datesInRange {
// 		uniqueVisitorCount, err := c.ui.store.VisitorCount(statsstore.VisitorQueryOptions{
// 			CreatedAtGte: date + " 00:00:00",
// 			CreatedAtLte: date + " 23:59:59",
// 			Distinct:     statsstore.COLUMN_IP_ADDRESS,
// 		})

// 		if err != nil {
// 			return nil, nil, nil, err
// 		}

// 		totalVisitorCount, err := c.ui.store.VisitorCount(statsstore.VisitorQueryOptions{
// 			CreatedAtGte: date + " 00:00:00",
// 			CreatedAtLte: date + " 23:59:59",
// 		})

// 		if err != nil {
// 			return nil, nil, nil, err
// 		}

// 		dates = append(dates, date)
// 		uniqueVisits = append(uniqueVisits, uniqueVisitorCount)
// 		totalVisits = append(totalVisits, totalVisitorCount)
// 	}

// 	return dates, uniqueVisits, totalVisits, nil
// }
