// Code generated by goderive DO NOT EDIT.

package testsrc

import (
	"database/sql"
	"sync"
	"time"

	"etop.vn/backend/pkg/common/sql/cmsql"
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

// Type Account represents table account
func sqlgenAccount(_ *Account) bool { return true }

type Accounts []*Account

const __sqlAccount_Table = "account"
const __sqlAccount_ListCols = "\"id\",\"first_name\",\"last_name\""
const __sqlAccount_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"first_name\" = EXCLUDED.\"first_name\",\"last_name\" = EXCLUDED.\"last_name\""
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

func init() {
	__sqlModels = append(__sqlModels, (*Account)(nil))
}

func (m *Account) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ID,
		core.String(m.FirstName),
		core.String(m.LastName),
	}
}

func (m *Account) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.FirstName),
		(*core.String)(&m.LastName),
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
	w.WriteMarkers(3)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Accounts) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlAccount_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(3)
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
	if m.FirstName != "" {
		flag = true
		w.WriteName("first_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.FirstName)
	}
	if m.LastName != "" {
		flag = true
		w.WriteName("last_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.LastName)
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
	w.WriteMarkers(3)
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

func (m AccountHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m AccountHistory) FirstName() core.Interface { return core.Interface{m["first_name"]} }
func (m AccountHistory) LastName() core.Interface  { return core.Interface{m["last_name"]} }

func (m *AccountHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 3)
	args := make([]interface{}, 3)
	for i := 0; i < 3; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(AccountHistory, 3)
	res["id"] = data[0]
	res["first_name"] = data[1]
	res["last_name"] = data[2]
	*m = res
	return nil
}

func (ms *AccountHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 3)
	args := make([]interface{}, 3)
	for i := 0; i < 3; i++ {
		args[i] = &data[i]
	}
	res := make(AccountHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(AccountHistory)
		m["id"] = data[0]
		m["first_name"] = data[1]
		m["last_name"] = data[2]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
