package statsstore

import (
	"context"
	"errors"

	"github.com/dromara/carbon/v2"
)

// SettingGet retrieves a setting value by key from the settings table.
// Returns an empty string and nil error if the key does not exist.
func (st *storeImplementation) SettingGet(ctx context.Context, key string) (string, error) {
	if st.settingsTableName == "" {
		return "", errors.New("settings table name is empty")
	}

	if key == "" {
		return "", errors.New("key is empty")
	}

	if !st.db.Schema().HasTable(st.settingsTableName) {
		return "", nil
	}

	var rows []settingsRow
	if err := st.db.Query().
		Table(st.settingsTableName).
		Where(COLUMN_KEY+" = ?", key).
		Get(&rows); err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "", nil
	}

	return rows[0].Value, nil
}

// SettingSet stores a setting value by key in the settings table, using
// upsert semantics — if the key exists, the value is updated; otherwise
// a new row is inserted.
func (st *storeImplementation) SettingSet(ctx context.Context, key, value string) error {
	if st.settingsTableName == "" {
		return errors.New("settings table name is empty")
	}

	if key == "" {
		return errors.New("key is empty")
	}

	if !st.db.Schema().HasTable(st.settingsTableName) {
		return errors.New("settings table does not exist")
	}

	now := carbon.Now(carbon.UTC).StdTime()

	// Check if the setting already exists.
	var existing []settingsRow
	if err := st.db.Query().
		Table(st.settingsTableName).
		Where(COLUMN_KEY+" = ?", key).
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
			Where(COLUMN_KEY+" = ?", key).
			Update(row)
		return err
	}

	row := map[string]any{
		COLUMN_KEY:        key,
		COLUMN_VALUE:      value,
		COLUMN_CREATED_AT: now,
		COLUMN_UPDATED_AT: now,
	}

	return st.db.Query().Table(st.settingsTableName).Create(row)
}

// SettingDelete removes a setting row by key. No error is returned if the
// key does not exist.
func (st *storeImplementation) SettingDelete(ctx context.Context, key string) error {
	if st.settingsTableName == "" {
		return errors.New("settings table name is empty")
	}

	if key == "" {
		return errors.New("key is empty")
	}

	if !st.db.Schema().HasTable(st.settingsTableName) {
		return nil
	}

	_, err := st.db.Query().
		Table(st.settingsTableName).
		Where(COLUMN_KEY+" = ?", key).
		Delete()
	return err
}

// SettingHas reports whether a setting key exists in the settings table.
func (st *storeImplementation) SettingHas(ctx context.Context, key string) (bool, error) {
	if st.settingsTableName == "" {
		return false, errors.New("settings table name is empty")
	}

	if key == "" {
		return false, errors.New("key is empty")
	}

	if !st.db.Schema().HasTable(st.settingsTableName) {
		return false, nil
	}

	var rows []settingsRow
	if err := st.db.Query().
		Table(st.settingsTableName).
		Where(COLUMN_KEY+" = ?", key).
		Get(&rows); err != nil {
		return false, err
	}

	return len(rows) > 0, nil
}

// SettingList returns all settings as a map of key to value. Returns an
// empty map if the settings table does not exist or is empty.
func (st *storeImplementation) SettingList(ctx context.Context) (map[string]string, error) {
	if st.settingsTableName == "" {
		return nil, errors.New("settings table name is empty")
	}

	if !st.db.Schema().HasTable(st.settingsTableName) {
		return map[string]string{}, nil
	}

	var rows []settingsRow
	if err := st.db.Query().
		Table(st.settingsTableName).
		Get(&rows); err != nil {
		return nil, err
	}

	result := make(map[string]string, len(rows))
	for _, r := range rows {
		result[r.Key] = r.Value
	}

	return result, nil
}
