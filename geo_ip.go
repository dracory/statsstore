package statsstore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"sync"
	"time"
)

// == CONSTANTS ================================================================

const (
	// GeoIPEndpointDefault is the default IP geolocation endpoint (ip2c.org).
	GeoIPEndpointDefault = "https://ip2c.org/"

	// GeoIPTimeoutDefault is the default timeout for IP geolocation lookups.
	GeoIPTimeoutDefault = 5 * time.Second

	// CountryUnknown is used for IPs that cannot be resolved (e.g. localhost).
	CountryUnknown = "UN"

	// GeoIPCacheTTLDefault is the default TTL for the in-memory IP cache.
	GeoIPCacheTTLDefault = 24 * time.Hour
)

// == INTERFACE ================================================================

// GeoIPResolver resolves an IP address to an ISO 3166-1 alpha-2 country code.
// Implementations may use an external API, a local MaxMind database, or any
// other source. A nil/empty result with no error means the IP could not be
// resolved and the country should be left empty for a later retry.
type GeoIPResolver interface {
	// Resolve returns an ISO2 country code (e.g. "US", "GB") for the given IP.
	// If the IP cannot be resolved, it returns CountryUnknown and nil error.
	// If the lookup itself fails (network error, bad response), it returns
	// an empty string and the error so the caller can leave the field empty
	// for a retry.
	Resolve(ctx context.Context, ip string) (string, error)
}

// == DEFAULT IMPLEMENTATION ===================================================

// DefaultGeoIPResolver implements GeoIPResolver using the ip2c.org service.
// It includes an in-memory cache with TTL to avoid duplicate lookups for the
// same IP within the cache window.
type DefaultGeoIPResolver struct {
	Endpoint   string        // default: GeoIPEndpointDefault
	Timeout    time.Duration // default: GeoIPTimeoutDefault
	HTTPClient *http.Client  // injectable for testing; if nil, a default client is used
	CacheTTL   time.Duration // default: GeoIPCacheTTLDefault; set to 0 to disable caching

	cache   map[string]cacheEntry
	cacheMu sync.RWMutex
}

type cacheEntry struct {
	country string
	expires time.Time
}

// NewDefaultGeoIPResolver creates a DefaultGeoIPResolver with sensible defaults.
func NewDefaultGeoIPResolver() *DefaultGeoIPResolver {
	return &DefaultGeoIPResolver{
		Endpoint: GeoIPEndpointDefault,
		Timeout:  GeoIPTimeoutDefault,
		CacheTTL: GeoIPCacheTTLDefault,
		cache:    make(map[string]cacheEntry),
	}
}

// Resolve looks up the country code for the given IP via ip2c.org.
func (r *DefaultGeoIPResolver) Resolve(ctx context.Context, ip string) (string, error) {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return CountryUnknown, nil
	}
	// Normalize IPv4-mapped IPv6 to plain IPv4.
	addr = addr.Unmap()

	// Localhost and private IPs cannot be geolocated
	if isLocalOrPrivateIP(addr) {
		return CountryUnknown, nil
	}

	canonicalIP := addr.String()

	// Check cache
	if r.CacheTTL > 0 {
		if country, ok := r.cacheGet(canonicalIP); ok {
			return country, nil
		}
	}

	endpoint := r.Endpoint
	if endpoint == "" {
		endpoint = GeoIPEndpointDefault
	}

	timeout := r.Timeout
	if timeout == 0 {
		timeout = GeoIPTimeoutDefault
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, endpoint+url.PathEscape(canonicalIP), nil)
	if err != nil {
		return "", err
	}

	client := r.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", errors.New("geo IP: nil response from HTTP client")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// ip2c.org response format: "status;countryCode;countryName;..."
	// status "1" = found, "0" = not found
	parts := strings.Split(string(body), ";")
	if len(parts) < 2 {
		return CountryUnknown, nil
	}

	status := strings.TrimSpace(parts[0])
	if status != "1" {
		return CountryUnknown, nil
	}

	code := strings.TrimSpace(parts[1])
	if code == "" {
		return CountryUnknown, nil
	}

	// Cache the result
	if r.CacheTTL > 0 {
		r.cacheSet(canonicalIP, code)
	}

	return code, nil
}

// == CACHE METHODS ============================================================

func (r *DefaultGeoIPResolver) cacheGet(ip string) (string, bool) {
	r.cacheMu.RLock()
	defer r.cacheMu.RUnlock()

	entry, ok := r.cache[ip]
	if !ok {
		return "", false
	}

	if time.Now().After(entry.expires) {
		return "", false
	}

	return entry.country, true
}

func (r *DefaultGeoIPResolver) cacheSet(ip, country string) {
	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()

	if r.cache == nil {
		r.cache = make(map[string]cacheEntry)
	}

	r.cache[ip] = cacheEntry{
		country: country,
		expires: time.Now().Add(r.CacheTTL),
	}
}

// == HELPERS ==================================================================

// isLocalOrPrivateIP returns true for addresses that cannot be geolocated
// to a real country: loopback, RFC 1918 / ULA private ranges, link-local,
// multicast, unspecified, and carrier-grade NAT (100.64.0.0/10).
func isLocalOrPrivateIP(ip netip.Addr) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return true
	}

	// CGNAT 100.64.0.0/10 is not reported as private by netip.
	if ip.Is4() {
		a4 := ip.As4()
		if a4[0] == 100 && a4[1]&0xc0 == 64 {
			return true
		}
	}

	return false
}
