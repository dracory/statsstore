package pageviewactivity

import (
	"fmt"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// CardPageViewActivity builds the page view activity card with filter toolbar,
// visitor list, and footer controls.
func CardPageViewActivity(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(cardHeader(data)).
		Child(cardBody(data, ui)).
		Child(filterModal(data))
}

func cardHeader(data ControllerData) hb.TagInterface {
	exportParams := shared.QueryParamsWith(data.Request, map[string]string{"action": "export"})
	exportURL := shared.UrlPageViewActivity(data.Request, exportParams)

	actions := hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(shared.ExportDropdown(exportURL)).
		Child(shared.FilterModalButton("pageViewActivityFilterModal"))

	return hb.Div().
		Class("card-header d-flex flex-wrap justify-content-between align-items-center gap-2").
		Child(hb.Heading4().
			Class("card-title mb-0").
			HTML("Page View Activity")).
		Child(actions)
}

func cardBody(data ControllerData, ui shared.ControllerOptions) hb.TagInterface {
	var list hb.TagInterface

	if len(data.Visitors) == 0 {
		list = hb.Div().
			Class("border rounded-3 p-5 text-center text-muted bg-light").
			Text("No page views recorded yet. Apply different filters or wait for new traffic.")
	} else {
		rows := lo.Map(data.Visitors, func(visitor statsstore.VisitorInterface, _ int) hb.TagInterface {
			return pageViewRow(data, ui, visitor)
		})

		list = hb.Div().
			Class("list-group list-group-flush border rounded-3 overflow-hidden").
			Children(rows)
	}

	return hb.Div().
		Class("card-body d-flex flex-column gap-4").
		Child(filterToolbar(data)).
		Child(list).
		Child(footerControls(data))
}

// == FILTER TOOLBAR ===========================================================

func filterToolbar(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-center justify-content-between gap-3").
		Child(shared.FilterModalButton("pageViewActivityFilterModal")).
		Child(activeFilterBadges(data.Filters))
}

func filterModal(data ControllerData) hb.TagInterface {
	return shared.FilterModal(shared.FilterModalConfig{
		ModalID: "pageViewActivityFilterModal",
		Title:   "Filter Page View Activity",
		Fields: []shared.FilterFieldDef{
			{Name: "range", Label: "Time Range", Type: "select", Options: shared.RangeFilterOptions(), Value: data.Filters.Range},
			{Name: "country", Label: "Country (ISO code or 'empty')", Type: "text", Value: data.Filters.Country},
			{Name: "device", Label: "Device Type", Type: "select", Options: shared.DeviceFilterOptions(), Value: data.Filters.Device},
			{Name: "browser", Label: "Browser", Type: "text", Value: data.Filters.Browser},
		},
	})
}

func activeFilterBadges(filters FilterOptions) hb.TagInterface {
	badges := []hb.TagInterface{}

	if filters.Range != "" {
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-primary").
			Text(fmt.Sprintf("Range: %s", shared.RangeLabel(filters.Range))))
	}

	if filters.From != "" && filters.To != "" && filters.Range == "" {
		badges = append(badges, hb.Span().
			Class("badge rounded-pill text-bg-info").
			Text(fmt.Sprintf("Custom Range: %s to %s", shared.ShortDate(filters.From), shared.ShortDate(filters.To))))
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

// == ROW RENDERING ===========================================================

func pageViewRow(data ControllerData, ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	header := hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-start justify-content-between gap-3").
		Child(rowHeaderLeft(ui, visitor)).
		Child(rowHeaderRight(visitor))

	body := hb.Div().
		Class("row gx-3 gy-1 align-items-start mt-2 small lh-sm").
		Child(hb.Div().
			Class("col-lg-3 d-flex flex-column gap-1").
			Child(timestampCell(visitor))).
		Child(hb.Div().
			Class("col-lg-3 d-flex flex-column gap-1").
			Child(locationCell(ui, visitor))).
		Child(hb.Div().
			Class("col-lg-3 d-flex flex-column gap-1").
			Child(referrerCell(visitor))).
		Child(hb.Div().
			Class("col-lg-3 d-flex flex-column gap-1").
			Child(pathCell(ui, visitor)))

	return hb.Div().
		Class("list-group-item p-2").
		Child(header).
		Child(body)
}

func rowHeaderLeft(ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex align-items-start gap-3").
		Child(shared.CountryBadge(ui, visitor)).
		Child(hb.Div().
			Class("d-flex flex-column gap-1").
			Child(hb.Span().Class("fw-semibold").Text(shared.FormatLocation(ui, visitor))))
}

func rowHeaderRight(visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-wrap justify-content-lg-end gap-2 align-items-center").
		Child(shared.DeviceBadge(visitor)).
		Child(shared.BrowserBadge(visitor)).
		Child(osBadge(visitor))
}

func timestampCell(visitor statsstore.VisitorInterface) hb.TagInterface {
	date, timeStr := splitTimestamp(visitor.GetCreatedAt())
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("Date", shared.InfoText(date))).
		Child(shared.InfoLine("Time", shared.InfoText(timeStr)))
}

func locationCell(ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	ip := visitor.GetIpAddress()
	if ip == "" {
		ip = "Unknown"
	}
	lang := visitor.GetUserAcceptLanguage()
	if lang == "" {
		lang = "Unknown"
	}
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("IP", shared.InfoText(ip))).
		Child(shared.InfoLine("Language", shared.InfoText(lang)))
}

func referrerCell(visitor statsstore.VisitorInterface) hb.TagInterface {
	referrer := visitor.GetUserReferrer()
	var value hb.TagInterface
	if referrer == "" {
		value = shared.InfoMuted("(No referring link)")
	} else {
		value = hb.A().
			Href(referrer).
			Class("text-success text-decoration-none").
			Attr("target", "_blank").
			Attr("data-bs-toggle", "tooltip").
			Attr("title", referrer).
			Text(referrer)
	}
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("Referrer", value))
}

func pathCell(ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("Page", shared.PathLink(ui, visitor.GetPath())))
}

// == FOOTER CONTROLS =========================================================

func footerControls(data ControllerData) hb.TagInterface {
	urlFunc := func(params map[string]string) string {
		return shared.UrlPageViewActivity(data.Request, params)
	}
	return hb.Div().
		Class("d-flex flex-column flex-xl-row align-items-xl-center justify-content-between gap-3").
		Child(shared.PaginationSummary(data.Page, data.PageSize, data.TotalCount, "page views")).
		Child(shared.QuickRangeButtons(data.Request, urlFunc)).
		Child(shared.PerPageSelector(data.Request, data.PageSize, urlFunc)).
		Child(pagination(data.Request, data.Page, data.TotalPages))
}

// == EXPORT & OPTIONS ========================================================

// == BADGES & ICONS ==========================================================

func osBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	os := strings.TrimSpace(visitor.GetUserOs() + " " + visitor.GetUserOsVersion())
	if os == "" {
		os = "Unknown OS"
	}

	return hb.Span().
		Class("badge bg-light text-dark border").
		Text(os)
}

// == INFO LINE HELPERS =======================================================

// == URL & QUERY HELPERS =====================================================

func queryParamsWith(data ControllerData, overrides map[string]string) map[string]string {
	return shared.QueryParamsWith(data.Request, overrides)
}

// == FORMATTING HELPERS ======================================================
