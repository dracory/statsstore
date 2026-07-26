package visitoractivity

import (
	"net/http"
	"strings"
	"time"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

type visitorJSON struct {
	ID            string `json:"id"`
	CreatedAt     string `json:"createdAt"`
	VisitTime     string `json:"visitTime"`
	Path          string `json:"path"`
	Country       string `json:"country"`
	CountryCode   string `json:"countryCode"`
	CountryName   string `json:"countryName"`
	Location      string `json:"location"`
	IPAddress     string `json:"ipAddress"`
	Referrer      string `json:"referrer"`
	Device        string `json:"device"`
	DeviceIcon    string `json:"deviceIcon"`
	OSIcon        string `json:"osIcon"`
	SystemSummary string `json:"systemSummary"`
	SessionLabel  string `json:"sessionLabel"`
	Duration      string `json:"duration"`
	Browser       string `json:"browser"`
	BrowserVer    string `json:"browserVer"`
	OS            string `json:"os"`
	OSVer         string `json:"osVer"`
	UserAgent     string `json:"userAgent"`
	Fingerprint   string `json:"fingerprint"`
}

// handleListAjax returns the visitor list as JSON for the Vue.js frontend
func (c *Controller) handleListAjax(w http.ResponseWriter, r *http.Request) string {
	page := shared.ParseIntWithDefault(req.GetString(r, "page"), 1)
	perPage := shared.ClampPerPage(shared.ParseIntWithDefault(req.GetString(r, "per_page"), 10))

	filters := parseFiltersFromReq(r)

	offset := (page - 1) * perPage

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

	visitors, err := c.UI.Store.VisitorList(r.Context(), options)
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
	if filters.Device != "" {
		countOptions = countOptions.SetDeviceType(filters.Device)
	}

	visitorCount, err := c.UI.Store.VisitorCount(r.Context(), countOptions)
	if err != nil {
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	totalPages := int(visitorCount) / perPage
	if int(visitorCount)%perPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	visitorList := make([]visitorJSON, 0, len(visitors))
	for i, v := range visitors {
		vj := visitorJSON{
			ID:          v.GetID(),
			CreatedAt:   v.GetCreatedAt(),
			VisitTime:   formatVisitorTimestamp(v.GetCreatedAt()),
			Path:        v.GetPath(),
			Country:     v.GetCountry(),
			CountryCode: strings.ToUpper(strings.TrimSpace(v.GetCountry())),
			CountryName: shared.ResolvedCountryName(c.UI, v.GetCountry()),
			Location:    shared.FormatLocation(c.UI, v),
			IPAddress:   v.GetIpAddress(),
			Referrer:    v.GetUserReferrer(),
			Device:      v.GetUserDevice(),
			Browser:     v.GetUserBrowser(),
			BrowserVer:  v.GetUserBrowserVersion(),
			OS:          v.GetUserOs(),
			OSVer:       v.GetUserOsVersion(),
			UserAgent:   v.GetUserAgent(),
			Fingerprint: v.GetFingerprint(),
			Duration:    formatVisitDuration(v, visitors, i),
		}

		vj.DeviceIcon = deviceIconClass(v)
		vj.OSIcon = osIconClass(v)

		systemText := strings.TrimSpace(v.GetUserBrowser() + " " + v.GetUserBrowserVersion())
		if systemText == "" {
			systemText = "Unknown Browser"
		}
		osText := strings.TrimSpace(v.GetUserOs() + " " + v.GetUserOsVersion())
		if osText == "" {
			osText = "Unknown OS"
		}
		vj.SystemSummary = systemText + " on " + osText

		fp := v.GetFingerprint()
		if len(fp) > 8 {
			fp = fp[:8]
		}
		if fp == "" {
			fp = "Session"
		}
		vj.SessionLabel = strings.ToUpper(fp)

		visitorList = append(visitorList, vj)
	}

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"visitors":   visitorList,
		"page":       page,
		"totalPages": totalPages,
		"pageSize":   perPage,
		"totalCount": visitorCount,
	}))

	return ""
}

func parseFiltersFromReq(r *http.Request) FilterOptions {
	get := func(key string) string {
		return strings.TrimSpace(req.GetString(r, key))
	}

	filters := FilterOptions{
		Range:   get("range"),
		Country: get("country"),
		Device:  get("device"),
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

func deviceIconClass(v statsstore.VisitorInterface) string {
	deviceType := strings.ToLower(v.GetUserDeviceType())
	switch {
	case strings.Contains(deviceType, "desktop"):
		return "bi bi-display text-primary"
	case strings.Contains(deviceType, "mobile"):
		return "bi bi-phone text-success"
	case strings.Contains(deviceType, "tablet"):
		return "bi bi-tablet text-info"
	case strings.Contains(deviceType, "bot"):
		return "bi bi-robot text-warning"
	default:
		return "bi bi-question-circle text-secondary"
	}
}

func osIconClass(v statsstore.VisitorInterface) string {
	os := strings.ToLower(v.GetUserOs())
	switch {
	case strings.Contains(os, "windows"):
		return "bi bi-windows text-primary"
	case strings.Contains(os, "mac"), strings.Contains(os, "ios"):
		return "bi bi-apple text-dark"
	case strings.Contains(os, "android"):
		return "bi bi-android2 text-success"
	case strings.Contains(os, "linux"):
		return "bi bi-ubuntu text-warning"
	default:
		return "bi bi-circle text-secondary"
	}
}
