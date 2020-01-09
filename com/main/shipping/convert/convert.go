package convert

import (
	"etop.vn/api/main/shipping"
	shippingtypes "etop.vn/api/main/shipping/types"
	orderconvert "etop.vn/backend/com/main/ordering/convert"
	shippingmodel "etop.vn/backend/com/main/shipping/model"
)

// +gen:convert: etop.vn/backend/com/main/shipping/model->etop.vn/api/main/shipping,etop.vn/api/main/shipping/types
// +gen:convert: etop.vn/backend/com/main/shipping/sharemodel->etop.vn/api/main/shipping
// +gen:convert: etop.vn/backend/com/main/shipping/modely->etop.vn/api/main/shipping
// +gen:convert: etop.vn/api/main/shipping

func Fulfillment(in *shippingmodel.Fulfillment, out *shipping.Fulfillment) {
	if in == nil {
		return
	}
	convert_shippingmodel_Fulfillment_shipping_Fulfillment(in, out)
	out.Lines = orderconvert.OrderLines(in.Lines)
	out.WeightInfo = shippingtypes.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
	out.ValueInfo = shippingtypes.ValueInfo{
		BasketValue:      in.BasketValue,
		CODAmount:        in.TotalCODAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
	return
}
