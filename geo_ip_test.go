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
		DB:              db,
		VisitorTableName: "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:   &mockGeoIPResolver{results: map[string]string{}},
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
			"8.8.8.8":     "US",
			"1.1.1.1":     "AU",
			"127.0.0.1":   CountryUnknown,
		},
	}

	store, err := NewStore(NewStoreOptions{
		DB:               db,
		VisitorTableName: "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:    resolver,
		EnhanceBatchSize: 10,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create visitors with empty country
	visitors := []VisitorInterface{
		NewVisitor().SetIpAddress("8.8.8.8").SetCountry(""),
		NewVisitor().SetIpAddress("1.1.1.1").SetCountry(""),
		NewVisitor().SetIpAddress("127.0.0.1").SetCountry(""),
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

	// Verify countries were set
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
		DB:               db,
		VisitorTableName: "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:    resolver,
		EnhanceBatchSize: 10,
	})
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	visitor := NewVisitor().SetIpAddress("8.8.8.8").SetCountry("")
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

	// Country should still be empty for retry
	found, err := store.VisitorFindByID(ctx, visitor.GetID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if found == nil {
		t.Fatal("visitor not found")
	}
	if found.GetCountry() != "" {
		t.Fatalf("expected empty country after failed lookup, got %q", found.GetCountry())
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
		DB:               db,
		VisitorTableName: "visitor_table",
		AutomigrateEnabled: true,
		GeoIPResolver:    resolver,
		EnhanceBatchSize: 2,
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
