// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"time"

	core "etop.vn/backend/pkg/common/sq/core"
)

type SQLWriter = core.SQLWriter

// Type ShipnowFulfillment represents table shipnow_fulfillment
func sqlgenShipnowFulfillment(_ *ShipnowFulfillment) bool { return true }

type ShipnowFulfillments []*ShipnowFulfillment

const __sqlShipnowFulfillment_Table = "shipnow_fulfillment"
const __sqlShipnowFulfillment_ListCols = "\"id\",\"shop_id\",\"partner_id\",\"order_ids\",\"pickup_address\",\"carrier\",\"shipping_service_code\",\"shipping_service_fee\",\"shipping_service_name\",\"shipping_service_description\",\"chargeable_weight\",\"gross_weight\",\"basket_value\",\"cod_amount\",\"shipping_note\",\"request_pickup_at\",\"delivery_points\",\"cancel_reason\",\"status\",\"confirm_status\",\"shipping_status\",\"etop_payment_status\",\"shipping_state\",\"shipping_code\",\"fee_lines\",\"carrier_fee_lines\",\"total_fee\",\"shipping_created_at\",\"shipping_picking_at\",\"shipping_delivering_at\",\"shipping_delivered_at\",\"shipping_cancelled_at\",\"sync_status\",\"sync_states\",\"created_at\",\"updated_at\",\"cod_etop_transfered_at\",\"shipping_shared_link\""
const __sqlShipnowFulfillment_Insert = "INSERT INTO \"shipnow_fulfillment\" (" + __sqlShipnowFulfillment_ListCols + ") VALUES"
const __sqlShipnowFulfillment_Select = "SELECT " + __sqlShipnowFulfillment_ListCols + " FROM \"shipnow_fulfillment\""
const __sqlShipnowFulfillment_Select_history = "SELECT " + __sqlShipnowFulfillment_ListCols + " FROM history.\"shipnow_fulfillment\""
const __sqlShipnowFulfillment_UpdateAll = "UPDATE \"shipnow_fulfillment\" SET (" + __sqlShipnowFulfillment_ListCols + ")"

func (m *ShipnowFulfillment) SQLTableName() string  { return "shipnow_fulfillment" }
func (m *ShipnowFulfillments) SQLTableName() string { return "shipnow_fulfillment" }
func (m *ShipnowFulfillment) SQLListCols() string   { return __sqlShipnowFulfillment_ListCols }

func (m *ShipnowFulfillment) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.Int64(m.ShopID),
		core.Int64(m.PartnerID),
		core.Array{m.OrderIDs, opts},
		core.JSON{m.PickupAddress},
		core.String(m.Carrier),
		core.String(m.ShippingServiceCode),
		core.Int32(m.ShippingServiceFee),
		core.String(m.ShippingServiceName),
		core.String(m.ShippingServiceDescription),
		core.Int32(m.ChargeableWeight),
		core.Int32(m.GrossWeight),
		core.Int32(m.BasketValue),
		core.Int32(m.CODAmount),
		core.String(m.ShippingNote),
		core.Time(m.RequestPickupAt),
		core.JSON{m.DeliveryPoints},
		core.String(m.CancelReason),
		core.Int(m.Status),
		core.Int(m.ConfirmStatus),
		core.Int(m.ShippingStatus),
		core.Int(m.EtopPaymentStatus),
		core.String(m.ShippingState),
		core.String(m.ShippingCode),
		core.JSON{m.FeeLines},
		core.JSON{m.CarrierFeeLines},
		core.Int(m.TotalFee),
		core.Time(m.ShippingCreatedAt),
		core.Time(m.ShippingPickingAt),
		core.Time(m.ShippingDeliveringAt),
		core.Time(m.ShippingDeliveredAt),
		core.Time(m.ShippingCancelledAt),
		core.Int(m.SyncStatus),
		core.JSON{m.SyncStates},
		core.Now(m.CreatedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.CODEtopTransferedAt),
		core.String(m.ShippingSharedLink),
	}
}

func (m *ShipnowFulfillment) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.Int64)(&m.ShopID),
		(*core.Int64)(&m.PartnerID),
		core.Array{&m.OrderIDs, opts},
		core.JSON{&m.PickupAddress},
		(*core.String)(&m.Carrier),
		(*core.String)(&m.ShippingServiceCode),
		(*core.Int32)(&m.ShippingServiceFee),
		(*core.String)(&m.ShippingServiceName),
		(*core.String)(&m.ShippingServiceDescription),
		(*core.Int32)(&m.ChargeableWeight),
		(*core.Int32)(&m.GrossWeight),
		(*core.Int32)(&m.BasketValue),
		(*core.Int32)(&m.CODAmount),
		(*core.String)(&m.ShippingNote),
		(*core.Time)(&m.RequestPickupAt),
		core.JSON{&m.DeliveryPoints},
		(*core.String)(&m.CancelReason),
		(*core.Int)(&m.Status),
		(*core.Int)(&m.ConfirmStatus),
		(*core.Int)(&m.ShippingStatus),
		(*core.Int)(&m.EtopPaymentStatus),
		(*core.String)(&m.ShippingState),
		(*core.String)(&m.ShippingCode),
		core.JSON{&m.FeeLines},
		core.JSON{&m.CarrierFeeLines},
		(*core.Int)(&m.TotalFee),
		(*core.Time)(&m.ShippingCreatedAt),
		(*core.Time)(&m.ShippingPickingAt),
		(*core.Time)(&m.ShippingDeliveringAt),
		(*core.Time)(&m.ShippingDeliveredAt),
		(*core.Time)(&m.ShippingCancelledAt),
		(*core.Int)(&m.SyncStatus),
		core.JSON{&m.SyncStates},
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.CODEtopTransferedAt),
		(*core.String)(&m.ShippingSharedLink),
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
	w.WriteMarkers(38)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms ShipnowFulfillments) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(38)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
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
	if m.Carrier != "" {
		flag = true
		w.WriteName("carrier")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(string(m.Carrier))
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
		w.WriteArg(int(m.Status))
	}
	if m.ConfirmStatus != 0 {
		flag = true
		w.WriteName("confirm_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.ConfirmStatus))
	}
	if m.ShippingStatus != 0 {
		flag = true
		w.WriteName("shipping_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.ShippingStatus))
	}
	if m.EtopPaymentStatus != 0 {
		flag = true
		w.WriteName("etop_payment_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.EtopPaymentStatus))
	}
	if m.ShippingState != "" {
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
		w.WriteArg(int(m.SyncStatus))
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
	if !m.UpdatedAt.IsZero() {
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
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *ShipnowFulfillment) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlShipnowFulfillment_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(38)
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

func (m *ShipnowFulfillmentHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 38)
	args := make([]interface{}, 38)
	for i := 0; i < 38; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(ShipnowFulfillmentHistory, 38)
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
	*m = res
	return nil
}

func (ms *ShipnowFulfillmentHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 38)
	args := make([]interface{}, 38)
	for i := 0; i < 38; i++ {
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
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
