package convertpb

import (
	"etop.vn/api/pb/etop/etc/shipping_fee_type"
	"etop.vn/backend/pkg/etop/model"
)

func PbShippingFeeType(s model.ShippingFeeLineType) shipping_fee_type.ShippingFeeType {
	st, ok := shipping_fee_type.ShippingFeeType_value[string(s)]
	if !ok {
		return shipping_fee_type.ShippingFeeType_unknown
	}
	return shipping_fee_type.ShippingFeeType(st)
}
