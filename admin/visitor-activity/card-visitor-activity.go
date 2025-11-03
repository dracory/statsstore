package visitoractivity

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/dracory/statsstore/geostore"
	"github.com/samber/lo"
)

// CardVisitorActivity builds the visitor activity card with detail modal
func CardVisitorActivity(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	card := hb.Div().
		Class("card shadow-sm mb-4").
		Child(cardHeader("Visitor Activity")).
		Child(cardBody(data, ui))

	return hb.Div().
		Child(card).
		Child(visitorDetailModal())
}

func footerControls(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-center justify-content-between gap-3 mt-4").
		Child(paginationSummary(data)).
		Child(quickRangeButtons(data)).
		Child(paginationControls(data))
}

func paginationSummary(data ControllerData) hb.TagInterface {
	if data.TotalCount == 0 {
		return hb.Span().Class("text-muted").Text("No visitors recorded yet")
	}

	start := (data.Page-1)*data.PageSize + 1
	end := data.Page * data.PageSize
	if int64(end) > data.TotalCount {
		end = int(data.TotalCount)
	}

	return hb.Span().
		Class("small text-muted").
		Text(fmt.Sprintf("Showing %d-%d of %d visitors", start, end, data.TotalCount))
}

func quickRangeButtons(data ControllerData) hb.TagInterface {
	btn := func(label, rng string) hb.TagInterface {
		params := map[string]string{"page": "1", "from": "", "to": ""}
		if rng != "" {
			params["range"] = rng
		}
		return hb.A().
			Class("btn btn-sm btn-outline-secondary").
			Href(shared.UrlVisitorActivity(data.Request, queryParamsWith(data, params))).
			Text(label)
	}

	return hb.Div().
		Class("btn-group").
		Attr("role", "group").
		Child(btn("All", "")).
		Child(btn("Last 24 Hours", "24h")).
		Child(btn("Today", "today"))
}

func paginationControls(data ControllerData) hb.TagInterface {
	urlFunc := func(page int) string {
		params := queryParamsWith(data, map[string]string{"page": fmt.Sprintf("%d", page)})
		return shared.UrlVisitorActivity(data.Request, params)
	}

	return shared.PaginationUI(data.Page, data.TotalPages, urlFunc)
}

func queryParamsWith(data ControllerData, overrides map[string]string) map[string]string {
	values := url.Values{}
	for key, val := range data.Request.URL.Query() {
		for _, v := range val {
			values.Add(key, v)
		}
	}

	for key, val := range overrides {
		if val == "" {
			values.Del(key)
			continue
		}
		values.Set(key, val)
	}

	result := map[string]string{}
	for key := range values {
		result[key] = values.Get(key)
	}

	return result
}

func cardHeader(title string) hb.TagInterface {
	actions := hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(exportDropdown()).
		Child(optionsButton())

	return hb.Div().
		Class("card-header d-flex flex-wrap justify-content-between align-items-center gap-2").
		Child(hb.Heading4().
			Class("card-title mb-0").
			HTML(title)).
		Child(actions)
}

func cardBody(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	return hb.Div().
		Class("card-body").
		Child(filterToolbar(data)).
		Child(hb.Div().
			Class("list-group list-group-flush border rounded-3 overflow-hidden").
			Children(lo.Map(data.Visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
				return visitorRow(data, ui, visitor, index)
			}))).
		Child(exportDataTable(data)).
		Child(footerControls(data))
}

func infoLine(label string, value hb.TagInterface) hb.TagInterface {
	labelTag := hb.Span().
		Class("text-muted text-uppercase fw-semibold small flex-shrink-0").
		Text(label)

	valueTag := hb.Div().
		Class("text-body fw-semibold text-break").
		Child(value)

	return hb.Div().
		Class("d-flex gap-2 align-items-baseline lh-sm").
		Child(labelTag).
		Child(valueTag)
}

func infoText(text string) hb.TagInterface {
	return hb.Span().Class("text-body").Text(text)
}

func infoMuted(text string) hb.TagInterface {
	return hb.Span().Class("text-muted fst-italic").Text(text)
}

func filterToolbar(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3").
		Child(addFilterDropdown(data)).
		Child(activeFilterBadges(data.Filters))
}

func addFilterDropdown(data ControllerData) hb.TagInterface {
	items := []struct {
		label  string
		params map[string]string
	}{
		{"Last 24 Hours", queryParamsWith(data, map[string]string{"range": "24h", "from": "", "to": "", "page": "1"})},
		{"Today", queryParamsWith(data, map[string]string{"range": "today", "from": "", "to": "", "page": "1"})},
		{"Country: Unknown", queryParamsWith(data, map[string]string{"country": "empty", "page": "1"})},
		{"Device: Desktop", queryParamsWith(data, map[string]string{"device": "desktop", "page": "1"})},
	}

	menu := hb.UL().Class("dropdown-menu")
	for _, item := range items {
		menu = menu.Child(hb.LI().
			Child(hb.A().
				Class("dropdown-item").
				Href(shared.UrlVisitorActivity(data.Request, item.params)).
				Text(item.label)))
	}

	return hb.Div().
		Class("dropdown").
		Child(hb.Button().
			Class("btn btn-outline-primary dropdown-toggle").
			Attr("type", "button").
			Attr("data-bs-toggle", "dropdown").
			Attr("aria-expanded", "false").
			HTML(`<i class="bi bi-funnel"></i> Add Filter`)).
		Child(menu)
}

func activeFilterBadges(filters FilterOptions) hb.TagInterface {
	tags := []hb.TagInterface{}

	if filters.Range != "" {
		tags = append(tags, hb.Span().Class("badge rounded-pill text-bg-primary").Text(fmt.Sprintf("Range: %s", rangeLabel(filters.Range))))
	}

	if filters.Country != "" {
		label := filters.Country
		if filters.Country == "empty" {
			label = "Unknown"
		}
		tags = append(tags, hb.Span().Class("badge rounded-pill text-bg-info").Text(fmt.Sprintf("Country: %s", strings.ToUpper(label))))
	}

	if filters.Device != "" {
		tags = append(tags, hb.Span().Class("badge rounded-pill text-bg-secondary").Text(fmt.Sprintf("Device: %s", strings.Title(filters.Device))))
	}

	if len(tags) == 0 {
		return hb.Span().Class("text-muted small").Text("No active filters")
	}

	return hb.Div().Class("d-flex flex-wrap gap-2").Children(tags)
}

func rangeLabel(value string) string {
	switch strings.ToLower(value) {
	case "24h", "last24hours", "last_24_hours":
		return "Last 24 Hours"
	case "today":
		return "Today"
	case "7d", "last7days":
		return "Last 7 Days"
	case "30d", "last30days":
		return "Last 30 Days"
	default:
		return value
	}
}

func visitorRow(data ControllerData, ui shared.ControllerOptions, visitor statsstore.VisitorInterface, index int) hb.TagInterface {
	header := hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-start justify-content-between gap-3")

	locationCol := hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold").Text(resolvedVisitorLocation(ui, visitor))).
		Child(hb.Span().Class("small text-muted").Text(visitor.IpAddress()))

	leftHeader := hb.Div().
		Class("d-flex align-items-start gap-2").
		Child(countryBadge(ui, visitor)).
		Child(locationCol)

	rightHeader := hb.Div().
		Class("d-flex flex-wrap gap-2 align-items-center").
		Child(sessionBadge(data.Visitors, visitor)).
		Child(systemSummary(visitor))

	body := hb.Div().
		Class("row gx-3 gy-1 align-items-start mt-2 small lh-sm").
		Child(hb.Div().
			Class("col-lg-5 d-flex flex-column gap-1").
			Child(infoLine("Visit", infoText(formatVisitorTimestamp(visitor.CreatedAt())))).
			Child(infoLine("Duration", infoText(formatVisitDuration(visitor, data.Visitors, index))))).
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-1").
			Child(activityReferrerRow(visitor))).
		Child(hb.Div().
			Class("col-lg-3 d-flex flex-column gap-1").
			Child(activityPathRow(visitor)))

	header = header.Child(leftHeader).Child(rightHeader)

	return hb.Div().
		Class("list-group-item p-2").
		Child(header).
		Child(body)
}

func activityReferrerRow(visitor statsstore.VisitorInterface) hb.TagInterface {
	referrer := visitor.UserReferrer()
	if referrer == "" {
		return infoLine("Referrer", infoMuted("(No referring link)"))
	}

	link := hb.A().
		Href(referrer).
		Class("text-success text-decoration-none").
		Attr("target", "_blank").
		Text(referrer)
	return infoLine("Referrer", link)
}

func activityPathRow(visitor statsstore.VisitorInterface) hb.TagInterface {
	return infoLine("Visited", hb.Raw(getVisitPageLink(visitor.Path())))
}

func resolvedVisitorLocation(visitor statsstore.VisitorInterface) string {
	country := visitor.Country()
	if country == "" {
		return "Unknown Location"
	}
	return strings.ToUpper(country)
}

func systemSummary(visitor statsstore.VisitorInterface) hb.TagInterface {
	systemText := strings.TrimSpace(fmt.Sprintf("%s %s", visitor.UserBrowser(), visitor.UserBrowserVersion()))
	if systemText == "" {
		systemText = "Unknown Browser"
	}
	osText := strings.TrimSpace(fmt.Sprintf("%s %s", visitor.UserOs(), visitor.UserOsVersion()))
	if osText == "" {
		osText = "Unknown OS"
	}

	return hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(deviceIcon(visitor)).
		Child(osIcon(visitor)).
		Child(hb.Span().Class("small").Text(systemText + " on " + osText))
}

func locationBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	country := visitor.Country()
	if country == "" {
		country = "Unknown"
	}
	ip := visitor.IpAddress()
	if ip == "" {
		ip = "Unknown"
	}

	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold").Text(country)).
		Child(hb.Span().Class("small text-muted").Text(fmt.Sprintf("IP Address: %s", ip)))
}

func referrerBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	referrer := visitor.UserReferrer()
	linkText := referrer
	if referrer == "" {
		linkText = "(No referring link)"
	}

	link := hb.Span().Class("text-success").Text(linkText)
	if referrer != "" {
		link = hb.A().
			Href(referrer).
			Class("text-success text-decoration-none").
			Attr("target", "_blank").
			Text(linkText)
	}

	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold small").Text("Referrer")).
		Child(link)
}

func pathBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold small").Text("Visited URL")).
		Child(hb.Raw(getVisitPageLink(visitor.Path())))
}

func sessionBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	fingerprint := visitor.Fingerprint()
	if len(fingerprint) > 8 {
		fingerprint = fingerprint[:8]
	}
	if fingerprint == "" {
		fingerprint = "Session"
	}

	return hb.Span().
		Class("badge text-bg-secondary").
		Text(fmt.Sprintf("Session %s", strings.ToUpper(fingerprint)))
}

func exportDropdown() hb.TagInterface {
	return hb.Div().
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
					Attr("onclick", "exportTableToCSV('visitor-activity-table', 'visitor_activity.csv')").
					Text("Export to CSV"))))
}

func optionsButton() hb.TagInterface {
	return hb.Button().
		Class("btn btn-sm btn-outline-secondary").
		Attr("type", "button").
		HTML(`<i class="bi bi-gear"></i>`)
}

func exportDataTable(data ControllerData) hb.TagInterface {
	head := hb.Thead().
		Child(hb.TR().Children([]hb.TagInterface{
			hb.TH().Text("Visit Time"),
			hb.TH().Text("Path"),
			hb.TH().Text("Country"),
			hb.TH().Text("IP Address"),
			hb.TH().Text("Referrer"),
			hb.TH().Text("Browser"),
			hb.TH().Text("OS"),
			hb.TH().Text("User Agent"),
		}))

	body := hb.Tbody().
		Children(lo.Map(data.Visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
			return hb.TR().Children([]hb.TagInterface{
				hb.TD().Text(formatVisitorTimestamp(visitor.CreatedAt())),
				hb.TD().Text(visitor.Path()),
				hb.TD().Text(strings.ToUpper(visitor.Country())),
				hb.TD().Text(visitor.IpAddress()),
				hb.TD().Text(visitor.UserReferrer()),
				hb.TD().Text(strings.TrimSpace(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())),
				hb.TD().Text(strings.TrimSpace(visitor.UserOs() + " " + visitor.UserOsVersion())),
				hb.TD().Text(visitor.UserAgent()),
			})
		}))

	return hb.Table().
		Class("table table-sm d-none").
		ID("visitor-activity-table").
		Child(head).
		Child(body)
}

func countryBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	code := strings.ToUpper(visitor.Country())
	flag := countryFlagEmoji(code)
	if code == "" {
		code = "--"
	}

	return hb.Span().
		Class("badge bg-light text-dark border").
		Text(fmt.Sprintf("%s %s", flag, code))
}

func formatLocation(visitor statsstore.VisitorInterface) string {
	country := visitor.Country()
	if country == "" {
		return "Unknown Location"
	}
	return strings.ToUpper(country)
}

func countryFlagEmoji(code string) string {
	if len(code) != 2 {
		return "🌐"
	}
	code = strings.ToUpper(code)
	r1 := rune(code[0])
	r2 := rune(code[1])
	if r1 < 'A' || r1 > 'Z' || r2 < 'A' || r2 > 'Z' {
		return "🌐"
	}
	return string(r1-65+0x1F1E6) + string(r2-65+0x1F1E6)
}
