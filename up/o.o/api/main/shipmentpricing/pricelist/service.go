package pricelist

import (
	"context"

	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShipmentPriceList(context.Context, *CreateShipmentPriceListArg) (*ShipmentPriceList, error)

	UpdateShipmentPriceList(context.Context, *UpdateShipmentPriceListArgs) error

	SetDefaultShipmentPriceList(ctx context.Context, ID dot.ID, connectionID dot.ID) error

	DeleteShipmentPriceList(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetShipmentPriceList(ctx context.Context, ID dot.ID) (*ShipmentPriceList, error)

	GetActiveShipmentPriceList(ctx context.Context, ConnectionID dot.ID) (*ShipmentPriceList, error)

	ListShipmentPriceLists(context.Context, *ListShipmentPriceListsArgs) ([]*ShipmentPriceList, error)
}

// +convert:create=ShipmentPriceList
type CreateShipmentPriceListArg struct {
	Name         string
	Description  string
	IsDefault    bool
	ConnectionID dot.ID
}

type ListShipmentPriceListsArgs struct {
	ConnectionID dot.ID
	IsDefault    dot.NullBool
}

// +convert:update=ShipmentPriceList
type UpdateShipmentPriceListArgs struct {
	ID          dot.ID
	Name        string
	Description string
}
