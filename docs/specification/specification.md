# statsstore — Product Specification

**License:** AGPL-3.0 (commercial license available separately)

---

## 1. What It Is

statsstore is a **privacy-first, cookieless website visitor analytics** library for Go. It lets you track, store, and analyze website traffic without cookies, without third-party scripts, and without sending visitor data to external services.

You embed it into your existing Go application — it runs on your infrastructure, stores data in your database, and serves a dashboard from your own domain.

### Who It's For

- Go developers who want built-in analytics without adding a SaaS dependency
- Self-hosted projects that need a privacy-respecting analytics dashboard
- Applications that already use `database/sql` and want to add visitor tracking

### Key Differentiators

- **Cookieless** — no cookies, no consent banners needed (GDPR-friendly)
- **Self-hosted** — your database, your server, your data
- **Embeddable** — it's a Go library, not a separate service to deploy
- **Fingerprint-based** — identifies unique visitors via IP + user-agent hash, enabling cross-day visitor linking
- **Bot-filtering** — built-in filtering for crawlers, referrer spam, and data-center traffic

---

## 2. Visitor Tracking

### 2.1 What Gets Tracked

Every page view captures:

- **Page path** — which URL the visitor accessed
- **IP address** — used for geo-location and unique visitor identification
- **Fingerprint** — MD5 hash of IP + user-agent, for unique visitor counting
- **User agent** — raw browser string
- **Browser** — name and version (Chrome, Firefox, Safari, Edge, etc.)
- **Operating system** — name and version (Windows, macOS, Android, iOS, etc.)
- **Device** — device name (iPhone, iPad, etc.) and type (desktop, mobile, tablet)
- **Referrer** — where the visitor came from
- **Language** — Accept-Language header
- **Country** — ISO 3166-1 alpha-2 code (populated via geo-IP enrichment)
- **Timestamp** — when the visit occurred

### 2.2 How Tracking Works

Tracking happens server-side. When a visitor hits your site, your application calls `VisitorRegister(request)` and statsstore:

1. Extracts the path, IP, user-agent, and referrer from the HTTP request
2. Parses the user-agent into browser, OS, and device fields
3. Creates a visitor record in the database

No JavaScript snippet is required for basic page-view tracking. The tracking script can be a simple server-side middleware or endpoint.

### 2.3 What Gets Filtered Out

When bot filtering is enabled, the following traffic is silently excluded:

- **Known bots and crawlers** — 65+ patterns including Googlebot, Bingbot, Ahrefs, SEMrush, curl, wget, Python-requests, headless browsers, Selenium, Puppeteer, etc.
- **Referrer spam** — 55+ known spam domains (semalt.com, darodar.com, etc.)
- **Data-center traffic** — requests from AWS, GCP, Azure, DigitalOcean, and Oracle Cloud IP ranges

### 2.4 Exclusions

You can configure:

- **Excluded path prefixes** — e.g. skip all `/admin/` traffic
- **Excluded IP addresses** — skip specific IPs (managed via dashboard settings page or programmatically)

---

## 3. Unique Visitor Identification

statsstore uses a **fingerprint** — an MD5 hash of the visitor's IP address + user-agent string — to identify unique visitors. This enables:

- **Unique visitor counts** — distinct fingerprints in a given period
- **First vs. return visits** — new fingerprints vs. previously seen ones
- **Cross-day linking** — the same visitor across multiple days (unlike daily-reset counters)

No cookies are used. The fingerprint is derived from data already present in the HTTP request.

---

## 4. Geo-IP Enrichment

Visitor records are initially saved without a country. Country is populated asynchronously via a geo-IP lookup.

### How It Works

1. A background task (cron job or goroutine) calls `VisitorEnhance` on a schedule
2. The system fetches visitor records with empty country fields (batch size configurable)
3. Each unique IP is resolved to a country code via a geo-IP service
4. All records sharing that IP are bulk-updated with the country
5. User-agent fields are also filled in if missing (browser, OS, device)

### Default Geo-IP Service

The built-in resolver uses [ip2c.org](https://ip2c.org), a free IP-to-country service:

- 24-hour in-memory cache (avoids duplicate lookups)
- 2-second rate limiting between API calls
- Localhost and private IPs return "Unknown" without an API call
- 5-second timeout per request

### Custom Geo-IP Providers

You can plug in any geo-IP provider (MaxMind, ipinfo, etc.) by implementing the `GeoIPResolver` interface. This is recommended for production traffic.

---

## 5. Dashboard

The admin dashboard is a **Vue.js single-page application** embedded directly in the Go binary — no separate frontend build step, no npm dependencies, no static file server needed.

### 5.1 Dashboard Home

The main analytics overview page with:

- **Stat cards** — total visits, unique visitors, first visits, return visits for the selected period
- **Period-over-period comparison** — current period vs. previous period with percentage change
- **Daily stats table** — day-by-day breakdown of total, unique, first, and return visits
- **Live visitors** — visitors active in the last 5 minutes
- **Traffic source cards** with tabbed views:
  - **Referrers** — top referring sites
  - **Pages** — most visited pages, entry pages, exit pages
  - **Browsers** — browser breakdown, device types, operating systems
  - **Countries** — geographic distribution, languages
- **Period selector** — today, yesterday, last 7 days, this week, last week, this month, last month
- **CSV export** — download the current dataset

### 5.2 Visitor Activity

A detailed per-visitor listing showing:

- Each visitor's path, location, device, browser, OS, referrer, and timestamp
- **Filters** — by time range, country, device type
- **Pagination** — configurable rows per page (10/25/50/100)
- **Detail modal** — click any visitor for full record details
- **CSV export** — download filtered results

### 5.3 Visitor Paths

Path-level aggregation showing:

- Unique paths visited, with visit counts and session counts
- Drill-down to visitor activity filtered by a specific path
- Filters and pagination
- CSV export

### 5.4 Page View Activity

Page-view-level listing with:

- Individual page view records with device/browser/OS labels, location, country
- Filters and pagination
- CSV export

### 5.5 Settings

Administrative controls:

- **Excluded IPs** — add or remove IP addresses from the exclusion list (persisted to database)
- **Delete visitors** — remove all visitor data by IP address, or delete all visitor records

---

## 6. Data Privacy

### 6.1 No Cookies

statsstore does not set or read any cookies. No cookie consent banner is required under GDPR/ePrivacy for the tracking itself (consult your legal advisor for your specific jurisdiction).

### 6.2 Data Ownership

All visitor data is stored in your database on your infrastructure. No data is sent to third-party services except the optional geo-IP lookup (which sends only the IP address to ip2c.org or your configured provider).

### 6.3 Data Retention

You control retention:

- **Soft delete** — mark records as deleted without removing them (recoverable)
- **Hard delete** — permanently remove individual records, all records by IP, or all records
- **Programmatic** — full CRUD access via the store API for custom retention policies

### 6.4 Fingerprint Privacy Trade-off

Unlike some privacy-first analytics tools (e.g. Plausible) that store no persistent identifiers and reset unique visitor counts daily, statsstore uses fingerprint-based tracking with IP storage. This enables cross-day visitor linking but is less privacy-strict. This is a deliberate design trade-off.

---

## 7. Data Export

- **CSV export** from the dashboard for visitor activity, visitor paths, and page view activity
- **Programmatic access** via the store API for custom export logic
- **Raw SQL access** via `GetDB()` for advanced queries and reporting

---

## 8. Integration Model

### How You Embed It

1. **Create a store** — pass your `*sql.DB` and table name; statsstore auto-migrates the schema
2. **Track visitors** — call `VisitorRegister(request)` from your HTTP handlers or middleware
3. **Mount the dashboard** — call `admin.New(options)` to get an `http.Handler`, mount it at any route
4. **Enrich data** — call `VisitorEnhance()` from a background ticker/cron on your preferred schedule

### Database Compatibility

Works with any SQL database supported by Go's `database/sql`:
- SQLite (pure-Go driver, no CGO)
- PostgreSQL
- MySQL / MariaDB
- SQL Server
- Any other driver that implements `database/sql`

### Layout Customization

The dashboard renders through a `LayoutInterface` that you implement. This gives you full control over:
- HTML page structure and `<head>` contents
- CSS frameworks (the demo uses Bootstrap 5.3)
- JavaScript libraries (the demo uses Vue.js 3.5 and Chart.js)
- Country name resolution (provide your own ISO2 → country name function)
- Branding and navigation

---

## 9. Demo

A runnable demo app is included at `examples/admin-demo/`. It provides:

- In-memory SQLite database (no setup needed)
- ~100 seeded demo visitor records across 7 days
- Full dashboard with all pages functional
- Reference layout implementation using Bootstrap, Vue.js, Chart.js, and flag-icons

Run it with `go run ./examples/admin-demo` and open `http://localhost:8080/admin/home`.

---

## 10. What's Not Included (By Design)

- **Session replay** — no recording of visitor screen activity
- **Web vitals** — no Core Web Vitals tracking
- **Error tracking** — no JS error collection
- **User profiles** — no persistent user identity beyond fingerprint
- **Heatmaps** — no click/scroll heatmaps
- **A/B testing** — no experiment framework
- **Email reports** — no scheduled report delivery (planned)
- **Public dashboards** — no shareable dashboard links (planned)
- **Multi-site management** — single site per store instance
- **Authentication** — no built-in auth; protect the dashboard route yourself

---

## 11. Roadmap (Proposed)

Based on competitive analysis against Plausible and other analytics tools:

### Near-term
- Bounce rate calculation (from existing data)
- Realtime/live visitor improvements
- UTM campaign tracking
- Visit duration / time on page
- Custom events

### Medium-term
- Goals and conversion tracking
- Funnels
- Outbound link click tracking
- File download tracking
- Scroll depth tracking
- 404 page tracking

### Long-term
- REST API layer (Stats API, Events API)
- Email/Slack reports
- Multi-site consolidated view
- Public/shareable dashboards

See `docs/proposals/todo/plausible.md` for the full feature gap analysis.

---

## 12. License

- **Open source:** AGPL-3.0 — free for self-hosted use, modifications must be shared under the same license
- **Commercial license:** available for use in proprietary/closed-source applications (contact via [lesichkov.co.uk/contact](https://lesichkov.co.uk/contact))
