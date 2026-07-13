# statsstore Admin Panel Overview

## Purpose
The statsstore admin panel delivers ready-made dashboards and tables for monitoring visitor activity. It wraps the core `statsstore` APIs in server-rendered interfaces that can be embedded into broader admin experiences.

## Architecture
- **Entry Point (`admin` package)**: `New` validates dependencies (response writer, request, store, layout) and returns an `http.Handler`. Each request is routed through `ServeHTTP`, which reads `path` from query parameters, injects context (endpoint, admin home URL), and dispatches to feature controllers via a map.
- **Shared Layer (`admin/shared`)**: Provides constants, a `LayoutInterface`, `ControllerOptions`, URL helpers, breadcrumb rendering, navigation header UI, stat cards, pagination builder, and shared helper functions (badges, formatting, query params, info lines, export dropdown, quick range buttons, per-page selector, pagination summary). Controllers rely on this package for consistent layout, navigation, and UI components.
- **Controllers**:
  - `home`: Dashboard view summarizing traffic over the last 31 days. It aggregates daily totals/unique visits using `VisitorCount`, renders stat cards, charts (via Chart.js), navigation cards, and tables. Scripts are injected lazily for Chart.js, HTMX, and SweetAlert2.
  - `visitor-activity`: Paginated table of individual visits. Uses `VisitorList` with limit/offset, formats timestamps/durations, and offers CSV export hooks plus modal scaffolding for detailed visitor data.
  - `visitor-paths`: Highlights most visited paths. Reuses visitor records ordered by timestamp, displays per-session path details with country badges, device/browser badges, drill-down links to visitor activity, quick range filters, per-page selector, and server-side CSV export.
  - `page-view-activity`: Paginated table of individual page views. Uses `VisitorList` with limit/offset, supports range/country/device/browser filters, renders timestamp, path, location, referrer, and user-agent columns, and offers server-side CSV export.

## Request Flow
1. Host application constructs the admin handler with `admin.New`, passing a statsstore implementation and a layout.
2. Incoming requests (e.g., `/admin/home?path=/admin/visitor-activity`) hit `admin.ServeHTTP`.
3. The handler resolves the controller based on `shared.Path*` constants and forwards the request with augmented context.
4. Controllers fetch data via the injected `StoreInterface`, populate the shared layout, and return the rendered HTML string to the response writer.

## Layout Expectations
- The provided layout must implement setters for title, scripts, styles, body, and a `Render` method. Controllers call these before returning.
- Navigation URLs use helpers that respect the original endpoint from context, preserving reverse-proxy or embedding scenarios.

## Extensibility
- Adding a new section involves creating a controller package adhering to `shared.ControllerOptions`, registering it in `admin.findHandlerFromPath`, and leveraging shared helpers for navigation/breadcrumbs.
- Shared components (cards, pagination, modals) can be reused or extended for consistent UI.
- Scripts are currently injected as inline strings; replace or extend via layout hooks if bundling assets differently.

## Shared Helpers

The `admin/shared/helpers.go` file provides reusable functions extracted from controller packages to eliminate duplication across visitor-paths, visitor-activity, and page-view-activity:

### Query & Pagination
- **`shared.QueryParamsWith(r *http.Request, overrides map[string]string) map[string]string`** — Merges current query params with overrides (empty values delete keys).
- **`shared.ParseIntWithDefault(value string, defaultValue int) int`** — Parses a string to int with a fallback default.
- **`shared.ClampPerPage(perPage int) int`** — Constrains per-page values to the 1–100 range.

### Formatting
- **`shared.RangeLabel(value string) string`** — Converts range codes (e.g., `24h`, `today`) to human-readable labels.
- **`shared.ShortDate(value string) string`** — Parses a timestamp and returns `YYYY-MM-DD`.
- **`shared.FormatTimestamp(value string) string`** — Parses a timestamp and returns a user-friendly formatted string.
- **`shared.TimeParse(value string) (time.Time, error)`** — Attempts to parse a timestamp using several common layouts.

### Country / Location
- **`shared.ResolvedCountryName(ui ControllerOptions, code string) string`** — Resolves an ISO code to a country name, falling back to the code or "Unknown".
- **`shared.FormatLocation(ui ControllerOptions, visitor statsstore.VisitorInterface) string`** — Returns a human-readable location string for a visitor.
- **`shared.CountryBadge(ui ControllerOptions, visitor statsstore.VisitorInterface) hb.TagInterface`** — Renders a badge with a CSS flag icon and country name tooltip.

### URL Helpers
- **`shared.FullPathURL(ui ControllerOptions, path string) string`** — Builds an absolute URL from a website base URL and a path.
- **`shared.PathLink(ui ControllerOptions, path string) hb.TagInterface`** — Renders an anchor tag with an external icon pointing to the full absolute URL.

### Info Lines
- **`shared.InfoLine(label string, value hb.TagInterface) hb.TagInterface`** — Renders a labelled value line with a muted label and bold value.
- **`shared.InfoText(text string) hb.TagInterface`** — Renders a plain body-text span.
- **`shared.InfoMuted(text string) hb.TagInterface`** — Renders a muted, italic span for placeholder values.

### Badges
- **`shared.DeviceBadge(visitor statsstore.VisitorInterface) hb.TagInterface`** — Renders a contextual badge for the visitor's device type.
- **`shared.BrowserBadge(visitor statsstore.VisitorInterface) hb.TagInterface`** — Renders a badge showing the visitor's browser and version.

### UI Components
- **`shared.OptionsButton() hb.TagInterface`** — Renders a gear-icon button placeholder.
- **`shared.ExportDropdown(exportURL string) hb.TagInterface`** — Renders a dropdown with a link to the export URL.
- **`shared.QuickRangeButtons(r *http.Request, urlFunc func(map[string]string) string) hb.TagInterface`** — Renders a button group of common time-range filters.
- **`shared.PerPageSelector(r *http.Request, currentPageSize int, urlFunc func(map[string]string) string) hb.TagInterface`** — Renders a button group for selecting rows per page.
- **`shared.PaginationSummary(page, pageSize int, totalCount int64, itemLabel string) hb.TagInterface`** — Renders a "Showing X-Y of Z" summary text.

## CSV Export

All admin controllers that offer CSV export use a shared server-side helper located in `admin/shared/export.go`. This replaces the former duplicated client-side JavaScript `exportTableToCSV` functions with a single, testable Go implementation.

### Usage

Controllers detect an `?action=export` query parameter and delegate to the shared helper:

```go
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := c.prepareData(r)

	if action := r.URL.Query().Get("action"); action == "export" {
		if errorMessage != "" {
			w.WriteHeader(http.StatusInternalServerError)
			return errorMessage
		}
		return c.exportCSV(w, data)
	}

	// ... normal page rendering
}
```

The `exportCSV` method builds headers and rows, then calls:

```go
shared.ExportCSV(w, shared.ExportFilename("visitor-paths"), headers, rows)
```

### Functions

- **`shared.ExportFilename(prefix string) string`** — Generates a CSV filename with the current UTC date (e.g., `visitor-paths-2024-01-15.csv`).
- **`shared.ExportCSV(w http.ResponseWriter, filename string, headers []string, rows [][]string) string`** — Writes CSV content with a UTF-8 BOM for Excel/Google Sheets compatibility, sets `Content-Type` and `Content-Disposition` headers, and returns the CSV string.

### Controllers Using Shared Export

- `home` — Exports daily stats summary (date, page views, unique visits, first-time visits, returning visits).
- `visitor-activity` — Exports individual visit records (visit time, path, country, IP, referrer, browser, OS, user agent).
- `visitor-paths` — Exports visitor path records (visit time, path, absolute URL, country, IP, referrer, session, device, browser).
- `page-view-activity` — Exports page view records (date, time, path, absolute URL, country, IP, referrer, device, browser, OS, user agent).

### Testing

Shared export tests are in `admin/shared/export_test.go`, covering filename generation, BOM presence, header/row correctness, and Content-Disposition headers.

## Dependencies
- HTML generation uses `github.com/dracory/hb`.
- Charts and interactions rely on external CDNs (Chart.js, HTMX, SweetAlert2) loaded on demand.
- Pagination utilities depend on `github.com/spf13/cast` for string conversion.
