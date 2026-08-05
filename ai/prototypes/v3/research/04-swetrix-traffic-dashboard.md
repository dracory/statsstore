# Swetrix Traffic Dashboard — Research Summary

## Overview
Swetrix is a privacy-focused, open-source analytics platform with cookieless tracking. The Traffic Analytics dashboard is the core view, providing comprehensive insights into visitors, behavior, and traffic sources.

## Essential Metrics
- **Unique visitors (Sessions Unique)**: Individual users within timeframe
- **Pageviews**: Total pages viewed
- **Session Duration**: Average time per session
- **Bounce Rate**: Percentage of single-page visits
- **Views per Unique**: Average pages per unique visitor
- **Revenue**: Total revenue (if revenue tracking enabled)

## Main Chart
- Visualizes various metrics over time
- Click any data point to drill down into sessions for that timeframe
- Time buckets: Hour, Day, Week, Month, Year

## Pages
- **Page Paths**: Specific URLs visited
- **Entry Pages**: First page visitors land on
- **Exit Pages**: Last page before leaving
- **User Flow**: Visualize paths through the site
- **Host**: Cross-domain tracking support

## Traffic Sources
- **Referrers**: External websites linking to yours
- **Source / Medium**: Detailed acquisition channel breakdown (e.g., `google / organic`, `newsletter / email`)
- **Campaigns**: UTM parameter tracking

## Network Intelligence
- **ISP**: Internet Service Provider serving the visitor
- **Organisation**: Registered owner of IP block
- **Usage type**: residential, business, hosting, cellular, government, school, search_engine_spider
- **Connection type**: Cable/DSL, Cellular, Corporate, Satellite, Dialup

## Custom Events
- Dedicated panel for tracked events (e.g., "Button Clicked", "Form Submitted")
- Metadata breakdowns for pageviews and custom events

## Time Periods
- Standard: Today, Yesterday, Last 7 days, Last 30 days, Month to Date
- Custom range selection
- Time bucket grouping: Hour, Day, Week, Month, Year

## Live Visitors
- Real-time visitor count (unique visitors active in last 5 minutes)
- Pulsing green live indicator
- Updates in real-time as visitors arrive/leave
- Historical concurrency charting (peak concurrent visitors)
- Browser tab title display option: `(12) Project Name | Swetrix`

## Segments & Export
- Save filter configurations as Segments
- Export current view data to CSV

## Design Principles
- Put most important stuff first (5-second health check)
- Group related metrics logically
- Match visualization type to metric (line chart for trends, bar chart for comparisons)
- Privacy-compliant metrics without invasive tracking
- Customizable alerts for traffic spikes, errors, significant events
- Drill-down capability from charts to individual sessions

## Key Design Takeaways
- "Who is visiting my site? Where are they coming from? What are they doing? Are they doing what I want?"
- Group metrics into categories that answer simple, direct questions
- Most critical KPIs at the top
- Privacy-first: cookieless, GDPR-compliant, no consent banner needed

## Sources
- https://swetrix.com/docs/analytics-dashboard/traffic
- https://swetrix.com/blog/traffic-dashboard
- https://swetrix.com/docs/analytics-dashboard/live-visitors
- https://swetrix.com/docs/traffic-sources
- https://swetrix.com/blog/web-analytics-dashboard
