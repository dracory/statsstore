# Rybbit.com — Research Summary

## Overview
Rybbit is an open-source (AGPL-3.0), cookieless, privacy-friendly alternative to Google Analytics. 18KB script, GDPR/CCPA compliant, no consent banner needed. Focuses on being intuitive and readable.

## Core Web Analytics
- **Visitor profiles**: Device, browser, OS, and location data per visitor
- **Traffic sources**: Referrer tracking with source/medium breakdown
- **Geographic data**: Country → region → city level (3-level location tracking)
- **Custom events**: Track sign-ups, purchases, downloads, any custom interaction
- **Custom data**: JSON properties attached to events
- **Advanced filtering**: Across 15+ dimensions

## Key Metrics
- Sessions, unique users, pageviews, bounce rate, session duration

## Advanced Analytics
- **Session replay**: Watch real user sessions (DOM mutations, not screenshots). Filter by duration, pageviews, events, country, device, custom property. Privacy: sensitive form inputs masked, async loading, zero performance impact.
- **Web vitals**: Core Web Vitals from real visits, breakdown by route, country, device
- **Funnels**: Visualize conversion paths, pinpoint where visitors drop off
- **User journeys**: Map how users navigate from landing to conversion (Sankey diagrams)
- **User sessions**: Follow complete user journeys from first visit to conversion
- **Goals**: Customizable goals with conversion tracking
- **Retention**: User retention analysis
- **Error tracking**: Capture and debug errors with session replay context

## Real-Time
- Real-time dashboard with live visitor count

## Maps
- Advanced map visualizations with 3-level location tracking

## Comparison (Rybbit vs others)
| Feature | Rybbit | GA4 | Plausible | Cloudflare |
|---|---|---|---|---|
| Open Source | ✅ | ❌ | ✅ | ❌ |
| Self-Hosting | ✅ | ❌ | ✅* | ❌ |
| Cookieless | ✅ | ❌ | ✅ | ✅ |
| Advanced Maps | ✅ | ❌ | ❌ | ❌ |
| Advanced Filters | ✅ | ⚠️ | ⚠️ | ❌ |
| Web Vitals | ✅** | ❌ | ❌ | ❌ |
| Session Details | ✅ | ❌ | ❌ | ❌ |
| User Profiles | ✅ | ❌ | ❌ | ❌ |
| Session Replays | ✅ | ❌ | ❌ | ❌ |
| Funnels | ✅ | ✅ | ✅** | ❌ |
| User Journeys | ✅ | ✅ | ❌ | ❌ |
| Retention | ✅ | ✅ | ❌ | ❌ |
| Goals & Events | ✅ | ✅ | ✅ | ❌ |
| Real-time | ✅ | ✅ | ✅ | ✅ |
| Custom Events (JSON) | ✅ | ✅ | ⚠️ | ❌ |
| Error Tracking | ✅ | ❌ | ❌ | ❌ |
| Public Dashboards | ✅ | ❌ | ✅ | ❌ |

## API & Developer Features
- REST API for all metrics (sessions, metrics, funnels, goals, errors, events)
- MCP server for agent-based access
- npm package (@rybbit/js)
- Docker Compose self-hosting

## Autocapture
- Clicks, form submits, copied text, outbound links, and errors captured automatically with zero instrumentation

## Design Principles
- One readable dashboard
- Intuitive and fast
- Privacy-first: no cookies, no consent banner
- Open source, self-hostable

## Sources
- https://rybbit.com/features
- https://github.com/rybbit-io/rybbit
- https://rybbit.com/
- https://rybbit.com/for-developers
- https://rybbit.com/features/session-replay
