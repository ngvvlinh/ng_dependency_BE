package convertpb

import (
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
)

func Pb3Ptr(s dot.NullInt) status3.NullStatus {
	if !s.Valid {
		return status3.NullStatus{}
	}
	return status3.Status(s.Int).Wrap()
}

func Pb4Ptr(s dot.NullInt) status4.NullStatus {
	if !s.Valid {
		return status4.NullStatus{}
	}
	return status4.Status(s.Int).Wrap()
}

func Pb5Ptr(s dot.NullInt) status5.NullStatus {
	if !s.Valid {
		return status5.NullStatus{}
	}
	return status5.Status(s.Int).Wrap()
}
