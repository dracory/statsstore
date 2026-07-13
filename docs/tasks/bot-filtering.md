# Task: Bot Filtering

## Source
Plausible + Rybbit feature research

## Status: Completed

## Objective
Filter out bot/crawler traffic and referrer spam to prevent data pollution.

## Features
- **User-agent bot detection** — match against known bot/crawler signatures (Googlebot, Bingbot, Slurp, DuckDuckBot, Baiduspider, YandexBot, AhrefsBot, SemrushBot, etc.)
- **Data center traffic filtering** — optionally filter requests from known data center IP ranges (AWS, GCP, Azure).
- **Referrer spam blocking** — maintain a blocklist of known referrer spam domains.
- **Configurable** — allow host application to enable/disable filtering and add custom patterns.

## Implementation Notes
- Add an `IsBot(userAgent string) bool` function in the `statsstore` package.
- Call `IsBot()` in `VisitorRegister` (or equivalent ingestion point) to skip bot visits.
- Maintain a static list of bot user-agent substrings (e.g., `bot`, `crawler`, `spider`, `scraper`, `fetch`, `monitor`).
- Consider using an existing Go library (e.g., `github.com/mssola/user_agent` or a curated blocklist).
- Referrer spam list can be a static set of domains, updated periodically.
- Optionally store bot visits in a separate table or flag them instead of dropping (for debugging).

## Dependencies
- No schema changes if bots are dropped at ingestion time.
- If storing bot visits with a flag, add `is_bot` boolean field to visitor record.

## Acceptance Criteria
- Known bots/crawlers are filtered out and do not appear in visitor lists or dashboard stats.
- Referrer spam domains are blocked.
- Filtering can be toggled on/off by host application.
- Bot detection patterns are maintainable (static list or external config).
