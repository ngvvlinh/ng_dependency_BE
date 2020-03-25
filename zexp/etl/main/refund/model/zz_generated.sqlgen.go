// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	time "time"

	cmsql "etop.vn/backend/pkg/common/sql/cmsql"
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

type Refunds []*Refund

const __sqlRefund_Table = "refund"
const __sqlRefund_ListCols = "\"id\",\"shop_id\",\"order_id\",\"code\",\"code_norm\",\"note\",\"lines\",\"adjustment_lines\",\"total_adjustment\",\"created_at\",\"updated_at\",\"cancelled_at\",\"confirmed_at\",\"created_by\",\"updated_by\",\"cancel_reason\",\"status\",\"customer_id\",\"total_amount\",\"basket_value\",\"rid\""
const __sqlRefund_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"order_id\" = EXCLUDED.\"order_id\",\"code\" = EXCLUDED.\"code\",\"code_norm\" = EXCLUDED.\"code_norm\",\"note\" = EXCLUDED.\"note\",\"lines\" = EXCLUDED.\"lines\",\"adjustment_lines\" = EXCLUDED.\"adjustment_lines\",\"total_adjustment\" = EXCLUDED.\"total_adjustment\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"cancelled_at\" = EXCLUDED.\"cancelled_at\",\"confirmed_at\" = EXCLUDED.\"confirmed_at\",\"created_by\" = EXCLUDED.\"created_by\",\"updated_by\" = EXCLUDED.\"updated_by\",\"cancel_reason\" = EXCLUDED.\"cancel_reason\",\"status\" = EXCLUDED.\"status\",\"customer_id\" = EXCLUDED.\"customer_id\",\"total_amount\" = EXCLUDED.\"total_amount\",\"basket_value\" = EXCLUDED.\"basket_value\",\"rid\" = EXCLUDED.\"rid\""
const __sqlRefund_Insert = "INSERT INTO \"refund\" (" + __sqlRefund_ListCols + ") VALUES"
const __sqlRefund_Select = "SELECT " + __sqlRefund_ListCols + " FROM \"refund\""
const __sqlRefund_Select_history = "SELECT " + __sqlRefund_ListCols + " FROM history.\"refund\""
const __sqlRefund_UpdateAll = "UPDATE \"refund\" SET (" + __sqlRefund_ListCols + ")"
const __sqlRefund_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT refund_pkey DO UPDATE SET"

func (m *Refund) SQLTableName() string  { return "refund" }
func (m *Refunds) SQLTableName() string { return "refund" }
func (m *Refund) SQLListCols() string   { return __sqlRefund_ListCols }

func (m *Refund) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlRefund_ListCols + " FROM \"refund\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*Refund)(nil))
}

func (m *Refund) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ID,
		m.ShopID,
		m.OrderID,
		core.String(m.Code),
		core.Int(m.CodeNorm),
		core.String(m.Note),
		core.JSON{m.Lines},
		core.JSON{m.AdjustmentLines},
		core.Int(m.TotalAdjustment),
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		core.Time(m.CancelledAt),
		core.Time(m.ConfirmedAt),
		m.CreatedBy,
		m.UpdatedBy,
		core.String(m.CancelReason),
		m.Status,
		m.CustomerID,
		core.Int(m.TotalAmount),
		core.Int(m.BasketValue),
		m.Rid,
	}
}

func (m *Refund) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShopID,
		&m.OrderID,
		(*core.String)(&m.Code),
		(*core.Int)(&m.CodeNorm),
		(*core.String)(&m.Note),
		core.JSON{&m.Lines},
		core.JSON{&m.AdjustmentLines},
		(*core.Int)(&m.TotalAdjustment),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.CancelledAt),
		(*core.Time)(&m.ConfirmedAt),
		&m.CreatedBy,
		&m.UpdatedBy,
		(*core.String)(&m.CancelReason),
		&m.Status,
		&m.CustomerID,
		(*core.Int)(&m.TotalAmount),
		(*core.Int)(&m.BasketValue),
		&m.Rid,
	}
}

func (m *Refund) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Refunds) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Refunds, 0, 128)
	for rows.Next() {
		m := new(Refund)
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

func (_ *Refund) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_Select)
	return nil
}

func (_ *Refunds) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_Select)
	return nil
}

func (m *Refund) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(21)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Refunds) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(21)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Refund) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlRefund_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlRefund_ListColsOnConflict)
	return nil
}

func (ms Refunds) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlRefund_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlRefund_ListColsOnConflict)
	return nil
}

func (m *Refund) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("refund")
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
	if m.OrderID != 0 {
		flag = true
		w.WriteName("order_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OrderID)
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
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
	}
	if m.Lines != nil {
		flag = true
		w.WriteName("lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Lines})
	}
	if m.AdjustmentLines != nil {
		flag = true
		w.WriteName("adjustment_lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.AdjustmentLines})
	}
	if m.TotalAdjustment != 0 {
		flag = true
		w.WriteName("total_adjustment")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalAdjustment)
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
		w.WriteArg(m.UpdatedAt)
	}
	if !m.CancelledAt.IsZero() {
		flag = true
		w.WriteName("cancelled_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelledAt)
	}
	if !m.ConfirmedAt.IsZero() {
		flag = true
		w.WriteName("confirmed_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConfirmedAt)
	}
	if m.CreatedBy != 0 {
		flag = true
		w.WriteName("created_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedBy)
	}
	if m.UpdatedBy != 0 {
		flag = true
		w.WriteName("updated_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.UpdatedBy)
	}
	if m.CancelReason != "" {
		flag = true
		w.WriteName("cancel_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelReason)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.CustomerID != 0 {
		flag = true
		w.WriteName("customer_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CustomerID)
	}
	if m.TotalAmount != 0 {
		flag = true
		w.WriteName("total_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalAmount)
	}
	if m.BasketValue != 0 {
		flag = true
		w.WriteName("basket_value")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.BasketValue)
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

func (m *Refund) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(21)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type RefundHistory map[string]interface{}
type RefundHistories []map[string]interface{}

func (m *RefundHistory) SQLTableName() string  { return "history.\"refund\"" }
func (m RefundHistories) SQLTableName() string { return "history.\"refund\"" }

func (m *RefundHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_Select_history)
	return nil
}

func (m RefundHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlRefund_Select_history)
	return nil
}

func (m RefundHistory) ID() core.Interface              { return core.Interface{m["id"]} }
func (m RefundHistory) ShopID() core.Interface          { return core.Interface{m["shop_id"]} }
func (m RefundHistory) OrderID() core.Interface         { return core.Interface{m["order_id"]} }
func (m RefundHistory) Code() core.Interface            { return core.Interface{m["code"]} }
func (m RefundHistory) CodeNorm() core.Interface        { return core.Interface{m["code_norm"]} }
func (m RefundHistory) Note() core.Interface            { return core.Interface{m["note"]} }
func (m RefundHistory) Lines() core.Interface           { return core.Interface{m["lines"]} }
func (m RefundHistory) AdjustmentLines() core.Interface { return core.Interface{m["adjustment_lines"]} }
func (m RefundHistory) TotalAdjustment() core.Interface { return core.Interface{m["total_adjustment"]} }
func (m RefundHistory) CreatedAt() core.Interface       { return core.Interface{m["created_at"]} }
func (m RefundHistory) UpdatedAt() core.Interface       { return core.Interface{m["updated_at"]} }
func (m RefundHistory) CancelledAt() core.Interface     { return core.Interface{m["cancelled_at"]} }
func (m RefundHistory) ConfirmedAt() core.Interface     { return core.Interface{m["confirmed_at"]} }
func (m RefundHistory) CreatedBy() core.Interface       { return core.Interface{m["created_by"]} }
func (m RefundHistory) UpdatedBy() core.Interface       { return core.Interface{m["updated_by"]} }
func (m RefundHistory) CancelReason() core.Interface    { return core.Interface{m["cancel_reason"]} }
func (m RefundHistory) Status() core.Interface          { return core.Interface{m["status"]} }
func (m RefundHistory) CustomerID() core.Interface      { return core.Interface{m["customer_id"]} }
func (m RefundHistory) TotalAmount() core.Interface     { return core.Interface{m["total_amount"]} }
func (m RefundHistory) BasketValue() core.Interface     { return core.Interface{m["basket_value"]} }
func (m RefundHistory) Rid() core.Interface             { return core.Interface{m["rid"]} }

func (m *RefundHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 21)
	args := make([]interface{}, 21)
	for i := 0; i < 21; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(RefundHistory, 21)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["order_id"] = data[2]
	res["code"] = data[3]
	res["code_norm"] = data[4]
	res["note"] = data[5]
	res["lines"] = data[6]
	res["adjustment_lines"] = data[7]
	res["total_adjustment"] = data[8]
	res["created_at"] = data[9]
	res["updated_at"] = data[10]
	res["cancelled_at"] = data[11]
	res["confirmed_at"] = data[12]
	res["created_by"] = data[13]
	res["updated_by"] = data[14]
	res["cancel_reason"] = data[15]
	res["status"] = data[16]
	res["customer_id"] = data[17]
	res["total_amount"] = data[18]
	res["basket_value"] = data[19]
	res["rid"] = data[20]
	*m = res
	return nil
}

func (ms *RefundHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 21)
	args := make([]interface{}, 21)
	for i := 0; i < 21; i++ {
		args[i] = &data[i]
	}
	res := make(RefundHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(RefundHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["order_id"] = data[2]
		m["code"] = data[3]
		m["code_norm"] = data[4]
		m["note"] = data[5]
		m["lines"] = data[6]
		m["adjustment_lines"] = data[7]
		m["total_adjustment"] = data[8]
		m["created_at"] = data[9]
		m["updated_at"] = data[10]
		m["cancelled_at"] = data[11]
		m["confirmed_at"] = data[12]
		m["created_by"] = data[13]
		m["updated_by"] = data[14]
		m["cancel_reason"] = data[15]
		m["status"] = data[16]
		m["customer_id"] = data[17]
		m["total_amount"] = data[18]
		m["basket_value"] = data[19]
		m["rid"] = data[20]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
