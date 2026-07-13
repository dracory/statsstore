package home

import (
	"sort"
	"strings"
	"time"

	"github.com/dracory/statsstore"
	"github.com/dromara/carbon/v2"
)

// trafficSourceEntry represents a single row in a traffic source breakdown
// table (e.g. top referrers, top pages, top browsers).
type trafficSourceEntry struct {
	Label    string
	Sessions string
}

// trafficSourcesData holds all computed traffic source breakdowns derived
// from the actual visitor records for the selected period.
type trafficSourcesData struct {
	Referrers []trafficSourceEntry
	Pages     []trafficSourceEntry
	Browsers  []trafficSourceEntry
	Countries []trafficSourceEntry
	Events    []trafficSourceEntry
	Heatmap   weeklyHeatmapData
}

// weeklyHeatmapData holds the computed weekly trends heatmap.
type weeklyHeatmapData struct {
	Slots       []string
	Days        []string
	Intensities [][]int
}

// computeTrafficSources derives all traffic-source breakdowns from the
// visitor list stored in ControllerData.  Every number shown in the admin
// dashboard is computed from real records – nothing is hardcoded.
func computeTrafficSources(data ControllerData) trafficSourcesData {
	visitors := data.visitors

	referrerCounts := map[string]int64{}
	pageCounts := map[string]int64{}
	browserCounts := map[string]int64{}
	countryCounts := map[string]int64{}

	for _, v := range visitors {
		referrer := normalizeReferrer(v.GetUserReferrer())
		referrerCounts[referrer]++

		page := v.GetPath()
		if page == "" {
			page = "/"
		}
		pageCounts[page]++

		browser := strings.TrimSpace(v.GetUserBrowser())
		if browser == "" {
			browser = "Unknown"
		}
		browserCounts[browser]++

		country := strings.ToUpper(strings.TrimSpace(v.GetCountry()))
		if country == "" || country == "UN" || country == "ZZ" {
			country = "Unknown"
		}
		if data.ui.CountryNameByIso2 != nil && country != "Unknown" {
			if name, err := data.ui.CountryNameByIso2(country); err == nil && name != "" {
				country = name
			}
		}
		countryCounts[country]++
	}

	return trafficSourcesData{
		Referrers: topEntries(referrerCounts, 10),
		Pages:     topEntries(pageCounts, 10),
		Browsers:  topEntries(browserCounts, 10),
		Countries: topEntries(countryCounts, 10),
		Events:    computeEvents(visitors),
		Heatmap:   computeHeatmap(visitors),
	}
}

// topEntries converts a count map into a sorted slice of trafficSourceEntry
// capped at maxItems, sorted by count descending.
func topEntries(counts map[string]int64, maxItems int) []trafficSourceEntry {
	rawEntries := make([]rawEntry, 0, len(counts))
	for label, count := range counts {
		rawEntries = append(rawEntries, rawEntry{label: label, count: count})
	}
	sort.Slice(rawEntries, func(i, j int) bool {
		if rawEntries[i].count == rawEntries[j].count {
			return rawEntries[i].label < rawEntries[j].label
		}
		return rawEntries[i].count > rawEntries[j].count
	})

	entries := make([]trafficSourceEntry, 0, maxItems)
	for i := 0; i < len(rawEntries) && i < maxItems; i++ {
		entries = append(entries, trafficSourceEntry{
			Label:    rawEntries[i].label,
			Sessions: formatCount(rawEntries[i].count),
		})
	}
	return entries
}

type rawEntry struct {
	label string
	count int64
}

// formatCount converts an int64 count into a human-readable string (e.g.
// 4000 -> "4K").
func formatCount(n int64) string {
	if n >= 1000000 {
		return formatDecimal(float64(n)/1000000) + "M"
	}
	if n >= 1000 {
		return formatDecimal(float64(n)/1000) + "K"
	}
	return formatInt(n)
}

func formatDecimal(f float64) string {
	rounded := float64(int(f*10)) / 10
	if rounded == float64(int(rounded)) {
		return formatInt(int64(rounded))
	}
	s := formatFloat(rounded)
	return s
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

func formatFloat(f float64) string {
	// Simple float formatting with 1 decimal place
	s := formatInt(int64(f))
	remainder := f - float64(int64(f))
	if remainder < 0 {
		remainder = -remainder
	}
	decDigit := int(remainder*10 + 0.5)
	if decDigit == 10 {
		decDigit = 0
		// Round up the integer part
		s = formatInt(int64(f) + 1)
	}
	return s + "." + string(byte('0'+decDigit))
}

// normalizeReferrer extracts the domain from a referrer URL and strips
// common prefixes.
func normalizeReferrer(referrer string) string {
	referrer = strings.TrimSpace(referrer)
	if referrer == "" {
		return "(Direct / None)"
	}
	// Strip scheme
	for _, prefix := range []string{"https://", "http://", "www."} {
		referrer = strings.TrimPrefix(referrer, prefix)
	}
	// Keep only the domain part
	if idx := strings.Index(referrer, "/"); idx > 0 {
		referrer = referrer[:idx]
	}
	return referrer
}

// computeEvents extracts custom event labels from visitor paths that look
// like event endpoints (e.g. /event/signup).  Visitors without an event
// pattern are ignored.
func computeEvents(visitors []statsstore.VisitorInterface) []trafficSourceEntry {
	eventCounts := map[string]int64{}
	for _, v := range visitors {
		path := strings.TrimSpace(v.GetPath())
		if path == "" {
			continue
		}
		// Detect event-style paths: /event/<name> or /track/<name>
		lower := strings.ToLower(path)
		for _, prefix := range []string{"/event/", "/track/"} {
			if strings.HasPrefix(lower, prefix) {
				name := strings.TrimSpace(path[len(prefix):])
				if name == "" {
					name = "unnamed"
				}
				eventCounts[name]++
				break
			}
		}
	}
	if len(eventCounts) == 0 {
		return []trafficSourceEntry{}
	}
	return topEntries(eventCounts, 10)
}

// computeHeatmap builds the weekly trends heatmap from visitor timestamps.
// It buckets visitors by day-of-week and 2-hour time slots, then normalises
// the counts to a 0-5 intensity scale.
func computeHeatmap(visitors []statsstore.VisitorInterface) weeklyHeatmapData {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	slots := []string{
		"1 AM", "3 AM", "5 AM", "7 AM", "9 AM", "11 AM",
		"1 PM", "3 PM", "5 PM", "7 PM", "9 PM", "11 PM",
	}

	// 12 slots x 7 days
	counts := make([][]int, len(slots))
	for i := range counts {
		counts[i] = make([]int, len(days))
	}

	// Map carbon weekday (0=Sunday, 1=Monday, ...) to our index (0=Monday)
	weekdayToIndex := func(weekday time.Weekday) int {
		// carbon: 0=Sunday, 1=Monday, ..., 6=Saturday
		// our index: 0=Monday, ..., 5=Saturday, 6=Sunday
		if weekday == time.Sunday {
			return 6
		}
		return int(weekday) - 1
	}

	// Map hour to slot index (each slot covers 2 hours: 0-1 -> slot 0, 2-3 -> slot 1, etc.)
	hourToSlot := func(hour int) int {
		return hour / 2
	}

	var maxCount int
	for _, v := range visitors {
		c := v.GetCreatedAtCarbon()
		if c == nil {
			continue
		}
		dayIdx := weekdayToIndex(c.StdTime().Weekday())
		slotIdx := hourToSlot(c.Hour())
		if dayIdx < 0 || dayIdx >= len(days) || slotIdx < 0 || slotIdx >= len(slots) {
			continue
		}
		counts[slotIdx][dayIdx]++
		if counts[slotIdx][dayIdx] > maxCount {
			maxCount = counts[slotIdx][dayIdx]
		}
	}

	// Normalise to 0-5 intensity scale
	intensities := make([][]int, len(slots))
	for i := range intensities {
		intensities[i] = make([]int, len(days))
		for j := range intensities[i] {
			if maxCount == 0 {
				intensities[i][j] = 0
				continue
			}
			ratio := float64(counts[i][j]) / float64(maxCount)
			intensities[i][j] = intensityLevel(ratio)
		}
	}

	return weeklyHeatmapData{
		Slots:       slots,
		Days:        days,
		Intensities: intensities,
	}
}

// intensityLevel maps a 0.0-1.0 ratio to a 0-5 intensity level.
func intensityLevel(ratio float64) int {
	switch {
	case ratio >= 0.85:
		return 5
	case ratio >= 0.65:
		return 4
	case ratio >= 0.45:
		return 3
	case ratio >= 0.25:
		return 2
	case ratio >= 0.05:
		return 1
	default:
		return 0
	}
}

// computeStatsOverview computes the extended statistics (sessions, pageviews,
// pages per session, bounce rate, session duration) from the visitor list.
type extendedStats struct {
	Sessions               string
	Pageviews              string
	PagesPerSession        string
	BounceRate             string
	BounceRateValue        float64
	SessionDuration        string
	SessionDurationSeconds float64
}

func computeStatsOverview(visitors []statsstore.VisitorInterface) extendedStats {
	totalPageviews := int64(len(visitors))

	// Group by fingerprint (calculating it on the fly if not stored) to
	// identify unique visitors/sessions for bounce rate and visit duration.
	sessions := map[string][]statsstore.VisitorInterface{}
	for _, v := range visitors {
		key := strings.TrimSpace(v.GetFingerprint())
		if key == "" {
			key = v.FingerprintCalculate()
		}
		if key == "" {
			key = "unknown"
		}
		sessions[key] = append(sessions[key], v)
	}

	sessionCount := int64(len(sessions))
	divisor := sessionCount
	if divisor == 0 {
		divisor = 1
	}

	pagesPerSession := float64(totalPageviews) / float64(divisor)

	var bounceSessions int64
	var totalIntervalSeconds float64
	var intervalCount int64

	for _, sessionVisitors := range sessions {
		if len(sessionVisitors) == 1 {
			bounceSessions++
			continue
		}

		times := make([]time.Time, 0, len(sessionVisitors))
		for _, sv := range sessionVisitors {
			c := sv.GetCreatedAtCarbon()
			if c != nil {
				times = append(times, c.StdTime())
			}
		}
		if len(times) < 2 {
			continue
		}

		sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
		for i := 1; i < len(times); i++ {
			interval := times[i].Sub(times[i-1]).Seconds()
			if interval > 0 {
				totalIntervalSeconds += interval
				intervalCount++
			}
		}
	}

	bounceRate := float64(0)
	if divisor > 0 {
		bounceRate = float64(bounceSessions) / float64(divisor) * 100
	}

	avgVisitDuration := float64(0)
	if intervalCount > 0 {
		avgVisitDuration = totalIntervalSeconds / float64(intervalCount)
	}

	return extendedStats{
		Sessions:               formatCount(sessionCount),
		Pageviews:              formatCount(totalPageviews),
		PagesPerSession:        formatFloat2(pagesPerSession),
		BounceRate:             formatFloat2(bounceRate) + "%",
		BounceRateValue:        bounceRate,
		SessionDuration:        formatDuration(avgVisitDuration),
		SessionDurationSeconds: avgVisitDuration,
	}
}

func formatFloat2(f float64) string {
	return formatFloat(f)
}

func formatDuration(seconds float64) string {
	if seconds <= 0 {
		return "0s"
	}
	totalSecs := int(seconds)
	hours := totalSecs / 3600
	minutes := (totalSecs % 3600) / 60
	secs := totalSecs % 60
	if hours > 0 {
		return formatInt(int64(hours)) + "h " + formatInt(int64(minutes)) + "m"
	}
	if minutes > 0 {
		return formatInt(int64(minutes)) + "m " + formatInt(int64(secs)) + "s"
	}
	return formatInt(int64(secs)) + "s"
}

// carbonNow returns the current carbon time – kept as a convenience to
// avoid importing carbon in files that don't already use it.
var _ = carbon.Now
