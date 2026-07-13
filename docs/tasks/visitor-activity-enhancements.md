# Task: Visitor Activity Enhancements

## Objective
Evolve the Visitor Activity page to resemble the StatCounter reference while staying within the current data model and ensuring the experience integrates cleanly with the host admin shell.

## Status: Mostly Implemented (~75%)

The visitor activity page is functional with filter toolbar, list-group rows, pagination, and CSV export. Several enhancements from the original plan remain unimplemented, and some dead code exists.

## Current State

### Implemented
- **Two-column row layout** (`card-visitor-activity.go:232-270`) — `visitorRow()` with left header (country badge + location + IP) and right header (session badge + system summary), plus 3-column body (visit time/duration, referrer, path)
- **Responsive stacking** — uses `flex-column flex-lg-row` classes for mobile/desktop adaptation
- **Country flag badge** (`countryBadge()`) — CSS flag icons via `fi fi-xx` class
- **System summary** (`systemSummary()`) — browser + OS text with `deviceIcon()` and `osIcon()` Bootstrap icons
- **Session badge** (`sessionBadge()`) — fingerprint-based session label
- **Filter toolbar** — `addFilterDropdown()` with Last 24 Hours, Today, Country: Unknown, Device: Desktop; `activeFilterBadges()` with removal display
- **Quick range buttons** — All, Last 24 Hours, Today (missing Last 7 Days)
- **Pagination** — `paginationControls()` via `shared.PaginationUI`, `paginationSummary()` showing "Showing X-Y of Z visitors"
- **CSV export** — client-side JS `exportTableToCSV` with hidden `exportDataTable` table (8 columns: Visit Time, Path, Country, IP, Referrer, Browser, OS, User Agent)
- **Visitor detail modal** (`modal_visitor_detail.go`) — Bootstrap modal with loading spinner placeholder
- **Scripts** — HTMX and SweetAlert2 loaders wired in controller
- **Tests** — `visitor_activity_controller_test.go` with success and error cases
- **Options button** — placeholder gear icon button

### Not Implemented
- **Page views count** — `VisitorInterface` has no page views count field; not shown in rows
- **Exit time** — no exit time concept in data model; `formatVisitDuration` approximates duration between consecutive visits, not true exit time
- **ISP/IP lookup link** — IP address shown as plain text, no lookup link
- **Per-page selector** — missing (visitor-paths has 10/25/50/100 selector; visitor-activity does not)
- **Custom date picker** — no date picker modal or custom range UI
- **Last 7 Days quick range** — missing (visitor-paths has it)
- **Documentation** — `docs/admin-overview.md` not updated

### Dead Code
- **`locationBlock()`** (`card-visitor-activity.go:341-355`) — defined but not called by `visitorRow()`
- **`referrerBlock()`** (`card-visitor-activity.go:357-377`) — defined but not called by `visitorRow()`
- **`pathBlock()`** (`card-visitor-activity.go:379-384`) — defined but not called by `visitorRow()`
- **`countryFlagEmoji()`** (`card-visitor-activity.go:460-471`) — defined but `countryBadge()` uses CSS flag icons instead
- **`pagination()` in `pagination.go`** — defined but `card-visitor-activity.go` uses `paginationControls()` instead

### Inconsistencies vs Visitor Paths
- **Device filter not applied** — `helpers.go:35-47` parses `filters.Device` but never calls `options.SetDeviceType()`. Visitor-paths and page-view-activity both apply it.
- **No per-page selector** — visitor-paths has `perPageSelector()` with 10/25/50/100 buttons; visitor-activity lacks this
- **No "Last 7 Days" button** — visitor-paths `quickRangeButtons` includes 7d; visitor-activity stops at Today
- **Client-side export** — visitor-activity uses duplicated JS `exportTableToCSV`; visitor-paths uses server-side `encoding/csv` with `?action=export`
- **No `?action=export` handler** — visitor-activity controller has no server-side export route

## Key UI Elements (Reference Alignment)
- **Visit summary stack** per row showing:
  - Page views count and session indicator badges.
  - Exit time and resolution (hide gracefully if unknown).
  - Browser + OS icons and version metadata.
- **Session detail panel** on the right side of each row presenting:
  - Country flag + location text.
    - ISP/IP address with optional lookup link.
  - Referrer link block: primary referrer URL (green when "No referring link").
- **Interactive controls** mirroring reference:
  - "Add Filter" dropdown for quick filters (date range, country, device type).
  - Export and Options buttons (Export wired to shared helper; Options placeholder for future features).
  - Footer controls: page indicator, quick date range buttons (All, 24 Hours, Today, custom range), date picker, results-per-page selector.

## Remaining Deliverables
- Apply device filter in `buildControllerData` (add `options.SetDeviceType(filters.Device)` — bug fix, not enhancement)
- Add per-page selector (10/25/50/100) matching visitor-paths pattern
- Add "Last 7 Days" quick range button
- Migrate CSV export from client-side JS to server-side shared helper (depends on shared-export-improvements task)
- Remove dead code: `locationBlock()`, `referrerBlock()`, `pathBlock()`, `countryFlagEmoji()`, `pagination()` in `pagination.go`
- Wire visitor detail modal to actual data (currently just a loading spinner placeholder)
- Add custom date picker if needed (descope if not required)
- Page views count: descope — no field in `VisitorInterface`, would require store-level changes
- Exit time: descope — no field in `VisitorInterface`, `formatVisitDuration` is a reasonable approximation
- ISP/IP lookup link: descope or implement as external link to a lookup service
- Update `docs/admin-overview.md` with current state

## Dependencies
- `statsstore.StoreInterface` must expose page views count if available; otherwise document placeholder logic and follow-up task. **Descoped** — no page views field exists in `VisitorInterface`.
- Shared export improvements task for unified CSV handling. **Not started** — see `shared-export-improvements.md`.
- Potential new helper for ISP/IP lookup links. **Descoped** unless external lookup service is specified.

## Acceptance Criteria
- Desktop layout closely matches the reference hierarchy; mobile layout remains readable via collapse/stack.
- Filters, pagination, and exports remain functional and stateful.
- No regressions to existing visitor activity data fetching.
- Documentation updated (`admin-overview.md`) with screenshots/wireframes once implemented.
