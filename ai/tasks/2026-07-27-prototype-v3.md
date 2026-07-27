# Task: Prototype v3 — Visitor Analytics Dashboard (from scratch)

## Objective

Create `/ai/prototypes/v3/` containing a complete static HTML/CSS/JS prototype of a visitor analytics dashboard, designed from scratch (not based on v1/v2 or the current admin UI).

## Steps

### 1. Research subfolder

Create `/ai/prototypes/v3/research/` and save each research source as a separate file:

- `research/01-plausible.md` — Plausible.io (privacy-first, cookieless, AGPL-3.0: bounce rate, visit duration, UTM tracking, custom events, scroll depth, outbound links, file downloads, form submissions, 404 tracking, funnels, goals, journeys, period comparison, public/shared dashboards, API, email reports)
- `research/02-rybbit.md` — Rybbit.com (open-source, cookieless: bounce rate, traffic sources, city-level location, custom events, UTM, link tracking, bot blocking, session replay, web vitals, funnels, goals, journeys, globe views, error tracking, user sessions, retention, API, email reports)
- `research/03-matomo-dashboard.md` — Matomo web analytics dashboard guide (dashboards, KPIs, traffic/engagement/audience/acquisition/behaviour metrics, goal tracking, funnel analysis)
- `research/04-swetrix-traffic-dashboard.md` — Swetrix traffic dashboard guide (essential metrics: unique visitors, pageviews, sessions, bounce rates, conversion data, channel attribution, real-time error monitoring, privacy-compliant metrics, customization, alerts, filters, drill-down)
- `research/05-kissmetrics-visitor-tracking.md` — KissMetrics visitor tracking guide (session-based vs person-level tracking, cookies/fingerprinting/login-based identity, acquisition source, conversion events, engagement milestones)
- `research/06-swetrix-visitor-tracking-2025.md` — Swetrix 2025 visitor tracking guide (cookieless tracking, GDPR compliance, event/goal tracking, heatmaps, session recordings, privacy regulations)
- `research/07-ga4-guide-2025.md` — GA4 complete guide (home dashboard, realtime, lifecycle reports, engagement reports, events, conversions, pages/screens, landing pages, explore section: funnel/path/segment/cohort explorations)
- `research/08-visitors-now.md` — Visitors.now (cookie-free, GDPR compliant, realtime visitor data, pageviews/sessions, custom events, outbound links, top pages, referrers, countries, devices)
- `research/09-peek-analytics.md` — Peek privacy-first analytics (live visitors, pages & sources, countries & devices, funnels, journeys, UTM campaigns, no cookies, under 1KB)
- `research/10-umami.md` — Umami (privacy-focused, no cookies, pageviews, unique visitors, referrers, browsers & devices, countries & regions, custom events, realtime data, multi-website support)
- `research/11-tinyvisits.md` — Tinyvisits (minimalistic, privacy-centric, cookie-free, anonymous pageview counts, GDPR/ePrivacy compliant, opt-out, no fingerprinting)
- `research/12-kobbe.md` — Kobbe (privacy-first, cookieless, traffic overview, realtime visitors, funnels, conversions, revenue attribution, 404 tracking, engagement KPIs, heatmap)

Each file should contain a structured summary of the features and design principles described by that source.

### 2. Features plan

Create `research/features-plan.md` synthesizing the research into a concrete feature set for v3. The plan must:

- Determine the page structure (which pages, what goes on each) based on research findings — do not pre-decide the number of pages
- List all proposed features grouped by page
- Justify each feature with which research sources support it
- Explicitly mark which features are **in scope** vs **out of scope** for the v3 prototype
- Follow the design principle: **keep pages minimal so users do not get overwhelmed**

### 3. Design principles for v3

Based on the research, the v3 prototype should follow these principles:

- **Privacy-first**: No cookies, no persistent visitor profiles, anonymous aggregated data
- **Minimal & calm**: Few pages, few metrics per page, no overwhelming dashboards (inspired by Peek, Umami, Fathom, tinyvisits)
- **Real-time**: Live visitor count visible on the dashboard
- **Essential metrics only**: Pageviews, unique visitors, sessions, bounce rate, avg session duration, top pages, top sources/referrers, countries, devices/browsers/OS
- **Responsive**: Works on mobile and desktop
- **Dark/light theme**: Toggle persisted in localStorage
- **Self-contained**: Static HTML + CSS + JS, no build step, CDN-loaded Bootstrap + Vue + Chart.js
- **Multi-page architecture**: Each page is a separate standalone HTML file (not a SPA). Navigation between pages is via regular `<a href>` links in the sidebar. No client-side routing, no Vue router, no shared state across pages beyond `shared.js`/`shared.css`

### 4. Page structure

The number of pages and what goes on each page will be **determined by the features plan** (step 2) based on research findings. Do not pre-decide the page count — let the research inform the structure. The guiding constraint is: **keep pages minimal so users do not get overwhelmed**.

### 5. Shared assets

- `shared.css` — Design system: CSS variables for theming, sidebar, topbar, cards, tables, toasts, responsive breakpoints
- `shared.js` — Shared Vue data (mock data), helpers, theme toggle, toast system, navigation config

### 6. Ask for approval

Before building any HTML files, present the features plan (including the proposed page structure) to the user and ask for approval. Do not proceed to step 7 until approved.

### 7. Build the prototype

After approval, create all files in `/ai/prototypes/v3/` as specified in the approved features plan:

Each page should:
- Use Vue.js 3 (CDN) for reactivity
- Use Bootstrap 5 (CDN) for layout/components
- Use Bootstrap Icons (CDN) for icons
- Use flag-icons (CDN) for country flags
- Use Chart.js (CDN) for charts
- Include the sidebar navigation shell
- Support dark/light theme toggle
- Be fully responsive
- Use mock data from `shared.js`
- Have `v-cloak` to prevent FOUC
- Show loading skeletons/spinners where appropriate

## Constraints

- **No build step** — everything is static HTML with CDN imports
- **No npm, no bundler, no TypeScript** — plain JS and CSS
- **No backend** — all data is mock data in `shared.js`
- **AGPL-3.0** — footer should reference the license
- **Minimal** — each page should fit on 1-2 screens without overwhelming the user