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

type Aggregate interface {
	CreateShipnowFulfillment(ctx context.Context, cmd *CreateShipnowFulfillmentCommand) (*ShipnowFulfillment, error)

	ConfirmShipnowFulfillment(ctx context.Context, cmd *ConfirmShipnowFulfillmentCommand) (*ShipnowFulfillment, error)

	CancelShipnowFulfillment(ctx context.Context, cmd *CancelShipnowFulfillmentCommand) (*meta.Empty, error)

	GetActiveShipnowFulfillments(ctx context.Context, cmd *GetActiveShipnowFulfillmentsCommand) ([]*ShipnowFulfillment, error)

	UpdateShipnowFulfillment(ctx context.Context, cmd *UpdateShipnowFulfillmentCommand) (*ShipnowFulfillment, error)
}

type QueryService interface {
	GetShipnowFulfillment(context.Context, *GetShipnowFulfillmentQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
	GetShipnowFulfillments(context.Context, *GetShipnowFulfillmentsQueryArgs) (*GetShipnowFulfillmentsQueryResult, error)
}

//-- Commands --//

type CreateShipnowFulfillmentCommand = shipnowv1.CreateShipnowFulfillmentCommand
type ConfirmShipnowFulfillmentCommand = shipnowv1.ConfirmShipnowFulfillmentCommand
type CancelShipnowFulfillmentCommand = shipnowv1.CancelShipnowFulfillmentCommand
type UpdateShipnowFulfillmentCommand = shipnowv1.UpdateShipnowFulfillmentCommand

type CommitCreateShipnowFulfillmentCommand struct {
	ID int64
}

type RejectCreateShipnowFulfillmentCommand struct {
	ID int64
}

type GetActiveShipnowFulfillmentsCommand struct {
	OrderID int64
}

//-- queries --//

type GetShipnowFulfillmentQueryArgs = shipnowv1.GetShipnowFulfillmentQueryArgs
type GetShipnowFulfillmentQueryResult = shipnowv1.GetShipnowFulfillmentQueryResult
type GetShipnowFulfillmentsQueryArgs = shipnowv1.GetShipnowFulfillmentsQueryArgs
type GetShipnowFulfillmentsQueryResult = shipnowv1.GetShipnowFulfillmentsQueryResult
