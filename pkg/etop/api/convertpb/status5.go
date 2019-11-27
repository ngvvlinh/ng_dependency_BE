package convertpb

import (
	"etop.vn/api/pb/etop/etc/status5"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
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

func Pb5Ptr(s dot.NullInt) *status5.Status {
	if !s.Valid {
		return nil
	}
	res := status5.Status(s.Int)
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
