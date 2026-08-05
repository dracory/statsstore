# Features Plan — Prototype v3

## Synthesis Approach

After analyzing 12 research sources covering privacy-first analytics tools (Plausible, Rybbit, Swetrix, Peek, Umami, Tinyvisits, Kobbe, Visitors.now), comprehensive platforms (Matomo, GA4), and person-level tracking (KissMetrics), the following plan synthesizes the findings into a concrete feature set.

### Key Findings from Research

1. **Minimal dashboards win**: Plausible, Peek, Umami, Tinyvisits, Kobbe all emphasize single-page or few-page dashboards with essential metrics only. No overwhelming multi-level menus.
2. **Essential metrics consensus**: All tools agree on: pageviews, unique visitors, sessions, bounce rate, avg session duration. These are non-negotiable.
3. **Real-time is expected**: 10 of 12 tools feature live visitor counts. It's table stakes.
4. **Traffic sources are universal**: Every tool shows referrers, sources, and geographic/device breakdowns.
5. **Privacy-first is the differentiator**: Cookieless, no fingerprinting, GDPR compliance is the core identity — not a feature.
6. **Funnels & journeys are advanced**: Available in Plausible, Rybbit, Peek, Kobbe, GA4 — but marked as "deeper" views, not main dashboard.
7. **Page-level detail matters**: Top pages, entry pages, exit pages appear across all tools.
8. **Period comparison is valued**: Plausible, Kobbe, GA4 all support comparing time periods.

---

## Proposed Page Structure (5 pages)

Based on the research, the principle of "keep pages minimal so users do not get overwhelmed" leads to a **5-page structure**:

### Page 1: Dashboard (Overview)
The main landing page — a calm, at-a-glance view of essential metrics.

### Page 2: Pages
Detailed page-by-page analytics — top pages, entry/exit pages.

### Page 3: Sources
Traffic acquisition — referrers, channels, campaigns, geographic data.

### Page 4: Visitors
Visitor/audience breakdown — devices, browsers, OS, countries. Plus real-time visitor activity.

### Page 5: Live
Real-time dashboard — live visitor count, current activity stream, active pages.

---

## Features by Page

### Page 1: Dashboard (Overview)

| Feature | In Scope | Research Support |
|---|---|---|
| KPI cards: Unique Visitors, Pageviews, Sessions, Bounce Rate, Avg Session Duration | ✅ In | Plausible, Swetrix, Umami, Kobbe, Matomo, GA4 — universal across all tools |
| Main trend chart (line chart, metric selectable) | ✅ In | Plausible, Swetrix, Kobbe, GA4 — all have a main chart with metric toggle |
| Period selector (Today, Yesterday, 7d, 30d, Month, Custom) | ✅ In | Plausible, Swetrix, Kobbe, Umami, GA4 — universal |
| Period comparison (vs previous period) | ✅ In | Plausible, Kobbe, GA4 — valued for trend analysis |
| Mini top pages list (top 5) | ✅ In | Plausible, Swetrix, Umami, Kobbe — quick glance at popular content |
| Mini top sources list (top 5) | ✅ In | Plausible, Swetrix, Umami, Kobbe — quick glance at traffic origin |
| Live visitor count badge in topbar | ✅ In | Swetrix, Plausible, Peek, Kobbe, Visitors.now — table stakes |
| Loading skeletons | ✅ In | Modern UX best practice, seen across tools |
| Chart annotations (pin notes to dates) | ❌ Out | Kobbe — useful but adds complexity; defer to future |
| Saved segments/filters | ❌ Out | Plausible, Swetrix, GA4 — too complex for minimal prototype |
| Consolidated multi-site view | ❌ Out | Plausible, Rybbit — out of scope for single-site prototype |

### Page 2: Pages

| Feature | In Scope | Research Support |
|---|---|---|
| Full top pages table (page path, pageviews, unique visitors, bounce rate, avg time on page) | ✅ In | Plausible, Swetrix, Umami, Kobbe, GA4 — all show detailed page tables |
| Entry pages table (landing pages, visits, bounce rate, session duration) | ✅ In | Plausible, Swetrix, GA4 — entry pages are universally tracked |
| Exit pages table (exit page, exits, exit rate) | ✅ In | Plausible, Swetrix — exit rate is a key engagement metric |
| Page trend chart (pageviews over time for selected page) | ✅ In | Swetrix, GA4 — drill-down into individual page performance |
| Search/filter pages | ✅ In | Swetrix, GA4 — necessary for sites with many pages |
| Scroll depth per page | ❌ Out | Plausible — requires custom event tracking; out of scope for mock |
| Page-level conversion rate | ❌ Out | Plausible, GA4 — requires goal setup; out of scope |

### Page 3: Sources

| Feature | In Scope | Research Support |
|---|---|---|
| Referrer table (source, visitors, visits, bounce rate, session duration) | ✅ In | Plausible, Swetrix, Umami, Kobbe, GA4 — universal |
| Source/Medium breakdown | ✅ In | Swetrix, GA4, Matomo — standard acquisition dimension |
| UTM campaign table (source, medium, campaign, content, term) | ✅ In | Plausible, Swetrix, Peek, Kobbe, GA4 — campaign tracking is essential |
| Countries table with flag icons | ✅ In | Plausible, Swetrix, Umami, Kobbe, Visitors.now — geographic data universal |
| Regions/cities breakdown | ❌ Out | Rybbit, Visitors.now, Umami — city-level adds complexity; country-level sufficient for prototype |
| Channel grouping (organic, social, direct, referral, paid) | ✅ In | Plausible, GA4, Matomo, Swetrix — automatic channel categorization |
| Network intelligence (ISP, connection type) | ❌ Out | Swetrix — too niche for minimal prototype |
| Conversion attribution by source | ❌ Out | Plausible, GA4, Matomo — requires goal tracking; out of scope |

### Page 4: Visitors

| Feature | In Scope | Research Support |
|---|---|---|
| Device type breakdown (desktop, mobile, tablet) with chart | ✅ In | Swetrix, Umami, Kobbe, Matomo, GA4 — universal |
| Browser breakdown (Chrome, Safari, Firefox, etc.) with chart | ✅ In | Plausible, Swetrix, Umami, Kobbe — universal |
| OS breakdown (Windows, macOS, Linux, iOS, Android) with chart | ✅ In | Plausible, Swetrix, Umami, Kobbe — universal |
| Countries breakdown with flag icons | ✅ In | All tools — geographic data is fundamental |
| New vs. returning visitors | ✅ In | Matomo, GA4, Swetrix — key audience metric |
| Visitor trend chart (unique visitors over time) | ✅ In | Plausible, Swetrix, GA4 — trend visualization |
| Visitor profiles (individual person timelines) | ❌ Out | KissMetrics, Rybbit, Visitors.now — requires person-level tracking; contradicts privacy-first anonymous model |
| Session replays | ❌ Out | Rybbit, Swetrix, Umami — complex feature, out of scope |
| Retention cohorts | ❌ Out | GA4, Rybbit — requires persistent tracking; out of scope for privacy-first model |
| Heatmaps | ❌ Out | Swetrix, Kobbe — requires session replay infrastructure |
| Web Vitals | ❌ Out | Rybbit, Peek, Kobbe, Visitors.now — performance monitoring is separate concern |

### Page 5: Live

| Feature | In Scope | Research Support |
|---|---|---|
| Live visitor count (large number, pulsing indicator) | ✅ In | Swetrix, Plausible, Peek, Kobbe, Visitors.now — 10/12 tools have this |
| Active pages list (pages currently being viewed) | ✅ In | Kobbe, Visitors.now, Swetrix — see what people are reading right now |
| Active referrers/sources | ✅ In | Kobbe, Visitors.now — where live traffic is coming from |
| Active countries | ✅ In | Kobbe, Visitors.now — geographic distribution of current visitors |
| Recent activity stream (events as they happen) | ✅ In | Visitors.now, Kobbe — streaming activity feed |
| Auto-refresh (simulated with mock data) | ✅ In | Plausible (30s), Swetrix (real-time), Visitors.now (WebSocket) — live updates |
| 3D globe visualization | ❌ Out | Visitors.now — impressive but heavy; simple list is more minimal |
| Live visitor profiles (click to see journey) | ❌ Out | Visitors.now — requires individual tracking; contradicts privacy model |
| Browser tab title counter | ❌ Out | Swetrix — nice touch but not a dashboard feature |

---

## Design Principles for v3

1. **Privacy-first**: No cookies, no persistent visitor profiles, anonymous aggregated data — the core identity (Plausible, Peek, Umami, Tinyvisits, Kobbe)
2. **Minimal & calm**: Few pages, few metrics per page, no overwhelming dashboards (Peek, Umami, Tinyvisits, Plausible)
3. **Real-time**: Live visitor count visible on every page (topbar badge), dedicated Live page (Swetrix, Plausible, Peek, Kobbe)
4. **Essential metrics only**: Pageviews, unique visitors, sessions, bounce rate, avg session duration, top pages, top sources, countries, devices/browsers/OS (consensus across all tools)
5. **Responsive**: Works on mobile and desktop (modern baseline)
6. **Dark/light theme**: Toggle persisted in localStorage (modern UX expectation)
7. **Self-contained**: Static HTML + CSS + JS, no build step, CDN-loaded Bootstrap + Vue + Chart.js
8. **Multi-page architecture**: Each page is a separate standalone HTML file. Navigation via `<a href>` links in sidebar. No client-side routing, no Vue router, no shared state beyond `shared.js`/`shared.css`

---

## Shared Assets

### `shared.css`
- CSS variables for theming (dark/light)
- Sidebar navigation styles
- Topbar styles (with live visitor badge)
- Card component styles (KPI cards)
- Table styles
- Toast notification styles
- Loading skeleton styles
- Responsive breakpoints (mobile sidebar collapse, etc.)
- Chart container styles

### `shared.js`
- Mock data (visitors, pageviews, sources, pages, countries, devices, browsers, OS, live activity)
- Vue app factory function (creates Vue app with shared data/methods)
- Theme toggle logic (localStorage persistence)
- Toast notification system
- Navigation config (sidebar items)
- Chart.js helper functions (create line chart, bar chart, doughnut chart)
- Date range helper functions
- Number formatting helpers (e.g., 1.2K, 3.4M)
- Flag icon helper (country code to flag CSS class)

---

## Out-of-Scope Summary

The following features are explicitly **out of scope** for the v3 prototype, based on the "keep pages minimal" principle and the privacy-first anonymous model:

- Individual visitor profiles or person-level tracking
- Session replays
- Heatmaps
- Web Vitals / Core Web Vitals
- Funnels and journeys (could be a future page)
- Goals and conversions (requires event tracking setup)
- Revenue attribution / ecommerce
- Retention cohorts
- Saved segments/filters
- Multi-site consolidated view
- 3D globe visualization
- Network intelligence (ISP, connection type)
- Scroll depth tracking
- Chart annotations
- Email/Slack reports
- API access

---

## Technology Stack (per task spec)

- Vue.js 3 (CDN) for reactivity
- Bootstrap 5 (CDN) for layout/components
- Bootstrap Icons (CDN) for icons
- flag-icons (CDN) for country flags
- Chart.js (CDN) for charts
- No build step, no npm, no bundler, no TypeScript
- All data is mock data in `shared.js`
- AGPL-3.0 license reference in footer
