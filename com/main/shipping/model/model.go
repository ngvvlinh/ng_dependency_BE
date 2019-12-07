package model

import (
	"encoding/json"
	"fmt"
	"time"

	"etop.vn/api/top/types/etc/status5"

	"etop.vn/api/top/types/etc/status4"

	"etop.vn/api/top/types/etc/status3"

	ordermodel "etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFulfillment(&Fulfillment{})

type Fulfillment struct {
	ID        dot.ID
	OrderID   dot.ID
	ShopID    dot.ID
	PartnerID dot.ID

	ShopConfirm   status3.Status
	ConfirmStatus status3.Status

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
	ShippingFeeShopLines     []*model.ShippingFeeLine
	ShippingServiceFee       int // copy from order
	ExternalShippingFee      int // provider charges eTop
	ProviderShippingFeeLines []*model.ShippingFeeLine
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
	Lines      ordermodel.OrderLinesList

	TypeFrom      model.FulfillmentEndpoint
	TypeTo        model.FulfillmentEndpoint
	AddressFrom   *model.Address
	AddressTo     *model.Address
	AddressReturn *model.Address

	AddressToProvinceCode string
	AddressToDistrictCode string
	AddressToWardCode     string

	CreatedAt                   time.Time `sq:"create"`
	UpdatedAt                   time.Time `sq:"update"`
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

	ShippingProvider  model.ShippingProvider
	ProviderServiceID string
	ShippingCode      string
	ShippingNote      string
	TryOn             model.TryOn
	IncludeInsurance  bool

	ExternalShippingName        string
	ExternalShippingID          string // it's shipping_service_code
	ExternalShippingCode        string // it's shipping_code
	ExternalShippingCreatedAt   time.Time
	ExternalShippingUpdatedAt   time.Time
	ExternalShippingCancelledAt time.Time
	ExternalShippingDeliveredAt time.Time
	ExternalShippingReturnedAt  time.Time
	ExternalShippingClosedAt    time.Time
	ExternalShippingState       string
	ExternalShippingStateCode   string
	ExternalShippingStatus      status5.Status
	ExternalShippingNote        string
	ExternalShippingSubState    string

	ExternalShippingData json.RawMessage

	ShippingState     model.ShippingState
	ShippingStatus    status5.Status
	EtopPaymentStatus status4.Status

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status status5.Status

	SyncStatus status4.Status // -1:error, 0:new, 1:created, 2:pending
	SyncStates *model.FulfillmentSyncStates

	// Updated by webhook or querying GHN API
	LastSyncAt time.Time

	ExternalShippingLogs []*model.ExternalShippingLog
	AdminNote            string
	IsPartialDelivery    bool
}

func (f *Fulfillment) SelfURL(baseURL string, accType int) string {
	switch accType {
	case model.TagEtop:
		return ""

	case model.TagShop:
		if baseURL == "" || f.ShopID == 0 || f.ID == 0 {
			return ""
		}
		return fmt.Sprintf("%v/s/%v/fulfillments/%v", baseURL, f.ShopID, f.ID)

	default:
		panic(fmt.Sprintf("unsupported account type: %v", accType))
	}
}

func (f *Fulfillment) BeforeInsert() error {
	return f.BeforeUpdate()
}

func (f *Fulfillment) BeforeUpdate() error {
	if f.AddressTo != nil {
		f.AddressToProvinceCode = f.AddressTo.ProvinceCode
		f.AddressToDistrictCode = f.AddressTo.DistrictCode
		f.AddressToWardCode = f.AddressTo.WardCode
	}
	return nil
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

func (f *Fulfillment) ApplyEtopPrice(price int) error {
	f.EtopPriceRule = true
	f.EtopAdjustedShippingFeeMain = price
	return nil
}
