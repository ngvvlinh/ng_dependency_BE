package convertpb

import (
	"etop.vn/api/top/types/etc/shipping_fee_type"
	"etop.vn/backend/pkg/etop/model"
)

func PbShippingFeeType(s model.ShippingFeeLineType) shipping_fee_type.ShippingFeeType {
	st, ok := shipping_fee_type.ParseShippingFeeType(string(s))
	if !ok {
		return shipping_fee_type.Unknown
	}
	return shipping_fee_type.ShippingFeeType(st)
}
