# Plan: Replace `admin/` with `statsadmin` in the statsstore library

## Goal

Replace the existing `admin/` package in the `github.com/dracory/statsstore`
library with the `statsadmin` package from `coursethread.com/pkg/statsadmin`,
because `statsadmin` works better (richer dashboard, sessions, IP details, bot
identification, settings management).

**Guiding principle:** the resulting `admin/` package must be a drop-in,
self-contained library that **any project** can embed via
`admin.New(admin.Options{...})` — the same reusability contract the old
`admin/` package already offers. Other projects must NOT have to reimplement
`statsadmin` themselves; they provide their own `LayoutInterface`, a
`statsstore.StoreInterface`, and a few optional callbacks, and get the full
dashboard for free.

The statsstore library is a **standalone public AGPL library**. It cannot import
`project/internal/*`. The current `statsadmin` is tightly coupled to
coursethread internals, so the swap requires **decoupling** it before it can
live inside the library. Decoupling is not a chore — it is the entire value of
this task: it turns a coursethread-private dashboard into a reusable library
feature.

## Current state

### Old `admin/` (to be removed)
- Location: `D:\PROJECTs\_modules_dracory\statsstore\admin\`
- Subpackages: `home/`, `visitor-activity/`, `visitor-paths/`,
  `page-view-activity/`, `settings/`, `shared/`
- Design: framework-agnostic. Consumer injects `shared.LayoutInterface` and
  `statsstore.StoreInterface` via `shared.ControllerOptions`. No dependency on
  any host application.
- Consumer: `examples/admin-demo/main.go` (implements `LayoutInterface`
  inline), plus README screenshots.

### New `statsadmin/` (to be ported in)
- Location: `D:\PROJECTs\coursethread.com\pkg\statsadmin\`
- Subpackages: `dashboard/`, `visitors/`, `sessions/`, `settings/`,
  `ipdetails/`, `shared/`
- Design: server-rendered Vue.js SPA per controller, AJAX data endpoints,
  bot-IP identification, excluded-IP management, threat flagging.
- **Problem:** depends on coursethread internals:
  - `project/internal/app` — `AppInterface` for `GetStatsStore()`,
    `GetGeoStore()`, `GetCacheStore()`
  - `project/internal/links` — `links.ADMIN_STATS`, `links.Admin().Stats()`,
    `links.Admin().Home()`
  - `project/internal/layouts` — `layouts.NewAdminLayout()`,
    `layouts.Breadcrumbs()`
  - `project/internal/helpers` — `helpers.GetAuthUser()`,
    `helpers.ToFlashError()`
  - `project/internal/config`, `project/internal/testutils` — tests only

## Decoupling strategy

Replace each `project/internal/*` dependency with a library-local abstraction,
mirroring the pattern the old `admin/` package already uses
(`LayoutInterface` + `ControllerOptions`). The end result is a package that
depends only on `github.com/dracory/statsstore` (its own root) plus a handful
of injected interfaces/callbacks — no host-application knowledge.

| Coupling | Replacement |
|---|---|
| `app.AppInterface` | Extend the existing `shared.ControllerOptions` (from `admin/shared/types.go`) with the new optional fields the richer dashboard needs. Do NOT rename existing fields/methods. Controllers take `shared.ControllerOptions` as before. |
| `internal/layouts` | Reuse the existing `shared.LayoutInterface` (already in `admin/shared/types.go`) **unchanged**. Port it into the new `shared/`. The demo implements it; consumers implement their own. No new methods added. |
| `internal/links` | `shared.NewLinks(baseURL)` already exists but falls back to `links.ADMIN_STATS`. Remove that fallback; require `baseURL` to be passed in via `ControllerOptions`. Build URLs from `baseURL` + query params inline. |
| `helpers.GetAuthUser` | Optional `AuthUserID func(*http.Request) string` callback in `ControllerOptions` (same pattern as `statsadmin.AdminOptions`). When nil, auth check is skipped. |
| `helpers.ToFlashError` | Optional `FlashError` callback in `ControllerOptions` matching the `helpers.ToFlashError` signature. When nil, controllers fall back to a plain `http.Redirect`. **No change to `LayoutInterface`.** |
| `geostore.CountryList` | Use the existing `CountryNameByIso2` callback already in `ControllerOptions`. Do NOT add `dracory/geostore` as a dependency. |
| `internal/config`, `internal/testutils` (tests) | Rewrite tests to construct a `statsstore.Store` with `:memory:` sqlite (same pattern as existing `admin/home/*_test.go` and `store_test.go`). Drop `testutils`/`config` imports. |

## Step-by-step plan

### Phase 1 — Prepare the new package inside the library

1. **Copy** `coursethread.com/pkg/statsadmin/` into
   `_modules_dracory/statsstore/admin/` (overwrite the old contents).
   - Keep the package name `admin` (not `statsadmin`) so the public import path
     `github.com/dracory/statsstore/admin` stays stable for existing consumers
     (Decision #1).

2. **Delete** the old subpackages that have no equivalent in statsadmin:
   `visitor-activity/`, `visitor-paths/`, `page-view-activity/`, and the old
   `home/` (replaced by `dashboard/`).

### Phase 2 — Decouple from coursethread internals

3. **`shared/types.go`** — port the existing `LayoutInterface` and
   `ControllerOptions` from the old `admin/shared/types.go` **unchanged**, then
   **extend** `ControllerOptions` with the new optional fields the richer
   dashboard needs:
   - `BaseURL string` (replaces `internal/links` dependency)
   - `AuthUserID func(*http.Request) string` (replaces `helpers.GetAuthUser`)
   - `FlashError func(store, w, r, message, redirectURL string, delay int) string` (replaces `helpers.ToFlashError`; nil → plain redirect)
   - `Logger *slog.Logger` (already in old struct)
   - Keep existing `Store`, `Layout`, `HomeURL`, `WebsiteUrl`,
     `CountryNameByIso2` fields as-is.
   - Do NOT add a `GeoStore` field or interface — country resolution uses the
     existing `CountryNameByIso2` callback.

4. **`shared/links.go`** — remove `import "project/internal/links"`. Build URLs
   from the configured `baseURL` only. No `links.ADMIN_STATS` fallback.

5. **`shared/helpers.go`** — remove `import "project/internal/app"`. Change
   `CountryNameResolver` to use the `CountryNameByIso2` callback from
   `ControllerOptions` instead of `AppInterface`/`geostore`.

6. **`shared/bot_ips.go`** — remove `import "project/internal/app"`. Change all
   functions to take `statsstore.StoreInterface` (which already has
   `SettingGet`/`SettingSet`) instead of `AppInterface`.

7. **`shared/visitor_bot_flag.go`** — remove `project/internal/app` import;
   use `statsstore.StoreInterface`.

8. **Each controller (`dashboard/`, `visitors/`, `sessions/`, `settings/`,
   `ipdetails/`)**:
   - Replace `controller.app` field with `controller.opts shared.ControllerOptions`.
   - Replace `controller.app.GetStatsStore()` with `controller.opts.Store`.
   - Replace `helpers.GetAuthUser(r)` with `controller.opts.AuthUserID(r)`
     (nil → skip auth check).
   - Replace `helpers.ToFlashError(...)` with `controller.opts.FlashError(...)`
     (nil → `http.Redirect` fallback).
   - Replace `layouts.NewAdminLayout(...)` / `layouts.Breadcrumbs(...)` with
     `controller.opts.Layout` calls (SetTitle/SetBody/SetScriptURLs/SetScripts/
     Render) + a local `shared.Breadcrumbs(...)` helper (port from old
     `admin/shared/components.go`).
   - Replace `links.Admin().Stats()` / `links.Admin().Home()` with
     `shared.NewLinks(controller.opts.BaseURL)` calls.

9. **`admin.go` / `routes.go`** — replace `app.AppInterface` with
   `shared.ControllerOptions`. `New()` returns `http.Handler` like the old
   package, keeping the same `Options`/`New` shape.

### Phase 3 — Tests

10. **Rewrite tests** to drop `project/internal/config` and
    `project/internal/testutils`. Use `modernc.org/sqlite` `:memory:` +
    `statsstore.NewStore` (same pattern as existing
    `admin/home/home_controller_test.go`). Keep the same test coverage
    (bot identification, flag bot, remove entries, load visitors, etc.).

### Phase 4 — Demo & docs

11. **Update `examples/admin-demo/main.go`** to construct the new
    `admin.New(admin.Options{...})` with the new `ControllerOptions`/`Layout`.
    The demo's inline `layout` struct already implements a compatible
    `LayoutInterface`; adjust method set if the interface changed.

12. **Update `README.md`** — refresh the admin section: new screenshots
    (dashboard, sessions, IP details, settings), updated usage snippet
    reflecting the new `Options`.

13. **Update `AGENTS.md`** in the statsstore repo to reflect the new
    `admin/` subpackages (`dashboard`, `visitors`, `sessions`, `settings`,
    `ipdetails`).

### Phase 5 — Verify

14. `go build ./...` in the statsstore repo.
15. `go test ./...` in the statsstore repo.
16. `go run ./examples/admin-demo` and exercise the dashboard end-to-end.
17. `task cover` for coverage sanity.

## What is NOT in scope (this task)

- The core `github.com/dracory/statsstore` store (store.go, visitor.go,
  geo_ip.go, bot_filter.go) is untouched.

## Follow-up (separate task, after this one lands)

- **coursethread.com** currently maintains its own `pkg/statsadmin` copy
  (decision 060, task 194 — completed). Once the library `admin/` package is
  the decoupled, reusable version, coursethread.com should switch
  `internal/controllers/admin/stats/routes.go` back to importing
  `github.com/dracory/statsstore/admin` and delete its private
  `pkg/statsadmin` — so it consumes the shared library like any other project,
  instead of keeping a fork. This is a separate task with its own decision
  record; it is explicitly out of scope here to keep this task focused on the
  library.

## Decisions (locked by GOD)

1. **Package name**: keep `admin`. No rename. The public import path
   `github.com/dracory/statsstore/admin` stays stable for existing consumers.
2. **Interface naming**: no changes. Reuse the existing `LayoutInterface` and
   `ControllerOptions` names/pattern from the old `admin/shared/types.go`
   verbatim. Do not invent new interface names. Extend `ControllerOptions` with
   the new optional fields the richer dashboard needs, but do not rename
   existing fields or methods.
3. **Flash messages**: callback. Do NOT add `SetFlash` to `LayoutInterface`
   (that would change the interface). Instead, accept an optional
   `FlashError` callback in `ControllerOptions` (signature matches the
   `helpers.ToFlashError` use sites: store, w, r, message, redirectURL,
   delaySeconds → string). When nil, controllers fall back to a plain
   `http.Redirect` to the redirect URL. Keeps the existing
   `LayoutInterface` untouched.
4. **GeoStore**: callback. Do NOT add `dracory/geostore` as a dependency.
   `CountryNameResolver` takes the existing `CountryNameByIso2`
   `func(iso2Code string) (string, error)` callback (already in
   `ControllerOptions`). Keeps the library lean.
5. **Auth**: callback. Reuse the existing `AuthUserID` pattern from
   `statsadmin.AdminOptions` as an optional `ControllerOptions` field. No
   `helpers.GetAuthUser` dependency.
6. **AGPL notice**: the ported code retains AGPL-3.0 headers; no license
   conflict since both are AGPL under the same author.

## File-level change summary

```
statsstore/
├── admin/                      # OVERWRITTEN with statsadmin contents
│   ├── admin.go                # new (from statsadmin.go, decoupled)
│   ├── routes.go               # new (decoupled)
│   ├── shared/                 # new types.go, links.go, helpers.go, bot_ips.go, ...
│   ├── dashboard/              # new
│   ├── visitors/               # new
│   ├── sessions/               # new
│   ├── settings/               # new
│   ├── ipdetails/              # new
│   ├── home/                   # DELETED (replaced by dashboard/)
│   ├── visitor-activity/       # DELETED
│   ├── visitor-paths/          # DELETED
│   └── page-view-activity/     # DELETED
├── examples/admin-demo/main.go # UPDATED wiring
├── README.md                   # UPDATED
└── AGENTS.md                   # UPDATED
```
