package types

import (
	"context"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow/carrier"
	shipnowcarriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	"o.o/capi/dot"
)

type Config struct {
	ThirdPartyHost string
	Paths          []ConfigPath
}

type ConfigPath struct {
	Carrier              shipnowcarriertypes.ShipnowCarrier
	PathUserVerification string
}

type ConfigPaths []ConfigPath

func (cp ConfigPaths) GetByCarrier(carrier shipnowcarriertypes.ShipnowCarrier) (string, bool) {
	for _, c := range cp {
		if c.Carrier == carrier {
			return c.PathUserVerification, true
		}
	}
	return "", false
}

type Driver interface {
	GetShipnowDriver(
		env string, locationQS location.QueryBus,
		connection *connectioning.Connection,
		shopConnection *connectioning.ShopConnection,
		identityQS identity.QueryBus,
		accountshipnowQS accountshipnow.QueryBus,
		pathCfgs Config,
	) (ShipnowCarrier, error)

	GetAffiliateShipnowDriver(
		env string, locationQS location.QueryBus,
		connection *connectioning.Connection,
		identityQS identity.QueryBus,
		accountshipnowQS accountshipnow.QueryBus,
	) (ShipnowCarrier, error)
}

type ShipnowCarrier interface {
	Code() shipnowcarriertypes.ShipnowCarrier
	CreateExternalShipnow(context.Context, *carrier.CreateExternalShipnowCommand, *shipnowtypes.ShipnowService) (*carrier.ExternalShipnow, error)
	CancelExternalShipnow(context.Context, *carrier.CancelExternalShipnowCommand) error
	GetShipnowServices(context.Context, GetShipnowServiceArgs) ([]*shipnowtypes.ShipnowService, error)
	GetServiceName(code string) (serviceName string, ok bool)
	ParseServiceCode(code string) (serviceID string, err error)

	RegisterExternalAccount(context.Context, *RegisterExternalAccountArgs) (*carrier.RegisterExternalAccountResult, error)
	GetExternalAccount(context.Context, *GetExternalAccountArgs) (*carrier.ExternalAccount, error)
	VerifyExternalAccount(context.Context, *VerifyExternalAccountArgs) (*carrier.VerifyExternalAccountResult, error)
}

type GetShipnowServiceArgs struct {
	ArbitraryID   dot.ID
	ShopID        dot.ID
	OwnerID       dot.ID
	PickupAddress *ordertypes.Address

	DeliveryPoints []*shipnowtypes.DeliveryPoint
	Coupon         string
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
