package sharemodel

import (
	"time"

	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

var ShippingStateMap = map[shipping.State]string{
	shipping.Default:       "Mặc định",
	shipping.Created:       "Mới",
	shipping.Picking:       "Đang lấy hàng",
	shipping.Holding:       "Đã lấy hàng",
	shipping.Delivering:    "Đang giao hàng",
	shipping.Returning:     "Đang trả hàng",
	shipping.Delivered:     "Đã giao hàng",
	shipping.Returned:      "Đã trả hàng",
	shipping.Cancelled:     "Hủy",
	shipping.Undeliverable: "Bồi hoàn",
	shipping.Unknown:       "Không xác định",
}

type ShippingFeeLine struct {
	ShippingFeeType          shipping_fee_type.ShippingFeeType `json:"shipping_fee_type"`
	Cost                     int                               `json:"cost"`
	ExternalServiceID        string                            `json:"external_service_id"`
	ExternalServiceName      string                            `json:"external_service_name"`
	ExternalServiceType      string                            `json:"external_service_type"`
	ExternalShippingOrderID  string                            `json:"external_order_id"`
	ExternalPaymentChannelID string                            `json:"external_payment_channel_id"`
	ExternalShippingCode     string                            `json:"external_shipping_code"`
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
	if contains(shippingtypes.ShippingFeeShopTypes, item.ShippingFeeType) {
		return &item
	}
	return nil
}

func GetReturnedFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		if item.ShippingFeeType == shipping_fee_type.Return {
			result = item.Cost
			break
		}
	}
	return result
}

func GetMainFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		if item.ShippingFeeType == shipping_fee_type.Main {
			result = item.Cost
			break
		}
	}
	return result
}

func GetTotalShippingFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		result += item.Cost
	}
	return result
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

type FulfillmentSyncStates struct {
	SyncAt    time.Time        `json:"sync_at"`
	TrySyncAt time.Time        `json:"try_sync_at"`
	Error     *etopmodel.Error `json:"error"`

	NextShippingState shipping.State `json:"next_shipping_state"`
}

type AvailableShippingService struct {
	Name string

	// ServiceFee: Tổng phí giao hàng (đã bao gồm phí chính + các phụ phí khác)
	ServiceFee int

	// ShippingFeeMain: Phí chính giao hàng
	ShippingFeeMain   int
	Provider          shipping_provider.ShippingProvider
	ProviderServiceID string

	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
	Source             etopmodel.ShippingPriceSource
	ConnectionInfo     *ConnectionInfo

	// Thông tin các gói được admin định nghĩa
	ShipmentServiceInfo *ShipmentServiceInfo
	ShipmentPriceInfo   *ShipmentPriceInfo
	ShippingFeeLines    []*ShippingFeeLine
}

type ConnectionInfo struct {
	ID       dot.ID
	Name     string
	ImageURL string
}

type ShipmentServiceInfo struct {
	ID           dot.ID
	Code         string
	Name         string
	IsAvailable  bool
	ErrorMessage string
}

type ShipmentPriceInfo struct {
	ShipmentPriceID     dot.ID `json:"shipment_price_id"`
	ShipmentPriceListID dot.ID `json:"shipment_price_list_id"`
	OriginFee           int    `json:"origin_fee,omitempty"`
	MakeupFee           int    `json:"makeup_fee,omitempty"`
}

func (service *AvailableShippingService) ApplyFeeMain(feeMain int) {
	service.ServiceFee = service.ServiceFee - service.ShippingFeeMain + feeMain
	service.ShippingFeeMain = feeMain
}
