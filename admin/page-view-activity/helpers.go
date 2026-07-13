package pageviewactivity

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dracory/statsstore"
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

// buildControllerData prepares the controller data and returns an optional error message.
func buildControllerData(r *http.Request, store statsstore.StoreInterface) (
	data ControllerData,
	errMessage string,
) {
	data = ControllerData{Request: r}

	query := r.URL.Query()

	page := shared.ParseIntWithDefault(query.Get("page"), 1)
	perPage := shared.ClampPerPage(shared.ParseIntWithDefault(query.Get("per_page"), 10))
	offset := (page - 1) * perPage

	filters := parseFilters(query)

	options := statsstore.VisitorQuery().
		SetLimit(perPage).
		SetOffset(offset).
		SetOrderBy(statsstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC")

	if filters.Country != "" {
		options = options.SetCountry(filters.Country)
	}
	if filters.From != "" {
		options = options.SetCreatedAtGte(filters.From)
	}
	if filters.To != "" {
		options = options.SetCreatedAtLte(filters.To)
	}
	if filters.Device != "" {
		options = options.SetDeviceType(filters.Device)
	}
	if filters.Browser != "" {
		// Browser-specific filtering not yet supported at store level; left for future enhancement.
	}

	visitors, err := store.VisitorList(r.Context(), options)
	if err != nil {
		return data, err.Error()
	}

	countOptions := statsstore.VisitorQuery()
	if filters.Country != "" {
		countOptions = countOptions.SetCountry(filters.Country)
	}
	if filters.From != "" {
		countOptions = countOptions.SetCreatedAtGte(filters.From)
	}
	if filters.To != "" {
		countOptions = countOptions.SetCreatedAtLte(filters.To)
	}
	if filters.Device != "" {
		countOptions = countOptions.SetDeviceType(filters.Device)
	}

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
