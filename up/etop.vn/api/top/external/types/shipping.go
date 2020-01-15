package types

import (
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type GetShippingServicesRequest struct {
	ConnectionIDs   []dot.ID         `json:"connection_ids"`
	PickupAddress   *LocationAddress `json:"pickup_address"`
	ShippingAddress *LocationAddress `json:"shipping_address"`
	// in gram (g)
	GrossWeight int `json:"gross_weight"`
	// in gram (g)
	ChargeableWeight int `json:"chargeable_weight"`
	// in centimetre (cm)
	Length int `json:"length"`
	// in centimetre (cm)
	Width int `json:"width"`
	// in centimetre (cm)
	Height           int          `json:"height"`
	BasketValue      int          `json:"basket_value"`
	CodAmount        int          `json:"cod_amount"`
	IncludeInsurance dot.NullBool `json:"include_insurance"`
}

func (m *GetShippingServicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShippingServicesResponse struct {
	Services []*ShippingService `json:"services"`
}

func (m *GetShippingServicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShippingService struct {
	Code string `json:"code"`
	// @deprecated use carrier info instead
	Name string `json:"name"`
	Fee  int    `json:"fee"`
	// @deprecated
	Carrier             shipping_provider.ShippingProvider `json:"carrier"`
	EstimatedPickupAt   dot.Time                           `json:"estimated_pickup_at"`
	EstimatedDeliveryAt dot.Time                           `json:"estimated_delivery_at"`
	CarrierInfo         *CarrierInfo                       `json:"carrier_info"`
}

func (m *ShippingService) String() string { return jsonx.MustMarshalToString(m) }

type CarrierInfo struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func (m *CarrierInfo) String() string { return jsonx.MustMarshalToString(m) }
