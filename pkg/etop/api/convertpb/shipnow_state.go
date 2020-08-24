package convertpb

import (
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/capi/dot"
)

func ShipnowNullState(s dot.NullString) shipnow_state.NullState {
	if s.Apply("") == "" {
		return shipnow_state.NullState{}
	}
	st, ok := shipnow_state.ParseState(s.String)
	if !ok {
		return shipnow_state.StateDefault.Wrap()
	}
	return st.Wrap()
}
