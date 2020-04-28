package types

import (
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/top/int/types"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
)

type Fulfillment struct {
	Id                       dot.ID                             `json:"id"`
	OrderId                  dot.ID                             `json:"order_id"`
	ShopId                   dot.ID                             `json:"shop_id"`
	SelfUrl                  dot.NullString                     `json:"self_url"`
	TotalItems               dot.NullInt                        `json:"total_items"`
	BasketValue              dot.NullInt                        `json:"basket_value"`
	CreatedAt                dot.Time                           `json:"created_at"`
	UpdatedAt                dot.Time                           `json:"updated_at"`
	ClosedAt                 dot.Time                           `json:"closed_at"`
	CancelledAt              dot.Time                           `json:"cancelled_at"`
	CancelReason             dot.NullString                     `json:"cancel_reason"`
	Carrier                  shipping_provider.ShippingProvider `json:"carrier"`
	ShippingServiceName      dot.NullString                     `json:"shipping_service_name"`
	ShippingServiceFee       dot.NullInt                        `json:"shipping_service_fee"`
	ActualShippingServiceFee dot.NullInt                        `json:"actual_shipping_service_fee"`
	ShippingServiceCode      dot.NullString                     `json:"shipping_service_code"`
	ShippingCode             dot.NullString                     `json:"shipping_code"`
	ShippingNote             dot.NullString                     `json:"shipping_note"`
	TryOn                    try_on.TryOnCode                   `json:"try_on"`
	IncludeInsurance         dot.NullBool                       `json:"include_insurance"`
	ConfirmStatus            status3.NullStatus                 `json:"confirm_status"`
	ShippingState            shipping.NullState                 `json:"shipping_state"`
	ShippingStatus           status5.NullStatus                 `json:"shipping_status"`
	Status                   status5.NullStatus                 `json:"status"`
	CodAmount                dot.NullInt                        `json:"cod_amount"`
	ActualCodAmount          dot.NullInt                        `json:"actual_cod_amount"`
	ChargeableWeight         dot.NullInt                        `json:"chargeable_weight"`
	PickupAddress            *OrderAddress                      `json:"pickup_address"`
	ReturnAddress            *OrderAddress                      `json:"return_address"`
	ShippingAddress          *OrderAddress                      `json:"shipping_address"`
	EtopPaymentStatus        status4.NullStatus                 `json:"etop_payment_status"`
	EstimatedDeliveryAt      dot.Time                           `json:"estimated_delivery_at"`
	EstimatedPickupAt        dot.Time                           `json:"estimated_pickup_at"`
}

func (m *Fulfillment) String() string { return jsonx.MustMarshalToString(m) }

type FulfillmentsResponse struct {
	Fulfillments []*Fulfillment         `json:"fulfillments"`
	Paging       *common.CursorPageInfo `json:"paging"`
}

func (m *FulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type FulfillmentIDRequest struct {
	Id           dot.ID `json:"id"`
	ShippingCode string `json:"shipping_code"`
}

func (m *FulfillmentIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListFulfillmentsFilter struct {
	OrderID filter.IDs `json:"order_id"`
}

func (m *ListFulfillmentsFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListFulfillmentsRequest struct {
	Filter ListFulfillmentsFilter `json:"filter"`
	Paging *common.CursorPaging   `json:"paging"`
}

func (m *ListFulfillmentsRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *Fulfillment) HasChanged() bool {
	return m.Status.Valid ||
		m.ShippingState.Valid ||
		m.EtopPaymentStatus.Valid ||
		m.ActualShippingServiceFee.Valid ||
		m.CodAmount.Valid ||
		m.ActualCodAmount.Valid ||
		m.ShippingNote.Valid ||
		m.ChargeableWeight.Valid
}

type CreateFulfillmentRequest struct {
	OrderID             dot.ID                  `json:"order_id"`
	ShippingType        ordertypes.ShippingType `json:"shipping_type"`
	ShippingServiceCode string                  `json:"shipping_service_code"`
	ShippingServiceFee  int                     `json:"shipping_service_fee"`
	ShippingServiceName string                  `json:"shipping_service_name"`
	ShippingNote        string                  `json:"shipping_note"`
	PickupAddress       *types.OrderAddress     `json:"pickup_address"`
	ReturnAddress       *types.OrderAddress     `json:"return_address"`
	ShippingAddress     *types.OrderAddress     `json:"shipping_address"`
	TryOn               try_on.TryOnCode        `json:"try_on"`
	ChargeableWeight    int                     `json:"chargeable_weight"`
	GrossWeight         int                     `json:"gross_weight"`
	Height              int                     `json:"height"`
	Width               int                     `json:"width"`
	Length              int                     `json:"length"`
	CODAmount           int                     `json:"cod_amount"`
	IncludeInsurance    bool                    `json:"include_insurance"`

	ShopCarrierID dot.ID `json:"shop_carrier_id"`
}

func (m *CreateFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }

type CancelFulfillmentRequest struct {
	FulfillmentID dot.ID `json:"fulfillment_id"`
	CancelReason  string `json:"cancel_reason"`
}

func (m *CancelFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }
