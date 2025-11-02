package home

import (
	"time"

	"github.com/dracory/hb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// tableStatsSummary creates the data table
func tableStatsSummary(data ControllerData) hb.TagInterface {
	table := hb.Table().
		ID("stats-table").
		Class("table table-striped table-hover table-sm").
		Children([]hb.TagInterface{
			hb.Thead().
				Class("table-light").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Date"),
						hb.TH().
							Class("text-end").
							Text("Page Views"),
						hb.TH().
							Class("text-end").
							Text("Unique Visits"),
						hb.TH().
							Class("text-end").
							Text("First Time Visits"),
						hb.TH().
							Class("text-end").
							Text("Returning Visits"),
					}),
				}),
			hb.Tbody().Children(lo.Map(data.dates, func(date string, index int) hb.TagInterface {
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(formatSummaryDate(date)),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.totalVisits[index])),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.uniqueVisits[index])),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.firstVisits[index])),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.returnVisits[index])),
				})
			})),
			hb.Tfoot().
				Class("table-light fw-bold").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Total"),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.totalVisits))),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.uniqueVisits))),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.firstVisits))),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.returnVisits))),
					}),
				}),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}

func formatSummaryDate(date string) string {
	if parsed, err := time.Parse("2006-01-02", date); err == nil {
		return parsed.Format("Mon, 2 Jan 2006")
	}
	return date
}
