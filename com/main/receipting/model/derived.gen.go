// Code generated by goderive DO NOT EDIT.

package model

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

// Type Receipt represents table receipt
func sqlgenReceipt(_ *Receipt) bool { return true }

type Receipts []*Receipt

const __sqlReceipt_Table = "receipt"
const __sqlReceipt_ListCols = "\"id\",\"shop_id\",\"trader_id\",\"code\",\"code_norm\",\"title\",\"type\",\"description\",\"trader_full_name_norm\",\"amount\",\"status\",\"ref_ids\",\"ref_type\",\"lines\",\"ledger_id\",\"trader\",\"cancelled_reason\",\"created_type\",\"created_by\",\"paid_at\",\"confirmed_at\",\"cancelled_at\",\"created_at\",\"updated_at\",\"deleted_at\""
const __sqlReceipt_Insert = "INSERT INTO \"receipt\" (" + __sqlReceipt_ListCols + ") VALUES"
const __sqlReceipt_Select = "SELECT " + __sqlReceipt_ListCols + " FROM \"receipt\""
const __sqlReceipt_Select_history = "SELECT " + __sqlReceipt_ListCols + " FROM history.\"receipt\""
const __sqlReceipt_UpdateAll = "UPDATE \"receipt\" SET (" + __sqlReceipt_ListCols + ")"

func (m *Receipt) SQLTableName() string  { return "receipt" }
func (m *Receipts) SQLTableName() string { return "receipt" }
func (m *Receipt) SQLListCols() string   { return __sqlReceipt_ListCols }

func (m *Receipt) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlReceipt_ListCols + " FROM \"receipt\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Receipt)(nil))
}

func (m *Receipt) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.ShopID,
		m.TraderID,
		core.String(m.Code),
		core.Int(m.CodeNorm),
		core.String(m.Title),
		m.Type,
		core.String(m.Description),
		core.String(m.TraderFullNameNorm),
		core.Int(m.Amount),
		m.Status,
		core.Array{m.RefIDs, opts},
		m.RefType,
		core.JSON{m.Lines},
		m.LedgerID,
		core.JSON{m.Trader},
		core.String(m.CancelledReason),
		core.String(m.CreatedType),
		m.CreatedBy,
		core.Time(m.PaidAt),
		core.Time(m.ConfirmedAt),
		core.Time(m.CancelledAt),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
	}
}

func (m *Receipt) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShopID,
		&m.TraderID,
		(*core.String)(&m.Code),
		(*core.Int)(&m.CodeNorm),
		(*core.String)(&m.Title),
		&m.Type,
		(*core.String)(&m.Description),
		(*core.String)(&m.TraderFullNameNorm),
		(*core.Int)(&m.Amount),
		&m.Status,
		core.Array{&m.RefIDs, opts},
		&m.RefType,
		core.JSON{&m.Lines},
		&m.LedgerID,
		core.JSON{&m.Trader},
		(*core.String)(&m.CancelledReason),
		(*core.String)(&m.CreatedType),
		&m.CreatedBy,
		(*core.Time)(&m.PaidAt),
		(*core.Time)(&m.ConfirmedAt),
		(*core.Time)(&m.CancelledAt),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
	}
}

func (m *Receipt) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Receipts) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Receipts, 0, 128)
	for rows.Next() {
		m := new(Receipt)
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

func (_ *Receipt) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_Select)
	return nil
}

func (_ *Receipts) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_Select)
	return nil
}

func (m *Receipt) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(25)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Receipts) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(25)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Receipt) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("receipt")
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
	if m.TraderID != 0 {
		flag = true
		w.WriteName("trader_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TraderID)
	}
	if m.Code != "" {
		flag = true
		w.WriteName("code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Code)
	}
	if m.CodeNorm != 0 {
		flag = true
		w.WriteName("code_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CodeNorm)
	}
	if m.Title != "" {
		flag = true
		w.WriteName("title")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Title)
	}
	if m.Type != 0 {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if m.Description != "" {
		flag = true
		w.WriteName("description")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Description)
	}
	if m.TraderFullNameNorm != "" {
		flag = true
		w.WriteName("trader_full_name_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TraderFullNameNorm)
	}
	if m.Amount != 0 {
		flag = true
		w.WriteName("amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Amount)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.RefIDs != nil {
		flag = true
		w.WriteName("ref_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.RefIDs, opts})
	}
	if m.RefType != 0 {
		flag = true
		w.WriteName("ref_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RefType)
	}
	if m.Lines != nil {
		flag = true
		w.WriteName("lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Lines})
	}
	if m.LedgerID != 0 {
		flag = true
		w.WriteName("ledger_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.LedgerID)
	}
	if m.Trader != nil {
		flag = true
		w.WriteName("trader")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Trader})
	}
	if m.CancelledReason != "" {
		flag = true
		w.WriteName("cancelled_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelledReason)
	}
	if m.CreatedType != "" {
		flag = true
		w.WriteName("created_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedType)
	}
	if m.CreatedBy != 0 {
		flag = true
		w.WriteName("created_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedBy)
	}
	if !m.PaidAt.IsZero() {
		flag = true
		w.WriteName("paid_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaidAt)
	}
	if !m.ConfirmedAt.IsZero() {
		flag = true
		w.WriteName("confirmed_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConfirmedAt)
	}
	if !m.CancelledAt.IsZero() {
		flag = true
		w.WriteName("cancelled_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelledAt)
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

func (m *Receipt) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(25)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ReceiptHistory map[string]interface{}
type ReceiptHistories []map[string]interface{}

func (m *ReceiptHistory) SQLTableName() string  { return "history.\"receipt\"" }
func (m ReceiptHistories) SQLTableName() string { return "history.\"receipt\"" }

func (m *ReceiptHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_Select_history)
	return nil
}

func (m ReceiptHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlReceipt_Select_history)
	return nil
}

func (m ReceiptHistory) ID() core.Interface          { return core.Interface{m["id"]} }
func (m ReceiptHistory) ShopID() core.Interface      { return core.Interface{m["shop_id"]} }
func (m ReceiptHistory) TraderID() core.Interface    { return core.Interface{m["trader_id"]} }
func (m ReceiptHistory) Code() core.Interface        { return core.Interface{m["code"]} }
func (m ReceiptHistory) CodeNorm() core.Interface    { return core.Interface{m["code_norm"]} }
func (m ReceiptHistory) Title() core.Interface       { return core.Interface{m["title"]} }
func (m ReceiptHistory) Type() core.Interface        { return core.Interface{m["type"]} }
func (m ReceiptHistory) Description() core.Interface { return core.Interface{m["description"]} }
func (m ReceiptHistory) TraderFullNameNorm() core.Interface {
	return core.Interface{m["trader_full_name_norm"]}
}
func (m ReceiptHistory) Amount() core.Interface          { return core.Interface{m["amount"]} }
func (m ReceiptHistory) Status() core.Interface          { return core.Interface{m["status"]} }
func (m ReceiptHistory) RefIDs() core.Interface          { return core.Interface{m["ref_ids"]} }
func (m ReceiptHistory) RefType() core.Interface         { return core.Interface{m["ref_type"]} }
func (m ReceiptHistory) Lines() core.Interface           { return core.Interface{m["lines"]} }
func (m ReceiptHistory) LedgerID() core.Interface        { return core.Interface{m["ledger_id"]} }
func (m ReceiptHistory) Trader() core.Interface          { return core.Interface{m["trader"]} }
func (m ReceiptHistory) CancelledReason() core.Interface { return core.Interface{m["cancelled_reason"]} }
func (m ReceiptHistory) CreatedType() core.Interface     { return core.Interface{m["created_type"]} }
func (m ReceiptHistory) CreatedBy() core.Interface       { return core.Interface{m["created_by"]} }
func (m ReceiptHistory) PaidAt() core.Interface          { return core.Interface{m["paid_at"]} }
func (m ReceiptHistory) ConfirmedAt() core.Interface     { return core.Interface{m["confirmed_at"]} }
func (m ReceiptHistory) CancelledAt() core.Interface     { return core.Interface{m["cancelled_at"]} }
func (m ReceiptHistory) CreatedAt() core.Interface       { return core.Interface{m["created_at"]} }
func (m ReceiptHistory) UpdatedAt() core.Interface       { return core.Interface{m["updated_at"]} }
func (m ReceiptHistory) DeletedAt() core.Interface       { return core.Interface{m["deleted_at"]} }

func (m *ReceiptHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 25)
	args := make([]interface{}, 25)
	for i := 0; i < 25; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ReceiptHistory, 25)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["trader_id"] = data[2]
	res["code"] = data[3]
	res["code_norm"] = data[4]
	res["title"] = data[5]
	res["type"] = data[6]
	res["description"] = data[7]
	res["trader_full_name_norm"] = data[8]
	res["amount"] = data[9]
	res["status"] = data[10]
	res["ref_ids"] = data[11]
	res["ref_type"] = data[12]
	res["lines"] = data[13]
	res["ledger_id"] = data[14]
	res["trader"] = data[15]
	res["cancelled_reason"] = data[16]
	res["created_type"] = data[17]
	res["created_by"] = data[18]
	res["paid_at"] = data[19]
	res["confirmed_at"] = data[20]
	res["cancelled_at"] = data[21]
	res["created_at"] = data[22]
	res["updated_at"] = data[23]
	res["deleted_at"] = data[24]
	*m = res
	return nil
}

func (ms *ReceiptHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 25)
	args := make([]interface{}, 25)
	for i := 0; i < 25; i++ {
		args[i] = &data[i]
	}
	res := make(ReceiptHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ReceiptHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["trader_id"] = data[2]
		m["code"] = data[3]
		m["code_norm"] = data[4]
		m["title"] = data[5]
		m["type"] = data[6]
		m["description"] = data[7]
		m["trader_full_name_norm"] = data[8]
		m["amount"] = data[9]
		m["status"] = data[10]
		m["ref_ids"] = data[11]
		m["ref_type"] = data[12]
		m["lines"] = data[13]
		m["ledger_id"] = data[14]
		m["trader"] = data[15]
		m["cancelled_reason"] = data[16]
		m["created_type"] = data[17]
		m["created_by"] = data[18]
		m["paid_at"] = data[19]
		m["confirmed_at"] = data[20]
		m["cancelled_at"] = data[21]
		m["created_at"] = data[22]
		m["updated_at"] = data[23]
		m["deleted_at"] = data[24]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
