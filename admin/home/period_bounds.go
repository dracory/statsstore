package home

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dromara/carbon/v2"
)

type periodBoundsData struct {
	selectedPeriod   string
	periodOptions    []periodOption
	createdAtGte     string
	createdAtLte     string
	prevCreatedAtGte string
	prevCreatedAtLte string
	prevLabel        string
	dateRange        []string
	prevDateRange    []string
}

// getPeriodBounds computes the date bounds for the selected period without loading any data.
func (c *Controller) getPeriodBounds(r *http.Request) (periodBoundsData, string) {
	periodOptions := []periodOption{
		{Value: "today", Label: "Today"},
		{Value: "yesterday", Label: "Yesterday"},
		{Value: "last-7-days", Label: "Last 7 Days"},
		{Value: "previous-7-days", Label: "Previous 7 Days"},
		{Value: "this-week", Label: "This Week"},
		{Value: "last-week", Label: "Last Week"},
		{Value: "this-month", Label: "This Month"},
		{Value: "last-month", Label: "Last Month"},
	}

	selectedPeriod := req.GetStringOr(r, "period", "this-week")

	now := carbon.Now(carbon.UTC)
	start := now.Copy()
	end := now.Copy()

	switch selectedPeriod {
	case "today":
		start = now.Copy().StartOfDay()
		end = now.Copy().EndOfDay()
	case "yesterday":
		start = now.Copy().SubDays(1).StartOfDay()
		end = start.Copy().EndOfDay()
	case "last-7-days":
		start = now.Copy().SubDays(6).StartOfDay()
		end = now.Copy().EndOfDay()
	case "previous-7-days":
		end = now.Copy().SubDays(7).EndOfDay()
		start = end.Copy().SubDays(6).StartOfDay()
	case "last-week":
		start = now.SubWeeks(1).StartOfWeek()
		end = start.Copy().EndOfWeek()
	case "this-month":
		start = now.StartOfMonth()
		end = now.EndOfMonth()
	case "last-month":
		start = now.SubMonths(1).StartOfMonth()
		end = start.Copy().EndOfMonth()
	default:
		start = now.StartOfWeek()
		end = now.EndOfWeek()
	}

	dateRange := datesInRange(start.Copy(), end.Copy())
	createdAtGte := start.ToDateTimeString(carbon.UTC)
	createdAtLte := end.ToDateTimeString(carbon.UTC)

	prevStart, prevEnd, prevLabel := previousPeriodBounds(selectedPeriod, start, end)
	prevCreatedAtGte := prevStart.ToDateTimeString(carbon.UTC)
	prevCreatedAtLte := prevEnd.ToDateTimeString(carbon.UTC)
	prevDateRange := datesInRange(prevStart.Copy(), prevEnd.Copy())

	return periodBoundsData{
		selectedPeriod:   selectedPeriod,
		periodOptions:    periodOptions,
		createdAtGte:     createdAtGte,
		createdAtLte:     createdAtLte,
		prevCreatedAtGte: prevCreatedAtGte,
		prevCreatedAtLte: prevCreatedAtLte,
		prevLabel:        prevLabel,
		dateRange:        dateRange,
		prevDateRange:    prevDateRange,
	}, ""
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
