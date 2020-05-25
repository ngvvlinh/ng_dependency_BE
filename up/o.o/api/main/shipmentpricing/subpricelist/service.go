package subpricelist

import (
	"context"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShipmentSubPriceList(context.Context, *CreateSubPriceListArgs) (*ShipmentSubPriceList, error)

	UpdateShipmentSubPriceList(context.Context, *UpdateSubPriceListArgs) error

	DeleteShipmentSubPriceList(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	ListShipmentSubPriceList(context.Context, *ListSubPriceListArgs) ([]*ShipmentSubPriceList, error)

	ListShipmentSubPriceListByIDs(ctx context.Context, IDs []dot.ID) ([]*ShipmentSubPriceList, error)

	GetShipmentSubPriceList(ctx context.Context, ID dot.ID) (*ShipmentSubPriceList, error)
}

// +convert:create=ShipmentSubPriceList
type CreateSubPriceListArgs struct {
	Name         string
	Description  string
	ConnectionID dot.ID
}

// +convert:update=ShipmentSubPriceList
type UpdateSubPriceListArgs struct {
	ID          dot.ID
	Name        string
	Description string
	Status      status3.NullStatus
}

type ListSubPriceListArgs struct {
	ConnectionID dot.ID
	Status       status3.NullStatus
}
