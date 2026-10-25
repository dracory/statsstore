package statsstore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVisitorRegisterCapturesUtmParams(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	req := httptest.NewRequest(http.MethodGet,
		"/landing?utm_source=google&utm_medium=cpc&utm_campaign=spring&utm_term=analytics&utm_content=banner-a",
		nil)
	req.Header.Set("Referer", "https://news.example.com/story/1?utm_source=rss")

	if err := store.VisitorRegister(context.Background(), req); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, err := store.VisitorList(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 visitor, got %d", len(list))
	}

	v := list[0]
	if v.GetUtmSource() != "google" {
		t.Fatalf("expected utm_source google, got %s", v.GetUtmSource())
	}
	if v.GetUtmMedium() != "cpc" {
		t.Fatalf("expected utm_medium cpc, got %s", v.GetUtmMedium())
	}
	if v.GetUtmCampaign() != "spring" {
		t.Fatalf("expected utm_campaign spring, got %s", v.GetUtmCampaign())
	}
	if v.GetUtmTerm() != "analytics" {
		t.Fatalf("expected utm_term analytics, got %s", v.GetUtmTerm())
	}
	if v.GetUtmContent() != "banner-a" {
		t.Fatalf("expected utm_content banner-a, got %s", v.GetUtmContent())
	}
	if v.GetReferrerHost() != "news.example.com" {
		t.Fatalf("expected referrer_host news.example.com, got %s", v.GetReferrerHost())
	}
}

func TestVisitorRegisterNoUtmParams(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/plain", nil)
	if err := store.VisitorRegister(context.Background(), req); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, _ := store.VisitorList(context.Background(), VisitorQuery())
	if len(list) != 1 {
		t.Fatalf("expected 1 visitor, got %d", len(list))
	}
	v := list[0]
	if v.GetUtmSource() != "" || v.GetReferrerHost() != "" {
		t.Fatal("expected empty utm/referrer fields")
	}
}

func TestVisitorRegisterReferrerHostNormalization(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Referer", "HTTPS://WWW.Example.COM:443/Path")
	if err := store.VisitorRegister(context.Background(), req); err != nil {
		t.Fatal("unexpected error:", err)
	}

	list, _ := store.VisitorList(context.Background(), VisitorQuery())
	if list[0].GetReferrerHost() != "example.com" {
		t.Fatalf("expected lowercase host without port and www, got %s", list[0].GetReferrerHost())
	}
}

func TestVisitorQueryUtmFilters(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	seed := func(fingerprint, source, campaign string) {
		v := NewVisitor().
			SetFingerprint(fingerprint).
			SetIpAddress("1.2.3.4").
			SetPath("/x").
			SetUtmSource(source).
			SetUtmCampaign(campaign)
		if err := store.VisitorCreate(context.Background(), v); err != nil {
			t.Fatal("unexpected error:", err)
		}
	}
	seed("fp-1", "google", "spring")
	seed("fp-2", "newsletter", "spring")
	seed("fp-3", "google", "summer")

	count, err := store.VisitorCount(context.Background(), VisitorQuery().SetUtmSource("google"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 google visitors, got %d", count)
	}

	count, err = store.VisitorCount(context.Background(), VisitorQuery().
		SetUtmSource("google").SetUtmCampaign("summer"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 google+summer visitor, got %d", count)
	}

	// Set-but-empty source matches nothing meaningful (filter skipped on empty)
	count, err = store.VisitorCount(context.Background(), VisitorQuery().SetUtmSource(""))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 3 {
		t.Fatalf("expected empty source filter to be skipped, got %d", count)
	}
}

func TestVisitorUtmFromExistingData(t *testing.T) {
	v := NewVisitorFromExistingData(map[string]string{
		COLUMN_UTM_SOURCE:    "newsletter",
		COLUMN_UTM_CAMPAIGN:  "launch",
		COLUMN_REFERRER_HOST: "mail.example.com",
	})
	if v.GetUtmSource() != "newsletter" || v.GetUtmCampaign() != "launch" || v.GetReferrerHost() != "mail.example.com" {
		t.Fatal("expected utm fields loaded from data map")
	}
}
