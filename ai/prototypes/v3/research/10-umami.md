# Umami — Research Summary

## Overview
Umami is a simple, fast, privacy-focused, open-source alternative to Google Analytics. No cookies, no fingerprinting, no personal data collection. GDPR compliant out of the box. Self-hosted with Docker. Tracking script under 2KB.

## Core Analytics
- **Pageviews**: Track every page visit with automatic path detection
- **Unique Visitors**: Count distinct visitors without invasive tracking
- **Bounce Rate**: Percentage of single-page sessions
- **Session Duration**: Time spent on site
- **Referrers**: Where traffic comes from (referral sites, search engines, direct)
- **Browsers**: Chrome, Firefox, Safari, Edge + versions + market share
- **Operating Systems**: OS breakdown
- **Devices**: Desktop, mobile, tablet + distribution percentages
- **Countries & Regions**: Geographic distribution (country, region, city level)
- **Map visualization**: Visual geographic heatmap

## Custom Events
- Track button clicks, form submissions, video plays, downloads, any interaction
- Custom event names
- Event properties/metadata
- Event counts and trends
- Event conversion funnels

## Real-Time Data
- Current active visitors on site
- Real-time pageview stream
- Live event tracking
- Instant metric updates

## Sessions
- View individual visitor activity
- Session properties without identifying personal information

## Multi-Website Management
- Track unlimited websites from single instance
- Individual dashboards per site
- Separate tracking codes
- Aggregate reporting across sites
- Team access controls

## Filtering
- Specific pages or page groups
- Traffic sources
- Countries or regions
- Browsers or devices
- Custom event properties

## Privacy Features
- **No cookies**: No cookie consent banners needed
- **No cross-site tracking**: Visitors not tracked across different websites
- **No persistent identifiers**: No cookies or localStorage tracking
- **IP address anonymization**: IPs hashed and not stored
- **No fingerprinting**: No browser fingerprinting techniques
- **GDPR compliant**: Meets EU data protection requirements
- **No data selling**: Data never shared or sold

## API
- Full REST API for programmatic access
- CRUD operations
- Real-time data access
- Custom integrations
- Automated reporting
- Third-party tool integration

## Additional Features
- User journeys: trace exact path visitors take through site
- Session replays: watch what happened (where they got confused, what they clicked)
- Goals and conversions
- Revenue tracking

## Dashboard Design
- Key metrics on one dashboard: pageviews, visitors, bounce rate, average visit time
- Pick any date range
- Filter by country or device
- Full picture in seconds
- Clean, intuitive interface

## Design Principles
- **Simple**: Deploy with Docker in minutes, no complex configuration
- **Fast**: Lightweight script, clean interface
- **Privacy-focused**: No cookies, no personal data, GDPR compliant
- **Own your data**: Self-host on your infrastructure, data never leaves servers
- **Lightweight**: Under 2KB tracking script

## Sources
- https://docs.umami.is/docs
- https://umami.is/website-analytics
- https://github.com/umami-software/umami/
- https://www.mintlify.com/umami-software/umami/features
