package convert

import (
	"etop.vn/api/main/etop"
	"etop.vn/backend/pkg/etop/model"
)

func Status3ToModel(s etop.Status3) model.Status3 {
	return model.Status3(s)
}

func Status4ToModel(s etop.Status3) model.Status4 {
	return model.Status4(s)
}

func Status5ToModel(s etop.Status3) model.Status5 {
	return model.Status5(s)
}
