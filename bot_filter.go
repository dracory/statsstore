package statsstore

import (
	"net/netip"
	"net/url"
	"strings"
)

// == BOT USER-AGENT PATTERNS ==================================================

// botUserAgentBroadPatterns contains lowercase substrings that indicate a bot,
// crawler, spider, or automated tool but are short enough to cause false
// positives if matched as plain substrings (e.g. "bot" inside "RoboBotonic").
// These are matched with a word-boundary check: the pattern must be followed
// by a non-alphanumeric character (or end of string).
var botUserAgentBroadPatterns = []string{
	"bot",
	"crawler",
	"spider",
	"scraper",
	"slurp",
}

// botUserAgentSpecificPatterns contains lowercase substrings that are specific
// enough to be safely matched as plain case-insensitive substrings.
var botUserAgentSpecificPatterns = []string{
	"baidu",
	"yandex",
	"ahrefs",
	"semrush",
	"duckduckbot",
	"facebookexternalhit",
	"twitterbot",
	"linkedinbot",
	"telegrambot",
	"applebot",
	"petalbot",
	"bytespider",
	"curl",
	"wget",
	"python-requests",
	"go-http-client",
	"okhttp",
	"headless",
	"phantom",
	"selenium",
	"puppeteer",
	"cypress",
	"lighthouse",
	"w3c_validator",
	"chrome-lighthouse",
	"google-structured-data-testing-tool",
	"google-page-speed-insights",
	"feedly",
	"uptime",
	"pingdom",
	"datadog",
	"newrelic",
	"site24x7",
	"node-fetch",
	"axios",
	"postman",
	"insomnia",
	"httpx",
	"scrapy",
	"mechanize",
	"mechanize-go",
	"colly",
	"heritrix",
	"nutch",
	"archive.org_bot",
	"ia_archiver",
	"wayback",
}

// == REFERRER SPAM DOMAINS ====================================================

// referrerSpamDomains contains lowercase domain names known to engage in
// referrer spam. Matching is case-insensitive against the referrer host.
var referrerSpamDomains = map[string]bool{
	"semalt.com":                       true,
	"semalt.semalt.com":                true,
	"darodar.com":                      true,
	"priceg.com":                       true,
	"7makemoneyonline.com":             true,
	"buttons-for-website.com":          true,
	"buttons-for-your-site.com":        true,
	"bestwebsitesawards.com":           true,
	"hulfingtonpost.com":               true,
	"best-seo-offer.com":               true,
	"offers.bycontext.com":             true,
	"www1.social-buttons.com":          true,
	"social-buttons.com":               true,
	"free-share-buttons.com":           true,
	"trafficmonetize.com":              true,
	"webmonetizer.net":                 true,
	"ranksonic.info":                   true,
	"ranksonic.org":                    true,
	"ranksonic.com":                    true,
	"site1.floating-share-buttons.com": true,
	"floating-share-buttons.com":       true,
	"get-free-traffic-now.com":         true,
	"quality-traffic.com":              true,
	"traffic2cash.com":                 true,
	"traffic2money.com":                true,
	"cyber-monday.ga":                  true,
	"cyber-monday.biz":                 true,
	"free-traffic.xyz":                 true,
	"buy-cheap-online.com":             true,
	"erot.co":                          true,
	"palvira.com":                      true,
	"gowildpass.com":                   true,
	"torture.ml":                       true,
	"xn--80adgbcm5aj1b5bfh.xn--p1ai":   true,
	"ilovevitaly.com":                  true,
	"ilovevitaly.ru":                   true,
	"ilovevitaly.org":                  true,
	"ilovevitaly.co":                   true,
	"econom.co":                        true,
	"blackhatworth.com":                true,
	"adviceforum.info":                 true,
	"hongfanji.com":                    true,
	"howtostopreferralspam.eu":         true,
	"humanorightswatch.org":            true,
	"o-o-6-o-o.com":                    true,
	"o-o-8-o-o.com":                    true,
	"rank-checker.online":              true,
	"referrerdisabler.com":             true,
	"success-seo.com":                  true,
	"videos-for-your-business.com":     true,
	"get-clicky.com":                   true,
	"snip.to":                          true,
	"snip.it":                          true,
	"adsterra.com":                     true,
}

// == DATA CENTER CIDR RANGES ==================================================

// dataCenterCIDRs contains CIDR ranges commonly associated with cloud
// providers and data centers. Traffic from these ranges is more likely
// to be automated.
var dataCenterCIDRs = []string{
	// AWS
	"3.0.0.0/9",
	"13.0.0.0/8",
	"15.0.0.0/8",
	"18.0.0.0/8",
	"34.0.0.0/8",
	"52.0.0.0/8",
	"54.0.0.0/8",
	"99.77.0.0/16",
	// GCP
	"35.184.0.0/13",
	"35.192.0.0/14",
	"35.196.0.0/15",
	"35.198.0.0/16",
	"35.199.0.0/17",
	"35.200.0.0/13",
	"35.208.0.0/12",
	"35.224.0.0/12",
	// Azure
	"4.128.0.0/12",
	"4.144.0.0/12",
	"4.160.0.0/12",
	"20.0.0.0/8",
	"40.0.0.0/8",
	// DigitalOcean
	"159.65.0.0/16",
	"159.203.0.0/16",
	"165.22.0.0/16",
	"167.99.0.0/16",
	"206.189.0.0/16",
	// Oracle Cloud
	"129.146.0.0/16",
	"129.148.0.0/16",
	"129.213.0.0/16",
	"140.238.0.0/16",
	"152.70.0.0/16",
}

// parsedDataCenterCIDRs is lazily initialized from dataCenterCIDRs.
var parsedDataCenterCIDRs []netip.Prefix

func init() {
	for _, cidr := range dataCenterCIDRs {
		prefix, err := netip.ParsePrefix(cidr)
		if err == nil {
			parsedDataCenterCIDRs = append(parsedDataCenterCIDRs, prefix)
		}
	}
}

// == PUBLIC FUNCTIONS =========================================================

// IsBot checks whether a user-agent string matches known bot/crawler patterns.
// Broad patterns (bot, crawler, spider, scraper, slurp) are matched with a
// word-boundary check to avoid false positives (e.g. "bot" inside a non-bot
// word). Specific patterns (semrush, curl, googlebot, etc.) are matched as
// plain case-insensitive substrings.
func IsBot(userAgent string) bool {
	if userAgent == "" {
		return false
	}
	uaLower := strings.ToLower(userAgent)

	// Check broad patterns with word-boundary matching.
	for _, pattern := range botUserAgentBroadPatterns {
		if matchWordBoundary(uaLower, pattern) {
			return true
		}
	}

	// Check specific patterns with plain substring matching.
	for _, pattern := range botUserAgentSpecificPatterns {
		if strings.Contains(uaLower, pattern) {
			return true
		}
	}

	return false
}

// matchWordBoundary reports whether s contains pattern followed by a
// non-alphanumeric character (or end of string). This prevents false
// positives like "bot" matching inside "RoboBotonic" or "iBot".
func matchWordBoundary(s, pattern string) bool {
	idx := 0
	for {
		pos := strings.Index(s[idx:], pattern)
		if pos < 0 {
			return false
		}
		pos += idx
		// Check the character after the match.
		end := pos + len(pattern)
		if end >= len(s) {
			return true // end of string is a valid boundary
		}
		next := s[end]
		if !isAlnum(next) {
			return true
		}
		// Not a boundary — keep searching after this position.
		idx = pos + 1
	}
}

// isAlnum reports whether b is an ASCII letter or digit.
func isAlnum(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// IsReferrerSpam checks whether a referrer URL host matches a known spam domain.
// The referrer can be a full URL or just a domain. Matching is case-insensitive.
func IsReferrerSpam(referrer string) bool {
	if referrer == "" {
		return false
	}

	host := referrer

	// Try to parse as URL to extract host
	if u, err := url.Parse(referrer); err == nil && u.Host != "" {
		host = u.Host
	}

	host = strings.ToLower(strings.TrimSpace(host))
	// Strip port if present (handle IPv6 brackets)
	if i := strings.LastIndex(host, "]"); i >= 0 {
		host = host[:i+1]
	} else if idx := strings.LastIndex(host, ":"); idx > 0 {
		host = host[:idx]
	}
	// Strip leading "www." for matching
	host = strings.TrimPrefix(host, "www.")

	if host == "" {
		return false
	}

	// Check exact match
	if referrerSpamDomains[host] {
		return true
	}

	// Check if any spam domain is a suffix (handles subdomains)
	for domain := range referrerSpamDomains {
		if strings.HasSuffix(host, "."+domain) {
			return true
		}
	}

	return false
}

// IsDataCenterIP checks whether an IP address falls within known data center
// CIDR ranges (AWS, GCP, Azure, DigitalOcean, Oracle Cloud).
func IsDataCenterIP(ip string) bool {
	if ip == "" {
		return false
	}

	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}

	for _, prefix := range parsedDataCenterCIDRs {
		if prefix.Contains(addr) {
			return true
		}
	}

	return false
}
