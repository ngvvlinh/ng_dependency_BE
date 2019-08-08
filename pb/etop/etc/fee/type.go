package fee

import "etop.vn/backend/com/main/ordering/model"

func Pb(s model.OrderFeeType) FeeType {
	st, ok := FeeType_value[string(s)]
	if !ok {
		return FeeType_other
	}
	return FeeType(st)
}

func (s *FeeType) ToModel() model.OrderFeeType {
	if s == nil || *s == 0 {
		return "other"
	}
	return model.OrderFeeType(s.String())
}
