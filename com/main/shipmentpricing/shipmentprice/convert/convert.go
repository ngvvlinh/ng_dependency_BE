package convert

import (
	"o.o/api/main/shipmentpricing/shipmentprice"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
)

// +gen:convert: o.o/backend/com/main/shipmentpricing/shipmentprice/model -> o.o/api/main/shipmentpricing/shipmentprice
// +gen:convert: o.o/api/main/shipmentpricing/shipmentprice

func Convert_shipmentprice_ShippingFee_To_shippingsharemodel_ShippingFeeLine(in *shipmentprice.ShippingFee) *shippingsharemodel.ShippingFeeLine {
	if in == nil {
		return nil
	}
	return &shippingsharemodel.ShippingFeeLine{
		ShippingFeeType: in.FeeType,
		Cost:            in.Price,
	}
}

func Convert_shipmentprice_ShippingFees_To_shippingsharemodel_ShippingFeeLines(items []*shipmentprice.ShippingFee) []*shippingsharemodel.ShippingFeeLine {
	result := make([]*shippingsharemodel.ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = Convert_shipmentprice_ShippingFee_To_shippingsharemodel_ShippingFeeLine(item)
	}
	return result
}
