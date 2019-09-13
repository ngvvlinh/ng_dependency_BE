package haravan

import (
	"context"

	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	SendUpdateExternalFulfillmentState(context.Context, *SendUpdateExternalFulfillmentStateArgs) (*meta.Empty, error)

	SendUpdateExternalPaymentStatus(context.Context, *SendUpdateExternalPaymentStatusArgs) (*meta.Empty, error)
}

// Call to update fulfillment's Haravan
type SendUpdateExternalFulfillmentStateArgs struct {
	FulfillmentID int64
}

type SendUpdateExternalPaymentStatusArgs struct {
	FulfillmentID int64
}
