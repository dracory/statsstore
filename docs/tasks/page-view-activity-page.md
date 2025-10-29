# Task: Page View Activity Page

## Objective
Introduce a dedicated "Page View Activity" screen that mirrors the StatCounter-style reference while leveraging current visitor data captured by `statsstore`.

## Key Features
- **Table columns**
  - Date and time (separate columns) derived from visitor `created_at`.
  - System column with browser and OS icons (reuse existing helper logic from visitor activity).
  - Location / language with country flag icon when available.
  - Host name / page / referrer block showing:
    - Entry path with external link to live site.
    - Reverse DNS / IP lookup link when host name available.
    - Green "(No referring link)" fallback messaging.
- **Row badges** for device category (desktop / mobile / tablet) and session indicators if present.
- **Filters toolbar**
  - Primary "Add Filter" button exposing dropdown for date range, country, device type, browser.
  - Quick range buttons (All, Today, 24 Hours, custom date picker) similar to reference component.
- **Pagination footer** with page indicator and results-per-page selector.
- **Upgrade banner slot** (optional) allowing host application to inject plan messaging.

## Data & Store Requirements
- Confirm `VisitorList` exposes IP, path, referrer, user agent, country, and device metadata.
- Add helper to derive host name from IP/address (using existing or new utility).
- Implement translation layer for browser/OS to icons.

## UI Implementation Notes
- Build page as new controller under `admin/page-view-activity` with shared layout wiring.
- Reuse `shared.AdminHeaderUI` for navigation and add breadcrumb entries.
- Create reusable components for:
  - Filter bar (shared with other pages over time).
  - Host/referrer cell with multi-line formatting.
- Ensure table supports hover tooltips for long URLs and host information.

## Acceptance Criteria
- Page renders with responsive table matching reference hierarchy.
- Filters update results via query parameters and persist across pagination.
- CSV export uses shared helper with same columns.
- No new authentication or user-role logic inside this module (handled by host system).
