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

type ShopCustomerGroups []*ShopCustomerGroup

const __sqlShopCustomerGroup_Table = "shop_customer_group"
const __sqlShopCustomerGroup_ListCols = "\"id\",\"partner_id\",\"name\",\"shop_id\",\"created_at\",\"updated_at\",\"rid\""
const __sqlShopCustomerGroup_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"partner_id\" = EXCLUDED.\"partner_id\",\"name\" = EXCLUDED.\"name\",\"shop_id\" = EXCLUDED.\"shop_id\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShopCustomerGroup_Insert = "INSERT INTO \"shop_customer_group\" (" + __sqlShopCustomerGroup_ListCols + ") VALUES"
const __sqlShopCustomerGroup_Select = "SELECT " + __sqlShopCustomerGroup_ListCols + " FROM \"shop_customer_group\""
const __sqlShopCustomerGroup_Select_history = "SELECT " + __sqlShopCustomerGroup_ListCols + " FROM history.\"shop_customer_group\""
const __sqlShopCustomerGroup_UpdateAll = "UPDATE \"shop_customer_group\" SET (" + __sqlShopCustomerGroup_ListCols + ")"
const __sqlShopCustomerGroup_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shop_customer_group_pkey DO UPDATE SET"

func (m *ShopCustomerGroup) SQLTableName() string  { return "shop_customer_group" }
func (m *ShopCustomerGroups) SQLTableName() string { return "shop_customer_group" }
func (m *ShopCustomerGroup) SQLListCols() string   { return __sqlShopCustomerGroup_ListCols }

func (m *ShopCustomerGroup) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShopCustomerGroup_ListCols + " FROM \"shop_customer_group\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShopCustomerGroup) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shop_customer_group"); err != nil {
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
		"partner_id": {
			ColumnName:       "partner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"name": {
			ColumnName:       "name",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
	if err := migration.Compare(db, "shop_customer_group", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShopCustomerGroup)(nil))
}

func (m *ShopCustomerGroup) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		m.ID,
		m.PartnerID,
		core.String(m.Name),
		m.ShopID,
		core.Time(m.CreatedAt),
		core.Time(m.UpdatedAt),
		m.Rid,
	}
}

func (m *ShopCustomerGroup) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.PartnerID,
		(*core.String)(&m.Name),
		&m.ShopID,
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		&m.Rid,
	}
}

func (m *ShopCustomerGroup) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShopCustomerGroups) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShopCustomerGroups, 0, 128)
	for rows.Next() {
		m := new(ShopCustomerGroup)
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

func (_ *ShopCustomerGroup) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_Select)
	return nil
}

func (_ *ShopCustomerGroups) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_Select)
	return nil
}

func (m *ShopCustomerGroup) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(7)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShopCustomerGroups) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(7)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShopCustomerGroup) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShopCustomerGroup_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopCustomerGroup_ListColsOnConflict)
	return nil
}

func (ms ShopCustomerGroups) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShopCustomerGroup_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShopCustomerGroup_ListColsOnConflict)
	return nil
}

func (m *ShopCustomerGroup) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shop_customer_group")
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
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
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

func (m *ShopCustomerGroup) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(7)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShopCustomerGroupHistory map[string]interface{}
type ShopCustomerGroupHistories []map[string]interface{}

func (m *ShopCustomerGroupHistory) SQLTableName() string  { return "history.\"shop_customer_group\"" }
func (m ShopCustomerGroupHistories) SQLTableName() string { return "history.\"shop_customer_group\"" }

func (m *ShopCustomerGroupHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_Select_history)
	return nil
}

func (m ShopCustomerGroupHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShopCustomerGroup_Select_history)
	return nil
}

func (m ShopCustomerGroupHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m ShopCustomerGroupHistory) PartnerID() core.Interface { return core.Interface{m["partner_id"]} }
func (m ShopCustomerGroupHistory) Name() core.Interface      { return core.Interface{m["name"]} }
func (m ShopCustomerGroupHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m ShopCustomerGroupHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m ShopCustomerGroupHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m ShopCustomerGroupHistory) Rid() core.Interface       { return core.Interface{m["rid"]} }

func (m *ShopCustomerGroupHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 7)
	args := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShopCustomerGroupHistory, 7)
	res["id"] = data[0]
	res["partner_id"] = data[1]
	res["name"] = data[2]
	res["shop_id"] = data[3]
	res["created_at"] = data[4]
	res["updated_at"] = data[5]
	res["rid"] = data[6]
	*m = res
	return nil
}

func (ms *ShopCustomerGroupHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 7)
	args := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		args[i] = &data[i]
	}
	res := make(ShopCustomerGroupHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShopCustomerGroupHistory)
		m["id"] = data[0]
		m["partner_id"] = data[1]
		m["name"] = data[2]
		m["shop_id"] = data[3]
		m["created_at"] = data[4]
		m["updated_at"] = data[5]
		m["rid"] = data[6]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
