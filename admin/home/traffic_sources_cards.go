package home

import (
	"fmt"
	"strings"

	"github.com/dracory/hb"
)

type tabData struct {
	Label   string
	Entries []trafficSourceEntry
}

func ensureEntries(entries []trafficSourceEntry, emptyLabel string) []trafficSourceEntry {
	if len(entries) == 0 {
		return []trafficSourceEntry{{Label: emptyLabel, Sessions: "0"}}
	}
	return entries
}

func trafficSourcesCards(data ControllerData) hb.TagInterface {
	tsd := computeTrafficSources(data)

	trafficTabs := []tabData{
		{Label: "Referrers", Entries: ensureEntries(tsd.Referrers, "(No data)")},
		{Label: "Channels", Entries: ensureEntries(tsd.Channels, "(No data)")},
		{Label: "Source", Entries: ensureEntries(tsd.Sources, "(No data)")},
		{Label: "Medium", Entries: ensureEntries(tsd.Mediums, "(No data)")},
		{Label: "Campaign", Entries: ensureEntries(tsd.Campaigns, "(No data)")},
		{Label: "Term", Entries: ensureEntries(tsd.Terms, "(No data)")},
	}

	pageTabs := []tabData{
		{Label: "Pages", Entries: ensureEntries(tsd.Pages, "(No data)")},
		{Label: "Page Titles", Entries: ensureEntries(nil, "(No data)")},
		{Label: "Entry Pages", Entries: ensureEntries(tsd.EntryPages, "(No data)")},
		{Label: "Exit Pages", Entries: ensureEntries(tsd.ExitPages, "(No data)")},
		{Label: "Hostnames", Entries: ensureEntries(nil, "(No data)")},
	}

	audienceTabs := []tabData{
		{Label: "Browsers", Entries: ensureEntries(tsd.Browsers, "(No data)")},
		{Label: "Devices", Entries: ensureEntries(tsd.Devices, "(No data)")},
		{Label: "Operating Systems", Entries: ensureEntries(tsd.OperatingSystems, "(No data)")},
		{Label: "Screen Dimensions", Entries: ensureEntries(nil, "(No data)")},
	}

	geoTabs := []tabData{
		{Label: "Countries", Entries: ensureEntries(tsd.Countries, "(No data)")},
		{Label: "Regions", Entries: ensureEntries(nil, "(No data)")},
		{Label: "Cities", Entries: ensureEntries(nil, "(No data)")},
		{Label: "Languages", Entries: ensureEntries(tsd.Languages, "(No data)")},
		{Label: "Map", Entries: ensureEntries(nil, "(No map data)")},
		{Label: "Timezones", Entries: ensureEntries(nil, "(No data)")},
	}

	eventTabs := []tabData{
		{Label: "Custom Events", Entries: ensureEntries(tsd.Events, "(No events)")},
		{Label: "Outbound Links", Entries: ensureEntries(tsd.OutboundLinks, "(No outbound links)")},
	}

	trafficRow := hb.Div().
		Class("row row-cols-1 row-cols-lg-2 g-4").
		Child(trafficSourceColumn("Referrers", "Sessions", trafficTabs)).
		Child(trafficSourceColumn("Pages", "Sessions", pageTabs))

	audienceRow := hb.Div().
		Class("row row-cols-1 row-cols-lg-2 g-4").
		Child(trafficSourceColumn("Browsers", "Sessions", audienceTabs)).
		Child(trafficSourceColumn("Countries", "Sessions", geoTabs))

	engagementRow := hb.Div().
		Class("row row-cols-1 row-cols-lg-2 g-4").
		Child(trafficSourceColumn("Custom Events", "Count", eventTabs)).
		Child(weeklyTrendsColumn(data))

	return hb.Div().
		Class("d-flex flex-column gap-4").
		Child(trafficRow).
		Child(audienceRow).
		Child(engagementRow)
}

func trafficSourceColumn(title, valueLabel string, tabs []tabData) hb.TagInterface {
	return hb.Div().
		Class("col").
		Child(trafficSourceCard(title, valueLabel, tabs))
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

func trafficSourceCard(title, valueLabel string, tabs []tabData) hb.TagInterface {
	cardSlug := slugify(title)

	navLinks := make([]hb.TagInterface, 0, len(tabs))
	tabPanes := make([]hb.TagInterface, 0, len(tabs))

	for i, tab := range tabs {
		tabID := fmt.Sprintf("%s-tab-%d", cardSlug, i)
		paneID := fmt.Sprintf("%s-pane-%d", cardSlug, i)

		linkClasses := "nav-link text-nowrap"
		paneClasses := "tab-pane fade"
		if i == 0 {
			linkClasses += " active"
			paneClasses += " show active"
		}

		navLinks = append(navLinks,
			hb.A().
				Class(linkClasses).
				Attr("id", tabID).
				Attr("data-bs-toggle", "tab").
				Attr("data-bs-target", "#"+paneID).
				Attr("type", "button").
				Attr("role", "tab").
				Attr("aria-controls", paneID).
				Attr("aria-selected", fmt.Sprintf("%t", i == 0)).
				Text(tab.Label),
		)

		rows := make([]hb.TagInterface, 0, len(tab.Entries))
		for _, entry := range tab.Entries {
			rows = append(rows,
				hb.TR().Children([]hb.TagInterface{
					hb.TD().Class("fw-medium").Text(entry.Label),
					hb.TD().Class("text-end").Text(entry.Sessions),
				}),
			)
		}

		tabPanes = append(tabPanes,
			hb.Div().
				Class(paneClasses).
				Attr("id", paneID).
				Attr("role", "tabpanel").
				Attr("aria-labelledby", tabID).
				Child(hb.Div().
					Class("table-responsive").
					Child(hb.Table().
						Class("table table-hover table-sm mb-0").
						Children([]hb.TagInterface{
							hb.Thead().
								Class("table-light").
								Children([]hb.TagInterface{
									hb.TR().Children([]hb.TagInterface{
										hb.TH().Text(tab.Label),
										hb.TH().Class("text-end").Text(valueLabel),
									}),
								}),
							hb.Tbody().Children(rows),
						}))),
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
				Attr("role", "tablist").
				Children(navLinks))).
		Child(hb.Div().
			Class("card-body p-0").
			Child(hb.Div().
				Class("tab-content").
				Children(tabPanes)))
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