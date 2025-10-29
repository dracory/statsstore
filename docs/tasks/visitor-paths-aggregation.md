# Task: Visitor Paths Experience

## Objective
Upgrade the Visitor Paths page to mirror the StatCounter reference by combining per-session details with quick path navigation aides, while using existing visitor data where possible.

## Key UI Elements (Reference Alignment)
- **Row header** showing country flag + location, host name, and page path with external link icon.
- **Timestamps** for entry/exit moments stacked vertically under each row header.
- **Referrer block** with primary link in green for "No referring link" fallback.
- **Session metadata** column indicating session number, magnifier icon for drill-down, and device/browser badges aligned to the right.
- **Filter / control bar** including:
  - “Add Filter” dropdown (date range, country, path contains, device type).
  - Export button using shared helper.
- **Footer controls**: pagination status, quick range buttons (All, 24 Hours, Today, custom range), date selection controls, and per-page selector.
- **Upgrade banner region** centred within list for host-provided messaging.

## Deliverables
- Restructure view layout to render list items matching reference hierarchy (two-column flex/table hybrid).
- Implement reusable partial for row header so future pages can share location/flag/host formatting.
- Enable drill-down action (link or modal) from session icon to open detailed visitor activity for that session (stubs acceptable for initial release).
- Introduce path-level filters (contains/exact) executed through query parameters and backed by store queries.
- Ensure export includes location, timestamps, path, referrer, session number, device/browser data.

## Data / Store Considerations
- `statsstore.StoreInterface` should provide path, timestamps, country, referrer, user agent, and session identifier; document missing fields and create follow-up tasks if necessary.
- If session identifiers are not currently stored, derive temporary grouping based on fingerprint+timestamps; note limitations in documentation.
- Leverage helper from visitor activity task for icons and flag rendering to keep visuals consistent.

## Dependencies
- Shared export improvements task.
- Potential store enhancements for session grouping/path filtering.

## Acceptance Criteria
- Page layout matches reference structure on desktop and remains usable on mobile via stacked/collapsible layout.
- Filters update URL/query params and persist through pagination/export.
- Export output matches on-screen data columns.
- Admin overview documentation updated with screenshots/wireframes after implementation.
