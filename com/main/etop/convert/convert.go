package convert

import (
	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/etop/model"
)

func Status3(in model.Status3) (out etop.Status3) {
	if in < 0 {
		return etop.Status3(in + 128)
	}
	return etop.Status3(in)
}

func Status3ToModel(in etop.Status3) (out model.Status3) {
	out = model.Status3(in)
	if out >= 64 {
		out -= 128
	}
	return out
}

func Status4(in model.Status4) (out etop.Status4) {
	if in < 0 {
		return etop.Status4(in + 128)
	}
	return etop.Status4(in)
}

func Status4ToModel(in etop.Status4) (out model.Status4) {
	out = model.Status4(in)
	if out >= 64 {
		out -= 128
	}
	return out
}

func Status5(in model.Status5) (out etop.Status5) {
	if in < 0 {
		return etop.Status5(in + 128)
	}
	return etop.Status5(in)
}

func Status5ToModel(in etop.Status5) (out model.Status5) {
	out = model.Status5(in)
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
