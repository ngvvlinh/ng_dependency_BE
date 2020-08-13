package convert

import (
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	orderconvert "o.o/backend/com/main/ordering/convert"
	shippingmodel "o.o/backend/com/main/shipping/model"
)

// +gen:convert: o.o/backend/com/main/shipping/model -> o.o/api/main/shipping, o.o/api/main/shipping/types
// +gen:convert: o.o/backend/com/main/shipping/sharemodel -> o.o/api/main/shipping/types
// +gen:convert: o.o/backend/com/main/shipping/sharemodel -> o.o/api/main/shipping
// +gen:convert: o.o/backend/com/main/shipping/modely -> o.o/api/main/shipping
// +gen:convert: o.o/api/main/shipping

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
		IncludeInsurance: in.IncludeInsurance.Apply(false),
		InsuranceValue:   in.InsuranceValue,
	}
	return
}
