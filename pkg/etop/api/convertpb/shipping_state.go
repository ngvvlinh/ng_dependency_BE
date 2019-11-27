package convertpb

import (
	"etop.vn/api/pb/etop/etc/shipping"
	"etop.vn/backend/pkg/etop/model"
)

func PbShippingState(s model.ShippingState) shipping.State {
	st, ok := shipping.State_value[string(s)]
	if !ok {
		return shipping.State_unknown
	}
	return shipping.State(st)
}

func PbPtrShippingState(s model.ShippingState) *shipping.State {
	res := PbShippingState(s)
	return &res
}

func PbPtr(s *string) *shipping.State {
	if s == nil || *s == "" {
		return nil
	}
	st := PbShippingState(model.ShippingState(*s))
	return &st
}

func ShippingStateToModel(s *shipping.State) model.ShippingState {
	if s == nil {
		return ""
	}
	return model.ShippingState(shipping.State_name[int(*s)])
}
