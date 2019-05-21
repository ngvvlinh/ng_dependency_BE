package model

import (
	"encoding/json"
	"fmt"
	"time"

	"etop.vn/backend/pkg/etop/model"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFulfillment(&Fulfillment{})

type Fulfillment struct {
	ID         int64
	OrderID    int64
	ShopID     int64
	SupplierID int64
	PartnerID  int64

	SupplierConfirm model.Status3
	ShopConfirm     model.Status3
	ConfirmStatus   model.Status3

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

	VariantIDs []int64
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

	MoneyTransactionID                 int64
	MoneyTransactionShippingExternalID int64

	CancelReason string

	// CreatedBy   int64
	// UpdatedBy   int64
	// CancelledBy int64

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
	ExternalShippingStatus      model.Status5
	ExternalShippingNote        string
	ExternalShippingSubState    string

	ExternalShippingData json.RawMessage

	ShippingState     model.ShippingState
	ShippingStatus    model.Status5
	EtopPaymentStatus model.Status4

	// -1:cancelled, 0:default, 1:delivered, 2:processing
	//
	// 0: just created, still error, can be edited
	// 2: processing, can not be edited
	// 1: done
	// -1: cancelled
	// -2: returned
	Status model.Status5

	SyncStatus model.Status4 // -1:error, 0:new, 1:created, 2:pending
	SyncStates *model.FulfillmentSyncStates

	// Updated by webhook or querying GHN API
	LastSyncAt time.Time

	ExternalShippingLogs []*model.ExternalShippingLog
	AdminNote            string
	IsPartialDelivery    bool
}

func (f *Fulfillment) SelfURL(baseURL string, accType int) string {
	switch accType {
	case model.TagEtop, model.TagSupplier:
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
