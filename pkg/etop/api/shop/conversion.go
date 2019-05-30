package shop

import (
	catalogtypes "etop.vn/api/main/catalog/v1/types"
	orderingv1types "etop.vn/api/main/ordering/v1/types"
	"etop.vn/api/main/shipnow"
	shippingtypes "etop.vn/api/main/shipping/v1/types"
	orderP "etop.vn/backend/pb/etop/order"
)

func Convert_api_ShipnowFulfillment_To_core_ShipnowFulfillment(in *orderP.ShipnowFulfillment) *shipnow.ShipnowFulfillment {
	return &shipnow.ShipnowFulfillment{
		Id:                  in.Id,
		ShopId:              in.ShopId,
		PartnerId:           in.PartnerId,
		PickupAddress:       Convert_api_OrderAddress_To_core_OrderAddress(in.PickupAddress),
		DeliveryPoints:      Conver_api_DeliveryPoints_To_core_DeliveryPoints(in.DeliveryPoints),
		Carrier:             in.Carrier,
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		WeightInfo:          Convert_api_WeightInfo_To_core_WeightInfo(in.WeightInfo),
		ValueInfo:           Convert_api_ValueInfo_To_core_ValueInfo(in.ValueInfo),
		ShippingNote:        in.ShippingNote,
		RequestPickupAt:     nil,
	}
}

func Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in *shipnow.ShipnowFulfillment) *orderP.ShipnowFulfillment {
	if in == nil {
		return nil
	}
	return &orderP.ShipnowFulfillment{
		Id:                  in.Id,
		ShopId:              in.ShopId,
		PartnerId:           in.PartnerId,
		PickupAddress:       Convert_core_OrderAddress_To_api_OrderAddress(in.PickupAddress),
		DeliveryPoints:      Convert_core_DeliveryPoints_To_api_DeliveryPoints(in.DeliveryPoints),
		Carrier:             in.Carrier,
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		WeightInfo:          Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:           Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		ShippingNote:        in.ShippingNote,
		RequestPickupAt:     nil,
	}
}

func Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(ins []*shipnow.ShipnowFulfillment) (outs []*orderP.ShipnowFulfillment) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in))
	}
	return
}

func Convert_api_DeliveryPoint_To_core_DeliveryPoint(in *orderP.DeliveryPoint) *shipnow.DeliveryPoint {
	return &shipnow.DeliveryPoint{
		ShippingAddress: Convert_api_OrderAddress_To_core_OrderAddress(in.ShippingAddress),
		Lines:           Convert_api_OrderLines_To_core_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_api_WeightInfo_To_core_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_api_ValueInfo_To_core_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Conver_api_DeliveryPoints_To_core_DeliveryPoints(ins []*orderP.DeliveryPoint) (outs []*shipnow.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_api_DeliveryPoint_To_core_DeliveryPoint(in))
	}
	return
}

func Convert_core_DeliveryPoint_To_api_DeliveryPoint(in *shipnow.DeliveryPoint) *orderP.DeliveryPoint {
	return &orderP.DeliveryPoint{
		ShippingAddress: Convert_core_OrderAddress_To_api_OrderAddress(in.ShippingAddress),
		Lines:           Convert_core_OrderLines_To_api_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Convert_core_DeliveryPoints_To_api_DeliveryPoints(ins []*shipnow.DeliveryPoint) (outs []*orderP.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_core_DeliveryPoint_To_api_DeliveryPoint(in))
	}
	return
}

func Convert_api_WeightInfo_To_core_WeightInfo(in orderP.WeightInfo) shippingtypes.WeightInfo {
	return shippingtypes.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_core_WeightInfo_To_api_WeightInfo(in shippingtypes.WeightInfo) orderP.WeightInfo {
	return orderP.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_api_ValueInfo_To_core_ValueInfo(in orderP.ValueInfo) shippingtypes.ValueInfo {
	return shippingtypes.ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_core_ValueInfo_To_api_ValueInfo(in shippingtypes.ValueInfo) orderP.ValueInfo {
	return orderP.ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_api_OrderLine_To_core_OrderLine(in *orderP.OrderLine) (out *orderingv1types.ItemLine) {
	return &orderingv1types.ItemLine{
		OrderId:   in.OrderId,
		ProductId: in.ProductId,
		Quantity:  in.Quantity,
		VariantId: in.VariantId,
		IsOutside: in.IsOutsideEtop,
		ProductInfo: orderingv1types.ProductInfo{
			ProductName: in.ProductName,
			ImageUrl:    in.ImageUrl,
			Attributes:  Convert_api_Attributes_To_core_Atributes(in.Attributes),
		},
	}
}

func Convert_api_OrderLines_To_core_OrderLines(ins []*orderP.OrderLine) (outs []*orderingv1types.ItemLine) {
	res := make([]*orderingv1types.ItemLine, len(ins))
	for i, in := range ins {
		res[i] = Convert_api_OrderLine_To_core_OrderLine(in)
	}
	return res
}

func Convert_core_OrderLine_To_api_OrderLine(in *orderingv1types.ItemLine) *orderP.OrderLine {
	return &orderP.OrderLine{
		OrderId:       in.OrderId,
		VariantId:     in.VariantId,
		ProductName:   in.ProductInfo.ProductName,
		IsOutsideEtop: in.IsOutside,
		Quantity:      in.Quantity,
		ImageUrl:      in.ProductInfo.ImageUrl,
		Attributes:    Convert_core_Attribute_To_api_Atribures(in.ProductInfo.Attributes),
		ProductId:     in.ProductId,
	}
}

func Convert_core_OrderLines_To_api_OrderLines(ins []*orderingv1types.ItemLine) (outs []*orderP.OrderLine) {
	for _, in := range ins {
		outs = append(outs, Convert_core_OrderLine_To_api_OrderLine(in))
	}
	return
}

func Convert_api_Attributes_To_core_Atributes(ins []*orderP.Attribute) []*catalogtypes.Attribute {
	res := make([]*catalogtypes.Attribute, len(ins))
	for i, in := range ins {
		res[i] = &catalogtypes.Attribute{
			Name:  in.Name,
			Value: in.Value,
		}
	}
	return res
}

func Convert_core_Attribute_To_api_Atribures(ins []*catalogtypes.Attribute) (outs []*orderP.Attribute) {
	for _, in := range ins {
		outs = append(outs, &orderP.Attribute{
			Name:  in.Name,
			Value: in.Value,
		})
	}
	return
}

func Convert_api_OrderAddress_To_core_OrderAddress(in *orderP.OrderAddress) *orderingv1types.Address {
	if in == nil {
		return nil
	}
	return &orderingv1types.Address{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Company:  in.Company,
		Address1: in.Address1,
		Address2: in.Address2,
		Location: orderingv1types.Location{
			ProvinceCode: in.ProvinceCode,
			DistrictCode: in.DistrictCode,
			WardCode:     in.WardCode,
			Coordinates:  Convert_api_Coordinates_To_core_Coordinates(in.Coordinates),
		},
	}
}

func Convert_core_OrderAddress_To_api_OrderAddress(in *orderingv1types.Address) *orderP.OrderAddress {
	if in == nil {
		return nil
	}
	return &orderP.OrderAddress{
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Address1:     in.Address1,
		Address2:     in.Address2,
		ProvinceCode: in.Location.ProvinceCode,
		DistrictCode: in.Location.DistrictCode,
		WardCode:     in.Location.WardCode,
		Coordinates:  Convert_core_Coordinates_To_api_Coordinates(in.Coordinates),
	}
}

func Convert_api_Coordinates_To_core_Coordinates(in *orderP.Coordinates) *orderingv1types.Coordinates {
	if in == nil {
		return nil
	}
	return &orderingv1types.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func Convert_core_Coordinates_To_api_Coordinates(in *orderingv1types.Coordinates) *orderP.Coordinates {
	if in == nil {
		return nil
	}
	return &orderP.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}
