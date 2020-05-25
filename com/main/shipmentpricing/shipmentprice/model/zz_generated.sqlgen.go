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

type ShipmentPrices []*ShipmentPrice

const __sqlShipmentPrice_Table = "shipment_price"
const __sqlShipmentPrice_ListCols = "\"id\",\"shipment_sub_price_list_id\",\"shipment_service_id\",\"name\",\"custom_region_types\",\"custom_region_ids\",\"region_types\",\"province_types\",\"urban_types\",\"details\",\"priority_point\",\"created_at\",\"updated_at\",\"deleted_at\",\"wl_partner_id\",\"status\""
const __sqlShipmentPrice_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"shipment_sub_price_list_id\" = EXCLUDED.\"shipment_sub_price_list_id\",\"shipment_service_id\" = EXCLUDED.\"shipment_service_id\",\"name\" = EXCLUDED.\"name\",\"custom_region_types\" = EXCLUDED.\"custom_region_types\",\"custom_region_ids\" = EXCLUDED.\"custom_region_ids\",\"region_types\" = EXCLUDED.\"region_types\",\"province_types\" = EXCLUDED.\"province_types\",\"urban_types\" = EXCLUDED.\"urban_types\",\"details\" = EXCLUDED.\"details\",\"priority_point\" = EXCLUDED.\"priority_point\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"wl_partner_id\" = EXCLUDED.\"wl_partner_id\",\"status\" = EXCLUDED.\"status\""
const __sqlShipmentPrice_Insert = "INSERT INTO \"shipment_price\" (" + __sqlShipmentPrice_ListCols + ") VALUES"
const __sqlShipmentPrice_Select = "SELECT " + __sqlShipmentPrice_ListCols + " FROM \"shipment_price\""
const __sqlShipmentPrice_Select_history = "SELECT " + __sqlShipmentPrice_ListCols + " FROM history.\"shipment_price\""
const __sqlShipmentPrice_UpdateAll = "UPDATE \"shipment_price\" SET (" + __sqlShipmentPrice_ListCols + ")"
const __sqlShipmentPrice_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shipment_price_pkey DO UPDATE SET"

func (m *ShipmentPrice) SQLTableName() string  { return "shipment_price" }
func (m *ShipmentPrices) SQLTableName() string { return "shipment_price" }
func (m *ShipmentPrice) SQLListCols() string   { return __sqlShipmentPrice_ListCols }

func (m *ShipmentPrice) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShipmentPrice_ListCols + " FROM \"shipment_price\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShipmentPrice) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shipment_price"); err != nil {
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
		"shipment_sub_price_list_id": {
			ColumnName:       "shipment_sub_price_list_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipment_service_id": {
			ColumnName:       "shipment_service_id",
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
		"custom_region_types": {
			ColumnName:       "custom_region_types",
			ColumnType:       "[]route_type.CustomRegionRouteType",
			ColumnDBType:     "[]enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"custom_region_ids": {
			ColumnName:       "custom_region_ids",
			ColumnType:       "[]dot.ID",
			ColumnDBType:     "[]int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"region_types": {
			ColumnName:       "region_types",
			ColumnType:       "[]route_type.RegionRouteType",
			ColumnDBType:     "[]enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"province_types": {
			ColumnName:       "province_types",
			ColumnType:       "[]route_type.ProvinceRouteType",
			ColumnDBType:     "[]enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"urban_types": {
			ColumnName:       "urban_types",
			ColumnType:       "[]route_type.UrbanType",
			ColumnDBType:     "[]enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"details": {
			ColumnName:       "details",
			ColumnType:       "[]*PricingDetail",
			ColumnDBType:     "[]*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"priority_point": {
			ColumnName:       "priority_point",
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
		"status": {
			ColumnName:       "status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "N"},
		},
	}
	if err := migration.Compare(db, "shipment_price", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShipmentPrice)(nil))
}

func (m *ShipmentPrice) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.ShipmentSubPriceListID,
		m.ShipmentServiceID,
		core.String(m.Name),
		core.Array{m.CustomRegionTypes, opts},
		core.Array{m.CustomRegionIDs, opts},
		core.Array{m.RegionTypes, opts},
		core.Array{m.ProvinceTypes, opts},
		core.Array{m.UrbanTypes, opts},
		core.JSON{m.Details},
		core.Int(m.PriorityPoint),
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		m.WLPartnerID,
		m.Status,
	}
}

func (m *ShipmentPrice) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShipmentSubPriceListID,
		&m.ShipmentServiceID,
		(*core.String)(&m.Name),
		core.Array{&m.CustomRegionTypes, opts},
		core.Array{&m.CustomRegionIDs, opts},
		core.Array{&m.RegionTypes, opts},
		core.Array{&m.ProvinceTypes, opts},
		core.Array{&m.UrbanTypes, opts},
		core.JSON{&m.Details},
		(*core.Int)(&m.PriorityPoint),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		&m.WLPartnerID,
		&m.Status,
	}
}

func (m *ShipmentPrice) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShipmentPrices) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShipmentPrices, 0, 128)
	for rows.Next() {
		m := new(ShipmentPrice)
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

func (_ *ShipmentPrice) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_Select)
	return nil
}

func (_ *ShipmentPrices) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_Select)
	return nil
}

func (m *ShipmentPrice) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(16)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShipmentPrices) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(16)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShipmentPrice) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShipmentPrice_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipmentPrice_ListColsOnConflict)
	return nil
}

func (ms ShipmentPrices) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShipmentPrice_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipmentPrice_ListColsOnConflict)
	return nil
}

func (m *ShipmentPrice) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shipment_price")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.ShipmentSubPriceListID != 0 {
		flag = true
		w.WriteName("shipment_sub_price_list_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShipmentSubPriceListID)
	}
	if m.ShipmentServiceID != 0 {
		flag = true
		w.WriteName("shipment_service_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShipmentServiceID)
	}
	if m.Name != "" {
		flag = true
		w.WriteName("name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Name)
	}
	if m.CustomRegionTypes != nil {
		flag = true
		w.WriteName("custom_region_types")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.CustomRegionTypes, opts})
	}
	if m.CustomRegionIDs != nil {
		flag = true
		w.WriteName("custom_region_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.CustomRegionIDs, opts})
	}
	if m.RegionTypes != nil {
		flag = true
		w.WriteName("region_types")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.RegionTypes, opts})
	}
	if m.ProvinceTypes != nil {
		flag = true
		w.WriteName("province_types")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.ProvinceTypes, opts})
	}
	if m.UrbanTypes != nil {
		flag = true
		w.WriteName("urban_types")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.UrbanTypes, opts})
	}
	if m.Details != nil {
		flag = true
		w.WriteName("details")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Details})
	}
	if m.PriorityPoint != 0 {
		flag = true
		w.WriteName("priority_point")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PriorityPoint)
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
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShipmentPrice) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(16)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShipmentPriceHistory map[string]interface{}
type ShipmentPriceHistories []map[string]interface{}

func (m *ShipmentPriceHistory) SQLTableName() string  { return "history.\"shipment_price\"" }
func (m ShipmentPriceHistories) SQLTableName() string { return "history.\"shipment_price\"" }

func (m *ShipmentPriceHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_Select_history)
	return nil
}

func (m ShipmentPriceHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPrice_Select_history)
	return nil
}

func (m ShipmentPriceHistory) ID() core.Interface { return core.Interface{m["id"]} }
func (m ShipmentPriceHistory) ShipmentSubPriceListID() core.Interface {
	return core.Interface{m["shipment_sub_price_list_id"]}
}
func (m ShipmentPriceHistory) ShipmentServiceID() core.Interface {
	return core.Interface{m["shipment_service_id"]}
}
func (m ShipmentPriceHistory) Name() core.Interface { return core.Interface{m["name"]} }
func (m ShipmentPriceHistory) CustomRegionTypes() core.Interface {
	return core.Interface{m["custom_region_types"]}
}
func (m ShipmentPriceHistory) CustomRegionIDs() core.Interface {
	return core.Interface{m["custom_region_ids"]}
}
func (m ShipmentPriceHistory) RegionTypes() core.Interface { return core.Interface{m["region_types"]} }
func (m ShipmentPriceHistory) ProvinceTypes() core.Interface {
	return core.Interface{m["province_types"]}
}
func (m ShipmentPriceHistory) UrbanTypes() core.Interface { return core.Interface{m["urban_types"]} }
func (m ShipmentPriceHistory) Details() core.Interface    { return core.Interface{m["details"]} }
func (m ShipmentPriceHistory) PriorityPoint() core.Interface {
	return core.Interface{m["priority_point"]}
}
func (m ShipmentPriceHistory) CreatedAt() core.Interface   { return core.Interface{m["created_at"]} }
func (m ShipmentPriceHistory) UpdatedAt() core.Interface   { return core.Interface{m["updated_at"]} }
func (m ShipmentPriceHistory) DeletedAt() core.Interface   { return core.Interface{m["deleted_at"]} }
func (m ShipmentPriceHistory) WLPartnerID() core.Interface { return core.Interface{m["wl_partner_id"]} }
func (m ShipmentPriceHistory) Status() core.Interface      { return core.Interface{m["status"]} }

func (m *ShipmentPriceHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 16)
	args := make([]interface{}, 16)
	for i := 0; i < 16; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShipmentPriceHistory, 16)
	res["id"] = data[0]
	res["shipment_sub_price_list_id"] = data[1]
	res["shipment_service_id"] = data[2]
	res["name"] = data[3]
	res["custom_region_types"] = data[4]
	res["custom_region_ids"] = data[5]
	res["region_types"] = data[6]
	res["province_types"] = data[7]
	res["urban_types"] = data[8]
	res["details"] = data[9]
	res["priority_point"] = data[10]
	res["created_at"] = data[11]
	res["updated_at"] = data[12]
	res["deleted_at"] = data[13]
	res["wl_partner_id"] = data[14]
	res["status"] = data[15]
	*m = res
	return nil
}

func (ms *ShipmentPriceHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 16)
	args := make([]interface{}, 16)
	for i := 0; i < 16; i++ {
		args[i] = &data[i]
	}
	res := make(ShipmentPriceHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShipmentPriceHistory)
		m["id"] = data[0]
		m["shipment_sub_price_list_id"] = data[1]
		m["shipment_service_id"] = data[2]
		m["name"] = data[3]
		m["custom_region_types"] = data[4]
		m["custom_region_ids"] = data[5]
		m["region_types"] = data[6]
		m["province_types"] = data[7]
		m["urban_types"] = data[8]
		m["details"] = data[9]
		m["priority_point"] = data[10]
		m["created_at"] = data[11]
		m["updated_at"] = data[12]
		m["deleted_at"] = data[13]
		m["wl_partner_id"] = data[14]
		m["status"] = data[15]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
