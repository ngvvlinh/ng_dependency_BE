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

type ShopStocktakes []*ShopStocktake

const __sqlShopStocktake_Table = "shop_stocktake"
const __sqlShopStocktake_ListCols = "\"id\",\"shop_id\",\"total_quantity\",\"created_by\",\"updated_by\",\"cancel_reason\",\"type\",\"code\",\"code_norm\",\"status\",\"created_at\",\"updated_at\",\"confirmed_at\",\"cancelled_at\",\"lines\",\"note\",\"product_ids\",\"rid\""
const __sqlShopStocktake_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"total_quantity\" = EXCLUDED.\"total_quantity\",\"created_by\" = EXCLUDED.\"created_by\",\"updated_by\" = EXCLUDED.\"updated_by\",\"cancel_reason\" = EXCLUDED.\"cancel_reason\",\"type\" = EXCLUDED.\"type\",\"code\" = EXCLUDED.\"code\",\"code_norm\" = EXCLUDED.\"code_norm\",\"status\" = EXCLUDED.\"status\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"confirmed_at\" = EXCLUDED.\"confirmed_at\",\"cancelled_at\" = EXCLUDED.\"cancelled_at\",\"lines\" = EXCLUDED.\"lines\",\"note\" = EXCLUDED.\"note\",\"product_ids\" = EXCLUDED.\"product_ids\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShopStocktake_Insert = "INSERT INTO \"shop_stocktake\" (" + __sqlShopStocktake_ListCols + ") VALUES"
const __sqlShopStocktake_Select = "SELECT " + __sqlShopStocktake_ListCols + " FROM \"shop_stocktake\""
const __sqlShopStocktake_Select_history = "SELECT " + __sqlShopStocktake_ListCols + " FROM history.\"shop_stocktake\""
const __sqlShopStocktake_UpdateAll = "UPDATE \"shop_stocktake\" SET (" + __sqlShopStocktake_ListCols + ")"
const __sqlShopStocktake_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_stocktake_pkey DO UPDATE SET"

func (m *ShopStocktake) SQLTableName() string  { return "shop_stocktake" }
func (m *ShopStocktakes) SQLTableName() string { return "shop_stocktake" }
func (m *ShopStocktake) SQLListCols() string   { return __sqlShopStocktake_ListCols }

func (m *ShopStocktake) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopStocktake_ListCols + " FROM \"shop_stocktake\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShopStocktake) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shop_stocktake"); err != nil {
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
		"shop_id": {
			ColumnName:       "shop_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"total_quantity": {
			ColumnName:       "total_quantity",
			ColumnType:       "int",
			ColumnDBType:     "int",
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
		"cancel_reason": {
			ColumnName:       "cancel_reason",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"type": {
			ColumnName:       "type",
			ColumnType:       "stocktake_type.StocktakeType",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"balance", "discard"},
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
		"lines": {
			ColumnName:       "lines",
			ColumnType:       "[]*StocktakeLine",
			ColumnDBType:     "[]*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"note": {
			ColumnName:       "note",
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
	if err := migration.Compare(db, "shop_stocktake", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopStocktake)(nil))
}

func (m *ShopStocktake) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.ShopID,
		core.Int(m.TotalQuantity),
		m.CreatedBy,
		m.UpdatedBy,
		core.String(m.CancelReason),
		m.Type,
		core.String(m.Code),
		core.Int(m.CodeNorm),
		m.Status,
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.ConfirmedAt),
		core.Time(m.CancelledAt),
		core.JSON{m.Lines},
		core.String(m.Note),
		core.Array{m.ProductIDs, opts},
		m.Rid,
	}
}

func (m *ShopStocktake) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShopID,
		(*core.Int)(&m.TotalQuantity),
		&m.CreatedBy,
		&m.UpdatedBy,
		(*core.String)(&m.CancelReason),
		&m.Type,
		(*core.String)(&m.Code),
		(*core.Int)(&m.CodeNorm),
		&m.Status,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.ConfirmedAt),
		(*core.Time)(&m.CancelledAt),
		core.JSON{&m.Lines},
		(*core.String)(&m.Note),
		core.Array{&m.ProductIDs, opts},
		&m.Rid,
	}
}

func (m *ShopStocktake) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopStocktakes) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopStocktakes, 0, 128)
	for rows.Next() {
		m := new(ShopStocktake)
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

func (_ *ShopStocktake) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_Select)
	return nil
}

func (_ *ShopStocktakes) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_Select)
	return nil
}

func (m *ShopStocktake) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(18)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopStocktakes) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(18)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopStocktake) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopStocktake_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopStocktake_ListColsOnConflict)
	return nil
}

func (ms ShopStocktakes) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopStocktake_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopStocktake_ListColsOnConflict)
	return nil
}

func (m *ShopStocktake) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_stocktake")
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
	if m.TotalQuantity != 0 {
		flag = true
		w.WriteName("total_quantity")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalQuantity)
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
	if m.Type != 0 {
		flag = true
		w.WriteName("type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Type)
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
	if m.Lines != nil {
		flag = true
		w.WriteName("lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Lines})
	}
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
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

func (m *ShopStocktake) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(18)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopStocktakeHistory map[string]interface{}
type ShopStocktakeHistories []map[string]interface{}

func (m *ShopStocktakeHistory) SQLTableName() string  { return "history.\"shop_stocktake\"" }
func (m ShopStocktakeHistories) SQLTableName() string { return "history.\"shop_stocktake\"" }

func (m *ShopStocktakeHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_Select_history)
	return nil
}

func (m ShopStocktakeHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopStocktake_Select_history)
	return nil
}

func (m ShopStocktakeHistory) ID() core.Interface     { return core.Interface{m["id"]} }
func (m ShopStocktakeHistory) ShopID() core.Interface { return core.Interface{m["shop_id"]} }
func (m ShopStocktakeHistory) TotalQuantity() core.Interface {
	return core.Interface{m["total_quantity"]}
}
func (m ShopStocktakeHistory) CreatedBy() core.Interface { return core.Interface{m["created_by"]} }
func (m ShopStocktakeHistory) UpdatedBy() core.Interface { return core.Interface{m["updated_by"]} }
func (m ShopStocktakeHistory) CancelReason() core.Interface {
	return core.Interface{m["cancel_reason"]}
}
func (m ShopStocktakeHistory) Type() core.Interface        { return core.Interface{m["type"]} }
func (m ShopStocktakeHistory) Code() core.Interface        { return core.Interface{m["code"]} }
func (m ShopStocktakeHistory) CodeNorm() core.Interface    { return core.Interface{m["code_norm"]} }
func (m ShopStocktakeHistory) Status() core.Interface      { return core.Interface{m["status"]} }
func (m ShopStocktakeHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m ShopStocktakeHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m ShopStocktakeHistory) ConfirmedAt() core.Interface { return core.Interface{m["confirmed_at"]} }
func (m ShopStocktakeHistory) CancelledAt() core.Interface { return core.Interface{m["cancelled_at"]} }
func (m ShopStocktakeHistory) Lines() core.Interface       { return core.Interface{m["lines"]} }
func (m ShopStocktakeHistory) Note() core.Interface        { return core.Interface{m["note"]} }
func (m ShopStocktakeHistory) ProductIDs() core.Interface  { return core.Interface{m["product_ids"]} }
func (m ShopStocktakeHistory) Rid() core.Interface         { return core.Interface{m["rid"]} }

func (m *ShopStocktakeHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 18)
	args := make([]interface{}, 18)
	for i := 0; i < 18; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopStocktakeHistory, 18)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["total_quantity"] = data[2]
	res["created_by"] = data[3]
	res["updated_by"] = data[4]
	res["cancel_reason"] = data[5]
	res["type"] = data[6]
	res["code"] = data[7]
	res["code_norm"] = data[8]
	res["status"] = data[9]
	res["created_at"] = data[10]
	res["updated_at"] = data[11]
	res["confirmed_at"] = data[12]
	res["cancelled_at"] = data[13]
	res["lines"] = data[14]
	res["note"] = data[15]
	res["product_ids"] = data[16]
	res["rid"] = data[17]
	*m = res
	return nil
}

func (ms *ShopStocktakeHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 18)
	args := make([]interface{}, 18)
	for i := 0; i < 18; i++ {
		args[i] = &data[i]
	}
	res := make(ShopStocktakeHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopStocktakeHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["total_quantity"] = data[2]
		m["created_by"] = data[3]
		m["updated_by"] = data[4]
		m["cancel_reason"] = data[5]
		m["type"] = data[6]
		m["code"] = data[7]
		m["code_norm"] = data[8]
		m["status"] = data[9]
		m["created_at"] = data[10]
		m["updated_at"] = data[11]
		m["confirmed_at"] = data[12]
		m["cancelled_at"] = data[13]
		m["lines"] = data[14]
		m["note"] = data[15]
		m["product_ids"] = data[16]
		m["rid"] = data[17]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
