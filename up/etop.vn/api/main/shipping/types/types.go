package types

import (
	"errors"

	orderingtypes "etop.vn/api/main/ordering/v1/types"
	shippingv1types "etop.vn/api/main/shipping/v1/types"
)

type TryOn = shippingv1types.TryOnCode
type FeeLineType = shippingv1types.FeeLineType

const (
	TryOnOpen = shippingv1types.TryOnCode_open
	TryOnNone = shippingv1types.TryOnCode_none
	TryOnTry  = shippingv1types.TryOnCode_try

	FeeLineTypeOther         = shippingv1types.FeeLineType_other
	FeeLineTypeMain          = shippingv1types.FeeLineType_main
	FeeLineTypeReturn        = shippingv1types.FeeLineType_return
	FeeLineTypeAdjustment    = shippingv1types.FeeLineType_adjustment
	FeeLineTypeCods          = shippingv1types.FeeLineType_cods
	FeeLineTypeInsurance     = shippingv1types.FeeLineType_insurance
	FeeLineTypeAddressChange = shippingv1types.FeeLineType_address_change
	FeeLineDiscount          = shippingv1types.FeeLineType_discount
)

func TryOnFromString(s string) (TryOn, error) {
	t := shippingv1types.TryOnCode_value[s]
	if t == 0 {
		return 0, errors.New("invalid tryon code")
	}
	return TryOn(t), nil
}

func FeelineTypeFromString(s string) FeeLineType {
	f, ok := shippingv1types.FeeLineType_value[s]
	if !ok {
		f = 0
	}
	return FeeLineType(f)
}

type ShippingInfo struct {
	PickupAddress       *orderingtypes.Address
	ReturnAddress       *orderingtypes.Address
	ShippingServiceName string
	ShippingServiceCode string
	ShippingServiceFee  int
	Carrier             string
	IncludeInsurance    bool
	TryOn               TryOn
	ShippingNote        string
	CODAmount           int
	GrossWeight         int
	Length              int
	Height              int
	ChargeableWeight    int
}

type WeightInfo = shippingv1types.WeightInfo
type ValueInfo = shippingv1types.ValueInfo

type FeeLine = shippingv1types.FeeLine

func TotalFee(lines []*FeeLine) int {
	res := 0
	for _, line := range lines {
		res += int(line.Cost)
	}
	return res
}
