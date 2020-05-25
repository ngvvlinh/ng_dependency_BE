package pricelist

import (
	"context"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShipmentPriceList(context.Context, *CreateShipmentPriceListArg) (*ShipmentPriceList, error)

	UpdateShipmentPriceList(context.Context, *UpdateShipmentPriceListArgs) error

	ActivateShipmentPriceList(ctx context.Context, ID dot.ID) error

	DeleteShipmentPriceList(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetShipmentPriceList(ctx context.Context, ID dot.ID) (*ShipmentPriceList, error)

	GetActiveShipmentPriceList(context.Context, *meta.Empty) (*ShipmentPriceList, error)

	ListShipmentPriceLists(context.Context, *ListShipmentPriceListsArgs) ([]*ShipmentPriceList, error)
}

// +convert:create=ShipmentPriceList
type CreateShipmentPriceListArg struct {
	Name                    string
	Description             string
	IsActive                bool
	ShipmentSubPriceListIDs []dot.ID
}

type ListShipmentPriceListsArgs struct {
	SubShipmentPriceListIDs []dot.ID
}

// +convert:update=ShipmentPriceList
type UpdateShipmentPriceListArgs struct {
	ID                      dot.ID
	Name                    string
	Description             string
	ShipmentSubPriceListIDs []dot.ID
}
