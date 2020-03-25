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

type ShopCategories []*ShopCategory

const __sqlShopCategory_Table = "shop_category"
const __sqlShopCategory_ListCols = "\"id\",\"partner_id\",\"shop_id\",\"parent_id\",\"name\",\"status\",\"created_at\",\"updated_at\",\"rid\""
const __sqlShopCategory_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"partner_id\" = EXCLUDED.\"partner_id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"parent_id\" = EXCLUDED.\"parent_id\",\"name\" = EXCLUDED.\"name\",\"status\" = EXCLUDED.\"status\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShopCategory_Insert = "INSERT INTO \"shop_category\" (" + __sqlShopCategory_ListCols + ") VALUES"
const __sqlShopCategory_Select = "SELECT " + __sqlShopCategory_ListCols + " FROM \"shop_category\""
const __sqlShopCategory_Select_history = "SELECT " + __sqlShopCategory_ListCols + " FROM history.\"shop_category\""
const __sqlShopCategory_UpdateAll = "UPDATE \"shop_category\" SET (" + __sqlShopCategory_ListCols + ")"
const __sqlShopCategory_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_category_pkey DO UPDATE SET"

func (m *ShopCategory) SQLTableName() string   { return "shop_category" }
func (m *ShopCategories) SQLTableName() string { return "shop_category" }
func (m *ShopCategory) SQLListCols() string    { return __sqlShopCategory_ListCols }

func (m *ShopCategory) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopCategory_ListCols + " FROM \"shop_category\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopCategory)(nil))
}

func (m *ShopCategory) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ID,
		m.PartnerID,
		m.ShopID,
		m.ParentID,
		core.String(m.Name),
		core.Int(m.Status),
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		m.Rid,
	}
}

func (m *ShopCategory) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.PartnerID,
		&m.ShopID,
		&m.ParentID,
		(*core.String)(&m.Name),
		(*core.Int)(&m.Status),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		&m.Rid,
	}
}

func (m *ShopCategory) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopCategories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopCategories, 0, 128)
	for rows.Next() {
		m := new(ShopCategory)
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

func (_ *ShopCategory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_Select)
	return nil
}

func (_ *ShopCategories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_Select)
	return nil
}

func (m *ShopCategory) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopCategories) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(9)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopCategory) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopCategory_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopCategory_ListColsOnConflict)
	return nil
}

func (ms ShopCategories) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopCategory_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopCategory_ListColsOnConflict)
	return nil
}

func (m *ShopCategory) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_category")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.PartnerID != 0 {
		flag = true
		w.WriteName("partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PartnerID)
	}
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.ParentID != 0 {
		flag = true
		w.WriteName("parent_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ParentID)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
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
		w.WriteArg(m.UpdatedAt)
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

func (m *ShopCategory) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(9)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopCategoryHistory map[string]interface{}
type ShopCategoryHistories []map[string]interface{}

func (m *ShopCategoryHistory) SQLTableName() string  { return "history.\"shop_category\"" }
func (m ShopCategoryHistories) SQLTableName() string { return "history.\"shop_category\"" }

func (m *ShopCategoryHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_Select_history)
	return nil
}

func (m ShopCategoryHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCategory_Select_history)
	return nil
}

func (m ShopCategoryHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m ShopCategoryHistory) PartnerID() core.Interface { return core.Interface{m["partner_id"]} }
func (m ShopCategoryHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m ShopCategoryHistory) ParentID() core.Interface  { return core.Interface{m["parent_id"]} }
func (m ShopCategoryHistory) Name() core.Interface      { return core.Interface{m["name"]} }
func (m ShopCategoryHistory) Status() core.Interface    { return core.Interface{m["status"]} }
func (m ShopCategoryHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m ShopCategoryHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m ShopCategoryHistory) Rid() core.Interface       { return core.Interface{m["rid"]} }

func (m *ShopCategoryHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopCategoryHistory, 9)
	res["id"] = data[0]
	res["partner_id"] = data[1]
	res["shop_id"] = data[2]
	res["parent_id"] = data[3]
	res["name"] = data[4]
	res["status"] = data[5]
	res["created_at"] = data[6]
	res["updated_at"] = data[7]
	res["rid"] = data[8]
	*m = res
	return nil
}

func (ms *ShopCategoryHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 9)
	args := make([]interface{}, 9)
	for i := 0; i < 9; i++ {
		args[i] = &data[i]
	}
	res := make(ShopCategoryHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopCategoryHistory)
		m["id"] = data[0]
		m["partner_id"] = data[1]
		m["shop_id"] = data[2]
		m["parent_id"] = data[3]
		m["name"] = data[4]
		m["status"] = data[5]
		m["created_at"] = data[6]
		m["updated_at"] = data[7]
		m["rid"] = data[8]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
