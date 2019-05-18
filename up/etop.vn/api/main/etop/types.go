package etop

import (
	etopv1types "etop.vn/api/main/etop/v1"
)

type (
	Status3 = etopv1types.Status3
	Status4 = etopv1types.Status4
	Status5 = etopv1types.Status5
)

const (
	S3Negative = etopv1types.Status3_N
	S3Zero     = etopv1types.Status3_Z
	S3Positive = etopv1types.Status3_P

	S4Negative = etopv1types.Status4_N
	S4Zero     = etopv1types.Status4_Z
	S4SuperPos = etopv1types.Status4_S
	S4Positive = etopv1types.Status4_P

	S5NegSuper = etopv1types.Status5_NS
	S5Negative = etopv1types.Status5_N
	S5Zero     = etopv1types.Status5_Z
	S5Positive = etopv1types.Status5_P
	S5SuperPos = etopv1types.Status5_S
)

func Status3FromInt(s int) Status3 {
	if s < 0 {
		return Status3(s + 128)
	}
	return Status3(s)
}

func Status4FromInt(s int) Status4 {
	if s < 0 {
		return Status4(s + 128)
	}
	return Status4(s)
}

func Status5FromInt(s int) Status5 {
	if s < 0 {
		return Status5(s + 128)
	}
	return Status5(s)
}
