package home

import "github.com/dromara/carbon/v2"

// datesInRange returns an array of dates between the start and end dates
func datesInRange(timeStart, timeEnd *carbon.Carbon) []string {
	rangeDates := []string{}

	if timeStart.Lte(timeEnd) {
		rangeDates = append(rangeDates, timeStart.ToDateString())
		for timeStart.Lt(timeEnd) {
			timeStart = timeStart.AddDays(1)
			rangeDates = append(rangeDates, timeStart.ToDateString())
		}
	}

	return rangeDates
}
