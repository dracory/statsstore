package visitoractivity

import (
	"fmt"
	"time"

	"github.com/dracory/statsstore"
)

func formatVisitorTimestamp(timestamp string) string {
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		return t.Format("2006-01-02 15:04:05 -0700 UTC")
	}
	return timestamp
}

func formatVisitDuration(visitor statsstore.VisitorInterface, visitors []statsstore.VisitorInterface, index int) string {
	if index < len(visitors)-1 {
		nextVisit := visitors[index+1]
		t1, err1 := time.Parse(time.RFC3339, visitor.GetCreatedAt())
		t2, err2 := time.Parse(time.RFC3339, nextVisit.GetCreatedAt())
		if err1 == nil && err2 == nil {
			durationSec := t1.Sub(t2).Seconds()
			if durationSec > 0 {
				return fmt.Sprintf("%.0f seconds", durationSec)
			}
		}
	}
	return "-"
}
