package home

import (
	"fmt"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// periodComparisonTable renders a side-by-side comparison of the current and
// previous periods for the key dashboard metrics.
func periodComparisonTable(data ControllerData) hb.TagInterface {
	current := periodComparisonRow{
		label:     "Current Period",
		unique:    lo.Sum(data.uniqueVisits),
		total:     lo.Sum(data.totalVisits),
		first:     lo.Sum(data.firstVisits),
		returning: lo.Sum(data.returnVisits),
		bounce:    data.currentStats.BounceRateValue,
		duration:  data.currentStats.SessionDurationSeconds,
	}
	previous := periodComparisonRow{
		label:     "Previous Period",
		unique:    data.previousPeriodUnique,
		total:     data.previousPeriodTotal,
		first:     data.previousPeriodFirst,
		returning: data.previousPeriodReturning,
		bounce:    data.previousStats.BounceRateValue,
		duration:  data.previousStats.SessionDurationSeconds,
	}

	bodyRows := []hb.TagInterface{
		comparisonCountRow("Total Unique Visitors", current.unique, previous.unique),
		comparisonCountRow("Total Visitors", current.total, previous.total),
		comparisonCountRow("First Time Visits", current.first, previous.first),
		comparisonCountRow("Returning Visits", current.returning, previous.returning),
		comparisonFloatRowInverted("Bounce Rate", current.bounce, previous.bounce, func(f float64) string { return formatFloat2(f) + "%" }),
		comparisonFloatRow("Avg. Visit Duration", current.duration, previous.duration, formatDuration),
	}

	headRow := hb.TR().
		Child(hb.TH().Text("Metric")).
		Child(hb.TH().Class("text-end").Text(current.label)).
		Child(hb.TH().Class("text-end").Text(previous.label)).
		Child(hb.TH().Class("text-end").Text("Change"))

	table := hb.Table().
		Class("table table-sm table-hover align-middle").
		Child(hb.Thead().Class("table-light").Child(headRow)).
		Child(hb.Tbody().Children(bodyRows))

	return hb.Div().
		Class("mb-4").
		Child(hb.Heading5().Class("mb-3").Text("Period Comparison")).
		Child(hb.Div().Class("table-responsive").Child(table)).
		Child(hb.Small().Class("text-muted").Text("Comparing against " + data.previousPeriodLabel + "."))
}

type periodComparisonRow struct {
	label     string
	unique    int64
	total     int64
	first     int64
	returning int64
	bounce    float64
	duration  float64
}

func comparisonCountRow(label string, current, previous int64) hb.TagInterface {
	return hb.TR().Children([]hb.TagInterface{
		hb.TD().Text(label),
		hb.TD().Class("text-end").Text(formatCount(current)),
		hb.TD().Class("text-end").Text(formatCount(previous)),
		hb.TD().Class("text-end").Child(changeBadgePercent(changePercentInt(current, previous))),
	})
}

func comparisonFloatRow(label string, current, previous float64, formatter func(float64) string) hb.TagInterface {
	return hb.TR().Children([]hb.TagInterface{
		hb.TD().Text(label),
		hb.TD().Class("text-end").Text(formatter(current)),
		hb.TD().Class("text-end").Text(formatter(previous)),
		hb.TD().Class("text-end").Child(changeBadgePercent(changePercentFloat(current, previous))),
	})
}

func comparisonFloatRowInverted(label string, current, previous float64, formatter func(float64) string) hb.TagInterface {
	return hb.TR().Children([]hb.TagInterface{
		hb.TD().Text(label),
		hb.TD().Class("text-end").Text(formatter(current)),
		hb.TD().Class("text-end").Text(formatter(previous)),
		hb.TD().Class("text-end").Child(changeBadgePercentInverted(changePercentFloat(current, previous))),
	})
}

func changePercentInt(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return float64(current-previous) / float64(previous) * 100
}

func changePercentFloat(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return (current - previous) / previous * 100
}

func changeBadgePercent(percent float64) hb.TagInterface {
	if percent > 0 {
		return hb.Span().Class("badge bg-success").HTML(fmt.Sprintf(`<i class="bi bi-arrow-up"></i> %.1f%%`, percent))
	}
	if percent < 0 {
		return hb.Span().Class("badge bg-danger").HTML(fmt.Sprintf(`<i class="bi bi-arrow-down"></i> %.1f%%`, -percent))
	}
	return hb.Span().Class("badge bg-light text-dark").Text("0.0%")
}

func changeBadgePercentInverted(percent float64) hb.TagInterface {
	if percent > 0 {
		return hb.Span().Class("badge bg-danger").HTML(fmt.Sprintf(`<i class="bi bi-arrow-up"></i> %.1f%%`, percent))
	}
	if percent < 0 {
		return hb.Span().Class("badge bg-success").HTML(fmt.Sprintf(`<i class="bi bi-arrow-down"></i> %.1f%%`, -percent))
	}
	return hb.Span().Class("badge bg-light text-dark").Text("0.0%")
}
