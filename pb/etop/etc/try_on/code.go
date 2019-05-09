package try_on

import "etop.vn/backend/pkg/etop/model"

func (x TryOnCode) MarshalJSON() ([]byte, error) {
	return []byte(x.String()), nil
}

func (x *TryOnCode) ToModel() model.TryOn {
	if x == nil || *x == 0 {
		return ""
	}
	return model.TryOn(x.String())
}

func PbTryOn(m model.TryOn) TryOnCode {
	return TryOnCode(TryOnCode_value[string(m)])
}

func PbPtrTryOn(m model.TryOn) *TryOnCode {
	res := PbTryOn(m)
	return &res
}

func PbTryOnPtr(m *string) *TryOnCode {
	if m == nil || *m == "" {
		return nil
	}
	res := PbTryOn(model.TryOn(*m))
	return &res
}
