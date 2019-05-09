package shipping

import "etop.vn/backend/pkg/etop/model"

func Pb(s model.ShippingState) State {
	st, ok := State_value[string(s)]
	if !ok {
		return State_unknown
	}
	return State(st)
}

func PbPtrShippingState(s model.ShippingState) *State {
	res := Pb(s)
	return &res
}

func PbPtr(s *string) *State {
	if s == nil || *s == "" {
		return nil
	}
	st := Pb(model.ShippingState(*s))
	return &st
}

func (s *State) ToModel() model.ShippingState {
	if s == nil {
		return ""
	}
	return model.ShippingState(State_name[int32(*s)])
}
