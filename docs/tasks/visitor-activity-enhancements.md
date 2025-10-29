# Task: Visitor Activity Enhancements

## Objective
Evolve the Visitor Activity page to resemble the StatCounter reference while staying within the current data model and ensuring the experience integrates cleanly with the host admin shell.

## Key UI Elements (Reference Alignment)
- **Visit summary stack** per row showing:
  - Page views count and session indicator badges.
  - Exit time and resolution (hide gracefully if unknown).
  - Browser + OS icons and version metadata.
- **Session detail panel** on the right side of each row presenting:
  - Country flag + location text.
  - ISP/IP address with optional lookup link.
  - Referrer link block: primary referrer URL (green when "No referring link").
- **Interactive controls** mirroring reference:
  - “Add Filter” dropdown for quick filters (date range, country, device type).
  - Export and Options buttons (Export wired to shared helper; Options placeholder for future features).
  - Footer controls: page indicator, quick date range buttons (All, 24 Hours, Today, custom range), date picker, results-per-page selector.

## Deliverables
- Restructure table markup (or card layout) to support two-column row design similar to screenshot.
- Implement expandable detail or stacked layout for smaller screens while preserving reference layout on desktop.
- Add flag/device/browser icon helpers to visually align with reference.
- Ensure filters load results via query params and persist through pagination/export.
- Integrate upgrade banner slot in footer (non-functional placeholder for host messaging).
- Update CSV export to include all surfaced fields (page views, exit time, system, location, referrer, IP).

## Dependencies
- `statsstore.StoreInterface` must expose page views count if available; otherwise document placeholder logic and follow-up task.
- Shared export improvements task for unified CSV handling.
- Potential new helper for ISP/IP lookup links.

## Acceptance Criteria
- Desktop layout closely matches the reference hierarchy; mobile layout remains readable via collapse/stack.
- Filters, pagination, and exports remain functional and stateful.
- No regressions to existing visitor activity data fetching.
- Documentation updated (`admin-overview.md`) with screenshots/wireframes once implemented.
