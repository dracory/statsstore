# AGENTS.md

Compact guidance for working in the `statsstore` repo (module `github.com/dracory/statsstore`).

Long-term, repo-specific context lives in [`ai/lessons-learnt.md`](ai/lessons-learnt.md) — consult and update it as you discover non-obvious truths.

## What this repo is
A Go library for website visitor analytics.
It is a **library**, not a standalone service: consumers bring their own `database/sql` DB.
The admin dashboard was extracted to the separate `github.com/dracory/statsadmin` module — this repo must NOT depend on it (dependency is one-way: statsadmin → statsstore).

## Developer commands
- Run all tests: `go test ./...` (also `task test`).
- Run a single test: `go test -run <TestName> .`
- Coverage report (opens HTML): `task cover` (runs `go test ./... -coverprofile=coverage.out`).
- Linters (not wired into CI, must `go install` first): `task errcheck`, `task nilaway`, `task gocritic`, `task golangci-lint`, `task gosec`. See `taskfile.yml` for install tasks.

## Architecture / package boundaries
- Root package `statsstore` — core library: `Store` (`store*.go`), `Visitor` (`visitor*.go`), geo-IP enrichment (`geo_ip.go`), bot filtering (`bot_filter.go`). Entrypoint: `NewStore(NewStoreOptions{...})`.
- The embeddable admin dashboard lives in `github.com/dracory/statsadmin` (`statsadmin.New(statsadmin.AdminOptions{...})` returns an `http.Handler`/`AdminInterface`; `FuncLayout` optional with a built-in default).
- `docs/` — design notes and proposals (some are aspirational; `docs/overview.md` references goqu/`sb` builders that are NOT in the current `go.mod`, so trust code over those docs).

## Important quirks & gotchas
- **Pure-Go SQLite, no CGO.** Tests use `modernc.org/sqlite` with a `:memory:` DSN (`store_test.go`); no external DB or CGO needed to build/test. Don't add `mattn/go-sqlite3` or other cgo drivers.
- **Geo-IP enrichment is manual.** `VisitorEnhance` is NOT called automatically — it must be invoked from a cron/goroutine, and only updates records with an empty `country`. The default `DefaultGeoIPResolver` calls the free ip2c.org service with a built-in 2s delay between calls (rate limiting) and a 24h in-memory cache. Provide a custom `GeoIPResolver` for production.
- **Toolchain mismatch:** `go.mod` declares `go 1.26.3`, but CI (`.github/workflows/tests.yml`) pins `go-version: '1.23'`. Use the local toolchain per `go.mod`; the CI pin is stale.
- `VisitorRegister(r *http.Request)` infers path/IP/user-agent from the request — the normal ingestion path. `VisitorQueryOptions` drives filtering/pagination for `VisitorList`/`VisitorCount`.

## License constraint
Licensed **AGPL-3.0**; commercial use requires a separate commercial license (see README). Keep license headers and don't advise folding this code into proprietary/closed distributions without flagging the AGPL implication.
