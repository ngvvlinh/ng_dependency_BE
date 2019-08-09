package gateway

import (
	"context"

	"etop.vn/api/external/haravan"
)

type Aggregate interface {
	GetShippingRate(context.Context, *GetShippingRateRequestArgs) (*GetShippingRateResponse, error)

	CreateOrder(context.Context, *CreateOrderRequestArgs) (*CreateOrderResponse, error)

	GetOrder(context.Context, *GetOrderRequestArgs) (*GetOrderResponse, error)

	CancelOrder(context.Context, *CancelOrderRequestArgs) (*GetOrderResponse, error)
}

type GetShippingRateRequestArgs struct {
	EtopShopID  int64
	Origin      *haravan.Address `json:"origin"`
	Destination *haravan.Address `json:"destination"`
	CodAmount   float32          `json:"cod_amount"`
	TotalGrams  float32          `json:"total_grams"`
}

type GetShippingRateResponse struct {
	ShippingRates []*haravan.ShippingRate `json:"rates"`
}

type CreateOrderRequestArgs struct {
	EtopShopID  int64
	Origin      *haravan.Address `json:"origin"`
	Destination *haravan.Address `json:"destination"`
	Items       []*haravan.Item  `json:"items"`
	CodAmount   float32          `json:"cod_amount"`
	TotalGrams  float32          `json:"total_grams"`
	// ExternalStoreID: Mã shop
	ExternalStoreID int `json:"external_store_id"`
	// ExternalOrderID: Order id của shop
	ExternalOrderID int `json:"external_order_id"`
	// ExternalFulfillmentID: Mã vận đơn
	ExternalFulfillmentID int `json:"external_fulfillment_id"`
	// ExternalCode = external_store_id + ”_” + external_order_id + “_” + external_fulfillment_id
	ExternalCode string `json:"external_code"`
	Note         string `json:"note"`
	// ShippingRateID: Là service_id nhận được khi post lấy danh sách gói vận chuyển
	ShippingRateID int32 `json:"shipping_rate_id"`
}

type CreateOrderResponse struct {
	TrackingNumber string `json:"tracking_number"`
	ShippingFee    int32  `json:"shipping_fee"`
	TrackingURL    string `json:"tracking_url"`
	CodAmount      int32  `json:"cod_amount"`
}

type GetOrderRequestArgs struct {
	EtopShopID     int64
	TrackingNumber string `json:"tracking_number"`
}

type GetOrderResponse struct {
	TrackingNumber string `json:"tracking_number"`
	ShippingFee    int32  `json:"shipping_fee"`
	TrackingURL    string `json:"tracking_url"`
	CodAmount      int32  `json:"cod_amount"`
	Status         string `json:"status"`
	CodStatus      string `json:"cod_status"`
}

type CancelOrderRequestArgs struct {
	EtopShopID     int64
	TrackingNumber string `json:"tracking_number"`
}

type CancelOrderResponse struct {
	Status bool `json:"status"`
}
