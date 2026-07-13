# statsstore Admin Panel – Upgrade Roadmap

## Remaining Work

### Phase 1 — Quick Wins (No Schema Changes)
- **Dashboard Enhancements** — bounce rate, visit duration, period comparison, live visitors → `dashboard-enhancements.md`
- **Bot Filtering** — user-agent bot detection, referrer spam blocking, data center filtering → `bot-filtering.md`

### Phase 2 — Schema Changes & Event Infrastructure
- **UTM Tracking** — parse and store UTM campaign params → `utm-tracking.md`
- **Custom Events** — named events with custom properties, event tracking endpoint → `custom-events.md`

### Phase 3 — Advanced Analytics
- **Goals & Funnels** — conversion goals, funnel drop-off visualization → `goals-funnels.md`
- **Automated Tracking** — outbound links, file downloads, form submissions, scroll depth, 404s → `automated-tracking.md`
- **Enhanced Sessions** — session grouping by fingerprint, navigation flow, cross-day linking → `enhanced-sessions.md`

### Phase 4 — Platform & Integration
- **API Layer** — REST/JSON endpoints for stats, events, and sites → `api-layer.md`
- **Email Reports** — scheduled reports, traffic spike notifications, Slack integration → `email-reports.md`

### Phase 5 — High Effort / Longer Term
- **Advanced Analytics** — web vitals, error tracking, session replay, retention cohorts, revenue tracking → `advanced-analytics.md`
- **Platform Features** — public dashboards, multi-site, white label, SSO, saved segments, GA import, AI assistant tracking → `platform-features.md`

### Misc
- **Browser filter** — `FilterOptions.Browser` parsed but not applied in page-view-activity (`VisitorQueryInterface` has no `SetBrowser`). Add store-level support or remove from UI.
- **Script loading** — Home controller uses `github.com/dracory/cdn` package for CDN URLs; visitor-activity and visitor-paths hardcode CDN URLs. Standardize on `cdn` package.
- **Documentation** — `docs/admin-overview.md` not updated to reflect current state.

## Defer / Out of Scope (Handled by Host or Future Consideration)

- Authentication, role management, and audit logging (managed by embedding application).
- Real-time streaming dashboards – current storage pattern is batch-friendly; adopt only if underlying ingestion pipeline supports it.

## Technical Considerations

- **Schema evolution**: Introduce migrations carefully; the store already supports automigrate for the visitor table, so new fields should preserve backward compatibility.
- **Performance**: Leverage existing `VisitorQueryOptions` (limit/offset, order) and add indexes when introducing new filters.
- **Extensibility**: Keep shared helpers generic so additional controllers can plug into navigation and breadcrumbs without duplication.
- **Integration**: Ensure URLs, scripts, and assets remain self-contained so the admin package continues to embed cleanly under various host domains.