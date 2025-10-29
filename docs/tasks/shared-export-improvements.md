# Task: Shared Export Improvements

## Objective
Provide a reusable CSV (and future JSON/PDF) export helper that can be applied across admin tables to keep formatting consistent.

## Deliverables
- Move existing inline CSV export script into a shared JS snippet referenced by all controllers.
- Define server-side helper for naming files (e.g., `visitor-activity-YYYYMMDD.csv`).
- Document how to hook table IDs into the shared exporter.
- Ensure exported files include column headers and UTF-8 encoding with BOM for spreadsheet compatibility.

## Dependencies
- Requires layout support for injecting shared scripts (already available via `Layout.SetScripts`).
- Coordinates with each controller task to adopt the helper.

## Acceptance Criteria
- Visitor Activity and Visitor Paths both use the shared helper without duplicating code.
- Exported files open cleanly in Google Sheets and Excel.
- Documentation updated in `docs/admin-overview.md` to reference the new helper.
