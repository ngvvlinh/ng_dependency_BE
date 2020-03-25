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

type ShopProductCollections []*ShopProductCollection

const __sqlShopProductCollection_Table = "shop_product_collection"
const __sqlShopProductCollection_ListCols = "\"shop_id\",\"product_id\",\"collection_id\",\"created_at\",\"updated_at\",\"rid\""
const __sqlShopProductCollection_ListColsOnConflict = "\"shop_id\" = EXCLUDED.\"shop_id\",\"product_id\" = EXCLUDED.\"product_id\",\"collection_id\" = EXCLUDED.\"collection_id\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShopProductCollection_Insert = "INSERT INTO \"shop_product_collection\" (" + __sqlShopProductCollection_ListCols + ") VALUES"
const __sqlShopProductCollection_Select = "SELECT " + __sqlShopProductCollection_ListCols + " FROM \"shop_product_collection\""
const __sqlShopProductCollection_Select_history = "SELECT " + __sqlShopProductCollection_ListCols + " FROM history.\"shop_product_collection\""
const __sqlShopProductCollection_UpdateAll = "UPDATE \"shop_product_collection\" SET (" + __sqlShopProductCollection_ListCols + ")"
const __sqlShopProductCollection_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_product_collection_pkey DO UPDATE SET"

func (m *ShopProductCollection) SQLTableName() string  { return "shop_product_collection" }
func (m *ShopProductCollections) SQLTableName() string { return "shop_product_collection" }
func (m *ShopProductCollection) SQLListCols() string   { return __sqlShopProductCollection_ListCols }

func (m *ShopProductCollection) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopProductCollection_ListCols + " FROM \"shop_product_collection\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopProductCollection)(nil))
}

func (m *ShopProductCollection) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ShopID,
		m.ProductID,
		m.CollectionID,
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		m.Rid,
	}
}

func (m *ShopProductCollection) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ShopID,
		&m.ProductID,
		&m.CollectionID,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		&m.Rid,
	}
}

func (m *ShopProductCollection) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopProductCollections) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopProductCollections, 0, 128)
	for rows.Next() {
		m := new(ShopProductCollection)
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

func (_ *ShopProductCollection) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_Select)
	return nil
}

func (_ *ShopProductCollections) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_Select)
	return nil
}

func (m *ShopProductCollection) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(6)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopProductCollections) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(6)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopProductCollection) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopProductCollection_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopProductCollection_ListColsOnConflict)
	return nil
}

func (ms ShopProductCollections) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopProductCollection_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopProductCollection_ListColsOnConflict)
	return nil
}

func (m *ShopProductCollection) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_product_collection")
	w.WriteRawString(" SET ")
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.ProductID != 0 {
		flag = true
		w.WriteName("product_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ProductID)
	}
	if m.CollectionID != 0 {
		flag = true
		w.WriteName("collection_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CollectionID)
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

func (m *ShopProductCollection) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(6)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopProductCollectionHistory map[string]interface{}
type ShopProductCollectionHistories []map[string]interface{}

func (m *ShopProductCollectionHistory) SQLTableName() string {
	return "history.\"shop_product_collection\""
}
func (m ShopProductCollectionHistories) SQLTableName() string {
	return "history.\"shop_product_collection\""
}

func (m *ShopProductCollectionHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_Select_history)
	return nil
}

func (m ShopProductCollectionHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProductCollection_Select_history)
	return nil
}

func (m ShopProductCollectionHistory) ShopID() core.Interface { return core.Interface{m["shop_id"]} }
func (m ShopProductCollectionHistory) ProductID() core.Interface {
	return core.Interface{m["product_id"]}
}
func (m ShopProductCollectionHistory) CollectionID() core.Interface {
	return core.Interface{m["collection_id"]}
}
func (m ShopProductCollectionHistory) CreatedAt() core.Interface {
	return core.Interface{m["created_at"]}
}
func (m ShopProductCollectionHistory) UpdatedAt() core.Interface {
	return core.Interface{m["updated_at"]}
}
func (m ShopProductCollectionHistory) Rid() core.Interface { return core.Interface{m["rid"]} }

func (m *ShopProductCollectionHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 6)
	args := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopProductCollectionHistory, 6)
	res["shop_id"] = data[0]
	res["product_id"] = data[1]
	res["collection_id"] = data[2]
	res["created_at"] = data[3]
	res["updated_at"] = data[4]
	res["rid"] = data[5]
	*m = res
	return nil
}

func (ms *ShopProductCollectionHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 6)
	args := make([]interface{}, 6)
	for i := 0; i < 6; i++ {
		args[i] = &data[i]
	}
	res := make(ShopProductCollectionHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopProductCollectionHistory)
		m["shop_id"] = data[0]
		m["product_id"] = data[1]
		m["collection_id"] = data[2]
		m["created_at"] = data[3]
		m["updated_at"] = data[4]
		m["rid"] = data[5]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
