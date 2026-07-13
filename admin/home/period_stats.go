package home

import (
	"github.com/dracory/statsstore"
)

// periodStats holds the aggregated daily and total metrics for a date range.
type periodStats struct {
	dates          []string
	uniqueVisits   []int64
	totalVisits    []int64
	firstVisits    []int64
	returnVisits   []int64
	totalUnique    int64
	totalTotal     int64
	totalFirst     int64
	totalReturning int64
}

// computePeriodStats aggregates visitor records into daily and total counts for
// the supplied date range. Unique visitors are identified by IP address.
func computePeriodStats(visitors []statsstore.VisitorInterface, dates []string) periodStats {
	dailyPageViews := map[string]int64{}
	dailyUniqueIPs := map[string]map[string]struct{}{}
	firstVisitByIP := map[string]string{}

	for _, visitor := range visitors {
		createdAt := visitor.GetCreatedAtCarbon()
		if createdAt == nil {
			continue
		}

		visitDate := createdAt.ToDateString()
		identifier := visitor.GetIpAddress()
		if identifier == "" {
			identifier = "unknown-ip"
		}

		dailyPageViews[visitDate]++

		if _, ok := dailyUniqueIPs[visitDate]; !ok {
			dailyUniqueIPs[visitDate] = map[string]struct{}{}
		}

		dailyUniqueIPs[visitDate][identifier] = struct{}{}

		if existingDate, ok := firstVisitByIP[identifier]; !ok || visitDate < existingDate {
			firstVisitByIP[identifier] = visitDate
		}
	}

	result := periodStats{dates: dates}

	for _, date := range dates {
		uniqueSet := dailyUniqueIPs[date]
		uniqueCount := int64(len(uniqueSet))

		var firstCount int64
		for ip := range uniqueSet {
			if firstVisitByIP[ip] == date {
				firstCount++
			}
		}

		returnCount := uniqueCount - firstCount
		if returnCount < 0 {
			returnCount = 0
		}

		result.uniqueVisits = append(result.uniqueVisits, uniqueCount)
		result.totalVisits = append(result.totalVisits, dailyPageViews[date])
		result.firstVisits = append(result.firstVisits, firstCount)
		result.returnVisits = append(result.returnVisits, returnCount)

		result.totalUnique += uniqueCount
		result.totalTotal += dailyPageViews[date]
		result.totalFirst += firstCount
		result.totalReturning += returnCount
	}

	return result
}
