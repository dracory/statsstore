package visitorpaths

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// CardVisitorPaths builds the visitor paths experience card.
func CardVisitorPaths(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(cardHeader()).
		Child(cardBody(data, ui))
}

func cardHeader() hb.TagInterface {
	actions := hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(exportDropdown()).
		Child(optionsButton())

	return hb.Div().
		Class("card-header d-flex flex-wrap justify-content-between align-items-center gap-2").
		Child(hb.Heading4().
			Class("card-title mb-0").
			HTML("Visitor Paths")).
		Child(actions)
}

func cardBody(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	var list hb.TagInterface

	if len(data.Paths) == 0 {
		list = hb.Div().
			Class("border rounded-3 p-5 text-center text-muted bg-light").
			Text("No visitor paths recorded yet. Apply different filters or wait for new traffic.")
	} else {
		rows := lo.Map(data.Paths, func(visitor statsstore.VisitorInterface, _ int) hb.TagInterface {
			return pathRow(data, ui, visitor)
		})

		list = hb.Div().
			Class("list-group list-group-flush border rounded-3 overflow-hidden").
			Children(rows)
	}

	return hb.Div().
		Class("card-body d-flex flex-column gap-4").
		Child(filterToolbar(data)).
		Child(list).
		Child(upgradeBanner()).
		Child(exportDataTable(data, ui)).
		Child(footerControls(data, ui))
}

func filterToolbar(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-center justify-content-between gap-3").
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
		{"Path contains '/pricing'", queryParamsWith(data, map[string]string{"path_contains": "pricing", "page": "1"})},
	}

	menu := hb.UL().Class("dropdown-menu")
	for _, item := range items {
		menu = menu.Child(hb.LI().
			Child(hb.A().
				Class("dropdown-item").
				Href(shared.UrlVisitorPaths(data.Request, item.params)).
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
	badges := []hb.TagInterface{}

	if filters.Range != "" {
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-primary").
			Text(fmt.Sprintf("Range: %s", rangeLabel(filters.Range))))
	}

	if filters.From != "" && filters.To != "" {
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-info").
			Text(fmt.Sprintf("Custom Range: %s to %s", shortDate(filters.From), shortDate(filters.To))))
	}

	if filters.Country != "" {
		label := filters.Country
		if filters.Country == "empty" {
			label = "Unknown"
		}
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-success").
			Text(fmt.Sprintf("Country: %s", strings.ToUpper(label))))
	}

	if filters.PathContains != "" {
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-secondary").
			Text(fmt.Sprintf("Path contains '%s'", filters.PathContains)))
	}

	if filters.PathExact != "" {
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-dark").
			Text(fmt.Sprintf("Path is '%s'", filters.PathExact)))
	}

	if filters.Device != "" {
		label := filters.Device
		if filters.Device == "empty" {
			label = "Unknown"
		}
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-warning").
			Text(fmt.Sprintf("Device: %s", strings.Title(label))))
	}

	if len(badges) == 0 {
		return hb.Div().
			Class("text-muted small").
			Text("No active filters")
	}

	return hb.Div().Class("d-flex flex-wrap gap-2").Children(badges)
}

func pathRow(data ControllerData, ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	header := hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-start justify-content-between gap-3").
		Child(pathHeaderLeft(ui, visitor)).
		Child(sessionMetadataColumn(data, visitor))

	body := hb.Div().
		Class("row g-3 align-items-start mt-2").
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-2").
			Child(timestampBlock(visitor)).
			Child(ipBlock(visitor))).
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-2").
			Child(referrerBlock(visitor)).
			Child(pathMetaBlock(ui, visitor))).
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-2").
			Child(userAgentBlock(visitor)))

	return hb.Div().
		Class("list-group-item p-3").
		Child(header).
		Child(body)
}

func pathHeaderLeft(ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	host := websiteHost(ui)

	return hb.Div().
		Class("d-flex align-items-start gap-3").
		Child(countryBadge(visitor)).
		Child(hb.Div().
			Class("d-flex flex-column gap-1").
			Child(hb.Div().
				Class("d-flex flex-wrap align-items-center gap-2").
				Child(hb.Span().Class("fw-semibold").Text(formatLocation(visitor))).
				Child(hb.Span().Class("badge text-bg-light").Text(host))).
			Child(pathLink(ui, visitor.Path())))
}

func sessionMetadataColumn(data ControllerData, visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-wrap justify-content-lg-end gap-2 align-items-center").
		Child(sessionBadge(visitor)).
		Child(deviceBadge(visitor)).
		Child(browserBadge(visitor)).
		Child(drillDownButton(data, visitor))
}

func timestampBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	created := formatTimestamp(visitor.CreatedAt())
	return hb.Div().
		Class("d-flex flex-column small text-muted gap-1").
		Child(hb.Span().Text(fmt.Sprintf("Entry: %s", created))).
		Child(hb.Span().Text("Exit: -"))
}

func ipBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	ip := visitor.IpAddress()
	if ip == "" {
		ip = "Unknown"
	}
	return hb.Div().
		Class("small text-muted").
		Text(fmt.Sprintf("IP Address: %s", ip))
}

func referrerBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	referrer := visitor.UserReferrer()
	if referrer == "" {
		return hb.Div().
			Class("d-flex flex-column gap-1").
			Child(hb.Span().Class("fw-semibold small").Text("Referrer")).
			Child(hb.Span().Class("text-muted small").Text("(No referring link)"))
	}

	link := hb.A().
		Href(referrer).
		Class("text-success text-decoration-none").
		Attr("target", "_blank").
		Text(referrer)

	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold small").Text("Referrer")).
		Child(link)
}

func pathMetaBlock(ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	absolute := fullPathURL(ui, visitor.Path())

	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold small").Text("Visited URL")).
		Child(hb.A().
			Href(absolute).
			Class("text-success text-decoration-none d-inline-flex align-items-center gap-1").
			Attr("target", "_blank").
			HTML(fmt.Sprintf("%s <i class=\"bi bi-box-arrow-up-right\"></i>", visitor.Path())))
}

func userAgentBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	ua := visitor.UserAgent()
	if ua == "" {
		ua = "Unknown"
	}
	return hb.Div().
		Class("small text-muted").
		Text(fmt.Sprintf("User Agent: %s", ua))
}

func drillDownButton(data ControllerData, visitor statsstore.VisitorInterface) hb.TagInterface {
	params := map[string]string{
		"path": visitor.Path(),
		"page": "1",
	}
	drillLink := shared.UrlVisitorActivity(data.Request, params)

	return hb.A().
		Class("btn btn-sm btn-outline-secondary d-flex align-items-center gap-1").
		Attr("href", drillLink).
		Attr("title", "View session in Visitor Activity").
		HTML(`<i class="bi bi-search"></i> View Session`)
}

func footerControls(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column flex-xl-row align-items-xl-center justify-content-between gap-3").
		Child(paginationSummary(data)).
		Child(quickRangeButtons(data)).
		Child(perPageSelector(data)).
		Child(pagination(data.Request, data.Page, data.TotalPages))
}

func paginationSummary(data ControllerData) hb.TagInterface {
	if data.TotalCount == 0 {
		return hb.Span().Class("text-muted small").Text("No visitor paths to display")
	}

	start := (data.Page-1)*data.PageSize + 1
	end := data.Page * data.PageSize
	if int64(end) > data.TotalCount {
		end = int(data.TotalCount)
	}

	return hb.Span().
		Class("small text-muted").
		Text(fmt.Sprintf("Showing %d-%d of %d paths", start, end, data.TotalCount))
}

func quickRangeButtons(data ControllerData) hb.TagInterface {
	btn := func(label, rng string) hb.TagInterface {
		params := map[string]string{"page": "1", "from": "", "to": ""}
		if rng != "" {
			params["range"] = rng
		}
		return hb.A().
			Class("btn btn-sm btn-outline-secondary").
			Href(shared.UrlVisitorPaths(data.Request, queryParamsWith(data, params))).
			Text(label)
	}

	return hb.Div().
		Class("btn-group").
		Attr("role", "group").
		Child(btn("All", "")).
		Child(btn("Last 24 Hours", "24h")).
		Child(btn("Today", "today")).
		Child(btn("Last 7 Days", "7d"))
}

func perPageSelector(data ControllerData) hb.TagInterface {
	options := []int{10, 25, 50, 100}
	group := hb.Div().Class("d-flex align-items-center gap-2")
	group = group.Child(hb.Span().Class("small text-muted").Text("Rows per page:"))

	buttons := hb.Div().Class("btn-group btn-group-sm")
	for _, size := range options {
		params := queryParamsWith(data, map[string]string{"per_page": fmt.Sprintf("%d", size), "page": "1"})
		classes := "btn btn-outline-secondary"
		if data.PageSize == size {
			classes = "btn btn-secondary"
		}
		buttons = buttons.Child(hb.A().
			Class(classes).
			Href(shared.UrlVisitorPaths(data.Request, params)).
			Text(fmt.Sprintf("%d", size)))
	}

	return group.Child(buttons)
}

func exportDropdown() hb.TagInterface {
	button := hb.Button().
		Class("btn btn-sm btn-outline-secondary dropdown-toggle").
		Attr("type", "button").
		Attr("data-bs-toggle", "dropdown").
		Attr("aria-expanded", "false").
		Text("Export")

	item := hb.A().
		Class("dropdown-item").
		Href("#").
		Attr("onclick", "exportTableToCSV('visitor-paths-table', 'visitor_paths.csv')").
		Text("Export to CSV")

	menu := hb.UL().
		Class("dropdown-menu").
		Child(hb.LI().Child(item))

	return hb.Div().
		Class("dropdown").
		Child(button).
		Child(menu)
}

func optionsButton() hb.TagInterface {
	return hb.Button().
		Class("btn btn-sm btn-outline-secondary").
		Attr("type", "button").
		HTML(`<i class="bi bi-gear"></i>`)
}

func upgradeBanner() hb.TagInterface {
	return hb.Div().
		Class("alert alert-info text-center mb-0").
		HTML("<strong>Upgrade Insight:</strong> Connect deeper analytics to unlock funnel visualisations and path grouping.")
}

func exportDataTable(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	head := hb.Thead().
		Child(hb.TR().Children([]hb.TagInterface{
			hb.TH().Text("Visit Time"),
			hb.TH().Text("Path"),
			hb.TH().Text("Absolute URL"),
			hb.TH().Text("Country"),
			hb.TH().Text("IP Address"),
			hb.TH().Text("Referrer"),
			hb.TH().Text("Session"),
			hb.TH().Text("Device"),
			hb.TH().Text("Browser"),
		}))

	body := hb.Tbody().
		Children(lo.Map(data.Paths, func(visitor statsstore.VisitorInterface, _ int) hb.TagInterface {
			absolute := fullPathURL(ui, visitor.Path())
			return hb.TR().Children([]hb.TagInterface{
				hb.TD().Text(formatTimestamp(visitor.CreatedAt())),
				hb.TD().Text(visitor.Path()),
				hb.TD().Text(absolute),
				hb.TD().Text(strings.ToUpper(visitor.Country())),
				hb.TD().Text(visitor.IpAddress()),
				hb.TD().Text(visitor.UserReferrer()),
				hb.TD().Text(sessionLabel(visitor)),
				hb.TD().Text(visitor.UserDevice()),
				hb.TD().Text(strings.TrimSpace(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())),
			})
		}))

	return hb.Table().
		Class("table table-sm d-none").
		ID("visitor-paths-table").
		Child(head).
		Child(body)
}

func sessionBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Span().
		Class("badge text-bg-secondary").
		Text(sessionLabel(visitor))
}

func sessionLabel(visitor statsstore.VisitorInterface) string {
	fingerprint := visitor.Fingerprint()
	if len(fingerprint) > 8 {
		fingerprint = fingerprint[:8]
	}
	if fingerprint == "" {
		fingerprint = "Session"
	}
	return fmt.Sprintf("Session %s", strings.ToUpper(fingerprint))
}

func deviceBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	deviceType := strings.ToLower(visitor.UserDeviceType())
	label := visitor.UserDeviceType()
	if label == "" {
		label = "Unknown"
	}

	classes := "badge bg-light text-dark border"
	switch {
	case strings.Contains(deviceType, "desktop"):
		classes = "badge bg-primary-subtle text-primary"
	case strings.Contains(deviceType, "mobile"):
		classes = "badge bg-success-subtle text-success"
	case strings.Contains(deviceType, "tablet"):
		classes = "badge bg-info-subtle text-info"
	case strings.Contains(deviceType, "bot"):
		classes = "badge bg-warning-subtle text-warning"
	}

	return hb.Span().Class(classes).Text(strings.Title(label))
}

func browserBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	browser := strings.TrimSpace(visitor.UserBrowser() + " " + visitor.UserBrowserVersion())
	if browser == "" {
		browser = "Unknown Browser"
	}

	return hb.Span().
		Class("badge bg-light text-dark border").
		Text(browser)
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

func formatLocation(visitor statsstore.VisitorInterface) string {
	country := visitor.Country()
	if country == "" {
		return "Unknown Location"
	}
	return strings.ToUpper(country)
}

func formatTimestamp(value string) string {
	if value == "" {
		return "Unknown"
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return value
}

func shortDate(value string) string {
	if value == "" {
		return "-"
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.Format("2006-01-02")
	}
	return value
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

func pathLink(ui shared.ControllerOptions, path string) hb.TagInterface {
	absolute := fullPathURL(ui, path)
	display := path
	if display == "" {
		display = "/"
	}

	return hb.A().
		Href(absolute).
		Class("text-success text-decoration-none d-inline-flex align-items-center gap-1").
		Attr("target", "_blank").
		HTML(fmt.Sprintf("%s <i class=\"bi bi-box-arrow-up-right\"></i>", display))
}

func fullPathURL(ui shared.ControllerOptions, path string) string {
	base := ui.WebsiteUrl
	if base == "" {
		return path
	}

	u, err := url.Parse(base)
	if err != nil {
		return base + path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u.Path = path
	return u.String()
}

func websiteHost(ui shared.ControllerOptions) string {
	if ui.WebsiteUrl == "" {
		return "host"
	}
	if parsed, err := url.Parse(ui.WebsiteUrl); err == nil && parsed.Host != "" {
		return parsed.Host
	}
	return ui.WebsiteUrl
}

func queryParamsWith(data ControllerData, overrides map[string]string) map[string]string {
	values := url.Values{}
	for key, vals := range data.Request.URL.Query() {
		for _, v := range vals {
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
