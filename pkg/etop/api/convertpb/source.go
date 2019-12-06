package convertpb

import (
	"etop.vn/api/top/types/etc/source"
	"etop.vn/backend/pkg/etop/model"
)

func PbSource(s model.OrderSourceType) source.Source {
	value, _ := source.ParseSource(string(s))
	return value
}

func SourceToModel(s source.Source) model.OrderSourceType {
	return model.OrderSourceType(s.Name())
}
