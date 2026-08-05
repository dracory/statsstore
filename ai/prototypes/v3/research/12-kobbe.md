# Kobbe — Research Summary

## Overview
Kobbe is a privacy-first, cookieless web analytics platform. 1.8KB default tracker (gzip). GDPR-ready, 0 raw IPs stored. Combines traffic and revenue in one place. No cookies, no fingerprinting, no persistent visitor profiles.

## Core Features

### Traffic Overview
- Visitors, visits, views, bounce rate, session time
- Trend changes in one dashboard
- Main chart with selected KPI over time
- Hover/click data points to inspect buckets
- Period comparison (previous period, week, month, quarter, year, custom)
- Comparison adds dashed line to main chart

### Realtime Visitors
- Active visitors count (approximate, online right now)
- Current pages being viewed
- Recent activity in live window
- Referrers / sources of active traffic
- Locations (country or region hints)
- Live map of online visitors
- Intentionally approximate count (useful for presence and direction, not billing-grade)

### Insights
- **Engagement KPIs**: Average daily visitors, single-visit share, custom events, core traffic metrics
- **Ranked breakdowns**: Top sources, countries, devices, browsers, pages, converting events in one view
- **Conversion peak heatmap**: Day-by-hour heatmap showing when custom events happen most often

### Pages
- Top pages with traffic data
- Click a page row to inspect single path

### Sources
- Referrer tracking
- Channel breakdown

### Locations
- Country/region geographic data

### Devices
- Device type breakdown

### Events
- Custom events: clicks, signups, purchases, downloads
- Events activity log with filtering and CSV export
- Auto-track: form submits, contact clicks, outbound links, messaging taps
- Filter dashboard by goal

### Funnels
- Build from pages and custom events
- Trace UTM campaigns, sources, mediums in same view
- Measure drop-off across page paths and events

### Conversions
- Auto-track common goals (contact clicks, form submits)
- Conversion tracking with UTM campaign connection

### Revenue Attribution
- Pass attribution through checkout
- Revenue beside traffic KPIs, pages, sources
- Tab-scoped attribution through payment metadata (not visitor profiling)
- Total revenue, order count, average order value
- Breakdowns by source, country, landing page

### 404 Tracking
- Flag not-found page once
- See broken URLs, hit counts, which internal page linked to them

### Chart Annotations
- Pin notes to specific days on traffic chart
- Launches, campaigns, incidents stay tied to numbers

### Web Vitals
- Core Web Vitals tracking
- Percentile trends, slow pages, environment breakdowns
- Track beside traffic data

### Search Console
- Organic queries, landing pages, clicks, impressions, rankings
- AI traffic referrals shown alongside other sources

## Privacy Model
- **No cookies**: Default tracker uses no cookies
- **No fingerprinting**: No canvas, WebGL, font enumeration techniques
- **No browser storage**: No `localStorage` or `sessionStorage` for measurement/profiling
- **No raw IPs stored**: Server computes short-lived anonymous hash from request metadata + daily-rotating secret
- **Hash not reversible**: Scoped to site, not shared across days
- **No full URLs**: Query strings stripped, only referrer origin sent
- **Respects GPC/DNT**: No analytics request when Global Privacy Control or Do Not Track enabled
- **Opt-out**: `localStorage.kobbe_ignore = "true"` flag

## Additional Features
- Monthly reports (opt-in email summaries)
- Traffic alerts (spike/drop notifications)
- Shared dashboard links (overview + realtime, overview only, realtime only)
- CSV export

## Design Principles
- **Every view answers a real question**: Traffic overview, realtime, insights, funnels, revenue
- **Privacy-first defaults**: No cookies, no fingerprinting, no raw IPs, no persistent profiles
- **Compact views**: Switch between visitors and revenue in one view
- **Context-aware**: Annotations, comparisons, breakdowns together
- **Lightweight**: 1.8KB tracker, async loading

## Sources
- https://kobbe.io/
- https://kobbe.io/docs
- https://kobbe.io/docs/privacy-and-cookieless-tracking
- https://kobbe.io/docs/dashboard-overview
- https://kobbe.io/docs/realtime-visitors
