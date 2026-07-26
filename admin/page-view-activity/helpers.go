package pageviewactivity

import (
	"github.com/dracory/statsstore/admin/shared"
)

// splitTimestamp parses a created_at timestamp and returns separate date and time strings.
func splitTimestamp(value string) (date string, timeStr string) {
	if value == "" {
		return "Unknown", "Unknown"
	}
	t, err := shared.TimeParse(value)
	if err != nil {
		return value, ""
	}
	return t.Format("2006-01-02"), t.Format("15:04:05")
}
