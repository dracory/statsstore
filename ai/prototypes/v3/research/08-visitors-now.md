# Visitors.now — Research Summary

## Overview
Visitors.now is a privacy-friendly, cookie-free, GDPR-compliant Google Analytics alternative. Features realtime visitor tracking with interactive 3D globe, streaming activity feed, and visitor profiles.

## Core Features

### Realtime Tracking
- Live visitor count
- Interactive 3D visitor globe (WebSocket streaming, no polling)
- Streaming activity feed
- City-level precision (without exact coordinates for privacy)
- Live visitor profiles (click marker to see journey, device, current page)

### Pageviews & Sessions
- Automatic tracking of every page visit
- Session duration tracking
- Bounce rates
- Total pages viewed across all sessions
- Top pages (entered, exited)
- Heartbeat events for session continuity and page duration tracking

### Custom Events
- Track button clicks, form submissions, purchases, any interaction

### Outbound Links
- See when visitors click links leading away from the site

### Traffic Sources
- Referrers, search engines, social platforms
- UTM campaigns (source, medium, campaign, content, term)
- Referrer type: search, direct, social, etc.

### Geographic Data
- Countries, regions, cities
- City center coordinates used (not exact locations) for privacy
- Geoname IDs stored, not precise coordinates

### Devices & Browsers
- Desktop, mobile, tablet
- Browser name and version
- OS name and version
- Device vendor and model
- CPU architecture
- Browser engine

### Visitor Profiles
- Full session history
- Cross-device tracking (optional)
- Identify logged-in users (optional)
- Anonymous signatures (salted hash, rotating)

## Privacy Model
- **Cookie-free by default**: No cookies set unless persist mode or identify function enabled
- **GDPR, CCPA, PECR compliant** out of the box
- **IP discarded immediately**: Used only for geo-lookup, then deleted
- **Anonymous signatures**: Salted hash of IP + UA + project token, rotating
- **Respects DNT and GPC**: Script won't send events if Do Not Track or Global Privacy Control enabled
- **No cross-site tracking**

## Event Types
| Event Type | Description |
|---|---|
| Pageview | Automatically tracked on page visit |
| Identify | Sent when user logs in, connects anonymous signature to profile |
| Outgoing | Tracked on external link clicks |
| Custom | Any custom interaction (button clicks, form submissions) |
| Heartbeat | Periodically sent for page durations and session continuity |
| Performance | Core Web Vitals (LCP, CLS, INP, FCP, TTFB) |

## Data Storage
- **Events**: Raw activity log with anonymous signature, browser info, location, page details
- **Sessions**: Aggregated session data (duration, pages viewed, bounce rates)
- **Profiles**: Anonymous visitor profiles (optionally with personal info if identify used)
- **Performance**: Core Web Vitals per page visit

## Comparison (from their site)
| Feature | Visitors | GA | Plausible | Fathom |
|---|---|---|---|---|
| GDPR compliant | ✅ | ❌ | ✅ | ✅ |
| No cookie banner | ✅ | ❌ | ✅ | ✅ |
| Revenue attribution | ✅ | ❌ | ❌ | ❌ |
| Realtime analytics | ✅ | ✅ | ✅ | ❌ |
| Visitor journeys | ✅ | ❌ | ❌ | ❌ |
| User identification | ✅ | ❌ | ❌ | ❌ |
| Performance insights | ✅ | ❌ | ❌ | ❌ |
| Funnel analysis | ✅ | ✅ | ✅ | ❌ |

## Design Principles
- Realtime-first: see visitors the moment they land
- Privacy-first: no cookies, no personal data without consent
- Interactive globe visualization for geographic data
- Streaming activity for live monitoring
- Simple setup, lightweight script

## Sources
- https://visitors.now/analytics
- https://visitors.now/
- https://visitors.now/data
- https://visitors.now/realtime
- https://visitors.now/docs/api/get-realtime
