package haravan

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	SendUpdateExternalFulfillmentState(context.Context, *SendUpdateExternalFulfillmentStateArgs) (*meta.Empty, error)

	SendUpdateExternalPaymentStatus(context.Context, *SendUpdateExternalPaymentStatusArgs) (*meta.Empty, error)
}

// Call to update fulfillment's Haravan
type SendUpdateExternalFulfillmentStateArgs struct {
	FulfillmentID dot.ID
}

type SendUpdateExternalPaymentStatusArgs struct {
	FulfillmentID dot.ID
}
