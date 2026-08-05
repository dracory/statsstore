# Swetrix 2025 Visitor Tracking Guide — Research Summary

## Overview
Swetrix is a privacy-focused, open-source analytics platform using cookieless tracking. This research focuses on their 2025 visitor tracking methodology and GDPR compliance approach.

## Cookieless Identification Method
Swetrix never sets cookies, never writes to `localStorage`, never reads client-side identifiers. Instead:

1. **Collect request attributes**: IP address, User-Agent string, project ID
2. **Add rotating salt**: Secret rotating string known only to backend
3. **Hash the combination**: One-way hash function (SHA-256) produces opaque ID like `anon_8214637194021987452`
4. **Discard raw data**: IP/UA never stored; only hash retained

### Salt Rotation
- **Session salt**: Rotates every 24 hours at UTC midnight — detects same-session pageviews
- **Profile salt**: Rotates monthly — detects returning visitors (DAU/WAU/MAU)
- When salt rotates, previous fingerprints become permanently unlinkable

### Accuracy
- IP+UA fingerprint is unique per visitor in 95%+ of cases for typical sites
- Edge cases: shared networks (NAT, corporate, schools) may undercount unique visitors
- Optional `profileId` for per-user accuracy when you control user identity

## GDPR/Privacy Compliance
- **No cookies**: Completely cookieless, no client-side identifiers
- **No PII**: No personal data or personally identifiable information collected
- **IP discarded**: Raw IP never stored; only salted hash retained temporarily (30 min or UTC midnight)
- **No consent banner needed**: No cookies means no consent required under ePrivacy/GDPR
- **Data minimization**: Only essential analytics data collected
- **Respects DNT**: `respectDNT` option to not collect data from Do Not Track users

## Tracking Methods Comparison
| Method | How It Works | Privacy Impact | Advantage | Disadvantage |
|---|---|---|---|---|
| Third-Party Cookies | File on browser, cross-site tracking | High | Cross-site retargeting | Blocked by modern browsers, invasive |
| Server-Side Tracking | Data to your server first, then analytics | Low-Medium | Improved accuracy/control | More complex setup |
| Cookieless Tracking | Anonymous fingerprinting from browser data | Very Low | GDPR/CCPA safe | Less accurate for long-term ID |

## Core Privacy Law Principles
- **Consent**: Clear, explicit permission before non-essential cookies/personal data
- **Data Minimization**: Collect only what's needed
- **Purpose Limitation**: Be upfront about why data is collected, stick to it
- **User Rights**: Visitors can see, correct, delete their data

## Features
- **Custom events**: Track any interaction (sign up, purchase, etc.)
- **Goals**: Pageview goals, custom event goals, multi-condition goals
- **Conversion tracking**: Conversion rates, time-to-convert metrics
- **Heatmaps**: (via session replays on cloud)
- **Session recordings**: Opt-in session replays (DOM mutations, not screenshots)
- **Metadata**: Key-value pairs on events for additional context

## Design Principles
- Privacy-first by architecture, not by policy
- Cookieless from the ground up
- GDPR/CCPA compliant without compromises on insights
- Anonymous aggregated data, no individual profiling
- Open source for transparency and auditability

## Sources
- https://swetrix.com/docs/visitor-identification
- https://swetrix.com/data-policy
- https://swetrix.com/blog/website-visitor-tracking
- https://swetrix.com/docs/swetrix-js-reference
- https://github.com/Swetrix/swetrix/blob/main/docs/content/docs/analytics-dashboard/goals.mdx
