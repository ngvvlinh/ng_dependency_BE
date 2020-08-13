package types

import (
	"time"

	orderingtypes "o.o/api/main/ordering/types"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
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

type ShippingInfo struct {
	PickupAddress       *orderingtypes.Address
	ReturnAddress       *orderingtypes.Address
	ShippingServiceName string
	ShippingServiceCode string
	ShippingServiceFee  int
	Carrier             shipping_provider.ShippingProvider
	IncludeInsurance    bool
	TryOn               try_on.TryOnCode
	ShippingNote        string
	CODAmount           int
	GrossWeight         int
	Length              int
	Width               int
	Height              int
	ChargeableWeight    int
}

type ShippingService struct {
	Carrier             string
	Code                string
	Fee                 int
	Name                string
	EstimatedPickupAt   time.Time
	EstimatedDeliveryAt time.Time
}

type WeightInfo struct {
	GrossWeight      int
	ChargeableWeight int
	Length           int
	Width            int
	Height           int
}

type ValueInfo struct {
	BasketValue      int
	CODAmount        int
	IncludeInsurance bool
	InsuranceValue   dot.NullInt
}

type ShippingFeeLine struct {
	ShippingFeeType     shipping_fee_type.ShippingFeeType
	Cost                int
	ExternalServiceID   string
	ExternalServiceName string
	ExternalServiceType string
}

func GetTotalShippingFee(items []*ShippingFeeLine) int {
	result := 0
	for _, item := range items {
		result += item.Cost
	}
	return result
}

func GetShippingFeeShopLines(items []*ShippingFeeLine, etopPriceRule bool, mainFee dot.NullInt) []*ShippingFeeLine {
	if items == nil {
		return nil
	}
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
