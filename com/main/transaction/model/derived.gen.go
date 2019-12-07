// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	"time"

	"etop.vn/backend/pkg/common/cmsql"
	core "etop.vn/backend/pkg/common/sq/core"
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

// Type Transaction represents table transaction
func sqlgenTransaction(_ *Transaction) bool { return true }

type Transactions []*Transaction

const __sqlTransaction_Table = "transaction"
const __sqlTransaction_ListCols = "\"id\",\"amount\",\"account_id\",\"status\",\"type\",\"note\",\"metadata\",\"created_at\",\"updated_at\""
const __sqlTransaction_Insert = "INSERT INTO \"transaction\" (" + __sqlTransaction_ListCols + ") VALUES"
const __sqlTransaction_Select = "SELECT " + __sqlTransaction_ListCols + " FROM \"transaction\""
const __sqlTransaction_Select_history = "SELECT " + __sqlTransaction_ListCols + " FROM history.\"transaction\""
const __sqlTransaction_UpdateAll = "UPDATE \"transaction\" SET (" + __sqlTransaction_ListCols + ")"

func (m *Transaction) SQLTableName() string  { return "transaction" }
func (m *Transactions) SQLTableName() string { return "transaction" }
func (m *Transaction) SQLListCols() string   { return __sqlTransaction_ListCols }

func (m *Transaction) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlTransaction_ListCols + " FROM \"transaction\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Transaction)(nil))
}

func (m *Transaction) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.Int(m.Amount),
		m.AccountID,
		m.Status,
		core.String(m.Type),
		core.String(m.Note),
		core.JSON{m.Metadata},
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
	}
}

func (m *Transaction) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.Int)(&m.Amount),
		&m.AccountID,
		&m.Status,
		(*core.String)(&m.Type),
		(*core.String)(&m.Note),
		core.JSON{&m.Metadata},
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
	}
}

func (m *Transaction) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Transactions) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Transactions, 0, 128)
	for rows.Next() {
		m := new(Transaction)
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

func (_ *Transaction) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_Select)
	return nil
}

func (_ *Transactions) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_Select)
	return nil
}

func (m *Transaction) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Transactions) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(9)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Transaction) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("transaction")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.Amount != 0 {
		flag = true
		w.WriteName("amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Amount)
	}
	if m.AccountID != 0 {
		flag = true
		w.WriteName("account_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AccountID)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.Type != "" {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
	}
	if m.Metadata != nil {
		flag = true
		w.WriteName("metadata")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Metadata})
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

func (m *Transaction) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type TransactionHistory map[string]interface{}
type TransactionHistories []map[string]interface{}

func (m *TransactionHistory) SQLTableName() string  { return "history.\"transaction\"" }
func (m TransactionHistories) SQLTableName() string { return "history.\"transaction\"" }

func (m *TransactionHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_Select_history)
	return nil
}

func (m TransactionHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlTransaction_Select_history)
	return nil
}

func (m TransactionHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m TransactionHistory) Amount() core.Interface    { return core.Interface{m["amount"]} }
func (m TransactionHistory) AccountID() core.Interface { return core.Interface{m["account_id"]} }
func (m TransactionHistory) Status() core.Interface    { return core.Interface{m["status"]} }
func (m TransactionHistory) Type() core.Interface      { return core.Interface{m["type"]} }
func (m TransactionHistory) Note() core.Interface      { return core.Interface{m["note"]} }
func (m TransactionHistory) Metadata() core.Interface  { return core.Interface{m["metadata"]} }
func (m TransactionHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m TransactionHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }

func (m *TransactionHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(TransactionHistory, 9)
	res["id"] = data[0]
	res["amount"] = data[1]
	res["account_id"] = data[2]
	res["status"] = data[3]
	res["type"] = data[4]
	res["note"] = data[5]
	res["metadata"] = data[6]
	res["created_at"] = data[7]
	res["updated_at"] = data[8]
	*m = res
	return nil
}

func (ms *TransactionHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	res := make(TransactionHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(TransactionHistory)
		m["id"] = data[0]
		m["amount"] = data[1]
		m["account_id"] = data[2]
		m["status"] = data[3]
		m["type"] = data[4]
		m["note"] = data[5]
		m["metadata"] = data[6]
		m["created_at"] = data[7]
		m["updated_at"] = data[8]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
