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

type ShipmentPriceLists []*ShipmentPriceList

const __sqlShipmentPriceList_Table = "shipment_price_list"
const __sqlShipmentPriceList_ListCols = "\"id\",\"name\",\"description\",\"is_active\",\"created_at\",\"updated_at\",\"deleted_at\",\"wl_partner_id\""
const __sqlShipmentPriceList_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"name\" = EXCLUDED.\"name\",\"description\" = EXCLUDED.\"description\",\"is_active\" = EXCLUDED.\"is_active\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"wl_partner_id\" = EXCLUDED.\"wl_partner_id\""
const __sqlShipmentPriceList_Insert = "INSERT INTO \"shipment_price_list\" (" + __sqlShipmentPriceList_ListCols + ") VALUES"
const __sqlShipmentPriceList_Select = "SELECT " + __sqlShipmentPriceList_ListCols + " FROM \"shipment_price_list\""
const __sqlShipmentPriceList_Select_history = "SELECT " + __sqlShipmentPriceList_ListCols + " FROM history.\"shipment_price_list\""
const __sqlShipmentPriceList_UpdateAll = "UPDATE \"shipment_price_list\" SET (" + __sqlShipmentPriceList_ListCols + ")"
const __sqlShipmentPriceList_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shipment_price_list_pkey DO UPDATE SET"

func (m *ShipmentPriceList) SQLTableName() string  { return "shipment_price_list" }
func (m *ShipmentPriceLists) SQLTableName() string { return "shipment_price_list" }
func (m *ShipmentPriceList) SQLListCols() string   { return __sqlShipmentPriceList_ListCols }

func (m *ShipmentPriceList) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShipmentPriceList_ListCols + " FROM \"shipment_price_list\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShipmentPriceList) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shipment_price_list"); err != nil {
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
		"name": {
			ColumnName:       "name",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"description": {
			ColumnName:       "description",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"is_active": {
			ColumnName:       "is_active",
			ColumnType:       "bool",
			ColumnDBType:     "bool",
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
		"deleted_at": {
			ColumnName:       "deleted_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"wl_partner_id": {
			ColumnName:       "wl_partner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
	}
	if err := migration.Compare(db, "shipment_price_list", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShipmentPriceList)(nil))
}

func (m *ShipmentPriceList) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		core.String(m.Name),
		core.String(m.Description),
		core.Bool(m.IsActive),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		m.WLPartnerID,
	}
}

func (m *ShipmentPriceList) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		(*core.String)(&m.Name),
		(*core.String)(&m.Description),
		(*core.Bool)(&m.IsActive),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		&m.WLPartnerID,
	}
}

func (m *ShipmentPriceList) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShipmentPriceLists) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShipmentPriceLists, 0, 128)
	for rows.Next() {
		m := new(ShipmentPriceList)
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

func (_ *ShipmentPriceList) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_Select)
	return nil
}

func (_ *ShipmentPriceLists) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_Select)
	return nil
}

func (m *ShipmentPriceList) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShipmentPriceLists) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(8)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShipmentPriceList) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShipmentPriceList_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipmentPriceList_ListColsOnConflict)
	return nil
}

func (ms ShipmentPriceLists) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShipmentPriceList_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipmentPriceList_ListColsOnConflict)
	return nil
}

func (m *ShipmentPriceList) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shipment_price_list")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
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
	if m.IsActive {
		flag = true
		w.WriteName("is_active")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IsActive)
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
	if m.WLPartnerID != 0 {
		flag = true
		w.WriteName("wl_partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.WLPartnerID)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShipmentPriceList) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(8)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShipmentPriceListHistory map[string]interface{}
type ShipmentPriceListHistories []map[string]interface{}

func (m *ShipmentPriceListHistory) SQLTableName() string  { return "history.\"shipment_price_list\"" }
func (m ShipmentPriceListHistories) SQLTableName() string { return "history.\"shipment_price_list\"" }

func (m *ShipmentPriceListHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_Select_history)
	return nil
}

func (m ShipmentPriceListHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceList_Select_history)
	return nil
}

func (m ShipmentPriceListHistory) ID() core.Interface   { return core.Interface{m["id"]} }
func (m ShipmentPriceListHistory) Name() core.Interface { return core.Interface{m["name"]} }
func (m ShipmentPriceListHistory) Description() core.Interface {
	return core.Interface{m["description"]}
}
func (m ShipmentPriceListHistory) IsActive() core.Interface  { return core.Interface{m["is_active"]} }
func (m ShipmentPriceListHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m ShipmentPriceListHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m ShipmentPriceListHistory) DeletedAt() core.Interface { return core.Interface{m["deleted_at"]} }
func (m ShipmentPriceListHistory) WLPartnerID() core.Interface {
	return core.Interface{m["wl_partner_id"]}
}

func (m *ShipmentPriceListHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShipmentPriceListHistory, 8)
	res["id"] = data[0]
	res["name"] = data[1]
	res["description"] = data[2]
	res["is_active"] = data[3]
	res["created_at"] = data[4]
	res["updated_at"] = data[5]
	res["deleted_at"] = data[6]
	res["wl_partner_id"] = data[7]
	*m = res
	return nil
}

func (ms *ShipmentPriceListHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 8)
	args := make([]interface{}, 8)
	for i := 0; i < 8; i++ {
		args[i] = &data[i]
	}
	res := make(ShipmentPriceListHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShipmentPriceListHistory)
		m["id"] = data[0]
		m["name"] = data[1]
		m["description"] = data[2]
		m["is_active"] = data[3]
		m["created_at"] = data[4]
		m["updated_at"] = data[5]
		m["deleted_at"] = data[6]
		m["wl_partner_id"] = data[7]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
