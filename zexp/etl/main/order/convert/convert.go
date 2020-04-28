package convert

import (
	orderingmodel "o.o/backend/com/main/ordering/model"
	"o.o/backend/zexp/etl/main/order/model"
)

// +gen:convert: o.o/backend/zexp/etl/main/order/model -> o.o/backend/com/main/ordering/model

func ConvertOrder(in *orderingmodel.Order, out *model.Order) {
	convert_orderingmodel_Order_ordermodel_Order(in, out)
	for _, feeLine := range in.FeeLines {
		out.FeeLines = append(out.FeeLines, model.OrderFeeLine{
			Amount: feeLine.Amount,
			Desc:   feeLine.Desc,
			Code:   feeLine.Code,
			Name:   feeLine.Name,
			Type:   feeLine.Type,
		})
	}

}
