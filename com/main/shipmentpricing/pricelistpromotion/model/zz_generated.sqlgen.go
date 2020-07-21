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

type ShipmentPriceListPromotions []*ShipmentPriceListPromotion

const __sqlShipmentPriceListPromotion_Table = "shipment_price_list_promotion"
const __sqlShipmentPriceListPromotion_ListCols = "\"id\",\"price_list_id\",\"name\",\"description\",\"status\",\"date_from\",\"date_to\",\"applied_rules\",\"created_at\",\"updated_at\",\"deleted_at\",\"wl_partner_id\",\"connection_id\",\"priority_point\""
const __sqlShipmentPriceListPromotion_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"price_list_id\" = EXCLUDED.\"price_list_id\",\"name\" = EXCLUDED.\"name\",\"description\" = EXCLUDED.\"description\",\"status\" = EXCLUDED.\"status\",\"date_from\" = EXCLUDED.\"date_from\",\"date_to\" = EXCLUDED.\"date_to\",\"applied_rules\" = EXCLUDED.\"applied_rules\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"deleted_at\" = EXCLUDED.\"deleted_at\",\"wl_partner_id\" = EXCLUDED.\"wl_partner_id\",\"connection_id\" = EXCLUDED.\"connection_id\",\"priority_point\" = EXCLUDED.\"priority_point\""
const __sqlShipmentPriceListPromotion_Insert = "INSERT INTO \"shipment_price_list_promotion\" (" + __sqlShipmentPriceListPromotion_ListCols + ") VALUES"
const __sqlShipmentPriceListPromotion_Select = "SELECT " + __sqlShipmentPriceListPromotion_ListCols + " FROM \"shipment_price_list_promotion\""
const __sqlShipmentPriceListPromotion_Select_history = "SELECT " + __sqlShipmentPriceListPromotion_ListCols + " FROM history.\"shipment_price_list_promotion\""
const __sqlShipmentPriceListPromotion_UpdateAll = "UPDATE \"shipment_price_list_promotion\" SET (" + __sqlShipmentPriceListPromotion_ListCols + ")"
const __sqlShipmentPriceListPromotion_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shipment_price_list_promotion_pkey DO UPDATE SET"

func (m *ShipmentPriceListPromotion) SQLTableName() string  { return "shipment_price_list_promotion" }
func (m *ShipmentPriceListPromotions) SQLTableName() string { return "shipment_price_list_promotion" }
func (m *ShipmentPriceListPromotion) SQLListCols() string {
	return __sqlShipmentPriceListPromotion_ListCols
}

func (m *ShipmentPriceListPromotion) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShipmentPriceListPromotion_ListCols + " FROM \"shipment_price_list_promotion\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShipmentPriceListPromotion) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shipment_price_list_promotion"); err != nil {
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
		"price_list_id": {
			ColumnName:       "price_list_id",
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
		"status": {
			ColumnName:       "status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "N"},
		},
		"date_from": {
			ColumnName:       "date_from",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"date_to": {
			ColumnName:       "date_to",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"applied_rules": {
			ColumnName:       "applied_rules",
			ColumnType:       "*AppliedRules",
			ColumnDBType:     "*struct",
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
		"connection_id": {
			ColumnName:       "connection_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
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
	}
	if err := migration.Compare(db, "shipment_price_list_promotion", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShipmentPriceListPromotion)(nil))
}

func (m *ShipmentPriceListPromotion) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.PriceListID,
		core.String(m.Name),
		core.String(m.Description),
		m.Status,
		core.Time(m.DateFrom),
		core.Time(m.DateTo),
		core.JSON{m.AppliedRules},
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.DeletedAt),
		m.WLPartnerID,
		m.ConnectionID,
		core.Int(m.PriorityPoint),
	}
}

func (m *ShipmentPriceListPromotion) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.PriceListID,
		(*core.String)(&m.Name),
		(*core.String)(&m.Description),
		&m.Status,
		(*core.Time)(&m.DateFrom),
		(*core.Time)(&m.DateTo),
		core.JSON{&m.AppliedRules},
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.DeletedAt),
		&m.WLPartnerID,
		&m.ConnectionID,
		(*core.Int)(&m.PriorityPoint),
	}
}

func (m *ShipmentPriceListPromotion) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShipmentPriceListPromotions) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShipmentPriceListPromotions, 0, 128)
	for rows.Next() {
		m := new(ShipmentPriceListPromotion)
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

func (_ *ShipmentPriceListPromotion) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_Select)
	return nil
}

func (_ *ShipmentPriceListPromotions) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_Select)
	return nil
}

func (m *ShipmentPriceListPromotion) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(14)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShipmentPriceListPromotions) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(14)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShipmentPriceListPromotion) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShipmentPriceListPromotion_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipmentPriceListPromotion_ListColsOnConflict)
	return nil
}

func (ms ShipmentPriceListPromotions) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShipmentPriceListPromotion_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipmentPriceListPromotion_ListColsOnConflict)
	return nil
}

func (m *ShipmentPriceListPromotion) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shipment_price_list_promotion")
	w.WriteRawString(" SET ")
	if m.ID != 0 {
		flag = true
		w.WriteName("id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ID)
	}
	if m.PriceListID != 0 {
		flag = true
		w.WriteName("price_list_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PriceListID)
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
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if !m.DateFrom.IsZero() {
		flag = true
		w.WriteName("date_from")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DateFrom)
	}
	if !m.DateTo.IsZero() {
		flag = true
		w.WriteName("date_to")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DateTo)
	}
	if m.AppliedRules != nil {
		flag = true
		w.WriteName("applied_rules")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.AppliedRules})
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
	if m.ConnectionID != 0 {
		flag = true
		w.WriteName("connection_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionID)
	}
	if m.PriorityPoint != 0 {
		flag = true
		w.WriteName("priority_point")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PriorityPoint)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShipmentPriceListPromotion) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(14)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShipmentPriceListPromotionHistory map[string]interface{}
type ShipmentPriceListPromotionHistories []map[string]interface{}

func (m *ShipmentPriceListPromotionHistory) SQLTableName() string {
	return "history.\"shipment_price_list_promotion\""
}
func (m ShipmentPriceListPromotionHistories) SQLTableName() string {
	return "history.\"shipment_price_list_promotion\""
}

func (m *ShipmentPriceListPromotionHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_Select_history)
	return nil
}

func (m ShipmentPriceListPromotionHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipmentPriceListPromotion_Select_history)
	return nil
}

func (m ShipmentPriceListPromotionHistory) ID() core.Interface { return core.Interface{m["id"]} }
func (m ShipmentPriceListPromotionHistory) PriceListID() core.Interface {
	return core.Interface{m["price_list_id"]}
}
func (m ShipmentPriceListPromotionHistory) Name() core.Interface { return core.Interface{m["name"]} }
func (m ShipmentPriceListPromotionHistory) Description() core.Interface {
	return core.Interface{m["description"]}
}
func (m ShipmentPriceListPromotionHistory) Status() core.Interface {
	return core.Interface{m["status"]}
}
func (m ShipmentPriceListPromotionHistory) DateFrom() core.Interface {
	return core.Interface{m["date_from"]}
}
func (m ShipmentPriceListPromotionHistory) DateTo() core.Interface {
	return core.Interface{m["date_to"]}
}
func (m ShipmentPriceListPromotionHistory) AppliedRules() core.Interface {
	return core.Interface{m["applied_rules"]}
}
func (m ShipmentPriceListPromotionHistory) CreatedAt() core.Interface {
	return core.Interface{m["created_at"]}
}
func (m ShipmentPriceListPromotionHistory) UpdatedAt() core.Interface {
	return core.Interface{m["updated_at"]}
}
func (m ShipmentPriceListPromotionHistory) DeletedAt() core.Interface {
	return core.Interface{m["deleted_at"]}
}
func (m ShipmentPriceListPromotionHistory) WLPartnerID() core.Interface {
	return core.Interface{m["wl_partner_id"]}
}
func (m ShipmentPriceListPromotionHistory) ConnectionID() core.Interface {
	return core.Interface{m["connection_id"]}
}
func (m ShipmentPriceListPromotionHistory) PriorityPoint() core.Interface {
	return core.Interface{m["priority_point"]}
}

func (m *ShipmentPriceListPromotionHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 14)
	args := make([]interface{}, 14)
	for i := 0; i < 14; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShipmentPriceListPromotionHistory, 14)
	res["id"] = data[0]
	res["price_list_id"] = data[1]
	res["name"] = data[2]
	res["description"] = data[3]
	res["status"] = data[4]
	res["date_from"] = data[5]
	res["date_to"] = data[6]
	res["applied_rules"] = data[7]
	res["created_at"] = data[8]
	res["updated_at"] = data[9]
	res["deleted_at"] = data[10]
	res["wl_partner_id"] = data[11]
	res["connection_id"] = data[12]
	res["priority_point"] = data[13]
	*m = res
	return nil
}

func (ms *ShipmentPriceListPromotionHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 14)
	args := make([]interface{}, 14)
	for i := 0; i < 14; i++ {
		args[i] = &data[i]
	}
	res := make(ShipmentPriceListPromotionHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShipmentPriceListPromotionHistory)
		m["id"] = data[0]
		m["price_list_id"] = data[1]
		m["name"] = data[2]
		m["description"] = data[3]
		m["status"] = data[4]
		m["date_from"] = data[5]
		m["date_to"] = data[6]
		m["applied_rules"] = data[7]
		m["created_at"] = data[8]
		m["updated_at"] = data[9]
		m["deleted_at"] = data[10]
		m["wl_partner_id"] = data[11]
		m["connection_id"] = data[12]
		m["priority_point"] = data[13]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
