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

// Type PurchaseOrder represents table purchase_order
func sqlgenPurchaseOrder(_ *PurchaseOrder) bool { return true }

type PurchaseOrders []*PurchaseOrder

const __sqlPurchaseOrder_Table = "purchase_order"
const __sqlPurchaseOrder_ListCols = "\"id\",\"shop_id\",\"supplier_id\",\"supplier\",\"basket_value\",\"total_discount\",\"total_amount\",\"code\",\"code_norm\",\"note\",\"status\",\"variant_ids\",\"lines\",\"created_by\",\"cancelled_reason\",\"confirmed_at\",\"cancelled_at\",\"created_at\",\"updated_at\",\"deleted_at\",\"supplier_full_name_norm\",\"supplier_phone_norm\""
const __sqlPurchaseOrder_Insert = "INSERT INTO \"purchase_order\" (" + __sqlPurchaseOrder_ListCols + ") VALUES"
const __sqlPurchaseOrder_Select = "SELECT " + __sqlPurchaseOrder_ListCols + " FROM \"purchase_order\""
const __sqlPurchaseOrder_Select_history = "SELECT " + __sqlPurchaseOrder_ListCols + " FROM history.\"purchase_order\""
const __sqlPurchaseOrder_UpdateAll = "UPDATE \"purchase_order\" SET (" + __sqlPurchaseOrder_ListCols + ")"

func (m *PurchaseOrder) SQLTableName() string  { return "purchase_order" }
func (m *PurchaseOrders) SQLTableName() string { return "purchase_order" }
func (m *PurchaseOrder) SQLListCols() string   { return __sqlPurchaseOrder_ListCols }

func (m *PurchaseOrder) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlPurchaseOrder_ListCols + " FROM \"purchase_order\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*PurchaseOrder)(nil))
}

func (m *PurchaseOrder) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.ShopID,
		m.SupplierID,
		core.JSON{m.Supplier},
		core.Int(m.BasketValue),
		core.Int(m.TotalDiscount),
		core.Int(m.TotalAmount),
		core.String(m.Code),
		core.Int(m.CodeNorm),
		core.String(m.Note),
		m.Status,
		core.Array{m.VariantIDs, opts},
		core.JSON{m.Lines},
		m.CreatedBy,
		core.String(m.CancelledReason),
		core.Time(m.ConfirmedAt),
		core.Time(m.CancelledAt),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		core.String(m.SupplierFullNameNorm),
		core.String(m.SupplierPhoneNorm),
	}
}

func (m *PurchaseOrder) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShopID,
		&m.SupplierID,
		core.JSON{&m.Supplier},
		(*core.Int)(&m.BasketValue),
		(*core.Int)(&m.TotalDiscount),
		(*core.Int)(&m.TotalAmount),
		(*core.String)(&m.Code),
		(*core.Int)(&m.CodeNorm),
		(*core.String)(&m.Note),
		&m.Status,
		core.Array{&m.VariantIDs, opts},
		core.JSON{&m.Lines},
		&m.CreatedBy,
		(*core.String)(&m.CancelledReason),
		(*core.Time)(&m.ConfirmedAt),
		(*core.Time)(&m.CancelledAt),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		(*core.String)(&m.SupplierFullNameNorm),
		(*core.String)(&m.SupplierPhoneNorm),
	}
}

func (m *PurchaseOrder) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *PurchaseOrders) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(PurchaseOrders, 0, 128)
	for rows.Next() {
		m := new(PurchaseOrder)
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

func (_ *PurchaseOrder) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_Select)
	return nil
}

func (_ *PurchaseOrders) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_Select)
	return nil
}

func (m *PurchaseOrder) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(22)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms PurchaseOrders) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(22)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *PurchaseOrder) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("purchase_order")
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
	if m.SupplierID != 0 {
		flag = true
		w.WriteName("supplier_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SupplierID)
	}
	if m.Supplier != nil {
		flag = true
		w.WriteName("supplier")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Supplier})
	}
	if m.BasketValue != 0 {
		flag = true
		w.WriteName("basket_value")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.BasketValue)
	}
	if m.TotalDiscount != 0 {
		flag = true
		w.WriteName("total_discount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalDiscount)
	}
	if m.TotalAmount != 0 {
		flag = true
		w.WriteName("total_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalAmount)
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
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.VariantIDs != nil {
		flag = true
		w.WriteName("variant_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.VariantIDs, opts})
	}
	if m.Lines != nil {
		flag = true
		w.WriteName("lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Lines})
	}
	if m.CreatedBy != 0 {
		flag = true
		w.WriteName("created_by")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedBy)
	}
	if m.CancelledReason != "" {
		flag = true
		w.WriteName("cancelled_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelledReason)
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
	if m.SupplierFullNameNorm != "" {
		flag = true
		w.WriteName("supplier_full_name_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SupplierFullNameNorm)
	}
	if m.SupplierPhoneNorm != "" {
		flag = true
		w.WriteName("supplier_phone_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SupplierPhoneNorm)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *PurchaseOrder) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(22)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type PurchaseOrderHistory map[string]interface{}
type PurchaseOrderHistories []map[string]interface{}

func (m *PurchaseOrderHistory) SQLTableName() string  { return "history.\"purchase_order\"" }
func (m PurchaseOrderHistories) SQLTableName() string { return "history.\"purchase_order\"" }

func (m *PurchaseOrderHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_Select_history)
	return nil
}

func (m PurchaseOrderHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlPurchaseOrder_Select_history)
	return nil
}

func (m PurchaseOrderHistory) ID() core.Interface          { return core.Interface{m["id"]} }
func (m PurchaseOrderHistory) ShopID() core.Interface      { return core.Interface{m["shop_id"]} }
func (m PurchaseOrderHistory) SupplierID() core.Interface  { return core.Interface{m["supplier_id"]} }
func (m PurchaseOrderHistory) Supplier() core.Interface    { return core.Interface{m["supplier"]} }
func (m PurchaseOrderHistory) BasketValue() core.Interface { return core.Interface{m["basket_value"]} }
func (m PurchaseOrderHistory) TotalDiscount() core.Interface {
	return core.Interface{m["total_discount"]}
}
func (m PurchaseOrderHistory) TotalAmount() core.Interface { return core.Interface{m["total_amount"]} }
func (m PurchaseOrderHistory) Code() core.Interface        { return core.Interface{m["code"]} }
func (m PurchaseOrderHistory) CodeNorm() core.Interface    { return core.Interface{m["code_norm"]} }
func (m PurchaseOrderHistory) Note() core.Interface        { return core.Interface{m["note"]} }
func (m PurchaseOrderHistory) Status() core.Interface      { return core.Interface{m["status"]} }
func (m PurchaseOrderHistory) VariantIDs() core.Interface  { return core.Interface{m["variant_ids"]} }
func (m PurchaseOrderHistory) Lines() core.Interface       { return core.Interface{m["lines"]} }
func (m PurchaseOrderHistory) CreatedBy() core.Interface   { return core.Interface{m["created_by"]} }
func (m PurchaseOrderHistory) CancelledReason() core.Interface {
	return core.Interface{m["cancelled_reason"]}
}
func (m PurchaseOrderHistory) ConfirmedAt() core.Interface { return core.Interface{m["confirmed_at"]} }
func (m PurchaseOrderHistory) CancelledAt() core.Interface { return core.Interface{m["cancelled_at"]} }
func (m PurchaseOrderHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m PurchaseOrderHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m PurchaseOrderHistory) DeletedAt() core.Interface   { return core.Interface{m["deleted_at"]} }
func (m PurchaseOrderHistory) SupplierFullNameNorm() core.Interface {
	return core.Interface{m["supplier_full_name_norm"]}
}
func (m PurchaseOrderHistory) SupplierPhoneNorm() core.Interface {
	return core.Interface{m["supplier_phone_norm"]}
}

func (m *PurchaseOrderHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 22)
	args := make([]interface{}, 22)
	for i := 0; i < 22; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(PurchaseOrderHistory, 22)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["supplier_id"] = data[2]
	res["supplier"] = data[3]
	res["basket_value"] = data[4]
	res["total_discount"] = data[5]
	res["total_amount"] = data[6]
	res["code"] = data[7]
	res["code_norm"] = data[8]
	res["note"] = data[9]
	res["status"] = data[10]
	res["variant_ids"] = data[11]
	res["lines"] = data[12]
	res["created_by"] = data[13]
	res["cancelled_reason"] = data[14]
	res["confirmed_at"] = data[15]
	res["cancelled_at"] = data[16]
	res["created_at"] = data[17]
	res["updated_at"] = data[18]
	res["deleted_at"] = data[19]
	res["supplier_full_name_norm"] = data[20]
	res["supplier_phone_norm"] = data[21]
	*m = res
	return nil
}

func (ms *PurchaseOrderHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 22)
	args := make([]interface{}, 22)
	for i := 0; i < 22; i++ {
		args[i] = &data[i]
	}
	res := make(PurchaseOrderHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(PurchaseOrderHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["supplier_id"] = data[2]
		m["supplier"] = data[3]
		m["basket_value"] = data[4]
		m["total_discount"] = data[5]
		m["total_amount"] = data[6]
		m["code"] = data[7]
		m["code_norm"] = data[8]
		m["note"] = data[9]
		m["status"] = data[10]
		m["variant_ids"] = data[11]
		m["lines"] = data[12]
		m["created_by"] = data[13]
		m["cancelled_reason"] = data[14]
		m["confirmed_at"] = data[15]
		m["cancelled_at"] = data[16]
		m["created_at"] = data[17]
		m["updated_at"] = data[18]
		m["deleted_at"] = data[19]
		m["supplier_full_name_norm"] = data[20]
		m["supplier_phone_norm"] = data[21]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
