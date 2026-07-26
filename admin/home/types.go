package home

import "encoding/json"

// statCardJSON represents a single stat card in the dashboard.
type statCardJSON struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// comparisonRowJSON represents a single row in the comparison table.
type comparisonRowJSON struct {
	Label    string  `json:"label"`
	Current  string  `json:"current"`
	Previous string  `json:"previous"`
	Change   float64 `json:"change"`
	Inverted bool    `json:"inverted"`
}

// dailyStatJSON represents a single day's stats in the daily table.
type dailyStatJSON struct {
	Date         string `json:"date"`
	TotalVisits  int64  `json:"totalVisits"`
	UniqueVisits int64  `json:"uniqueVisits"`
	FirstVisits  int64  `json:"firstVisits"`
	ReturnVisits int64  `json:"returnVisits"`
}

// totalsJSON represents the totals row in the daily table.
type totalsJSON struct {
	TotalVisits  int64 `json:"totalVisits"`
	UniqueVisits int64 `json:"uniqueVisits"`
	FirstVisits  int64 `json:"firstVisits"`
	ReturnVisits int64 `json:"returnVisits"`
}

// trafficCardJSON represents a traffic source card with tabs.
type trafficCardJSON struct {
	Title      string           `json:"title"`
	ValueLabel string           `json:"valueLabel"`
	Tabs       []trafficTabJSON `json:"tabs"`
}

// trafficTabJSON represents a single tab within a traffic card.
type trafficTabJSON struct {
	Label   string               `json:"label"`
	Entries []trafficSourceEntry `json:"entries"`
}

// heatmapJSON represents the weekly heatmap data.
type heatmapJSON struct {
	Days        []string `json:"days"`
	Slots       []string `json:"slots"`
	Intensities [][]int  `json:"intensities"`
}

// jsonMarshal marshals v to JSON string, returning an error JSON on failure.
func jsonMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return `{"error":"failed to marshal JSON"}`
	}
	return string(b)
}
