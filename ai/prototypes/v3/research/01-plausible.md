# Plausible.io — Research Summary

## Overview
Plausible is a privacy-first, cookieless, open-source (AGPL-3.0) web analytics platform. It focuses on simplicity — one dashboard with essential stats, no sub-menus, no custom reports needed.

## Core Metrics
- **Unique visitors**: Distinct visitor count
- **Total visits (sessions)**: Session count
- **Total page views**: All pageviews
- **Views per visit**: Average pages per session
- **Bounce rate**: Percentage of sessions with only one page view (unless custom event triggered)
- **Visit duration**: Average session length (bounces counted as 0 seconds)
- **Scroll depth**: Percentage of page height reached on average (auto-tracked, 1-100%)

## Traffic Sources
- Referral sources ranked by unique visitors
- UTM parameter support: `utm_source`, `utm_medium`, `utm_campaign`, `utm_content`, `utm_term`
- Also supports `ref`, `source` query parameters
- Automatic channel grouping (Paid Search, Affiliates, etc.)
- Campaigns tab for UTM-tagged traffic
- Click-ID detection (gclid, fbclid, etc.)

## Pages
- Top pages ranked by unique visitors
- Entry pages (landing pages) with visit duration
- Exit pages with exit rate percentage
- Page-level metrics: pageviews, bounce rate, time on page, scroll depth

## Goals & Conversions
- Goals based on page visits or custom events
- Conversion rate, unique conversions, total conversions
- Revenue tracking (ecommerce)
- Custom properties (custom dimensions) attached to events

## Funnels & Journeys
- **Funnels**: Define sequence of steps (pageviews or events), see drop-off. Sequential or strict order mode.
- **User Journeys**: Explore actual paths visitors take — no predefined sequence. Start from any page/event, follow forward or backward.

## Auto-Tracked Features
- 404 error page tracking
- File downloads
- Outbound link clicks
- Form submissions
- Scroll depth (automatic, no setup)

## Real-Time
- Real-time dashboard with same metrics as main dashboard
- Updates every 30 seconds without page refresh

## Dashboard Design
- Single page with all essential stats
- No layers of menus
- Click any metric to display in top graph
- Expandable sections for full lists with additional metrics
- Period comparison between different time ranges
- Saved segments for audience filtering
- Annotations on chart

## Other Features
- Stats API for custom dashboards
- Events API for server-side tracking
- Email/Slack reports
- Google Analytics import
- Consolidated view for multiple sites
- Public/shared dashboards
- Custom properties for enriched event data

## Design Principles
- **Simplicity**: One dashboard, essential insights, no noise
- **Lightweight**: Script is 54x smaller than GA (saves ~135KB JS per visitor)
- **Privacy**: No cookies, no cross-site tracking, no cross-device tracking, EU-hosted
- **Speed**: Fast-loading dashboard

## Sources
- https://plausible.io/docs/your-plausible-experience
- https://plausible.io/
- https://plausible.io/docs/metrics-definitions
- https://plausible.io/docs/custom-event-goals
- https://plausible.io/docs/guided-tour
