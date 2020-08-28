package types

import (
	inttypes "o.o/api/top/int/types"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type ShipnowAddressShortVersion struct {
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Address  string `json:"address"`
	// Vui lòng cung cấp lat. , long. để lấy được giá chính xác nhất
	Coordinates *Coordinates `json:"coordinates"`
}

func (m *ShipnowAddressShortVersion) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowAddress struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Address  string `json:"address"`
	// Vui lòng cung cấp lat. , long. để lấy được giá chính xác nhất
	Coordinates *Coordinates `json:"coordinates"`
}

func (m *ShipnowAddress) String() string { return jsonx.MustMarshalToString(m) }

type GetShipnowServicesRequest struct {
	PickupAddress  *ShipnowAddressShortVersion         `json:"pickup_address"`
	DeliveryPoints []*ShipnowDeliveryPointShortVersion `json:"delivery_points"`
	Coupon         string                              `json:"coupon"`
}

func (m *GetShipnowServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowDeliveryPointShortVersion struct {
	ShippingAddress *ShipnowAddressShortVersion `json:"shipping_address"`
	CODAmount       inttypes.Int                `json:"cod_amount"`
}

func (m *ShipnowDeliveryPointShortVersion) String() string { return jsonx.MustMarshalToString(m) }

type GetShipnowServicesResponse struct {
	Services []*ShipnowService `json:"services"`
}

func (m *GetShipnowServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowService struct {
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Fee         int          `json:"fee"`
	CarrierInfo *CarrierInfo `json:"carrier_info"`
	Description string       `json:"description"`
}

func (m *ShipnowService) String() string { return jsonx.MustMarshalToString(m) }

type CreateShipnowFulfillmentRequest struct {
	// @required
	PickupAddress *ShipnowAddress `json:"pickup_address"`
	// @required
	DeliveryPoints []*ShipnowDeliveryPointRequest `json:"delivery_points"`
	ExternalID     string                         `json:"external_id"`
	// @required
	ShippingServiceCode string `json:"shipping_service_code"`
	// @required
	ShippingServiceFee inttypes.Int `json:"shipping_service_fee"`
	ShippingNote       string       `json:"shipping_note"`
	Coupon             string       `json:"coupon"`
}

func (m *CreateShipnowFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowDeliveryPointRequest struct {
	ChargeableWeight inttypes.Int `json:"chargeable_weight"`
	// @required
	GrossWeight inttypes.Int `json:"gross_weight"`
	// @required
	CODAmount       inttypes.Int    `json:"cod_amount"`
	ShippingNote    string          `json:"shipping_note"`
	ShippingAddress *ShipnowAddress `json:"shipping_address"`
	// @required
	BasketValue inttypes.Int `json:"basket_value"`
	// @required
	Lines []*OrderLine `json:"lines"`
}

func (m *ShipnowDeliveryPointRequest) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowDeliveryPoint struct {
	ChargeableWeight inttypes.Int `json:"chargeable_weight"`
	// @required
	GrossWeight inttypes.Int `json:"gross_weight"`
	// @required
	CODAmount       inttypes.Int    `json:"cod_amount"`
	ShippingNote    string          `json:"shipping_note"`
	ShippingAddress *ShipnowAddress `json:"shipping_address"`
	BasketValue     inttypes.Int    `json:"basket_value"`
	Lines           []*OrderLine    `json:"lines"`
	// Map với trạng thái của đơn hàng
	ShippingState shipnow_state.State `json:"shipping_state"`
}

func (m *ShipnowDeliveryPoint) String() string { return jsonx.MustMarshalToString(m) }

type ShipnowFulfillment struct {
	ID                         dot.ID                  `json:"id"`
	ShopID                     dot.ID                  `json:"shop_id"`
	PickupAddress              *ShipnowAddress         `json:"pickup_address"`
	DeliveryPoints             []*ShipnowDeliveryPoint `json:"delivery_points"`
	ShippingServiceCode        dot.NullString          `json:"shipping_service_code"`
	ShippingServiceFee         dot.NullInt             `json:"shipping_service_fee"`
	ActualShippingServiceFee   dot.NullInt             `json:"actual_shipping_service_fee"`
	ShippingServiceName        dot.NullString          `json:"shipping_service_name"`
	ShippingServiceDescription dot.NullString          `json:"shipping_service_description"`
	GrossWeight                dot.NullInt             `json:"gross_weight"`
	ChargeableWeight           dot.NullInt             `json:"chargeable_weight"`
	BasketValue                dot.NullInt             `json:"basket_value"`
	CODAmount                  dot.NullInt             `json:"cod_amount"`
	ShippingNote               dot.NullString          `json:"shipping_note"`
	Status                     status5.NullStatus      `json:"status"`
	ShippingStatus             status5.NullStatus      `json:"shipping_status"`
	ShippingCode               dot.NullString          `json:"shipping_code"`
	ShippingState              shipnow_state.NullState `json:"shipping_state"`
	ConfirmStatus              status3.NullStatus      `json:"confirm_status"`
	OrderIDs                   []dot.ID                `json:"order_ids"`
	CreatedAt                  dot.Time                `json:"created_at"`
	UpdatedAt                  dot.Time                `json:"updated_at"`
	ShippingSharedLink         dot.NullString          `json:"shipping_shared_link"`
	CancelReason               dot.NullString          `json:"cancel_reason"`
	CarrierInfo                *CarrierInfo            `json:"carrier_info"`
	ExternalID                 dot.NullString          `json:"external_id"`
	Coupon                     dot.NullString          `json:"coupon"`
}

func (m *ShipnowFulfillment) String() string { return jsonx.MustMarshalToString(m) }

func (m *ShipnowFulfillment) HasChanged() bool {
	return m.Status.Valid ||
		m.ShippingState.Valid ||
		m.CODAmount.Valid ||
		m.ActualShippingServiceFee.Valid ||
		m.ShippingNote.Valid ||
		m.ChargeableWeight.Valid ||
		m.DeliveryPoints != nil
}

type CancelShipnowFulfillmentRequest struct {
	ID           dot.ID `json:"id"`
	ShippingCode string `json:"shipping_code"`
	ExternalID   string `json:"external_id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelShipnowFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }
