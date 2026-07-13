package statsstore

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// == IsBot TESTS =============================================================

func TestIsBot_KnownBots(t *testing.T) {
	botUAs := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
		"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
		"Mozilla/5.0 (compatible; DuckDuckBot/1.1; +http://duckduckgo.com/duckduckbot.html)",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
		"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)",
		"Mozilla/5.0 (compatible; PetalBot;+https://webmaster.petalsearch.com/site/petalbot)",
		"Mozilla/5.0 (compatible; Bytespider; spider'bot, +https://www.google.com/bot.html)",
		"curl/7.88.1",
		"Wget/1.21.4",
		"python-requests/2.31.0",
		"Go-http-client/1.1",
		"okhttp/4.12.0",
		"HeadlessChrome/120.0.6099.109",
		"phantomjs",
		"selenium/4.15.0",
		"Puppeteer/1.20.0",
		"Cypress/13.6.0",
		"Chrome-Lighthouse",
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"Twitterbot/1.0",
		"LinkedInBot/1.0 (compatible; LinkedInApp/1.0)",
		"TelegramBot (like TwitterBot)",
		"AppleBot/0.1; +http://www.apple.com/go/applebot.html)",
		"feedly/1.0",
		"Pingdom.com_bot_version_1.4_(http://www.pingdom.com)",
		"Datadog/Synthetic HTTP/1.0",
		"archive.org_bot",
		"IA_Archiver",
	}

	for _, ua := range botUAs {
		if !IsBot(ua) {
			t.Errorf("IsBot(%q) = false, expected true", ua)
		}
	}
}

func TestIsBot_RealBrowsers(t *testing.T) {
	realUAs := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	for _, ua := range realUAs {
		if IsBot(ua) {
			t.Errorf("IsBot(%q) = true, expected false", ua)
		}
	}
}

func TestIsBot_EmptyAndEdgeCases(t *testing.T) {
	if IsBot("") {
		t.Error("IsBot(\"\") = true, expected false")
	}
	if IsBot("Mozilla/5.0") {
		t.Error("IsBot(\"Mozilla/5.0\") = true, expected false")
	}
}

func TestIsBot_CaseInsensitive(t *testing.T) {
	if !IsBot("GOOGLEBOT/2.1") {
		t.Error("IsBot(\"GOOGLEBOT/2.1\") = false, expected true (case-insensitive)")
	}
	if !IsBot("CURL/7.88") {
		t.Error("IsBot(\"CURL/7.88\") = false, expected true (case-insensitive)")
	}
}

// == IsReferrerSpam TESTS ====================================================

func TestIsReferrerSpam_KnownSpam(t *testing.T) {
	spamReferrers := []string{
		"https://semalt.com/",
		"http://darodar.com/",
		"https://www.priceg.com/",
		"https://best-seo-offer.com/seo.php",
		"http://ilovevitaly.com/",
		"https://social-buttons.com/",
		"http://hulfingtonpost.com/",
		"https://adsterra.com/",
	}

	for _, ref := range spamReferrers {
		if !IsReferrerSpam(ref) {
			t.Errorf("IsReferrerSpam(%q) = false, expected true", ref)
		}
	}
}

func TestIsReferrerSpam_LegitimateReferrers(t *testing.T) {
	legitReferrers := []string{
		"https://www.google.com/search?q=test",
		"https://www.bing.com/search?q=test",
		"https://duckduckgo.com/",
		"https://github.com/user/repo",
		"https://twitter.com/user/status/123",
		"https://www.facebook.com/page",
		"https://www.linkedin.com/in/user",
		"https://stackoverflow.com/questions/123",
		"https://news.ycombinator.com/item?id=123",
		"",
	}

	for _, ref := range legitReferrers {
		if IsReferrerSpam(ref) {
			t.Errorf("IsReferrerSpam(%q) = true, expected false", ref)
		}
	}
}

func TestIsReferrerSpam_SubdomainMatch(t *testing.T) {
	if !IsReferrerSpam("https://sub.semalt.com/") {
		t.Error("IsReferrerSpam(\"https://sub.semalt.com/\") = false, expected true (subdomain match)")
	}
	if !IsReferrerSpam("https://www1.social-buttons.com/") {
		t.Error("IsReferrerSpam(\"https://www1.social-buttons.com/\") = false, expected true (subdomain match)")
	}
}

func TestIsReferrerSpam_BareDomain(t *testing.T) {
	if !IsReferrerSpam("semalt.com") {
		t.Error("IsReferrerSpam(\"semalt.com\") = false, expected true")
	}
	if !IsReferrerSpam("darodar.com") {
		t.Error("IsReferrerSpam(\"darodar.com\") = false, expected true")
	}
}

// == IsDataCenterIP TESTS ====================================================

func TestIsDataCenterIP_KnownRanges(t *testing.T) {
	ips := []string{
		"3.1.2.3",       // AWS
		"13.107.21.200", // AWS/Microsoft
		"52.1.2.3",      // AWS
		"35.192.1.2",    // GCP
		"35.200.1.2",    // GCP
		"20.1.2.3",      // Azure
		"40.1.2.3",      // Azure
		"159.65.1.2",    // DigitalOcean
		"167.99.1.2",    // DigitalOcean
		"129.146.1.2",   // Oracle Cloud
	}

	for _, ip := range ips {
		if !IsDataCenterIP(ip) {
			t.Errorf("IsDataCenterIP(%q) = false, expected true", ip)
		}
	}
}

func TestIsDataCenterIP_ResidentialIPs(t *testing.T) {
	ips := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"8.8.8.8",
		"1.1.1.1",
		"203.0.113.1",
		"198.51.100.1",
		"192.0.2.1",
	}

	for _, ip := range ips {
		if IsDataCenterIP(ip) {
			t.Errorf("IsDataCenterIP(%q) = true, expected false", ip)
		}
	}
}

func TestIsDataCenterIP_EmptyAndInvalid(t *testing.T) {
	if IsDataCenterIP("") {
		t.Error("IsDataCenterIP(\"\") = true, expected false")
	}
	if IsDataCenterIP("not-an-ip") {
		t.Error("IsDataCenterIP(\"not-an-ip\") = true, expected false")
	}
}

// == INTEGRATION: VisitorRegister with bot filtering =========================

func TestVisitorRegister_BotFiltered(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store.SetBotFilterEnabled(true)

	req := httptest.NewRequest(http.MethodGet, "/test-page", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	err = store.VisitorRegister(context.Background(), req)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, err := store.VisitorCount(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 visitors after bot filtering, got %d", count)
	}
}

func TestVisitorRegister_BotFilterDisabled(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Bot filter disabled by default
	req := httptest.NewRequest(http.MethodGet, "/test-page", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	err = store.VisitorRegister(context.Background(), req)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, err := store.VisitorCount(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 visitor (bot filter disabled), got %d", count)
	}
}

func TestVisitorRegister_RealBrowserNotFiltered(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store.SetBotFilterEnabled(true)

	req := httptest.NewRequest(http.MethodGet, "/test-page", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	err = store.VisitorRegister(context.Background(), req)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, err := store.VisitorCount(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 visitor (real browser not filtered), got %d", count)
	}
}

func TestVisitorRegister_ReferrerSpamFiltered(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store.SetBotFilterEnabled(true)

	req := httptest.NewRequest(http.MethodGet, "/test-page", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://semalt.com/")

	err = store.VisitorRegister(context.Background(), req)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, err := store.VisitorCount(context.Background(), VisitorQuery())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 visitors after referrer spam filtering, got %d", count)
	}
}

func TestVisitorRegister_ToggleOnOff(t *testing.T) {
	store, err := initStore()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Start with filtering enabled
	store.SetBotFilterEnabled(true)
	if !store.IsBotFilterEnabled() {
		t.Fatal("expected bot filter enabled")
	}

	req := httptest.NewRequest(http.MethodGet, "/page1", nil)
	req.Header.Set("User-Agent", "curl/7.88.1")

	err = store.VisitorRegister(context.Background(), req)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, _ := store.VisitorCount(context.Background(), VisitorQuery())
	if count != 0 {
		t.Fatalf("expected 0 visitors with filter on, got %d", count)
	}

	// Disable filtering, same bot UA should now be stored
	store.SetBotFilterEnabled(false)
	if store.IsBotFilterEnabled() {
		t.Fatal("expected bot filter disabled")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/page2", nil)
	req2.Header.Set("User-Agent", "curl/7.88.1")

	err = store.VisitorRegister(context.Background(), req2)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	count, _ = store.VisitorCount(context.Background(), VisitorQuery())
	if count != 1 {
		t.Fatalf("expected 1 visitor with filter off, got %d", count)
	}
}
