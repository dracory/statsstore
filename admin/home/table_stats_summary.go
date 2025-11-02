package home

import (
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
							Text("Unique Visitors"),
						hb.TH().
							Class("text-end").
							Text("Total Visitors"),
					}),
				}),
			hb.Tbody().Children(lo.Map(data.dates, func(date string, index int) hb.TagInterface {
				return hb.TR().Children([]hb.TagInterface{
					hb.TD().Text(date),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.uniqueVisits[index])),
					hb.TD().
						Class("text-end").
						Text(cast.ToString(data.totalVisits[index])),
				})
			})),
			hb.Tfoot().
				Class("table-light fw-bold").
				Children([]hb.TagInterface{
					hb.TR().Children([]hb.TagInterface{
						hb.TH().Text("Total"),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.uniqueVisits))),
						hb.TH().
							Class("text-end").
							Text(cast.ToString(lo.Sum(data.totalVisits))),
					}),
				}),
		})

	return hb.Div().
		Class("table-responsive").
		Child(table)
}
