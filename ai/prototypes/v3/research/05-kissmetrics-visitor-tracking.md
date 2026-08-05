# KissMetrics Visitor Tracking — Research Summary

## Overview
KissMetrics is a person-based analytics platform. Unlike session-based tools, it treats each person as the primary unit, stitching together all actions across sessions, devices, and channels into a single timeline.

## Data Model
Three components:
1. **People**: Representation of a user. Assigned anonymous identifier first, then aliased to known identity (e.g., email) when they log in.
2. **Events**: Actions with name, associated person, timestamp, and properties.
3. **Properties**: Key-value pairs describing events and people.

## Session-Based vs. Person-Level Tracking

### Session-Based Tracking
- Each visit is independent (starts on arrival, ends after ~30 min inactivity)
- Self-contained data point: source, pages viewed, actions, duration
- Same person visiting 3 times = 3 separate, unconnected sessions
- How most analytics tools work (GA, etc.)
- Good for: "How many sessions came from organic search?"
- Bad for: "What is the average sessions before a visitor becomes a customer?"

### Person-Level Tracking
- Every action across sessions, devices, channels stitched into one timeline
- First visit from ad → return via organic → signup → trial → upgrade — all connected
- Enables cross-device and cross-browser tracking

## Tracking Methods

### Cookie-Based Tracking
- Small text files with unique identifier in visitor's browser
- First-party cookies are standard mechanism
- Widely supported, reasonably reliable

### Browser Fingerprinting
- Combines browser/device attributes: screen resolution, fonts, browser version, OS, timezone, language
- Cannot be easily deleted by user
- More invasive than cookies

### Login-Based Identity
- Gold standard for accuracy
- Connects pre-login anonymous behavior to identified profile
- Eliminates cross-device, cookie expiration, and many privacy concerns
- Limitation: only fraction of visitors log in (10-30% for e-commerce)

## Automatically Tracked Events
- **Visited Site**: Triggered on each visit (30-min inactivity window)
- **Search Engine Hit**: Arrival from search engine
- **Ad Campaign Hit**: From UTM-tagged or Adwords campaigns

## Automatically Tracked Properties
- Campaign Source/medium/name (from UTM parameters)
- Referrer URL
- KM Landing Page (first page of session)
- New or Returning visitor status

## Key Events to Track (Engagement Milestones)
1. **Visited Site** — Entry point, record source/campaign/landing page
2. **Signed Up** — Visitor becomes known user (signup method, plan)
3. **Completed Onboarding** — Each step + overall completion
4. **Reached Aha Moment** — Event correlating with long-term retention
5. **Used Core Feature** — Track usage of major features
6. **Logged In** — Login frequency as retention predictor
7. **Started Trial** — Funnel step between signup and payment
8. **Upgraded Plan / Completed Purchase** — Revenue event (plan, billing cycle, amount)
9. **Churned / Cancelled** — Include reason if collected

## User Properties
- Plan, company, role, signup date
- Attached to identified person

## Design Principles
- Person-based, not session-based
- Track engagement milestones as a funnel
- Connect anonymous → identified behavior
- Properties add context for segmentation

## Sources
- https://support.kissmetrics.io/docs/kissmetrics-technical-implementation-overview
- https://www.kissmetrics.io/blog/website-visitor-tracking
- https://support.kissmetrics.io/docs/identities
- https://www.kissmetrics.io/blog/event-tracking-101-getting-started-with-kissmetrics
- https://support.kissmetrics.io/docs/understanding-people-events-and-properties-within-kissmetrics
