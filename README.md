# Stats Store

[![Tests Status](https://github.com/dracory/statsstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/statsstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/statsstore)](https://goreportcard.com/report/github.com/dracory/statsstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/statsstore)](https://pkg.go.dev/github.com/dracory/statsstore)

Stats Store - a visitor stats storage implementation for Go.

## License
This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation
```
go get -u github.com/dracory/statsstore
```

## Setup

```golang
store, err := NewStore(NewStoreOptions{
	VisitorTableName:     "stats_visitor",
	DB:                 databaseInstance,
	AutomigrateEnabled: true,
})

```

## Geo-IP Enrichment

Visitor records are saved with an empty `country` field by default. To populate country codes (ISO 3166-1 alpha-2), configure a `GeoIPResolver` and call `VisitorEnhance` from a background task on your preferred schedule (e.g. every 5 minutes).

### Setup

```golang
store, err := NewStore(NewStoreOptions{
	VisitorTableName:   "stats_visitor",
	DB:                 databaseInstance,
	AutomigrateEnabled: true,
	GeoIPResolver:      NewDefaultGeoIPResolver(), // uses ip2c.org
	EnhanceBatchSize:   10,                        // records per call (default: 10)
})
```

### Running Enrichment

Call `VisitorEnhance` from a cron job, task scheduler, or goroutine ticker:

```golang
processed, err := store.VisitorEnhance(context.Background())
// processed = number of records that were successfully enriched
```

`VisitorEnhance` will:
1. Fetch up to `EnhanceBatchSize` visitor records where `country` is empty
2. For each record, parse the user agent to fill in browser, OS, device, and device type (if those fields are empty)
3. Look up the country via the configured `GeoIPResolver`
4. Update the record with the enriched data
5. Return the count of fully processed records (country + UA)

UA fields are updated even if the geo-IP lookup fails, but the country stays empty so the record gets retried on the next call. This makes `VisitorEnhance` a complete replacement for any custom post-processing task — it handles both UA parsing and country enrichment.

### Default Resolver (ip2c.org)

`DefaultGeoIPResolver` uses the free [ip2c.org](https://ip2c.org) service. It includes:

- **In-memory cache** with 24h TTL — avoids duplicate lookups for the same IP
- **Localhost/private IP detection** — returns `"UN"` (unknown) without making an HTTP call
- **Configurable timeout** (default: 5s) and HTTP client (for testing)

```golang
resolver := &DefaultGeoIPResolver{
	Endpoint:   "https://ip2c.org/",     // default
	Timeout:    5 * time.Second,         // default
	CacheTTL:   24 * time.Hour,          // default; set to 0 to disable caching
	HTTPClient: myCustomClient,          // optional; nil uses default
}
```

### Custom Resolver

Implement the `GeoIPResolver` interface to use any geo-IP provider (MaxMind, ipinfo, etc.):

```golang
type GeoIPResolver interface {
	Resolve(ctx context.Context, ip string) (string, error)
}
```

- Return an ISO2 country code (e.g. `"US"`, `"GB"`) on success
- Return `""` + error on failure (record stays empty for retry)
- Return `"UN"` + nil error for unresolvable IPs (localhost, private ranges, etc.)

Example with a local MaxMind database:

```golang
type maxmindResolver struct {
	db *geoip2.Reader
}

func (r *maxmindResolver) Resolve(ctx context.Context, ip string) (string, error) {
	addr := net.ParseIP(ip)
	if addr == nil {
		return statsstore.CountryUnknown, nil
	}
	country, err := r.db.Country(addr)
	if err != nil {
		return "", err
	}
	return country.Country.IsoCode, nil
}

store, _ := NewStore(NewStoreOptions{
	// ...
	GeoIPResolver: &maxmindResolver{db: myGeoIPDB},
})
```

## Admin Dashboard

The `admin/` subpackage provides a self-contained, reusable visitor analytics
dashboard. Any project can embed it by injecting a `LayoutInterface`, a
`statsstore.StoreInterface`, and optional callbacks:

```go
import (
    "github.com/dracory/statsstore"
    "github.com/dracory/statsstore/admin"
)

// Implement shared.LayoutInterface (SetTitle, SetBody, SetScripts, Render, etc.)
type myLayout struct { /* ... */ }

handler, err := admin.New(admin.Options{
    Store:             store,        // required
    Layout:            &myLayout{},  // required
    HomeURL:           "/admin",     // required
    BaseURL:           "/admin/stats", // defaults to "/admin/stats"
    CountryNameByIso2: countryNameByIso2, // optional
    AuthUserID:        func(r *http.Request) string { /* ... */ }, // optional
    FlashError:        func(w, r, msg, url string, delay int) string { /* ... */ }, // optional
})

http.Handle("/admin/stats", handler)
```

### Features

- **Dashboard**: Overview with total visitors, unique IPs, top paths/countries/browsers/OS/device types
- **Visitors**: Paginated visitor list with filtering (IP, country, browser, OS, device type, path, date range, bot/threat flags)
- **Sessions**: Visitor sessions grouped by IP with visit history
- **IP Details**: Per-IP deep dive with path history, flag as bot/threat, delete entries
- **Settings**: Excluded IP management, bot IP management, automated bot identification scan
- **Vue.js SPA**: Each page is a Vue.js single-page app with AJAX data loading
- **No host dependencies**: The package depends only on `github.com/dracory/statsstore` and injected interfaces

### Controllers

| Controller | Query param | Description |
|---|---|---|
| Dashboard | `controller=dashboard` (default) | Overview stats and top-N breakdowns |
| Visitors | `controller=visitors` | Paginated, filterable visitor list |
| Sessions | `controller=sessions` | Sessions grouped by IP |
| Settings | `controller=settings` | Excluded IPs, bot IPs, bot identification |
| IP Details | `controller=ip-details&ip=X` | Per-IP details and actions |

See `examples/admin-demo/` for a complete working example.

## Screenshots

### Dashboard

![Dashboard](examples/admin-demo/screenshots/screenshot-dashboard.png)

### Visitors

![Visitors](examples/admin-demo/screenshots/screenshot-visitors.png)

### Sessions

![Sessions](examples/admin-demo/screenshots/screenshot-sessions.png)

### Settings

![Settings](examples/admin-demo/screenshots/screenshot-settings.png)

### IP Details

![IP Details](examples/admin-demo/screenshots/screenshot-ip-details.png)
