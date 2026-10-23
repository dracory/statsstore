package statsstore

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
)

// ExcludedPathList returns the effective excluded path prefixes — the merge of
// prefixes stored in the database and any supplied via NewStoreOptions.
func (st *storeImplementation) ExcludedPathList(ctx context.Context) ([]string, error) {
	return st.GetExcludedPathPrefixes(), nil
}

// ExcludedPathAdd adds a path prefix to the exclusion list in the database and
// updates the in-memory cache. Duplicate prefixes are silently ignored.
// Prefixes supplied via NewStoreOptions are never written to the database.
func (st *storeImplementation) ExcludedPathAdd(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()

	if slices.Contains(st.excludedPathPrefixes, prefix) {
		return nil
	}

	dbList := append(stringListSubtract(st.excludedPathPrefixes, st.excludedPathPrefixesOpts), prefix)
	if err := st.pathListSaveToDB(ctx, SETTING_EXCLUDED_PATHS, dbList); err != nil {
		return err
	}

	st.excludedPathPrefixes = append(slices.Clone(st.excludedPathPrefixesOpts), dbList...)
	return nil
}

// ExcludedPathRemove removes a path prefix from the exclusion list in the
// database and updates the in-memory cache. Prefixes supplied via
// NewStoreOptions cannot be removed this way — they stay active for the
// lifetime of the process.
func (st *storeImplementation) ExcludedPathRemove(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()

	dbList := stringListRemove(
		stringListSubtract(st.excludedPathPrefixes, st.excludedPathPrefixesOpts), prefix)
	if err := st.pathListSaveToDB(ctx, SETTING_EXCLUDED_PATHS, dbList); err != nil {
		return err
	}

	st.excludedPathPrefixes = append(slices.Clone(st.excludedPathPrefixesOpts), dbList...)
	return nil
}

// IncludedPathList returns the effective included path prefixes — the merge of
// prefixes stored in the database and any supplied via NewStoreOptions.
// When non-empty the list acts as a whitelist — only matching paths are tracked.
func (st *storeImplementation) IncludedPathList(ctx context.Context) ([]string, error) {
	return st.GetIncludedPathPrefixes(), nil
}

// IncludedPathAdd adds a path prefix to the inclusion list in the database and
// updates the in-memory cache. Duplicate prefixes are silently ignored.
// Prefixes supplied via NewStoreOptions are never written to the database.
func (st *storeImplementation) IncludedPathAdd(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()

	if slices.Contains(st.includedPathPrefixes, prefix) {
		return nil
	}

	dbList := append(stringListSubtract(st.includedPathPrefixes, st.includedPathPrefixesOpts), prefix)
	if err := st.pathListSaveToDB(ctx, SETTING_INCLUDED_PATHS, dbList); err != nil {
		return err
	}

	st.includedPathPrefixes = append(slices.Clone(st.includedPathPrefixesOpts), dbList...)
	return nil
}

// IncludedPathRemove removes a path prefix from the inclusion list in the
// database and updates the in-memory cache. Prefixes supplied via
// NewStoreOptions cannot be removed this way — they stay active for the
// lifetime of the process.
func (st *storeImplementation) IncludedPathRemove(ctx context.Context, prefix string) error {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return errors.New("path prefix is empty")
	}

	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()

	dbList := stringListRemove(
		stringListSubtract(st.includedPathPrefixes, st.includedPathPrefixesOpts), prefix)
	if err := st.pathListSaveToDB(ctx, SETTING_INCLUDED_PATHS, dbList); err != nil {
		return err
	}

	st.includedPathPrefixes = append(slices.Clone(st.includedPathPrefixesOpts), dbList...)
	return nil
}

// stringListSubtract returns the elements of list that are not present in exclude.
func stringListSubtract(list, exclude []string) []string {
	out := make([]string, 0, len(list))
	for _, s := range list {
		if !slices.Contains(exclude, s) {
			out = append(out, s)
		}
	}
	return out
}

// stringListRemove returns a copy of list with all occurrences of s removed.
func stringListRemove(list []string, s string) []string {
	out := make([]string, 0, len(list))
	for _, existing := range list {
		if existing != s {
			out = append(out, existing)
		}
	}
	return out
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
