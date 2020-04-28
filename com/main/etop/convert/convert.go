package convert

import (
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/pkg/etop/model"
)

// +gen:convert: o.o/backend/pkg/etop/model -> o.o/api/meta

func Status3(in status3.Status) (out status3.Status) {
	if in < 0 {
		return status3.Status(in + 128)
	}
	return status3.Status(in)
}

func Status3ToModel(in status3.Status) (out status3.Status) {
	out = status3.Status(in)
	if out >= 64 {
		out -= 128
	}
	return out
}

func Status4(in status4.Status) (out status4.Status) {
	if in < 0 {
		return status4.Status(in + 128)
	}
	return status4.Status(in)
}

func Status4ToModel(in status4.Status) (out status4.Status) {
	out = status4.Status(in)
	if out >= 64 {
		out -= 128
	}
	return out
}

func Status5(in status5.Status) (out status5.Status) {
	if in < 0 {
		return status5.Status(in + 128)
	}
	return status5.Status(in)
}

func Status5ToModel(in status5.Status) (out status5.Status) {
	out = status5.Status(in)
	if out >= 64 {
		out -= 128
	}
	return out
}

func ErrorToModel(in *meta.Error) (out *model.Error) {
	if in == nil {
		return nil
	}
	return &model.Error{
		Code: in.Code,
		Msg:  in.Msg,
		Meta: in.Meta,
	}
}

func Error(in *model.Error) (out *meta.Error) {
	if in == nil {
		return nil
	}
	return &meta.Error{
		Code: in.Code,
		Msg:  in.Msg,
		Meta: in.Meta,
	}
}
