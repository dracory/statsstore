package visitoractivity

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dracory/statsstore"
	"github.com/gouniverse/hb"
)

// Data helpers

func buildControllerData(r *http.Request, store statsstore.StoreInterface) (ControllerData, string) {
	data := ControllerData{Request: r}

	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	perPage := 10
	offset := (pageInt - 1) * perPage

	visitors, err := store.VisitorList(r.Context(), statsstore.VisitorQueryOptions{
		Limit:     perPage,
		Offset:    offset,
		OrderBy:   statsstore.COLUMN_CREATED_AT,
		SortOrder: "DESC",
	})
	if err != nil {
		return data, err.Error()
	}

	visitorCount, err := store.VisitorCount(r.Context(), statsstore.VisitorQueryOptions{})
	if err != nil {
		return data, err.Error()
	}

	totalPages := (int(visitorCount) + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	data.Visitors = visitors
	data.Page = pageInt
	data.TotalPages = totalPages

	return data, ""
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
