package shipmentservice

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShipmentService(context.Context, *CreateShipmentServiceArgs) (*ShipmentService, error)

	UpdateShipmentService(context.Context, *UpdateShipmentServiceArgs) error

	DeleteShipmentService(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetShipmentService(ctx context.Context, ID dot.ID) (*ShipmentService, error)

	GetShipmentServiceByServiceID(ctx context.Context, ServiceID string, ConnID dot.ID) (*ShipmentService, error)

	ListShipmentServices(context.Context, *meta.Empty) ([]*ShipmentService, error)
}

// +convert:create=ShipmentService
type CreateShipmentServiceArgs struct {
	ConnectionID dot.ID
	Name         string
	EdCode       string
	ServiceIDs   []string
	Description  string
	ImageURL     string
}

// +convert:update=ShipmentService
type UpdateShipmentServiceArgs struct {
	ID           dot.ID
	ConnectionID dot.ID
	Name         string
	EdCode       string
	ServiceIDs   []string
	Description  string
	ImageURL     string
	Status       status3.NullStatus
}
