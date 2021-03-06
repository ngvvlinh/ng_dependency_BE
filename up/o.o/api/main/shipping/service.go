package shipping

import (
	"context"
	"encoding/json"
	"time"

	"o.o/api/main/connectioning"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping/types"
	shippingtypes "o.o/api/main/shipping/types"
	shippingstate "o.o/api/top/types/etc/shipping"
	shippingsubstate "o.o/api/top/types/etc/shipping/substate"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_payment_type"
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

	CreateFulfillmentsFromImport(context.Context, *CreateFulfillmentsFromImportArgs) ([]*CreateFullfillmentsFromImportResult, error)

	CreatePartialFulfillment(context.Context, *CreatePartialFulfillmentArgs) (fulfillmentID dot.ID, _ error)

	UpdateFulfillmentShippingState(context.Context, *UpdateFulfillmentShippingStateArgs) (updated int, _ error)

	UpdateFulfillmentShippingSubstate(context.Context, *UpdateFulfillmentShippingSubstateArgs) (updated int, _ error)

	UpdateFulfillmentShippingFees(context.Context, *UpdateFulfillmentShippingFeesArgs) (updated int, err error)

	UpdateFulfillmentCODAmount(context.Context, *UpdateFulfillmentCODAmountArgs) error

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

	UpdateFulfillmentShippingCode(context.Context, *UpdateFulfillmentShippingCodeArgs) error
}

type QueryService interface {
	GetFulfillmentByIDOrShippingCode(context.Context, *GetFulfillmentByIDOrShippingCodeArgs) (*Fulfillment, error)

	ListCustomerReturnRates(context.Context, *ListCustomerReturnRatesArgs) ([]*CustomerReturnRateExtended, error)

	ListFulfillmentsByIDs(ctx context.Context, IDs []dot.ID, shopID dot.ID) ([]*Fulfillment, error)

	ListFulfillmentsByMoneyTx(context.Context, *ListFullfillmentsByMoneyTxArgs) ([]*Fulfillment, error)

	/*
		ListFulfillmentsForMoneyTx: L???c t???t c??? ffms th??a ??i???u ki???n ????? th??m v??o phi??n chuy???n ti???n shop
	*/
	ListFulfillmentsForMoneyTx(context.Context, *ListFulfillmentForMoneyTxArgs) ([]*Fulfillment, error)

	GetFulfillmentExtended(ctx context.Context, ID dot.ID, ShippingCode string) (*FulfillmentExtended, error)

	ListFulfillmentExtendedsByIDs(ctx context.Context, IDs []dot.ID, ShopID dot.ID) ([]*FulfillmentExtended, error)

	ListFulfillmentExtendedsByMoneyTxShippingID(ctx context.Context, shopID dot.ID, moneyTxShippingID dot.ID) ([]*FulfillmentExtended, error)

	ListFulfillmentsByShippingCodes(ctx context.Context, Codes []string) ([]*Fulfillment, error)
}

//-- Commands --//

type CreatePartialFulfillmentArgs struct {
	FulfillmentID dot.ID
	ShopID        dot.ID

	InfoChanges *InfoChanges
}

type InfoChanges struct {
	ShippingCode dot.NullString
	Weight       dot.NullInt
	Length       dot.NullInt
	Height       dot.NullInt
	Width        dot.NullInt
}

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

	ShippingPaymentType shipping_payment_type.NullShippingPaymentType

	ShippingNote string

	ConnectionID dot.ID

	ShopCarrierID dot.ID

	Coupon string
}

type CreateFulfillmentsFromImportArgs struct {
	Fulfillments []*CreateFulfillmentFromImportArgs
}

type CreateFulfillmentFromImportArgs struct {
	ID dot.ID

	ShopID dot.ID

	ConnectionID dot.ID

	ShippingServiceCode string

	ShippingServiceFee int

	ShippingServiceName string

	EdCode string

	PickupAddress *ordertypes.Address

	ShippingAddress *ordertypes.Address

	TotalWeight int

	BasketValue int

	CODAmount int

	ProductDescription string

	IncludeInsurance bool

	TryOn try_on.TryOnCode

	ShippingNote string

	CreatedBy dot.ID
}

type CreateFullfillmentsFromImportResult struct {
	FulfillmentID dot.ID
	Error         error
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
	AdminNote                string
}

type UpdateFulfillmentShippingFeesArgs struct {
	FulfillmentID            dot.ID
	ShippingCode             string
	ProviderShippingFeeLines []*shippingtypes.ShippingFeeLine
	ShippingFeeLines         []*shippingtypes.ShippingFeeLine
	// @deprecated TotalCODAmount
	TotalCODAmount    dot.NullInt
	UpdatedBy         dot.ID
	AdminNote         string
	ShipmentPriceInfo *ShipmentPriceInfo
}

type UpdateFulfillmentCODAmountArgs struct {
	FulfillmentID     dot.ID
	ShippingCode      string
	TotalCODAmount    dot.NullInt
	IsPartialDelivery dot.NullBool
	AdminNote         string
	UpdatedBy         dot.ID
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
	FulfillmentID       dot.ID
	AddressTo           *ordertypes.Address
	AddressFrom         *ordertypes.Address
	IncludeInsurance    dot.NullBool
	InsuranceValue      dot.NullInt
	GrossWeight         dot.NullInt
	TryOn               try_on.TryOnCode
	ShippingPaymentType shipping_payment_type.NullShippingPaymentType
	ShippingNote        dot.NullString
}

type UpdateFulfillmentsCODTransferedAtArgs struct {
	FulfillmentIDs     []dot.ID
	MoneyTxShippingIDs []dot.ID
	CODTransferedAt    time.Time
}

type UpdateFfmExternalShippingInfoArgs struct {
	FulfillmentID             dot.ID
	ShippingState             shippingstate.State
	ShippingSubstate          shippingsubstate.NullSubstate
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
	ShippingHoldingAt         time.Time
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
	ProviderFeeLines []*shippingtypes.ShippingFeeLine
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

type ListCustomerReturnRatesArgs struct {
	ConnectionIDs []dot.ID
	ShopID        dot.ID
	Phone         string
}

type CustomerReturnRateExtended struct {
	Connection         *connectioning.Connection
	CustomerReturnRate *CustomerReturnRate
}

type CustomerReturnRate struct {
	Level     string  `json:"level"`
	LevelCode string  `json:"level_code"`
	Rate      float64 `json:"rate"`
}

type UpdateFulfillmentShippingSubstateArgs struct {
	FulfillmentID    dot.ID
	ShippingSubstate shippingsubstate.Substate
}

type UpdateFulfillmentShippingCodeArgs struct {
	FulfillmentID dot.ID
	ShippingCode  string
}
