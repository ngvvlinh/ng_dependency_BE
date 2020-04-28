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

	ListShipmentPriceList(context.Context, *meta.Empty) ([]*ShipmentPriceList, error)
}

// +convert:create=ShipmentPriceList
type CreateShipmentPriceListArg struct {
	Name        string
	Description string
	IsActive    bool
}

// +convert:update=ShipmentPriceList
type UpdateShipmentPriceListArgs struct {
	ID          dot.ID
	Name        string
	Description string
}
