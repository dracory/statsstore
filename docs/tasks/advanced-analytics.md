# Task: Advanced Analytics — Web Vitals, Error Tracking, Session Replay, Retention, Revenue

## Source
Rybbit + Plausible feature research

## Status: Not Started

## Objective
Add advanced analytics features that require significant infrastructure investment. These are longer-term items.

## Features

### Web Vitals (Rybbit)
- Capture Core Web Vitals: LCP (Largest Contentful Paint), FID (First Input Delay), CLS (Cumulative Layout Shift), INP (Interaction to Next Paint).
- JS snippet collects metrics via `PerformanceObserver` API.
- Store as fields on visitor record or separate `web_vitals` table.
- Display in admin with page-level performance breakdown.

### Error Tracking (Rybbit)
- Capture client-side JavaScript errors.
- JS error handler (`window.onerror` + `unhandledrejection`) sends errors to endpoint.
- Store in `error_logs` table with stack trace, URL, user agent, timestamp.
- Admin page: `/admin/errors` with error list and grouping.

### Session Replay (Rybbit)
- Record DOM mutations and user interactions for replay.
- Requires a JS recording library (e.g., rrweb).
- Storage for replay data (can be large — consider compression or sampling).
- Replay viewer in admin.
- **High effort** — significant infrastructure investment.

### Retention / Cohort Analysis (Rybbit)
- Cohort analysis: track returning visitors over time.
- Group visitors by first-visit date (cohort), track return visits in subsequent days/weeks.
- Display as cohort retention heatmap (rows = cohorts, columns = days/weeks, cells = % returning).
- Dashboard already tracks first/return visits — expand to full cohort table.
- New admin page: `/admin/retention`.

### Revenue Tracking (Plausible)
- Track ecommerce revenue alongside conversion goals.
- Send revenue value with custom events (e.g., `purchase` event with `amount: 99`).
- Revenue attribution: which referrer/campaign drove the purchase.
- Revenue breakdown in reports: by source, by goal, by time period.
- Depends on `custom-events.md` and `goals-funnels.md`.

## Implementation Notes
- Web vitals and error tracking require client-side JS snippets and new endpoints.
- Session replay is the heaviest feature — consider deferring unless specifically needed.
- Retention is pure computation from existing fingerprint + timestamp data (no schema change).
- Revenue tracking requires custom events infrastructure + optional `revenue` field on events.

## Dependencies
- `custom-events.md` for revenue tracking.
- `goals-funnels.md` for revenue attribution.
- Client-side JS snippet for web vitals and error tracking.
- Significant storage for session replay data.

## Acceptance Criteria
- Web vitals captured per page view and visible in admin.
- Client-side errors logged and visible in admin with grouping.
- Session replay available (if implemented) with playback viewer.
- Retention cohort heatmap shows return rates by cohort.
- Revenue tracked and attributed to sources/campaigns.
