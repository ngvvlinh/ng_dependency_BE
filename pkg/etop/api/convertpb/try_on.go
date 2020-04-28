package convertpb

import (
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

func PbPtrTryOn(m try_on.TryOnCode) try_on.NullTryOnCode {
	return m.Wrap()
}

func PbTryOnPtr(m dot.NullString) try_on.NullTryOnCode {
	if m.Apply("") == "" {
		return try_on.NullTryOnCode{}
	}
	code, _ := try_on.ParseTryOnCode(m.String)
	return code.Wrap()
}
