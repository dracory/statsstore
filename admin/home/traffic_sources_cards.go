package home

import (
	"fmt"

	"github.com/dracory/hb"
)

func trafficSourcesCards(data ControllerData) hb.TagInterface {
	tsd := computeTrafficSources(data)

	referrerEntries := tsd.Referrers
	if len(referrerEntries) == 0 {
		referrerEntries = []trafficSourceEntry{{Label: "(No data)", Sessions: "0"}}
	}

	pageEntries := tsd.Pages
	if len(pageEntries) == 0 {
		pageEntries = []trafficSourceEntry{{Label: "(No data)", Sessions: "0"}}
	}

	eventEntries := tsd.Events
	if len(eventEntries) == 0 {
		eventEntries = []trafficSourceEntry{{Label: "(No events)", Sessions: "0"}}
	}

	browserEntries := tsd.Browsers
	if len(browserEntries) == 0 {
		browserEntries = []trafficSourceEntry{{Label: "(No data)", Sessions: "0"}}
	}

	countryEntries := tsd.Countries
	if len(countryEntries) == 0 {
		countryEntries = []trafficSourceEntry{{Label: "(No data)", Sessions: "0"}}
	}

	trafficRow := hb.Div().
		Class("row row-cols-1 row-cols-lg-2 g-4").
		Child(trafficSourceColumn("Referrers", "Sessions", referrerEntries, []string{"Referrers", "Channels", "Source", "Medium", "Campaign", "Term"})).
		Child(trafficSourceColumn("Pages", "Sessions", pageEntries, []string{"Pages", "Page Titles", "Entry Pages", "Exit Pages", "Hostnames"}))

	audienceRow := hb.Div().
		Class("row row-cols-1 row-cols-lg-2 g-4").
		Child(trafficSourceColumn("Browsers", "Sessions", browserEntries, []string{"Browsers", "Devices", "Operating Systems", "Screen Dimensions"})).
		Child(trafficSourceColumn("Countries", "Sessions", countryEntries, []string{"Countries", "Regions", "Cities", "Languages", "Map", "Timezones"}))

	engagementRow := hb.Div().
		Class("row row-cols-1 row-cols-lg-2 g-4").
		Child(trafficSourceColumn("Custom Events", "Count", eventEntries, []string{"Custom Events", "Outbound Links"})).
		Child(weeklyTrendsColumn(data))

	return hb.Div().
		Class("d-flex flex-column gap-4").
		Child(trafficRow).
		Child(audienceRow).
		Child(engagementRow)
}

func trafficSourceColumn(title, valueLabel string, entries []trafficSourceEntry, tabs []string) hb.TagInterface {
	return hb.Div().
		Class("col").
		Child(trafficSourceCard(title, valueLabel, entries, tabs))
}

func trafficSourceCard(title, valueLabel string, entries []trafficSourceEntry, tabs []string) hb.TagInterface {
	navLinks := make([]hb.TagInterface, 0, len(tabs))
	for i, tab := range tabs {
		classes := "nav-link text-nowrap"
		if i == 0 {
			classes += " active"
		}
		navLinks = append(navLinks,
			hb.A().
				Class(classes).
				Attr("href", "#").
				Attr("onclick", "return false;").
				Text(tab),
		)
	}

	rows := make([]hb.TagInterface, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows,
			hb.TR().Children([]hb.TagInterface{
				hb.TD().Class("fw-medium").Text(entry.Label),
				hb.TD().Class("text-end").Text(entry.Sessions),
			}),
		)
	}

	return hb.Div().
		Class("card shadow-sm border-0 h-100").
		Child(hb.Div().
			Class("card-header bg-transparent border-bottom-0").
			Child(hb.Div().
				Class("d-flex align-items-center justify-content-between gap-3 flex-wrap").
				Child(hb.Span().
					Class("fw-semibold text-uppercase small text-muted letter-spacing-1").
					Text(title)).
				Child(hb.Button().
					Class("btn btn-sm btn-outline-secondary").
					Attr("type", "button").
					Attr("onclick", "return false;").
					Child(hb.I().Class("bi bi-arrows-fullscreen"))))).
		Child(hb.Div().
			Class("card-header pt-0 bg-transparent border-bottom-0 pb-0").
			Child(hb.Div().
				Class("nav nav-tabs card-header-tabs small overflow-auto flex-nowrap").
				Children(navLinks))).
		Child(hb.Div().
			Class("card-body p-0").
			Child(hb.Div().
				Class("table-responsive").
				Child(hb.Table().
					Class("table table-hover table-sm mb-0").
					Children([]hb.TagInterface{
						hb.Thead().
							Class("table-light").
							Children([]hb.TagInterface{
								hb.TR().Children([]hb.TagInterface{
									hb.TH().Text(title),
									hb.TH().Class("text-end").Text(valueLabel),
								}),
							}),
						hb.Tbody().Children(rows),
					}))))
}

func weeklyTrendsColumn(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("col").
		Child(weeklyTrendsCard(data))
}

func weeklyTrendsCard(data ControllerData) hb.TagInterface {
	tsd := computeTrafficSources(data)
	hm := tsd.Heatmap

	metrics := []string{"Unique Visitors", "Pageviews", "Sessions", "Bounce Rate", "Pages per Session", "Session Duration"}
	selectedMetric := metrics[0]

	dropdownItems := make([]hb.TagInterface, 0, len(metrics))
	for _, metric := range metrics {
		item := hb.A().
			Class("dropdown-item").
			Attr("href", "#").
			Attr("onclick", "return false;").
			Text(metric)
		if metric == selectedMetric {
			item = item.Class("active")
		}
		dropdownItems = append(dropdownItems, item)
	}

	weeklyDays := hm.Days
	weeklySlots := hm.Slots
	weeklyIntensities := hm.Intensities

	headRowCells := []hb.TagInterface{hb.TH().Class("text-muted small fw-normal").Text("")}
	for _, day := range weeklyDays {
		headRowCells = append(headRowCells, hb.TH().Class("text-muted small fw-normal text-center").Text(day))
	}

	headRow := hb.TR().Children(headRowCells)

	bodyRows := make([]hb.TagInterface, 0, len(weeklySlots))
	for slotIndex, slot := range weeklySlots {
		cells := []hb.TagInterface{hb.TH().Class("text-muted small fw-normal text-nowrap").Text(slot)}
		for dayIndex := range weeklyDays {
			level := weeklyIntensities[slotIndex][dayIndex]
			cells = append(cells,
				hb.TD().
					Class("p-1").
					Child(
						hb.Div().
							Class("rounded-1").
							Attr("data-level", fmt.Sprintf("%d", level)).
							Style(fmt.Sprintf("height: 26px; background-color: %s;", heatmapColor(level))),
					),
			)
		}
		bodyRows = append(bodyRows, hb.TR().Children(cells))
	}

	heatmap := hb.Table().
		Class("table table-borderless align-middle mb-0").
		Children([]hb.TagInterface{
			headRow,
			hb.Tbody().Children(bodyRows),
		})

	return hb.Div().
		Class("card shadow-sm border-0 h-100").
		Child(hb.Div().
			Class("card-header bg-transparent border-bottom-0").
			Child(hb.Div().
				Class("d-flex align-items-center justify-content-between gap-3 flex-wrap").
				Child(hb.Span().
					Class("fw-semibold text-uppercase small text-muted letter-spacing-1").
					Text("Weekly Trends")).
				Child(hb.Div().
					Class("dropdown").
					Child(hb.Button().
						Class("btn btn-sm btn-outline-secondary dropdown-toggle").
						Attr("type", "button").
						Attr("data-bs-toggle", "dropdown").
						Attr("aria-expanded", "false").
						Text(selectedMetric)).
					Child(hb.Div().
						Class("dropdown-menu dropdown-menu-end").
						Children(dropdownItems))))).
		Child(hb.Div().
			Class("card-body p-0").
			Child(hb.Div().
				Class("table-responsive").
				Child(heatmap)))
}

func heatmapColor(level int) string {
	switch level {
	case 5:
		return "#1f8254"
	case 4:
		return "#1a9a65"
	case 3:
		return "#17b176"
	case 2:
		return "#14c987"
	case 1:
		return "#12e198"
	default:
		return "#1c2333"
	}
}
