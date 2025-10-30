package pageviewactivity

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dracory/statsstore"
)

// buildControllerData prepares the controller data and returns an optional error message.
func buildControllerData(r *http.Request, store statsstore.StoreInterface) (ControllerData, string) {
	data := ControllerData{Request: r}

	query := r.URL.Query()

	page := parseIntWithDefault(query.Get("page"), 1)
	perPage := clampPerPage(parseIntWithDefault(query.Get("per_page"), 10))
	offset := (page - 1) * perPage

	filters := parseFilters(query)

	options := statsstore.VisitorQueryOptions{
		Limit:     perPage,
		Offset:    offset,
		OrderBy:   statsstore.COLUMN_CREATED_AT,
		SortOrder: "DESC",
	}

	if filters.Country != "" {
		options.Country = filters.Country
	}
	if filters.From != "" {
		options.CreatedAtGte = filters.From
	}
	if filters.To != "" {
		options.CreatedAtLte = filters.To
	}
	if filters.Device != "" {
		options.DeviceType = filters.Device
	}
	if filters.Browser != "" {
		// Browser-specific filtering not yet supported at store level; left for future enhancement.
	}

	visitors, err := store.VisitorList(r.Context(), options)
	if err != nil {
		return data, err.Error()
	}

	countOptions := options
	countOptions.Limit = 0
	countOptions.Offset = 0
	countOptions.CountOnly = true

	totalCount, err := store.VisitorCount(r.Context(), countOptions)
	if err != nil {
		return data, err.Error()
	}

	totalPages := int(totalCount) / perPage
	if int(totalCount)%perPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	data.Visitors = visitors
	data.Page = page
	data.TotalPages = totalPages
	data.PageSize = perPage
	data.TotalCount = totalCount
	data.Filters = filters

	return data, ""
}

func parseFilters(values url.Values) FilterOptions {
	get := func(key string) string {
		return strings.TrimSpace(values.Get(key))
	}

	filters := FilterOptions{
		Range:   get("range"),
		From:    get("from"),
		To:      get("to"),
		Country: get("country"),
		Device:  get("device"),
		Browser: get("browser"),
	}

	if filters.Range != "" {
		now := time.Now().UTC()
		switch strings.ToLower(filters.Range) {
		case "24h", "last24hours", "last_24_hours":
			filters.From = now.Add(-24 * time.Hour).Format(time.RFC3339)
			filters.To = now.Format(time.RFC3339)
		case "today":
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
			filters.From = start.Format(time.RFC3339)
			filters.To = start.Add(24 * time.Hour).Format(time.RFC3339)
		case "7d", "last7days":
			filters.From = now.Add(-7 * 24 * time.Hour).Format(time.RFC3339)
			filters.To = now.Format(time.RFC3339)
		case "30d", "last30days":
			filters.From = now.Add(-30 * 24 * time.Hour).Format(time.RFC3339)
			filters.To = now.Format(time.RFC3339)
		}
	}

	return filters
}

func parseIntWithDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
		return parsed
	}
	return defaultValue
}

func clampPerPage(perPage int) int {
	switch {
	case perPage < 1:
		return 10
	case perPage > 100:
		return 100
	default:
		return perPage
	}
}
