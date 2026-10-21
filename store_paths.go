package statsstore

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

// ExcludedPathList returns all excluded path prefixes stored in the database.
func (st *storeImplementation) ExcludedPathList(ctx context.Context) ([]string, error) {
	return st.pathListLoadFromDB(ctx, SETTING_EXCLUDED_PATHS)
}

// ExcludedPathAdd adds a path prefix to the exclusion list in the database and
// updates the in-memory cache. Duplicate prefixes are silently ignored.
func (st *storeImplementation) ExcludedPathAdd(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	for _, existing := range st.excludedPathPrefixes {
		if existing == prefix {
			return nil
		}
	}

	st.excludedPathPrefixes = append(st.excludedPathPrefixes, prefix)
	return st.pathListSaveToDB(ctx, SETTING_EXCLUDED_PATHS, st.excludedPathPrefixes)
}

// ExcludedPathRemove removes a path prefix from the exclusion list in the
// database and updates the in-memory cache.
func (st *storeImplementation) ExcludedPathRemove(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	filtered := make([]string, 0, len(st.excludedPathPrefixes))
	for _, existing := range st.excludedPathPrefixes {
		if existing != prefix {
			filtered = append(filtered, existing)
		}
	}
	st.excludedPathPrefixes = filtered

	return st.pathListSaveToDB(ctx, SETTING_EXCLUDED_PATHS, st.excludedPathPrefixes)
}

// IncludedPathList returns all included path prefixes stored in the database.
// When non-empty the list acts as a whitelist — only matching paths are tracked.
func (st *storeImplementation) IncludedPathList(ctx context.Context) ([]string, error) {
	return st.pathListLoadFromDB(ctx, SETTING_INCLUDED_PATHS)
}

// IncludedPathAdd adds a path prefix to the inclusion list in the database and
// updates the in-memory cache. Duplicate prefixes are silently ignored.
func (st *storeImplementation) IncludedPathAdd(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	for _, existing := range st.includedPathPrefixes {
		if existing == prefix {
			return nil
		}
	}

	st.includedPathPrefixes = append(st.includedPathPrefixes, prefix)
	return st.pathListSaveToDB(ctx, SETTING_INCLUDED_PATHS, st.includedPathPrefixes)
}

// IncludedPathRemove removes a path prefix from the inclusion list in the
// database and updates the in-memory cache.
func (st *storeImplementation) IncludedPathRemove(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	filtered := make([]string, 0, len(st.includedPathPrefixes))
	for _, existing := range st.includedPathPrefixes {
		if existing != prefix {
			filtered = append(filtered, existing)
		}
	}
	st.includedPathPrefixes = filtered

	return st.pathListSaveToDB(ctx, SETTING_INCLUDED_PATHS, st.includedPathPrefixes)
}

// pathPatternMatchAny reports whether s matches any of the given patterns and
// returns the first matching pattern (for logging purposes).
func pathPatternMatchAny(patterns []string, s string) (string, bool) {
	for _, pattern := range patterns {
		if pathPatternMatch(pattern, s) {
			return pattern, true
		}
	}
	return "", false
}

// pathPatternMatch matches a path pattern against s. Patterns without a "*"
// act as prefix matches (e.g. "/admin/" covers "/admin/users"). "*" matches
// any sequence of characters, including "/" — so "/api/*", "/*.php" and
// "/files/*/download" all work.
func pathPatternMatch(pattern, s string) bool {
	if !strings.Contains(pattern, "*") {
		return strings.HasPrefix(s, pattern)
	}

	parts := strings.Split(pattern, "*")
	pos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}
		idx := strings.Index(s[pos:], part)
		if idx < 0 {
			return false
		}
		// The first literal segment must anchor at the start unless the
		// pattern begins with "*".
		if i == 0 && idx != 0 {
			return false
		}
		pos += idx + len(part)
	}
	// The last literal segment must anchor at the end unless the pattern
	// ends with "*".
	if !strings.HasSuffix(pattern, "*") && pos != len(s) {
		return false
	}
	return true
}

// pathListLoadFromDB reads a JSON array of path prefixes from the settings table.
func (st *storeImplementation) pathListLoadFromDB(ctx context.Context, key string) ([]string, error) {
	if st.settingsTableName == "" {
		return nil, nil
	}

	value, err := st.SettingGet(ctx, key)
	if err != nil {
		return nil, err
	}

	if value == "" {
		return nil, nil
	}

	var prefixes []string
	if err := json.Unmarshal([]byte(value), &prefixes); err != nil {
		return nil, err
	}

	return prefixes, nil
}

// pathListSaveToDB serializes a path prefix list to the settings table as a
// JSON array, using upsert semantics.
func (st *storeImplementation) pathListSaveToDB(ctx context.Context, key string, prefixes []string) error {
	data, err := json.Marshal(prefixes)
	if err != nil {
		return err
	}

	return st.SettingSet(ctx, key, string(data))
}
