package convertpb

import (
	"etop.vn/api/pb/etop/order/source"
	"etop.vn/backend/pkg/etop/model"
)

func PbSource(s model.OrderSourceType) source.Source {
	return source.Source(source.Source_value[string(s)])
}

func SourceToModel(s source.Source) model.OrderSourceType {
	return model.OrderSourceType(source.Source_name[int(s)])
}
