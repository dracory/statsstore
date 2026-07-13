# Task: Visitor Paths Experience

## Objective
Upgrade the Visitor Paths page to mirror the StatCounter reference by combining per-session details with quick path navigation aides, while using existing visitor data where possible.

## Status: Mostly Implemented (~85%)

The visitor paths page is the most complete admin analytics page. It has filter toolbar, list-group rows with session metadata, pagination, per-page selector, server-side CSV export, and tests. A few items remain and some dead code exists.

## Current State

### Implemented
- **Two-column row layout** (`card_visitor_paths.go:179-202`) — `pathRow()` with left header (country badge + location + path link) and right header (session badge + device/browser badges + drill-down button), plus 3-column body (timestamp+IP, referrer, user agent)
- **Responsive stacking** — `flex-column flex-lg-row` classes for mobile/desktop
- **Country flag badge** (`countryBadge()`) — CSS flag icons via `fi fi-xx` class with country name tooltip
- **Path link with external icon** (`pathLink()`) — absolute URL via `fullPathURL()` using `ui.WebsiteUrl`, with `bi-box-arrow-up-right` icon
- **Timestamp block** (`timestampBlock()`) — entry time formatted; exit time hardcoded to "-"
- **Referrer block** (`referrerBlock()`) — green link or "(No referring link)" muted fallback
- **Session metadata column** (`sessionMetadataColumn()`) — session count badge, device badge, browser badge, drill-down button
- **Session count** (`sessionCount()`) — groups by fingerprint (or ID fallback) within current page results
- **Device/browser badges** (`deviceBadge()`, `browserBadge()`) — Bootstrap contextual color classes
- **Drill-down button** (`drillDownButton()`) — links to visitor activity page with path filter
- **Filter toolbar** — `addFilterDropdown()` with Last 24 Hours, Today, Country: Unknown, Device: Desktop, Path contains '/pricing'; `activeFilterBadges()` with all filter types
- **Path-level filters** — `path_contains` and `path_exact` query params, backed by `VisitorQuery().SetPathContains()` / `SetPathExact()` at store level
- **Quick range buttons** — All, Last 24 Hours, Today, Last 7 Days
- **Per-page selector** — 10/25/50/100 buttons with active state
- **Pagination** — `shared.PaginationUI` via `pagination()`
- **Pagination summary** — "Showing X-Y of Z paths"
- **Server-side CSV export** — `exportCSV()` in controller using `encoding/csv`, triggered via `?action=export`, with proper Content-Type and Content-Disposition headers
- **Export dropdown** — links to `?action=export` URL
- **Scripts** — HTMX and SweetAlert2 loaders wired in controller
- **Tests** — `visitor_paths_controller_test.go` with export CSV test (headers, data, session count, absolute URL) and store error test
- **Options button** — placeholder gear icon
- **Empty state** — friendly "No visitor paths recorded yet" message

### Not Implemented
- **Exit time** — hardcoded to "-" in `timestampBlock()`. No exit time concept in `VisitorInterface`; would require store-level changes. **Descope.**
- **Host name** — `VisitorInterface` has no host name / reverse DNS field. Path link uses `fullPathURL()` with `WebsiteUrl` + path. **Descope.**
- **Custom date picker** — no date picker modal or custom range UI (quick range buttons cover common cases)
- **Reusable row header partial** — `pathHeaderLeft()` is module-local, not shared. `countryBadge()`, `formatLocation()`, `resolvedCountryName()` are duplicated between visitor-activity and visitor-paths.
- **Documentation** — `docs/admin-overview.md` not updated

### Dead Code
- **`tableVisitorPaths()`** (`tsble_visitor_paths.go:13-51`) — hidden export table defined but never called. Export is server-side via `exportCSV()` in controller, not client-side JS. This file is leftover from a previous client-side export approach.
- **`countryFlagEmoji()`** (`card_visitor_paths.go:474-485`) — defined but `countryBadge()` uses CSS flag icons instead

### Inconsistencies vs Other Pages
- **No `Layout` dependency** — `visitor_paths_controller.go:47-50` creates controller without `Layout`, while visitor-activity requires it. Controller works because `Handler()` checks `c.ui.Layout` is non-nil before calling `SetTitle`/`SetBody`/`Render`. However, tests pass `nil` layout and only test export mode, not normal page rendering.
- **No page render test** — tests only cover CSV export and store error. No test for normal HTML rendering (unlike visitor-activity which has `TestVisitorActivityControllerHandlerSuccess`).
- **Typo in filename** — `tsble_visitor_paths.go` should be `table_visitor_paths.go`

## Key UI Elements (Reference Alignment)
- **Row header** showing country flag + location, host name, and page path with external link icon.
- **Timestamps** for entry/exit moments stacked vertically under each row header.
- **Referrer block** with primary link in green for "No referring link" fallback.
- **Session metadata** column indicating session number, magnifier icon for drill-down, and device/browser badges aligned to the right.
- **Filter / control bar** including:
  - "Add Filter" dropdown (date range, country, path contains, device type).
  - Export button using shared helper.
- **Footer controls**: pagination status, quick range buttons (All, 24 Hours, Today, custom range), date selection controls, and per-page selector.

## Remaining Deliverables
- Remove dead code: `tableVisitorPaths()` in `tsble_visitor_paths.go` (delete file), `countryFlagEmoji()`
- Fix filename typo: `tsble_visitor_paths.go` → `table_visitor_paths.go` (or delete if removing dead code)
- Extract shared helpers: `countryBadge()`, `formatLocation()`, `resolvedCountryName()`, `rangeLabel()`, `queryParamsWith()`, `quickRangeButtons()` pattern — move to `admin/shared` to avoid duplication with visitor-activity and page-view-activity
- Add page render test (normal HTML rendering, not just export)
- Migrate CSV export to shared helper once `shared-export-improvements.md` is implemented
- Add custom date picker if needed (descope if quick ranges are sufficient)
- Exit time: **descope** — no field in `VisitorInterface`
- Host name: **descope** — no field in `VisitorInterface`
- Update `docs/admin-overview.md` with current state

## Data / Store Considerations
- `statsstore.StoreInterface` provides path, timestamps, country, referrer, user agent, and fingerprint — **confirmed available** via `visitor_interface.go`.
- Session grouping derived from fingerprint via `sessionCount()` — counts visits with same fingerprint within current page results only (not global). Limitation noted.
- Path filtering via `SetPathContains()` and `SetPathExact()` — **confirmed available** at store level.
- No host name or exit time fields — **descoped**.

## Dependencies
- Shared export improvements task — **not started**, see `shared-export-improvements.md`. Current server-side export works independently but should migrate to shared helper for consistency.
- No store enhancements needed — all required fields are available.

## Acceptance Criteria
- Page layout matches reference structure on desktop and remains usable on mobile via stacked/collapsible layout.
- Filters update URL/query params and persist through pagination/export.
- Export output matches on-screen data columns.
- Admin overview documentation updated with screenshots/wireframes after implementation.
