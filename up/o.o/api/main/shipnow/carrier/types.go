package carrier

import (
	"context"
	"time"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/capi/dot"
)

type Manager interface {
	CreateExternalShipnow(ctx context.Context, cmd *CreateExternalShipnowCommand) (*ExternalShipnow, error)
	CancelExternalShipping(ctx context.Context, cmd *CancelExternalShipnowCommand) error
	GetExternalShipnowServices(ctx context.Context, cmd *GetExternalShipnowServicesCommand) ([]*shipnowtypes.ShipnowService, error)

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
	Carrier              types.ShipnowCarrier
	ConnectionID         dot.ID
}

type ExternalShipnow struct {
	ID         string
	UserID     string
	Duration   int
	Distance   float32
	State      shipnow_state.State
	TotalFee   int
	FeeLines   []*shippingtypes.ShippingFeeLine
	CreatedAt  time.Time
	SharedLink string
	Service    *shipnowtypes.ShipnowService
}

type GetExternalShipnowServicesCommand struct {
	ShopID         dot.ID
	PickupAddress  *ordertypes.Address
	DeliveryPoints []*shipnow.DeliveryPoint
	ConnectionIDs  []dot.ID
}

type RegisterExternalAccountCommand struct {
	Phone   string
	Name    string
	Address string
	Carrier types.ShipnowCarrier
}

type GetExternalServiceNameCommand struct {
	Code    string
	Carrier types.ShipnowCarrier
}

type RegisterExternalAccountResult struct {
	Token string
}

type GetExternalAccountCommand struct {
	OwnerID dot.ID
	Carrier types.ShipnowCarrier
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
	Carrier types.ShipnowCarrier
}

type VerifyExternalAccountResult struct {
	TicketID    string `json:"ticket_id"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
