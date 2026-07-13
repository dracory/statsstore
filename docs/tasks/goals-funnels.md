# Task: Goals and Funnels

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Define conversion goals and funnels to measure visitor progression through defined steps.

## Features

### Goals
- Define conversion goals:
  - **Page visit goals**: "visited /thank-you page"
  - **Custom event goals**: "completed signup event" (depends on `custom-events.md`)
  - **Automated goals**: 404 visits, outbound link clicks, file downloads (depends on `automated-tracking.md`)
- Track goal completion rate: `conversions / total visitors * 100`.
- Revenue goals: attach monetary value to goal completions (Plausible feature).
- Goal configuration admin page: create, edit, delete goals.
- Goal dashboard: conversion rate, total conversions, trend over time.

### Funnels
- Define a sequence of pages or events as funnel steps.
- Measure drop-off at each step: `visitors at step N / visitors at step N-1`.
- Plausible supports strict order funnels (consecutive steps, no other actions in between).
- Funnel visualization: bar chart or step-by-step breakdown.
- Funnel configuration admin page.
- Built on top of visitor paths data (already available).

## Implementation Notes
- Goals require a `goals` config table: `id`, `name`, `type` (page_visit, event, automated), `condition` (path or event name), `value` (optional revenue).
- Funnel steps stored in a `funnels` + `funnel_steps` table.
- Goal completion check: query visitor records for matching path/event within selected period.
- Funnel computation: for each fingerprint, check if they visited all steps in order.
- New admin pages: `/admin/goals`, `/admin/funnels`.
- New dashboard cards: goal conversion rate, funnel drop-off.

## Dependencies
- `custom-events.md` for event-based goals.
- `automated-tracking.md` for automated goals (404, downloads, outbound links).
- Schema changes: goals and funnels config tables.

## Acceptance Criteria
- Goals can be created, edited, and deleted from admin.
- Goal conversion rate shown on dashboard.
- Funnels can be defined with ordered steps.
- Funnel drop-off visualized in admin.
- Goal and funnel data filterable by time period.
