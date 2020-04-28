// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	"time"

	cmsql "o.o/backend/pkg/common/sql/cmsql"
	migration "o.o/backend/pkg/common/sql/migration"
	core "o.o/backend/pkg/common/sql/sq/core"
)

var __sqlModels []interface{ SQLVerifySchema(db *cmsql.Database) }
var __sqlonce sync.Once

func SQLVerifySchema(db *cmsql.Database) {
	__sqlonce.Do(func() {
		for _, m := range __sqlModels {
			m.SQLVerifySchema(db)
		}
	})
}

type SQLWriter = core.SQLWriter

type Accounts []*Account

const __sqlAccount_Table = "account"
const __sqlAccount_ListCols = "\"id\",\"owner_id\",\"name\",\"type\",\"image_url\",\"url_slug\",\"rid\""
const __sqlAccount_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"owner_id\" = EXCLUDED.\"owner_id\",\"name\" = EXCLUDED.\"name\",\"type\" = EXCLUDED.\"type\",\"image_url\" = EXCLUDED.\"image_url\",\"url_slug\" = EXCLUDED.\"url_slug\",\"rid\" = EXCLUDED.\"rid\""
const __sqlAccount_Insert = "INSERT INTO \"account\" (" + __sqlAccount_ListCols + ") VALUES"
const __sqlAccount_Select = "SELECT " + __sqlAccount_ListCols + " FROM \"account\""
const __sqlAccount_Select_history = "SELECT " + __sqlAccount_ListCols + " FROM history.\"account\""
const __sqlAccount_UpdateAll = "UPDATE \"account\" SET (" + __sqlAccount_ListCols + ")"
const __sqlAccount_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT account_pkey DO UPDATE SET"

func (m *Account) SQLTableName() string  { return "account" }
func (m *Accounts) SQLTableName() string { return "account" }
func (m *Account) SQLListCols() string   { return __sqlAccount_ListCols }

func (m *Account) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlAccount_ListCols + " FROM \"account\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Account) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "account"); err != nil {
		db.RecordError(err)
		return
	} else {
		mDBColumnNameAndType = val
	}
	mModelColumnNameAndType := map[string]migration.ColumnDef{
		"id": {
			ColumnName:       "id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"owner_id": {
			ColumnName:       "owner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"name": {
			ColumnName:       "name",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"type": {
			ColumnName:       "type",
			ColumnType:       "account_type.AccountType",
			ColumnDBType:     "enum",
			ColumnTag:        "enum(account_type)",
			ColumnEnumValues: []string{"unknown", "partner", "shop", "affiliate", "etop"},
		},
		"image_url": {
			ColumnName:       "image_url",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"url_slug": {
			ColumnName:       "url_slug",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"rid": {
			ColumnName:       "rid",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "account", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Account)(nil))
}

func (m *Account) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ID,
		m.OwnerID,
		core.String(m.Name),
		m.Type,
		core.String(m.ImageURL),
		core.String(m.URLSlug),
		m.Rid,
	}
}

func (m *Account) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.OwnerID,
		(*core.String)(&m.Name),
		&m.Type,
		(*core.String)(&m.ImageURL),
		(*core.String)(&m.URLSlug),
		&m.Rid,
	}
}

func (m *Account) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Accounts) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Accounts, 0, 128)
	for rows.Next() {
		m := new(Account)
		args := m.SQLScanArgs(opts)
		if err := rows.Scan(args...); err != nil {
			return err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

func (_ *Account) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Select)
	return nil
}

func (_ *Accounts) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Select)
	return nil
}

func (m *Account) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(7)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Accounts) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(7)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Account) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlAccount_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlAccount_ListColsOnConflict)
	return nil
}

func (ms Accounts) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlAccount_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlAccount_ListColsOnConflict)
	return nil
}

func (m *Account) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("account")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.OwnerID != 0 {
		flag = true
		w.WriteName("owner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OwnerID)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.Type != 0 {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if m.ImageURL != "" {
		flag = true
		w.WriteName("image_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ImageURL)
	}
	if m.URLSlug != "" {
		flag = true
		w.WriteName("url_slug")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.URLSlug)
	}
	if m.Rid != 0 {
		flag = true
		w.WriteName("rid")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Rid)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Account) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(7)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type AccountHistory map[string]interface{}
type AccountHistories []map[string]interface{}

func (m *AccountHistory) SQLTableName() string  { return "history.\"account\"" }
func (m AccountHistories) SQLTableName() string { return "history.\"account\"" }

func (m *AccountHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Select_history)
	return nil
}

func (m AccountHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Select_history)
	return nil
}

func (m AccountHistory) ID() core.Interface       { return core.Interface{m["id"]} }
func (m AccountHistory) OwnerID() core.Interface  { return core.Interface{m["owner_id"]} }
func (m AccountHistory) Name() core.Interface     { return core.Interface{m["name"]} }
func (m AccountHistory) Type() core.Interface     { return core.Interface{m["type"]} }
func (m AccountHistory) ImageURL() core.Interface { return core.Interface{m["image_url"]} }
func (m AccountHistory) URLSlug() core.Interface  { return core.Interface{m["url_slug"]} }
func (m AccountHistory) Rid() core.Interface      { return core.Interface{m["rid"]} }

func (m *AccountHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 7)
	args := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(AccountHistory, 7)
	res["id"] = data[0]
	res["owner_id"] = data[1]
	res["name"] = data[2]
	res["type"] = data[3]
	res["image_url"] = data[4]
	res["url_slug"] = data[5]
	res["rid"] = data[6]
	*m = res
	return nil
}

func (ms *AccountHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 7)
	args := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		args[i] = &data[i]
	}
	res := make(AccountHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(AccountHistory)
		m["id"] = data[0]
		m["owner_id"] = data[1]
		m["name"] = data[2]
		m["type"] = data[3]
		m["image_url"] = data[4]
		m["url_slug"] = data[5]
		m["rid"] = data[6]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
