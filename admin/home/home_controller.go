package home

import (
	"fmt"
	"net/http"

	"github.com/dromara/carbon/v2"

	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/statsstore/admin/shared"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONSTRUCTOR ==============================================================

// New creates a new home page controller instance
func New(ui shared.UIContext) http.Handler {
	return &Controller{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

// Controller handles the dashboard home page
type Controller struct {
	ui shared.UIContext
}

// ControllerData contains the data needed for the home page
type ControllerData struct {
	Request      *http.Request
	dates        []string
	uniqueVisits []int64
	totalVisits  []int64
}

// ServeHTTP implements the http.Handler interface
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(c.Handle(w, r)))
}

// Handle renders the controller to an HTML tag
func (c *Controller) Handle(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)

	c.ui.GetLayout().SetTitle("Dashboard | Visitor Analytics")

	if errorMessage != "" {
		c.ui.GetLayout().
			SetBody(hb.Div().Class("alert alert-danger").Text(errorMessage).ToHTML())

		return c.ui.GetLayout().Render(w, r)
	}

	// Load required scripts asynchronously
	scripts := []string{
		// Load Chart.js
		`
		if (!window.Chart) {
			const loadChartJS = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
				await new Promise(resolve => script.onload = resolve);
				console.log('Chart.js loaded');
			};
			loadChartJS();
		}
		`,
		// Load HTMX
		`
		if (!window.htmx) {
			const loadHtmx = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = '` + cdn.Htmx_2_0_0() + `';
				await new Promise(resolve => script.onload = resolve);
				console.log('HTMX loaded');
			};
			loadHtmx();
		}
		`,
		// Load SweetAlert2
		`
		if (!window.Swal) {
			const loadSwal = async () => {
				let script = document.createElement('script');
				document.head.appendChild(script);
				script.type = 'text/javascript';
				script.src = '` + cdn.Sweetalert2_11() + `';
				await new Promise(resolve => script.onload = resolve);
				console.log('SweetAlert2 loaded');
			};
			loadSwal();
		}
		`,
	}

	c.ui.GetLayout().SetBody(c.page(data).ToHTML())
	c.ui.GetLayout().SetScripts(scripts)

	return c.ui.GetLayout().Render(c.ui.GetResponse(), c.ui.GetRequest())
}

// == PRIVATE METHODS ==========================================================

// prepareData prepares the data for the home page
func (c *Controller) prepareData(r *http.Request) (data ControllerData, errorMessage string) {
	data.Request = r

	datesInRange := c.datesInRange(carbon.Now().SubDays(31), carbon.Now())

	dates := []string{}
	uniqueVisits := []int64{}
	totalVisits := []int64{}

	for _, date := range datesInRange {
		uniqueVisitorCount, err := c.ui.GetStore().VisitorCount(r.Context(), statsstore.VisitorQueryOptions{
			CreatedAtGte: date + " 00:00:00",
			CreatedAtLte: date + " 23:59:59",
			Distinct:     statsstore.COLUMN_IP_ADDRESS,
		})

		if err != nil {
			return data, err.Error()
		}

		totalVisitorCount, err := c.ui.GetStore().VisitorCount(r.Context(), statsstore.VisitorQueryOptions{
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

	data.dates = dates
	data.uniqueVisits = uniqueVisits
	data.totalVisits = totalVisits

	return data, ""
}

// page builds the main page layout
func (c *Controller) page(data ControllerData) hb.TagInterface {
	breadcrumbs := c.ui.Breadcrumbs([]shared.Breadcrumb{
		{
			Name: "Home",
			URL:  c.ui.URL(c.ui.GetHomeURL(), nil),
		},
		{
			Name: "Visitor Analytics",
			URL:  c.ui.URL(c.ui.GetPathHome(), nil),
		},
	})

	title := hb.Heading1().
		Class("mt-3 mb-4 text-primary").
		HTML("Visitor Analytics Dashboard")

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(c.ui.AdminHeader()).
		Child(hb.HR()).
		Child(title).
		Child(c.navigationPanel(data)).
		Child(c.cardStatsSummary(data))
}

// navigationPanel creates the navigation options panel
func (c *Controller) navigationPanel(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card mb-4 shadow-sm").
		Child(hb.Div().
			Class("card-header bg-light").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Navigation"))).
		Child(hb.Div().
			Class("card-body").
			Child(hb.Div().
				Class("row").
				Child(hb.Div().
					Class("col-md-6").
					Child(shared.NavCardUI("Visitor Activity", shared.UrlVisitorActivity(data.Request), "bi bi-activity", "Track visitor interactions"))).
				Child(hb.Div().
					Class("col-md-6").
					Child(shared.NavCardUI("Visitor Paths", shared.UrlVisitorPaths(data.Request), "bi bi-signpost-split", "Analyze visitor navigation paths")))))
}

// cardStatsSummary creates the stats summary card
func (c *Controller) cardStatsSummary(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header bg-light d-flex justify-content-between align-items-center").
			Child(hb.Heading4().
				Class("card-title mb-0").
				HTML("Stats Summary")).
			Child(hb.Div().
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
							Attr("onclick", "exportTableToCSV('stats-table', 'visitor_stats.csv')").
							Text("Export to CSV"))).
					Child(hb.LI().
						Child(hb.A().
							Class("dropdown-item").
							Href("#").
							Attr("onclick", "exportTableToPDF('stats-table', 'visitor_stats.pdf')").
							Text("Export to PDF")))))).
		Child(hb.Div().
			Class("card-body").
			Child(c.statsOverview(data)).
			Child(hb.HR().Class("my-4")).
			Child(c.chartStatsSummary(data)).
			Child(hb.HR().Class("my-4")).
			Child(c.tableStatsSummary(data)))
}

// statsOverview creates a summary of key statistics
func (c *Controller) statsOverview(data ControllerData) hb.TagInterface {
	totalUniqueVisitors := lo.Sum(data.uniqueVisits)
	totalVisitors := lo.Sum(data.totalVisits)
	avgUniqueVisits := float64(totalUniqueVisitors) / float64(len(data.dates))
	avgTotalVisits := float64(totalVisitors) / float64(len(data.dates))

	return hb.Div().
		Class("row g-4 text-center").
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Total Unique Visitors", fmt.Sprintf("%d", totalUniqueVisitors), "bi bi-person", "primary"))).
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Total Visitors", fmt.Sprintf("%d", totalVisitors), "bi bi-people", "success"))).
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Avg. Unique Visitors", fmt.Sprintf("%.2f", avgUniqueVisits), "bi bi-graph-up", "info"))).
		Child(hb.Div().
			Class("col-md-3").
			Child(shared.StatCardUI("Avg. Total Visitors", fmt.Sprintf("%.2f", avgTotalVisits), "bi bi-bar-chart", "warning")))
}

// chartStatsSummary creates the chart visualization
func (c *Controller) chartStatsSummary(data ControllerData) hb.TagInterface {
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
		document.addEventListener('DOMContentLoaded', function() {
			// Wait for Chart.js to load
			const checkChartInterval = setInterval(function() {
				if (window.Chart) {
					clearInterval(checkChartInterval);
					generateVisitorsChart();
				}
			}, 100);
		});

		function generateVisitorsChart() {
			const ctx = document.getElementById('statsChart').getContext('2d');
			
			const visitorData = {
				labels: ` + labelsJSON + `,
				datasets: [
					{
						label: "Unique Visitors",
						backgroundColor: "rgba(59, 130, 246, 0.5)",
						borderColor: "rgb(59, 130, 246)",
						borderWidth: 2,
						borderRadius: 4,
						data: ` + uniqueVisitvaluesJSON + `
					},
					{
						label: "Total Visitors",
						backgroundColor: "rgba(16, 185, 129, 0.5)",
						borderColor: "rgb(16, 185, 129)",
						borderWidth: 2,
						borderRadius: 4,
						data: ` + totalVisitValuesJSON + `
					}
				]
			};
			
			new Chart(ctx, {
				type: 'bar',
				data: visitorData,
				options: {
					responsive: true,
					maintainAspectRatio: false,
					plugins: {
						legend: {
							position: 'top',
							labels: {
								usePointStyle: true,
								padding: 20
							}
						},
						tooltip: {
							mode: 'index',
							intersect: false,
							padding: 10,
							bodySpacing: 5,
							backgroundColor: 'rgba(0, 0, 0, 0.8)'
						}
					},
					scales: {
						y: {
							beginAtZero: true,
							grid: {
								drawBorder: false
							}
						},
						x: {
							grid: {
								display: false
							}
						}
					}
				}
			});

			// Add chart toggle functionality
			document.getElementById('toggleChartType').addEventListener('click', function() {
				const chart = Chart.getChart('statsChart');
				if (!chart) return;
				
				const newType = chart.config.type === 'bar' ? 'line' : 'bar';
				chart.config.type = newType;
				
				if (newType === 'line') {
					chart.data.datasets.forEach(dataset => {
						dataset.backgroundColor = dataset.borderColor;
						dataset.pointBackgroundColor = dataset.borderColor;
						dataset.pointRadius = 4;
						dataset.tension = 0.2;
					});
					this.innerHTML = '<i class="bi bi-bar-chart"></i> Switch to Bar';
				} else {
					chart.data.datasets.forEach(dataset => {
						dataset.backgroundColor = dataset.borderColor.replace('rgb', 'rgba').replace(')', ', 0.5)');
						dataset.pointRadius = 0;
					});
					this.innerHTML = '<i class="bi bi-graph-up"></i> Switch to Line';
				}
				
				chart.update();
			});
		}

		// Export functions
		function exportTableToCSV(tableId, filename) {
			const table = document.getElementById(tableId);
			if (!table) return;
			
			let csv = [];
			const rows = table.querySelectorAll('tr');
			
			for (let i = 0; i < rows.length; i++) {
				const row = [], cols = rows[i].querySelectorAll('td, th');
				
				for (let j = 0; j < cols.length; j++) {
					row.push('"' + cols[j].innerText.replace(/"/g, '""') + '"');
				}
				
				csv.push(row.join(','));
			}
			
			const csvContent = csv.join('\\n');
			const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
			const link = document.createElement('a');
			
			link.href = URL.createObjectURL(blob);
			link.setAttribute('download', filename);
			link.click();
		}

		function exportTableToPDF(tableId, filename) {
			// This is a placeholder - in a real implementation you would use a library like jsPDF
			alert('PDF export would be implemented with jsPDF or similar library');
		}
	`)

	return hb.Div().
		Class("chart-container").
		Child(hb.Div().
			Class("d-flex justify-content-between align-items-center mb-3").
			Child(hb.Heading5().
				Class("mb-0").
				Text("Visitor Statistics")).
			Child(hb.Button().
				ID("toggleChartType").
				Class("btn btn-sm btn-outline-primary").
				Attr("type", "button").
				Child(hb.I().
					Class("bi bi-graph-up").
					Attr("style", "margin-right: 5px")).
				Text("Switch to Line"))).
		Child(hb.Div().
			Class("position-relative").
			Style("height: 350px;").
			Child(hb.Canvas().
				ID("statsChart").
				Attr("width", "100%").
				Attr("height", "350"))).
		Child(script)
}

// tableStatsSummary creates the data table
func (c *Controller) tableStatsSummary(data ControllerData) hb.TagInterface {
	table := hb.Table().
		ID("stats-table").
		Class("table table-striped table-hover table-sm").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-light").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Date"),
						hb.TH().
							Class("text-end").
							Text("Unique Visitors"),
						hb.TH().
							Class("text-end").
							Text("Total Visitors"),
					}),
				}),
			hb.Tbody().Children(lo.Map(data.dates, func(date string, index int) hb.TagInterface {
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(date),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.uniqueVisits[index])),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.totalVisits[index])),
				})
			})),
			hb.Tfoot().
				Class("table-light fw-bold").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Total"),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.uniqueVisits))),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.totalVisits))),
					}),
				}),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}

// datesInRange returns an array of dates between the start and end dates
func (c *Controller) datesInRange(timeStart, timeEnd *carbon.Carbon) []string {
	rangeDates := []string{}

	if timeStart.Lte(timeEnd) {
		rangeDates = append(rangeDates, timeStart.ToDateString())
		for timeStart.Lt(timeEnd) {
			timeStart = timeStart.AddDays(1)
			rangeDates = append(rangeDates, timeStart.ToDateString())
		}
	}

	return rangeDates
}
