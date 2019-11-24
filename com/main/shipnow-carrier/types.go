package shipnow_carrier

import (
	"context"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	"etop.vn/capi/dot"
)

type ShipnowCarrier interface {
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
