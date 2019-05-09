package source

import "etop.vn/backend/pkg/etop/model"

func PbSource(s model.OrderSourceType) Source {
	return Source(Source_value[string(s)])
}

func (s Source) ToModel() model.OrderSourceType {
	return model.OrderSourceType(Source_name[int32(s)])
}
