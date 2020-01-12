package shipping

import (
	"context"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipping/types"
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateFulfillments(context.Context, *CreateFulfillmentsArgs) (fulfillmentID []dot.ID, _ error)

	UpdateFulfillmentShippingState(context.Context, *UpdateFulfillmentShippingStateArgs) (updated int, err error)

	UpdateFulfillmentShippingFees(context.Context, *UpdateFulfillmentShippingFeesArgs) (updated int, err error)

	UpdateFulfillmentsMoneyTxShippingExternalID(context.Context, *UpdateFulfillmentsMoneyTxShippingExternalIDArgs) (updated int, _ error)

	UpdateFulfillmentsStatus(context.Context, *UpdateFulfillmentsStatusArgs) error

	CancelFulfillment(context.Context, *CancelFulfillmentArgs) error
}

type QueryService interface {
	GetFulfillmentByIDOrShippingCode(ctx context.Context, ID dot.ID, ShippingCode string) (*Fulfillment, error)
}

//-- Commands --//

type CreateFulfillmentsArgs struct {
	ShopID dot.ID

	OrderID dot.ID

	PickupAddress *ordertypes.Address

	ShippingAddress *ordertypes.Address

	ReturnAddress *ordertypes.Address

	ShippingType ordertypes.ShippingType

	ShippingServiceCode string

	ShippingServiceFee int

	ShippingServiceName string

	types.WeightInfo

	types.ValueInfo

	TryOn try_on.TryOnCode

	ShippingNote string

	ConnectionID dot.ID

	ShopCarrierID dot.ID
}

type ConfirmFulfillmentArgs struct {
	FulfillmentID dot.ID
}

type CancelFulfillmentArgs struct {
	FulfillmentID dot.ID

	CancelReason string
}

//-- Queries --//

type GetFulfillmentByIDQueryArgs struct {
	FulfillmentID dot.ID
}

type UpdateFulfillmentShippingStateArgs struct {
	FulfillmentID            dot.ID
	ShippingCode             string
	ShippingState            shipping.State
	ActualCompensationAmount dot.NullInt
}

type UpdateFulfillmentShippingFeesArgs struct {
	FulfillmentID    dot.ID
	ShippingCode     string
	ShippingFeeLines []*ShippingFeeLine
}

type UpdateFulfillmentsMoneyTxShippingExternalIDArgs struct {
	FulfillmentIDs            []dot.ID
	MoneyTxShippingExternalID dot.ID
}

type UpdateFulfillmentsStatusArgs struct {
	FulfillmentIDs []dot.ID
	Status         status4.NullStatus
	ShopConfirm    status3.NullStatus
	SyncStatus     status4.NullStatus
}
