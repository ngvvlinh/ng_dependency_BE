package convertpb

import (
	"etop.vn/api/top/types/etc/source"
	"etop.vn/backend/pkg/etop/model"
)

func PbSource(s model.OrderSourceType) source.Source {
	return source.Source(source.Source_value[string(s)])
}

func SourceToModel(s source.Source) model.OrderSourceType {
	return model.OrderSourceType(source.Source_name[int(s)])
}
