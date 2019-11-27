package haravan

import (
	"time"

	"etop.vn/capi/dot"
)

var HaravanPartnerID dot.ID = 1000421281650350414

type Address struct {
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	Province     string `json:"province"`
	ProvinceCode string `json:"province_code"`
	District     string `json:"district"`
	DistrictCode string `json:"district_code"`
	Ward         string `json:"ward"`
	WardCode     string `json:"ward_code"`
	Address1     string `json:"address1"`
	Address2     string `json:"address2"`
	Zip          string `json:"zip"`
	City         string `json:"city"`
	Phone        string `json:"phone"`
	Name         string `json:"name"`
}

type ShippingRate struct {
	ServiceID       int       `json:"service_id"` // use int
	ServiceName     string    `json:"service_name"`
	ServiceCode     string    `json:"service_code"`
	Currency        string    `json:"currency"`
	TotalPrice      int       `json:"total_price"`
	PhoneRequired   bool      `json:"phone_required"`
	MinDeliveryDate time.Time `json:"min_delivery_date"`
	MaxDeliveryDate time.Time `json:"max_delivery_date"`
	Description     string    `json:"description"`
}

type Item struct {
	Name      string  `json:"name"`
	Sku       string  `json:"sku"`
	Quantity  int     `json:"quantity"`
	Grams     float32 `json:"grams"`
	Price     float32 `json:"price"`
	ProductID int     `json:"product_id"`
	VariantID int     `json:"variant_id"`
}

type ExternalMeta struct {
	ExternalOrderID       string `json:"external_order_id"`
	ExternalFulfillmentID string `json:"external_fulfillment_id"`
}
