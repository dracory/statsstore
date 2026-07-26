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

## Dev workflow
- **Use `air` for hot-reload during development.** Run `air` from `examples/admin-demo/` (has `.air.toml`). It watches `.go`, `.html`, `.js`, `.css` files and auto-rebuilds/restarts the demo server on changes. This avoids manual restarts and repeated command approvals.
- **Chrome DevTools MCP for visual verification.** Use the `chrome-devtools` MCP server to navigate to `http://localhost:8080/admin/home`, take snapshots, inspect network requests, and verify AJAX responses. The a11y snapshot may show stale `{{ }}` Vue template tags — check actual Vue state via `evaluate_script` on `app._container._vnode.component.setupState` to confirm data is loaded.
- **Screenshots can't be saved to subdirectories.** The `mcp3_take_screenshot` tool with `filePath` only works within workspace root, not subdirectories. Use `fullPage: true` without `filePath` to capture inline.

## Settings page conventions
- The `admin/settings/` page follows the **same Vue.js SPA pattern** as `admin/home/`: embedded `settings.html` + `settings.js`, AJAX POST endpoints via `fetch` + `FormData`, `api.Respond` for all JSON responses, JS checks `data.status === 'success'` and unwraps `data.data`.
- Each AJAX handler is in its own file (`handle_list_ajax.go`, `handle_add_ip_ajax.go`, `handle_remove_ip_ajax.go`, `handle_delete_visitors_ajax.go`), matching the `admin/home/` convention.
- No HTML `<form>` elements — use `<div>` with `v-model` + `@click`/`@keyup.enter` for input, all actions go through Vue + AJAX.
- `<template v-if="loaded">` wraps all content so nothing renders in the DOM until Vue has mounted and loaded data.
- After every mutation (add/remove/delete), `loadIps()` is called to reload data from the server and confirm the change persisted.
- No Notiflix, no Axios, no PRG pattern — pure Vue + fetch.

## Visitor activity page conventions
- The `admin/visitor-activity/` page follows the same Vue SPA pattern: embedded `visitor_activity.html` + `visitor_activity.js`, AJAX POST for list data, `api.Respond` for JSON.
- The AJAX handler (`handle_list_ajax.go`) returns visitor rows as JSON with all formatting computed server-side (icon classes, location string, session label, duration, system summary). The Vue template just renders the pre-formatted strings — no client-side formatting logic.
- Pagination, filters, and per-page selector are all Vue-managed (reactive refs), not URL-based. Changing page/filters calls `fetchList()` which POSTs to `list-ajax` with `page`, `per_page`, `range`, `country`, `device` in the FormData body.
- The detail modal is a Bootstrap modal shown via `bootstrap.Modal` instance (ref to modal element), not via `data-bs-toggle` attributes. Visitor data is passed through `selectedVisitor` ref.
- CSV export stays as GET (`action=export`) — it's a download link, not AJAX. The export handler reads filters from `req.GetString` (same as list-ajax).
- Old server-rendered files (`card-visitor-activity.go`, `modal_visitor_detail.go`) were deleted; `helpers.go` was stripped to just `formatVisitorTimestamp` and `formatVisitDuration`. `ControllerData` was removed from `types.go`.
- `parseFiltersFromReq(r)` in `handle_list_ajax.go` replaces the old `parseFilters(url.Values)` — uses `req.GetString` instead of `r.URL.Query().Get`.
- HTMX and SweetAlert2 are no longer loaded — replaced by Vue reactivity and native `confirm()`.

## Behavior quirks
- **Geo-IP enrichment is manual.** `VisitorEnhance` is NOT auto-invoked. It must be driven from a cron/goroutine, only updates rows where `country` is empty, and returns the count of fully processed (country+UA) rows. UA fields update even if geo lookup fails; country stays empty for retry.
- Default `DefaultGeoIPResolver` calls the free ip2c.org service with a built-in **2s delay between calls** (rate limiting) and a 24h in-memory cache; localhost/private IPs return `"UN"` without a call. Swap in a custom `GeoIPResolver` for production.
- Normal ingestion path is `VisitorRegister(r *http.Request)` (infers path/IP/UA). `VisitorQueryOptions` drives filtering/pagination for `VisitorList`/`VisitorCount`.

## License
- **AGPL-3.0**; commercial use needs a separate commercial license (see README). Keep license headers intact. Do NOT advise embedding this code into proprietary/closed distributions without flagging the AGPL implication.
