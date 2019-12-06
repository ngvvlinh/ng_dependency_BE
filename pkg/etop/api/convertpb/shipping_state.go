package convertpb

import (
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func PbShippingState(s model.ShippingState) shipping.State {
	st, ok := shipping.ParseState(string(s))
	if !ok {
		return shipping.Unknown
	}
	return st
}

func PbPtrShippingState(s model.ShippingState) *shipping.State {
	res := PbShippingState(s)
	return &res
}

func PbPtr(s dot.NullString) *shipping.State {
	if s.Apply("") == "" {
		return nil
	}
	st := PbShippingState(model.ShippingState(s.String))
	return &st
}

func ShippingStateToModel(s *shipping.State) model.ShippingState {
	if s == nil {
		return ""
	}
	return model.ShippingState(s.String())
}
