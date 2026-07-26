package visitorpaths

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

type pathJSON struct {
	ID               string `json:"id"`
	Path             string `json:"path"`
	AbsoluteURL      string `json:"absoluteUrl"`
	VisitTime        string `json:"visitTime"`
	Country          string `json:"country"`
	CountryCode      string `json:"countryCode"`
	CountryName      string `json:"countryName"`
	Location         string `json:"location"`
	IPAddress        string `json:"ipAddress"`
	Referrer         string `json:"referrer"`
	UserAgent        string `json:"userAgent"`
	SessionLabel     string `json:"sessionLabel"`
	DeviceLabel      string `json:"deviceLabel"`
	DeviceBadgeClass string `json:"deviceBadgeClass"`
	BrowserLabel     string `json:"browserLabel"`
	DrillDownURL     string `json:"drillDownUrl"`
}

// handleListAjax returns the visitor paths list as JSON for the Vue.js frontend
func (c *visitorPathsController) handleListAjax(w http.ResponseWriter, r *http.Request) string {
	page := shared.ParseIntWithDefault(req.GetString(r, "page"), 1)
	perPage := shared.ClampPerPage(shared.ParseIntWithDefault(req.GetString(r, "per_page"), 10))
	offset := (page - 1) * perPage

	filters := parseFiltersFromReq(r)

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

	visitors, err := c.ui.Store.VisitorList(r.Context(), options)
	if err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
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

	totalCount, err := c.ui.Store.VisitorCount(r.Context(), countOptions)
	if err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	totalPages := int(totalCount) / perPage
	if int(totalCount)%perPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	pathList := make([]pathJSON, 0, len(visitors))
	for _, v := range visitors {
		pj := pathJSON{
			ID:           v.GetID(),
			Path:         shared.StripMethodPrefix(v.GetPath()),
			AbsoluteURL:  shared.FullPathURL(c.ui, v.GetPath()),
			VisitTime:    shared.FormatTimestamp(v.GetCreatedAt()),
			Country:      v.GetCountry(),
			CountryCode:  strings.ToUpper(strings.TrimSpace(v.GetCountry())),
			CountryName:  shared.ResolvedCountryName(c.ui, v.GetCountry()),
			Location:     shared.FormatLocation(c.ui, v),
			IPAddress:    v.GetIpAddress(),
			Referrer:     v.GetUserReferrer(),
			UserAgent:    v.GetUserAgent(),
			SessionLabel: fmt.Sprintf("Sessions: %d", sessionCount(visitors, v)),
		}

		pj.DeviceLabel = strings.Title(strings.ToLower(v.GetUserDeviceType()))
		if pj.DeviceLabel == "" {
			pj.DeviceLabel = "Unknown"
		}
		pj.DeviceBadgeClass = deviceBadgeClass(v)

		browser := strings.TrimSpace(v.GetUserBrowser() + " " + v.GetUserBrowserVersion())
		if browser == "" {
			browser = "Unknown Browser"
		}
		pj.BrowserLabel = browser

		drillParams := map[string]string{"page": "1"}
		pj.DrillDownURL = shared.UrlVisitorActivity(r, drillParams)

		pathList = append(pathList, pj)
	}

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"paths":      pathList,
		"page":       page,
		"totalPages": totalPages,
		"pageSize":   perPage,
		"totalCount": totalCount,
	}))

	return ""
}

func parseFiltersFromReq(r *http.Request) FilterOptions {
	filters := FilterOptions{
		Range:        strings.TrimSpace(req.GetString(r, "range")),
		Country:      strings.TrimSpace(req.GetString(r, "country")),
		PathContains: strings.TrimSpace(req.GetString(r, "path_contains")),
		PathExact:    strings.TrimSpace(req.GetString(r, "path_exact")),
		Device:       strings.TrimSpace(req.GetString(r, "device")),
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

func deviceBadgeClass(v statsstore.VisitorInterface) string {
	deviceType := strings.ToLower(v.GetUserDeviceType())
	switch {
	case strings.Contains(deviceType, "desktop"):
		return "text-bg-primary"
	case strings.Contains(deviceType, "mobile"):
		return "text-bg-success"
	case strings.Contains(deviceType, "tablet"):
		return "text-bg-info"
	case strings.Contains(deviceType, "bot"):
		return "text-bg-warning"
	default:
		return "text-bg-secondary"
	}
}
