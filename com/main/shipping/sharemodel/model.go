package sharemodel

import (
	"time"

	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/shipping_fee_type"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

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

// +convert:type=shipping.ShippingFeeLine
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
	if contains(ShippingFeeShopTypes, item.ShippingFeeType) {
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