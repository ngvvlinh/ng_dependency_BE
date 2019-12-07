package convertpb

import (
	"etop.vn/api/top/types/etc/shipping"
	"etop.vn/capi/dot"
)

func PbPtr(s dot.NullString) shipping.NullState {
	if s.Apply("") == "" {
		return shipping.NullState{}
	}
	st, ok := shipping.ParseState(s.String)
	if !ok {
		return shipping.Unknown.Wrap()
	}
	return st.Wrap()
}
