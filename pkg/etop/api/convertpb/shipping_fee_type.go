package convertpb

import (
	"etop.vn/api/top/types/etc/shipping_fee_type"
)

func PbShippingFeeType(s shipping_fee_type.ShippingFeeType) shipping_fee_type.ShippingFeeType {
	st, _ := shipping_fee_type.ParseShippingFeeType(string(s))
	return st
}
