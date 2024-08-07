package blindindexstore

import "github.com/gouniverse/sb"

func (store *Store) sqlTableCreate() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.tableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_SOURCE_REFERENCE_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_SEARCH_VALUE,
			Type: sb.COLUMN_TYPE_LONGTEXT,
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
