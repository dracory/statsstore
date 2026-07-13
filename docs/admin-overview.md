# statsstore Admin Panel Overview

## Purpose
The statsstore admin panel delivers ready-made dashboards and tables for monitoring visitor activity. It wraps the core `statsstore` APIs in server-rendered interfaces that can be embedded into broader admin experiences.

## Architecture
- **Entry Point (`admin` package)**: `New` validates dependencies (response writer, request, store, layout) and returns an `http.Handler`. Each request is routed through `ServeHTTP`, which reads `path` from query parameters, injects context (endpoint, admin home URL), and dispatches to feature controllers via a map.
- **Shared Layer (`admin/shared`)**: Provides constants, a `LayoutInterface`, `ControllerOptions`, URL helpers, breadcrumb rendering, navigation header UI, stat cards, and pagination builder. Controllers rely on this package for consistent layout and navigation.
- **Controllers**:
  - `home`: Dashboard view summarizing traffic over the last 31 days. It aggregates daily totals/unique visits using `VisitorCount`, renders stat cards, charts (via Chart.js), navigation cards, and tables. Scripts are injected lazily for Chart.js, HTMX, and SweetAlert2.
  - `visitor-activity`: Paginated table of individual visits. Uses `VisitorList` with limit/offset, formats timestamps/durations, and offers CSV export hooks plus modal scaffolding for detailed visitor data.
  - `visitor-paths`: Intended to highlight most visited paths. Currently reuses visitor records ordered by timestamp while counting distinct paths; UI mirrors Visitor Activity with tailored labels.

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

### Testing

Shared export tests are in `admin/shared/export_test.go`, covering filename generation, BOM presence, header/row correctness, and Content-Disposition headers.

## Dependencies
- HTML generation uses `github.com/dracory/hb`.
- Charts and interactions rely on external CDNs (Chart.js, HTMX, SweetAlert2) loaded on demand.
- Pagination utilities depend on `github.com/spf13/cast` for string conversion.
