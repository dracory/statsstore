package statsstore

import (
	"context"
	"errors"
	"strconv"
	"time"

	"uuid"
)

// sessionRow is the scan target for the session aggregate query.
// started_at/ended_at are scanned as `any` because MIN()/MAX() aggregates do
// not carry a declared column type, so drivers may return time.Time, string,
// or []byte depending on the dialect.
type sessionRow struct {
	SessionID      string `db:"session_id"`
	Fingerprint    string `db:"fingerprint"`
	IPAddress      string `db:"ip_address"`
	Country        string `db:"country"`
	UserDeviceType string `db:"user_device_type"`
	UserBrowser    string `db:"user_browser"`
	UserOs         string `db:"user_os"`
	StartedAt      any    `db:"started_at"`
	EndedAt        any    `db:"ended_at"`
	PageViews      int64  `db:"page_views"`
	BotRows        int64  `db:"bot_rows"`
}

// sessionBackfillRow is the scan target for SessionIDBackfill. CreatedAt is
// `any` for the same reason as sessionRow.StartedAt/EndedAt.
type sessionBackfillRow struct {
	ID          string `db:"id"`
	Fingerprint string `db:"fingerprint"`
	CreatedAt   any    `db:"created_at"`
}

// sessionPathRow is the scan target for the per-session path listing used to
// derive entry and exit paths.
type sessionPathRow struct {
	SessionID string `db:"session_id"`
	Path      string `db:"path"`
	CreatedAt any    `db:"created_at"`
}

// SessionList aggregates visitor rows into sessions. Sessions are groups of
// page views sharing a session_id (assigned at registration: same fingerprint,
// consecutive views within SESSION_GAP_MINUTES).
//
// All VisitorQuery filters (country, device type, path, date range, bot,
// threat, IP, ...) apply to the underlying rows before grouping. Limit and
// offset paginate sessions (ordered by last activity, newest first), not rows.
// Entry and exit paths are the first and last path of each returned session
// within the filtered range. Rows predating the session_id column (empty
// session id) are excluded.
func (st *storeImplementation) SessionList(ctx context.Context, query VisitorQueryInterface) ([]SessionInterface, error) {
	if query == nil {
		return nil, errors.New("session list: query is nil")
	}

	q := st.db.Query().Model(&visitorImplementation{}).Table(st.visitorTableName)
	q = st.applyVisitorFilters(q, query)
	q = q.Where(COLUMN_SESSION_ID+" != ?", "")

	q = q.Select(COLUMN_SESSION_ID+" AS session_id"+
		", MAX("+COLUMN_FINGERPRINT+") AS fingerprint"+
		", MAX("+COLUMN_IP_ADDRESS+") AS ip_address"+
		", MAX("+COLUMN_COUNTRY+") AS country"+
		", MAX("+COLUMN_USER_DEVICE_TYPE+") AS user_device_type"+
		", MAX("+COLUMN_USER_BROWSER+") AS user_browser"+
		", MAX("+COLUMN_USER_OS+") AS user_os"+
		", MIN("+COLUMN_CREATED_AT+") AS started_at"+
		", MAX("+COLUMN_CREATED_AT+") AS ended_at"+
		", COUNT(*) AS page_views"+
		", SUM(CASE WHEN "+COLUMN_BOT+" = '"+VALUE_YES+"' THEN 1 ELSE 0 END) AS bot_rows").
		Group(COLUMN_SESSION_ID).
		OrderBy("ended_at", "desc")

	if query.HasLimit() && query.Limit() > 0 {
		q = q.Limit(query.Limit())
	}
	if query.HasOffset() && query.Offset() > 0 {
		q = q.Offset(query.Offset())
	}

	var rows []sessionRow
	if err := q.Get(&rows); err != nil {
		return []SessionInterface{}, err
	}

	sessions := make([]SessionInterface, 0, len(rows))
	sessionIDs := make([]any, 0, len(rows))
	byID := make(map[string]*sessionImplementation, len(rows))

	for _, r := range rows {
		s := &sessionImplementation{
			sessionID:      r.SessionID,
			fingerprint:    r.Fingerprint,
			ipAddress:      r.IPAddress,
			country:        r.Country,
			userDeviceType: r.UserDeviceType,
			userBrowser:    r.UserBrowser,
			userOs:         r.UserOs,
			startedAt:      sessionTime(r.StartedAt),
			endedAt:        sessionTime(r.EndedAt),
			pageViews:      r.PageViews,
			botRows:        r.BotRows,
		}
		sessions = append(sessions, s)
		sessionIDs = append(sessionIDs, r.SessionID)
		byID[r.SessionID] = s
	}

	if err := st.fillSessionEntryExitPaths(sessionIDs, byID); err != nil {
		return sessions, err
	}

	return sessions, nil
}

// SessionCount counts distinct session ids matching the visitor query
// filters. Rows with an empty session id are excluded.
func (st *storeImplementation) SessionCount(ctx context.Context, query VisitorQueryInterface) (int64, error) {
	if query == nil {
		return 0, errors.New("session count: query is nil")
	}

	q := st.db.Query().Model(&visitorImplementation{}).Table(st.visitorTableName)
	q = st.applyVisitorFilters(q, query)
	q = q.Where(COLUMN_SESSION_ID+" != ?", "")

	var results []map[string]any
	err := q.Select("COUNT(DISTINCT " + COLUMN_SESSION_ID + ") AS n").Get(&results)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}

	return toInt64(results[0]["n"]), nil
}

// SessionIDBackfill assigns session ids to rows that predate the session_id
// column (session_id = ”), processing them oldest-first so sessions chain
// correctly across the SESSION_GAP_MINUTES window. Returns the number of
// updated rows. Intended for one-off upgrades; batch with a date-range
// VisitorQuery on very large tables.
func (st *storeImplementation) SessionIDBackfill(ctx context.Context, query VisitorQueryInterface) (int64, error) {
	if query == nil {
		return 0, errors.New("session backfill: query is nil")
	}

	q := st.db.Query().Model(&visitorImplementation{}).Table(st.visitorTableName)
	q = st.applyVisitorFilters(q, query)
	q = q.Where("("+COLUMN_SESSION_ID+" = '' OR "+COLUMN_SESSION_ID+" IS NULL)").OrderBy(COLUMN_CREATED_AT, "asc")

	var rows []sessionBackfillRow
	if err := q.Select(COLUMN_ID + " AS id, " + COLUMN_FINGERPRINT + " AS fingerprint, " + COLUMN_CREATED_AT + " AS created_at").Get(&rows); err != nil {
		return 0, err
	}

	var updated int64
	// lastSeen[fingerprint] = created_at of the previous row in that session
	lastSeen := map[string]time.Time{}
	lastSession := map[string]string{}

	for _, r := range rows {
		createdAt := sessionTime(r.CreatedAt)
		if createdAt.IsZero() {
			continue
		}

		sessionID := ""
		if prev, exists := lastSeen[r.Fingerprint]; exists && createdAt.Sub(prev) <= time.Duration(SESSION_GAP_MINUTES)*time.Minute {
			sessionID = lastSession[r.Fingerprint]
		} else {
			sessionID = uuid.NewV7().String()
		}
		lastSeen[r.Fingerprint] = createdAt
		lastSession[r.Fingerprint] = sessionID

		if _, err := st.db.Query().Table(st.visitorTableName).Where(COLUMN_ID+" = ?", r.ID).Update(map[string]any{
			COLUMN_SESSION_ID: sessionID,
		}); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

// fillSessionEntryExitPaths loads all paths of the given sessions ordered by
// created_at and records the first path as the session's entry and the last
// as its exit.
func (st *storeImplementation) fillSessionEntryExitPaths(sessionIDs []any, byID map[string]*sessionImplementation) error {
	if len(sessionIDs) == 0 {
		return nil
	}

	var rows []sessionPathRow
	err := st.db.Query().Model(&visitorImplementation{}).Table(st.visitorTableName).
		Select(COLUMN_SESSION_ID+" AS session_id, "+COLUMN_PATH+" AS path, "+COLUMN_CREATED_AT+" AS created_at").
		WhereIn(COLUMN_SESSION_ID, sessionIDs).
		OrderBy(COLUMN_CREATED_AT, "asc").
		Get(&rows)
	if err != nil {
		return err
	}

	for _, r := range rows {
		s, ok := byID[r.SessionID]
		if !ok {
			continue
		}
		if s.entryPath == "" {
			s.entryPath = r.Path
		}
		s.exitPath = r.Path
	}

	return nil
}

// sessionTime converts a scanned MIN()/MAX() aggregate value into time.Time.
// SQLite stores datetimes as TEXT, so aggregate results arrive as strings
// rather than time.Time like declared DATETIME columns do.
func sessionTime(v any) time.Time {
	switch t := v.(type) {
	case time.Time:
		return t.UTC()
	case string:
		if parsed, ok := parseCreatedAt(t); ok {
			return parsed
		}
	case []byte:
		if parsed, ok := parseCreatedAt(string(t)); ok {
			return parsed
		}
	}
	return time.Time{}
}

// toInt64 converts a scanned numeric aggregate value into int64.
func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	case []byte:
		if parsed, err := strconv.ParseInt(string(n), 10, 64); err == nil {
			return parsed
		}
	case string:
		if parsed, err := strconv.ParseInt(n, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}
