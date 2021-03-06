// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package model

import (
	"database/sql"
	"sync"
	time "time"

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

type InventoryVariants []*InventoryVariant

const __sqlInventoryVariant_Table = "inventory_variant"
const __sqlInventoryVariant_ListCols = "\"shop_id\",\"variant_id\",\"quantity_on_hand\",\"quantity_picked\",\"cost_price\",\"created_at\",\"updated_at\",\"rid\""
const __sqlInventoryVariant_ListColsOnConflict = "\"shop_id\" = EXCLUDED.\"shop_id\",\"variant_id\" = EXCLUDED.\"variant_id\",\"quantity_on_hand\" = EXCLUDED.\"quantity_on_hand\",\"quantity_picked\" = EXCLUDED.\"quantity_picked\",\"cost_price\" = EXCLUDED.\"cost_price\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"rid\" = EXCLUDED.\"rid\""
const __sqlInventoryVariant_Insert = "INSERT INTO \"inventory_variant\" (" + __sqlInventoryVariant_ListCols + ") VALUES"
const __sqlInventoryVariant_Select = "SELECT " + __sqlInventoryVariant_ListCols + " FROM \"inventory_variant\""
const __sqlInventoryVariant_Select_history = "SELECT " + __sqlInventoryVariant_ListCols + " FROM history.\"inventory_variant\""
const __sqlInventoryVariant_UpdateAll = "UPDATE \"inventory_variant\" SET (" + __sqlInventoryVariant_ListCols + ")"
const __sqlInventoryVariant_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT inventory_variant_pkey DO UPDATE SET"

func (m *InventoryVariant) SQLTableName() string  { return "inventory_variant" }
func (m *InventoryVariants) SQLTableName() string { return "inventory_variant" }
func (m *InventoryVariant) SQLListCols() string   { return __sqlInventoryVariant_ListCols }

func (m *InventoryVariant) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlInventoryVariant_ListCols + " FROM \"inventory_variant\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *InventoryVariant) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "inventory_variant"); err != nil {
		db.RecordError(err)
		return
	} else {
		mDBColumnNameAndType = val
	}
	mModelColumnNameAndType := map[string]migration.ColumnDef{
		"shop_id": {
			ColumnName:       "shop_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"variant_id": {
			ColumnName:       "variant_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"quantity_on_hand": {
			ColumnName:       "quantity_on_hand",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"quantity_picked": {
			ColumnName:       "quantity_picked",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"cost_price": {
			ColumnName:       "cost_price",
			ColumnType:       "int",
			ColumnDBType:     "int",
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
		"rid": {
			ColumnName:       "rid",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "inventory_variant", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*InventoryVariant)(nil))
}

func (m *InventoryVariant) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ShopID,
		m.VariantID,
		core.Int(m.QuantityOnHand),
		core.Int(m.QuantityPicked),
		core.Int(m.CostPrice),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		m.Rid,
	}
}

func (m *InventoryVariant) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ShopID,
		&m.VariantID,
		(*core.Int)(&m.QuantityOnHand),
		(*core.Int)(&m.QuantityPicked),
		(*core.Int)(&m.CostPrice),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		&m.Rid,
	}
}

func (m *InventoryVariant) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *InventoryVariants) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(InventoryVariants, 0, 128)
	for rows.Next() {
		m := new(InventoryVariant)
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

func (_ *InventoryVariant) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_Select)
	return nil
}

func (_ *InventoryVariants) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_Select)
	return nil
}

func (m *InventoryVariant) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms InventoryVariants) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(8)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *InventoryVariant) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlInventoryVariant_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInventoryVariant_ListColsOnConflict)
	return nil
}

func (ms InventoryVariants) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlInventoryVariant_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInventoryVariant_ListColsOnConflict)
	return nil
}

func (m *InventoryVariant) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("inventory_variant")
	w.WriteRawString(" SET ")
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.VariantID != 0 {
		flag = true
		w.WriteName("variant_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.VariantID)
	}
	if m.QuantityOnHand != 0 {
		flag = true
		w.WriteName("quantity_on_hand")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.QuantityOnHand)
	}
	if m.QuantityPicked != 0 {
		flag = true
		w.WriteName("quantity_picked")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.QuantityPicked)
	}
	if m.CostPrice != 0 {
		flag = true
		w.WriteName("cost_price")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CostPrice)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if true { // always update time
		flag = true
		w.WriteName("updated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Now(m.UpdatedAt, time.Now(), true))
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

func (m *InventoryVariant) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type InventoryVariantHistory map[string]interface{}
type InventoryVariantHistories []map[string]interface{}

func (m *InventoryVariantHistory) SQLTableName() string  { return "history.\"inventory_variant\"" }
func (m InventoryVariantHistories) SQLTableName() string { return "history.\"inventory_variant\"" }

func (m *InventoryVariantHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_Select_history)
	return nil
}

func (m InventoryVariantHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVariant_Select_history)
	return nil
}

func (m InventoryVariantHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m InventoryVariantHistory) VariantID() core.Interface { return core.Interface{m["variant_id"]} }
func (m InventoryVariantHistory) QuantityOnHand() core.Interface {
	return core.Interface{m["quantity_on_hand"]}
}
func (m InventoryVariantHistory) QuantityPicked() core.Interface {
	return core.Interface{m["quantity_picked"]}
}
func (m InventoryVariantHistory) CostPrice() core.Interface { return core.Interface{m["cost_price"]} }
func (m InventoryVariantHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m InventoryVariantHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m InventoryVariantHistory) Rid() core.Interface       { return core.Interface{m["rid"]} }

func (m *InventoryVariantHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(InventoryVariantHistory, 8)
	res["shop_id"] = data[0]
	res["variant_id"] = data[1]
	res["quantity_on_hand"] = data[2]
	res["quantity_picked"] = data[3]
	res["cost_price"] = data[4]
	res["created_at"] = data[5]
	res["updated_at"] = data[6]
	res["rid"] = data[7]
	*m = res
	return nil
}

func (ms *InventoryVariantHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	res := make(InventoryVariantHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(InventoryVariantHistory)
		m["shop_id"] = data[0]
		m["variant_id"] = data[1]
		m["quantity_on_hand"] = data[2]
		m["quantity_picked"] = data[3]
		m["cost_price"] = data[4]
		m["created_at"] = data[5]
		m["updated_at"] = data[6]
		m["rid"] = data[7]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

type InventoryVouchers []*InventoryVoucher

const __sqlInventoryVoucher_Table = "inventory_voucher"
const __sqlInventoryVoucher_ListCols = "\"shop_id\",\"id\",\"created_by\",\"updated_by\",\"code\",\"code_norm\",\"status\",\"trader_id\",\"trader\",\"total_amount\",\"type\",\"lines\",\"variant_ids\",\"ref_id\",\"ref_code\",\"ref_type\",\"title\",\"created_at\",\"updated_at\",\"confirmed_at\",\"cancelled_at\",\"cancel_reason\",\"product_ids\",\"rid\""
const __sqlInventoryVoucher_ListColsOnConflict = "\"shop_id\" = EXCLUDED.\"shop_id\",\"id\" = EXCLUDED.\"id\",\"created_by\" = EXCLUDED.\"created_by\",\"updated_by\" = EXCLUDED.\"updated_by\",\"code\" = EXCLUDED.\"code\",\"code_norm\" = EXCLUDED.\"code_norm\",\"status\" = EXCLUDED.\"status\",\"trader_id\" = EXCLUDED.\"trader_id\",\"trader\" = EXCLUDED.\"trader\",\"total_amount\" = EXCLUDED.\"total_amount\",\"type\" = EXCLUDED.\"type\",\"lines\" = EXCLUDED.\"lines\",\"variant_ids\" = EXCLUDED.\"variant_ids\",\"ref_id\" = EXCLUDED.\"ref_id\",\"ref_code\" = EXCLUDED.\"ref_code\",\"ref_type\" = EXCLUDED.\"ref_type\",\"title\" = EXCLUDED.\"title\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"confirmed_at\" = EXCLUDED.\"confirmed_at\",\"cancelled_at\" = EXCLUDED.\"cancelled_at\",\"cancel_reason\" = EXCLUDED.\"cancel_reason\",\"product_ids\" = EXCLUDED.\"product_ids\",\"rid\" = EXCLUDED.\"rid\""
const __sqlInventoryVoucher_Insert = "INSERT INTO \"inventory_voucher\" (" + __sqlInventoryVoucher_ListCols + ") VALUES"
const __sqlInventoryVoucher_Select = "SELECT " + __sqlInventoryVoucher_ListCols + " FROM \"inventory_voucher\""
const __sqlInventoryVoucher_Select_history = "SELECT " + __sqlInventoryVoucher_ListCols + " FROM history.\"inventory_voucher\""
const __sqlInventoryVoucher_UpdateAll = "UPDATE \"inventory_voucher\" SET (" + __sqlInventoryVoucher_ListCols + ")"
const __sqlInventoryVoucher_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT inventory_voucher_pkey DO UPDATE SET"

func (m *InventoryVoucher) SQLTableName() string  { return "inventory_voucher" }
func (m *InventoryVouchers) SQLTableName() string { return "inventory_voucher" }
func (m *InventoryVoucher) SQLListCols() string   { return __sqlInventoryVoucher_ListCols }

func (m *InventoryVoucher) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlInventoryVoucher_ListCols + " FROM \"inventory_voucher\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *InventoryVoucher) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "inventory_voucher"); err != nil {
		db.RecordError(err)
		return
	} else {
		mDBColumnNameAndType = val
	}
	mModelColumnNameAndType := map[string]migration.ColumnDef{
		"shop_id": {
			ColumnName:       "shop_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"id": {
			ColumnName:       "id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"created_by": {
			ColumnName:       "created_by",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"updated_by": {
			ColumnName:       "updated_by",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"code": {
			ColumnName:       "code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"code_norm": {
			ColumnName:       "code_norm",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"status": {
			ColumnName:       "status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "N"},
		},
		"trader_id": {
			ColumnName:       "trader_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"trader": {
			ColumnName:       "trader",
			ColumnType:       "*Trader",
			ColumnDBType:     "*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"total_amount": {
			ColumnName:       "total_amount",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"type": {
			ColumnName:       "type",
			ColumnType:       "inventory_type.InventoryVoucherType",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"unknown", "in", "out"},
		},
		"lines": {
			ColumnName:       "lines",
			ColumnType:       "[]*InventoryVoucherItem",
			ColumnDBType:     "[]*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"variant_ids": {
			ColumnName:       "variant_ids",
			ColumnType:       "[]dot.ID",
			ColumnDBType:     "[]int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"ref_id": {
			ColumnName:       "ref_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"ref_code": {
			ColumnName:       "ref_code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"ref_type": {
			ColumnName:       "ref_type",
			ColumnType:       "inventory_voucher_ref.InventoryVoucherRef",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"unknown", "refund", "purchase_refund", "stocktake", "purchase_order", "order"},
		},
		"title": {
			ColumnName:       "title",
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
		"confirmed_at": {
			ColumnName:       "confirmed_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"cancelled_at": {
			ColumnName:       "cancelled_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"cancel_reason": {
			ColumnName:       "cancel_reason",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"product_ids": {
			ColumnName:       "product_ids",
			ColumnType:       "[]dot.ID",
			ColumnDBType:     "[]int64",
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
	if err := migration.Compare(db, "inventory_voucher", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*InventoryVoucher)(nil))
}

func (m *InventoryVoucher) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ShopID,
		m.ID,
		m.CreatedBy,
		m.UpdatedBy,
		core.String(m.Code),
		core.Int(m.CodeNorm),
		m.Status,
		m.TraderID,
		core.JSON{m.Trader},
		core.Int(m.TotalAmount),
		m.Type,
		core.JSON{m.Lines},
		core.Array{m.VariantIDs, opts},
		m.RefID,
		core.String(m.RefCode),
		m.RefType,
		core.String(m.Title),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.ConfirmedAt),
		core.Time(m.CancelledAt),
		core.String(m.CancelReason),
		core.Array{m.ProductIDs, opts},
		m.Rid,
	}
}

func (m *InventoryVoucher) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ShopID,
		&m.ID,
		&m.CreatedBy,
		&m.UpdatedBy,
		(*core.String)(&m.Code),
		(*core.Int)(&m.CodeNorm),
		&m.Status,
		&m.TraderID,
		core.JSON{&m.Trader},
		(*core.Int)(&m.TotalAmount),
		&m.Type,
		core.JSON{&m.Lines},
		core.Array{&m.VariantIDs, opts},
		&m.RefID,
		(*core.String)(&m.RefCode),
		&m.RefType,
		(*core.String)(&m.Title),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.ConfirmedAt),
		(*core.Time)(&m.CancelledAt),
		(*core.String)(&m.CancelReason),
		core.Array{&m.ProductIDs, opts},
		&m.Rid,
	}
}

func (m *InventoryVoucher) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *InventoryVouchers) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(InventoryVouchers, 0, 128)
	for rows.Next() {
		m := new(InventoryVoucher)
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

func (_ *InventoryVoucher) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_Select)
	return nil
}

func (_ *InventoryVouchers) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_Select)
	return nil
}

func (m *InventoryVoucher) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(24)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms InventoryVouchers) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(24)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *InventoryVoucher) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlInventoryVoucher_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInventoryVoucher_ListColsOnConflict)
	return nil
}

func (ms InventoryVouchers) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlInventoryVoucher_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlInventoryVoucher_ListColsOnConflict)
	return nil
}

func (m *InventoryVoucher) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("inventory_voucher")
	w.WriteRawString(" SET ")
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
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
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.TraderID != 0 {
		flag = true
		w.WriteName("trader_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TraderID)
	}
	if m.Trader != nil {
		flag = true
		w.WriteName("trader")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Trader})
	}
	if m.TotalAmount != 0 {
		flag = true
		w.WriteName("total_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalAmount)
	}
	if m.Type != 0 {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
	}
	if m.Lines != nil {
		flag = true
		w.WriteName("lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Lines})
	}
	if m.VariantIDs != nil {
		flag = true
		w.WriteName("variant_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.VariantIDs, opts})
	}
	if m.RefID != 0 {
		flag = true
		w.WriteName("ref_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RefID)
	}
	if m.RefCode != "" {
		flag = true
		w.WriteName("ref_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RefCode)
	}
	if m.RefType != 0 {
		flag = true
		w.WriteName("ref_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RefType)
	}
	if m.Title != "" {
		flag = true
		w.WriteName("title")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Title)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if true { // always update time
		flag = true
		w.WriteName("updated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Now(m.UpdatedAt, time.Now(), true))
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
	if m.CancelReason != "" {
		flag = true
		w.WriteName("cancel_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelReason)
	}
	if m.ProductIDs != nil {
		flag = true
		w.WriteName("product_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.ProductIDs, opts})
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

func (m *InventoryVoucher) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(24)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type InventoryVoucherHistory map[string]interface{}
type InventoryVoucherHistories []map[string]interface{}

func (m *InventoryVoucherHistory) SQLTableName() string  { return "history.\"inventory_voucher\"" }
func (m InventoryVoucherHistories) SQLTableName() string { return "history.\"inventory_voucher\"" }

func (m *InventoryVoucherHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_Select_history)
	return nil
}

func (m InventoryVoucherHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlInventoryVoucher_Select_history)
	return nil
}

func (m InventoryVoucherHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m InventoryVoucherHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m InventoryVoucherHistory) CreatedBy() core.Interface { return core.Interface{m["created_by"]} }
func (m InventoryVoucherHistory) UpdatedBy() core.Interface { return core.Interface{m["updated_by"]} }
func (m InventoryVoucherHistory) Code() core.Interface      { return core.Interface{m["code"]} }
func (m InventoryVoucherHistory) CodeNorm() core.Interface  { return core.Interface{m["code_norm"]} }
func (m InventoryVoucherHistory) Status() core.Interface    { return core.Interface{m["status"]} }
func (m InventoryVoucherHistory) TraderID() core.Interface  { return core.Interface{m["trader_id"]} }
func (m InventoryVoucherHistory) Trader() core.Interface    { return core.Interface{m["trader"]} }
func (m InventoryVoucherHistory) TotalAmount() core.Interface {
	return core.Interface{m["total_amount"]}
}
func (m InventoryVoucherHistory) Type() core.Interface       { return core.Interface{m["type"]} }
func (m InventoryVoucherHistory) Lines() core.Interface      { return core.Interface{m["lines"]} }
func (m InventoryVoucherHistory) VariantIDs() core.Interface { return core.Interface{m["variant_ids"]} }
func (m InventoryVoucherHistory) RefID() core.Interface      { return core.Interface{m["ref_id"]} }
func (m InventoryVoucherHistory) RefCode() core.Interface    { return core.Interface{m["ref_code"]} }
func (m InventoryVoucherHistory) RefType() core.Interface    { return core.Interface{m["ref_type"]} }
func (m InventoryVoucherHistory) Title() core.Interface      { return core.Interface{m["title"]} }
func (m InventoryVoucherHistory) CreatedAt() core.Interface  { return core.Interface{m["created_at"]} }
func (m InventoryVoucherHistory) UpdatedAt() core.Interface  { return core.Interface{m["updated_at"]} }
func (m InventoryVoucherHistory) ConfirmedAt() core.Interface {
	return core.Interface{m["confirmed_at"]}
}
func (m InventoryVoucherHistory) CancelledAt() core.Interface {
	return core.Interface{m["cancelled_at"]}
}
func (m InventoryVoucherHistory) CancelReason() core.Interface {
	return core.Interface{m["cancel_reason"]}
}
func (m InventoryVoucherHistory) ProductIDs() core.Interface { return core.Interface{m["product_ids"]} }
func (m InventoryVoucherHistory) Rid() core.Interface        { return core.Interface{m["rid"]} }

func (m *InventoryVoucherHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 24)
	args := make([]interface{}, 24)
	for i := 0; i < 24; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(InventoryVoucherHistory, 24)
	res["shop_id"] = data[0]
	res["id"] = data[1]
	res["created_by"] = data[2]
	res["updated_by"] = data[3]
	res["code"] = data[4]
	res["code_norm"] = data[5]
	res["status"] = data[6]
	res["trader_id"] = data[7]
	res["trader"] = data[8]
	res["total_amount"] = data[9]
	res["type"] = data[10]
	res["lines"] = data[11]
	res["variant_ids"] = data[12]
	res["ref_id"] = data[13]
	res["ref_code"] = data[14]
	res["ref_type"] = data[15]
	res["title"] = data[16]
	res["created_at"] = data[17]
	res["updated_at"] = data[18]
	res["confirmed_at"] = data[19]
	res["cancelled_at"] = data[20]
	res["cancel_reason"] = data[21]
	res["product_ids"] = data[22]
	res["rid"] = data[23]
	*m = res
	return nil
}

func (ms *InventoryVoucherHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 24)
	args := make([]interface{}, 24)
	for i := 0; i < 24; i++ {
		args[i] = &data[i]
	}
	res := make(InventoryVoucherHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(InventoryVoucherHistory)
		m["shop_id"] = data[0]
		m["id"] = data[1]
		m["created_by"] = data[2]
		m["updated_by"] = data[3]
		m["code"] = data[4]
		m["code_norm"] = data[5]
		m["status"] = data[6]
		m["trader_id"] = data[7]
		m["trader"] = data[8]
		m["total_amount"] = data[9]
		m["type"] = data[10]
		m["lines"] = data[11]
		m["variant_ids"] = data[12]
		m["ref_id"] = data[13]
		m["ref_code"] = data[14]
		m["ref_type"] = data[15]
		m["title"] = data[16]
		m["created_at"] = data[17]
		m["updated_at"] = data[18]
		m["confirmed_at"] = data[19]
		m["cancelled_at"] = data[20]
		m["cancel_reason"] = data[21]
		m["product_ids"] = data[22]
		m["rid"] = data[23]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
