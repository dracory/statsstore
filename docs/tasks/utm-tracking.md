# Task: UTM Campaign Tracking

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Parse and store UTM campaign parameters from visitor URLs to enable campaign performance reporting.

## Features
- Parse UTM parameters from request URL query string:
  - `utm_source` — traffic source (Google, newsletter, etc.)
  - `utm_medium` — marketing medium (cpc, email, social, etc.)
  - `utm_campaign` — campaign name
  - `utm_term` — search term (paid search)
  - `utm_content` — ad variant or content identifier
- Also parse Plausible-supported alternatives: `ref` and `source` params.
- Store as fields on the visitor record (schema change).
- Display UTM data in visitor activity and page view activity rows.
- Add UTM-based filtering to query options.
- Add UTM breakdown to dashboard traffic sources.

## Implementation Notes
- Add 5 new fields to `VisitorInterface`: `UtmSource`, `UtmMedium`, `UtmCampaign`, `UtmTerm`, `UtmContent`.
- Add corresponding columns to the visitor table schema (automigrate).
- Parse UTM params in `VisitorRegister` from the request URL.
- Add `SetUtmSource()`, `SetUtmMedium()`, `SetUtmCampaign()` etc. to `VisitorQueryInterface`.
- Add UTM columns to CSV export.
- Add UTM filter options to filter toolbar.

## Dependencies
- Schema migration (new fields on visitor table).
- Store-level query support for UTM fields.

## Acceptance Criteria
- UTM parameters are parsed from URL and stored on visitor records.
- UTM data visible in visitor activity rows and CSV export.
- Dashboard traffic sources include UTM campaign breakdown.
- UTM-based filtering works in admin pages.
