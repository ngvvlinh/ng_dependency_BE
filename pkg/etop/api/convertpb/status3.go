package convertpb

import (
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func Pb3(s model.Status3) status3.Status {
	if s < 0 {
		return status3.Status(s + 128)
	}
	return status3.Status(s)
}

func Pb3PtrStatus(s model.Status3) *status3.Status {
	res := Pb3(s)
	return &res
}

func Pb3Ptr(s dot.NullInt) *status3.Status {
	if !s.Valid {
		return nil
	}
	res := status3.Status(s.Int)
	if res < 0 {
		res += 128
	}
	return &res
}

func Status3ToModel(s *status3.Status) *model.Status3 {
	if s == nil {
		return nil
	}
	i := model.Status3(*s)
	if i >= 64 {
		i = i - 128
	}
	return &i
}
