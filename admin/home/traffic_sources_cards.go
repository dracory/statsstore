package home

import (
	"fmt"

	"github.com/dracory/hb"
)

type trafficSourceEntry struct {
	Label    string
	Sessions string
}

func trafficSourcesCards(data ControllerData) hb.TagInterface {
	referrerEntries := []trafficSourceEntry{
		{Label: "google.com", Sessions: "671"},
		{Label: "github.com", Sessions: "370"},
		{Label: "com.reddit.frontpage", Sessions: "159"},
		{Label: "producthunt.com", Sessions: "98"},
		{Label: "duckduckgo.com", Sessions: "55"},
		{Label: "alternativeto.net", Sessions: "28"},
		{Label: "trustmrr.com", Sessions: "28"},
		{Label: "bing.com", Sessions: "25"},
		{Label: "openalternative.co", Sessions: "20"},
		{Label: "selfh.st", Sessions: "19"},
	}

	pageEntries := []trafficSourceEntry{
		{Label: "/", Sessions: "4K"},
		{Label: "/docs", Sessions: "727"},
		{Label: "/pricing", Sessions: "719"},
		{Label: "/docs/self-hosting", Sessions: "603"},
		{Label: "/docs/self-host-vs-cloud", Sessions: "272"},
		{Label: "/docs/script", Sessions: "256"},
		{Label: "/docs/roadmap", Sessions: "243"},
		{Label: "/docs/self-hosting-guides/self-hosting-manual", Sessions: "240"},
		{Label: "/features", Sessions: "222"},
		{Label: "/docs/track-events", Sessions: "218"},
	}

	eventEntries := []trafficSourceEntry{
		{Label: "demo", Sessions: "478"},
		{Label: "signup", Sessions: "231"},
		{Label: "hello world", Sessions: "2"},
		{Label: "custom event", Sessions: "2"},
		{Label: "signup button clicked", Sessions: "1"},
	}

	browserEntries := []trafficSourceEntry{
		{Label: "Chrome", Sessions: "2.8K"},
		{Label: "Mobile Safari", Sessions: "802"},
		{Label: "Mobile Chrome", Sessions: "650"},
		{Label: "Firefox", Sessions: "648"},
		{Label: "Edge", Sessions: "246"},
		{Label: "Safari", Sessions: "217"},
		{Label: "Mobile Firefox", Sessions: "63"},
		{Label: "Opera", Sessions: "43"},
		{Label: "Chrome Headless", Sessions: "30"},
		{Label: "Android Browser", Sessions: "28"},
	}

	countryEntries := []trafficSourceEntry{
		{Label: "United States", Sessions: "1.8K"},
		{Label: "Germany", Sessions: "399"},
		{Label: "India", Sessions: "299"},
		{Label: "Canada", Sessions: "243"},
		{Label: "United Kingdom", Sessions: "234"},
		{Label: "Italy", Sessions: "230"},
		{Label: "France", Sessions: "180"},
		{Label: "Australia", Sessions: "126"},
		{Label: "Netherlands", Sessions: "121"},
		{Label: "Brazil", Sessions: "97"},
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
		Child(weeklyTrendsColumn())

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

func weeklyTrendsColumn() hb.TagInterface {
	return hb.Div().
		Class("col").
		Child(weeklyTrendsCard())
}

func weeklyTrendsCard() hb.TagInterface {
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

	weeklyDays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	weeklySlots := []string{"1 AM", "3 AM", "5 AM", "7 AM", "9 AM", "11 AM", "1 PM", "3 PM", "5 PM", "7 PM", "9 PM", "11 PM"}
	weeklyIntensities := [][]int{
		{5, 4, 3, 5, 3, 1, 1},
		{4, 3, 2, 4, 2, 1, 1},
		{3, 2, 2, 3, 2, 1, 0},
		{2, 2, 3, 4, 3, 1, 0},
		{2, 3, 4, 5, 4, 2, 1},
		{3, 4, 5, 5, 4, 2, 1},
		{4, 5, 5, 5, 4, 2, 1},
		{3, 4, 5, 4, 3, 2, 1},
		{2, 3, 4, 3, 2, 1, 1},
		{1, 2, 3, 2, 2, 1, 1},
		{1, 1, 2, 2, 1, 1, 1},
		{0, 1, 1, 1, 1, 0, 0},
	}

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
