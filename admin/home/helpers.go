package home

import "time"

// formatSummaryDate converts a "2006-01-02" date string into a human-readable
// format like "Mon, 2 Jan 2006".
func formatSummaryDate(date string) string {
	if parsed, err := time.Parse("2006-01-02", date); err == nil {
		return parsed.Format("Mon, 2 Jan 2006")
	}
	return date
}

// changePercentInt calculates the percentage change between two int64 values.
func changePercentInt(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return float64(current-previous) / float64(previous) * 100
}

// changePercentFloat calculates the percentage change between two float64 values.
func changePercentFloat(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return (current - previous) / previous * 100
}
