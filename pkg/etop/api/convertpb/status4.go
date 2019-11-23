package convertpb

import (
	"etop.vn/api/pb/etop/etc/status4"
	"etop.vn/backend/pkg/etop/model"
)

func Pb4(s model.Status4) status4.Status {
	if s < 0 {
		return status4.Status(s + 128)
	}
	return status4.Status(s)
}

func Pb4PtrStatus(s model.Status4) *status4.Status {
	res := Pb4(s)
	return &res
}

func Pb4Ptr(s *int) *status4.Status {
	if s == nil {
		return nil
	}
	res := status4.Status(*s)
	if res < 0 {
		res += 128
	}
	return &res
}

func Status4ToModel(s *status4.Status) *model.Status4 {
	if s == nil {
		return nil
	}
	i := model.Status4(*s)
	if i >= 64 {
		i = i - 128
	}
	return &i
}
