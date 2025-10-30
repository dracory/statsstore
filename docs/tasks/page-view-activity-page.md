# Task: Page View Activity Page

## Objective
Introduce a dedicated "Page View Activity" screen that mirrors the StatCounter-style reference while leveraging current visitor data captured by `statsstore`.

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
- **Upgrade banner slot** (optional) allowing host application to inject plan messaging.

## Data & Store Requirements
- Confirm `VisitorList` exposes IP, path, referrer, user agent, country, and device metadata.
- Add helper to derive host name from IP/address (using existing or new utility).
- Implement translation layer for browser/OS to icons.

## UI Implementation Notes
- Build page as new controller under `admin/page-view-activity` with shared layout wiring.
- Reuse `shared.AdminHeaderUI` for navigation and add breadcrumb entries.
- Create reusable components for:
  - Filter bar (shared with other pages over time).
  - Host/referrer cell with multi-line formatting.
- Ensure table supports hover tooltips for long URLs and host information.

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
- **Shared utilities (`admin/shared`)**
  - Breadcrumb + header wiring (`shared.AdminHeaderUI`).
  - URL builders (`shared.UrlVisitorActivity`, `shared.UrlVisitorPaths`) as patterns for future `UrlPageViewActivity` helper.
  - Pagination component (`shared.PaginationUI`).
- **Statsstore accessors (`visitor.go`, `visitor_interface.go`)**
  - Provide IP, path, referrer, browser/OS/device metadata required for table columns.

## Implementation Plan
1. **Controller & Routing Scaffold**
   - Create `admin/page-view-activity` controller with `ControllerData` model mirroring the Visitor Activity pattern.
   - Register route in the admin router and wire breadcrumbs + `shared.AdminHeaderUI` title setup.
   - Ensure controller accepts shared dependencies (`ControllerOptions`) for layout, store, and URLs.

2. **Query & Store Enhancements**
   - Extend `statsstore.VisitorQueryOptions` with browser/device filters if missing; reuse existing date and country filters.
   - Add helper to derive host name / reverse lookup string (prefer existing utility in `admin/shared`; otherwise add new helper module).
   - Confirm `VisitorList` exposes required fields; surface missing ones via accessor additions if needed.

3. **Filter Toolbar**
   - Build reusable filter toolbar component (consider placing in `admin/shared`) with:
     - Primary “Add Filter” dropdown for range, country, device type, browser.
     - Quick range buttons (All, 24 Hours, Today, Custom -> links to date picker modal).
     - Active filter badges with removal links that update query parameters.
   - Ensure all filters sync with URL parameters for pagination/export consistency.

4. **Table Composition**
   - Implement HB components for table header + body, reusing visitor activity helpers for system/device badges.
   - Break row rendering into subcomponents:
     - Timestamp cells (date and time split).
     - System cell (browser + OS icons using shared helper).
     - Location / language cell with flag + locale text.
     - Host / page / referrer block with multi-line layout and tooltips for long URLs.
   - Add hover tooltips using Bootstrap `data-bs-toggle="tooltip"` where appropriate.

5. **Footer Controls & Export**
   - Reuse pagination component (or enhance existing `pagination` helper) with current page indicator.
   - Add per-page selector (10/25/50/100) mirroring Visitor Paths implementation.
   - Embed hidden export table (`visitor-page-view-export`) and hook `exportTableToCSV` helper from visitor paths.

6. **Scripts & Assets**
   - Ensure HTMX/SweetAlert loaders are only added once; reuse global helper if introduced for Visitor Paths.
   - Add tooltip initialization script (Bootstrap) scoped to this page.

7. **Testing & QA**
   - Write controller unit tests covering filter combinations, empty states, and pagination metadata.
   - Add snapshot/HTML assertions for table structure (focus on column presence and badge rendering).
   - Verify CSV export columns align with on-screen table.

## Follow-ups & Questions
- Should the filter toolbar be extracted into a shared component for reuse across activity screens? (Recommended.)
- Confirm availability of browser/OS icon mapping helper; if absent, specify asset requirements before implementation.
- Determine whether reverse DNS lookup is synchronous or deferred; may require async job or cached lookup helper.
