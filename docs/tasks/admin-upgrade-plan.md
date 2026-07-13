# statsstore Admin Panel – Revised Upgrade Roadmap

## Remaining Work

### Phase 3 – Advanced Insights (Longer Term)
- **Session Reconstruction**
  - Explore storing session identifiers to enable chronological playback of visits.
  - Visualize time-on-page and navigation flow within a visit.
- **Reporting & Scheduling**
  - Generate downloadable summary reports (CSV/PDF) for key metrics, triggered manually inside the admin.
- **Map Visualization**
  - Integrate lightweight mapping (static choropleth or chart) using available country codes; avoid relying on precise geo IP unless already provided upstream.

### Misc
- **Browser filter** — `FilterOptions.Browser` parsed but not applied in page-view-activity (`VisitorQueryInterface` has no `SetBrowser`). Add store-level support or remove from UI.
- **Script loading** — Home controller uses `github.com/dracory/cdn` package for CDN URLs; visitor-activity and visitor-paths hardcode CDN URLs. Standardize on `cdn` package.
- **Documentation** — `docs/admin-overview.md` not updated to reflect current state.

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