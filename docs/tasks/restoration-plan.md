# Task: Restore Advanced Analytics Features

## Source
Product Specification + Feature Audit

## Status: In Progress

## Objective
Restore missing analytics features described in the product specification and task documentation, including fingerprint-based tracking, period comparisons, bounce rates, and enhanced session reconstruction.

## Summary
The current implementation of `statsstore` is missing several key features that were previously planned or existed in earlier versions. This plan aims to restore these features to align the codebase with the `specification.md` and `docs/tasks/` documentation.

## Implementation Steps

### 1. Core Tracking Enhancements
- [x] **Fingerprint Calculation**: Fixed `VisitorRegister` in `store.go` to correctly calculate and save the `fingerprint` field (MD5 of IP + User-Agent).
- **Fingerprint-based Analytics**:
    - [x] Updated `VisitorCount` (already supported `SetDistinct`) and dashboard logic to use fingerprints instead of IP addresses for "Unique Visitors".
    - **Separate Audience Page**: To avoid cluttering the main dashboard, implement a new "Audience" page for "First vs Return" counts, retention cohorts, and detailed visitor insights.

### 2. Dashboard Advanced Metrics
- **Period-over-Period Comparison**:
    - Update `handle_load_dashboard.go` to perform a second query for the previous period.
    - Calculate percentage changes for Total Visitors, Unique Visitors, and other metrics.
    - Update `dashboard.html` and `dashboard.js` to display these comparison indicators.
- **Realtime / Live Visitors**:
    - Add a card to the dashboard showing the count of unique fingerprints active in the last 15 minutes.
    - Implement a small AJAX endpoint or update the main dashboard load to include this.
- **Bounce Rate & Duration**:
    - Implement bounce rate calculation: `(sessions with 1 visit) / (total sessions)`.
    - Implement average visit duration calculation based on session start/end timestamps.

### 3. Enhanced Session Reconstruction
- **30-Minute Inactivity Window**: Update `handle_load_sessions.go` to group visits into sessions based on the `fingerprint` and a 30-minute inactivity gap, rather than just grouping all visits by IP.
- **Navigation Flow**:
    - Identify Entry and Exit pages for each session.
    - Update the session detail view to show a clear chronological timeline of paths.

### 4. Data Export
- [x] **CSV Export**:
    - Added `action=export-csv` handlers to `Dashboard`, `Visitors`, and `Sessions` controllers in both `statsstore` and `statsadmin`.
    - Added "Export CSV" buttons to the respective UI pages.
    - Used Go's `encoding/csv` to generate the files.
    - Implemented frontend download logic via `fetch` and `Blob`.

### 5. UI/UX Polishing
- **Daily Stats Table/Chart**: Restore the day-by-day breakdown table/chart mentioned in the spec.
- **Standardize CDN usage**: Use the `cdn` package for all external scripts/styles in the admin panel as noted in `admin-upgrade-plan.md`.
- [x] **IP Details Navigation**: Made visitor paths clickable on the IP Details page to align with other dashboard views.

## Files to Modify
- `store.go` - `VisitorRegister` and potentially `VisitorCount` helpers.
- `admin/dashboard/handle_load_dashboard.go` - Main dashboard data aggregation logic.
- `admin/dashboard/dashboard.html` & `dashboard.js` - UI updates for new metrics and comparison.
- `admin/sessions/handle_load_sessions.go` - Session grouping logic.
- `admin/sessions/sessions.html` & `sessions.js` - UI updates for navigation flow.
- `admin/visitors/handle_load_visitors.go` - Add CSV export.
- `admin/routes.go` - Register new export routes.

## Verification
- [ ] Run `examples/admin-demo` and verify new cards (Bounce Rate, Duration, Live Visitors) appear.
- [x] Verify Unique Visitor counts match fingerprint-based logic (verified via local tests).
- [ ] Test CSV exports for all pages.
- [ ] Verify session grouping correctly splits visits after 30 minutes of inactivity.
- [ ] Check period comparison indicators with seeded historical data.

## Risks/Considerations
- **Performance**: Fetching 10,000 visitors to group them in-memory might be slow for high-traffic sites. We should ensure we are using indexes on `fingerprint` and `created_at`.
- **Database Compatibility**: Ensure that any new aggregation queries remain compatible with SQLite, PostgreSQL, and MySQL.
