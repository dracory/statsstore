package visitorpaths

import (
	"strings"

	"github.com/dracory/statsstore"
)

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
