package visitoractivity

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// CardVisitorActivity builds the visitor activity card with detail modal
func CardVisitorActivity(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	card := hb.Div().
		Class("card shadow-sm mb-4").
		Child(cardHeader("Visitor Activity", data)).
		Child(cardBody(data, ui)).
		Child(filterModal(data))

	return hb.Div().
		Child(card).
		Child(visitorDetailModal())
}

func footerControls(data ControllerData) hb.TagInterface {
	urlFunc := func(params map[string]string) string {
		return shared.UrlVisitorActivity(data.Request, params)
	}
	return hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-center justify-content-between gap-3 mt-4").
		Child(shared.PaginationSummary(data.Page, data.PageSize, data.TotalCount, "visitors")).
		Child(shared.QuickRangeButtons(data.Request, urlFunc)).
		Child(shared.PerPageSelector(data.Request, data.PageSize, urlFunc)).
		Child(paginationControls(data))
}

func paginationControls(data ControllerData) hb.TagInterface {
	urlFunc := func(page int) string {
		params := shared.QueryParamsWith(data.Request, map[string]string{"page": fmt.Sprintf("%d", page)})
		return shared.UrlVisitorActivity(data.Request, params)
	}

	return shared.PaginationUI(data.Page, data.TotalPages, urlFunc)
}

func queryParamsWith(data ControllerData, overrides map[string]string) map[string]string {
	return shared.QueryParamsWith(data.Request, overrides)
}

func cardHeader(title string, data ControllerData) hb.TagInterface {
	exportParams := shared.QueryParamsWith(data.Request, map[string]string{"action": "export"})
	exportURL := shared.UrlVisitorActivity(data.Request, exportParams)

	actions := hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(shared.ExportDropdown(exportURL)).
		Child(shared.FilterModalButton("visitorActivityFilterModal"))

	return hb.Div().
		Class("card-header d-flex flex-wrap justify-content-between align-items-center gap-2").
		Child(hb.Heading4().
			Class("card-title mb-0").
			HTML(title)).
		Child(actions)
}

func cardBody(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	list := hb.Div()
	if len(data.Visitors) == 0 {
		list = hb.Div().
			Class("border rounded-3 p-5 text-center text-muted bg-light").
			Text("No visitors recorded yet. Apply different filters or wait for new traffic.")
	} else {
		list = hb.Div().
			Class("list-group list-group-flush border rounded-3 overflow-hidden").
			Children(lo.Map(data.Visitors, func(visitor statsstore.VisitorInterface, index int) hb.TagInterface {
				return visitorRow(data, ui, visitor, index)
			}))
	}

	return hb.Div().
		Class("card-body").
		Child(filterToolbar(data)).
		Child(list).
		Child(footerControls(data))
}

func infoLine(label string, value hb.TagInterface) hb.TagInterface {
	return shared.InfoLine(label, value)
}

func infoText(text string) hb.TagInterface {
	return shared.InfoText(text)
}

func infoMuted(text string) hb.TagInterface {
	return shared.InfoMuted(text)
}

func filterToolbar(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-wrap align-items-center justify-content-between gap-2 mb-3").
		Child(shared.FilterModalButton("visitorActivityFilterModal")).
		Child(activeFilterBadges(data.Filters))
}

func filterModal(data ControllerData) hb.TagInterface {
	return shared.FilterModal(shared.FilterModalConfig{
		ModalID: "visitorActivityFilterModal",
		Title:   "Filter Visitor Activity",
		Fields: []shared.FilterFieldDef{
			{Name: "range", Label: "Time Range", Type: "select", Options: shared.RangeFilterOptions(), Value: data.Filters.Range},
			{Name: "country", Label: "Country (ISO code or 'empty')", Type: "text", Value: data.Filters.Country},
			{Name: "device", Label: "Device Type", Type: "select", Options: shared.DeviceFilterOptions(), Value: data.Filters.Device},
		},
	})
}

func activeFilterBadges(filters FilterOptions) hb.TagInterface {
	tags := []hb.TagInterface{}

	if filters.Range != "" {
		tags = append(tags, hb.Span().Class("badge rounded-pill text-bg-primary").Text(fmt.Sprintf("Range: %s", shared.RangeLabel(filters.Range))))
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
	return shared.RangeLabel(value)
}

func visitorRow(data ControllerData, ui shared.ControllerOptions, visitor statsstore.VisitorInterface, index int) hb.TagInterface {
	header := hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-start justify-content-between gap-3")

	locationCol := hb.Div().
		Class("d-flex flex-column gap-1").
		Child(hb.Span().Class("fw-semibold").Text(shared.FormatLocation(ui, visitor))).
		Child(hb.Span().Class("small text-muted").Text(visitor.GetIpAddress()))

	leftHeader := hb.Div().
		Class("d-flex align-items-start gap-2").
		Child(shared.CountryBadge(ui, visitor)).
		Child(locationCol)

	rightHeader := hb.Div().
		Class("d-flex flex-wrap gap-2 align-items-center").
		Child(sessionBadge(visitor)).
		Child(systemSummary(visitor)).
		Child(detailButton(visitor))

	body := hb.Div().
		Class("row gx-3 gy-1 align-items-start mt-2 small lh-sm").
		Child(hb.Div().
			Class("col-lg-5 d-flex flex-column gap-1").
			Child(infoLine("Visit", infoText(formatVisitorTimestamp(visitor.GetCreatedAt())))).
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
	referrer := visitor.GetUserReferrer()
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
	path := visitor.GetPath()
	if path == "" {
		path = "/"
	}

	link := hb.A().
		Href(path).
		Class("text-primary text-decoration-none").
		Attr("target", "_blank").
		Text(path)
	return infoLine("Path", link)
}

func systemSummary(visitor statsstore.VisitorInterface) hb.TagInterface {
	systemText := strings.TrimSpace(fmt.Sprintf("%s %s", visitor.GetUserBrowser(), visitor.GetUserBrowserVersion()))
	if systemText == "" {
		systemText = "Unknown Browser"
	}

	osText := strings.TrimSpace(fmt.Sprintf("%s %s", visitor.GetUserOs(), visitor.GetUserOsVersion()))
	if osText == "" {
		osText = "Unknown OS"
	}

	return hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(deviceIcon(visitor)).
		Child(osIcon(visitor)).
		Child(hb.Span().Class("small").Text(systemText + " on " + osText))
}

func sessionBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	fingerprint := visitor.GetFingerprint()
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

func detailButton(visitor statsstore.VisitorInterface) hb.TagInterface {
	data, _ := json.Marshal(map[string]string{
		"id":          visitor.GetID(),
		"createdAt":   visitor.GetCreatedAt(),
		"path":        visitor.GetPath(),
		"country":     visitor.GetCountry(),
		"ipAddress":   visitor.GetIpAddress(),
		"referrer":    visitor.GetUserReferrer(),
		"device":      visitor.GetUserDevice(),
		"browser":     visitor.GetUserBrowser(),
		"browserVer":  visitor.GetUserBrowserVersion(),
		"os":          visitor.GetUserOs(),
		"osVer":       visitor.GetUserOsVersion(),
		"userAgent":   visitor.GetUserAgent(),
		"fingerprint": visitor.GetFingerprint(),
	})

	return hb.Button().
		Class("btn btn-sm btn-outline-primary").
		Attr("type", "button").
		Attr("data-bs-toggle", "modal").
		Attr("data-bs-target", "#visitorDetailModal").
		Attr("data-visitor", string(data)).
		HTML(`<i class="bi bi-info-circle"></i> Details`)
}
