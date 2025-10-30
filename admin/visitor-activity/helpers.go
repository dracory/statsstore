package visitoractivity

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
)

// Data helpers

func buildControllerData(r *http.Request, store statsstore.StoreInterface) (ControllerData, string) {
	data := ControllerData{Request: r}

	query := r.URL.Query()
	page := query.Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPageParam := query.Get("per_page")
	perPage := 10
	if perPageParam != "" {
		if val, errConv := strconv.Atoi(perPageParam); errConv == nil && val > 0 && val <= 100 {
			perPage = val
		}
	}

	filters := parseFilters(query)

	offset := (pageInt - 1) * perPage

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

	visitors, err := store.VisitorList(r.Context(), options)
	if err != nil {
		return data, err.Error()
	}

	countOptions := options
	countOptions.Limit = 0
	countOptions.Offset = 0
	countOptions.CountOnly = true

	visitorCount, err := store.VisitorCount(r.Context(), countOptions)
	if err != nil {
		return data, err.Error()
	}

	totalPages := int(visitorCount) / perPage
	if int(visitorCount)%perPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	data.Visitors = visitors
	data.Page = pageInt
	data.TotalPages = totalPages
	data.PageSize = perPage
	data.TotalCount = visitorCount
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

// Helper Functions

func formatVisitorTimestamp(timestamp string) string {
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		return t.Format("2006-01-02 15:04:05 -0700 UTC")
	}
	return timestamp
}

func formatVisitDuration(visitor statsstore.VisitorInterface, visitors []statsstore.VisitorInterface, index int) string {
	if index < len(visitors)-1 {
		nextVisit := visitors[index+1]
		t1, err1 := time.Parse(time.RFC3339, visitor.CreatedAt())
		t2, err2 := time.Parse(time.RFC3339, nextVisit.CreatedAt())
		if err1 == nil && err2 == nil {
			durationSec := t1.Sub(t2).Seconds()
			if durationSec > 0 {
				return fmt.Sprintf("%.0f seconds", durationSec)
			}
		}
	}
	return "-"
}

func deviceIcon(visitor statsstore.VisitorInterface) hb.TagInterface {
	deviceType := strings.ToLower(visitor.UserDeviceType())

	iconClass := "bi bi-question-circle"
	color := "text-secondary"

	switch {
	case strings.Contains(deviceType, "desktop"):
		iconClass = "bi bi-display"
		color = "text-primary"
	case strings.Contains(deviceType, "mobile"):
		iconClass = "bi bi-phone"
		color = "text-success"
	case strings.Contains(deviceType, "tablet"):
		iconClass = "bi bi-tablet"
		color = "text-info"
	case strings.Contains(deviceType, "bot"):
		iconClass = "bi bi-robot"
		color = "text-warning"
	}

	return hb.I().Class(iconClass+" "+color).Attr("title", visitor.UserDevice())
}

// osIcon returns an icon representing the operating system
func osIcon(visitor statsstore.VisitorInterface) hb.TagInterface {
	os := strings.ToLower(visitor.UserOs())

	iconClass := "bi bi-circle"
	color := "text-secondary"

	switch {
	case strings.Contains(os, "windows"):
		iconClass = "bi bi-windows"
		color = "text-primary"
	case strings.Contains(os, "mac"), strings.Contains(os, "ios"):
		iconClass = "bi bi-apple"
		color = "text-dark"
	case strings.Contains(os, "android"):
		iconClass = "bi bi-android2"
		color = "text-success"
	case strings.Contains(os, "linux"):
		iconClass = "bi bi-ubuntu"
		color = "text-warning"
	}

	return hb.I().
		Class(iconClass + " " + color).
		Title(visitor.UserOs() + " " + visitor.UserOsVersion())
}

// getVisitPageLink returns HTML for visit page with link
func getVisitPageLink(path string) string {
	if path == "" {
		path = "/"
	}

	if strings.Contains(path, "snippet/55") {
		return `<a href="https://www.kuikie.com/snippet/55/cpp-how-to-check-if-a-qstring-is-base64-encoded" target="_blank" 
			style="color: #28a745; text-decoration: none; display: flex; align-items: center;">
			https://www.kuikie.com/snippet/55/cpp-how-to-check-if-a-qstring-is-base64-encoded
			<i class="bi bi-box-arrow-up-right" style="margin-left: 5px;"></i>
		</a>`
	}

	siteURL := "https://www.example.com" // Replace with actual site URL
	return `<a href="` + siteURL + path + `" target="_blank" 
		style="color: #28a745; text-decoration: none; display: flex; align-items: center;">
		` + path + `
		<i class="bi bi-box-arrow-up-right" style="margin-left: 5px;"></i>
	</a>`
}
