package shared

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dracory/hb"
	"github.com/dracory/statsstore"
)

// == QUERY & PAGINATION HELPERS ================================================

// QueryParamsWith takes the current request's query parameters, applies the
// given overrides (empty values delete the key), and returns a map suitable
// for passing to URL builder functions.
func QueryParamsWith(r *http.Request, overrides map[string]string) map[string]string {
	values := url.Values{}
	for key, vals := range r.URL.Query() {
		for _, v := range vals {
			values.Add(key, v)
		}
	}

	for key, val := range overrides {
		if val == "" {
			values.Del(key)
			continue
		}
		values.Set(key, val)
	}

	result := map[string]string{}
	for key := range values {
		result[key] = values.Get(key)
	}
	return result
}

// ParseIntWithDefault parses a string to int, returning the default on error
// or non-positive values.
func ParseIntWithDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
		return parsed
	}
	return defaultValue
}

// ClampPerPage constrains a per-page value to the 1–100 range.
func ClampPerPage(perPage int) int {
	switch {
	case perPage < 1:
		return 10
	case perPage > 100:
		return 100
	default:
		return perPage
	}
}

// == FORMATTING HELPERS ========================================================

// RangeLabel converts a range code to a human-readable label.
func RangeLabel(value string) string {
	switch strings.ToLower(value) {
	case "24h", "last24hours", "last_24_hours":
		return "Last 24 Hours"
	case "today":
		return "Today"
	case "7d", "last7days":
		return "Last 7 Days"
	case "30d", "last30days":
		return "Last 30 Days"
	default:
		return value
	}
}

// ShortDate parses a timestamp and returns a YYYY-MM-DD string.
func ShortDate(value string) string {
	if value == "" {
		return "-"
	}
	if t, err := TimeParse(value); err == nil {
		return t.Format("2006-01-02")
	}
	return value
}

// TimeParse attempts to parse a timestamp using several common layouts.
func TimeParse(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse timestamp: %s", value)
}

// FormatTimestamp parses a timestamp and returns a user-friendly formatted string.
func FormatTimestamp(value string) string {
	if value == "" {
		return "Unknown"
	}
	if t, err := TimeParse(value); err == nil {
		return t.Format("Mon, 02 Jan 2006, 15:04")
	}
	return value
}

// == COUNTRY / LOCATION HELPERS ================================================

// countryNameCache caches ISO2->name lookups at the package level so that
// country name resolution is a single map lookup after the first call,
// regardless of which controller triggers it. Country data is static.
type countryNameCache struct {
	mu   sync.RWMutex
	data map[string]string
}

var globalCountryNameCache = &countryNameCache{data: map[string]string{}}

// ResolvedCountryName resolves a country code to a human-readable name,
// falling back to the ISO code or "Unknown". Results are cached package-wide
// so the CountryNameByIso2 callback is invoked at most once per ISO code.
func ResolvedCountryName(ui ControllerOptions, code string) string {
	trimmed := strings.TrimSpace(code)
	if trimmed == "" {
		return "Unknown"
	}
	iso := strings.ToUpper(trimmed)

	if iso == "UN" || iso == "ZZ" {
		return "Unknown"
	}

	globalCountryNameCache.mu.RLock()
	if name, ok := globalCountryNameCache.data[iso]; ok {
		globalCountryNameCache.mu.RUnlock()
		return name
	}
	globalCountryNameCache.mu.RUnlock()

	name := iso
	if ui.CountryNameByIso2 != nil {
		if resolved, err := ui.CountryNameByIso2(iso); err == nil && resolved != "" {
			name = resolved
		}
		globalCountryNameCache.mu.Lock()
		globalCountryNameCache.data[iso] = name
		globalCountryNameCache.mu.Unlock()
	}

	return name
}

// FormatLocation returns a human-readable location string for a visitor.
func FormatLocation(ui ControllerOptions, visitor statsstore.VisitorInterface) string {
	name := ResolvedCountryName(ui, visitor.GetCountry())
	if name == "" || name == "Unknown" {
		return "Unknown Location"
	}
	return name
}

// == URL HELPERS ===============================================================

// StripMethodPrefix removes any "[METHOD] " prefix from a path string.
// Some consumers store paths like "[GET] /shop/product/1" — this returns "/shop/product/1".
func StripMethodPrefix(path string) string {
	if idx := strings.Index(path, "] "); idx >= 0 && strings.HasPrefix(path, "[") {
		return strings.TrimSpace(path[idx+2:])
	}
	return path
}

// FullPathURL builds an absolute URL from a website base URL and a path.
func FullPathURL(ui ControllerOptions, path string) string {
	base := ui.WebsiteUrl
	if base == "" {
		return path
	}

	path = StripMethodPrefix(path)

	u, err := url.Parse(base)
	if err != nil {
		return base + path
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u.Path = path
	return u.String()
}

// PathLink renders an anchor tag with an external icon pointing to the full
// absolute URL for the given path.
func PathLink(ui ControllerOptions, path string) hb.TagInterface {
	absolute := FullPathURL(ui, path)
	display := path
	if display == "" {
		display = "/"
	}

	return hb.A().
		Href(absolute).
		Class("text-success text-decoration-none d-inline-flex align-items-center gap-1").
		Attr("target", "_blank").
		HTML(fmt.Sprintf("%s <i class=\"bi bi-box-arrow-up-right\"></i>", display))
}

// == INFO LINE HELPERS =========================================================

// InfoLine renders a labelled value line with a muted label and bold value.
func InfoLine(label string, value hb.TagInterface) hb.TagInterface {
	labelTag := hb.Span().
		Class("text-muted text-uppercase fw-semibold small flex-shrink-0").
		Text(label)

	valueTag := hb.Div().
		Class("text-body fw-semibold text-break").
		Child(value)

	return hb.Div().
		Class("d-flex gap-2 align-items-baseline lh-sm").
		Child(labelTag).
		Child(valueTag)
}

// InfoText renders a plain body-text span.
func InfoText(text string) hb.TagInterface {
	return hb.Span().Class("text-body").Text(text)
}

// InfoMuted renders a muted, italic span — used for placeholder values.
func InfoMuted(text string) hb.TagInterface {
	return hb.Span().Class("text-muted fst-italic").Text(text)
}

// == BADGES ====================================================================

// CountryBadge renders a badge with a CSS flag icon and country name tooltip.
func CountryBadge(ui ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface {
	code := strings.ToUpper(strings.TrimSpace(visitor.GetCountry()))
	name := ResolvedCountryName(ui, visitor.GetCountry())

	badge := hb.Span().
		Class("badge bg-light text-dark border d-inline-flex align-items-center gap-1")

	if code != "" && len(code) == 2 {
		badge = badge.Child(hb.Span().Class("fi fi-" + strings.ToLower(code)))
	} else {
		badge = badge.Text("🌐")
	}

	if name != "" {
		badge = badge.Attr("title", name)
	}

	return badge
}

// DeviceBadge renders a contextual badge for the visitor's device type.
func DeviceBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	deviceType := strings.ToLower(visitor.GetUserDeviceType())
	label := visitor.GetUserDeviceType()
	if label == "" {
		label = "Unknown"
	}

	classes := "badge bg-light text-dark border"
	switch {
	case strings.Contains(deviceType, "desktop"):
		classes = "badge bg-primary-subtle text-primary"
	case strings.Contains(deviceType, "mobile"):
		classes = "badge bg-success-subtle text-success"
	case strings.Contains(deviceType, "tablet"):
		classes = "badge bg-info-subtle text-info"
	case strings.Contains(deviceType, "bot"):
		classes = "badge bg-warning-subtle text-warning"
	}

	return hb.Span().Class(classes).Text(strings.Title(label))
}

// BrowserBadge renders a badge showing the visitor's browser and version.
func BrowserBadge(visitor statsstore.VisitorInterface) hb.TagInterface {
	browser := strings.TrimSpace(visitor.GetUserBrowser() + " " + visitor.GetUserBrowserVersion())
	if browser == "" {
		browser = "Unknown Browser"
	}

	return hb.Span().
		Class("badge bg-light text-dark border").
		Text(browser)
}

// == SHARED UI COMPONENTS ======================================================

// OptionsButton renders a gear-icon button placeholder used in card headers.
func OptionsButton() hb.TagInterface {
	return hb.Button().
		Class("btn btn-sm btn-outline-secondary").
		Attr("type", "button").
		HTML(`<i class="bi bi-gear"></i>`)
}

// ExportDropdown renders a dropdown with a link to the given export URL.
func ExportDropdown(exportURL string) hb.TagInterface {
	button := hb.Button().
		Class("btn btn-sm btn-outline-secondary dropdown-toggle").
		Attr("type", "button").
		Attr("data-bs-toggle", "dropdown").
		Attr("aria-expanded", "false").
		Text("Export")

	item := hb.A().
		Class("dropdown-item").
		Href(exportURL).
		Attr("target", "_blank").
		Attr("rel", "noopener").
		Text("Export to CSV")

	menu := hb.UL().
		Class("dropdown-menu").
		Child(hb.LI().Child(item))

	return hb.Div().
		Class("dropdown").
		Child(button).
		Child(menu)
}

// QuickRangeButtons renders a button group of common time-range filters.
// The urlFunc parameter receives the query-params map and should return the
// fully qualified URL for that range.
func QuickRangeButtons(r *http.Request, urlFunc func(params map[string]string) string) hb.TagInterface {
	btn := func(label, rng string) hb.TagInterface {
		params := map[string]string{"page": "1", "from": "", "to": ""}
		if rng != "" {
			params["range"] = rng
		}
		return hb.A().
			Class("btn btn-sm btn-outline-secondary").
			Href(urlFunc(QueryParamsWith(r, params))).
			Text(label)
	}

	return hb.Div().
		Class("btn-group").
		Attr("role", "group").
		Child(btn("All", "")).
		Child(btn("Last 24 Hours", "24h")).
		Child(btn("Today", "today")).
		Child(btn("Last 7 Days", "7d"))
}

// PerPageSelector renders a button group for selecting rows per page.
// The urlFunc parameter receives the query-params map and should return the
// fully qualified URL for that per-page value.
func PerPageSelector(r *http.Request, currentPageSize int, urlFunc func(params map[string]string) string) hb.TagInterface {
	options := []int{10, 25, 50, 100}
	group := hb.Div().Class("d-flex align-items-center gap-2")
	group = group.Child(hb.Span().Class("small text-muted").Text("Rows per page:"))

	buttons := hb.Div().Class("btn-group btn-group-sm")
	for _, size := range options {
		params := QueryParamsWith(r, map[string]string{"per_page": fmt.Sprintf("%d", size), "page": "1"})
		classes := "btn btn-outline-secondary"
		if currentPageSize == size {
			classes = "btn btn-secondary"
		}
		buttons = buttons.Child(hb.A().
			Class(classes).
			Href(urlFunc(params)).
			Text(fmt.Sprintf("%d", size)))
	}

	return group.Child(buttons)
}

// PaginationSummary renders a "Showing X-Y of Z" summary text.
func PaginationSummary(page, pageSize int, totalCount int64, itemLabel string) hb.TagInterface {
	if totalCount == 0 {
		return hb.Span().Class("text-muted small").Text(fmt.Sprintf("No %s to display", itemLabel))
	}

	start := (page-1)*pageSize + 1
	end := page * pageSize
	if int64(end) > totalCount {
		end = int(totalCount)
	}

	return hb.Span().
		Class("small text-muted").
		Text(fmt.Sprintf("Showing %d-%d of %d %s", start, end, totalCount, itemLabel))
}
