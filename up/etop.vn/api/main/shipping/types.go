package shipping

import (
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/shipping_fee_type"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi/dot"
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
	ProviderShippingFeeLines []*ShippingFeeLine
	ShippingFeeShopLines     []*ShippingFeeLine

	TotalItems               int
	TotalWeight              int
	TotalDiscount            int
	TotalAmount              int
	TotalCODAmount           int
	ActualCompensationAmount int
	EtopDiscount             int

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
	IncludeInsurance    bool

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

type FulfillmentCreatingEvent struct {
	meta.EventMeta
	ShopID      dot.ID
	ShippingFee int
}

type FulfillmentUpdatingEvent struct {
	meta.EventMeta
	FulfillmentID dot.ID
}

type FulfillmentShippingFeeChangedEvent struct {
	meta.EventMeta
	FulfillmentID dot.ID
}
