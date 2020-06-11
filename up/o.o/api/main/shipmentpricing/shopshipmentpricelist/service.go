package shopshipmentpricelist

import (
	"context"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShopShipmentPriceList(context.Context, *CreateShopShipmentPriceListArgs) (*ShopShipmentPriceList, error)

	UpdateShopShipmentPriceList(context.Context, *UpdateShopShipmentPriceListArgs) error

	DeleteShopShipmentPriceList(ctx context.Context, ShopID, ConnectionID dot.ID) error
}

type QueryService interface {
	ListShopShipmentPriceLists(context.Context, *GetShopShipmentPriceListsArgs) (*GetShopShipmentPriceListsResponse, error)

	ListShopShipmentPriceListsByPriceListIDs(ctx context.Context, PriceListIDs []dot.ID) ([]*ShopShipmentPriceList, error)

	GetShopShipmentPriceList(ctx context.Context, ShopID, ConnectionID dot.ID) (*ShopShipmentPriceList, error)
}

// +convert:create=ShopShipmentPriceList
type CreateShopShipmentPriceListArgs struct {
	ShopID              dot.ID
	ShipmentPriceListID dot.ID
	Note                string
	UpdatedBy           dot.ID
	ConnectionID        dot.ID
}

// +convert:update=ShopShipmentPriceList
type UpdateShopShipmentPriceListArgs struct {
	ShopID              dot.ID
	ShipmentPriceListID dot.ID
	ConnectionID        dot.ID
	Note                string
	UpdatedBy           dot.ID
}

type GetShopShipmentPriceListsArgs struct {
	ShipmentPriceListID dot.ID
	ConnectionID        dot.ID
	ShopID              dot.ID
	Paging              meta.Paging
}

type GetShopShipmentPriceListsResponse struct {
	ShopShipmentPriceLists []*ShopShipmentPriceList
	Paging                 meta.PageInfo
}
