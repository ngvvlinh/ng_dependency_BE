package shipnow

import (
	"context"

	shipnowv1 "etop.vn/api/main/shipnow/v1"
	"etop.vn/api/meta"
)

type Bus struct {
	meta.Bus
}

type EventBus struct {
	meta.EventBus
}

type ShipnowProcessManager interface {
	HandleConfirmationRequested(ctx context.Context, event ShipnowEvent) error
}

type ShipnowAggregate interface {
	CreateShipnowFulfillment(ctx context.Context, cmd *CreateShipnowFulfillmentCommand) (*ShipnowFulfillment, error)

	ConfirmShipnowFulfillment(ctx context.Context, cmd *ConfirmShipnowFulfillmentCommand) (*meta.Empty, error)

	CancelShipnowFulfillment(ctx context.Context, cmd *CancelShipnowFulfillmentCommand) (*meta.Empty, error)
}

type ShipnowQueryService interface {
	GetShipnowFulfillment(context.Context, *GetShipnowFulfillmentQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
}

//-- Commands --//

type CreateShipnowFulfillmentCommand = shipnowv1.CreateShipnowFulfillmentCommand
type ConfirmShipnowFulfillmentCommand = shipnowv1.ConfirmShipnowFulfillmentCommand
type CancelShipnowFulfillmentCommand = shipnowv1.CancelShipnowFulfillmentCommand

type CommitCreateShipnowFulfillmentCommand struct {
	ID int64
}

type RejectCreateShipnowFulfillmentCommand struct {
	ID int64
}

//-- queries --//

type GetShipnowFulfillmentQueryArgs = shipnowv1.GetShipnowFulfillmentQueryArgs
type GetShipnowFulfillmentQueryResult = shipnowv1.GetShipnowFulfillmentQueryResult
