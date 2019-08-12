package shipnow

import (
	"context"

	shipnowv1 "etop.vn/api/main/shipnow/v1"
	"etop.vn/api/meta"
)

type Aggregate interface {
	CreateShipnowFulfillment(context.Context, *CreateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	ConfirmShipnowFulfillment(context.Context, *ConfirmShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	CancelShipnowFulfillment(context.Context, *CancelShipnowFulfillmentArgs) (*meta.Empty, error)

	UpdateShipnowFulfillment(context.Context, *UpdateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	UpdateShipnowFulfillmentCarrierInfo(context.Context, *UpdateShipnowFulfillmentCarrierInfoArgs) (*ShipnowFulfillment, error)

	UpdateShipnowFulfillmentState(context.Context, *UpdateShipnowFulfillmentStateArgs) (*ShipnowFulfillment, error)

	GetShipnowServices(context.Context, *GetShipnowServicesArgs) (*GetShipnowServicesResult, error)
}

type QueryService interface {
	GetShipnowFulfillment(context.Context, *GetShipnowFulfillmentQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
	GetShipnowFulfillments(context.Context, *GetShipnowFulfillmentsQueryArgs) (*GetShipnowFulfillmentsQueryResult, error)
	GetShipnowFulfillmentByShippingCode(context.Context, *GetShipnowFulfillmentByShippingCodeQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
}

//-- commands --//

type CreateShipnowFulfillmentArgs = shipnowv1.CreateShipnowFulfillmentCommand
type ConfirmShipnowFulfillmentArgs = shipnowv1.ConfirmShipnowFulfillmentCommand
type CancelShipnowFulfillmentArgs = shipnowv1.CancelShipnowFulfillmentCommand
type UpdateShipnowFulfillmentArgs = shipnowv1.UpdateShipnowFulfillmentCommand
type UpdateShipnowFulfillmentCarrierInfoArgs = shipnowv1.UpdateShipnowFulfillmentCarrierInfoCommand
type UpdateShipnowFulfillmentStateArgs = shipnowv1.UpdateShipnowFulfullmentStateCommand
type GetShipnowServicesArgs = shipnowv1.GetShipnowServicesCommand
type GetShipnowServicesResult = shipnowv1.GetShipnowServicesCommandResult

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
type GetShipnowFulfillmentByShippingCodeQueryArgs = shipnowv1.GetShipnowFulfillmentByShippingCodeQueryArgs
