package shipping_fee_type

import "etop.vn/backend/pkg/etop/model"

func Pb(s model.ShippingFeeLineType) ShippingFeeType {
	st, ok := ShippingFeeType_value[string(s)]
	if !ok {
		return ShippingFeeType_unknown
	}
	return ShippingFeeType(st)
}
