package convert

import (
	api "o.o/backend/zexp/sample/counter/api"
	model "o.o/backend/zexp/sample/counter/model"
)

// +gen:convert: o.o/backend/zexp/sample/counter/model  -> o.o/backend/zexp/sample/counter/api

func convertCounter(arg *model.Counter, out *api.Counter) {
	convert_countermodel_Counter_api_Counter(arg, out)
	out.ValueOne = arg.Value
}

func convertCounterModel(arg *api.Counter, out *model.Counter) {
	convert_api_Counter_countermodel_Counter(arg, out)
	out.Value = arg.ValueOne
}
