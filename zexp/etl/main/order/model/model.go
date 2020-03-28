package model

import (
	"encoding/json"
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/top/types/etc/fee"
	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/order_source"
	"etop.vn/api/top/types/etc/payment_method"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	addressmodel "etop.vn/backend/com/main/address/model"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/capi/dot"
)

// +sqlgen
type Order struct {
	ID         dot.ID
	ShopID     dot.ID
	Code       string
	EdCode     string
	ProductIDs []dot.ID
	VariantIDs []dot.ID
	PartnerID  dot.ID

	Currency      string
	PaymentMethod payment_method.PaymentMethod

	Customer        *OrderCustomer
	CustomerAddress *OrderAddress
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	CustomerName    string
	CustomerPhone   string
	CustomerEmail   string

	CreatedAt    time.Time
	ProcessedAt  time.Time
	UpdatedAt    time.Time
	ClosedAt     time.Time
	ConfirmedAt  time.Time
	CancelledAt  time.Time
	CancelReason string

	CustomerConfirm status3.Status
	ShopConfirm     status3.Status
	ConfirmStatus   status3.Status

	FulfillmentShippingStatus status5.Status
	EtopPaymentStatus         status4.Status

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status status5.Status

	FulfillmentShippingStates  []string
	FulfillmentPaymentStatuses []int
	FulfillmentStatuses        []int

	Lines           OrderLinesList
	Discounts       []*OrderDiscount
	TotalItems      int
	BasketValue     int
	TotalWeight     int
	TotalTax        int
	OrderDiscount   int
	TotalDiscount   int
	ShopShippingFee int
	TotalFee        int
	FeeLines        OrderFeeLines
	ShopCOD         int
	TotalAmount     int

	OrderNote       string
	ShopNote        string
	ShippingNote    string
	OrderSourceType order_source.Source
	OrderSourceID   dot.ID
	ExternalOrderID string
	ReferenceURL    string
	ExternalURL     string
	ShopShipping    *OrderShipping
	IsOutsideEtop   bool

	// @deprecated: use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode
	TryOn       try_on.TryOnCode

	CustomerNameNorm string
	ProductNameNorm  string
	FulfillmentType  ordertypes.ShippingType
	FulfillmentIDs   []dot.ID
	ExternalMeta     json.RawMessage
	TradingShopID    dot.ID

	// payment
	PaymentStatus status4.Status
	PaymentID     dot.ID

	ReferralMeta json.RawMessage

	CustomerID dot.ID
	CreatedBy  dot.ID

	Rid dot.ID
}

type OrderAddress struct {
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`

	Country  string `json:"country"`
	City     string `json:"city"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
	Zip      string `json:"zip"`

	DistrictCode string `json:"district_code"`
	ProvinceCode string `json:"province_code"`
	WardCode     string `json:"ward_code"`

	Company     string                    `json:"company"`
	Address1    string                    `json:"address1"`
	Address2    string                    `json:"address2"`
	Coordinates *addressmodel.Coordinates `json:"coordinates"`
}

type OrderCustomer struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Gender        string `json:"gender"`
	Birthday      string `json:"birthday"`
	VerifiedEmail bool   `json:"verified_email"`

	ExternalID string `json:"external_id"`
}

type OrderLinesList []*OrderLine

type OrderLine struct {
	OrderID     dot.ID `json:"order_id"`
	VariantID   dot.ID `json:"variant_id"`
	ProductName string `json:"product_name"`
	ProductID   dot.ID `json:"product_id"`
	ShopID      dot.ID `json:"shop_id"`

	Weight       int `json:"weight"`
	Quantity     int `json:"quantity"`
	ListPrice    int `json:"list_price"`
	RetailPrice  int `json:"retail_price"`
	PaymentPrice int `json:"payment_price"`

	LineAmount      int `json:"line_amount"`
	TotalDiscount   int `json:"discount"`
	TotalLineAmount int `json:"total_line_amount"`

	ImageURL      string                           `json:"image_url" `
	Attributes    []*catalogmodel.ProductAttribute `json:"attributes" sq:"-"`
	IsOutsideEtop bool                             `json:"is_outside_etop"`
	Code          string                           `json:"code"`
	MetaFields    []*MetaField                     `json:"meta_fields"`
}

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Name  string `json:"name"`
}

type OrderDiscount struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type OrderFeeLines []OrderFeeLine

type OrderFeeLine struct {
	Amount int         `json:"amount"`
	Desc   string      `json:"desc"`
	Code   string      `json:"code"`
	Name   string      `json:"name"`
	Type   fee.FeeType `json:"type"`
}

type OrderShipping struct {
	ShopAddress         *OrderAddress                      `json:"shop_address"`
	ReturnAddress       *OrderAddress                      `json:"return_address"`
	ExternalServiceID   string                             `json:"external_service_id"`
	ExternalShippingFee int                                `json:"external_shipping_fee"`
	ExternalServiceName string                             `json:"external_service_name"`
	ShippingProvider    shipping_provider.ShippingProvider `json:"shipping_provider"`
	ProviderServiceID   string                             `json:"provider_service_id"`
	IncludeInsurance    bool                               `json:"include_insurance"`

	Length int `json:"length"`
	Width  int `json:"width"`
	Height int `json:"height"`

	GrossWeight      int `json:"gross_weight"`
	ChargeableWeight int `json:"chargeable_weight"`
}
