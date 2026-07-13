# Task: Automated Tracking — Outbound Links, Downloads, Forms, Scroll Depth, 404s

## Source
Plausible feature research

## Status: Not Started

## Objective
Automatically track visitor interactions beyond page views: outbound link clicks, file downloads, form submissions, scroll depth, and 404 error pages.

## Features

### Outbound Link Click Tracking
- Track clicks on external links automatically.
- JS snippet detects clicks on `<a>` tags with `target="_blank"` or external domain.
- Send click event with link URL to tracking endpoint.

### File Download Tracking
- Track downloads of common file types (PDF, ZIP, DOCX, XLSX, etc.).
- JS snippet detects clicks on links ending with known file extensions.
- Send download event with file URL to tracking endpoint.

### Form Submission Tracking
- Track form completions automatically.
- JS snippet hooks into `submit` events on `<form>` elements.
- Send form submission event with form ID/action to tracking endpoint.

### Scroll Depth Tracking
- Track how far visitors scroll on each page (0–100%).
- JS snippet measures scroll position and reports max scroll depth on page leave (via `visibilitychange` or `beforeunload`).
- Report scroll depth as a percentage on the visitor record or as a custom event.

### 404 Error Page Tracking
- Track visits to 404 pages.
- Detect server-side (HTTP 404 status) or client-side (page content matching).
- Store as a special event type or flag on the visitor record.

## Implementation Notes
- All automated tracking requires a client-side JS snippet loaded on the tracked site.
- JS snippet sends events to a tracking endpoint (same as page view tracking, with `event_type` parameter).
- Depends on `custom-events.md` infrastructure for event storage.
- Plausible offers these as plug-and-play toggles — consider a config object for enabling/disabling each tracker.
- Scroll depth can be stored as a custom property on the page view event.

## Dependencies
- `custom-events.md` for event storage and tracking endpoint.
- Client-side JS snippet (tracking script) deployed on the monitored site.

## Acceptance Criteria
- Outbound link clicks tracked and visible in admin.
- File downloads tracked and visible in admin.
- Form submissions tracked and visible in admin.
- Scroll depth recorded per page view and visible in admin.
- 404 page visits tracked and visible in admin.
- Each tracker can be toggled on/off via configuration.
