# Task: Page View Activity Page

## Objective
Introduce a dedicated "Page View Activity" screen that mirrors the StatCounter-style reference while leveraging current visitor data captured by `statsstore`.

## Status: Partially Implemented (~20%)

### Done
- Controller scaffold (`admin/page-view-activity/controller.go`) with routing registered in `admin.go`
- Breadcrumbs + `shared.AdminHeaderUI` nav entry wired
- `UrlPageViewActivity` URL helper in `shared/utils.go`
- `ControllerData` and `FilterOptions` types in `types.go` (includes `Browser` field)
- `buildControllerData` in `helpers.go` — queries store with country/date/device filters, computes pagination metadata
- `parseFilters` — handles range presets (24h, today, 7d, 30d), custom from/to, country, device, browser
- `clampPerPage` and `parseIntWithDefault` helpers

### Not Done
- **Filter toolbar UI** — no dropdown, badges, or quick range buttons rendered; `page()` shows placeholder alert div
- **Table/list composition** — no row rendering, no timestamp/system/location/host-referrer cells
- **Footer controls** — no pagination summary, per-page selector, or pagination controls
- **CSV export** — no export dropdown, no server-side CSV handler
- **Scripts** — `SetScripts` never called (no HTMX/SweetAlert2/tooltip init)
- **Tests** — no `*_test.go` in `admin/page-view-activity/`

## Key Features
- **Table columns**
  - Date and time (separate columns) derived from visitor `created_at`.
  - System column with browser and OS icons (reuse existing helper logic from visitor activity).
  - Location / language with country flag icon when available.
  - Host name / page / referrer block showing:
    - Entry path with external link to live site.
    - Reverse DNS / IP lookup link when host name available.
    - Green "(No referring link)" fallback messaging.
- **Row badges** for device category (desktop / mobile / tablet) and session indicators if present.
- **Filters toolbar**
  - Primary "Add Filter" button exposing dropdown for date range, country, device type, browser.
  - Quick range buttons (All, Today, 24 Hours, custom date picker) similar to reference component.
- **Pagination footer** with page indicator and results-per-page selector.

## Data & Store Requirements
- `VisitorList` exposes IP, path, referrer, user agent, country, and device metadata — **confirmed available** via `visitor_interface.go`.
- Browser filter: `FilterOptions.Browser` is parsed but not applied (`helpers.go:47`). `VisitorQueryInterface` has no `SetBrowser` method. Either add store-level support or remove from UI.
- Host name / reverse DNS: `VisitorInterface` has no host name field. **Descoped** — display IP address directly via `GetIpAddress()`.

## UI Implementation Notes
- Build page as new controller under `admin/page-view-activity` with shared layout wiring.
- Reuse `shared.AdminHeaderUI` for navigation and add breadcrumb entries.
- Create reusable components for:
  - Filter bar (shared with other pages over time).
  - Host/referrer cell with multi-line formatting.
- Ensure table supports hover tooltips for long URLs and host information.

## Design Decisions to Resolve
- **Table vs list-group layout**: Plan describes a traditional `<table>` with separate Date/Time columns. Both existing pages (visitor-activity, visitor-paths) use a card/list-group layout with `infoLine` rows. Confirm whether a true `<table>` is wanted or the list-group pattern should be reused for consistency.
- **CSV export approach**: Visitor-paths uses server-side CSV (`encoding/csv` in controller). Visitor-activity uses client-side JS (`exportTableToCSV`). Page-view-activity should use the server-side approach (cleaner, matches visitor-paths).
- **Shared filter toolbar**: Both visitor-activity and visitor-paths have near-identical `filterToolbar`, `addFilterDropdown`, `activeFilterBadges`, `quickRangeButtons`, `queryParamsWith`, `rangeLabel`. Consider extracting to `admin/shared` before adding a third copy.
- **Browser filter**: Implement at store level (`VisitorQueryInterface.SetBrowser`) or remove from filter UI to avoid dead code.

## Acceptance Criteria
- Page renders with responsive table matching reference hierarchy.
- Filters update results via query parameters and persist across pagination.
- CSV export uses shared helper with same columns.
- No new authentication or user-role logic inside this module (handled by host system).

## Existing Components to Reuse
- **Visitor Activity card (`admin/visitor-activity/card-visitor-activity.go`)**
  - Filter toolbar scaffolding (`addFilterDropdown`, `activeFilterBadges`) with URL-sync behaviour.
  - Pagination helpers (`quickRangeButtons`, `paginationControls`) wired through `shared.PaginationUI`.
  - Export dropdown + hidden table pattern (`exportDataTable`) already wired to `exportTableToCSV` helper.
  - Device/session badges and system summary utilities for browser/OS text + icons.
- **Visitor Paths card (`admin/visitor-paths/card_visitor_paths.go`)**
  - Modernized list layout with three-column detail body that can inform row composition.
  - Footer controls combining summary text, quick ranges, per-page selector, and pagination call.
  - Server-side CSV export (`exportCSV` method using `encoding/csv`).
- **Shared utilities (`admin/shared`)**
  - Breadcrumb + header wiring (`shared.AdminHeaderUI`).
  - URL builder `shared.UrlPageViewActivity` — already exists.
  - Pagination component (`shared.PaginationUI`).
- **Statsstore accessors (`visitor.go`, `visitor_interface.go`)**
  - Provide IP, path, referrer, browser/OS/device metadata required for table columns.

## Remaining Implementation Plan
1. **Filter Toolbar**
   - Build filter toolbar component (consider placing in `admin/shared`) with:
     - Primary "Add Filter" dropdown for range, country, device type.
     - Quick range buttons (All, 24 Hours, Today, Last 7 Days) using `UrlPageViewActivity`.
     - Active filter badges with removal links that update query parameters.
   - Ensure all filters sync with URL parameters for pagination/export consistency.
   - Port `queryParamsWith` and `rangeLabel` helpers (or extract to shared).

2. **Table/List Composition**
   - Implement HB components for row rendering, reusing visitor activity helpers for system/device badges.
   - Break row rendering into subcomponents:
     - Timestamp cells (date and time split).
     - System cell (browser + OS icons using shared helper).
     - Location / language cell with flag + locale text.
     - Host / page / referrer block with multi-line layout and tooltips for long URLs.
   - Add hover tooltips using Bootstrap `data-bs-toggle="tooltip"` where appropriate.
   - Decide: `<table>` element or list-group pattern (see Design Decisions).

3. **Footer Controls & Export**
   - Add pagination summary ("Showing X-Y of Z visitors").
   - Add per-page selector (10/25/50/100) mirroring Visitor Paths implementation.
   - Add pagination controls using `shared.PaginationUI` with `UrlPageViewActivity` URL func.
   - Implement server-side CSV export (`exportCSV` method with `encoding/csv`, matching visitor-paths pattern).
   - Add export dropdown in card header.

4. **Scripts & Assets**
   - Wire `SetScripts` with HTMX/SweetAlert2 loaders (copy pattern from visitor-paths controller).
   - Add tooltip initialization script (Bootstrap) scoped to this page.

5. **Testing & QA**
   - Write controller unit tests covering filter combinations, empty states, and pagination metadata.
   - Add snapshot/HTML assertions for table structure (focus on column presence and badge rendering).
   - Verify CSV export columns align with on-screen table.
   - Match pattern from `visitor_paths_controller_test.go`.
