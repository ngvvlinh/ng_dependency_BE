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

type ShipnowFulfillments []*ShipnowFulfillment

const __sqlShipnowFulfillment_Table = "shipnow_fulfillment"
const __sqlShipnowFulfillment_ListCols = "\"id\",\"shop_id\",\"partner_id\",\"order_ids\",\"pickup_address\",\"carrier\",\"shipping_service_code\",\"shipping_service_fee\",\"shipping_service_name\",\"shipping_service_description\",\"chargeable_weight\",\"gross_weight\",\"basket_value\",\"cod_amount\",\"shipping_note\",\"request_pickup_at\",\"delivery_points\",\"cancel_reason\",\"status\",\"confirm_status\",\"shipping_status\",\"etop_payment_status\",\"shipping_state\",\"shipping_code\",\"fee_lines\",\"carrier_fee_lines\",\"total_fee\",\"shipping_created_at\",\"shipping_picking_at\",\"shipping_delivering_at\",\"shipping_delivered_at\",\"shipping_cancelled_at\",\"sync_status\",\"sync_states\",\"created_at\",\"updated_at\",\"cod_etop_transfered_at\",\"shipping_shared_link\",\"address_to_province_code\",\"address_to_district_code\",\"address_to_phone\",\"address_to_full_name_norm\",\"connection_id\",\"connection_method\",\"external_id\",\"coupon\",\"driver_phone\",\"driver_name\",\"rid\""
const __sqlShipnowFulfillment_ListColsOnConflict = "\"id\" = EXCLUDED.\"id\",\"shop_id\" = EXCLUDED.\"shop_id\",\"partner_id\" = EXCLUDED.\"partner_id\",\"order_ids\" = EXCLUDED.\"order_ids\",\"pickup_address\" = EXCLUDED.\"pickup_address\",\"carrier\" = EXCLUDED.\"carrier\",\"shipping_service_code\" = EXCLUDED.\"shipping_service_code\",\"shipping_service_fee\" = EXCLUDED.\"shipping_service_fee\",\"shipping_service_name\" = EXCLUDED.\"shipping_service_name\",\"shipping_service_description\" = EXCLUDED.\"shipping_service_description\",\"chargeable_weight\" = EXCLUDED.\"chargeable_weight\",\"gross_weight\" = EXCLUDED.\"gross_weight\",\"basket_value\" = EXCLUDED.\"basket_value\",\"cod_amount\" = EXCLUDED.\"cod_amount\",\"shipping_note\" = EXCLUDED.\"shipping_note\",\"request_pickup_at\" = EXCLUDED.\"request_pickup_at\",\"delivery_points\" = EXCLUDED.\"delivery_points\",\"cancel_reason\" = EXCLUDED.\"cancel_reason\",\"status\" = EXCLUDED.\"status\",\"confirm_status\" = EXCLUDED.\"confirm_status\",\"shipping_status\" = EXCLUDED.\"shipping_status\",\"etop_payment_status\" = EXCLUDED.\"etop_payment_status\",\"shipping_state\" = EXCLUDED.\"shipping_state\",\"shipping_code\" = EXCLUDED.\"shipping_code\",\"fee_lines\" = EXCLUDED.\"fee_lines\",\"carrier_fee_lines\" = EXCLUDED.\"carrier_fee_lines\",\"total_fee\" = EXCLUDED.\"total_fee\",\"shipping_created_at\" = EXCLUDED.\"shipping_created_at\",\"shipping_picking_at\" = EXCLUDED.\"shipping_picking_at\",\"shipping_delivering_at\" = EXCLUDED.\"shipping_delivering_at\",\"shipping_delivered_at\" = EXCLUDED.\"shipping_delivered_at\",\"shipping_cancelled_at\" = EXCLUDED.\"shipping_cancelled_at\",\"sync_status\" = EXCLUDED.\"sync_status\",\"sync_states\" = EXCLUDED.\"sync_states\",\"created_at\" = EXCLUDED.\"created_at\",\"updated_at\" = EXCLUDED.\"updated_at\",\"cod_etop_transfered_at\" = EXCLUDED.\"cod_etop_transfered_at\",\"shipping_shared_link\" = EXCLUDED.\"shipping_shared_link\",\"address_to_province_code\" = EXCLUDED.\"address_to_province_code\",\"address_to_district_code\" = EXCLUDED.\"address_to_district_code\",\"address_to_phone\" = EXCLUDED.\"address_to_phone\",\"address_to_full_name_norm\" = EXCLUDED.\"address_to_full_name_norm\",\"connection_id\" = EXCLUDED.\"connection_id\",\"connection_method\" = EXCLUDED.\"connection_method\",\"external_id\" = EXCLUDED.\"external_id\",\"coupon\" = EXCLUDED.\"coupon\",\"driver_phone\" = EXCLUDED.\"driver_phone\",\"driver_name\" = EXCLUDED.\"driver_name\",\"rid\" = EXCLUDED.\"rid\""
const __sqlShipnowFulfillment_Insert = "INSERT INTO \"shipnow_fulfillment\" (" + __sqlShipnowFulfillment_ListCols + ") VALUES"
const __sqlShipnowFulfillment_Select = "SELECT " + __sqlShipnowFulfillment_ListCols + " FROM \"shipnow_fulfillment\""
const __sqlShipnowFulfillment_Select_history = "SELECT " + __sqlShipnowFulfillment_ListCols + " FROM history.\"shipnow_fulfillment\""
const __sqlShipnowFulfillment_UpdateAll = "UPDATE \"shipnow_fulfillment\" SET (" + __sqlShipnowFulfillment_ListCols + ")"
const __sqlShipnowFulfillment_UpdateOnConflict = " ON CONFLICT ON CONSTRAINT shipnow_fulfillment_pkey DO UPDATE SET"

func (m *ShipnowFulfillment) SQLTableName() string  { return "shipnow_fulfillment" }
func (m *ShipnowFulfillments) SQLTableName() string { return "shipnow_fulfillment" }
func (m *ShipnowFulfillment) SQLListCols() string   { return __sqlShipnowFulfillment_ListCols }

func (m *ShipnowFulfillment) SQLVerifySchema(db *cmsql.Database) {
	query := "SELECT " + __sqlShipnowFulfillment_ListCols + " FROM \"shipnow_fulfillment\" WHERE false"
	if _, err := db.SQL(query).Exec(); err != nil {
		db.RecordError(err)
	}
}

func (m *ShipnowFulfillment) Migration(db *cmsql.Database) {
	var mDBColumnNameAndType map[string]string
	if val, err := migration.GetColumnNamesAndTypes(db, "shipnow_fulfillment"); err != nil {
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
		"partner_id": {
			ColumnName:       "partner_id",
			ColumnType:       "dot.ID",
			ColumnDBType:     "int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"order_ids": {
			ColumnName:       "order_ids",
			ColumnType:       "[]dot.ID",
			ColumnDBType:     "[]int64",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"pickup_address": {
			ColumnName:       "pickup_address",
			ColumnType:       "*orderingmodel.OrderAddress",
			ColumnDBType:     "*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"carrier": {
			ColumnName:       "carrier",
			ColumnType:       "carriertypes.ShipnowCarrier",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"default", "ahamove"},
		},
		"shipping_service_code": {
			ColumnName:       "shipping_service_code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_service_fee": {
			ColumnName:       "shipping_service_fee",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_service_name": {
			ColumnName:       "shipping_service_name",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_service_description": {
			ColumnName:       "shipping_service_description",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"chargeable_weight": {
			ColumnName:       "chargeable_weight",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"gross_weight": {
			ColumnName:       "gross_weight",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"basket_value": {
			ColumnName:       "basket_value",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"cod_amount": {
			ColumnName:       "cod_amount",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_note": {
			ColumnName:       "shipping_note",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"request_pickup_at": {
			ColumnName:       "request_pickup_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"delivery_points": {
			ColumnName:       "delivery_points",
			ColumnType:       "[]*DeliveryPoint",
			ColumnDBType:     "[]*struct",
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
		"status": {
			ColumnName:       "status",
			ColumnType:       "status5.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "S", "N", "NS"},
		},
		"confirm_status": {
			ColumnName:       "confirm_status",
			ColumnType:       "status3.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "N"},
		},
		"shipping_status": {
			ColumnName:       "shipping_status",
			ColumnType:       "status5.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "S", "N", "NS"},
		},
		"etop_payment_status": {
			ColumnName:       "etop_payment_status",
			ColumnType:       "status4.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "S", "N"},
		},
		"shipping_state": {
			ColumnName:       "shipping_state",
			ColumnType:       "shipnow_state.State",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"default", "created", "assigning", "picking", "delivering", "delivered", "returning", "returned", "unknown", "undeliverable", "cancelled"},
		},
		"shipping_code": {
			ColumnName:       "shipping_code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"fee_lines": {
			ColumnName:       "fee_lines",
			ColumnType:       "[]*sharemodel.ShippingFeeLine",
			ColumnDBType:     "[]*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"carrier_fee_lines": {
			ColumnName:       "carrier_fee_lines",
			ColumnType:       "[]*sharemodel.ShippingFeeLine",
			ColumnDBType:     "[]*struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"total_fee": {
			ColumnName:       "total_fee",
			ColumnType:       "int",
			ColumnDBType:     "int",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_created_at": {
			ColumnName:       "shipping_created_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_picking_at": {
			ColumnName:       "shipping_picking_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_delivering_at": {
			ColumnName:       "shipping_delivering_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_delivered_at": {
			ColumnName:       "shipping_delivered_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_cancelled_at": {
			ColumnName:       "shipping_cancelled_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"sync_status": {
			ColumnName:       "sync_status",
			ColumnType:       "status4.Status",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"Z", "P", "S", "N"},
		},
		"sync_states": {
			ColumnName:       "sync_states",
			ColumnType:       "*sharemodel.FulfillmentSyncStates",
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
		"cod_etop_transfered_at": {
			ColumnName:       "cod_etop_transfered_at",
			ColumnType:       "time.Time",
			ColumnDBType:     "struct",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"shipping_shared_link": {
			ColumnName:       "shipping_shared_link",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"address_to_province_code": {
			ColumnName:       "address_to_province_code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"address_to_district_code": {
			ColumnName:       "address_to_district_code",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"address_to_phone": {
			ColumnName:       "address_to_phone",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"address_to_full_name_norm": {
			ColumnName:       "address_to_full_name_norm",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
		"connection_method": {
			ColumnName:       "connection_method",
			ColumnType:       "connection_type.ConnectionMethod",
			ColumnDBType:     "enum",
			ColumnTag:        "",
			ColumnEnumValues: []string{"unknown", "builtin", "topship", "direct"},
		},
		"external_id": {
			ColumnName:       "external_id",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"coupon": {
			ColumnName:       "coupon",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"driver_phone": {
			ColumnName:       "driver_phone",
			ColumnType:       "string",
			ColumnDBType:     "string",
			ColumnTag:        "",
			ColumnEnumValues: []string{},
		},
		"driver_name": {
			ColumnName:       "driver_name",
			ColumnType:       "string",
			ColumnDBType:     "string",
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
	if err := migration.Compare(db, "shipnow_fulfillment", mModelColumnNameAndType, mDBColumnNameAndType); err != nil {
		db.RecordError(err)
	}
}

func init() {
	__sqlModels = append(__sqlModels, (*ShipnowFulfillment)(nil))
}

func (m *ShipnowFulfillment) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		m.ID,
		m.ShopID,
		m.PartnerID,
		core.Array{m.OrderIDs, opts},
		core.JSON{m.PickupAddress},
		m.Carrier,
		core.String(m.ShippingServiceCode),
		core.Int(m.ShippingServiceFee),
		core.String(m.ShippingServiceName),
		core.String(m.ShippingServiceDescription),
		core.Int(m.ChargeableWeight),
		core.Int(m.GrossWeight),
		core.Int(m.BasketValue),
		core.Int(m.CODAmount),
		core.String(m.ShippingNote),
		core.Time(m.RequestPickupAt),
		core.JSON{m.DeliveryPoints},
		core.String(m.CancelReason),
		m.Status,
		m.ConfirmStatus,
		m.ShippingStatus,
		m.EtopPaymentStatus,
		m.ShippingState,
		core.String(m.ShippingCode),
		core.JSON{m.FeeLines},
		core.JSON{m.CarrierFeeLines},
		core.Int(m.TotalFee),
		core.Time(m.ShippingCreatedAt),
		core.Time(m.ShippingPickingAt),
		core.Time(m.ShippingDeliveringAt),
		core.Time(m.ShippingDeliveredAt),
		core.Time(m.ShippingCancelledAt),
		m.SyncStatus,
		core.JSON{m.SyncStates},
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.CODEtopTransferedAt),
		core.String(m.ShippingSharedLink),
		core.String(m.AddressToProvinceCode),
		core.String(m.AddressToDistrictCode),
		core.String(m.AddressToPhone),
		core.String(m.AddressToFullNameNorm),
		m.ConnectionID,
		m.ConnectionMethod,
		core.String(m.ExternalID),
		core.String(m.Coupon),
		core.String(m.DriverPhone),
		core.String(m.DriverName),
		m.Rid,
	}
}

func (m *ShipnowFulfillment) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		&m.ID,
		&m.ShopID,
		&m.PartnerID,
		core.Array{&m.OrderIDs, opts},
		core.JSON{&m.PickupAddress},
		&m.Carrier,
		(*core.String)(&m.ShippingServiceCode),
		(*core.Int)(&m.ShippingServiceFee),
		(*core.String)(&m.ShippingServiceName),
		(*core.String)(&m.ShippingServiceDescription),
		(*core.Int)(&m.ChargeableWeight),
		(*core.Int)(&m.GrossWeight),
		(*core.Int)(&m.BasketValue),
		(*core.Int)(&m.CODAmount),
		(*core.String)(&m.ShippingNote),
		(*core.Time)(&m.RequestPickupAt),
		core.JSON{&m.DeliveryPoints},
		(*core.String)(&m.CancelReason),
		&m.Status,
		&m.ConfirmStatus,
		&m.ShippingStatus,
		&m.EtopPaymentStatus,
		&m.ShippingState,
		(*core.String)(&m.ShippingCode),
		core.JSON{&m.FeeLines},
		core.JSON{&m.CarrierFeeLines},
		(*core.Int)(&m.TotalFee),
		(*core.Time)(&m.ShippingCreatedAt),
		(*core.Time)(&m.ShippingPickingAt),
		(*core.Time)(&m.ShippingDeliveringAt),
		(*core.Time)(&m.ShippingDeliveredAt),
		(*core.Time)(&m.ShippingCancelledAt),
		&m.SyncStatus,
		core.JSON{&m.SyncStates},
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.CODEtopTransferedAt),
		(*core.String)(&m.ShippingSharedLink),
		(*core.String)(&m.AddressToProvinceCode),
		(*core.String)(&m.AddressToDistrictCode),
		(*core.String)(&m.AddressToPhone),
		(*core.String)(&m.AddressToFullNameNorm),
		&m.ConnectionID,
		&m.ConnectionMethod,
		(*core.String)(&m.ExternalID),
		(*core.String)(&m.Coupon),
		(*core.String)(&m.DriverPhone),
		(*core.String)(&m.DriverName),
		&m.Rid,
	}
}

func (m *ShipnowFulfillment) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *ShipnowFulfillments) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(ShipnowFulfillments, 0, 128)
	for rows.Next() {
		m := new(ShipnowFulfillment)
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

func (_ *ShipnowFulfillment) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Select)
	return nil
}

func (_ *ShipnowFulfillments) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Select)
	return nil
}

func (m *ShipnowFulfillment) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(49)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShipnowFulfillments) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(49)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *ShipnowFulfillment) SQLUpsert(w SQLWriter) error {
	m.SQLInsert(w)
	w.WriteQueryString(__sqlShipnowFulfillment_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipnowFulfillment_ListColsOnConflict)
	return nil
}

func (ms ShipnowFulfillments) SQLUpsert(w SQLWriter) error {
	ms.SQLInsert(w)
	w.WriteQueryString(__sqlShipnowFulfillment_UpdateOnConflict)
	w.WriteQueryString(" ")
	w.WriteQueryString(__sqlShipnowFulfillment_ListColsOnConflict)
	return nil
}

func (m *ShipnowFulfillment) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("shipnow_fulfillment")
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
	if m.PartnerID != 0 {
		flag = true
		w.WriteName("partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PartnerID)
	}
	if m.OrderIDs != nil {
		flag = true
		w.WriteName("order_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.OrderIDs, opts})
	}
	if m.PickupAddress != nil {
		flag = true
		w.WriteName("pickup_address")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.PickupAddress})
	}
	if m.Carrier != 0 {
		flag = true
		w.WriteName("carrier")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Carrier)
	}
	if m.ShippingServiceCode != "" {
		flag = true
		w.WriteName("shipping_service_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingServiceCode)
	}
	if m.ShippingServiceFee != 0 {
		flag = true
		w.WriteName("shipping_service_fee")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingServiceFee)
	}
	if m.ShippingServiceName != "" {
		flag = true
		w.WriteName("shipping_service_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingServiceName)
	}
	if m.ShippingServiceDescription != "" {
		flag = true
		w.WriteName("shipping_service_description")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingServiceDescription)
	}
	if m.ChargeableWeight != 0 {
		flag = true
		w.WriteName("chargeable_weight")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ChargeableWeight)
	}
	if m.GrossWeight != 0 {
		flag = true
		w.WriteName("gross_weight")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.GrossWeight)
	}
	if m.BasketValue != 0 {
		flag = true
		w.WriteName("basket_value")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.BasketValue)
	}
	if m.CODAmount != 0 {
		flag = true
		w.WriteName("cod_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CODAmount)
	}
	if m.ShippingNote != "" {
		flag = true
		w.WriteName("shipping_note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingNote)
	}
	if !m.RequestPickupAt.IsZero() {
		flag = true
		w.WriteName("request_pickup_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.RequestPickupAt)
	}
	if m.DeliveryPoints != nil {
		flag = true
		w.WriteName("delivery_points")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.DeliveryPoints})
	}
	if m.CancelReason != "" {
		flag = true
		w.WriteName("cancel_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelReason)
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Status)
	}
	if m.ConfirmStatus != 0 {
		flag = true
		w.WriteName("confirm_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConfirmStatus)
	}
	if m.ShippingStatus != 0 {
		flag = true
		w.WriteName("shipping_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingStatus)
	}
	if m.EtopPaymentStatus != 0 {
		flag = true
		w.WriteName("etop_payment_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.EtopPaymentStatus)
	}
	if m.ShippingState != 0 {
		flag = true
		w.WriteName("shipping_state")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingState)
	}
	if m.ShippingCode != "" {
		flag = true
		w.WriteName("shipping_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingCode)
	}
	if m.FeeLines != nil {
		flag = true
		w.WriteName("fee_lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.FeeLines})
	}
	if m.CarrierFeeLines != nil {
		flag = true
		w.WriteName("carrier_fee_lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.CarrierFeeLines})
	}
	if m.TotalFee != 0 {
		flag = true
		w.WriteName("total_fee")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalFee)
	}
	if !m.ShippingCreatedAt.IsZero() {
		flag = true
		w.WriteName("shipping_created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingCreatedAt)
	}
	if !m.ShippingPickingAt.IsZero() {
		flag = true
		w.WriteName("shipping_picking_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingPickingAt)
	}
	if !m.ShippingDeliveringAt.IsZero() {
		flag = true
		w.WriteName("shipping_delivering_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingDeliveringAt)
	}
	if !m.ShippingDeliveredAt.IsZero() {
		flag = true
		w.WriteName("shipping_delivered_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingDeliveredAt)
	}
	if !m.ShippingCancelledAt.IsZero() {
		flag = true
		w.WriteName("shipping_cancelled_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingCancelledAt)
	}
	if m.SyncStatus != 0 {
		flag = true
		w.WriteName("sync_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.SyncStatus)
	}
	if m.SyncStates != nil {
		flag = true
		w.WriteName("sync_states")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.SyncStates})
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
	if !m.CODEtopTransferedAt.IsZero() {
		flag = true
		w.WriteName("cod_etop_transfered_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CODEtopTransferedAt)
	}
	if m.ShippingSharedLink != "" {
		flag = true
		w.WriteName("shipping_shared_link")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingSharedLink)
	}
	if m.AddressToProvinceCode != "" {
		flag = true
		w.WriteName("address_to_province_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AddressToProvinceCode)
	}
	if m.AddressToDistrictCode != "" {
		flag = true
		w.WriteName("address_to_district_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AddressToDistrictCode)
	}
	if m.AddressToPhone != "" {
		flag = true
		w.WriteName("address_to_phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AddressToPhone)
	}
	if m.AddressToFullNameNorm != "" {
		flag = true
		w.WriteName("address_to_full_name_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.AddressToFullNameNorm)
	}
	if m.ConnectionID != 0 {
		flag = true
		w.WriteName("connection_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionID)
	}
	if m.ConnectionMethod != 0 {
		flag = true
		w.WriteName("connection_method")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ConnectionMethod)
	}
	if m.ExternalID != "" {
		flag = true
		w.WriteName("external_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalID)
	}
	if m.Coupon != "" {
		flag = true
		w.WriteName("coupon")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Coupon)
	}
	if m.DriverPhone != "" {
		flag = true
		w.WriteName("driver_phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DriverPhone)
	}
	if m.DriverName != "" {
		flag = true
		w.WriteName("driver_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.DriverName)
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

func (m *ShipnowFulfillment) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(49)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type ShipnowFulfillmentHistory map[string]interface{}
type ShipnowFulfillmentHistories []map[string]interface{}

func (m *ShipnowFulfillmentHistory) SQLTableName() string  { return "history.\"shipnow_fulfillment\"" }
func (m ShipnowFulfillmentHistories) SQLTableName() string { return "history.\"shipnow_fulfillment\"" }

func (m *ShipnowFulfillmentHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Select_history)
	return nil
}

func (m ShipnowFulfillmentHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Select_history)
	return nil
}

func (m ShipnowFulfillmentHistory) ID() core.Interface        { return core.Interface{m["id"]} }
func (m ShipnowFulfillmentHistory) ShopID() core.Interface    { return core.Interface{m["shop_id"]} }
func (m ShipnowFulfillmentHistory) PartnerID() core.Interface { return core.Interface{m["partner_id"]} }
func (m ShipnowFulfillmentHistory) OrderIDs() core.Interface  { return core.Interface{m["order_ids"]} }
func (m ShipnowFulfillmentHistory) PickupAddress() core.Interface {
	return core.Interface{m["pickup_address"]}
}
func (m ShipnowFulfillmentHistory) Carrier() core.Interface { return core.Interface{m["carrier"]} }
func (m ShipnowFulfillmentHistory) ShippingServiceCode() core.Interface {
	return core.Interface{m["shipping_service_code"]}
}
func (m ShipnowFulfillmentHistory) ShippingServiceFee() core.Interface {
	return core.Interface{m["shipping_service_fee"]}
}
func (m ShipnowFulfillmentHistory) ShippingServiceName() core.Interface {
	return core.Interface{m["shipping_service_name"]}
}
func (m ShipnowFulfillmentHistory) ShippingServiceDescription() core.Interface {
	return core.Interface{m["shipping_service_description"]}
}
func (m ShipnowFulfillmentHistory) ChargeableWeight() core.Interface {
	return core.Interface{m["chargeable_weight"]}
}
func (m ShipnowFulfillmentHistory) GrossWeight() core.Interface {
	return core.Interface{m["gross_weight"]}
}
func (m ShipnowFulfillmentHistory) BasketValue() core.Interface {
	return core.Interface{m["basket_value"]}
}
func (m ShipnowFulfillmentHistory) CODAmount() core.Interface { return core.Interface{m["cod_amount"]} }
func (m ShipnowFulfillmentHistory) ShippingNote() core.Interface {
	return core.Interface{m["shipping_note"]}
}
func (m ShipnowFulfillmentHistory) RequestPickupAt() core.Interface {
	return core.Interface{m["request_pickup_at"]}
}
func (m ShipnowFulfillmentHistory) DeliveryPoints() core.Interface {
	return core.Interface{m["delivery_points"]}
}
func (m ShipnowFulfillmentHistory) CancelReason() core.Interface {
	return core.Interface{m["cancel_reason"]}
}
func (m ShipnowFulfillmentHistory) Status() core.Interface { return core.Interface{m["status"]} }
func (m ShipnowFulfillmentHistory) ConfirmStatus() core.Interface {
	return core.Interface{m["confirm_status"]}
}
func (m ShipnowFulfillmentHistory) ShippingStatus() core.Interface {
	return core.Interface{m["shipping_status"]}
}
func (m ShipnowFulfillmentHistory) EtopPaymentStatus() core.Interface {
	return core.Interface{m["etop_payment_status"]}
}
func (m ShipnowFulfillmentHistory) ShippingState() core.Interface {
	return core.Interface{m["shipping_state"]}
}
func (m ShipnowFulfillmentHistory) ShippingCode() core.Interface {
	return core.Interface{m["shipping_code"]}
}
func (m ShipnowFulfillmentHistory) FeeLines() core.Interface { return core.Interface{m["fee_lines"]} }
func (m ShipnowFulfillmentHistory) CarrierFeeLines() core.Interface {
	return core.Interface{m["carrier_fee_lines"]}
}
func (m ShipnowFulfillmentHistory) TotalFee() core.Interface { return core.Interface{m["total_fee"]} }
func (m ShipnowFulfillmentHistory) ShippingCreatedAt() core.Interface {
	return core.Interface{m["shipping_created_at"]}
}
func (m ShipnowFulfillmentHistory) ShippingPickingAt() core.Interface {
	return core.Interface{m["shipping_picking_at"]}
}
func (m ShipnowFulfillmentHistory) ShippingDeliveringAt() core.Interface {
	return core.Interface{m["shipping_delivering_at"]}
}
func (m ShipnowFulfillmentHistory) ShippingDeliveredAt() core.Interface {
	return core.Interface{m["shipping_delivered_at"]}
}
func (m ShipnowFulfillmentHistory) ShippingCancelledAt() core.Interface {
	return core.Interface{m["shipping_cancelled_at"]}
}
func (m ShipnowFulfillmentHistory) SyncStatus() core.Interface {
	return core.Interface{m["sync_status"]}
}
func (m ShipnowFulfillmentHistory) SyncStates() core.Interface {
	return core.Interface{m["sync_states"]}
}
func (m ShipnowFulfillmentHistory) CreatedAt() core.Interface { return core.Interface{m["created_at"]} }
func (m ShipnowFulfillmentHistory) UpdatedAt() core.Interface { return core.Interface{m["updated_at"]} }
func (m ShipnowFulfillmentHistory) CODEtopTransferedAt() core.Interface {
	return core.Interface{m["cod_etop_transfered_at"]}
}
func (m ShipnowFulfillmentHistory) ShippingSharedLink() core.Interface {
	return core.Interface{m["shipping_shared_link"]}
}
func (m ShipnowFulfillmentHistory) AddressToProvinceCode() core.Interface {
	return core.Interface{m["address_to_province_code"]}
}
func (m ShipnowFulfillmentHistory) AddressToDistrictCode() core.Interface {
	return core.Interface{m["address_to_district_code"]}
}
func (m ShipnowFulfillmentHistory) AddressToPhone() core.Interface {
	return core.Interface{m["address_to_phone"]}
}
func (m ShipnowFulfillmentHistory) AddressToFullNameNorm() core.Interface {
	return core.Interface{m["address_to_full_name_norm"]}
}
func (m ShipnowFulfillmentHistory) ConnectionID() core.Interface {
	return core.Interface{m["connection_id"]}
}
func (m ShipnowFulfillmentHistory) ConnectionMethod() core.Interface {
	return core.Interface{m["connection_method"]}
}
func (m ShipnowFulfillmentHistory) ExternalID() core.Interface {
	return core.Interface{m["external_id"]}
}
func (m ShipnowFulfillmentHistory) Coupon() core.Interface { return core.Interface{m["coupon"]} }
func (m ShipnowFulfillmentHistory) DriverPhone() core.Interface {
	return core.Interface{m["driver_phone"]}
}
func (m ShipnowFulfillmentHistory) DriverName() core.Interface {
	return core.Interface{m["driver_name"]}
}
func (m ShipnowFulfillmentHistory) Rid() core.Interface { return core.Interface{m["rid"]} }

func (m *ShipnowFulfillmentHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 49)
	args := make([]interface{}, 49)
	for i := 0; i < 49; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShipnowFulfillmentHistory, 49)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["partner_id"] = data[2]
	res["order_ids"] = data[3]
	res["pickup_address"] = data[4]
	res["carrier"] = data[5]
	res["shipping_service_code"] = data[6]
	res["shipping_service_fee"] = data[7]
	res["shipping_service_name"] = data[8]
	res["shipping_service_description"] = data[9]
	res["chargeable_weight"] = data[10]
	res["gross_weight"] = data[11]
	res["basket_value"] = data[12]
	res["cod_amount"] = data[13]
	res["shipping_note"] = data[14]
	res["request_pickup_at"] = data[15]
	res["delivery_points"] = data[16]
	res["cancel_reason"] = data[17]
	res["status"] = data[18]
	res["confirm_status"] = data[19]
	res["shipping_status"] = data[20]
	res["etop_payment_status"] = data[21]
	res["shipping_state"] = data[22]
	res["shipping_code"] = data[23]
	res["fee_lines"] = data[24]
	res["carrier_fee_lines"] = data[25]
	res["total_fee"] = data[26]
	res["shipping_created_at"] = data[27]
	res["shipping_picking_at"] = data[28]
	res["shipping_delivering_at"] = data[29]
	res["shipping_delivered_at"] = data[30]
	res["shipping_cancelled_at"] = data[31]
	res["sync_status"] = data[32]
	res["sync_states"] = data[33]
	res["created_at"] = data[34]
	res["updated_at"] = data[35]
	res["cod_etop_transfered_at"] = data[36]
	res["shipping_shared_link"] = data[37]
	res["address_to_province_code"] = data[38]
	res["address_to_district_code"] = data[39]
	res["address_to_phone"] = data[40]
	res["address_to_full_name_norm"] = data[41]
	res["connection_id"] = data[42]
	res["connection_method"] = data[43]
	res["external_id"] = data[44]
	res["coupon"] = data[45]
	res["driver_phone"] = data[46]
	res["driver_name"] = data[47]
	res["rid"] = data[48]
	*m = res
	return nil
}

func (ms *ShipnowFulfillmentHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 49)
	args := make([]interface{}, 49)
	for i := 0; i < 49; i++ {
		args[i] = &data[i]
	}
	res := make(ShipnowFulfillmentHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(ShipnowFulfillmentHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["partner_id"] = data[2]
		m["order_ids"] = data[3]
		m["pickup_address"] = data[4]
		m["carrier"] = data[5]
		m["shipping_service_code"] = data[6]
		m["shipping_service_fee"] = data[7]
		m["shipping_service_name"] = data[8]
		m["shipping_service_description"] = data[9]
		m["chargeable_weight"] = data[10]
		m["gross_weight"] = data[11]
		m["basket_value"] = data[12]
		m["cod_amount"] = data[13]
		m["shipping_note"] = data[14]
		m["request_pickup_at"] = data[15]
		m["delivery_points"] = data[16]
		m["cancel_reason"] = data[17]
		m["status"] = data[18]
		m["confirm_status"] = data[19]
		m["shipping_status"] = data[20]
		m["etop_payment_status"] = data[21]
		m["shipping_state"] = data[22]
		m["shipping_code"] = data[23]
		m["fee_lines"] = data[24]
		m["carrier_fee_lines"] = data[25]
		m["total_fee"] = data[26]
		m["shipping_created_at"] = data[27]
		m["shipping_picking_at"] = data[28]
		m["shipping_delivering_at"] = data[29]
		m["shipping_delivered_at"] = data[30]
		m["shipping_cancelled_at"] = data[31]
		m["sync_status"] = data[32]
		m["sync_states"] = data[33]
		m["created_at"] = data[34]
		m["updated_at"] = data[35]
		m["cod_etop_transfered_at"] = data[36]
		m["shipping_shared_link"] = data[37]
		m["address_to_province_code"] = data[38]
		m["address_to_district_code"] = data[39]
		m["address_to_phone"] = data[40]
		m["address_to_full_name_norm"] = data[41]
		m["connection_id"] = data[42]
		m["connection_method"] = data[43]
		m["external_id"] = data[44]
		m["coupon"] = data[45]
		m["driver_phone"] = data[46]
		m["driver_name"] = data[47]
		m["rid"] = data[48]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
