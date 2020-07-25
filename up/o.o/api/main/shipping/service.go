package shipping

import (
	"context"
	"encoding/json"
	"time"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping/types"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
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

	UpdateFulfillmentShippingState(context.Context, *UpdateFulfillmentShippingStateArgs) (updated int, _ error)

	UpdateFulfillmentShippingFees(context.Context, *UpdateFulfillmentShippingFeesArgs) (updated int, err error)

	UpdateFulfillmentsMoneyTxID(context.Context, *UpdateFulfillmentsMoneyTxIDArgs) (updated int, _ error)

	UpdateFulfillmentsCODTransferedAt(context.Context, *UpdateFulfillmentsCODTransferedAtArgs) error

	RemoveFulfillmentsMoneyTxID(context.Context, *RemoveFulfillmentsMoneyTxIDArgs) (updated int, _ error)

	UpdateFulfillmentsStatus(context.Context, *UpdateFulfillmentsStatusArgs) error

	ShopUpdateFulfillmentInfo(context.Context, *UpdateFulfillmentInfoArgs) (updated int, _ error)

	ShopUpdateFulfillmentCOD(context.Context, *ShopUpdateFulfillmentCODArgs) (updated int, err error)

	UpdateFulfillmentInfo(context.Context, *UpdateFulfillmentInfoByAdminArgs) (updated int, _ error)

	CancelFulfillment(context.Context, *CancelFulfillmentArgs) error

	UpdateFulfillmentExternalShippingInfo(context.Context, *UpdateFfmExternalShippingInfoArgs) (updated int, _ error)

	UpdateFulfillmentShippingFeesFromWebhook(context.Context, *UpdateFulfillmentShippingFeesFromWebhookArgs) error

	AddFulfillmentShippingFee(context.Context, *AddFulfillmentShippingFeeArgs) error
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

	ListFulfillmentsByShippingCodes(ctx context.Context, Codes []string) ([]*Fulfillment, error)
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

	Coupon string
}

type ConfirmFulfillmentArgs struct {
	FulfillmentID dot.ID
}

type CancelFulfillmentArgs struct {
	FulfillmentID dot.ID

	CancelReason string
}

type UpdateFulfillmentInfoByAdminArgs struct {
	FulfillmentID dot.ID
	ShippingCode  string
	FullName      dot.NullString
	Phone         dot.NullString
	AdminNote     string
}

//-- Queries --//

type GetFulfillmentByIDQueryArgs struct {
	FulfillmentID dot.ID
}

type UpdateFulfillmentShippingStateArgs struct {
	PartnerID                dot.ID
	FulfillmentID            dot.ID
	ShippingCode             string
	ShippingState            shippingstate.State
	ActualCompensationAmount dot.NullInt
	UpdatedBy                dot.ID
}

type UpdateFulfillmentShippingFeesArgs struct {
	FulfillmentID            dot.ID
	ShippingCode             string
	ProviderShippingFeeLines []*ShippingFeeLine
	ShippingFeeLines         []*ShippingFeeLine
	TotalCODAmount           dot.NullInt
	UpdatedBy                dot.ID
}

type ShopUpdateFulfillmentCODArgs struct {
	FulfillmentID  dot.ID
	ShippingCode   string
	TotalCODAmount dot.NullInt
	UpdatedBy      dot.ID
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

// +convert:update=Fulfillment(ID,ShopID)
type UpdateFulfillmentInfoArgs struct {
	FulfillmentID    dot.ID
	AddressTo        *ordertypes.Address
	AddressFrom      *ordertypes.Address
	IncludeInsurance dot.NullBool
	InsuranceValue   dot.NullInt
	GrossWeight      int
	TryOn            try_on.TryOnCode
	ShippingNote     dot.NullString
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

type AddFulfillmentShippingFeeArgs struct {
	FulfillmentID   dot.ID
	ShippingCode    string
	ShippingFeeType shipping_fee_type.ShippingFeeType
	UpdatedBy       dot.ID
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
