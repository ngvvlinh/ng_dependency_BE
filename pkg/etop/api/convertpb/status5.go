package convertpb

import (
	"etop.vn/api/pb/etop/etc/status5"
	"etop.vn/backend/pkg/etop/model"
)

func Pb5(s model.Status5) status5.Status {
	if s < 0 {
		return status5.Status(s + 128)
	}
	return status5.Status(s)
}

func Pb5PtrStatus(s model.Status5) *status5.Status {
	res := Pb5(s)
	return &res
}

func Pb5Ptr(s *int) *status5.Status {
	if s == nil {
		return nil
	}
	res := status5.Status(*s)
	if res < 0 {
		res += 128
	}
	return &res
}

func Status5ToModel(s *status5.Status) *model.Status5 {
	if s == nil {
		return nil
	}
	i := model.Status5(*s)
	if i >= 64 {
		i = i - 128
	}
	return &i
}
