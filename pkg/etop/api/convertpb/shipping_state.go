package convertpb

import (
	"o.o/api/top/types/etc/shipping"
	"o.o/capi/dot"
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
