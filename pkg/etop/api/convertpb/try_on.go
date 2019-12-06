package convertpb

import (
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func TryOnCodeToModel(x *try_on.TryOnCode) model.TryOn {
	if x == nil || *x == 0 {
		return ""
	}
	return model.TryOn(x.String())
}

func PbTryOn(m model.TryOn) try_on.TryOnCode {
	value, _ := try_on.ParseTryOnCode(string(m))
	return value
}

func PbPtrTryOn(m model.TryOn) *try_on.TryOnCode {
	res := PbTryOn(m)
	return &res
}

func PbTryOnPtr(m dot.NullString) *try_on.TryOnCode {
	if m.Apply("") == "" {
		return nil
	}
	res := PbTryOn(model.TryOn(m.String))
	return &res
}
