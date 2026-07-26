package pageviewactivity

import (
	"net/http"
	"strings"
	"time"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/statsstore"
	"github.com/dracory/statsstore/admin/shared"
)

type pageViewJSON struct {
	ID           string `json:"id"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Path         string `json:"path"`
	AbsoluteURL  string `json:"absoluteUrl"`
	Country      string `json:"country"`
	CountryCode  string `json:"countryCode"`
	CountryName  string `json:"countryName"`
	Location     string `json:"location"`
	IPAddress    string `json:"ipAddress"`
	Referrer     string `json:"referrer"`
	DeviceLabel  string `json:"deviceLabel"`
	BrowserLabel string `json:"browserLabel"`
	OSLabel      string `json:"osLabel"`
	Language     string `json:"language"`
	UserAgent    string `json:"userAgent"`
}

// handleListAjax returns the page view list as JSON for the Vue.js frontend
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

	totalCount, err := c.UI.Store.VisitorCount(r.Context(), countOptions)
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

	pageViewList := make([]pageViewJSON, 0, len(visitors))
	for _, v := range visitors {
		date, timeStr := splitTimestamp(v.GetCreatedAt())

		pv := pageViewJSON{
			ID:          v.GetID(),
			Date:        date,
			Time:        timeStr,
			Path:        v.GetPath(),
			AbsoluteURL: shared.FullPathURL(c.UI, v.GetPath()),
			Country:     v.GetCountry(),
			CountryCode: strings.ToUpper(strings.TrimSpace(v.GetCountry())),
			CountryName: shared.ResolvedCountryName(c.UI, v.GetCountry()),
			Location:    shared.FormatLocation(c.UI, v),
			IPAddress:   v.GetIpAddress(),
			Referrer:    v.GetUserReferrer(),
			Language:    v.GetUserAcceptLanguage(),
			UserAgent:   v.GetUserAgent(),
		}

		if pv.IPAddress == "" {
			pv.IPAddress = "Unknown"
		}
		if pv.Language == "" {
			pv.Language = "Unknown"
		}

		pv.DeviceLabel = strings.Title(strings.ToLower(v.GetUserDeviceType()))
		if pv.DeviceLabel == "" {
			pv.DeviceLabel = "Unknown"
		}

		browser := strings.TrimSpace(v.GetUserBrowser() + " " + v.GetUserBrowserVersion())
		if browser == "" {
			browser = "Unknown Browser"
		}
		pv.BrowserLabel = browser

		osLabel := strings.TrimSpace(v.GetUserOs() + " " + v.GetUserOsVersion())
		if osLabel == "" {
			osLabel = "Unknown OS"
		}
		pv.OSLabel = osLabel

		pageViewList = append(pageViewList, pv)
	}

	api.Respond(w, r, api.SuccessWithData("success", map[string]any{
		"pageViews":   pageViewList,
		"page":        page,
		"totalPages":  totalPages,
		"pageSize":    perPage,
		"totalCount":  totalCount,
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
