// Code generated by goderive DO NOT EDIT.

package model

import (
	"database/sql"
	"time"

	core "etop.vn/backend/pkg/common/sq/core"
)

type SQLWriter = core.SQLWriter

// Type Order represents table order
func sqlgenOrder(_ *Order) bool { return true }

type Orders []*Order

const __sqlOrder_Table = "order"
const __sqlOrder_ListCols = "\"id\",\"shop_id\",\"code\",\"ed_code\",\"product_ids\",\"variant_ids\",\"partner_id\",\"currency\",\"payment_method\",\"customer\",\"customer_address\",\"billing_address\",\"shipping_address\",\"customer_name\",\"customer_phone\",\"customer_email\",\"created_at\",\"processed_at\",\"updated_at\",\"closed_at\",\"confirmed_at\",\"cancelled_at\",\"cancel_reason\",\"customer_confirm\",\"shop_confirm\",\"confirm_status\",\"fulfillment_shipping_status\",\"customer_payment_status\",\"etop_payment_status\",\"status\",\"fulfillment_shipping_states\",\"fulfillment_payment_statuses\",\"lines\",\"discounts\",\"total_items\",\"basket_value\",\"total_weight\",\"total_tax\",\"order_discount\",\"total_discount\",\"shop_shipping_fee\",\"total_fee\",\"fee_lines\",\"shop_cod\",\"total_amount\",\"order_note\",\"shop_note\",\"shipping_note\",\"order_source_type\",\"order_source_id\",\"external_order_id\",\"reference_url\",\"external_url\",\"shop_shipping\",\"is_outside_etop\",\"ghn_note_code\",\"try_on\",\"customer_name_norm\",\"product_name_norm\",\"fulfillment_type\",\"fulfillment_ids\",\"external_meta\""
const __sqlOrder_Insert = "INSERT INTO \"order\" (" + __sqlOrder_ListCols + ") VALUES"
const __sqlOrder_Select = "SELECT " + __sqlOrder_ListCols + " FROM \"order\""
const __sqlOrder_Select_history = "SELECT " + __sqlOrder_ListCols + " FROM history.\"order\""
const __sqlOrder_UpdateAll = "UPDATE \"order\" SET (" + __sqlOrder_ListCols + ")"

func (m *Order) SQLTableName() string  { return "order" }
func (m *Orders) SQLTableName() string { return "order" }
func (m *Order) SQLListCols() string   { return __sqlOrder_ListCols }

func (m *Order) SQLArgs(opts core.Opts, create bool) []interface{} {
	now := time.Now()
	return []interface{}{
		core.Int64(m.ID),
		core.Int64(m.ShopID),
		core.String(m.Code),
		core.String(m.EdCode),
		core.Array{m.ProductIDs, opts},
		core.Array{m.VariantIDs, opts},
		core.Int64(m.PartnerID),
		core.String(m.Currency),
		core.String(m.PaymentMethod),
		core.JSON{m.Customer},
		core.JSON{m.CustomerAddress},
		core.JSON{m.BillingAddress},
		core.JSON{m.ShippingAddress},
		core.String(m.CustomerName),
		core.String(m.CustomerPhone),
		core.String(m.CustomerEmail),
		core.Time(m.CreatedAt),
		core.Now(m.ProcessedAt, now, create),
		core.Now(m.UpdatedAt, now, true),
		core.Time(m.ClosedAt),
		core.Time(m.ConfirmedAt),
		core.Time(m.CancelledAt),
		core.String(m.CancelReason),
		core.Int(m.CustomerConfirm),
		core.Int(m.ShopConfirm),
		core.Int(m.ConfirmStatus),
		core.Int(m.FulfillmentShippingStatus),
		core.Int(m.CustomerPaymentStatus),
		core.Int(m.EtopPaymentStatus),
		core.Int(m.Status),
		core.Array{m.FulfillmentShippingStates, opts},
		core.Array{m.FulfillmentPaymentStatuses, opts},
		core.JSON{m.Lines},
		core.JSON{m.Discounts},
		core.Int(m.TotalItems),
		core.Int(m.BasketValue),
		core.Int(m.TotalWeight),
		core.Int(m.TotalTax),
		core.Int(m.OrderDiscount),
		core.Int(m.TotalDiscount),
		core.Int(m.ShopShippingFee),
		core.Int(m.TotalFee),
		core.JSON{m.FeeLines},
		core.Int(m.ShopCOD),
		core.Int(m.TotalAmount),
		core.String(m.OrderNote),
		core.String(m.ShopNote),
		core.String(m.ShippingNote),
		core.String(m.OrderSourceType),
		core.Int64(m.OrderSourceID),
		core.String(m.ExternalOrderID),
		core.String(m.ReferenceURL),
		core.String(m.ExternalURL),
		core.JSON{m.ShopShipping},
		core.Bool(m.IsOutsideEtop),
		core.String(m.GhnNoteCode),
		core.String(m.TryOn),
		core.String(m.CustomerNameNorm),
		core.String(m.ProductNameNorm),
		core.Int32(m.FulfillmentType),
		core.Array{m.FulfillmentIDs, opts},
		core.JSON{m.ExternalMeta},
	}
}

func (m *Order) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.ID),
		(*core.Int64)(&m.ShopID),
		(*core.String)(&m.Code),
		(*core.String)(&m.EdCode),
		core.Array{&m.ProductIDs, opts},
		core.Array{&m.VariantIDs, opts},
		(*core.Int64)(&m.PartnerID),
		(*core.String)(&m.Currency),
		(*core.String)(&m.PaymentMethod),
		core.JSON{&m.Customer},
		core.JSON{&m.CustomerAddress},
		core.JSON{&m.BillingAddress},
		core.JSON{&m.ShippingAddress},
		(*core.String)(&m.CustomerName),
		(*core.String)(&m.CustomerPhone),
		(*core.String)(&m.CustomerEmail),
		(*core.Time)(&m.CreatedAt),
		(*core.Time)(&m.ProcessedAt),
		(*core.Time)(&m.UpdatedAt),
		(*core.Time)(&m.ClosedAt),
		(*core.Time)(&m.ConfirmedAt),
		(*core.Time)(&m.CancelledAt),
		(*core.String)(&m.CancelReason),
		(*core.Int)(&m.CustomerConfirm),
		(*core.Int)(&m.ShopConfirm),
		(*core.Int)(&m.ConfirmStatus),
		(*core.Int)(&m.FulfillmentShippingStatus),
		(*core.Int)(&m.CustomerPaymentStatus),
		(*core.Int)(&m.EtopPaymentStatus),
		(*core.Int)(&m.Status),
		core.Array{&m.FulfillmentShippingStates, opts},
		core.Array{&m.FulfillmentPaymentStatuses, opts},
		core.JSON{&m.Lines},
		core.JSON{&m.Discounts},
		(*core.Int)(&m.TotalItems),
		(*core.Int)(&m.BasketValue),
		(*core.Int)(&m.TotalWeight),
		(*core.Int)(&m.TotalTax),
		(*core.Int)(&m.OrderDiscount),
		(*core.Int)(&m.TotalDiscount),
		(*core.Int)(&m.ShopShippingFee),
		(*core.Int)(&m.TotalFee),
		core.JSON{&m.FeeLines},
		(*core.Int)(&m.ShopCOD),
		(*core.Int)(&m.TotalAmount),
		(*core.String)(&m.OrderNote),
		(*core.String)(&m.ShopNote),
		(*core.String)(&m.ShippingNote),
		(*core.String)(&m.OrderSourceType),
		(*core.Int64)(&m.OrderSourceID),
		(*core.String)(&m.ExternalOrderID),
		(*core.String)(&m.ReferenceURL),
		(*core.String)(&m.ExternalURL),
		core.JSON{&m.ShopShipping},
		(*core.Bool)(&m.IsOutsideEtop),
		(*core.String)(&m.GhnNoteCode),
		(*core.String)(&m.TryOn),
		(*core.String)(&m.CustomerNameNorm),
		(*core.String)(&m.ProductNameNorm),
		(*core.Int32)(&m.FulfillmentType),
		core.Array{&m.FulfillmentIDs, opts},
		core.JSON{&m.ExternalMeta},
	}
}

func (m *Order) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *Orders) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(Orders, 0, 128)
	for rows.Next() {
		m := new(Order)
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

func (_ *Order) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_Select)
	return nil
}

func (_ *Orders) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_Select)
	return nil
}

func (m *Order) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(62)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms Orders) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(62)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *Order) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("order")
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
	if m.Code != "" {
		flag = true
		w.WriteName("code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Code)
	}
	if m.EdCode != "" {
		flag = true
		w.WriteName("ed_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.EdCode)
	}
	if m.ProductIDs != nil {
		flag = true
		w.WriteName("product_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.ProductIDs, opts})
	}
	if m.VariantIDs != nil {
		flag = true
		w.WriteName("variant_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.VariantIDs, opts})
	}
	if m.PartnerID != 0 {
		flag = true
		w.WriteName("partner_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PartnerID)
	}
	if m.Currency != "" {
		flag = true
		w.WriteName("currency")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Currency)
	}
	if m.PaymentMethod != "" {
		flag = true
		w.WriteName("payment_method")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaymentMethod)
	}
	if m.Customer != nil {
		flag = true
		w.WriteName("customer")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Customer})
	}
	if m.CustomerAddress != nil {
		flag = true
		w.WriteName("customer_address")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.CustomerAddress})
	}
	if m.BillingAddress != nil {
		flag = true
		w.WriteName("billing_address")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.BillingAddress})
	}
	if m.ShippingAddress != nil {
		flag = true
		w.WriteName("shipping_address")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.ShippingAddress})
	}
	if m.CustomerName != "" {
		flag = true
		w.WriteName("customer_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CustomerName)
	}
	if m.CustomerPhone != "" {
		flag = true
		w.WriteName("customer_phone")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CustomerPhone)
	}
	if m.CustomerEmail != "" {
		flag = true
		w.WriteName("customer_email")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CustomerEmail)
	}
	if !m.CreatedAt.IsZero() {
		flag = true
		w.WriteName("created_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CreatedAt)
	}
	if !m.ProcessedAt.IsZero() {
		flag = true
		w.WriteName("processed_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ProcessedAt)
	}
	if !m.UpdatedAt.IsZero() {
		flag = true
		w.WriteName("updated_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Now(m.UpdatedAt, time.Now(), true))
	}
	if !m.ClosedAt.IsZero() {
		flag = true
		w.WriteName("closed_at")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ClosedAt)
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
	if m.CancelReason != "" {
		flag = true
		w.WriteName("cancel_reason")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CancelReason)
	}
	if m.CustomerConfirm != 0 {
		flag = true
		w.WriteName("customer_confirm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.CustomerConfirm))
	}
	if m.ShopConfirm != 0 {
		flag = true
		w.WriteName("shop_confirm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.ShopConfirm))
	}
	if m.ConfirmStatus != 0 {
		flag = true
		w.WriteName("confirm_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.ConfirmStatus))
	}
	if m.FulfillmentShippingStatus != 0 {
		flag = true
		w.WriteName("fulfillment_shipping_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.FulfillmentShippingStatus))
	}
	if m.CustomerPaymentStatus != 0 {
		flag = true
		w.WriteName("customer_payment_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.CustomerPaymentStatus))
	}
	if m.EtopPaymentStatus != 0 {
		flag = true
		w.WriteName("etop_payment_status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.EtopPaymentStatus))
	}
	if m.Status != 0 {
		flag = true
		w.WriteName("status")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int(m.Status))
	}
	if m.FulfillmentShippingStates != nil {
		flag = true
		w.WriteName("fulfillment_shipping_states")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.FulfillmentShippingStates, opts})
	}
	if m.FulfillmentPaymentStatuses != nil {
		flag = true
		w.WriteName("fulfillment_payment_statuses")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.FulfillmentPaymentStatuses, opts})
	}
	if m.Lines != nil {
		flag = true
		w.WriteName("lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Lines})
	}
	if m.Discounts != nil {
		flag = true
		w.WriteName("discounts")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.Discounts})
	}
	if m.TotalItems != 0 {
		flag = true
		w.WriteName("total_items")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalItems)
	}
	if m.BasketValue != 0 {
		flag = true
		w.WriteName("basket_value")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.BasketValue)
	}
	if m.TotalWeight != 0 {
		flag = true
		w.WriteName("total_weight")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalWeight)
	}
	if m.TotalTax != 0 {
		flag = true
		w.WriteName("total_tax")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalTax)
	}
	if m.OrderDiscount != 0 {
		flag = true
		w.WriteName("order_discount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OrderDiscount)
	}
	if m.TotalDiscount != 0 {
		flag = true
		w.WriteName("total_discount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalDiscount)
	}
	if m.ShopShippingFee != 0 {
		flag = true
		w.WriteName("shop_shipping_fee")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopShippingFee)
	}
	if m.TotalFee != 0 {
		flag = true
		w.WriteName("total_fee")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalFee)
	}
	if m.FeeLines != nil {
		flag = true
		w.WriteName("fee_lines")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.FeeLines})
	}
	if m.ShopCOD != 0 {
		flag = true
		w.WriteName("shop_cod")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopCOD)
	}
	if m.TotalAmount != 0 {
		flag = true
		w.WriteName("total_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalAmount)
	}
	if m.OrderNote != "" {
		flag = true
		w.WriteName("order_note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OrderNote)
	}
	if m.ShopNote != "" {
		flag = true
		w.WriteName("shop_note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopNote)
	}
	if m.ShippingNote != "" {
		flag = true
		w.WriteName("shipping_note")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShippingNote)
	}
	if m.OrderSourceType != "" {
		flag = true
		w.WriteName("order_source_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(string(m.OrderSourceType))
	}
	if m.OrderSourceID != 0 {
		flag = true
		w.WriteName("order_source_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OrderSourceID)
	}
	if m.ExternalOrderID != "" {
		flag = true
		w.WriteName("external_order_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalOrderID)
	}
	if m.ReferenceURL != "" {
		flag = true
		w.WriteName("reference_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ReferenceURL)
	}
	if m.ExternalURL != "" {
		flag = true
		w.WriteName("external_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ExternalURL)
	}
	if m.ShopShipping != nil {
		flag = true
		w.WriteName("shop_shipping")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.ShopShipping})
	}
	if m.IsOutsideEtop {
		flag = true
		w.WriteName("is_outside_etop")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IsOutsideEtop)
	}
	if m.GhnNoteCode != "" {
		flag = true
		w.WriteName("ghn_note_code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.GhnNoteCode)
	}
	if m.TryOn != "" {
		flag = true
		w.WriteName("try_on")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(string(m.TryOn))
	}
	if m.CustomerNameNorm != "" {
		flag = true
		w.WriteName("customer_name_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.CustomerNameNorm)
	}
	if m.ProductNameNorm != "" {
		flag = true
		w.WriteName("product_name_norm")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ProductNameNorm)
	}
	if m.FulfillmentType != 0 {
		flag = true
		w.WriteName("fulfillment_type")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(int32(m.FulfillmentType))
	}
	if m.FulfillmentIDs != nil {
		flag = true
		w.WriteName("fulfillment_ids")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.Array{m.FulfillmentIDs, opts})
	}
	if m.ExternalMeta != nil {
		flag = true
		w.WriteName("external_meta")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(core.JSON{m.ExternalMeta})
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *Order) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(62)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type OrderHistory map[string]interface{}
type OrderHistories []map[string]interface{}

func (m *OrderHistory) SQLTableName() string  { return "history.\"order\"" }
func (m OrderHistories) SQLTableName() string { return "history.\"order\"" }

func (m *OrderHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_Select_history)
	return nil
}

func (m OrderHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrder_Select_history)
	return nil
}

func (m OrderHistory) ID() core.Interface              { return core.Interface{m["id"]} }
func (m OrderHistory) ShopID() core.Interface          { return core.Interface{m["shop_id"]} }
func (m OrderHistory) Code() core.Interface            { return core.Interface{m["code"]} }
func (m OrderHistory) EdCode() core.Interface          { return core.Interface{m["ed_code"]} }
func (m OrderHistory) ProductIDs() core.Interface      { return core.Interface{m["product_ids"]} }
func (m OrderHistory) VariantIDs() core.Interface      { return core.Interface{m["variant_ids"]} }
func (m OrderHistory) PartnerID() core.Interface       { return core.Interface{m["partner_id"]} }
func (m OrderHistory) Currency() core.Interface        { return core.Interface{m["currency"]} }
func (m OrderHistory) PaymentMethod() core.Interface   { return core.Interface{m["payment_method"]} }
func (m OrderHistory) Customer() core.Interface        { return core.Interface{m["customer"]} }
func (m OrderHistory) CustomerAddress() core.Interface { return core.Interface{m["customer_address"]} }
func (m OrderHistory) BillingAddress() core.Interface  { return core.Interface{m["billing_address"]} }
func (m OrderHistory) ShippingAddress() core.Interface { return core.Interface{m["shipping_address"]} }
func (m OrderHistory) CustomerName() core.Interface    { return core.Interface{m["customer_name"]} }
func (m OrderHistory) CustomerPhone() core.Interface   { return core.Interface{m["customer_phone"]} }
func (m OrderHistory) CustomerEmail() core.Interface   { return core.Interface{m["customer_email"]} }
func (m OrderHistory) CreatedAt() core.Interface       { return core.Interface{m["created_at"]} }
func (m OrderHistory) ProcessedAt() core.Interface     { return core.Interface{m["processed_at"]} }
func (m OrderHistory) UpdatedAt() core.Interface       { return core.Interface{m["updated_at"]} }
func (m OrderHistory) ClosedAt() core.Interface        { return core.Interface{m["closed_at"]} }
func (m OrderHistory) ConfirmedAt() core.Interface     { return core.Interface{m["confirmed_at"]} }
func (m OrderHistory) CancelledAt() core.Interface     { return core.Interface{m["cancelled_at"]} }
func (m OrderHistory) CancelReason() core.Interface    { return core.Interface{m["cancel_reason"]} }
func (m OrderHistory) CustomerConfirm() core.Interface { return core.Interface{m["customer_confirm"]} }
func (m OrderHistory) ShopConfirm() core.Interface     { return core.Interface{m["shop_confirm"]} }
func (m OrderHistory) ConfirmStatus() core.Interface   { return core.Interface{m["confirm_status"]} }
func (m OrderHistory) FulfillmentShippingStatus() core.Interface {
	return core.Interface{m["fulfillment_shipping_status"]}
}
func (m OrderHistory) CustomerPaymentStatus() core.Interface {
	return core.Interface{m["customer_payment_status"]}
}
func (m OrderHistory) EtopPaymentStatus() core.Interface {
	return core.Interface{m["etop_payment_status"]}
}
func (m OrderHistory) Status() core.Interface { return core.Interface{m["status"]} }
func (m OrderHistory) FulfillmentShippingStates() core.Interface {
	return core.Interface{m["fulfillment_shipping_states"]}
}
func (m OrderHistory) FulfillmentPaymentStatuses() core.Interface {
	return core.Interface{m["fulfillment_payment_statuses"]}
}
func (m OrderHistory) Lines() core.Interface           { return core.Interface{m["lines"]} }
func (m OrderHistory) Discounts() core.Interface       { return core.Interface{m["discounts"]} }
func (m OrderHistory) TotalItems() core.Interface      { return core.Interface{m["total_items"]} }
func (m OrderHistory) BasketValue() core.Interface     { return core.Interface{m["basket_value"]} }
func (m OrderHistory) TotalWeight() core.Interface     { return core.Interface{m["total_weight"]} }
func (m OrderHistory) TotalTax() core.Interface        { return core.Interface{m["total_tax"]} }
func (m OrderHistory) OrderDiscount() core.Interface   { return core.Interface{m["order_discount"]} }
func (m OrderHistory) TotalDiscount() core.Interface   { return core.Interface{m["total_discount"]} }
func (m OrderHistory) ShopShippingFee() core.Interface { return core.Interface{m["shop_shipping_fee"]} }
func (m OrderHistory) TotalFee() core.Interface        { return core.Interface{m["total_fee"]} }
func (m OrderHistory) FeeLines() core.Interface        { return core.Interface{m["fee_lines"]} }
func (m OrderHistory) ShopCOD() core.Interface         { return core.Interface{m["shop_cod"]} }
func (m OrderHistory) TotalAmount() core.Interface     { return core.Interface{m["total_amount"]} }
func (m OrderHistory) OrderNote() core.Interface       { return core.Interface{m["order_note"]} }
func (m OrderHistory) ShopNote() core.Interface        { return core.Interface{m["shop_note"]} }
func (m OrderHistory) ShippingNote() core.Interface    { return core.Interface{m["shipping_note"]} }
func (m OrderHistory) OrderSourceType() core.Interface { return core.Interface{m["order_source_type"]} }
func (m OrderHistory) OrderSourceID() core.Interface   { return core.Interface{m["order_source_id"]} }
func (m OrderHistory) ExternalOrderID() core.Interface { return core.Interface{m["external_order_id"]} }
func (m OrderHistory) ReferenceURL() core.Interface    { return core.Interface{m["reference_url"]} }
func (m OrderHistory) ExternalURL() core.Interface     { return core.Interface{m["external_url"]} }
func (m OrderHistory) ShopShipping() core.Interface    { return core.Interface{m["shop_shipping"]} }
func (m OrderHistory) IsOutsideEtop() core.Interface   { return core.Interface{m["is_outside_etop"]} }
func (m OrderHistory) GhnNoteCode() core.Interface     { return core.Interface{m["ghn_note_code"]} }
func (m OrderHistory) TryOn() core.Interface           { return core.Interface{m["try_on"]} }
func (m OrderHistory) CustomerNameNorm() core.Interface {
	return core.Interface{m["customer_name_norm"]}
}
func (m OrderHistory) ProductNameNorm() core.Interface { return core.Interface{m["product_name_norm"]} }
func (m OrderHistory) FulfillmentType() core.Interface { return core.Interface{m["fulfillment_type"]} }
func (m OrderHistory) FulfillmentIDs() core.Interface  { return core.Interface{m["fulfillment_ids"]} }
func (m OrderHistory) ExternalMeta() core.Interface    { return core.Interface{m["external_meta"]} }

func (m *OrderHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 62)
	args := make([]interface{}, 62)
	for i := 0; i < 62; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(OrderHistory, 62)
	res["id"] = data[0]
	res["shop_id"] = data[1]
	res["code"] = data[2]
	res["ed_code"] = data[3]
	res["product_ids"] = data[4]
	res["variant_ids"] = data[5]
	res["partner_id"] = data[6]
	res["currency"] = data[7]
	res["payment_method"] = data[8]
	res["customer"] = data[9]
	res["customer_address"] = data[10]
	res["billing_address"] = data[11]
	res["shipping_address"] = data[12]
	res["customer_name"] = data[13]
	res["customer_phone"] = data[14]
	res["customer_email"] = data[15]
	res["created_at"] = data[16]
	res["processed_at"] = data[17]
	res["updated_at"] = data[18]
	res["closed_at"] = data[19]
	res["confirmed_at"] = data[20]
	res["cancelled_at"] = data[21]
	res["cancel_reason"] = data[22]
	res["customer_confirm"] = data[23]
	res["shop_confirm"] = data[24]
	res["confirm_status"] = data[25]
	res["fulfillment_shipping_status"] = data[26]
	res["customer_payment_status"] = data[27]
	res["etop_payment_status"] = data[28]
	res["status"] = data[29]
	res["fulfillment_shipping_states"] = data[30]
	res["fulfillment_payment_statuses"] = data[31]
	res["lines"] = data[32]
	res["discounts"] = data[33]
	res["total_items"] = data[34]
	res["basket_value"] = data[35]
	res["total_weight"] = data[36]
	res["total_tax"] = data[37]
	res["order_discount"] = data[38]
	res["total_discount"] = data[39]
	res["shop_shipping_fee"] = data[40]
	res["total_fee"] = data[41]
	res["fee_lines"] = data[42]
	res["shop_cod"] = data[43]
	res["total_amount"] = data[44]
	res["order_note"] = data[45]
	res["shop_note"] = data[46]
	res["shipping_note"] = data[47]
	res["order_source_type"] = data[48]
	res["order_source_id"] = data[49]
	res["external_order_id"] = data[50]
	res["reference_url"] = data[51]
	res["external_url"] = data[52]
	res["shop_shipping"] = data[53]
	res["is_outside_etop"] = data[54]
	res["ghn_note_code"] = data[55]
	res["try_on"] = data[56]
	res["customer_name_norm"] = data[57]
	res["product_name_norm"] = data[58]
	res["fulfillment_type"] = data[59]
	res["fulfillment_ids"] = data[60]
	res["external_meta"] = data[61]
	*m = res
	return nil
}

func (ms *OrderHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 62)
	args := make([]interface{}, 62)
	for i := 0; i < 62; i++ {
		args[i] = &data[i]
	}
	res := make(OrderHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(OrderHistory)
		m["id"] = data[0]
		m["shop_id"] = data[1]
		m["code"] = data[2]
		m["ed_code"] = data[3]
		m["product_ids"] = data[4]
		m["variant_ids"] = data[5]
		m["partner_id"] = data[6]
		m["currency"] = data[7]
		m["payment_method"] = data[8]
		m["customer"] = data[9]
		m["customer_address"] = data[10]
		m["billing_address"] = data[11]
		m["shipping_address"] = data[12]
		m["customer_name"] = data[13]
		m["customer_phone"] = data[14]
		m["customer_email"] = data[15]
		m["created_at"] = data[16]
		m["processed_at"] = data[17]
		m["updated_at"] = data[18]
		m["closed_at"] = data[19]
		m["confirmed_at"] = data[20]
		m["cancelled_at"] = data[21]
		m["cancel_reason"] = data[22]
		m["customer_confirm"] = data[23]
		m["shop_confirm"] = data[24]
		m["confirm_status"] = data[25]
		m["fulfillment_shipping_status"] = data[26]
		m["customer_payment_status"] = data[27]
		m["etop_payment_status"] = data[28]
		m["status"] = data[29]
		m["fulfillment_shipping_states"] = data[30]
		m["fulfillment_payment_statuses"] = data[31]
		m["lines"] = data[32]
		m["discounts"] = data[33]
		m["total_items"] = data[34]
		m["basket_value"] = data[35]
		m["total_weight"] = data[36]
		m["total_tax"] = data[37]
		m["order_discount"] = data[38]
		m["total_discount"] = data[39]
		m["shop_shipping_fee"] = data[40]
		m["total_fee"] = data[41]
		m["fee_lines"] = data[42]
		m["shop_cod"] = data[43]
		m["total_amount"] = data[44]
		m["order_note"] = data[45]
		m["shop_note"] = data[46]
		m["shipping_note"] = data[47]
		m["order_source_type"] = data[48]
		m["order_source_id"] = data[49]
		m["external_order_id"] = data[50]
		m["reference_url"] = data[51]
		m["external_url"] = data[52]
		m["shop_shipping"] = data[53]
		m["is_outside_etop"] = data[54]
		m["ghn_note_code"] = data[55]
		m["try_on"] = data[56]
		m["customer_name_norm"] = data[57]
		m["product_name_norm"] = data[58]
		m["fulfillment_type"] = data[59]
		m["fulfillment_ids"] = data[60]
		m["external_meta"] = data[61]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

// Type OrderLine represents table order_line
func sqlgenOrderLine(_ *OrderLine) bool { return true }

type OrderLines []*OrderLine

const __sqlOrderLine_Table = "order_line"
const __sqlOrderLine_ListCols = "\"order_id\",\"variant_id\",\"product_name\",\"product_id\",\"shop_id\",\"weight\",\"quantity\",\"list_price\",\"retail_price\",\"payment_price\",\"line_amount\",\"total_discount\",\"total_line_amount\",\"image_url\",\"is_outside_etop\",\"code\""
const __sqlOrderLine_Insert = "INSERT INTO \"order_line\" (" + __sqlOrderLine_ListCols + ") VALUES"
const __sqlOrderLine_Select = "SELECT " + __sqlOrderLine_ListCols + " FROM \"order_line\""
const __sqlOrderLine_Select_history = "SELECT " + __sqlOrderLine_ListCols + " FROM history.\"order_line\""
const __sqlOrderLine_UpdateAll = "UPDATE \"order_line\" SET (" + __sqlOrderLine_ListCols + ")"

func (m *OrderLine) SQLTableName() string  { return "order_line" }
func (m *OrderLines) SQLTableName() string { return "order_line" }
func (m *OrderLine) SQLListCols() string   { return __sqlOrderLine_ListCols }

func (m *OrderLine) SQLArgs(opts core.Opts, create bool) []interface{} {
	return []interface{}{
		core.Int64(m.OrderID),
		core.Int64(m.VariantID),
		core.String(m.ProductName),
		core.Int64(m.ProductID),
		core.Int64(m.ShopID),
		core.Int(m.Weight),
		core.Int(m.Quantity),
		core.Int(m.ListPrice),
		core.Int(m.RetailPrice),
		core.Int(m.PaymentPrice),
		core.Int(m.LineAmount),
		core.Int(m.TotalDiscount),
		core.Int(m.TotalLineAmount),
		core.String(m.ImageURL),
		core.Bool(m.IsOutsideEtop),
		core.String(m.Code),
	}
}

func (m *OrderLine) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.Int64)(&m.OrderID),
		(*core.Int64)(&m.VariantID),
		(*core.String)(&m.ProductName),
		(*core.Int64)(&m.ProductID),
		(*core.Int64)(&m.ShopID),
		(*core.Int)(&m.Weight),
		(*core.Int)(&m.Quantity),
		(*core.Int)(&m.ListPrice),
		(*core.Int)(&m.RetailPrice),
		(*core.Int)(&m.PaymentPrice),
		(*core.Int)(&m.LineAmount),
		(*core.Int)(&m.TotalDiscount),
		(*core.Int)(&m.TotalLineAmount),
		(*core.String)(&m.ImageURL),
		(*core.Bool)(&m.IsOutsideEtop),
		(*core.String)(&m.Code),
	}
}

func (m *OrderLine) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *OrderLines) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(OrderLines, 0, 128)
	for rows.Next() {
		m := new(OrderLine)
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

func (_ *OrderLine) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_Select)
	return nil
}

func (_ *OrderLines) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_Select)
	return nil
}

func (m *OrderLine) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_Insert)
	w.WriteRawString(" (")
	w.WriteMarkers(16)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), true))
	return nil
}

func (ms OrderLines) SQLInsert(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_Insert)
	w.WriteRawString(" (")
	for i := 0; i < len(ms); i++ {
		w.WriteMarkers(16)
		w.WriteArgs(ms[i].SQLArgs(w.Opts(), true))
		w.WriteRawString("),(")
	}
	w.TrimLast(2)
	return nil
}

func (m *OrderLine) SQLUpdate(w SQLWriter) error {
	now, opts := time.Now(), w.Opts()
	_, _ = now, opts // suppress unuse error
	var flag bool
	w.WriteRawString("UPDATE ")
	w.WriteName("order_line")
	w.WriteRawString(" SET ")
	if m.OrderID != 0 {
		flag = true
		w.WriteName("order_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.OrderID)
	}
	if m.VariantID != 0 {
		flag = true
		w.WriteName("variant_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.VariantID)
	}
	if m.ProductName != "" {
		flag = true
		w.WriteName("product_name")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ProductName)
	}
	if m.ProductID != 0 {
		flag = true
		w.WriteName("product_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ProductID)
	}
	if m.ShopID != 0 {
		flag = true
		w.WriteName("shop_id")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ShopID)
	}
	if m.Weight != 0 {
		flag = true
		w.WriteName("weight")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Weight)
	}
	if m.Quantity != 0 {
		flag = true
		w.WriteName("quantity")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Quantity)
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
	if m.PaymentPrice != 0 {
		flag = true
		w.WriteName("payment_price")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.PaymentPrice)
	}
	if m.LineAmount != 0 {
		flag = true
		w.WriteName("line_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.LineAmount)
	}
	if m.TotalDiscount != 0 {
		flag = true
		w.WriteName("total_discount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalDiscount)
	}
	if m.TotalLineAmount != 0 {
		flag = true
		w.WriteName("total_line_amount")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.TotalLineAmount)
	}
	if m.ImageURL != "" {
		flag = true
		w.WriteName("image_url")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.ImageURL)
	}
	if m.IsOutsideEtop {
		flag = true
		w.WriteName("is_outside_etop")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.IsOutsideEtop)
	}
	if m.Code != "" {
		flag = true
		w.WriteName("code")
		w.WriteByte('=')
		w.WriteMarker()
		w.WriteByte(',')
		w.WriteArg(m.Code)
	}
	if !flag {
		return core.ErrNoColumn
	}
	w.TrimLast(1)
	return nil
}

func (m *OrderLine) SQLUpdateAll(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_UpdateAll)
	w.WriteRawString(" = (")
	w.WriteMarkers(16)
	w.WriteByte(')')
	w.WriteArgs(m.SQLArgs(w.Opts(), false))
	return nil
}

type OrderLineHistory map[string]interface{}
type OrderLineHistories []map[string]interface{}

func (m *OrderLineHistory) SQLTableName() string  { return "history.\"order_line\"" }
func (m OrderLineHistories) SQLTableName() string { return "history.\"order_line\"" }

func (m *OrderLineHistory) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_Select_history)
	return nil
}

func (m OrderLineHistories) SQLSelect(w SQLWriter) error {
	w.WriteQueryString(__sqlOrderLine_Select_history)
	return nil
}

func (m OrderLineHistory) OrderID() core.Interface       { return core.Interface{m["order_id"]} }
func (m OrderLineHistory) VariantID() core.Interface     { return core.Interface{m["variant_id"]} }
func (m OrderLineHistory) ProductName() core.Interface   { return core.Interface{m["product_name"]} }
func (m OrderLineHistory) ProductID() core.Interface     { return core.Interface{m["product_id"]} }
func (m OrderLineHistory) ShopID() core.Interface        { return core.Interface{m["shop_id"]} }
func (m OrderLineHistory) Weight() core.Interface        { return core.Interface{m["weight"]} }
func (m OrderLineHistory) Quantity() core.Interface      { return core.Interface{m["quantity"]} }
func (m OrderLineHistory) ListPrice() core.Interface     { return core.Interface{m["list_price"]} }
func (m OrderLineHistory) RetailPrice() core.Interface   { return core.Interface{m["retail_price"]} }
func (m OrderLineHistory) PaymentPrice() core.Interface  { return core.Interface{m["payment_price"]} }
func (m OrderLineHistory) LineAmount() core.Interface    { return core.Interface{m["line_amount"]} }
func (m OrderLineHistory) TotalDiscount() core.Interface { return core.Interface{m["total_discount"]} }
func (m OrderLineHistory) TotalLineAmount() core.Interface {
	return core.Interface{m["total_line_amount"]}
}
func (m OrderLineHistory) ImageURL() core.Interface      { return core.Interface{m["image_url"]} }
func (m OrderLineHistory) IsOutsideEtop() core.Interface { return core.Interface{m["is_outside_etop"]} }
func (m OrderLineHistory) Code() core.Interface          { return core.Interface{m["code"]} }

func (m *OrderLineHistory) SQLScan(opts core.Opts, row *sql.Row) error {
	data := make([]interface{}, 16)
	args := make([]interface{}, 16)
	for i := 0; i < 16; i++ {
		args[i] = &data[i]
	}
	if err := row.Scan(args...); err != nil {
		return err
	}
	res := make(OrderLineHistory, 16)
	res["order_id"] = data[0]
	res["variant_id"] = data[1]
	res["product_name"] = data[2]
	res["product_id"] = data[3]
	res["shop_id"] = data[4]
	res["weight"] = data[5]
	res["quantity"] = data[6]
	res["list_price"] = data[7]
	res["retail_price"] = data[8]
	res["payment_price"] = data[9]
	res["line_amount"] = data[10]
	res["total_discount"] = data[11]
	res["total_line_amount"] = data[12]
	res["image_url"] = data[13]
	res["is_outside_etop"] = data[14]
	res["code"] = data[15]
	*m = res
	return nil
}

func (ms *OrderLineHistories) SQLScan(opts core.Opts, rows *sql.Rows) error {
	data := make([]interface{}, 16)
	args := make([]interface{}, 16)
	for i := 0; i < 16; i++ {
		args[i] = &data[i]
	}
	res := make(OrderLineHistories, 0, 128)
	for rows.Next() {
		if err := rows.Scan(args...); err != nil {
			return err
		}
		m := make(OrderLineHistory)
		m["order_id"] = data[0]
		m["variant_id"] = data[1]
		m["product_name"] = data[2]
		m["product_id"] = data[3]
		m["shop_id"] = data[4]
		m["weight"] = data[5]
		m["quantity"] = data[6]
		m["list_price"] = data[7]
		m["retail_price"] = data[8]
		m["payment_price"] = data[9]
		m["line_amount"] = data[10]
		m["total_discount"] = data[11]
		m["total_line_amount"] = data[12]
		m["image_url"] = data[13]
		m["is_outside_etop"] = data[14]
		m["code"] = data[15]
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}
