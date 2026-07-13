# Task: Dashboard Enhancements — Bounce Rate, Visit Duration, Period Comparison, Realtime

## Source
Plausible + Rybbit feature research (`docs/proposals/todo/plausible.md`, `docs/proposals/todo/rybbit.md`)

## Status: Not Started

## Objective
Add four high-value dashboard metrics that require no schema changes — pure computation and presentation on existing visitor data.

## Features

### 1. Bounce Rate
- A visitor who viewed only 1 page = bounce.
- Compute: `(visitors with 1 page view) / (total visitors) * 100` for selected period.
- Group by fingerprint to count unique visitors and their page view count.
- Add to dashboard stats summary card.

### 2. Visit Duration / Time on Page
- Approximate from consecutive page view timestamps per fingerprint.
- For each visitor (by fingerprint), sort page views by `created_at`, compute time between consecutive visits.
- Last page view in a session has unknown duration (show as "-" or use average).
- Add to dashboard stats summary card as average visit duration.

### 3. Period-over-Period Comparison
- Show current period vs. previous period side-by-side.
- Dashboard already supports period selection (today, yesterday, 7d, week, month).
- Add a comparison column showing current vs. previous period count and % change.
- Pure presentation logic, no schema changes.

### 4. Realtime / Live Visitors
- Show count of visitors active in the last 5–15 minutes.
- Reuse existing `VisitorList` with narrow `created_at_gte` filter.
- Add a "Live" indicator card to the dashboard.
- Optional: auto-refresh via HTMX or meta refresh every 30s.

## Implementation Notes
- Bounce rate and visit duration require aggregation queries grouped by fingerprint.
- Period comparison needs a second query for the previous period (same duration shifted back).
- Realtime card can use existing `VisitorCount` with `SetCreatedAtGte(now - 15min)`.
- All four can be added to the home dashboard `home_controller.go` and stats summary cards.

## Dependencies
- No schema changes required.
- May benefit from a `VisitorCountByFingerprint` or `VisitorPageViewCount` store method for bounce rate (or compute in-memory from `VisitorList` results).

## Acceptance Criteria
- Bounce rate shown as percentage on dashboard for selected period.
- Average visit duration shown on dashboard for selected period.
- Period comparison shows current vs. previous period with % change indicator.
- Live visitor count card shows on dashboard, updates without full page reload.
