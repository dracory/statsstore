# statsstore Admin Panel – Revised Upgrade Roadmap

This roadmap reflects the current implementation inside `github.com/dracory/statsstore/admin` and recognizes that the admin experience is embedded inside a host application, which already handles authentication, authorization, and global security controls.

## Current State Snapshot

- **Available screens**: Home dashboard, Visitor Activity table, Visitor Paths list.
- **Data source**: `statsstore.StoreInterface` providing CRUD, counting, and soft-deletion of visitor records.
- **UI stack**: Server-rendered templates via `github.com/dracory/hb`, Bootstrap styling, and lightweight script injection for Chart.js, HTMX, and SweetAlert2.
- **Navigation**: Shared header/breadcrumb helpers with URL generation that respects the hosting endpoint.

## Guiding Principles

1. **Leverage existing store capabilities first** – extend schema or tracking only when required by a new UI surface.
2. **Keep embedding simple** – avoid duplicating host responsibilities (auth, user management, plan upgrades).
3. **Ship iteratively** – prioritize enhancements that unlock actionable insights with minimal upstream dependencies.
4. **Favor progressive enhancement** – dashboards should render with core data without relying on heavy client-side tooling.

## Roadmap

### Phase 1 – Strengthen Existing Views (Near Term)
- **Visitor Activity**
  - Add country/region/IP columns using existing visitor fields.
  - Introduce row expansion or modal detail using current modal scaffold.
  - Provide quick date filter presets (Today, 7 days, 30 days) using query options already supported by the store.
- **Visitor Paths**
  - Replace simple visitor list with aggregated path summaries (group by `path`, count, last seen).
  - Expose filter for minimum visits within selected range.
- **Shared Enhancements**
  - Standardize CSV export helpers across tables.
  - Add loading indicators and empty-state messaging.

### Phase 2 – Broader Analytics (Mid Term)
- **Page View Activity**
  - Extend store or introduce a lightweight view to aggregate hits per path and render sortable table with counts and last visit.
- **Traffic Sources**
  - Capture and surface `user_referrer` and basic categorization (direct vs external).
  - Provide simple charts highlighting top referrers over chosen period.
- **Filtering Framework**
  - Centralize date, country, and referrer filters in shared helper to keep behavior consistent across screens.

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

## Next Steps

1. Validate data availability for Phase 1 columns (country, referrer, IP) and add missing collectors if needed.
2. Draft detailed UI/UX updates for Visitor Activity and Visitor Paths, including wireframes or HB component breakdowns.
3. Estimate effort for store-level changes (aggregation queries) and add tasks to implementation backlog.