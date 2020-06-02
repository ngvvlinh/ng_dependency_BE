package shipnowcarrier

import (
	"context"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow/carrier"
	shipnowcarriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	"o.o/capi/dot"
)

type ShipnowCarrier interface {
	Code() shipnowcarriertypes.Carrier
	CreateExternalShipnow(context.Context, *carrier.CreateExternalShipnowCommand, *shipnowtypes.ShipnowService) (*carrier.ExternalShipnow, error)
	CancelExternalShipnow(context.Context, *carrier.CancelExternalShipnowCommand) error
	GetShippingServices(context.Context, GetShippingServiceArgs) ([]*shipnowtypes.ShipnowService, error)
	GetServiceName(code string) (serviceName string, ok bool)
	ParseServiceCode(code string) (serviceID string, err error)
}

type ShipnowCarrierAccount interface {
	RegisterExternalAccount(context.Context, *RegisterExternalAccountArgs) (*carrier.RegisterExternalAccountResult, error)
	GetExternalAccount(context.Context, *GetExternalAccountArgs) (*carrier.ExternalAccount, error)
	VerifyExternalAccount(context.Context, *VerifyExternalAccountArgs) (*carrier.VerifyExternalAccountResult, error)
}

type GetShippingServiceArgs struct {
	ShopID        dot.ID
	PickupAddress *ordertypes.Address

	DeliveryPoints []*shipnowtypes.DeliveryPoint
}

type RegisterExternalAccountArgs struct {
	Phone   string
	Name    string
	Address string
}

type GetExternalAccountArgs struct {
	OwnerID dot.ID
}

type VerifyExternalAccountArgs struct {
	OwnerID dot.ID
}
