# Plausible Feature Research

Source: https://plausible.io/

Plausible is an open-source, privacy-first, cookieless web analytics platform (AGPL-3.0).
It is deliberately simpler than tools like Rybbit or PostHog — no session replay, no web vitals,
no error tracking, no user profiles — but excels in clean UX, platform/sharing features, and
strict privacy. This document compares Plausible's features against the current statsstore
project and identifies features we could add.

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

## Plausible Feature Comparison

### Core Web Analytics

| Feature | Plausible | statsstore | Gap? |
|---|---|---|---|
| Page views | Yes | Yes (page view activity page) | No |
| Visitors | Yes (daily counter, resets every 24h) | Yes (visitor activity page) | No |
| Bounce rate | Yes | No | **Yes** |
| Visit duration | Yes (time on page, engagement-based) | No | **Yes** |
| Traffic sources | Yes (referrers, search, social, campaigns, channels, AI assistants) | Partial (referrer field only) | **Yes** |
| Location | Yes (country-level) | Partial (country only) | No |
| Devices | Yes | Yes (device, browser, OS) | No |
| Languages | No | Yes (accept-language field) | No |
| Filtering | Yes (dashboard filters, saved segments) | Yes (query filters) | No |
| Realtime data | Yes (real-time dashboard, updates every 30s) | No | **Yes** |
| Custom events | Yes | No | **Yes** |
| Custom data | Yes (custom properties) | No | **Yes** |
| UTM tracking | Yes (utm_source, utm_medium, utm_campaign, utm_content, utm_term, ref, source) | No | **Yes** |
| Links (click tracking) | Yes (outbound link clicks, automated) | No | **Yes** |
| Bot blocking | Yes (built-in bot filtering, referrer spam, data center traffic) | No | **Yes** |
| Scroll depth tracking | Yes (automatic, 1–100%) | No | **Yes** |
| File downloads | Yes (automatic tracking) | No | **Yes** |
| Form submissions | Yes (automatic tracking) | No | **Yes** |
| 404 error pages | Yes (tracking) | No | **Yes** |

### Advanced Analytics

| Feature | Plausible | statsstore | Gap? |
|---|---|---|---|
| Session replay | No (deliberately not offered) | No | No |
| Web vitals | No | No | No |
| Funnels | Yes (strict order funnels) | No | **Yes** |
| Goals | Yes (codeless, custom events, revenue) | No | **Yes** |
| Journey | Yes (user journeys, forward & backward) | Partial (visitor paths page) | **Yes** |
| Globe views | No | No | No |
| Error tracking | No | No | No |
| User sessions | Partial (within-day only, no cross-day linking) | Partial (visitor activity) | No |
| Google Search Console | Yes (search keywords integration) | No | **Yes** |
| Compare (period-over-period) | Yes | No | **Yes** |
| User profiles | No (deliberately not offered) | No | No |
| Retention | No (no cohort analysis) | Partial (first/return visit tracking) | No |
| Revenue tracking | Yes (ecommerce revenue attribution) | No | **Yes** |

### Access & Sharing

| Feature | Plausible | statsstore | Gap? |
|---|---|---|---|
| Organizations | Yes (teams, multi-site management) | No | **Yes** |
| Public dashboards | Yes (shareable public links, embeddable) | No | **Yes** |
| Private link sharing | Yes (password-protected, segment-limited) | No | **Yes** |
| RBAC | Yes (team roles, 2FA enforcement, SSO) | No | **Yes** |
| Shared segments | Yes (saved & shared segments) | No | **Yes** |
| Embedded dashboards | Yes (embeddable without branding) | No | **Yes** |
| White label | Yes | No | **Yes** |

### Privacy

| Feature | Plausible | statsstore | Gap? |
|---|---|---|---|
| GDPR & CCPA | Yes (no cookies, no personal data) | Yes (cookieless) | No |
| Data anonymization | Yes (no IP, no fingerprint, no persistent ID) | Yes (fingerprint, no PII storage policy) | No |
| No cookies | Yes | Yes | No |
| Data ownership | Yes (self-hosted Community Edition) | Yes (self-hosted library) | No |

### Cloud / Platform

| Feature | Plausible | statsstore | Gap? |
|---|---|---|---|
| Email reports | Yes (email + Slack, weekly/monthly, traffic spike notifications) | No | **Yes** |
| API access | Yes (Stats API, Sites API, Events API) | No (library only, no HTTP API layer) | **Yes** |
| GA import | Yes (Google Analytics historical import) | No | **Yes** |
| Consolidated view | Yes (multi-site aggregate dashboard) | No | **Yes** |
| Looker Studio | Yes (connector) | No | **Yes** |
| SSO | Yes (Enterprise) | No | **Yes** |
| Managed proxy | Yes (Enterprise) | No | **Yes** |
| Script proxying | Yes (first-party script from own domain) | No | **Yes** |

---

## Recommended Features to Add (Prioritized)

### Tier 1 — High Value, Moderate Effort

1. **Bounce Rate Calculation**
   - Compute from existing visitor data: a visitor who viewed only 1 page = bounce.
   - No schema changes needed; pure query/aggregation logic.
   - Add to dashboard stats summary card.

2. **Realtime Data / Live Visitors**
   - Show visitors active in the last 5–15 minutes.
   - Reuse existing `VisitorList` with a narrow `created_at_gte` filter.
   - Add a "Live" indicator card to the dashboard.

3. **Bot Blocking**
   - Filter requests by user-agent patterns (known bot/crawler signatures).
   - Add an `IsBot()` check in `VisitorRegister` or a configurable blocklist.
   - Plausible also filters data center traffic and referrer spam — consider all three.
   - Prevents data pollution with minimal effort.

4. **UTM Tracking**
   - Parse `utm_source`, `utm_medium`, `utm_campaign`, `utm_term`, `utm_content` from request URL.
   - Also parse `ref` and `source` params (Plausible-supported alternatives).
   - Store as fields on the visitor record (schema change).
   - Enables campaign performance reporting.

5. **Compare (Period-over-Period)**
   - Show current period vs. previous period side-by-side.
   - The dashboard already supports period selection; add a comparison column.
   - No schema changes, pure presentation logic.

6. **Visit Duration / Time on Page**
   - Plausible uses engagement signals to calculate time on page.
   - Can approximate from consecutive page view timestamps per fingerprint.
   - Add to dashboard stats summary card.

### Tier 2 — Medium Value, Medium Effort

7. **Custom Events**
   - Allow tracking of named events (sign-up, purchase, download) beyond page views.
   - Requires an `event_type` field on the visitor record or a separate events table.
   - Add event tracking endpoint and event-based filtering.

8. **Goals**
   - Define conversion goals (e.g., "visited /thank-you page", "clicked signup button").
   - Track goal completion rate against total visitors.
   - Plausible supports codeless goals (page visit), custom event goals, and automated goals (404, downloads, outbound links, form submissions).
   - New admin page + goal configuration.

9. **Funnels**
   - Define a sequence of pages/events and measure drop-off at each step.
   - Plausible supports strict order funnels (consecutive steps with no other actions in between).
   - Built on top of visitor paths data (already available).
   - New admin page with funnel visualization.

10. **Outbound Link Click Tracking**
    - Track clicks on external links automatically.
    - Requires a JS snippet to capture click events and an endpoint to receive them.
    - Plausible offers this as a plug-and-play toggle.

11. **File Download Tracking**
    - Track file downloads (PDF, ZIP, etc.) automatically.
    - Requires JS snippet to detect download link clicks.
    - Plausible offers this as a plug-and-play toggle.

12. **Form Submission Tracking**
    - Track form completions automatically.
    - Requires JS snippet to hook into form submit events.
    - Plausible offers this as a plug-and-play toggle.

13. **Scroll Depth Tracking**
    - Track how far visitors scroll on each page (1–100%).
    - Requires JS snippet to measure scroll position.
    - Plausible tracks this automatically on every page with no setup.

14. **404 Error Page Tracking**
    - Track visits to 404 pages as a special event type.
    - Can be detected server-side (HTTP 404 status) or client-side.
    - Plausible offers this as an automated goal.

15. **User Sessions (Enhanced)**
    - Group page views by visitor fingerprint within a time window (e.g., 30 min).
    - Show session timeline per visitor.
    - Enhances the existing visitor activity page.
    - Note: Plausible only links sessions within the same day by design; statsstore can go further with fingerprint-based cross-day linking.

### Tier 3 — High Value, High Effort

16. **API Layer**
    - Expose visitor data via REST/JSON endpoints.
    - Plausible offers Stats API (query metrics), Sites API (manage sites), and Events API (record events without JS).
    - Enables third-party integrations and custom dashboards.
    - New HTTP handler package alongside the admin UI.

17. **Email Reports**
    - Scheduled summary reports (daily/weekly/monthly).
    - Plausible also offers traffic spike notifications via email and Slack.
    - Requires a scheduler/cron integration and email sending.
    - Template-based report generation from existing stats.

18. **Revenue Tracking**
    - Track ecommerce revenue alongside conversion goals.
    - Send revenue value with custom events (e.g., purchase amount).
    - Plausible offers revenue attribution and revenue breakdown in reports.
    - Requires schema change (revenue field) and event tracking infrastructure.

19. **Google Search Console Integration**
    - Import search keyword data from Google Search Console.
    - Display organic search queries in the dashboard.
    - Requires OAuth integration with Google Search Console API.

20. **Consolidated View (Multi-Site)**
    - Aggregate stats across multiple sites/dashboards into a single view.
    - Requires multi-site support (site ID field on visitor records).
    - New admin page showing combined metrics.

### Tier 4 — Nice to Have

21. **Public Dashboards** — Shareable read-only dashboard URLs (access control layer).
22. **Private Link Sharing** — Password-protected dashboard links with segment filtering.
23. **Embedded Dashboards** — Embeddable dashboard widgets without branding.
24. **Organizations** — Multi-tenant site/team management (architectural change).
25. **RBAC** — Role-based access control (auth system required).
26. **SSO** — Single sign-on integration (enterprise feature).
27. **Saved Segments** — Save and share filter presets across the team.
28. **GA Import** — Import historical Google Analytics data.
29. **Looker Studio Connector** — Data connector for Google Looker Studio.
30. **White Label** — Remove all branding for reseller/agency use.
31. **Script Proxying** — Serve tracking script as first-party from own domain.
32. **Managed Proxy** — Managed first-party proxy service (enterprise feature).
33. **AI Assistant Traffic Tracking** — Detect and categorize traffic from AI tools (ChatGPT, Claude, Perplexity, etc.) as a dedicated channel.

---

## Summary

Plausible and statsstore share the same privacy fundamentals (cookieless, GDPR-friendly,
self-hosted). Plausible is deliberately simpler than Rybbit — it does **not** offer session
replay, web vitals, error tracking, user profiles, or retention/cohort analysis. This means
statsstore has fewer gaps vs. Plausible than vs. Rybbit.

However, Plausible excels in two areas where statsstore has significant gaps:

1. **Automated tracking** — Plausible offers plug-and-play tracking for outbound links, file
   downloads, form submissions, scroll depth, and 404 pages. statsstore only tracks page views.
2. **Platform & sharing** — Plausible has a mature API (Stats/Sites/Events), email/Slack reports,
   team management, public/private dashboard sharing, SSO, embedded dashboards, and GA import.
   statsstore has none of these.

A key philosophical difference: Plausible stores **no IP addresses, no fingerprints, and no
persistent identifiers** — its daily unique visitor counter resets every 24 hours and cannot
link visits across days. statsstore uses fingerprint-based tracking with IP storage, which
enables cross-day visitor linking but is less privacy-strict. This is a design trade-off, not
a gap.

The most impactful additions with least effort are:
1. Bounce rate (pure computation from existing data)
2. Realtime/live visitors (narrow time filter on existing query)
3. Bot blocking (user-agent + data center + referrer spam filtering)
4. UTM tracking (parse URL params, store as fields)
5. Period comparison (presentation layer only)
6. Visit duration (compute from consecutive page view timestamps)
