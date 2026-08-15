package shared

import (
	"sort"
	"strings"

	"github.com/dracory/statsstore"
)

// BotReasonDisplay is a single reason an IP was flagged as a bot, formatted
// for display in the admin UI. It mirrors the JSON shape expected by the
// frontend (pattern, hits, sample paths).
type BotReasonDisplay struct {
	Pattern string   `json:"pattern"`
	Hits    int      `json:"hits"`
	Paths   []string `json:"paths"`
}

// Reason pattern labels for signals that are not path substrings. These are
// shown in the UI as the reason "pattern" so the user understands which
// heuristic triggered the bot flag.
const (
	// BotReasonUserAgent indicates the visitor's User-Agent string
	// self-identified as a bot (matched statsstore.IsBot).
	BotReasonUserAgent = "user-agent:bot"

	// BotReasonDataCenterIP indicates the visitor's IP falls within a known
	// data-center CIDR range (matched statsstore.IsDataCenterIP).
	BotReasonDataCenterIP = "datacenter-ip"

	// BotReasonReferrerSpam indicates the visitor's referrer matched a known
	// referrer-spam domain (matched statsstore.IsReferrerSpam).
	BotReasonReferrerSpam = "referrer-spam"
)

// ComputeBotReasons derives the bot reasons for a set of visitor records
// belonging to a single IP. It checks all signals used at ingestion time
// (see store.VisitorRegister): user-agent self-identification, data-center
// IP, referrer spam, bot-only paths, and malicious paths. extraPathPatterns
// is an optional list of consumer-specific path substrings to also check
// (e.g. ".php", "/wp-admin"); pass nil to check only the universal
// statsstore patterns.
//
// This is the single source of truth for why an IP was flagged, ensuring
// the displayed reasons match the ingestion logic. The returned slice is
// sorted by pattern for stable display.
func ComputeBotReasons(visitors []statsstore.VisitorInterface, extraPathPatterns []string) []BotReasonDisplay {
	type reasonAgg struct {
		Pattern string
		Hits    int
		Paths   []string
		seen    map[string]bool
	}

	reasons := map[string]*reasonAgg{}

	addReason := func(pattern, path string) {
		r, ok := reasons[pattern]
		if !ok {
			r = &reasonAgg{Pattern: pattern, seen: map[string]bool{}}
			reasons[pattern] = r
		}
		r.Hits++
		if path != "" && !r.seen[path] && len(r.Paths) < 3 {
			r.Paths = append(r.Paths, path)
			r.seen[path] = true
		}
	}

	// Combine universal bot path patterns, universal malicious path patterns,
	// and any consumer-specific extra patterns into one list for substring
	// matching. Matching is case-insensitive against the lowercased path,
	// consistent with the existing handleLoadBots logic.
	allPathPatterns := append([]string{}, statsstore.BotPathPatterns()...)
	allPathPatterns = append(allPathPatterns, statsstore.MaliciousPathPatterns()...)
	allPathPatterns = append(allPathPatterns, extraPathPatterns...)

	// Pre-lowercase the patterns once for case-insensitive matching.
	pathPatternsLower := make([]string, len(allPathPatterns))
	for i, p := range allPathPatterns {
		pathPatternsLower[i] = strings.ToLower(p)
	}

	// dataCenterIP is the same for all visitors from the same IP, so only
	// record it once. We still iterate all visitors for the other signals.
	dataCenterChecked := false

	for _, v := range visitors {
		ua := v.GetUserAgent()
		ip := v.GetIpAddress()
		referrer := v.GetUserReferrer()
		path := v.GetPath()
		pathLower := strings.ToLower(path)

		// User-agent self-identification (e.g. "Googlebot", "curl").
		if ua != "" && statsstore.IsBot(ua) {
			addReason(BotReasonUserAgent, path)
		}

		// Data-center IP — check once per IP since it's IP-level, not per-visit.
		if !dataCenterChecked && ip != "" && statsstore.IsDataCenterIP(ip) {
			addReason(BotReasonDataCenterIP, path)
			dataCenterChecked = true
		}

		// Referrer spam.
		if referrer != "" && statsstore.IsReferrerSpam(referrer) {
			addReason(BotReasonReferrerSpam, path)
		}

		// Path-based patterns (bot files + malicious paths + consumer extras).
		for i, patternLower := range pathPatternsLower {
			if strings.Contains(pathLower, patternLower) {
				addReason(allPathPatterns[i], path)
			}
		}
	}

	list := make([]BotReasonDisplay, 0, len(reasons))
	for _, r := range reasons {
		if r.Paths == nil {
			r.Paths = []string{}
		}
		list = append(list, BotReasonDisplay{
			Pattern: r.Pattern,
			Hits:    r.Hits,
			Paths:   r.Paths,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Pattern < list[j].Pattern
	})

	return list
}
