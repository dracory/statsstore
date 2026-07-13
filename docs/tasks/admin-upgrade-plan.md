# statsstore Admin Panel – Revised Upgrade Roadmap

This roadmap reflects the current implementation inside `github.com/dracory/statsstore/admin` and recognizes that the admin experience is embedded inside a host application, which already handles authentication, authorization, and global security controls.

## Current State Snapshot

- **Available screens**: Home dashboard, Visitor Activity, Visitor Paths, Page View Activity (placeholder UI — controller and data layer ready, UI not implemented).
- **Data source**: `statsstore.StoreInterface` providing CRUD, counting, and soft-deletion of visitor records. `VisitorQueryInterface` supports filtering by country, date range, device type, path contains/exact, with limit/offset/order.
- **UI stack**: Server-rendered templates via `github.com/dracory/hb`, Bootstrap styling, and lightweight script injection for Chart.js, HTMX, and SweetAlert2.
- **Navigation**: Shared header/breadcrumb helpers (`shared.AdminHeaderUI`, `shared.Breadcrumbs`) with URL generation that respects the hosting endpoint. All 4 screens registered in nav.
- **Home dashboard**: Period selector (today, yesterday, last 7 days, previous 7 days, this week, last week, this month, last month), stats summary cards, Chart.js line chart, traffic sources breakdown (referrers, pages, browsers, countries, events, weekly heatmap). Uses `github.com/dracory/cdn` for script URLs.
- **Visitor Activity**: List-group rows with country flag, location, IP, session badge, system summary (browser/OS icons), referrer, path, duration. Filter toolbar, quick ranges, pagination, client-side CSV export, visitor detail modal scaffold (not wired). Tests exist.
- **Visitor Paths**: List-group rows with country flag, location, path link, timestamps, referrer, user agent, session count, device/browser badges, drill-down button. Filter toolbar with path filters, quick ranges, per-page selector, pagination, server-side CSV export. Tests exist.
- **Known issues**: 3 duplicated copies of client-side `exportTableToCSV` JS, device filter parsed but not applied in visitor-activity, dead code in visitor-activity and visitor-paths, no shared filter toolbar helper.

## Guiding Principles

1. **Leverage existing store capabilities first** – extend schema or tracking only when required by a new UI surface.
2. **Keep embedding simple** – avoid duplicating host responsibilities (auth, user management, plan upgrades).
3. **Ship iteratively** – prioritize enhancements that unlock actionable insights with minimal upstream dependencies.
4. **Favor progressive enhancement** – dashboards should render with core data without relying on heavy client-side tooling.

## Roadmap

### Phase 1 – Strengthen Existing Views (Near Term)

**Visitor Activity** — ~75% done, see `visitor-activity-enhancements.md`
- ~~Add country/region/IP columns using existing visitor fields.~~ **Done** — `visitorRow()` shows country badge, location, IP.
- ~~Introduce row expansion or modal detail using current modal scaffold.~~ **Partial** — modal scaffold exists (`modal_visitor_detail.go`) with loading spinner, but not wired to data.
- ~~Provide quick date filter presets (Today, 7 days, 30 days).~~ **Partial** — has All, Last 24 Hours, Today. Missing Last 7 Days and Last 30 Days buttons (though `parseFilters` supports them). Device filter parsed but not applied (bug).
- **Remaining**: Fix device filter bug, add 7d/30d quick ranges, add per-page selector, wire modal to data, remove dead code, migrate to server-side export.

**Visitor Paths** — ~85% done, see `visitor-paths-aggregation.md`
- ~~Replace simple visitor list with aggregated path summaries (group by `path`, count, last seen).~~ **Not done** — still shows individual visitor records, not aggregated by path. Session count is per-fingerprint within page results.
- ~~Expose filter for minimum visits within selected range.~~ **Not done** — no minimum visits filter.
- **Remaining**: Path aggregation (if still desired — current per-visit view may be sufficient), remove dead code, extract shared helpers, add page render test.

**Shared Enhancements** — not started, see `shared-export-improvements.md`
- ~~Standardize CSV export helpers across tables.~~ **Not done** — 3 copies of client-side JS + 1 server-side implementation, no shared helper.
- ~~Add loading indicators and empty-state messaging.~~ **Partial** — visitor-paths has empty state. Visitor-activity does not. Modal has loading spinner.
- **Remaining**: Create shared server-side CSV helper with UTF-8 BOM, migrate all controllers, extract shared filter toolbar/pagination/quick-range helpers to `admin/shared`.

### Phase 2 – Broader Analytics (Mid Term)

**Page View Activity** — ~20% done, see `page-view-activity-page.md`
- ~~Extend store or introduce a lightweight view to aggregate hits per path.~~ **Partial** — controller scaffold, data fetching, and filter parsing exist. UI is placeholder. Store already supports `SetPathContains`/`SetPathExact` for path-level queries.
- **Remaining**: Build filter toolbar, table/list composition, footer controls, CSV export, scripts, tests.

**Traffic Sources** — **Done** (exceeds original plan)
- ~~Capture and surface `user_referrer` and basic categorization (direct vs external).~~ **Done** — `traffic_sources_data.go` computes referrer breakdown with `normalizeReferrer()` categorization.
- ~~Provide simple charts highlighting top referrers over chosen period.~~ **Done** — `trafficSourcesCards()` renders tabular breakdowns for referrers, pages, browsers, countries, events, plus weekly heatmap. Not Chart.js charts but data tables with session counts.
- **Note**: Home dashboard already has rich traffic source analytics beyond the original plan scope.

**Filtering Framework**
- ~~Centralize date, country, and referrer filters in shared helper.~~ **Not done** — each controller has its own `parseFilters`, `queryParamsWith`, `rangeLabel`, `quickRangeButtons`. Near-identical implementations duplicated across 3 controllers.
- **Remaining**: Extract to `admin/shared` before implementing page-view-activity UI (avoid 4th copy).

### Phase 3 – Advanced Insights (Longer Term)
- **Session Reconstruction**
  - Explore storing session identifiers to enable chronological playback of visits.
  - Visualize time-on-page and navigation flow within a visit.
- **Reporting & Scheduling**
  - Generate downloadable summary reports (CSV/PDF) for key metrics, triggered manually inside the admin.
- **Map Visualization**
  - Integrate lightweight mapping (static choropleth or chart) using available country codes; avoid relying on precise geo IP unless already provided upstream.

## Defer / Out of Scope (Handled by Host or Future Consideration)

- Authentication, role management, and audit logging (managed by embedding application).
- Heatmaps, session replay, and full campaign attribution – require significant additional data capture and storage strategy changes.
- Notification systems, plan upgrades, and subscription management – assumed to be owned by host platform.
- Real-time streaming dashboards – current storage pattern is batch-friendly; adopt only if underlying ingestion pipeline supports it.

## Technical Considerations

- **Schema evolution**: Introduce migrations carefully; the store already supports automigrate for the visitor table, so new fields should preserve backward compatibility.
- **Performance**: Leverage existing `VisitorQueryOptions` (limit/offset, order) and add indexes when introducing new filters.
- **Extensibility**: Keep shared helpers generic so additional controllers can plug into navigation and breadcrumbs without duplication.
- **Integration**: Ensure URLs, scripts, and assets remain self-contained so the admin package continues to embed cleanly under various host domains.
- **Script loading**: Home controller uses `github.com/dracory/cdn` package for CDN URLs; visitor-activity and visitor-paths hardcode CDN URLs. Should standardize on `cdn` package.

## Next Steps

1. ~~Validate data availability for Phase 1 columns (country, referrer, IP)~~ — **Done**, all fields available via `VisitorInterface`.
2. ~~Draft detailed UI/UX updates for Visitor Activity and Visitor Paths~~ — **Done**, task docs exist for all 4 features.
3. **Extract shared helpers** (filter toolbar, CSV export, quick ranges, per-page selector) to `admin/shared` — highest priority, blocks clean page-view-activity implementation.
4. **Fix visitor-activity device filter bug** — `helpers.go` parses `filters.Device` but never calls `options.SetDeviceType()`.
5. **Remove dead code** across visitor-activity and visitor-paths modules.
6. **Implement page-view-activity UI** — consume shared helpers, avoid third copy of filter toolbar.
7. **Wire visitor detail modal** to actual data or remove if not needed.
8. Estimate effort for path aggregation queries (Phase 1 Visitor Paths) if per-visit view is insufficient.