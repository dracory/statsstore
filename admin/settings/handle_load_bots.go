package settings

import (
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/statsstore"
)

// reasonAgg tracks per-pattern evidence for a bot IP during aggregation.
type reasonAgg struct {
	Pattern string
	Hits    int
	Paths   []string
}

func (controller *settingsController) handleLoadBots(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	store := controller.opts.Store
	if store == nil {
		api.Respond(w, r, api.Error("Stats store not available"))
		return ""
	}

	visitors, err := store.VisitorList(r.Context(), statsstore.VisitorQuery().
		SetBot(statsstore.VALUE_YES).
		SetLimit(100000))
	if err != nil {
		api.Respond(w, r, api.Error("Failed to load bot visitors"))
		return ""
	}

	// Aggregate by IP: hit count, last seen, and reasons (pattern matches).
	type ipAgg struct {
		IP       string
		Hits     int
		LastSeen string
		Reasons  map[string]*reasonAgg
	}
	agg := map[string]*ipAgg{}
	for _, v := range visitors {
		ip := v.GetIpAddress()
		if ip == "" {
			continue
		}
		a, ok := agg[ip]
		if !ok {
			a = &ipAgg{IP: ip, Reasons: map[string]*reasonAgg{}}
			agg[ip] = a
		}
		a.Hits++
		created := v.GetCreatedAt()
		if created > a.LastSeen {
			a.LastSeen = created
		}

		// Check visitor path against bot/malicious patterns to derive reasons.
		path := v.GetPath()
		for _, pattern := range botPagePatterns {
			if strings.Contains(path, pattern) {
				reason, ok := a.Reasons[pattern]
				if !ok {
					reason = &reasonAgg{Pattern: pattern}
					a.Reasons[pattern] = reason
				}
				reason.Hits++
				if len(reason.Paths) < 3 && !slices.Contains(reason.Paths, path) {
					reason.Paths = append(reason.Paths, path)
				}
			}
		}

		// Check user agent for self-identifying bots.
		ua := v.GetUserAgent()
		if ua != "" && statsstore.IsBot(ua) {
			reason, ok := a.Reasons[userAgentPattern]
			if !ok {
				reason = &reasonAgg{Pattern: userAgentPattern}
				a.Reasons[userAgentPattern] = reason
			}
			reason.Hits++
		}
	}

	list := make([]map[string]any, 0, len(agg))
	for _, a := range agg {
		list = append(list, map[string]any{
			FieldIP:       a.IP,
			FieldHits:     a.Hits,
			FieldLastSeen: a.LastSeen,
			FieldReasons:  reasonsToList(a.Reasons),
		})
	}
	sort.Slice(list, func(i, j int) bool {
		hi, _ := list[i][FieldHits].(int)
		hj, _ := list[j][FieldHits].(int)
		if hi != hj {
			return hi > hj
		}
		pi, _ := list[i][FieldIP].(string)
		pj, _ := list[j][FieldIP].(string)
		return pi < pj
	})

	api.Respond(w, r, api.SuccessWithData("Bot IPs loaded", map[string]any{
		FieldBots:  list,
		FieldTotal: len(list),
	}))
	return ""
}

// reasonsToList converts the reason aggregation map into a sorted slice for
// the API response, matching the format expected by the frontend.
func reasonsToList(reasons map[string]*reasonAgg) []map[string]any {
	if len(reasons) == 0 {
		return []map[string]any{}
	}
	out := make([]map[string]any, 0, len(reasons))
	for _, r := range reasons {
		out = append(out, map[string]any{
			"pattern": r.Pattern,
			"hits":    r.Hits,
			"paths":   r.Paths,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		pi, _ := out[i]["pattern"].(string)
		pj, _ := out[j]["pattern"].(string)
		return pi < pj
	})
	return out
}
