package carrier

import (
	"context"
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	carrierv1 "etop.vn/api/main/shipnow/carrier/v1"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
)

type Carrier = carrierv1.Carrier

const (
	Ahamove = carrierv1.Carrier_ahamove
)

func CarrierToString(s Carrier) string {
	if s == 0 {
		return ""
	}
	return s.String()
}

func CarrierFromString(s string) Carrier {
	st := carrierv1.Carrier_value[s]
	return Carrier(st)
}

type Manager interface {
	CreateExternalShipping(ctx context.Context, cmd *CreateExternalShipnowCommand) (*ExternalShipnow, error)
	CancelExternalShipping(ctx context.Context, cmd *CancelExternalShipnowCommand) error
	GetExternalShippingServices(ctx context.Context, cmd *GetExternalShipnowServicesCommand) ([]*shipnowtypes.ShipnowService, error)
	GetExternalServiceName(ctx context.Context, cmd *GetExternalServiceNameCommand) (string, error)

	RegisterExternalAccount(ctx context.Context, cmd *RegisterExternalAccountCommand) (*RegisterExternalAccountResult, error)
	GetExternalAccount(ctx context.Context, cmd *GetExternalAccountCommand) (*ExternalAccount, error)
	VerifyExternalAccount(ctx context.Context, cmd *VerifyExternalAccountCommand) (*VerifyExternalAccountResult, error)
}

type CreateExternalShipnowCommand struct {
	ShopID               int64
	ShipnowFulfillmentID int64
	PickupAddress        *ordertypes.Address
	DeliveryPoints       []*shipnow.DeliveryPoint
	ShippingNote         string
}

type CancelExternalShipnowCommand struct {
	ShopID               int64
	ShipnowFulfillmentID int64
	ExternalShipnowID    string
	CarrierServiceCode   string
	CancelReason         string
	Carrier              Carrier
}

type ExternalShipnow struct {
	ID         string
	UserID     string
	Duration   int
	Distance   float32
	State      shipnowtypes.State
	TotalFee   int
	FeeLines   []*shippingtypes.FeeLine
	CreatedAt  time.Time
	SharedLink string
	Service    *shipnowtypes.ShipnowService
}

type GetExternalShipnowServicesCommand struct {
	ShopID         int64
	PickupAddress  *ordertypes.Address
	DeliveryPoints []*shipnow.DeliveryPoint
}

type RegisterExternalAccountCommand struct {
	Phone   string
	Name    string
	Address string
	Carrier Carrier
}

type GetExternalServiceNameCommand struct {
	Code    string
	Carrier Carrier
}

type RegisterExternalAccountResult struct {
	Token string
}

type GetExternalAccountCommand struct {
	OwnerID int64
	Carrier Carrier
}

type ExternalAccount struct {
	ID        string
	Name      string
	Email     string
	Verified  bool
	CreatedAt time.Time
}

type VerifyExternalAccountCommand struct {
	OwnerID int64
	Carrier Carrier
}

type VerifyExternalAccountResult struct {
	TicketID    string
	Subject     string
	Description string
	CreatedAt   string
}
