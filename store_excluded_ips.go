package statsstore

import (
	"context"
	"encoding/json"
	"errors"
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

// ExcludedIPList returns all excluded IP addresses stored in the database.
func (st *storeImplementation) ExcludedIPList(ctx context.Context) ([]string, error) {
	return st.excludedIPsLoadFromDB(ctx)
}

// ExcludedIPAdd adds an IP address to the exclusion list in the database and
// updates the in-memory cache. Duplicate IPs are silently ignored.
func (st *storeImplementation) ExcludedIPAdd(ctx context.Context, ip string) error {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return errors.New("ip address is empty")
	}

	// Check if already exists in the in-memory list
	for _, existing := range st.excludedIPs {
		if existing == ip {
			return nil
		}
	}

	st.excludedIPs = append(st.excludedIPs, ip)
	return st.excludedIPsSaveToDB(ctx)
}

// ExcludedIPRemove removes an IP address from the exclusion list in the database
// and updates the in-memory cache.
func (st *storeImplementation) ExcludedIPRemove(ctx context.Context, ip string) error {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return errors.New("ip address is empty")
	}

	filtered := make([]string, 0, len(st.excludedIPs))
	for _, existing := range st.excludedIPs {
		if existing != ip {
			filtered = append(filtered, existing)
		}
	}
	st.excludedIPs = filtered

	return st.excludedIPsSaveToDB(ctx)
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

// excludedIPsSaveToDB serializes the in-memory excluded IPs list to the settings
// table as a JSON array, using upsert semantics.
func (st *storeImplementation) excludedIPsSaveToDB(ctx context.Context) error {
	if st.settingsTableName == "" {
		return errors.New("settings table name is empty")
	}

	data, err := json.Marshal(st.excludedIPs)
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
