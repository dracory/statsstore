# Task: Enhanced Session Reconstruction

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Group page views by visitor fingerprint within a time window to reconstruct sessions and visualize navigation flow.

## Features
- **Session grouping**: group page views by fingerprint within a 30-minute inactivity window.
- **Session timeline**: show chronological list of page views within a session.
- **Session metrics**: duration, page count, bounce status, entry/exit pages.
- **Cross-day linking**: statsstore uses fingerprint-based tracking, enabling cross-day session linking (Plausible cannot do this by design).
- **Navigation flow**: visualize path sequence within a session (forward and backward).
- **Session detail view**: admin page showing individual session timelines.

## Implementation Notes
- Query all page views for a fingerprint within selected period, sort by `created_at`.
- Split into sessions: if gap between consecutive page views > 30 min, start new session.
- Compute session metrics: start time, end time, duration, page count, entry page, exit page, bounce (1 page view).
- Enhance existing visitor activity page with session grouping toggle.
- New admin page: `/admin/sessions` with session list and detail view.
- Session detail: chronological page view list with timestamps and paths.

## Dependencies
- No schema changes — uses existing fingerprint and timestamp fields.
- May benefit from a `VisitorListByFingerprint` store method for efficient session queries.

## Acceptance Criteria
- Page views grouped into sessions by fingerprint with 30-min inactivity threshold.
- Session timeline visible in admin with chronological page views.
- Session metrics (duration, page count, bounce, entry/exit) shown.
- Cross-day sessions linked correctly via fingerprint.
- Session list filterable by date range and fingerprint.
