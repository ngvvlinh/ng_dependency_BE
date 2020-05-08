package partnercarrier

import (
	"context"

	"o.o/api/top/external/types"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/partnercarrier

// +apix:path=/carrier.Misc
type MiscService interface {
	GetLocationList(context.Context, *cm.Empty) (*types.LocationResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*types.Partner, error)
}

// +apix:path=/carrier.ShipmentConnection
type ShipmentConnectionService interface {
	GetConnections(context.Context, *cm.Empty) (*GetConnectionsResponse, error)
	CreateConnection(context.Context, *CreateConnectionRequest) (*ShipmentConnection, error)
	UpdateConnection(context.Context, *UpdateConnectionRequest) (*ShipmentConnection, error)
	DeleteConnection(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/carrier.Shipment
type ShipmentService interface {
	UpdateFulfillment(context.Context, *UpdateFulfillmentRequest) (*cm.UpdatedResponse, error)
}
