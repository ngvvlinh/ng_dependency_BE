package carrier

import (
	"context"
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier/types"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/top/types/etc/shipnow_state"
	"etop.vn/capi/dot"
)

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
	ShopID               dot.ID
	ShipnowFulfillmentID dot.ID
	PickupAddress        *ordertypes.Address
	DeliveryPoints       []*shipnow.DeliveryPoint
	ShippingNote         string
}

type CancelExternalShipnowCommand struct {
	ShopID               dot.ID
	ShipnowFulfillmentID dot.ID
	ExternalShipnowID    string
	CarrierServiceCode   string
	CancelReason         string
	Carrier              types.Carrier
}

type ExternalShipnow struct {
	ID         string
	UserID     string
	Duration   int
	Distance   float32
	State      shipnow_state.State
	TotalFee   int
	FeeLines   []*shippingtypes.FeeLine
	CreatedAt  time.Time
	SharedLink string
	Service    *shipnowtypes.ShipnowService
}

type GetExternalShipnowServicesCommand struct {
	ShopID         dot.ID
	PickupAddress  *ordertypes.Address
	DeliveryPoints []*shipnow.DeliveryPoint
}

type RegisterExternalAccountCommand struct {
	Phone   string
	Name    string
	Address string
	Carrier types.Carrier
}

type GetExternalServiceNameCommand struct {
	Code    string
	Carrier types.Carrier
}

type RegisterExternalAccountResult struct {
	Token string
}

type GetExternalAccountCommand struct {
	OwnerID dot.ID
	Carrier types.Carrier
}

type ExternalAccount struct {
	ID        string
	Name      string
	Email     string
	Verified  bool
	CreatedAt time.Time
}

type VerifyExternalAccountCommand struct {
	OwnerID dot.ID
	Carrier types.Carrier
}

type VerifyExternalAccountResult struct {
	TicketID    string
	Subject     string
	Description string
	CreatedAt   string
}
