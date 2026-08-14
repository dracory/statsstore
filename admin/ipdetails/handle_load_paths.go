package ipdetails

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"strings"

	"github.com/dracory/statsstore/admin/shared"

	"github.com/dracory/api"
	"github.com/dracory/req"
	"github.com/dracory/statsstore"
)

func (controller *ipDetailsController) handleLoadPaths(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	store := controller.opts.Store
	if store == nil {
		api.Respond(w, r, api.Error("Stats store not available"))
		return ""
	}

	ip := req.GetStringTrimmed(r, shared.ParamIP)
	if ip == "" {
		// Try JSON body if query param wasn't set.
		var body struct {
			IP string `json:"ip"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		ip = strings.TrimSpace(body.IP)
	}
	if ip == "" {
		api.Respond(w, r, api.Error("IP address is required"))
		return ""
	}

	ctx := r.Context()

	// Fetch all visitors for this IP. The store query doesn't support
	// filtering by IP directly, so we fetch a large batch and filter
	// in-memory. PathContains with empty string matches all paths.
	visitors, err := store.VisitorList(ctx, statsstore.VisitorQuery().
		SetOrderBy("created_at").
		SetSortOrder("DESC").
		SetLimit(10000))
	if err != nil {
		slog.Error("ipdetails: failed to load visitors", "error", err)
		api.Respond(w, r, api.Error("Failed to load visitor data"))
		return ""
	}

	// Filter to just this IP.
	var ipVisitors []statsstore.VisitorInterface
	for _, v := range visitors {
		if v.GetIpAddress() == ip {
			ipVisitors = append(ipVisitors, v)
		}
	}

	if len(ipVisitors) == 0 {
		api.Respond(w, r, api.SuccessWithData("No visits found", map[string]any{
			FieldDetails:    map[string]any{},
			FieldPaths:      []any{},
			FieldVisitCount: 0,
			FieldIsBot:      false,
			FieldIsThreat:   false,
		}))
		return ""
	}

	// Build details from the most recent record (DESC sort).
	latest := ipVisitors[0]
	first := ipVisitors[len(ipVisitors)-1]

	// Build path history (already DESC by created_at).
	paths := make([]map[string]any, 0, len(ipVisitors))
	for _, v := range ipVisitors {
		paths = append(paths, map[string]any{
			FieldPath:          v.GetPath(),
			FieldPathCreatedAt: v.GetCreatedAt(),
		})
	}

	// Unique path count.
	uniquePaths := map[string]bool{}
	for _, v := range ipVisitors {
		uniquePaths[v.GetPath()] = true
	}

	details := map[string]any{
		FieldIP:          ip,
		FieldCountry:     latest.GetCountry(),
		FieldBrowser:     latest.GetUserBrowser(),
		FieldBrowserVer:  latest.GetUserBrowserVersion(),
		FieldOS:          latest.GetUserOs(),
		FieldOSVer:       latest.GetUserOsVersion(),
		FieldDevice:      latest.GetUserDevice(),
		FieldDeviceType:  latest.GetUserDeviceType(),
		FieldUserAgent:   latest.GetUserAgent(),
		FieldAcceptLang:  latest.GetUserAcceptLanguage(),
		FieldAcceptEnc:   latest.GetUserAcceptEncoding(),
		FieldReferrer:    latest.GetUserReferrer(),
		FieldFingerprint: latest.GetFingerprint(),
		FieldFirstSeen:   first.GetCreatedAt(),
		FieldLastSeen:    latest.GetCreatedAt(),
	}

	// Sort paths by created_at descending (already sorted, but ensure).
	sort.Slice(paths, func(i, j int) bool {
		ci, _ := paths[i][FieldPathCreatedAt].(string)
		cj, _ := paths[j][FieldPathCreatedAt].(string)
		return ci > cj
	})

	api.Respond(w, r, api.SuccessWithData("IP details loaded", map[string]any{
		FieldDetails:    details,
		FieldPaths:      paths,
		FieldVisitCount: len(ipVisitors),
		FieldIsBot:      latest.GetBot() == statsstore.VALUE_YES,
		FieldIsThreat:   latest.GetThreat() == statsstore.VALUE_YES,
		FieldTotal:      len(paths),
	}))
	return ""
}
