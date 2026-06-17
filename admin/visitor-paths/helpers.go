package visitorpaths

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dracory/statsstore"
)

// buildControllerData builds the controller data
func buildControllerData(r *http.Request, store statsstore.StoreInterface) (visitorPathsControllerData, string) {
	data := visitorPathsControllerData{Request: r}

	query := r.URL.Query()
	page := parseIntWithDefault(query.Get("page"), 1)
	perPage := clampPerPage(parseIntWithDefault(query.Get("per_page"), 10))
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
	if filters.PathContains != "" {
		options = options.SetPathContains(filters.PathContains)
	}
	if filters.PathExact != "" {
		options = options.SetPathExact(filters.PathExact)
	}
	if filters.Device != "" {
		options = options.SetDeviceType(filters.Device)
	}

	paths, err := store.VisitorList(r.Context(), options)
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
	if filters.PathContains != "" {
		countOptions = countOptions.SetPathContains(filters.PathContains)
	}
	if filters.PathExact != "" {
		countOptions = countOptions.SetPathExact(filters.PathExact)
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

	data.Paths = paths
	data.Page = page
	data.PageSize = perPage
	data.TotalPages = totalPages
	data.TotalCount = totalCount
	data.Filters = filters

	return data, ""
}

func parseFilters(values url.Values) FilterOptions {
	filters := FilterOptions{
		Range:        strings.TrimSpace(values.Get("range")),
		From:         strings.TrimSpace(values.Get("from")),
		To:           strings.TrimSpace(values.Get("to")),
		Country:      strings.TrimSpace(values.Get("country")),
		PathContains: strings.TrimSpace(values.Get("path_contains")),
		PathExact:    strings.TrimSpace(values.Get("path_exact")),
		Device:       strings.TrimSpace(values.Get("device")),
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
