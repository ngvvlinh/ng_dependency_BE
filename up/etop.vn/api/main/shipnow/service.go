package shipnow

import (
	"context"

	shipnowv1 "etop.vn/api/main/shipnow/v1"
	"etop.vn/api/meta"
)

type AggregateBus struct{ meta.Bus }
type QueryBus struct{ meta.Bus }

type Aggregate interface {
	CreateShipnowFulfillment(ctx context.Context, cmd *CreateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	ConfirmShipnowFulfillment(ctx context.Context, cmd *ConfirmShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	CancelShipnowFulfillment(ctx context.Context, cmd *CancelShipnowFulfillmentArgs) (*meta.Empty, error)

	UpdateShipnowFulfillment(ctx context.Context, cmd *UpdateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)
}

type QueryService interface {
	GetShipnowFulfillment(context.Context, *GetShipnowFulfillmentQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
	GetShipnowFulfillments(context.Context, *GetShipnowFulfillmentsQueryArgs) (*GetShipnowFulfillmentsQueryResult, error)
}

//-- Commands --//

type CreateShipnowFulfillmentArgs = shipnowv1.CreateShipnowFulfillmentCommand
type ConfirmShipnowFulfillmentArgs = shipnowv1.ConfirmShipnowFulfillmentCommand
type CancelShipnowFulfillmentArgs = shipnowv1.CancelShipnowFulfillmentCommand
type UpdateShipnowFulfillmentArgs = shipnowv1.UpdateShipnowFulfillmentCommand

type CommitCreateShipnowFulfillmentArgs struct {
	ID int64
}

type RejectCreateShipnowFulfillmentArgs struct {
	ID int64
}

type GetActiveShipnowFulfillmentsArgs struct {
	OrderID int64
}

//-- queries --//

type GetShipnowFulfillmentQueryArgs = shipnowv1.GetShipnowFulfillmentQueryArgs
type GetShipnowFulfillmentQueryResult = shipnowv1.GetShipnowFulfillmentQueryResult
type GetShipnowFulfillmentsQueryArgs = shipnowv1.GetShipnowFulfillmentsQueryArgs
type GetShipnowFulfillmentsQueryResult = shipnowv1.GetShipnowFulfillmentsQueryResult
