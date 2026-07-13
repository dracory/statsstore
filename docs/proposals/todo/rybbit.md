# Rybbit Feature Research

Source: https://rybbit.com/

Rybbit is an open-source, cookieless, GDPR/CCPA-compliant web & product analytics platform.
This document compares Rybbit's features against the current statsstore project and identifies
features we could add.

---

## Current statsstore Features

- Visitor tracking: IP, fingerprint, user agent, browser, OS, device type, country, path, referrer, accept-language
- Admin dashboard with 4 pages: Home (stats summary + charts), Visitor Activity, Visitor Paths, Page View Activity
- Period filtering: today, yesterday, last 7 days, this/last week, this/last month
- Query filtering: by country, date range, device type, path (exact/contains), distinct, ID, limit/offset, order by
- CSV export (visitor activity, visitor paths)
- Cookieless (fingerprint-based)
- Soft delete support
- SQL-based storage (Go library, embeddable)

---

## Rybbit Feature Comparison

### Core Web Analytics

| Feature | Rybbit | statsstore | Gap? |
|---|---|---|---|
| Page views | Yes | Yes (page view activity page) | No |
| Visitors | Yes | Yes (visitor activity page) | No |
| Bounce rate | Yes | No | **Yes** |
| Traffic sources | Yes | Partial (referrer field only) | **Yes** |
| Location | Yes (city-level) | Partial (country only) | **Yes** |
| Devices | Yes | Yes (device, browser, OS) | No |
| Languages | Yes | Yes (accept-language field) | No |
| Filtering | Yes | Yes (query filters) | No |
| Realtime data | Yes | No | **Yes** |
| Custom events | Yes | No | **Yes** |
| Custom data | Yes | No | **Yes** |
| UTM tracking | Yes | No | **Yes** |
| Links (click tracking) | Yes | No | **Yes** |
| Bot blocking | Yes | No | **Yes** |

### Advanced Analytics

| Feature | Rybbit | statsstore | Gap? |
|---|---|---|---|
| Session replay | Yes | No | **Yes** |
| Web vitals | Yes | No | **Yes** |
| Funnels | Yes | No | **Yes** |
| Goals | Yes | No | **Yes** |
| Journey | Yes | Partial (visitor paths page) | **Yes** |
| Globe views | Yes | No | **Yes** |
| Error tracking | Yes | No | **Yes** |
| User sessions | Yes | Partial (visitor activity) | **Yes** |
| Google Search Console | Yes | No | **Yes** |
| Compare (period-over-period) | Yes | No | **Yes** |
| User profiles | Yes | No | **Yes** |
| Retention | Yes | Partial (first/return visit tracking) | **Yes** |

### Access & Sharing

| Feature | Rybbit | statsstore | Gap? |
|---|---|---|---|
| Organizations | Yes | No | **Yes** |
| Public dashboards | Yes | No | **Yes** |
| Private link sharing | Yes | No | **Yes** |
| RBAC | Yes | No | **Yes** |

### Privacy

| Feature | Rybbit | statsstore | Gap? |
|---|---|---|---|
| GDPR & CCPA | Yes | Yes (cookieless) | No |
| Data anonymization | Yes | Yes (fingerprint, no PII storage policy) | No |
| No cookies | Yes | Yes | No |
| Data ownership | Yes | Yes (self-hosted library) | No |

### Cloud / Platform

| Feature | Rybbit | statsstore | Gap? |
|---|---|---|---|
| Email reports | Yes | No | **Yes** |
| API access | Yes | No (library only, no HTTP API layer) | **Yes** |

---

## Recommended Features to Add (Prioritized)

### Tier 1 — High Value, Moderate Effort

1. **Bounce Rate Calculation**
   - Compute from existing visitor data: a visitor who viewed only 1 page = bounce.
   - No schema changes needed; pure query/aggregation logic.
   - Add to dashboard stats summary card.

2. **Realtime Data / Live Visitors**
   - Show visitors active in the last 5-15 minutes.
   - Reuse existing `VisitorList` with a narrow `created_at_gte` filter.
   - Add a "Live" indicator card to the dashboard.

3. **Bot Blocking**
   - Filter requests by user-agent patterns (known bot/crawler signatures).
   - Add a `IsBot()` check in `VisitorRegister` or a configurable blocklist.
   - Prevents data pollution with minimal effort.

4. **UTM Tracking**
   - Parse `utm_source`, `utm_medium`, `utm_campaign`, `utm_term`, `utm_content` from request URL.
   - Store as fields on the visitor record (schema change).
   - Enables campaign performance reporting.

5. **Compare (Period-over-Period)**
   - Show current period vs. previous period side-by-side.
   - The dashboard already supports period selection; add a comparison column.
   - No schema changes, pure presentation logic.

### Tier 2 — Medium Value, Medium Effort

6. **Custom Events**
   - Allow tracking of named events (sign-up, purchase, download) beyond page views.
   - Requires an `event_type` field on the visitor record or a separate events table.
   - Add event tracking endpoint and event-based filtering.

7. **Funnels**
   - Define a sequence of pages/events and measure drop-off at each step.
   - Built on top of visitor paths data (already available).
   - New admin page with funnel visualization.

8. **Goals**
   - Define conversion goals (e.g., "visited /thank-you page").
   - Track goal completion rate against total visitors.
   - New admin page + goal configuration.

9. **User Sessions (Enhanced)**
   - Group page views by visitor fingerprint within a time window (e.g., 30 min).
   - Show session timeline per visitor.
   - Enhances the existing visitor activity page.

10. **Retention Tracking (Enhanced)**
    - Cohort analysis: track returning visitors over time.
    - The dashboard already tracks first/return visits; expand to cohort tables.
    - New admin page with retention heatmap.

### Tier 3 — High Value, High Effort

11. **API Layer**
    - Expose visitor data via REST/JSON endpoints.
    - Enables third-party integrations and custom dashboards.
    - New HTTP handler package alongside the admin UI.

12. **Email Reports**
    - Scheduled summary reports (daily/weekly/monthly).
    - Requires a scheduler/cron integration and email sending.
    - Template-based report generation from existing stats.

13. **Web Vitals**
    - Capture Core Web Vitals (LCP, FID, CLS, INP) from the client.
    - Requires a JS snippet to collect metrics and an API endpoint to receive them.
    - New fields on visitor record or separate metrics table.

14. **Error Tracking**
    - Capture client-side JavaScript errors.
    - Requires JS error handler + reporting endpoint.
    - New error log table and admin page.

15. **Session Replay**
    - Record DOM mutations and user interactions for replay.
    - Heavy: requires a JS recording library (e.g., rrweb), storage for replay data, and a replay viewer.
    - Significant infrastructure investment.

### Tier 4 — Nice to Have

16. **Globe Views** — 3D traffic visualization (cosmetic, high effort).
17. **Organizations** — Multi-tenant site/team management (architectural change).
18. **Public Dashboards** — Shareable read-only dashboard URLs (access control layer).
19. **Private Link Sharing** — Password-protected dashboard links (access control).
20. **RBAC** — Role-based access control (auth system required).
21. **Google Search Console Integration** — Import GSC data (external API dependency).
22. **Link Click Tracking** — Track outbound link clicks (JS snippet + endpoint).
23. **City-Level Location** — GeoIP database integration for city resolution (currently country-only).

---

## Summary

statsstore already covers the privacy fundamentals (cookieless, GDPR-friendly, self-hosted) and
basic visitor analytics (page views, visitors, devices, traffic sources, paths). The biggest gaps
vs. Rybbit are in **advanced analytics** (funnels, goals, session replay, web vitals) and
**platform features** (API, email reports, multi-tenant access).

The most impactful additions with least effort are:
1. Bounce rate (pure computation from existing data)
2. Realtime/live visitors (narrow time filter on existing query)
3. Bot blocking (user-agent filtering)
4. UTM tracking (parse URL params, store as fields)
5. Period comparison (presentation layer only)
