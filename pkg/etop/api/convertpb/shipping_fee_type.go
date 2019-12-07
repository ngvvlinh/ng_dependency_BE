package convertpb

import (
	"etop.vn/api/top/types/etc/shipping_fee_type"
	"etop.vn/backend/pkg/etop/model"
)

func PbShippingFeeType(s model.ShippingFeeLineType) shipping_fee_type.ShippingFeeType {
	st, _ := shipping_fee_type.ParseShippingFeeType(string(s))
	return st
}
