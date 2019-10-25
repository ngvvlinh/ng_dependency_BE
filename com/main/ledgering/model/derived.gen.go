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

// Type ShopLedger represents table shop_ledger
func sqlgenShopLedger(_ *ShopLedger) bool { return true }

type ShopLedgers []*ShopLedger

const __sqlShopLedger_Table = "shop_ledger"
const __sqlShopLedger_ListCols = "\"id\",\"shop_id\",\"name\",\"bank_account\",\"note\",\"type\",\"status\",\"created_by\",\"created_at\",\"updated_at\",\"deleted_at\""
const __sqlShopLedger_Insert = "INSERT INTO \"shop_ledger\" (" + __sqlShopLedger_ListCols + ") VALUES"
const __sqlShopLedger_Select = "SELECT " + __sqlShopLedger_ListCols + " FROM \"shop_ledger\""
const __sqlShopLedger_Select_history = "SELECT " + __sqlShopLedger_ListCols + " FROM history.\"shop_ledger\""
const __sqlShopLedger_UpdateAll = "UPDATE \"shop_ledger\" SET (" + __sqlShopLedger_ListCols + ")"

func (m *ShopLedger) SQLTableName() string  { return "shop_ledger" }
func (m *ShopLedgers) SQLTableName() string { return "shop_ledger" }
func (m *ShopLedger) SQLListCols() string   { return __sqlShopLedger_ListCols }

func (m *ShopLedger) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopLedger_ListCols + " FROM \"shop_ledger\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopLedger)(nil))
}

func (m *ShopLedger) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.Int64(m.ShopID),
		core.String(m.Name),
		core.JSON{m.BankAccount},
		core.String(m.Note),
		core.String(m.Type),
		core.Int32(m.Status),
		core.Int64(m.CreatedBy),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
	}
}

func (m *ShopLedger) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.Int64)(&m.ShopID),
		(*core.String)(&m.Name),
		core.JSON{&m.BankAccount},
		(*core.String)(&m.Note),
		(*core.String)(&m.Type),
		(*core.Int32)(&m.Status),
		(*core.Int64)(&m.CreatedBy),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
	}
}

func (m *ShopLedger) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopLedgers) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopLedgers, 0, 128)
	for rows.Next() {
		m := new(ShopLedger)
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

func (_ *ShopLedger) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_Select)
	return nil
}

func (_ *ShopLedgers) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_Select)
	return nil
}

func (m *ShopLedger) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopLedgers) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(11)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopLedger) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_ledger")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.BankAccount != nil {
		flag = true
		w.WriteName("bank_account")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.BankAccount})
	}
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
	}
	if m.Type != "" {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.CreatedBy != 0 {
		flag = true
		w.WriteName("created_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedBy)
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
	if !m.DeletedAt.IsZero() {
		flag = true
		w.WriteName("deleted_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DeletedAt)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShopLedger) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(11)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopLedgerHistory map[string]interface{}
type ShopLedgerHistories []map[string]interface{}

func (m *ShopLedgerHistory) SQLTableName() string  { return "history.\"shop_ledger\"" }
func (m ShopLedgerHistories) SQLTableName() string { return "history.\"shop_ledger\"" }

func (m *ShopLedgerHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_Select_history)
	return nil
}

func (m ShopLedgerHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopLedger_Select_history)
	return nil
}

func (m ShopLedgerHistory) ID() core.Interface          { return core.Interface{m["id"]} }
func (m ShopLedgerHistory) ShopID() core.Interface      { return core.Interface{m["shop_id"]} }
func (m ShopLedgerHistory) Name() core.Interface        { return core.Interface{m["name"]} }
func (m ShopLedgerHistory) BankAccount() core.Interface { return core.Interface{m["bank_account"]} }
func (m ShopLedgerHistory) Note() core.Interface        { return core.Interface{m["note"]} }
func (m ShopLedgerHistory) Type() core.Interface        { return core.Interface{m["type"]} }
func (m ShopLedgerHistory) Status() core.Interface      { return core.Interface{m["status"]} }
func (m ShopLedgerHistory) CreatedBy() core.Interface   { return core.Interface{m["created_by"]} }
func (m ShopLedgerHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m ShopLedgerHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m ShopLedgerHistory) DeletedAt() core.Interface   { return core.Interface{m["deleted_at"]} }

func (m *ShopLedgerHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopLedgerHistory, 11)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["name"] = data[2]
	res["bank_account"] = data[3]
	res["note"] = data[4]
	res["type"] = data[5]
	res["status"] = data[6]
	res["created_by"] = data[7]
	res["created_at"] = data[8]
	res["updated_at"] = data[9]
	res["deleted_at"] = data[10]
	*m = res
	return nil
}

func (ms *ShopLedgerHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 11)
	args := make([]interface{}, 11)
	for i := 0; i < 11; i++ {
		args[i] = &data[i]
	}
	res := make(ShopLedgerHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopLedgerHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["name"] = data[2]
		m["bank_account"] = data[3]
		m["note"] = data[4]
		m["type"] = data[5]
		m["status"] = data[6]
		m["created_by"] = data[7]
		m["created_at"] = data[8]
		m["updated_at"] = data[9]
		m["deleted_at"] = data[10]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}