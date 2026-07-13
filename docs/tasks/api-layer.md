# Task: REST API Layer

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Expose visitor data and stats via REST/JSON endpoints to enable third-party integrations and custom dashboards.

## Features
- **Stats API**: query metrics (page views, visitors, bounce rate, top pages, top referrers) with date range and filter parameters.
- **Events API**: record events without JS (server-to-server event tracking).
- **Sites API**: manage tracked sites (if multi-site support is added).
- **Authentication**: API key or token-based auth.
- **Rate limiting**: prevent abuse.
- **Documentation**: OpenAPI/Swagger spec or simple API docs.

## Implementation Notes
- New HTTP handler package alongside admin UI (e.g., `admin/api/` or separate `api/` package).
- Reuse existing `StoreInterface` methods for data access.
- JSON responses with consistent error format.
- API key stored in config, validated via middleware.
- Endpoints:
  - `GET /api/v1/stats/summary?from=&to=&country=&device=`
  - `GET /api/v1/stats/pages?from=&to=`
  - `GET /api/v1/stats/referrers?from=&to=`
  - `GET /api/v1/visitors?from=&to=&limit=&offset=`
  - `POST /api/v1/events` (server-side event tracking)
- Consider embedding API routes in the existing `admin.go` router or a separate router.

## Dependencies
- No schema changes for basic stats endpoints.
- `custom-events.md` for the Events API endpoint.
- API key management (config or database).

## Acceptance Criteria
- Stats endpoints return JSON with visitor metrics for specified date ranges.
- Events endpoint accepts server-side event tracking requests.
- API authentication via API key.
- Rate limiting prevents abuse.
- API documentation available.
