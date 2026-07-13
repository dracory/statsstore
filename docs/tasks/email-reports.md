# Task: Email Reports and Notifications

## Source
Plausible + Rybbit feature research

## Status: Not Started

## Objective
Send scheduled summary reports and traffic spike notifications via email and/or Slack.

## Features
- **Scheduled reports**: daily, weekly, or monthly summary emails.
- **Report content**: page views, visitors, bounce rate, top pages, top referrers, top countries, device breakdown.
- **Traffic spike notifications**: alert when traffic exceeds a threshold (e.g., 2x average).
- **Slack integration**: send reports and alerts to Slack webhook.
- **Template-based**: HTML email templates with stats tables.
- **Configurable**: host application defines recipients, schedule, and thresholds.

## Implementation Notes
- Requires a scheduler/cron mechanism (Go ticker, or host app's job queue).
- Email sending via SMTP or host app's email service.
- Report generation: query stats for the period, render HTML template.
- Slack: POST to incoming webhook URL.
- New package: `statsstore/reports/` or integrate with host app's notification system.
- Config: report schedule, recipients, Slack webhook URL, spike threshold.

## Dependencies
- Email sending infrastructure (SMTP or host app integration).
- Scheduler/cron mechanism.
- HTML email template rendering.

## Acceptance Criteria
- Scheduled reports sent daily/weekly/monthly with correct stats.
- Traffic spike notifications triggered when threshold exceeded.
- Reports delivered via email and/or Slack.
- Report content matches dashboard data for the same period.
- Schedule and recipients configurable.
