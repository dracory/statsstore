package statsstore

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
)

// settingsRow is the DB row mapping for the generic settings table.
type settingsRow struct {
	Key       string    `db:"key"`
	Value     string    `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ExcludedIPList returns the effective excluded IP addresses — the merge of
// IPs stored in the database and any supplied via NewStoreOptions.
func (st *storeImplementation) ExcludedIPList(ctx context.Context) ([]string, error) {
	return st.GetExcludedIPs(), nil
}

// ExcludedIPAdd adds an IP address to the exclusion list in the database and
// updates the in-memory cache. Duplicate IPs are silently ignored.
// IPs supplied via NewStoreOptions are never written to the database.
func (st *storeImplementation) ExcludedIPAdd(ctx context.Context, ip string) error {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return errors.New("ip address is empty")
	}

	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()

	if slices.Contains(st.excludedIPs, ip) {
		return nil
	}

	dbList := append(stringListSubtract(st.excludedIPs, st.excludedIPsOpts), ip)
	if err := st.excludedIPsSaveToDB(ctx, dbList); err != nil {
		return err
	}

	st.excludedIPs = append(slices.Clone(st.excludedIPsOpts), dbList...)
	return nil
}

// ExcludedIPRemove removes an IP address from the exclusion list in the database
// and updates the in-memory cache. IPs supplied via NewStoreOptions cannot be
// removed this way — they stay active for the lifetime of the process.
func (st *storeImplementation) ExcludedIPRemove(ctx context.Context, ip string) error {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return errors.New("ip address is empty")
	}

	st.filtersMu.Lock()
	defer st.filtersMu.Unlock()

	dbList := stringListRemove(
		stringListSubtract(st.excludedIPs, st.excludedIPsOpts), ip)
	if err := st.excludedIPsSaveToDB(ctx, dbList); err != nil {
		return err
	}

	st.excludedIPs = append(slices.Clone(st.excludedIPsOpts), dbList...)
	return nil
}

// excludedIPsLoadFromDB reads the excluded IPs JSON array from the settings table.
func (st *storeImplementation) excludedIPsLoadFromDB(ctx context.Context) ([]string, error) {
	if st.settingsTableName == "" {
		return nil, nil
	}

	if !st.db.Schema().HasTable(st.settingsTableName) {
		return nil, nil
	}

	var rows []settingsRow
	if err := st.db.Query().
		Table(st.settingsTableName).
		Where(COLUMN_KEY+" = ?", SETTING_EXCLUDED_IPS).
		Get(&rows); err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	var ips []string
	if err := json.Unmarshal([]byte(rows[0].Value), &ips); err != nil {
		return nil, err
	}

	return ips, nil
}

// excludedIPsSaveToDB serializes the given excluded IPs list to the settings
// table as a JSON array, using upsert semantics.
func (st *storeImplementation) excludedIPsSaveToDB(ctx context.Context, ips []string) error {
	if st.settingsTableName == "" {
		return errors.New("settings table name is empty")
	}

	data, err := json.Marshal(ips)
	if err != nil {
		return err
	}

	now := carbon.Now(carbon.UTC).StdTime()
	value := string(data)

	// Check if the setting already exists
	var existing []settingsRow
	if err := st.db.Query().
		Table(st.settingsTableName).
		Where(COLUMN_KEY+" = ?", SETTING_EXCLUDED_IPS).
		Get(&existing); err != nil {
		return err
	}

	if len(existing) > 0 {
		row := map[string]any{
			COLUMN_VALUE:      value,
			COLUMN_UPDATED_AT: now,
		}
		_, err := st.db.Query().
			Table(st.settingsTableName).
			Where(COLUMN_KEY+" = ?", SETTING_EXCLUDED_IPS).
			Update(row)
		return err
	}

	row := map[string]any{
		COLUMN_KEY:        SETTING_EXCLUDED_IPS,
		COLUMN_VALUE:      value,
		COLUMN_CREATED_AT: now,
		COLUMN_UPDATED_AT: now,
	}

	return st.db.Query().Table(st.settingsTableName).Create(row)
}
