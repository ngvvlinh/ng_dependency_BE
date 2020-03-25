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

type ShopProducts []*ShopProduct

const __sqlShopProduct_Table = "shop_product"
const __sqlShopProduct_ListCols = "\"external_id\",\"external_code\",\"partner_id\",\"external_brand_id\",\"external_category_id\",\"shop_id\",\"product_id\",\"code\",\"name\",\"description\",\"desc_html\",\"short_desc\",\"image_urls\",\"note\",\"tags\",\"unit\",\"category_id\",\"cost_price\",\"list_price\",\"retail_price\",\"brand_id\",\"status\",\"created_at\",\"updated_at\",\"product_type\",\"meta_fields\",\"rid\""
const __sqlShopProduct_ListColsOnConflict = "\"external_id\" = EXCLUDED.\"external_id\",\"external_code\" = EXCLUDED.\"external_code\",\"partner_id\" = EXCLUDED.\"partner_id\",\"external_brand_id\" = EXCLUDED.\"external_brand_id\",\"external_category_id\" = EXCLUDED.\"external_category_id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"product_id\" = EXCLUDED.\"product_id\",\"code\" = EXCLUDED.\"code\",\"name\" = EXCLUDED.\"name\",\"description\" = EXCLUDED.\"description\",\"desc_html\" = EXCLUDED.\"desc_html\",\"short_desc\" = EXCLUDED.\"short_desc\",\"image_urls\" = EXCLUDED.\"image_urls\",\"note\" = EXCLUDED.\"note\",\"tags\" = EXCLUDED.\"tags\",\"unit\" = EXCLUDED.\"unit\",\"category_id\" = EXCLUDED.\"category_id\",\"cost_price\" = EXCLUDED.\"cost_price\",\"list_price\" = EXCLUDED.\"list_price\",\"retail_price\" = EXCLUDED.\"retail_price\",\"brand_id\" = EXCLUDED.\"brand_id\",\"status\" = EXCLUDED.\"status\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"product_type\" = EXCLUDED.\"product_type\",\"meta_fields\" = EXCLUDED.\"meta_fields\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShopProduct_Insert = "INSERT INTO \"shop_product\" (" + __sqlShopProduct_ListCols + ") VALUES"
const __sqlShopProduct_Select = "SELECT " + __sqlShopProduct_ListCols + " FROM \"shop_product\""
const __sqlShopProduct_Select_history = "SELECT " + __sqlShopProduct_ListCols + " FROM history.\"shop_product\""
const __sqlShopProduct_UpdateAll = "UPDATE \"shop_product\" SET (" + __sqlShopProduct_ListCols + ")"
const __sqlShopProduct_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_product_pkey DO UPDATE SET"

func (m *ShopProduct) SQLTableName() string  { return "shop_product" }
func (m *ShopProducts) SQLTableName() string { return "shop_product" }
func (m *ShopProduct) SQLListCols() string   { return __sqlShopProduct_ListCols }

func (m *ShopProduct) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopProduct_ListCols + " FROM \"shop_product\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopProduct)(nil))
}

func (m *ShopProduct) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		core.String(m.ExternalID),
		core.String(m.ExternalCode),
		m.PartnerID,
		core.String(m.ExternalBrandID),
		core.String(m.ExternalCategoryID),
		m.ShopID,
		m.ProductID,
		core.String(m.Code),
		core.String(m.Name),
		core.String(m.Description),
		core.String(m.DescHTML),
		core.String(m.ShortDesc),
		core.Array{m.ImageURLs, opts},
		core.String(m.Note),
		core.Array{m.Tags, opts},
		core.String(m.Unit),
		m.CategoryID,
		core.Int(m.CostPrice),
		core.Int(m.ListPrice),
		core.Int(m.RetailPrice),
		m.BrandID,
		m.Status,
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		m.ProductType,
		core.JSON{m.MetaFields},
		m.Rid,
	}
}

func (m *ShopProduct) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.String)(&m.ExternalID),
		(*core.String)(&m.ExternalCode),
		&m.PartnerID,
		(*core.String)(&m.ExternalBrandID),
		(*core.String)(&m.ExternalCategoryID),
		&m.ShopID,
		&m.ProductID,
		(*core.String)(&m.Code),
		(*core.String)(&m.Name),
		(*core.String)(&m.Description),
		(*core.String)(&m.DescHTML),
		(*core.String)(&m.ShortDesc),
		core.Array{&m.ImageURLs, opts},
		(*core.String)(&m.Note),
		core.Array{&m.Tags, opts},
		(*core.String)(&m.Unit),
		&m.CategoryID,
		(*core.Int)(&m.CostPrice),
		(*core.Int)(&m.ListPrice),
		(*core.Int)(&m.RetailPrice),
		&m.BrandID,
		&m.Status,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		&m.ProductType,
		core.JSON{&m.MetaFields},
		&m.Rid,
	}
}

func (m *ShopProduct) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopProducts) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopProducts, 0, 128)
	for rows.Next() {
		m := new(ShopProduct)
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

func (_ *ShopProduct) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_Select)
	return nil
}

func (_ *ShopProducts) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_Select)
	return nil
}

func (m *ShopProduct) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(27)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopProducts) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(27)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopProduct) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopProduct_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopProduct_ListColsOnConflict)
	return nil
}

func (ms ShopProducts) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopProduct_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopProduct_ListColsOnConflict)
	return nil
}

func (m *ShopProduct) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_product")
	w.WriteRawString(" SET ")
	if m.ExternalID != "" {
		flag = true
		w.WriteName("external_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalID)
	}
	if m.ExternalCode != "" {
		flag = true
		w.WriteName("external_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalCode)
	}
	if m.PartnerID != 0 {
		flag = true
		w.WriteName("partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PartnerID)
	}
	if m.ExternalBrandID != "" {
		flag = true
		w.WriteName("external_brand_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalBrandID)
	}
	if m.ExternalCategoryID != "" {
		flag = true
		w.WriteName("external_category_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalCategoryID)
	}
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
	if m.Code != "" {
		flag = true
		w.WriteName("code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Code)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.Description != "" {
		flag = true
		w.WriteName("description")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Description)
	}
	if m.DescHTML != "" {
		flag = true
		w.WriteName("desc_html")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DescHTML)
	}
	if m.ShortDesc != "" {
		flag = true
		w.WriteName("short_desc")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShortDesc)
	}
	if m.ImageURLs != nil {
		flag = true
		w.WriteName("image_urls")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.ImageURLs, opts})
	}
	if m.Note != "" {
		flag = true
		w.WriteName("note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Note)
	}
	if m.Tags != nil {
		flag = true
		w.WriteName("tags")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.Tags, opts})
	}
	if m.Unit != "" {
		flag = true
		w.WriteName("unit")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Unit)
	}
	if m.CategoryID != 0 {
		flag = true
		w.WriteName("category_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CategoryID)
	}
	if m.CostPrice != 0 {
		flag = true
		w.WriteName("cost_price")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CostPrice)
	}
	if m.ListPrice != 0 {
		flag = true
		w.WriteName("list_price")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ListPrice)
	}
	if m.RetailPrice != 0 {
		flag = true
		w.WriteName("retail_price")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RetailPrice)
	}
	if m.BrandID != 0 {
		flag = true
		w.WriteName("brand_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.BrandID)
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
	if m.ProductType != 0 {
		flag = true
		w.WriteName("product_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ProductType)
	}
	if m.MetaFields != nil {
		flag = true
		w.WriteName("meta_fields")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.MetaFields})
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

func (m *ShopProduct) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(27)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopProductHistory map[string]interface{}
type ShopProductHistories []map[string]interface{}

func (m *ShopProductHistory) SQLTableName() string  { return "history.\"shop_product\"" }
func (m ShopProductHistories) SQLTableName() string { return "history.\"shop_product\"" }

func (m *ShopProductHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_Select_history)
	return nil
}

func (m ShopProductHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopProduct_Select_history)
	return nil
}

func (m ShopProductHistory) ExternalID() core.Interface   { return core.Interface{m["external_id"]} }
func (m ShopProductHistory) ExternalCode() core.Interface { return core.Interface{m["external_code"]} }
func (m ShopProductHistory) PartnerID() core.Interface    { return core.Interface{m["partner_id"]} }
func (m ShopProductHistory) ExternalBrandID() core.Interface {
	return core.Interface{m["external_brand_id"]}
}
func (m ShopProductHistory) ExternalCategoryID() core.Interface {
	return core.Interface{m["external_category_id"]}
}
func (m ShopProductHistory) ShopID() core.Interface      { return core.Interface{m["shop_id"]} }
func (m ShopProductHistory) ProductID() core.Interface   { return core.Interface{m["product_id"]} }
func (m ShopProductHistory) Code() core.Interface        { return core.Interface{m["code"]} }
func (m ShopProductHistory) Name() core.Interface        { return core.Interface{m["name"]} }
func (m ShopProductHistory) Description() core.Interface { return core.Interface{m["description"]} }
func (m ShopProductHistory) DescHTML() core.Interface    { return core.Interface{m["desc_html"]} }
func (m ShopProductHistory) ShortDesc() core.Interface   { return core.Interface{m["short_desc"]} }
func (m ShopProductHistory) ImageURLs() core.Interface   { return core.Interface{m["image_urls"]} }
func (m ShopProductHistory) Note() core.Interface        { return core.Interface{m["note"]} }
func (m ShopProductHistory) Tags() core.Interface        { return core.Interface{m["tags"]} }
func (m ShopProductHistory) Unit() core.Interface        { return core.Interface{m["unit"]} }
func (m ShopProductHistory) CategoryID() core.Interface  { return core.Interface{m["category_id"]} }
func (m ShopProductHistory) CostPrice() core.Interface   { return core.Interface{m["cost_price"]} }
func (m ShopProductHistory) ListPrice() core.Interface   { return core.Interface{m["list_price"]} }
func (m ShopProductHistory) RetailPrice() core.Interface { return core.Interface{m["retail_price"]} }
func (m ShopProductHistory) BrandID() core.Interface     { return core.Interface{m["brand_id"]} }
func (m ShopProductHistory) Status() core.Interface      { return core.Interface{m["status"]} }
func (m ShopProductHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m ShopProductHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m ShopProductHistory) ProductType() core.Interface { return core.Interface{m["product_type"]} }
func (m ShopProductHistory) MetaFields() core.Interface  { return core.Interface{m["meta_fields"]} }
func (m ShopProductHistory) Rid() core.Interface         { return core.Interface{m["rid"]} }

func (m *ShopProductHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 27)
	args := make([]interface{}, 27)
	for i := 0; i < 27; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopProductHistory, 27)
	res["external_id"] = data[0]
	res["external_code"] = data[1]
	res["partner_id"] = data[2]
	res["external_brand_id"] = data[3]
	res["external_category_id"] = data[4]
	res["shop_id"] = data[5]
	res["product_id"] = data[6]
	res["code"] = data[7]
	res["name"] = data[8]
	res["description"] = data[9]
	res["desc_html"] = data[10]
	res["short_desc"] = data[11]
	res["image_urls"] = data[12]
	res["note"] = data[13]
	res["tags"] = data[14]
	res["unit"] = data[15]
	res["category_id"] = data[16]
	res["cost_price"] = data[17]
	res["list_price"] = data[18]
	res["retail_price"] = data[19]
	res["brand_id"] = data[20]
	res["status"] = data[21]
	res["created_at"] = data[22]
	res["updated_at"] = data[23]
	res["product_type"] = data[24]
	res["meta_fields"] = data[25]
	res["rid"] = data[26]
	*m = res
	return nil
}

func (ms *ShopProductHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 27)
	args := make([]interface{}, 27)
	for i := 0; i < 27; i++ {
		args[i] = &data[i]
	}
	res := make(ShopProductHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopProductHistory)
		m["external_id"] = data[0]
		m["external_code"] = data[1]
		m["partner_id"] = data[2]
		m["external_brand_id"] = data[3]
		m["external_category_id"] = data[4]
		m["shop_id"] = data[5]
		m["product_id"] = data[6]
		m["code"] = data[7]
		m["name"] = data[8]
		m["description"] = data[9]
		m["desc_html"] = data[10]
		m["short_desc"] = data[11]
		m["image_urls"] = data[12]
		m["note"] = data[13]
		m["tags"] = data[14]
		m["unit"] = data[15]
		m["category_id"] = data[16]
		m["cost_price"] = data[17]
		m["list_price"] = data[18]
		m["retail_price"] = data[19]
		m["brand_id"] = data[20]
		m["status"] = data[21]
		m["created_at"] = data[22]
		m["updated_at"] = data[23]
		m["product_type"] = data[24]
		m["meta_fields"] = data[25]
		m["rid"] = data[26]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
