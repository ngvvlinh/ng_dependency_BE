package shipping

import (
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipping"
	shippingsubstate "o.o/api/top/types/etc/shipping/substate"
	"o.o/api/top/types/etc/shipping_payment_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +gen:event:topic=event/shipping

type ShippingService struct {
	Code string

	Name string

	Fee int

	Carrier string

	EstimatedPickupAt time.Time

	EstimatedDeliveryAt time.Time
}

type Fulfillment struct {
	ID        dot.ID
	OrderID   dot.ID
	ShopID    dot.ID
	PartnerID dot.ID

	Lines []*ordertypes.ItemLine

	ShopConfirm       status3.Status
	ConfirmStatus     status3.Status
	Status            status5.Status
	ShippingState     shipping.State
	ShippingStatus    status5.Status
	EtopPaymentStatus status4.Status

	ShippingFeeShop          int
	ProviderShippingFeeLines []*types.ShippingFeeLine
	ShippingFeeShopLines     []*types.ShippingFeeLine

	TotalItems               int
	TotalWeight              int
	TotalDiscount            int
	TotalAmount              int
	TotalCODAmount           int
	ActualCompensationAmount int
	EtopDiscount             int
	EtopFeeAdjustment        int

	types.WeightInfo
	types.ValueInfo

	CODEtopTransferedAt                time.Time
	MoneyTransactionID                 dot.ID
	MoneyTransactionShippingExternalID dot.ID

	ShippingType        ordertypes.ShippingType
	ConnectionID        dot.ID
	ConnectionMethod    connection_type.ConnectionMethod
	ShopCarrierID       dot.ID
	ProviderServiceID   string
	ShippingCode        string
	ShippingServiceName string
	ShippingNote        string
	TryOn               try_on.TryOnCode
	ShippingPaymentType shipping_payment_type.ShippingPaymentType
	IncludeInsurance    dot.NullBool
	InsuranceValue      dot.NullInt

	Coupon string

	CreatedAt           time.Time
	UpdatedAt           time.Time
	ClosedAt            time.Time
	ShippingCancelledAt time.Time
	CancelReason        string

	ShippingProvider            shipping_provider.ShippingProvider
	ExternalShippingName        string
	ExternalShippingFee         int
	ShippingFeeCustomer         int
	ExternalShippingID          string
	ExternalShippingCode        string
	ExternalShippingCreatedAt   time.Time
	ExternalShippingUpdatedAt   time.Time
	ExternalShippingCancelledAt time.Time
	ExternalShippingDeliveredAt time.Time
	ExternalShippingReturnedAt  time.Time

	ExternalShippingState    string
	ExternalShippingStatus   status5.Status
	ExternalShippingNote     dot.NullString
	ExternalShippingSubState dot.NullString
	ExternalShippingLogs     []*ExternalShippingLog
	SyncStatus               status4.Status
	SyncStates               *FulfillmentSyncStates

	ExpectedDeliveryAt          time.Time
	ExpectedPickAt              time.Time
	ShippingFeeShopTransferedAt time.Time

	AddressTo   *ordertypes.Address
	AddressFrom *ordertypes.Address

	// EtopAdjustedShippingFeeMain: eTop điều chỉnh cước phí chính
	EtopAdjustedShippingFeeMain int
	// EtopPriceRule: true khi áp dụng bảng giá eTop, với giá `EtopAdjustedShippingFeeMain`
	EtopPriceRule     bool
	ShipmentPriceInfo *ShipmentPriceInfo

	LinesContent     string
	EdCode           string
	ShippingSubstate shippingsubstate.NullSubstate
}

type FulfillmentSyncStates struct {
	SyncAt            time.Time
	TrySyncAt         time.Time
	Error             *meta.Error
	NextShippingState shipping.State
}

type ExternalShippingLog struct {
	StateText string
	Time      string
	Message   string
}

type ExternalShipmentData struct {
	State string

	ShippingFee int

	// ShippingData

	// ShippingLogs

	ShippingFeeLine []*types.ShippingFeeLine

	CreatedAt time.Time

	UpdatedAt time.Time

	PickingAt time.Time

	PickedAt time.Time

	HoldingAt time.Time

	DeliveringAt time.Time

	DeliveredAt time.Time

	ReturningAt time.Time

	ReturnedAt time.Time
}

type FulfillmentExtended struct {
	*Fulfillment
	Shop  *identity.Shop
	Order *ordering.Order
	// MoneyTxShipping *moneytx.MoneyTransactionShipping
}

type FulfillmentsCreatingEvent struct {
	meta.EventMeta
	ShopID      dot.ID
	OrderID     dot.ID
	ShippingFee int
}

type FulfillmentsCreatedEvent struct {
	meta.EventMeta
	FulfillmentIDs []dot.ID
	ShippingType   ordertypes.ShippingType
	OrderID        dot.ID
}

type FulfillmentUpdatedEvent struct {
	meta.EventMeta
	FulfillmentID     dot.ID
	MoneyTxShippingID dot.ID
}

type FulfillmentUpdatedInfoEvent struct {
	meta.EventMeta
	OrderID  dot.ID
	FullName dot.NullString
	Phone    dot.NullString
}

type SingleFulfillmentCreatingEvent struct {
	meta.EventMeta
	ShopID       dot.ID
	FromAddress  *ordertypes.Address
	ShippingFee  int
	ConnectionID dot.ID
}

type FulfillmentFromImportCreatedEvent struct {
	meta.EventMeta
	ShopID        dot.ID
	FulfillmentID dot.ID
}

type ShipmentPriceInfo struct {
	ShipmentPriceID     dot.ID
	ShipmentPriceListID dot.ID
	OriginFee           int
	MakeupFee           int
}

func CalcShopShippingFee(externalFee int, ffm *Fulfillment) int {
	if ffm == nil {
		return externalFee
	}
	fee := externalFee + ffm.EtopFeeAdjustment - ffm.EtopDiscount
	if fee < 0 {
		return 0
	}
	return fee
}

func GetConnectionID(connectionID dot.ID, carrier shipping_provider.ShippingProvider) dot.ID {
	if connectionID != 0 {
		return connectionID
	}

	// backward-compatible
	switch carrier {
	case shipping_provider.GHN:
		return connectioning.DefaultTopshipGHNConnectionID
	case shipping_provider.GHTK:
		return connectioning.DefaultTopshipGHTKConnectionID
	case shipping_provider.VTPost:
		return connectioning.DefaultTopshipVTPostConnectionID
	default:
		return 0
	}
}

func IsStateReturn(state shipping.State) bool {
	return state == shipping.Returning || state == shipping.Returned
}
