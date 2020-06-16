package shipmentprice

import (
	"context"

	"o.o/api/top/types/etc/route_type"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShipmentPrice(context.Context, *CreateShipmentPriceArgs) (*ShipmentPrice, error)
	UpdateShipmentPrice(context.Context, *UpdateShipmentPriceArgs) (*ShipmentPrice, error)
	DeleteShipmentPrice(ctx context.Context, ID dot.ID) error
	UpdateShipmentPricesPriorityPoint(context.Context, *UpdateShipmentPricesPriorityPointArgs) (updated int, err error)
}

type QueryService interface {
	ListShipmentPrices(context.Context, *ListShipmentPricesArgs) ([]*ShipmentPrice, error)
	GetShipmentPrice(ctx context.Context, ID dot.ID) (*ShipmentPrice, error)
	CalculatePrice(context.Context, *CalculatePriceArgs) (*CalculatePriceResult, error)
	CalculateShippingFees(context.Context, *CalculatePriceArgs) (*CalculateShippingFeesResponse, error)
}

// +convert:create=ShipmentPrice
type CreateShipmentPriceArgs struct {
	Name                string
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID
	CustomRegionTypes   []route_type.CustomRegionRouteType
	CustomRegionIDs     []dot.ID
	RegionTypes         []route_type.RegionRouteType
	ProvinceTypes       []route_type.ProvinceRouteType
	UrbanTypes          []route_type.UrbanType
	PriorityPoint       int
	Details             []*PricingDetail
	AdditionalFees      []*AdditionalFee
}

// +convert:update=ShipmentPrice
type UpdateShipmentPriceArgs struct {
	ID                  dot.ID
	Name                string
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID
	CustomRegionTypes   []route_type.CustomRegionRouteType
	CustomRegionIDs     []dot.ID
	RegionTypes         []route_type.RegionRouteType
	ProvinceTypes       []route_type.ProvinceRouteType
	UrbanTypes          []route_type.UrbanType
	PriorityPoint       int
	Details             []*PricingDetail
	AdditionalFees      []*AdditionalFee
	Status              status3.Status
}

type ListShipmentPricesArgs struct {
	ShipmentPriceListID dot.ID
	ShipmentServiceID   dot.ID
}

type CalculatePriceArgs struct {
	AccountID           dot.ID
	FromProvince        string
	FromProvinceCode    string
	FromDistrict        string
	FromDistrictCode    string
	ToProvince          string
	ToProvinceCode      string
	ToDistrict          string
	ToDistrictCode      string
	ShipmentServiceID   dot.ID
	ConnectionID        dot.ID
	ShipmentPriceListID dot.ID
	Weight              int
	BasketValue         int
	CODAmount           int
	IncludeInsurance    bool
}

type CalculatePriceResult struct {
	ShipmentPriceID dot.ID
	Price           int
}

type CalculateShippingFeesResponse struct {
	ShipmentPriceID dot.ID
	TotalFee        int
	FeeLines        []*ShippingFee
}

type ShippingFee struct {
	FeeType shipping_fee_type.ShippingFeeType
	Price   int
}

type UpdateShipmentPricePriorityPointArgs struct {
	ID            dot.ID
	PriorityPoint int
}

type UpdateShipmentPricesPriorityPointArgs struct {
	ShipmentPrices []*UpdateShipmentPricePriorityPointArgs
}
