# Task: Custom Events Tracking

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Allow tracking of named events beyond page views (sign-up, purchase, download, etc.) with optional custom properties.

## Features
- Track named events: `event_name` (e.g., "signup", "purchase", "download").
- Custom properties: key-value pairs attached to events (e.g., `plan: "pro"`, `amount: 99`).
- Event tracking endpoint: receive events from client-side JS or server-side calls.
- Event-based filtering in admin pages.
- Event breakdown in dashboard traffic sources.
- Event-based goals (see `goals-funnels.md`).

## Implementation Notes
- **Option A**: Add `event_name` and `event_props` (JSON) fields to existing visitor record. Simpler, no new table.
- **Option B**: Create a separate `events` table with `visitor_id`, `event_name`, `event_props`, `created_at`. More flexible, supports multiple events per visitor.
- Recommend Option B for proper event tracking.
- Client-side JS snippet: `statsstore.track("signup", {plan: "pro"})` → POST to event endpoint.
- Server-side: `store.EventRegister(ctx, event)` method.
- Add `EventList`, `EventCount` store methods with query options.

## Dependencies
- Schema change: new `events` table (Option B) or new fields on visitor table (Option A).
- Client-side JS snippet for event tracking.
- Event endpoint in HTTP handler.

## Acceptance Criteria
- Custom events can be tracked from client-side JS and server-side calls.
- Events stored with name, properties, timestamp, and visitor association.
- Events visible in admin dashboard with breakdown by event name.
- Events filterable by name, date range, and custom properties.
