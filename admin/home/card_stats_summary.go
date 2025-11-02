package home

import (
	"github.com/dracory/hb"
	"github.com/dracory/statsstore/admin/shared"
)

// cardStatsSummary creates the stats summary card
func cardStatsSummary(data ControllerData) hb.TagInterface {
	return hb.Div().
		Class("card shadow-sm mb-4").
		Child(hb.Div().
			Class("card-header d-flex flex-column flex-xl-row gap-3 align-items-xl-center justify-content-between").
			Child(hb.Div().
				Class("d-flex flex-column flex-lg-row align-items-lg-center gap-3 w-100").
				Child(hb.Heading4().
					Class("card-title mb-0 flex-grow-1 text-uppercase fw-semibold letter-spacing-1").
					HTML("Stats Summary")).
				Child(hb.Form().
					Class("d-inline-flex align-items-center gap-2 ms-lg-auto").
					Attr("method", "get").
					Attr("action", shared.UrlHome(data.Request)).
					Child(hb.Input().
						Attr("type", "hidden").
						Attr("name", "path").
						Attr("value", shared.PathHome)).
					Child(hb.Label().
						Class("form-label mb-0 text-muted small").
						Attr("for", "stats-period-select").
						Text("Period")).
					Child(hb.Select().
						ID("stats-period-select").
						Class("form-select form-select-sm").
						Attr("name", "period").
						Children(periodOptionsToOptions(data.periodOptions, data.selectedPeriod))).
					Child(hb.Button().
						Class("btn btn-sm btn-outline-primary").
						Attr("type", "submit").
						Text("Apply")).
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
									Text("Export to PDF")))))))).
		Child(hb.Div().
			Class("card-body").
			Child(statsOverview(data)).
			Child(hb.HR().Class("my-4")).
			Child(chartStatsSummary(data)).
			Child(hb.HR().Class("my-4")).
			Child(tableStatsSummary(data)))
}

func periodOptionsToOptions(options []periodOption, selected string) []hb.TagInterface {
	optionTags := make([]hb.TagInterface, 0, len(options))

	for _, option := range options {
		tag := hb.Option().
			Attr("value", option.Value).
			Text(option.Label)

		if option.Value == selected {
			tag = tag.Attr("selected", "selected")
		}

		optionTags = append(optionTags, tag)
	}

	return optionTags
}
