package convert

import (
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	orderconvert "o.o/backend/com/main/ordering/convert"
	shippingmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
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
	out.ExternalShippingCreatedAt = cm.CoalesceTime(in.ExternalShippingCreatedAt, in.ShippingCreatedAt)
	out.ExternalShippingCancelledAt = cm.CoalesceTime(in.ExternalShippingCancelledAt, in.ShippingCancelledAt)
	out.ExternalShippingDeliveredAt = cm.CoalesceTime(in.ExternalShippingDeliveredAt, in.ShippingDeliveredAt)
	out.ExternalShippingReturningAt = in.ShippingReturningAt
	out.ExternalShippingReturnedAt = cm.CoalesceTime(in.ExternalShippingReturnedAt, in.ShippingReturnedAt)
	return
}

func Convert_shipping_Fulfillment_To_shippingmodel_Fulfillment(in *shipping.Fulfillment, out *shippingmodel.Fulfillment) {
	if in == nil {
		return
	}
	convert_shipping_Fulfillment_shippingmodel_Fulfillment(in, out)
	out.Lines = orderconvert.OrderLinesToModel(in.Lines)
	out.GrossWeight = in.GrossWeight
	out.ChargeableWeight = in.ChargeableWeight
	out.Length = in.Length
	out.Width = in.Width
	out.Height = in.Height
	out.BasketValue = in.BasketValue
	out.TotalCODAmount = in.CODAmount
	out.IncludeInsurance = in.IncludeInsurance
	out.InsuranceValue = in.InsuranceValue
}
