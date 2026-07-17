package statsstore

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// == MOCK RESOLVER ============================================================

type mockGeoIPResolver struct {
	results map[string]string // ip -> country code
	errs    map[string]error  // ip -> error
	calls   int
}

func (m *mockGeoIPResolver) Resolve(ctx context.Context, ip string) (string, error) {
	m.calls++
	if err, ok := m.errs[ip]; ok {
		return "", err
	}
	if country, ok := m.results[ip]; ok {
		return country, nil
	}
	return CountryUnknown, nil
}

// == RESOLVER TESTS ===========================================================

func TestDefaultGeoIPResolverLocalhost(t *testing.T) {
	r := NewDefaultGeoIPResolver()

	country, err := r.Resolve(context.Background(), "127.0.0.1")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country != CountryUnknown {
		t.Fatalf("expected %q for localhost, got %q", CountryUnknown, country)
	}
}

func TestDefaultGeoIPResolverPrivateIP(t *testing.T) {
	r := NewDefaultGeoIPResolver()

	ips := []string{"10.0.0.1", "192.168.1.1", "172.16.0.1"}
	for _, ip := range ips {
		country, err := r.Resolve(context.Background(), ip)
		if err != nil {
			t.Fatal("unexpected error for "+ip+":", err)
		}
		if country != CountryUnknown {
			t.Fatalf("expected %q for %s, got %q", CountryUnknown, ip, country)
		}
	}
}

func TestDefaultGeoIPResolverEmptyIP(t *testing.T) {
	r := NewDefaultGeoIPResolver()

	country, err := r.Resolve(context.Background(), "")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country != CountryUnknown {
		t.Fatalf("expected %q for empty IP, got %q", CountryUnknown, country)
	}
}

func TestDefaultGeoIPResolverHTTPSuccess(t *testing.T) {
	// ip2c.org format: "1;US;United States"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("1;US;United States"))
	}))
	defer server.Close()

	r := &DefaultGeoIPResolver{
		Endpoint:   server.URL + "/",
		HTTPClient: server.Client(),
		CacheTTL:   0, // disable cache for this test
	}

	country, err := r.Resolve(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country != "US" {
		t.Fatalf("expected US, got %q", country)
	}
}

func TestDefaultGeoIPResolverHTTPNotFound(t *testing.T) {
	// ip2c.org format: "0" means not found
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("0"))
	}))
	defer server.Close()

	r := &DefaultGeoIPResolver{
		Endpoint:   server.URL + "/",
		HTTPClient: server.Client(),
		CacheTTL:   0,
	}

	country, err := r.Resolve(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country != CountryUnknown {
		t.Fatalf("expected %q for not-found, got %q", CountryUnknown, country)
	}
}

func TestDefaultGeoIPResolverHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	r := &DefaultGeoIPResolver{
		Endpoint:   server.URL + "/",
		HTTPClient: server.Client(),
		CacheTTL:   0,
	}

	_, err := r.Resolve(context.Background(), "8.8.8.8")
	if err != nil {
		t.Fatal("expected nil error for non-200 status (should return empty string, nil error), got:", err)
	}
}

func TestDefaultGeoIPResolverCache(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("1;GB;United Kingdom"))
	}))
	defer server.Close()

	r := &DefaultGeoIPResolver{
		Endpoint:   server.URL + "/",
		HTTPClient: server.Client(),
		CacheTTL:   GeoIPCacheTTLDefault,
	}

	// First call hits the server
	country1, err := r.Resolve(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country1 != "GB" {
		t.Fatalf("expected GB, got %q", country1)
	}

	// Second call should hit the cache
	country2, err := r.Resolve(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if country2 != "GB" {
		t.Fatalf("expected GB from cache, got %q", country2)
	}

	if calls != 1 {
		t.Fatalf("expected 1 HTTP call (second should be cached), got %d", calls)
	}
}

// == VISITOR ENHANCE TESTS ====================================================

func TestVisitorEnhanceNoResolver(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	_, err = store.VisitorEnhance(context.Background())
	if err == nil {
		t.Fatal("expected error when GeoIPResolver is not configured")
	}
}

func TestVisitorEnhanceNoRecords(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:      &mockGeoIPResolver{results: map[string]string{}},
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	processed, err := store.VisitorEnhance(context.Background())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if processed != 0 {
		t.Fatalf("expected 0 processed, got %d", processed)
	}
}

func TestVisitorEnhanceWithRecords(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	resolver := &mockGeoIPResolver{
		results: map[string]string{
			"8.8.8.8":   "US",
			"1.1.1.1":   "AU",
			"127.0.0.1": CountryUnknown,
		},
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:      resolver,
		EnhanceBatchSize:   10,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create visitors with empty country and empty UA fields
	visitors := []VisitorInterface{
		NewVisitor().SetIpAddress("8.8.8.8").SetCountry("").SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		NewVisitor().SetIpAddress("1.1.1.1").SetCountry("").SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"),
		NewVisitor().SetIpAddress("127.0.0.1").SetCountry("").SetUserAgent(""),
	}

	for _, v := range visitors {
		if err := store.VisitorCreate(ctx, v); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	processed, err := store.VisitorEnhance(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if processed != 3 {
		t.Fatalf("expected 3 processed, got %d", processed)
	}

	// Verify countries and UA fields were set
	for _, v := range visitors {
		found, err := store.VisitorFindByID(ctx, v.GetID())
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if found == nil {
			t.Fatal("visitor not found")
		}

		expected := resolver.results[v.GetIpAddress()]
		if found.GetCountry() != expected {
			t.Fatalf("expected country %q for IP %s, got %q",
				expected, v.GetIpAddress(), found.GetCountry())
		}

		// UA fields should be populated for visitors with a UA string
		if v.GetUserAgent() != "" {
			if found.GetUserBrowser() == "" {
				t.Fatalf("expected non-empty browser for IP %s", v.GetIpAddress())
			}
			if found.GetUserOs() == "" {
				t.Fatalf("expected non-empty OS for IP %s", v.GetIpAddress())
			}
		}
	}
}

func TestVisitorEnhanceLookupFailureLeavesEmpty(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	resolver := &mockGeoIPResolver{
		errs: map[string]error{
			"8.8.8.8": errors.New("network timeout"),
		},
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:      resolver,
		EnhanceBatchSize:   10,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	visitor := NewVisitor().
		SetIpAddress("8.8.8.8").
		SetCountry("").
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	if err := store.VisitorCreate(ctx, visitor); err != nil {
		t.Fatal("unexpected error:", err)
	}

	processed, err := store.VisitorEnhance(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if processed != 0 {
		t.Fatalf("expected 0 processed (lookup failed), got %d", processed)
	}

	found, err := store.VisitorFindByID(ctx, visitor.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("visitor not found")
	}

	// Country should still be empty for retry
	if found.GetCountry() != "" {
		t.Fatalf("expected empty country after failed lookup, got %q", found.GetCountry())
	}

	// UA fields should be populated even though geo-IP failed
	if found.GetUserBrowser() != "Chrome" {
		t.Fatalf("expected browser 'Chrome', got %q", found.GetUserBrowser())
	}
	if found.GetUserOs() != "Windows" {
		t.Fatalf("expected OS 'Windows', got %q", found.GetUserOs())
	}
	if found.GetUserDeviceType() != "desktop" {
		t.Fatalf("expected device type 'desktop', got %q", found.GetUserDeviceType())
	}
}

func TestVisitorEnhanceBulkUpdateByIP(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	resolver := &mockGeoIPResolver{
		results: map[string]string{
			"8.8.8.8": "US",
		},
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:      resolver,
		EnhanceBatchSize:   2,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create 5 visitors with the same IP but empty country
	visitorIDs := []string{}
	for i := 0; i < 5; i++ {
		v := NewVisitor().SetIpAddress("8.8.8.8").SetCountry("")
		if err := store.VisitorCreate(ctx, v); err != nil {
			t.Fatal("unexpected error:", err)
		}
		visitorIDs = append(visitorIDs, v.GetID())
	}

	processed, err := store.VisitorEnhance(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Only 2 records are in the batch, so only 2 are "fully processed"
	if processed != 2 {
		t.Fatalf("expected 2 processed (batch size), got %d", processed)
	}

	// But ALL 5 records should have country set via the bulk update
	for _, id := range visitorIDs {
		found, err := store.VisitorFindByID(ctx, id)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if found == nil {
			t.Fatal("visitor not found:", id)
		}
		if found.GetCountry() != "US" {
			t.Fatalf("expected country 'US' for visitor %s (bulk-updated), got %q",
				id, found.GetCountry())
		}
	}

	// Resolver should have been called only once for the unique IP
	if resolver.calls != 1 {
		t.Fatalf("expected 1 resolve call (1 unique IP), got %d", resolver.calls)
	}
}

func TestVisitorEnhanceBatchSize(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	resolver := &mockGeoIPResolver{
		results: map[string]string{
			"8.8.8.8": "US",
		},
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		VisitorTableName:   "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:      resolver,
		EnhanceBatchSize:   2,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create 5 visitors with empty country
	for i := 0; i < 5; i++ {
		v := NewVisitor().SetIpAddress("8.8.8.8").SetCountry("")
		if err := store.VisitorCreate(ctx, v); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	processed, err := store.VisitorEnhance(ctx)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if processed != 2 {
		t.Fatalf("expected 2 processed (batch size), got %d", processed)
	}
}
