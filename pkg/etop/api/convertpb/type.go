package convertpb

import (
	"etop.vn/api/pb/etop/etc/fee"
	"etop.vn/backend/com/main/ordering/model"
)

func Pb(s model.OrderFeeType) fee.FeeType {
	st, ok := fee.FeeType_value[string(s)]
	if !ok {
		return fee.FeeType_other
	}
	return fee.FeeType(st)
}

func FeeTypeToModel(s *fee.FeeType) model.OrderFeeType {
	if s == nil || *s == 0 {
		return "other"
	}
	return model.OrderFeeType(s.String())
}
