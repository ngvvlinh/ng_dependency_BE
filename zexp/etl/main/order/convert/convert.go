package convert

import (
	orderingmodel "etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/zexp/etl/main/order/model"
)

// +gen:convert: etop.vn/backend/zexp/etl/main/order/model->etop.vn/backend/com/main/ordering/model

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
