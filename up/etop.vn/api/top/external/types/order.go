package types

import (
	catalogtypes "etop.vn/api/main/catalog/types"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/gender"
	shipping_provider "etop.vn/api/top/types/etc/shipping_provider"
	status3 "etop.vn/api/top/types/etc/status3"
	status4 "etop.vn/api/top/types/etc/status4"
	status5 "etop.vn/api/top/types/etc/status5"
	try_on "etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
	"etop.vn/common/jsonx"
)

type OrderCustomer struct {
	FullName string        `json:"full_name"`
	Email    string        `json:"email"`
	Phone    string        `json:"phone"`
	Gender   gender.Gender `json:"gender"`
}

func (m *OrderCustomer) String() string { return jsonx.MustMarshalToString(m) }

type OrderAddress struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Company  string `json:"company"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
}

func (m *OrderAddress) String() string { return jsonx.MustMarshalToString(m) }

func (m *Order) HasChanged() bool {
	return m.Status.Valid ||
		m.ConfirmStatus.Valid ||
		m.FulfillmentShippingStatus.Valid ||
		m.EtopPaymentStatus.Valid ||
		m.BasketValue.Valid ||
		m.TotalAmount.Valid ||
		m.Shipping != nil ||
		m.CustomerAddress != nil || m.ShippingAddress != nil
}

type CancelOrderRequest struct {
	Id           dot.ID `json:"id"`
	Code         string `json:"code"`
	ExternalId   string `json:"external_id"`
	CancelReason string `json:"cancel_reason"`
}

func (m *CancelOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderIDRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *OrderIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListOrdersFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *ListOrdersFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListOrdersRequest struct {
	Filter ListOrdersFilter     `json:"filter"`
	Paging *common.CursorPaging `json:"paging"`
}

func (m *ListOrdersRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrdersResponse struct {
	Orders []*Order               `json:"orders"`
	Paging *common.CursorPageInfo `json:"paging"`
}

func (m *OrdersResponse) String() string { return jsonx.MustMarshalToString(m) }

type Order struct {
	Id                        dot.ID                `json:"id"`
	ShopId                    dot.ID                `json:"shop_id"`
	Code                      dot.NullString        `json:"code"`
	ExternalId                dot.NullString        `json:"external_id"`
	ExternalCode              dot.NullString        `json:"external_code"`
	ExternalUrl               dot.NullString        `json:"external_url"`
	SelfUrl                   dot.NullString        `json:"self_url"`
	CustomerAddress           *OrderAddress         `json:"customer_address"`
	ShippingAddress           *OrderAddress         `json:"shipping_address"`
	CreatedAt                 dot.Time              `json:"created_at"`
	ProcessedAt               dot.Time              `json:"processed_at"`
	UpdatedAt                 dot.Time              `json:"updated_at"`
	ClosedAt                  dot.Time              `json:"closed_at"`
	ConfirmedAt               dot.Time              `json:"confirmed_at"`
	CancelledAt               dot.Time              `json:"cancelled_at"`
	CancelReason              dot.NullString        `json:"cancel_reason"`
	ConfirmStatus             status3.NullStatus    `json:"confirm_status"`
	Status                    status5.NullStatus    `json:"status"`
	FulfillmentShippingStatus status5.NullStatus    `json:"fulfillment_shipping_status"`
	EtopPaymentStatus         status4.NullStatus    `json:"etop_payment_status"`
	Lines                     []*OrderLine          `json:"lines"`
	TotalItems                dot.NullInt           `json:"total_items"`
	BasketValue               dot.NullInt           `json:"basket_value"`
	OrderDiscount             dot.NullInt           `json:"order_discount"`
	TotalDiscount             dot.NullInt           `json:"total_discount"`
	TotalFee                  dot.NullInt           `json:"total_fee"`
	FeeLines                  []*types.OrderFeeLine `json:"fee_lines"`
	TotalAmount               dot.NullInt           `json:"total_amount"`
	OrderNote                 dot.NullString        `json:"order_note"`
	Shipping                  *OrderShipping        `json:"shipping"`
}

func (m *Order) String() string { return jsonx.MustMarshalToString(m) }

type OrderShipping struct {
	PickupAddress       *OrderAddress  `json:"pickup_address"`
	ReturnAddress       *OrderAddress  `json:"return_address"`
	ShippingServiceName dot.NullString `json:"shipping_service_name"`
	ShippingServiceCode dot.NullString `json:"shipping_service_code"`
	ShippingServiceFee  dot.NullInt    `json:"shipping_service_fee"`
	// @Deprecated use connection_id instead
	Carrier          shipping_provider.ShippingProvider `json:"carrier"`
	IncludeInsurance dot.NullBool                       `json:"include_insurance"`
	TryOn            try_on.TryOnCode                   `json:"try_on"`
	ShippingNote     dot.NullString                     `json:"shipping_note"`
	CodAmount        dot.NullInt                        `json:"cod_amount"`
	GrossWeight      dot.NullInt                        `json:"gross_weight"`
	Length           dot.NullInt                        `json:"length"`
	Width            dot.NullInt                        `json:"width"`
	Height           dot.NullInt                        `json:"height"`
	ChargeableWeight dot.NullInt                        `json:"chargeable_weight"`
	ConnectionID     dot.ID                             `json:"connection_id"`
}

func (m *OrderShipping) String() string { return jsonx.MustMarshalToString(m) }

type CreateOrderRequest struct {
	ExternalId      string        `json:"external_id"`
	ExternalCode    string        `json:"external_code"`
	ExternalUrl     string        `json:"external_url"`
	CustomerAddress *OrderAddress `json:"customer_address"`
	ShippingAddress *OrderAddress `json:"shipping_address"`
	Lines           []*OrderLine  `json:"lines"`
	TotalItems      int           `json:"total_items"`
	// basket_value = SUM(lines.retail_price)
	BasketValue   int                   `json:"basket_value"`
	OrderDiscount int                   `json:"order_discount"`
	TotalDiscount int                   `json:"total_discount"`
	TotalFee      int                   `json:"total_fee"`
	FeeLines      []*types.OrderFeeLine `json:"fee_lines"`
	TotalAmount   int                   `json:"total_amount"`
	OrderNote     string                `json:"order_note"`
	Shipping      *OrderShipping        `json:"shipping"`
	ExternalMeta  map[string]string     `json:"external_meta"`
}

func (m *CreateOrderRequest) String() string { return jsonx.MustMarshalToString(m) }

type OrderLine struct {
	VariantId   dot.ID `json:"variant_id"`
	ProductId   dot.ID `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	ListPrice   int    `json:"list_price"`
	RetailPrice int    `json:"retail_price"`
	// payment_price = retail_price - discount_per_item
	PaymentPrice int                       `json:"payment_price"`
	ImageUrl     string                    `json:"image_url"`
	Attributes   []*catalogtypes.Attribute `json:"attributes"`
}

func (m *OrderLine) String() string { return jsonx.MustMarshalToString(m) }

type OrderAndFulfillments struct {
	Order             *Order          `json:"order"`
	Fulfillments      []*Fulfillment  `json:"fulfillments"`
	FulfillmentErrors []*common.Error `json:"fulfillment_errors"`
}

func (m *OrderAndFulfillments) String() string { return jsonx.MustMarshalToString(m) }

type ConfirmOrderRequest struct {
	OrderId dot.ID `json:"order_id"`
	// enum ('create', 'confirm')
	AutoInventoryVoucher dot.NullString `json:"auto_inventory_voucher"`
	// enum ('obey', 'ignore')
	InventoryPolicy bool `json:"inventory_policy"`
}

func (m *ConfirmOrderRequest) String() string { return jsonx.MustMarshalToString(m) }
