package model

import (
	"encoding/json"
	"time"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/top/types/etc/fee"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	addressmodel "o.o/backend/com/main/address/model"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/capi/dot"
)

// +sqlgen
type Order struct {
	ID         dot.ID
	ShopID     dot.ID
	Code       string
	EdCode     string
	ProductIDs []dot.ID
	VariantIDs []dot.ID

	PaymentMethod payment_method.PaymentMethod `sql_type:"text"`

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

	ShopConfirm   status3.Status `sql_type:"int2"`
	ConfirmStatus status3.Status `sql_type:"int2"`

	FulfillmentShippingStatus status5.Status `sql_type:"int2"`
	EtopPaymentStatus         status4.Status `sql_type:"int2"`

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status status5.Status `sql_type:"int2"`

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
	ExternalOrderID string
	ReferenceURL    string
	ExternalURL     string
	ShopShipping    *OrderShipping
	IsOutsideEtop   bool

	// @deprecated: use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode `sql_type:"enum(ghn_note_code)"`
	TryOn       try_on.TryOnCode          `sql_type:"enum(try_on)"`

	FulfillmentType ordertypes.ShippingType `sql_type:"int2"`
	FulfillmentIDs  []dot.ID
	ExternalMeta    json.RawMessage

	// payment
	PaymentStatus status4.Status `sql_type:"int2"`
	PaymentID     dot.ID

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
