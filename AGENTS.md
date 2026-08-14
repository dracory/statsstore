# AGENTS.md

Compact guidance for working in the `statsstore` repo (module `github.com/dracory/statsstore`).

Long-term, repo-specific context lives in [`ai/lessons-learnt.md`](ai/lessons-learnt.md) — consult and update it as you discover non-obvious truths.

## What this repo is
A Go library plus an embeddable admin dashboard for website visitor analytics.
It is a **library**, not a standalone service: consumers embed `admin.New(...)` as an `http.Handler` and bring their own `database/sql` DB.

## Developer commands
- Run all tests: `go test ./...` (also `task test`).
- Run a single test / package: `go test -run <TestName> ./admin/dashboard` etc.
- Coverage report (opens HTML): `task cover` (runs `go test ./... -coverprofile=coverage.out`).
- Linters (not wired into CI, must `go install` first): `task errcheck`, `task nilaway`, `task gocritic`, `task golangci-lint`, `task gosec`. See `taskfile.yml` for install tasks.
- Runnable demo: `go run ./examples/admin-demo` (`admin-demo.exe` is committed but gitignored as `*.exe`; rebuild from that dir).

## Architecture / package boundaries
- Root package `statsstore` — core library: `Store` (`store*.go`), `Visitor` (`visitor*.go`), geo-IP enrichment (`geo_ip.go`), bot filtering (`bot_filter.go`). Entrypoint: `NewStore(NewStoreOptions{...})`.
- `admin/` — framework-agnostic, self-contained admin dashboard. `admin.New(admin.Options{...})` returns an `http.Handler`. Consumers inject a `shared.LayoutInterface`, a `statsstore.StoreInterface`, and optional callbacks (`AuthUserID`, `FlashError`, `CountryNameByIso2`). The package depends only on `github.com/dracory/statsstore` — no host-application knowledge. The dashboard UI is a **Vue.js SPA** embedded per controller (`.html` + `.js` via `go:embed`); there is NO separate JS/frontend build step.
- `admin/dashboard/` — overview dashboard with total visitors, unique IPs, top-N breakdowns (paths, countries, browsers, OS, device types), recent visitors. AJAX endpoint: `action=load-dashboard`.
- `admin/visitors/` — paginated visitor list with filtering (IP, country, browser, OS, device type, path, date range, bot/threat flags). AJAX endpoint: `action=load-visitors`.
- `admin/sessions/` — visitor sessions grouped by IP with visit history. AJAX endpoint: `action=load-sessions`.
- `admin/settings/` — excluded IP management, bot IP management, automated bot identification scan. AJAX endpoints: `load-ips`, `add-ip`, `remove-ip`, `delete-visitors`, `load-bots`, `add-bot`, `remove-bot`, `identify-bots`, `delete-bots`, `delete-threats`.
- `admin/ipdetails/` — per-IP deep dive with path history, flag as bot/threat, delete entries. AJAX endpoints: `load-paths`, `flag-bot`, `flag-threat`, `remove-entries`.
- `admin/shared/` — shared types (`LayoutInterface`, `ControllerOptions`, `Breadcrumb`), URL builder (`Links`), period helpers, country name resolver, bot IP store helpers, breadcrumbs component, flash-or-redirect helper, test helpers.
- `examples/admin-demo/` — a self-contained demo app wiring the store + admin together; use it to exercise the library end-to-end.
- `docs/` — design notes and proposals (some are aspirational; `docs/overview.md` references goqu/`sb` builders that are NOT in the current `go.mod`, so trust code over those docs).

## Important quirks & gotchas
- **Pure-Go SQLite, no CGO.** Tests use `modernc.org/sqlite` with a `:memory:` DSN (`store_test.go`); no external DB or CGO needed to build/test. Don't add `mattn/go-sqlite3` or other cgo drivers.
- **Geo-IP enrichment is manual.** `VisitorEnhance` is NOT called automatically — it must be invoked from a cron/goroutine, and only updates records with an empty `country`. The default `DefaultGeoIPResolver` calls the free ip2c.org service with a built-in 2s delay between calls (rate limiting) and a 24h in-memory cache. Provide a custom `GeoIPResolver` for production.
- **Toolchain mismatch:** `go.mod` declares `go 1.26.3`, but CI (`.github/workflows/tests.yml`) pins `go-version: '1.23'`. Use the local toolchain per `go.mod`; the CI pin is stale.
- `VisitorRegister(r *http.Request)` infers path/IP/user-agent from the request — the normal ingestion path. `VisitorQueryOptions` drives filtering/pagination for `VisitorList`/`VisitorCount`.

## License constraint
Licensed **AGPL-3.0**; commercial use requires a separate commercial license (see README). Keep license headers and don't advise folding this code into proprietary/closed distributions without flagging the AGPL implication.
