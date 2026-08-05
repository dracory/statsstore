package statsstore

const (
	COLUMN_ID                   = "id"
	COLUMN_COUNTRY              = "country"
	COLUMN_CREATED_AT           = "created_at"
	COLUMN_SOFT_DELETED_AT      = "soft_deleted_at"
	COLUMN_IP_ADDRESS           = "ip_address"
	COLUMN_PATH                 = "path"
	COLUMN_UPDATED_AT           = "updated_at"
	COLUMN_FINGERPRINT          = "fingerprint"
	COLUMN_USER_AGENT           = "user_agent"
	COLUMN_USER_ACCEPT_LANGUAGE = "user_accept_language"
	COLUMN_USER_ACCEPT_ENCODING = "user_accept_encoding"
	COLUMN_USER_BROWSER         = "user_browser"
	COLUMN_USER_BROWSER_VERSION = "user_browser_version"
	COLUMN_USER_DEVICE          = "user_device"
	COLUMN_USER_DEVICE_TYPE     = "user_device_type"
	COLUMN_USER_OS              = "user_os"
	COLUMN_USER_OS_VERSION      = "user_os_version"
	COLUMN_USER_REFERRER        = "user_referrer"
	COLUMN_BOT                  = "bot"
	COLUMN_THREAT               = "threat"
)

// Yes/No string values used for boolean-like columns (bot, threat).
// Stored as VARCHAR(3) to match the existing all-string column convention.
const (
	VALUE_YES = "yes"
	VALUE_NO  = "no"
)

// Default table name for key-value settings.
const DEFAULT_SETTINGS_TABLE = "statsstore_settings"

// Settings table column names.
const (
	COLUMN_KEY           = "key"
	COLUMN_VALUE         = "value"
	SETTING_EXCLUDED_IPS = "excluded_ips"
)

// MAX_DATETIME is a far-future datetime used as the default soft-delete sentinel.
const MAX_DATETIME = "9999-12-31 23:59:59"
