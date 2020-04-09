// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	time "time"

	cmsql "etop.vn/backend/pkg/common/sql/cmsql"
	migration "etop.vn/backend/pkg/common/sql/migration"
	core "etop.vn/backend/pkg/common/sql/sq/core"
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
const __sqlAccount_ListCols = "\"id\",\"name\""
const __sqlAccount_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"name\" = EXCLUDED.\"name\""
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
		"name": {
			ColumnName:       "name",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
		core.String(m.Name),
	}
}

func (m *Account) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.Name),
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
	w.WriteMarkers(2)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Accounts) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(2)
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
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
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
	w.WriteMarkers(2)
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

func (m AccountHistory) ID() core.Interface   { return core.Interface{m["id"]} }
func (m AccountHistory) Name() core.Interface { return core.Interface{m["name"]} }

func (m *AccountHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 2)
	args := make([]interface{}, 2)
	for i := 0; i < 2; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(AccountHistory, 2)
	res["id"] = data[0]
	res["name"] = data[1]
	*m = res
	return nil
}

func (ms *AccountHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 2)
	args := make([]interface{}, 2)
	for i := 0; i < 2; i++ {
		args[i] = &data[i]
	}
	res := make(AccountHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(AccountHistory)
		m["id"] = data[0]
		m["name"] = data[1]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

type Foos []*Foo

const __sqlFoo_Table = "foo"
const __sqlFoo_ListCols = "\"id\",\"account_id\",\"abc_2\",\"def2\",\"created_at\",\"updated_at\""
const __sqlFoo_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"account_id\" = EXCLUDED.\"account_id\",\"abc_2\" = EXCLUDED.\"abc_2\",\"def2\" = EXCLUDED.\"def2\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\""
const __sqlFoo_Insert = "INSERT INTO \"foo\" (" + __sqlFoo_ListCols + ") VALUES"
const __sqlFoo_Select = "SELECT " + __sqlFoo_ListCols + " FROM \"foo\""
const __sqlFoo_Select_history = "SELECT " + __sqlFoo_ListCols + " FROM history.\"foo\""
const __sqlFoo_UpdateAll = "UPDATE \"foo\" SET (" + __sqlFoo_ListCols + ")"
const __sqlFoo_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT foo_pkey DO UPDATE SET"

func (m *Foo) SQLTableName() string  { return "foo" }
func (m *Foos) SQLTableName() string { return "foo" }
func (m *Foo) SQLListCols() string   { return __sqlFoo_ListCols }

func (m *Foo) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlFoo_ListCols + " FROM \"foo\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *Foo) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "foo"); err != nil {
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
		"account_id": {
			ColumnName:       "account_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"abc_2": {
			ColumnName:       "abc_2",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"def2": {
			ColumnName:       "def2",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"created_at": {
			ColumnName:       "created_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"updated_at": {
			ColumnName:       "updated_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "foo", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Foo)(nil))
}

func (m *Foo) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.AccountID,
		core.String(m.ABC),
		core.String(m.Def2),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
	}
}

func (m *Foo) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.AccountID,
		(*core.String)(&m.ABC),
		(*core.String)(&m.Def2),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
	}
}

func (m *Foo) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Foos) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Foos, 0, 128)
	for rows.Next() {
		m := new(Foo)
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

func (_ *Foo) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_Select)
	return nil
}

func (_ *Foos) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_Select)
	return nil
}

func (m *Foo) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(6)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Foos) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(6)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Foo) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlFoo_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlFoo_ListColsOnConflict)
	return nil
}

func (ms Foos) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlFoo_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlFoo_ListColsOnConflict)
	return nil
}

func (m *Foo) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("foo")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.AccountID != 0 {
		flag = true
		w.WriteName("account_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AccountID)
	}
	if m.ABC != "" {
		flag = true
		w.WriteName("abc_2")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ABC)
	}
	if m.Def2 != "" {
		flag = true
		w.WriteName("def2")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Def2)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if !m.UpdatedAt.IsZero() {
		flag = true
		w.WriteName("updated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Now(m.UpdatedAt, time.Now(), true))
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Foo) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(6)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type FooHistory map[string]interface{}
type FooHistories []map[string]interface{}

func (m *FooHistory) SQLTableName() string  { return "history.\"foo\"" }
func (m FooHistories) SQLTableName() string { return "history.\"foo\"" }

func (m *FooHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_Select_history)
	return nil
}

func (m FooHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlFoo_Select_history)
	return nil
}

func (m FooHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m FooHistory) AccountID() core.Interface { return core.Interface{m["account_id"]} }
func (m FooHistory) ABC() core.Interface       { return core.Interface{m["abc_2"]} }
func (m FooHistory) Def2() core.Interface      { return core.Interface{m["def2"]} }
func (m FooHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m FooHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }

func (m *FooHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 6)
	args := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(FooHistory, 6)
	res["id"] = data[0]
	res["account_id"] = data[1]
	res["abc_2"] = data[2]
	res["def2"] = data[3]
	res["created_at"] = data[4]
	res["updated_at"] = data[5]
	*m = res
	return nil
}

func (ms *FooHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 6)
	args := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		args[i] = &data[i]
	}
	res := make(FooHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(FooHistory)
		m["id"] = data[0]
		m["account_id"] = data[1]
		m["abc_2"] = data[2]
		m["def2"] = data[3]
		m["created_at"] = data[4]
		m["updated_at"] = data[5]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

type FooWithAccounts []*FooWithAccount

func (m *FooWithAccount) SQLTableName() string  { return "foo" }
func (m *FooWithAccounts) SQLTableName() string { return "foo" }

func (m *FooWithAccount) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *FooWithAccounts) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(FooWithAccounts, 0, 128)
	for rows.Next() {
		m := new(FooWithAccount)
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

func (m *FooWithAccount) SQLSelect(w SQLWriter) error {
	(*FooWithAccount)(nil).__sqlSelect(w)
	w.WriteByte(' ')
	(*FooWithAccount)(nil).__sqlJoin(w)
	return nil
}

func (m *FooWithAccounts) SQLSelect(w SQLWriter) error {
	return (*FooWithAccount)(nil).SQLSelect(w)
}

func (m *FooWithAccount) SQLJoin(w SQLWriter) error {
	m.__sqlJoin(w)
	return nil
}

func (m *FooWithAccounts) SQLJoin(w SQLWriter) error {
	return (*FooWithAccount)(nil).SQLJoin(w)
}

func (m *FooWithAccount) __sqlSelect(w SQLWriter) {
	w.WriteRawString("SELECT ")
	core.WriteCols(w, "foo", (*Foo)(nil).SQLListCols())
	w.WriteByte(',')
	core.WriteCols(w, "a", (*Account)(nil).SQLListCols())
}

func (m *FooWithAccount) __sqlJoin(w SQLWriter) {
	w.WriteRawString("FROM ")
	w.WriteName("foo")
	w.WriteRawString(" AS ")
	w.WriteName("foo")
	w.WriteRawString(" JOIN ")
	w.WriteName((*Account)(nil).SQLTableName())
	w.WriteRawString(" AS a ON")
	w.WriteQueryString(" foo.account_id = a.id")
}

func (m *FooWithAccount) SQLScanArgs(opts core.Opts) []interface{} {
	args := make([]interface{}, 0, 64) // TODO: pre-calculate length
	m.Foo = new(Foo)
	args = append(args, m.Foo.SQLScanArgs(opts)...)
	m.Account = new(Account)
	args = append(args, m.Account.SQLScanArgs(opts)...)
	return args
}
