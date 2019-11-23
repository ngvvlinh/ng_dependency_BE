package convertpb

import (
	"etop.vn/api/pb/etop/etc/try_on"
	"etop.vn/backend/pkg/etop/model"
)

func TryOnCodeToModel(x *try_on.TryOnCode) model.TryOn {
	if x == nil || *x == 0 {
		return ""
	}
	return model.TryOn(x.String())
}

func PbTryOn(m model.TryOn) try_on.TryOnCode {
	return try_on.TryOnCode(try_on.TryOnCode_value[string(m)])
}

func PbPtrTryOn(m model.TryOn) *try_on.TryOnCode {
	res := PbTryOn(m)
	return &res
}

func PbTryOnPtr(m *string) *try_on.TryOnCode {
	if m == nil || *m == "" {
		return nil
	}
	res := PbTryOn(model.TryOn(*m))
	return &res
}
