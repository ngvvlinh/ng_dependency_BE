package status5

import (
	"etop.vn/backend/pkg/etop/model"
)

func Pb(s model.Status5) Status {
	if s < 0 {
		return Status(s + 128)
	}
	return Status(s)
}

func PbPtrStatus(s model.Status5) *Status {
	res := Pb(s)
	return &res
}

func PbPtr(s *int) *Status {
	if s == nil {
		return nil
	}
	res := Status(*s)
	if res < 0 {
		res += 128
	}
	return &res
}

func (s *Status) ToModel() *model.Status5 {
	if s == nil {
		return nil
	}
	i := model.Status5(*s)
	if i >= 64 {
		i = i - 128
	}
	return &i
}
