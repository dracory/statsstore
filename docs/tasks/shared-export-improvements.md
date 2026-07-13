# Task: Shared Export Improvements

## Objective
Provide a reusable CSV (and future JSON/PDF) export helper that can be applied across admin tables to keep formatting consistent.

## Status: Not Started

No shared export helper exists. Export is implemented independently (and inconsistently) in each controller.

## Current State

### Client-side JS export (3 copies of `exportTableToCSV`)
- **`admin/visitor-activity/visitor_activity_controller.go:71`** — inline JS function, called via `onclick` from `card-visitor-activity.go:415`
- **`admin/home/chart_stats_summary.go:133`** — identical copy of the same JS function
- **`admin/home/card_stats_summary.go:55`** — calls `exportTableToCSV('stats-table', 'visitor_stats.csv')` via onclick (relies on the function defined in `chart_stats_summary.go`)
- All three copies are byte-for-byte identical (read table rows, quote fields, join with comma, create Blob, trigger download)
- No UTF-8 BOM in any copy

### Server-side CSV export (1 implementation)
- **`admin/visitor-paths/visitor_paths_controller.go:96`** — `exportCSV` method using `encoding/csv`, triggered via `?action=export` query parameter
- Sets `Content-Type: text/csv; charset=utf-8` and `Content-Disposition: attachment; filename="visitor-paths-YYYY-MM-DD.csv"`
- No UTF-8 BOM
- Has test coverage: `TestVisitorPathsControllerExportCSV` and `TestVisitorPathsControllerExportCSVStoreError` in `visitor_paths_controller_test.go`

### Page View Activity
- No export implementation at all (placeholder page)

## Problems
- **3 duplicated copies** of the client-side `exportTableToCSV` JS function across visitor-activity and home
- **2 incompatible approaches**: client-side JS (visitor-activity, home) vs server-side Go (visitor-paths)
- **No UTF-8 BOM** anywhere — Excel may misinterpret encoding for non-ASCII content
- **No shared helper** — each controller re-implements export independently
- **No documentation** in `docs/admin-overview.md` about export

## Deliverables
- Choose one export approach (server-side recommended — cleaner, testable, already proven in visitor-paths)
- Create shared server-side CSV helper in `admin/shared` (e.g., `shared.ExportCSV(w, filename, headers, rows)`)
- Include UTF-8 BOM (`\xEF\xBB\xBF`) at start of buffer for Excel/Sheets compatibility
- Define shared filename helper (e.g., `shared.ExportFilename(prefix)` → `prefix-2006-01-02.csv`)
- Migrate visitor-paths `exportCSV` to use the shared helper
- Migrate visitor-activity from client-side JS to server-side shared helper
- Migrate home stats table from client-side JS to server-side shared helper
- Remove all 3 copies of inline `exportTableToCSV` JS from controllers
- Document usage in `docs/admin-overview.md`

## Dependencies
- Requires layout support for injecting shared scripts (already available via `Layout.SetScripts`) — only needed if client-side approach is kept for any controller
- Coordinates with page-view-activity task (should adopt shared helper from the start)

## Acceptance Criteria
- All controllers (visitor-activity, visitor-paths, home, page-view-activity) use the shared server-side CSV helper
- No duplicated `exportTableToCSV` JS function remains in any controller
- Exported files include UTF-8 BOM and open cleanly in Google Sheets and Excel
- Export tests cover the shared helper directly (not per-controller reimplementation)
- Documentation updated in `docs/admin-overview.md` to reference the shared helper
