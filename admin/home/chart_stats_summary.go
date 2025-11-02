package home

import (
	"encoding/json"

	"github.com/dracory/hb"
)

// chartStatsSummary creates the chart visualization
func chartStatsSummary(data ControllerData) hb.TagInterface {
	labels := data.dates
	uniqueVisitValues := data.uniqueVisits
	totalVisitValues := data.totalVisits

	labelsJSON, err := json.Marshal(labels)
	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	uniqueVisitvaluesJSON, err := json.Marshal(uniqueVisitValues)
	if err != nil {
		return hb.Div().Class("alert alert-danger").Text(err.Error())
	}

	totalVisitValuesJSON, err := json.Marshal(totalVisitValues)
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
				labels: ` + string(labelsJSON) + `,
				datasets: [
					{
						label: "Unique Visitors",
						backgroundColor: "rgba(59, 130, 246, 0.5)",
						borderColor: "rgb(59, 130, 246)",
						borderWidth: 2,
						borderRadius: 4,
						data: ` + string(uniqueVisitvaluesJSON) + `
					},
					{
						label: "Total Visitors",
						backgroundColor: "rgba(16, 185, 129, 0.5)",
						borderColor: "rgb(16, 185, 129)",
						borderWidth: 2,
						borderRadius: 4,
						data: ` + string(totalVisitValuesJSON) + `
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
