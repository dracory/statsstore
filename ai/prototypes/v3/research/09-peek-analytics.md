# Peek Analytics — Research Summary

## Overview
Peek is a privacy-first, cookie-free web analytics platform. Tracker script is under 1KB (gzipped), 30x lighter than Google Analytics. GDPR & CCPA clean. Running on 2,400+ sites.

## Core Features

### Real-time Visitors
- See who's on the site and what pages they're reading, the moment it happens
- Every visit appears live, no refresh needed
- Live visitor count

### Pages & Sources
- Top pages, entries, exits
- Where every visit came from (referrers)
- Spot best pages and channels in seconds
- Nothing buried — all visible at a glance

### Countries & Devices
- Countries behind traffic at a glance
- Browsers and devices breakdown
- Know exactly who you're reaching and on what

### Funnels
- Follow the steps that matter
- See exactly where people drop off
- Find the leak before it costs signups
- Steps must be visited in order within same session
- Pages autocomplete from real traffic data
- Add Purchase step for revenue tracking
- Sequential matching with real page data

### Journeys
- Track how visitors move through the site, step by step
- See the real path people take to convert
- No tag soup

### UTM Campaigns
- Know which links, posts and emails actually brought people
- Tag a link once and Peek does the rest
- UTM attribution: source, medium, campaign, content, term
- Last non-direct attribution model

### Revenue Tracking
- Connect Stripe or Shopify
- Attribute sales to campaigns and traffic sources
- Total revenue, order count, average order value
- Breakdowns by source, country, landing page

### Forms
- Embedded lead capture with full UTM attribution
- Automatic capture of UTM parameters and landing page URL

## Privacy Model
- **No cookies**: No cookies, no localStorage, no client-side storage of any kind
- **No consent banner needed**: No data stored on visitor's device
- **Server-side identification**: Salted, rotating hash from IP + User-Agent + domain
- **Monthly salt rotation**: Returning visitors recognized within calendar month
- **Raw data discarded**: IP address never stored, logged, or written to database
- **GDPR & CCPA clean**: No personal data processed
- **Respects opt-out**: `peek_ignore` localStorage flag for developer self-exclusion

## Tracker Script
- Under 1KB gzipped
- Loads asynchronously
- Single script tag in `<head>`
- No SDK, no config, no cookie banner
- Automatic: page URL capture, referrer capture, SPA navigation interception, outbound link detection, Web Vitals measurement

## Server-Side Processing
- Geo-IP lookup (country, city, region)
- User-agent parsing (browser, OS, device)
- Salted visitor hash generation
- Session creation & grouping
- Site key validation
- Data persistence

## Design Principles
- **Minimal & calm**: Under 1KB, no clutter, just the numbers you need
- **Privacy by default**: No cookies, no consent, no personal data, every visit anonymous
- **Real-time**: See visitors the moment they land
- **Simple setup**: One script tag, no config
- **You own your data**: No data selling

## What You Get (Cookieless)
- Page views and unique visitors (accurately counted without consent bias)
- Traffic sources (referrer data, UTM parameters, campaign tracking)
- Geographic data (country, region from IP, not stored)
- Device and browser breakdown (from User-Agent)
- Top pages and entry/exit pages
- Session duration and bounce rate

## What You Don't Get
- Individual user journeys across multiple sessions
- Cross-device tracking
- Remarketing audiences

## Sources
- https://trypeek.dev/
- https://trypeek.dev/docs
- https://trypeek.dev/funnels
- https://trypeek.dev/blog/track-website-visitors-without-cookies
