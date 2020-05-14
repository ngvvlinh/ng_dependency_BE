package model

import (
	"encoding/json"
	"time"

	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	addressmodel "o.o/backend/com/main/address/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	etlordermodel "o.o/backend/zexp/etl/main/order/model"
	"o.o/capi/dot"
)

// +sqlgen
type Fulfillment struct {
	ID        dot.ID `paging:"id"`
	OrderID   dot.ID
	ShopID    dot.ID
	PartnerID dot.ID

	ShopConfirm   status3.Status `sql_type:"int2"`
	ConfirmStatus status3.Status `sql_type:"int2"`

	TotalItems        int
	TotalWeight       int
	BasketValue       int
	TotalDiscount     int
	TotalAmount       int
	TotalCODAmount    int
	OriginalCODAmount int

	// If a fulfillment is compensated, this field contains the actual amount
	// that carrier payback to shop.
	ActualCompensationAmount int

	ShippingFeeCustomer      int // shop charges customer/shop
	ShippingFeeShop          int // etop charges shop, actual_shipping_service_fee
	ShippingFeeShopLines     []*shippingsharemodel.ShippingFeeLine
	ShippingServiceFee       int // copy from order
	ExternalShippingFee      int // provider charges eTop
	ProviderShippingFeeLines []*shippingsharemodel.ShippingFeeLine
	EtopDiscount             int
	EtopFeeAdjustment        int // eTop điều chỉnh phi (phần thêm)

	ShippingFeeMain       int
	ShippingFeeReturn     int
	ShippingFeeInsurance  int
	ShippingFeeAdjustment int
	ShippingFeeCODS       int
	ShippingFeeInfoChange int
	ShippingFeeOther      int

	// EtopAdjustedShippingFeeMain: eTop điều chỉnh cước phí chính
	EtopAdjustedShippingFeeMain int
	// EtopPriceRule: true khi áp dụng bảng giá eTop, với giá `EtopAdjustedShippingFeeMain`
	EtopPriceRule bool

	VariantIDs []dot.ID
	Lines      []*etlordermodel.OrderLine

	AddressFrom   *addressmodel.Address
	AddressTo     *addressmodel.Address
	AddressReturn *addressmodel.Address

	AddressToProvinceCode string
	AddressToDistrictCode string
	AddressToWardCode     string

	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	ClosedAt                    time.Time
	ExpectedDeliveryAt          time.Time
	ExpectedPickAt              time.Time
	CODEtopTransferedAt         time.Time
	ShippingFeeShopTransferedAt time.Time
	ShippingCancelledAt         time.Time
	ShippingDeliveredAt         time.Time
	ShippingReturnedAt          time.Time
	ShippingCreatedAt           time.Time
	ShippingPickingAt           time.Time
	ShippingHoldingAt           time.Time
	ShippingDeliveringAt        time.Time
	ShippingReturningAt         time.Time

	MoneyTransactionID                 dot.ID
	MoneyTransactionShippingExternalID dot.ID

	CancelReason string

	// CreatedBy   dot.ID
	// UpdatedBy   dot.ID
	// CancelledBy dot.ID

	ShippingProvider    shipping_provider.ShippingProvider `sql_type:"enum(shipping_provider)"`
	ProviderServiceID   string
	ShippingCode        string
	ShippingNote        string
	TryOn               try_on.TryOnCode `sql_type:"enum(try_on)"`
	IncludeInsurance    bool
	ConnectionID        dot.ID
	ConnectionMethod    connection_type.ConnectionMethod `sql_type:"text"`
	ShippingServiceName string

	ExternalShippingName        string
	ExternalShippingID          string
	ExternalShippingCode        string // it's shipping_code
	ExternalShippingCreatedAt   time.Time
	ExternalShippingUpdatedAt   time.Time
	ExternalShippingCancelledAt time.Time
	ExternalShippingDeliveredAt time.Time
	ExternalShippingReturnedAt  time.Time
	ExternalShippingClosedAt    time.Time
	ExternalShippingState       string
	ExternalShippingStateCode   string
	ExternalShippingStatus      status5.Status `sql_type:"int2"`
	ExternalShippingNote        string
	ExternalShippingSubState    string

	ExternalShippingData json.RawMessage

	ShippingState     shipping.State
	ShippingStatus    status5.Status `sql_type:"int2"`
	EtopPaymentStatus status4.Status `sql_type:"int2"`

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status status5.Status `sql_type:"int2"`

	SyncStatus status4.Status `sql_type:"int2"` // -1:error, 0:new, 1:created, 2:pending
	SyncStates *shippingsharemodel.FulfillmentSyncStates

	// Updated by webhook or querying GHN API
	LastSyncAt time.Time

	ExternalShippingLogs []*shipmodel.ExternalShippingLog
	AdminNote            string
	IsPartialDelivery    bool
	CreatedBy            dot.ID

	GrossWeight      int
	ChargeableWeight int
	Length           int
	Width            int
	Height           int

	DeliveryRoute string

	Rid dot.ID
}
