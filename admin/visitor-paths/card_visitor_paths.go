package visitorpaths

import (
	"fmt"
	"strings"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
	"github.com/samber/lo"
)

// CardVisitorPaths builds the visitor paths experience card.
func CardVisitorPaths(data visitorPathsControllerData, ui shared.ControllerOptions) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(cardHeader(data)).
		Child(cardBody(data, ui)).
		Child(filterModal(data))
}

func cardHeader(data visitorPathsControllerData) hb.TagInterface {
	exportParams := shared.QueryParamsWith(data.Request, map[string]string{"action": "export"})
	exportURL := shared.UrlVisitorPaths(data.Request, exportParams)

	actions := hb.Div().
		Class("d-flex align-items-center gap-2").
		Child(shared.ExportDropdown(exportURL)).
		Child(shared.FilterModalButton("visitorPathsFilterModal"))

	return hb.Div().
		Class("card-header d-flex flex-wrap justify-content-between align-items-center gap-2").
		Child(hb.Heading4().
			Class("card-title mb-0").
			HTML("Visitor Paths")).
		Child(actions)
}

func cardBody(data visitorPathsControllerData, ui shared.ControllerOptions) hb.TagInterface {
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
		Child(footerControls(data, ui))
}

func filterToolbar(data visitorPathsControllerData) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-center justify-content-between gap-3").
		Child(shared.FilterModalButton("visitorPathsFilterModal")).
		Child(activeFilterBadges(data.Filters))
}

func filterModal(data visitorPathsControllerData) hb.TagInterface {
	return shared.FilterModal(shared.FilterModalConfig{
		ModalID: "visitorPathsFilterModal",
		Title:   "Filter Visitor Paths",
		Fields: []shared.FilterFieldDef{
			{Name: "range", Label: "Time Range", Type: "select", Options: shared.RangeFilterOptions(), Value: data.Filters.Range},
			{Name: "country", Label: "Country (ISO code or 'empty')", Type: "text", Value: data.Filters.Country},
			{Name: "device", Label: "Device Type", Type: "select", Options: shared.DeviceFilterOptions(), Value: data.Filters.Device},
			{Name: "path_contains", Label: "Path Contains", Type: "text", Value: data.Filters.PathContains},
			{Name: "path_exact", Label: "Path Exact", Type: "text", Value: data.Filters.PathExact},
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

	if filters.From != "" && filters.To != "" {
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

func pathRow(data visitorPathsControllerData, ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	header := hb.Div().
		Class("d-flex flex-column flex-lg-row align-items-lg-start justify-content-between gap-3").
		Child(pathHeaderLeft(ui, visitor)).
		Child(sessionMetadataColumn(data, visitor))

	body := hb.Div().
		Class("row gx-3 gy-1 align-items-start mt-2 small lh-sm").
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-1").
			Child(timestampBlock(visitor)).
			Child(ipBlock(visitor))).
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-1").
			Child(referrerBlock(visitor))).
		Child(hb.Div().
			Class("col-lg-4 d-flex flex-column gap-1").
			Child(userAgentBlock(visitor)))

	return hb.Div().
		Class("list-group-item p-2").
		Child(header).
		Child(body)
}

func pathHeaderLeft(ui shared.ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex align-items-start gap-3").
		Child(shared.CountryBadge(ui, visitor)).
		Child(hb.Div().
			Class("d-flex flex-column gap-1").
			Child(hb.Span().Class("fw-semibold").Text(shared.FormatLocation(ui, visitor))).
			Child(shared.PathLink(ui, visitor.GetPath())))
}

func sessionMetadataColumn(data visitorPathsControllerData, visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Div().
		Class("d-flex flex-wrap justify-content-lg-end gap-2 align-items-center").
		Child(sessionBadge(data, visitor)).
		Child(shared.DeviceBadge(visitor)).
		Child(shared.BrowserBadge(visitor)).
		Child(drillDownButton(data, visitor))
}

func timestampBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	created := shared.FormatTimestamp(visitor.GetCreatedAt())
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("Entry", shared.InfoText(created))).
		Child(shared.InfoLine("Exit", shared.InfoText("-")))
}

func ipBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	ip := visitor.GetIpAddress()
	if ip == "" {
		ip = "Unknown"
	}
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("IP", shared.InfoText(ip)))
}

func referrerBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	referrer := visitor.GetUserReferrer()
	var value hb.TagInterface
	if referrer == "" {
		value = shared.InfoMuted("(No referring link)")
	} else {
		value = hb.A().
			Href(referrer).
			Class("text-success text-decoration-none").
			Attr("target", "_blank").
			Text(referrer)
	}
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("Referrer", value))
}

func userAgentBlock(visitor statsstore.VisitorInterface) hb.TagInterface {
	ua := visitor.GetUserAgent()
	if ua == "" {
		ua = "Unknown"
	}
	return hb.Div().
		Class("d-flex flex-column gap-1").
		Child(shared.InfoLine("User Agent", hb.Span().Class("text-body text-break").Text(ua)))
}

func drillDownButton(data visitorPathsControllerData, visitor statsstore.VisitorInterface) hb.TagInterface {
	params := map[string]string{
		"path": visitor.GetPath(),
		"page": "1",
	}
	drillLink := shared.UrlVisitorActivity(data.Request, params)

	return hb.A().
		Class("btn btn-sm btn-outline-secondary d-flex align-items-center gap-1").
		Attr("href", drillLink).
		Attr("title", "View session in Visitor Activity").
		HTML(`<i class="bi bi-search"></i> View Session`)
}

func footerControls(data visitorPathsControllerData, ui shared.ControllerOptions) hb.TagInterface {
	urlFunc := func(params map[string]string) string {
		return shared.UrlVisitorPaths(data.Request, params)
	}
	return hb.Div().
		Class("d-flex flex-column flex-xl-row align-items-xl-center justify-content-between gap-3").
		Child(shared.PaginationSummary(data.Page, data.PageSize, data.TotalCount, "paths")).
		Child(shared.QuickRangeButtons(data.Request, urlFunc)).
		Child(shared.PerPageSelector(data.Request, data.PageSize, urlFunc)).
		Child(pagination(data.Request, data.Page, data.TotalPages))
}

func sessionBadge(data visitorPathsControllerData, visitor statsstore.VisitorInterface) hb.TagInterface {
	return hb.Span().
		Class("badge text-bg-secondary").
		Text(sessionLabel(data.Paths, visitor))
}

func sessionLabel(visitors []statsstore.VisitorInterface, visitor statsstore.VisitorInterface) string {
	count := sessionCount(visitors, visitor)
	return fmt.Sprintf("Sessions: %d", count)
}

func sessionCount(visitors []statsstore.VisitorInterface, visitor statsstore.VisitorInterface) int {
	targetFingerprint := strings.TrimSpace(visitor.GetFingerprint())
	targetID := strings.TrimSpace(visitor.GetID())

	count := 0

	for _, item := range visitors {
		if targetFingerprint != "" {
			if strings.TrimSpace(item.GetFingerprint()) == targetFingerprint {
				count++
			}
			continue
		}

		if targetID != "" && strings.TrimSpace(item.GetID()) == targetID {
			count++
		}
	}

	if count == 0 {
		count = 1
	}

	return count
}
