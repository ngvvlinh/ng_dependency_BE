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
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +gen:event:topic=event/shipping

var ShippingFeeShopTypes = []shipping_fee_type.ShippingFeeType{
	shipping_fee_type.Main,
	shipping_fee_type.Return,
	shipping_fee_type.Adjustment,
	shipping_fee_type.AddressChange,
	shipping_fee_type.Cods,
	shipping_fee_type.Insurance,
	shipping_fee_type.Other,
	shipping_fee_type.Discount,
}

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
	ProviderShippingFeeLines []*ShippingFeeLine
	ShippingFeeShopLines     []*ShippingFeeLine

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
	ExternalShippingNote     string
	ExternalShippingSubState string
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
	EtopPriceRule       bool
	ShipmentPriceListID dot.ID
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

	ShippingFeeLine []*ShippingFeeLine

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

type ShippingFeeLine struct {
	ShippingFeeType shipping_fee_type.ShippingFeeType

	Cost int

	ExternalServiceID string

	ExternalServiceName string

	ExternalServiceType string
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

func GetTotalShippingFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		result += item.Cost
	}
	return result
}

func GetShippingFeeShopLines(items []*ShippingFeeLine, etopPriceRule bool, mainFee dot.NullInt) []*ShippingFeeLine {
	res := make([]*ShippingFeeLine, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		line := GetShippingFeeShopLine(*item, etopPriceRule, mainFee)
		if line != nil {
			res = append(res, line)
		}
	}
	return res
}

func GetShippingFeeShopLine(item ShippingFeeLine, etopPriceRule bool, mainFee dot.NullInt) *ShippingFeeLine {
	if item.ShippingFeeType == shipping_fee_type.Main && etopPriceRule {
		item.Cost = mainFee.Apply(item.Cost)
	}
	if contains(ShippingFeeShopTypes, item.ShippingFeeType) {
		return &item
	}
	return nil
}

func ApplyShippingFeeLine(lines []*ShippingFeeLine, item *ShippingFeeLine) []*ShippingFeeLine {
	if item == nil {
		return lines
	}
	for _, line := range lines {
		if line.ShippingFeeType == item.ShippingFeeType {
			line.Cost = item.Cost
			return lines
		}
	}
	lines = append(lines, item)
	return lines
}

func GetShippingFeeLine(lines []*ShippingFeeLine, _type shipping_fee_type.ShippingFeeType) *ShippingFeeLine {
	for _, line := range lines {
		if line.ShippingFeeType == _type {
			return line
		}
	}
	return nil
}

func GetShippingFee(lines []*ShippingFeeLine, _type shipping_fee_type.ShippingFeeType) int {
	line := GetShippingFeeLine(lines, _type)
	if line == nil {
		return 0
	}
	return line.Cost
}

func UpdateShippingFees(items []*ShippingFeeLine, fee int, shippingFeeType shipping_fee_type.ShippingFeeType) []*ShippingFeeLine {
	if fee == 0 {
		return items
	}
	found := false
	for _, item := range items {
		if item.ShippingFeeType == shippingFeeType {
			item.Cost = fee
			found = true
		}
	}
	if !found {
		items = append(items, &ShippingFeeLine{
			ShippingFeeType: shippingFeeType,
			Cost:            fee,
		})
	}
	return items
}

func contains(lines []shipping_fee_type.ShippingFeeType, feeType shipping_fee_type.ShippingFeeType) bool {
	for _, line := range lines {
		if feeType == line {
			return true
		}
	}
	return false
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
