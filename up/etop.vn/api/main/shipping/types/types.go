package types

import (
	"errors"

	orderingtypes "etop.vn/api/main/ordering/v1/types"
	shippingv1types "etop.vn/api/main/shipping/v1/types"
)

type TryOn = shippingv1types.TryOnCode

const (
	TryOnOpen = shippingv1types.TryOnCode_open
	TryOnNone = shippingv1types.TryOnCode_none
	TryOnTry  = shippingv1types.TryOnCode_try
)

func TryOnFromString(s string) (TryOn, error) {
	t := shippingv1types.TryOnCode_value[s]
	if t == 0 {
		return 0, errors.New("invalid tryon code")
	}
	return TryOn(t), nil
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
