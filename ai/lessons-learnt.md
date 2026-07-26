# Lessons Learnt (AI long-term memory)

This file captures hard-won, repo-specific context for future AI/agent sessions.
Prefer this over re-deriving facts. Update it when you discover new non-obvious truths.

## Build / toolchain
- **No CGO.** Everything builds with the pure-Go `modernc.org/sqlite` driver. Never introduce a cgo SQLite driver (`mattn/go-sqlite3`, etc.) — it breaks the build for consumers.
- **Go toolchain mismatch:** `go.mod` declares `go 1.26.3` but `.github/workflows/tests.yml` pins `go-version: '1.23'`. The CI pin is stale; align with `go.mod`.
- **Task runner is `task`** (Taskfile). `task test` == `go test ./...`; `task cover` generates `coverage.out` and opens HTML. Linters (errcheck, nilaway, gocritic, golangci-lint, gosec) must each be `go install`ed before first run — they are not in CI.

## Testing
- Tests are self-contained: they open an in-memory SQLite DB via `initDB()` (`store_test.go`, DSN `:memory:?parseTime=true`). No external DB, no fixtures, no services required.
- Run a focused test: `go test -run <TestName> ./<pkg>`.
- The `admin` package has heavy controller/handler test coverage (`*_test.go` alongside each handler) — follow that pattern when adding endpoints.

## Architecture
- Three layers, each a distinct `package`/dir:
  1. Root `statsstore` — core lib (`Store`, `Visitor`, geo-IP, bot filter). Entrypoint `NewStore(NewStoreOptions{...})`.
  2. `admin/` — embeddable dashboard returned as an `http.Handler` via `admin.New(options)`. UI is a **Vue.js SPA embedded as `admin/home/home.html` + `home.js`**; there is NO frontend build step. Don't add npm/bundler tooling.
  3. `examples/admin-demo/` — runnable demo (`go run ./examples/admin-demo`). `admin-demo.exe` is committed but gitignored (`*.exe`); rebuild from that dir.
- `docs/` is partly aspirational: `docs/overview.md` references goqu/`sb` builders that are NOT in `go.mod`. Trust code over prose in docs.

## Conventions (admin dashboard)
- **Use `dracory/req` for request params.** `req.GetString(r, "key")` and `req.GetStringOr(r, "key", "default")` read from both GET and POST. Do NOT use `r.URL.Query().Get()` or `r.FormValue()` directly.
- **Use `dracory/api` for JSON responses.** `api.Respond(w, r, api.SuccessWithData("msg", data))` / `api.Respond(w, r, api.Error("msg"))`. Always returns HTTP 200 — the `status` field in the JSON body (`"success"`/`"error"`) is what matters, not the HTTP status code. Do NOT use `fmt.Sprintf` for JSON or set `Content-Type` manually.
- **AJAX endpoints use POST.** All dashboard AJAX requests (`overview-ajax`, `comparison-ajax`, `dashboard-data-ajax`, `live-ajax`) use POST with `application/x-www-form-urlencoded` or `FormData` body. The `path` param stays in the URL query string; `action` and `period` go in the POST body.
- **Export is NOT AJAX.** CSV export (`action=export`) stays as GET — it's a download link (`<a :href="...">`), not a fetch call.
- **JS response envelope.** `fetchSection()` checks `data.status === 'success'` (not `=== 'error'`) and unwraps `data.data`. Never check `resp.ok` — HTTP status is always 200.

## Behavior quirks
- **Geo-IP enrichment is manual.** `VisitorEnhance` is NOT auto-invoked. It must be driven from a cron/goroutine, only updates rows where `country` is empty, and returns the count of fully processed (country+UA) rows. UA fields update even if geo lookup fails; country stays empty for retry.
- Default `DefaultGeoIPResolver` calls the free ip2c.org service with a built-in **2s delay between calls** (rate limiting) and a 24h in-memory cache; localhost/private IPs return `"UN"` without a call. Swap in a custom `GeoIPResolver` for production.
- Normal ingestion path is `VisitorRegister(r *http.Request)` (infers path/IP/UA). `VisitorQueryOptions` drives filtering/pagination for `VisitorList`/`VisitorCount`.

## License
- **AGPL-3.0**; commercial use needs a separate commercial license (see README). Keep license headers intact. Do NOT advise embedding this code into proprietary/closed distributions without flagging the AGPL implication.
