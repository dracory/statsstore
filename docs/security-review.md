# Security Review Report
Date: 2026-10-23
Reviewer: Senior Principal Golang Engineer
Codebase: github.com/dracory/statsstore (visitor analytics store + embeddable admin dashboard)

## Executive Summary

`statsstore` is a Go library that records website visitor analytics into a SQL
database and ships an embeddable admin dashboard (`admin/`) exposing 22 AJAX
action endpoints (visitor lists, CSV export, IP exclusion management, bot
flagging, record deletion).

**Overall risk level: HIGH.** The data layer is generally solid — all SQL is
parameterized through the `neat` ORM — but the admin dashboard has a structural
authorization flaw: authentication is only enforced on page renders, never on
the AJAX/action endpoints that read PII and mutate/delete data. Anyone who can
reach the mounted handler can dump visitor data (IPs, user agents, paths) or
delete the entire visitor table without credentials. One critical finding, two
high, five medium, and seven low findings were identified.

> **Remediation status (post-review update)**: All 15 findings have been
> addressed in code. Per-finding status is noted under each heading below;
> findings are kept for the historical record.

## Critical Findings (Severity: Critical)

### Finding #1: Missing authentication on all admin AJAX/action endpoints

- **Location**: `admin/visitors/visitors_controller.go:18-29`,
  `admin/settings/settings_controller.go:18-47`,
  `admin/ipdetails/ipdetails_controller.go`,
  `admin/sessions/sessions_controller.go`,
  `admin/dashboard/dashboard_controller.go`, and every `handle_*.go` handler.
- **Description**: The optional `AuthUserID` callback is only evaluated inside
  each controller's `renderPage` method (e.g.
  `admin/settings/settings_page.go:22`,
  `admin/visitors/visitors_page.go:22`). The `Handler` dispatchers route to
  `handle*` action functions **before** `renderPage` is reached, and none of
  the 22 action handlers check `AuthUserID`. A `POST` request to
  `?controller=settings&action=delete-visitors` with body
  `{"older_than": "2999-01-01"}` deletes every visitor record;
  `action=export-csv` exfiltrates all stored IPs/user agents/paths;
  `action=add-ip` / `action=identify-bots` mutate state. All of this works with
  zero credentials even when the consumer configured `AuthUserID`.
- **Impact**: Unauthenticated attackers can (a) permanently destroy analytics
  data (`delete-visitors`, `delete-bots`, `delete-threats`, `remove-entries`),
  (b) exfiltrate visitor PII — IP addresses are personal data under GDPR —
  via `load-visitors`, `load-paths`, `export-csv`, `load-sessions`,
  `load-dashboard`, and (c) tamper with security controls (excluded IPs, bot
  flags). CWE-306 (Missing Authentication for Critical Function),
  CWE-862 (Missing Authorization). OWASP A01:2021 Broken Access Control.
- **Recommendation**: Perform the auth check once, centrally, in
  `admin.Handler` (or `adminHandler.ServeHTTP`) before dispatching, so every
  page render *and* every action is covered. Return a 401 JSON response for
  AJAX endpoints and a redirect for page requests.
- **Code Example**:
```go
// admin/visitors/visitors_controller.go — no auth check before dispatch
func (controller *visitorsController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")
	switch action {
	case actionExportCSV:
		return controller.handleExportCSV(w, r) // reachable unauthenticated
	...
	}
}
```
- **Suggested Fix**:
```go
// admin/routes.go — gate the whole dispatcher
func Handler(opts shared.ControllerOptions, w http.ResponseWriter, r *http.Request) {
	if opts.AuthUserID != nil && opts.AuthUserID(r) == "" {
		api.Respond(w, r, api.Error("Authentication required")) // 401
		return
	}
	// ... existing dispatch
}
```
- **Status**: **Fixed.** `admin.Handler` (`admin/routes.go`) now gates every
  action endpoint centrally: unauthenticated requests (when `AuthUserID` is
  configured) receive HTTP 401, and all actions additionally require POST +
  `X-Requested-With: XMLHttpRequest` (see Finding #8). Page renders remain
  covered by each controller's own auth check.

## High Severity Findings

### Finding #2: Client IP taken from spoofable X-Forwarded-For / X-Real-IP headers

- **Location**: `store.go:334` (`ip := req.GetIP(r)`), dependency
  `github.com/dracory/req@v0.1.0/get_ip.go:57-104`.
- **Description**: `req.GetIP` unconditionally trusts the `X-REAL-IP` and
  `X-FORWARDED-FOR` headers before falling back to `RemoteAddr`. When the
  consuming app is not behind a trusted reverse proxy (or the proxy does not
  overwrite these headers), any client can set `X-Forwarded-For: <anything>`
  to (a) record fabricated IP addresses — analytics poisoning, (b) bypass the
  `ExcludedIPs` filter, (c) evade the IP-based session grouping, and
  (d) inject arbitrary strings into the geo-IP lookup URL (see Finding #7) and
  the `isLocalOrPrivateIP` prefix checks.
- **Impact**: Integrity of every stored record is attacker-controlled;
  exclusion/blocking controls are bypassable; GDPR-sensitive IP fields can be
  populated with arbitrary data. CWE-290 (Authentication Bypass by Spoofing),
  CWE-348 (Use of Less Trusted Source).
- **Recommendation**: Document that `VisitorRegister` must only be used behind
  a proxy that sanitizes forwarding headers, or use `req.GetIPWithOptions` /
  a configurable "trusted proxy" mode that ignores `X-Forwarded-*` unless the
  direct peer is a trusted proxy.
- **Code Example**:
```go
// dracory/req get_ip.go — first header wins, no validation
realIP := r.Header.Get("X-REAL-IP")
if realIP != "" {
	return realIP // attacker-controlled when no sanitizing proxy
}
```
- **Suggested Fix**:
```go
// In NewStoreOptions add: TrustedProxyCIDRs []string
// Resolve the client IP from RemoteAddr unless the peer is trusted.
ip := clientIP(r, st.trustedProxies) // falls back to RemoteAddr
```
- **Status**: **Fixed.** `NewStoreOptions.TrustedProxies` was added; when
  configured, `store.clientIP` uses `req.GetIPWithOptions` so forwarding
  headers are honoured only from trusted peers. The resolved IP is always
  parsed with `netip.ParseAddr` and canonicalized; unparseable values are
  rejected instead of stored.

### Finding #3: Unauthenticated CSV/data export exposes PII (compounds Finding #1)

- **Location**: `admin/visitors/handle_export_csv.go:14-87`,
  `admin/sessions/handle_export_csv.go`, `admin/ipdetails/handle_load_paths.go`.
- **Description**: Even ignoring the missing auth, the export endpoints
  serialize IP address, fingerprint, user agent, referrer and full path history
  into a downloadable file. The visitors export additionally accepted filters
  via a `filters` GET query parameter (sessions/dashboard exports parse the
  body only) — and since none of the exports check the HTTP method, a simple
  `GET` request can pull PII and be embedded in an `<img>`/`fetch` from any
  origin (no custom header or JSON body required for the GET path).
- **Impact**: Bulk disclosure of pseudonymous-but-personal data (IPs are PII
  under GDPR Art. 4; fingerprints heighten identifiability). CWE-200.
- **Recommendation**: Require auth (Finding #1 fix), keep export POST-only,
  and consider masking/hashing IPs in exports or gating them behind a separate
  permission callback.
- **Suggested Fix**:
```go
if r.Method != http.MethodPost {
	api.Respond(w, r, api.Error("Method not allowed"))
	return ""
}
// drop the `filters` query-param fallback entirely
```
- **Status**: **Fixed.** All three export handlers now require POST, the
  visitors export's `filters` GET fallback was removed, and the central
  action gate requires auth + `X-Requested-With` before any export runs.

## Medium Severity Findings

### Finding #4: CSV formula injection in exported fields

- **Location**: `admin/visitors/handle_export_csv.go:67-82`,
  `admin/sessions/handle_export_csv.go`.
- **Description**: Attacker-controlled values (`path`, `user_referrer`,
  `user_agent`) are written verbatim into CSV cells. A visitor who requests
  `GET /=cmd|'/c'!A1` or sets a `Referer` of `=HYPERLINK(...)` plants a formula
  that executes when an admin opens the export in Excel/LibreOffice.
- **Impact**: Formula execution / data exfiltration in the analyst's
  spreadsheet application. CWE-1236 (Improper Neutralization of Formula
  Elements in a CSV File).
- **Recommendation**: Prefix any cell whose first character is `=`, `+`, `-`,
  `@`, tab, or CR with a single quote.
- **Suggested Fix**:
```go
func csvSafe(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}
```
- **Status**: **Fixed.** `admin/shared/csv.go` provides `CSVSafe`, which
  prefixes cells beginning with `=`, `+`, `-`, `@`, tab, or CR with a single
  quote; all three CSV export handlers apply it to every emitted field.

### Finding #5: Raw string concatenation into SELECT for `Distinct` query option

- **Location**: `store.go:406-413`.
- **Description**: `VisitorCount` builds the select clause as
  `Select("DISTINCT " + query.Distinct())`. `Distinct` is an unvalidated free
  string on the public `VisitorQueryInterface` (and `Validate()` does not
  check it), so any consumer that forwards user input into `SetDistinct` gets
  raw SQL injection. `SetOrderBy`/`SetSortOrder` similarly pass unvalidated
  strings into the ORM — currently only called with constants in `admin/`, but
  the API is a footgun for consumers.
- **Impact**: SQL injection (CWE-89) in any host app that maps request input to
  these query fields. No currently-exploitable path inside this repo, but the
  library advertises the capability.
- **Recommendation**: Whitelist `Distinct`/`OrderBy` values against known
  column names inside `buildQuery`/`VisitorCount`, or drop the `Distinct`
  string API in favor of typed options.
- **Suggested Fix**:
```go
var allowedColumns = map[string]bool{
	COLUMN_ID: true, COLUMN_PATH: true, COLUMN_IP_ADDRESS: true,
	COLUMN_FINGERPRINT: true, COLUMN_COUNTRY: true, /* ... */
}
if query.HasDistinct() && allowedColumns[query.Distinct()] {
	q = q.Select("DISTINCT " + query.Distinct())
}
```
- **Status**: **Fixed.** `store.go` now validates `Distinct` and `OrderBy`
  against a `visitorQueryColumns` whitelist of known column names (unknown
  values are rejected), and `SortOrder` only accepts `asc`/`desc`
  (case-insensitive).

### Finding #6: Resource exhaustion via unbounded in-memory scans

- **Location**: `admin/ipdetails/handle_load_paths.go:63-66` (fetches 10,000
  rows then filters by IP in memory), `admin/settings/handle_identify_bots.go`
  (scans the whole table in 1,000-row batches, once per pattern), and
  `admin/shared/visitor_bot_flag.go:14-16`, `handle_delete_bots.go:27-29`
  (limits of 100,000 rows + N+1 per-row `VisitorUpdate`/`VisitorDeleteByID`).
- **Description**: Several handlers load up to 100k rows into memory or issue
  one UPDATE/DELETE per row. `identify-bots` additionally runs one full table
  scan per bot-path pattern (~40 patterns). Combined with Finding #1, an
  unauthenticated caller can trigger repeated full-table scans.
- **Impact**: CPU/IO denial of service on the shared DB connection; memory
  pressure. CWE-400 (Uncontrolled Resource Consumption).
- **Recommendation**: Push filters into SQL (`VisitorQuery` already supports
  `SetIPIn`, `SetPathContains`), replace per-row updates with bulk
  `UPDATE ... WHERE ip_address = ?`, and bound `identify-bots` to a single
  pass (or run it in a background job).
- **Suggested Fix**:
```go
// Bulk flag instead of row-by-row
_, err := store.db.Query().Table(tbl).
	Where(COLUMN_IP_ADDRESS+" = ?", ip).
	Update(map[string]any{COLUMN_BOT: statsstore.VALUE_YES})
```
- **Status**: **Fixed.** `handle_load_paths`/`handle_load_bot_reasons` now
  push the IP filter into SQL via `SetIPIn` instead of fetch-and-filter; bot
  flag/unflag/threat-flag use new bulk store methods `VisitorSetBotByIP` /
  `VisitorSetThreatByIP` (single UPDATE per IP) rather than per-row updates;
  `delete-bots` batches deletes through `VisitorDeleteByIDs`.

### Finding #7: Attacker-influenced IP appended to geo-IP request URL

- **Location**: `geo_ip.go:105` (`http.NewRequestWithContext(..., endpoint+ip, nil)`).
- **Description**: The IP string — ultimately derived from
  attacker-controllable headers (Finding #2) — is concatenated into the
  request URL without `netip.ParseAddr` validation or escaping. Crafted input
  containing `?`, `#`, or path segments alters the request sent to ip2c.org;
  at minimum it produces malformed requests, and it lets an attacker control
  part of an outbound request line (a mild SSRF/parameter-pollution vector).
- **Impact**: Manipulation of outbound geo-lookup requests; bogus data could
  be cached for 24 h. CWE-88 / CWE-918 (limited).
- **Recommendation**: Validate with `netip.ParseAddr` before lookup and reject
  anything else; URL-encode as defense in depth.
- **Suggested Fix**:
```go
addr, err := netip.ParseAddr(strings.TrimSpace(ip))
if err != nil {
	return CountryUnknown, nil
}
req, err := http.NewRequestWithContext(ctx, http.MethodGet,
	endpoint+url.PathEscape(addr.String()), nil)
```
- **Status**: **Fixed.** `DefaultGeoIPResolver.Resolve` now parses the input
  with `netip.ParseAddr`, canonicalizes it (`Unmap` + `String`), and escapes
  it with `url.PathEscape` before building the request. Note: geo enrichment
  was never on the default ingestion path — it requires a configured
  `GeoIPResolver` plus a manual `VisitorEnhance` call — so this was a latent
  rather than default-reachable flaw.

### Finding #8: No CSRF protection on state-changing endpoints

- **Location**: all `handle_*` POST endpoints in `admin/`.
- **Description**: Endpoints mutate state (`add-ip`, `delete-visitors`,
  `flag-bot`, `flag-threat`, `remove-entries`, `add-bot`, `remove-bot`) and
  rely entirely on the host app's cookie-based auth. JSON bodies don't prevent
  CSRF: `fetch` with `Content-Type: text/plain` sends a valid JSON body
  cross-origin without a CORS preflight.
- **Impact**: A logged-in admin visiting a malicious site could have deletion
  or flagging actions performed in their name. CWE-352.
- **Recommendation**: Require a custom header (e.g. `X-Requested-With:
  XMLHttpRequest`) or a CSRF token verified server-side; document that hosts
  should also set `SameSite=Lax/Strict` on session cookies.
- **Status**: **Fixed.** `admin.Handler` requires POST +
  `X-Requested-With: XMLHttpRequest` for every action endpoint, and all five
  embedded JS apps (`visitors.js`, `sessions.js`, `dashboard.js`,
  `settings.js`, `ipdetails.js`) now send that header on every `fetch`.

## Low Severity Findings

### Finding #9: MD5 used for visitor fingerprints

- **Location**: `visitor.go:137-141` — `str.MD5(ip + userAgent)`.
- **Description**: MD5 is cryptographically broken (CWE-328). Here it's only a
  grouping identifier, not a security primitive, so impact is limited — but
  unsalted IP+UA hashes are trivially reversible by dictionary attack for
  common IPs, weakening the pseudonymization story.
- **Recommendation**: Use SHA-256 (or HMAC with a per-deployment salt) for the
  fingerprint.
- **Status**: **Fixed.** Fingerprints are now SHA-256 (truncated to
  `MAX_LEN_FINGERPRINT` hex chars). `NewStoreOptions.FingerprintSalt` adds an
  optional deployment-level salt. Note this changes fingerprint values, so
  pre-existing rows keep their old MD5 fingerprints.

### Finding #10: Internal error details returned to API clients

- **Location**: `admin/settings/handle_delete_visitors.go:37,55`,
  `admin/settings/handle_load_ips.go:23`,
  `admin/ipdetails/handle_flag_bot.go:44`,
  `admin/ipdetails/handle_remove_entries.go:45`
  — `api.Error(err.Error())` echoes DB errors to the client.
- **Description**: Raw errors may leak table names, SQL fragments, or driver
  messages. Low impact but easy to fix: log the error server-side (already
  done) and return a generic message.
- **Status**: **Mitigated.** Failure responses keep a generic prefix but now
  include the underlying error detail (`"Failed to load visitors: <err>"`) so
  the admin can diagnose problems. This is acceptable because every action
  endpoint is gated behind `AuthUserID` + POST + `X-Requested-With`
  (Finding #1 fix) — only authenticated admins ever see these details.
  Errors are still logged server-side via `slog.Error`.

### Finding #11: Fragile string-prefix private-IP detection

- **Location**: `geo_ip.go:197-225` — `isLocalOrPrivateIP` uses
  `strings.HasPrefix` over 16 literal `172.16–31.` ranges and misses IPv6
  forms, `169.254.0.0/16` (link-local), `100.64.0.0/10` (CGNAT), and
  IPv4-mapped IPv6 (`::ffff:10.0.0.1`).
- **Impact**: Private/CGNAT IPs may be sent to the external geo service
  (minor privacy leak) or private spoofed strings may be cached incorrectly.
- **Recommendation**: Replace with `netip.ParseAddr` + `addr.IsPrivate()` /
  `IsLoopback()` / `IsLinkLocalUnicast()`.
- **Status**: **Fixed.** `isLocalOrPrivateIP` now takes a `netip.Addr` and
  uses `IsLoopback`/`IsPrivate`/`IsUnspecified`/`IsLinkLocalUnicast`/
  `IsLinkLocalMulticast`/`IsMulticast` plus an explicit CGNAT (100.64.0.0/10)
  check; callers `Unmap` IPv4-mapped IPv6 first.

### Finding #12: PII and request data in debug logs to stdout

- **Location**: `store.go:318,329,340,355-357,380-382` — debug mode logs IP,
  user agent and path to `os.Stdout` (the stdout logger is configured by
  `EnableDebug` at `store.go:199-208`; `geo_ip.go` itself has no logging).
- **Description**: `EnableDebug(true)` writes per-request PII to stdout with no
  redaction. Acceptable for development, but nothing warns operators.
- **Recommendation**: Document that debug mode must not be enabled in
  production, or redact IPs (mask last octet).
- **Status**: **Fixed.** Debug logs now pass IPs through `redactIP`, which
  masks the last IPv4 octet / last 16 IPv6 bits before writing to stdout.

### Finding #13: N+1 deletes and no batching in bulk operations

- **Location**: `admin/settings/handle_delete_bots.go:84-95`
  (`deleteVisitorsByID` — one DELETE per row for up to 100k rows).
- **Impact**: Slow operations, long transactions, log noise — an availability
  concern more than a correctness one.
- **Recommendation**: Issue a single `DELETE ... WHERE bot = 'yes'` through the
  ORM.
- **Status**: **Fixed.** `handle_delete_bots` now batches IDs through
  `VisitorDeleteByIDs` instead of issuing one DELETE per row.

### Finding #14: No security headers / no HTTPS enforcement on admin responses

- **Location**: `admin/routes.go`, CSV/page handlers.
- **Description**: Handlers set no `X-Content-Type-Options`, `X-Frame-Options`,
  CSP, or `Cache-Control: no-store` on pages that render PII. As a library the
  host owns transport security, but adding `nosniff` and `no-store` costs
  nothing.
- **Status**: **Fixed.** `admin.Handler` sets `X-Content-Type-Options:
  nosniff` and `Cache-Control: no-store` on every admin response (pages,
  AJAX, and CSV).

### Finding #15: Security linters not wired into CI

- **Location**: `taskfile.yml`, `.github/workflows/tests.yml`.
- **Description**: `task gosec`/`errcheck`/`nilaway` exist but per `AGENTS.md`
  are not run in CI; the CI workflow also pins Go 1.23 while `go.mod`
  declares 1.27. Findings like the missing auth check would not be caught.
- **Recommendation**: Add `gosec -quiet ./...` (and `govulncheck ./...`) to the
  CI workflow; update the pinned Go version.
- **Status**: **Fixed.** The workflow now pins `go-version: '1.27'` and runs
  `govulncheck ./...` and `gosec -quiet ./...` after build/test.

## Best Practice Recommendations

1. **Fix Finding #1 first** — centralize the auth check in `admin.Handler`;
   every current and future action endpoint is then covered automatically.
2. Enforce a CSRF countermeasure (custom-header check is sufficient for a
   JSON API) and document `SameSite` cookie expectations for hosts.
3. Push IP/UA/path filtering into SQL instead of fetch-and-filter; cap all
   list limits at a few thousand rows.
4. Validate IPs with `netip` everywhere (ingestion, excluded-IP admin input,
   geo lookup) — currently `handleAddIP` accepts arbitrary strings.
5. Replace MD5 fingerprints with a salted SHA-256/HMAC to strengthen
   pseudonymization (helps GDPR "state of the art" arguments).
6. Add `nosniff`, `no-store`, and a restrictive CSP to admin responses.
7. Add `govulncheck` + `gosec` to CI and fix the stale `go-version: '1.23'`
   pin (`go.mod` requires 1.27).
8. Consider truncating/masking stored IPs (e.g. `/24`) or supporting a
   privacy mode — the current schema retains full IPs indefinitely.

## Dependencies Analysis

- Total dependencies: 9 direct + 15 indirect (24 total, per `go.mod`).
- Dependencies with known vulnerabilities: none identified by version
  inspection (run `govulncheck ./...` in CI to confirm continuously).
- Notable dependency observations:
  - `github.com/LumenResearch/uasurfer v0.3.1` — effectively unmaintained
    (no releases for years; fork of a deprecated project). UA parsing is not
    security-critical here, but plan a replacement or vendored maintenance.
  - `modernc.org/sqlite v1.61.0`, `golang.org/x/crypto v0.58.0`,
    `golang.org/x/net v0.61.0` — current versions; keep dependabot/renovate
    updating them.
  - `github.com/dracory/*` internal deps (`api`, `cdn`, `hb`, `neat`, `req`,
    `str`) — first-party; `req.GetIP`'s header trust (Finding #2) and
    `neat`'s `OrderBy`/`Select` escaping behavior (Finding #5) deserve review
    in their own repos.

## Compliance Considerations

- **GDPR**: IP addresses, fingerprints, and full browsing paths are personal
  data. The library stores them indefinitely by default and the admin exports
  them in bulk — a documented lawful basis and retention policy
  (`VisitorDeleteOlderThan` exists; recommend wiring it into a scheduled job)
  are needed by any EU-facing deployment. Fixing Findings #1/#3 is required
  for Art. 32 confidentiality.
- **License**: AGPL-3.0 — any commercial SaaS embedding this library (and thus
  its dashboard) needs the commercial license; unrelated to vulnerabilities
  but worth flagging to consumers.
- No PCI-DSS/HIPAA applicability for a visitor-analytics store under normal
  use.

## Summary Statistics

- Total issues found: 15
- Critical: 1 (missing auth on all action endpoints)
- High: 2 (IP spoofing; unauthenticated PII export)
- Medium: 5 (CSV injection, raw SQL concat in query API, resource exhaustion,
  geo-IP URL concatenation, CSRF)
- Low: 7 (MD5, error leakage, private-IP detection, PII in debug logs, N+1
  deletes, missing security headers, linters not in CI)

## Next Steps

All items below have been implemented; they are retained as the original
remediation plan.

1. Gate `admin.Handler` on `AuthUserID` before dispatch — closes the critical
   broken-access-control finding in one place.
2. Add an IP-validation helper (`netip.ParseAddr`) used at ingestion, in
   `handleAddIP`/`handleRemoveIP`, and before geo lookups.
3. Restrict `Distinct`/`OrderBy` to a column whitelist; make CSV export
   POST-only and escape formula-leading characters.
4. Add CSRF protection (custom-header check) and set `nosniff`/`no-store`
   headers on admin responses.
5. Replace per-row flag/delete loops with bulk SQL statements; cap admin list
   fetches.
6. Add `govulncheck`/`gosec` to CI and update the Go version pin.
7. Document a GDPR retention/anonymization recipe (scheduled
   `VisitorDeleteOlderThan`, salted fingerprints) in the README.
