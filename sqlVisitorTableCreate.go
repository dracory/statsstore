package statsstore

import (
	"github.com/gouniverse/sb"
)

// sqlVisitorTableCreate returns a SQL string for creating the visitor table
func (st *Store) sqlVisitorTableCreate() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(st.db)).
		Table(st.visitorTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			PrimaryKey: true,
			Length:     40,
		}).
		Column(sb.Column{
			Name:   COLUMN_PATH,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 512,
		}).
		Column(sb.Column{
			Name:   COLUMN_IP_ADDRESS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_COUNTRY,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 2,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_ACCEPT_LANGUAGE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_ACCEPT_ENCODING,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_AGENT,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_OS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 12,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_OS_VERSION,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 12,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_DEVICE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_DEVICE_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 12,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_BROWSER,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_BROWSER_VERSION,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 24,
		}).
		Column(sb.Column{
			Name:   COLUMN_USER_REFERRER,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql
}
