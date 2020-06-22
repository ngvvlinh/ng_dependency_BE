package shipping

import (
	"context"
	"encoding/json"
	"time"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping/types"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateFulfillments(context.Context, *CreateFulfillmentsArgs) (fulfillmentID []dot.ID, _ error)

	UpdateFulfillmentShippingState(context.Context, *UpdateFulfillmentShippingStateArgs) (updated int, err error)

	UpdateFulfillmentShippingFees(context.Context, *UpdateFulfillmentShippingFeesArgs) (updated int, err error)

	UpdateFulfillmentsMoneyTxID(context.Context, *UpdateFulfillmentsMoneyTxIDArgs) (updated int, _ error)

	UpdateFulfillmentsCODTransferedAt(context.Context, *UpdateFulfillmentsCODTransferedAtArgs) error

	RemoveFulfillmentsMoneyTxID(context.Context, *RemoveFulfillmentsMoneyTxIDArgs) (updated int, _ error)

	UpdateFulfillmentsStatus(context.Context, *UpdateFulfillmentsStatusArgs) error

	CancelFulfillment(context.Context, *CancelFulfillmentArgs) error

	UpdateFulfillmentExternalShippingInfo(context.Context, *UpdateFfmExternalShippingInfoArgs) (updated int, _ error)

	UpdateFulfillmentShippingFeesFromWebhook(context.Context, *UpdateFulfillmentShippingFeesFromWebhookArgs) error
}

type QueryService interface {
	GetFulfillmentByIDOrShippingCode(context.Context, *GetFulfillmentByIDOrShippingCodeArgs) (*Fulfillment, error)

	ListFulfillmentsByIDs(ctx context.Context, IDs []dot.ID, shopID dot.ID) ([]*Fulfillment, error)

	ListFulfillmentsByMoneyTx(context.Context, *ListFullfillmentsByMoneyTxArgs) ([]*Fulfillment, error)

	/*
		ListFulfillmentsForMoneyTx: Lọc tất cả ffms thõa điều kiện để thêm vào phiên chuyển tiền shop
	*/
	ListFulfillmentsForMoneyTx(context.Context, *ListFulfillmentForMoneyTxArgs) ([]*Fulfillment, error)

	GetFulfillmentExtended(ctx context.Context, ID dot.ID, ShippingCode string) (*FulfillmentExtended, error)

	ListFulfillmentExtendedsByIDs(ctx context.Context, IDs []dot.ID, ShopID dot.ID) ([]*FulfillmentExtended, error)

	ListFulfillmentExtendedsByMoneyTxShippingID(ctx context.Context, shopID dot.ID, moneyTxShippingID dot.ID) ([]*FulfillmentExtended, error)

	ListFulfillmentsByShippingCodes(ctx context.Context, IDs []string) ([]*Fulfillment, error)
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
	PartnerID                dot.ID
	FulfillmentID            dot.ID
	ShippingState            shippingstate.State
	ActualCompensationAmount dot.NullInt
}

type UpdateFulfillmentShippingFeesArgs struct {
	FulfillmentID               dot.ID
	ShippingCode                string
	EtopPriceRule               dot.NullBool
	EtopAdjustedShippingFeeMain dot.NullInt
	ProviderShippingFeeLines    []*ShippingFeeLine
	ShippingFeeLines            []*ShippingFeeLine
}

type UpdateFulfillmentsMoneyTxIDArgs struct {
	FulfillmentIDs            []dot.ID
	MoneyTxShippingExternalID dot.ID
	MoneyTxShippingID         dot.ID
}

type RemoveFulfillmentsMoneyTxIDArgs struct {
	FulfillmentIDs            []dot.ID
	MoneyTxShippingID         dot.ID
	MoneyTxShippingExternalID dot.ID
}

type UpdateFulfillmentsStatusArgs struct {
	FulfillmentIDs []dot.ID
	Status         status4.NullStatus
	ShopConfirm    status3.NullStatus
	SyncStatus     status4.NullStatus
}

type UpdateFulfillmentsCODTransferedAtArgs struct {
	FulfillmentIDs     []dot.ID
	MoneyTxShippingIDs []dot.ID
	CODTransferedAt    time.Time
}

type UpdateFfmExternalShippingInfoArgs struct {
	FulfillmentID             dot.ID
	ShippingState             shippingstate.State
	ShippingStatus            status5.Status
	ExternalShippingData      json.RawMessage
	ExternalShippingState     string
	ExternalShippingSubState  dot.NullString
	ExternalShippingStatus    status5.Status
	ExternalShippingNote      dot.NullString
	ExternalShippingUpdatedAt time.Time
	ExternalShippingLogs      []*ExternalShippingLog
	ExternalShippingStateCode string
	Weight                    int
	ClosedAt                  time.Time
	LastSyncAt                time.Time
	ShippingCreatedAt         time.Time
	ShippingPickingAt         time.Time
	ShippingDeliveringAt      time.Time
	ShippingDeliveredAt       time.Time
	ShippingReturningAt       time.Time
	ShippingReturnedAt        time.Time
	ShippingCancelledAt       time.Time
}

type UpdateFulfillmentShippingFeesFromWebhookArgs struct {
	FulfillmentID    dot.ID
	NewWeight        int
	NewState         shippingstate.State
	ProviderFeeLines []*ShippingFeeLine
}

type GetFulfillmentByIDOrShippingCodeArgs struct {
	ID            dot.ID
	ShippingCode  string
	ConnectionIDs []dot.ID
}

type ListFulfillmentForMoneyTxArgs struct {
	ShippingProvider shipping_provider.ShippingProvider
	ShippingStates   []shippingstate.State
	IsNoneCOD        dot.NullBool
}

type ListFullfillmentsByMoneyTxArgs struct {
	MoneyTxShippingIDs        []dot.ID
	MoneyTxShippingExternalID dot.ID
}
