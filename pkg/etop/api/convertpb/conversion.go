package convertpb

import (
	catalogtypes "etop.vn/api/main/catalog/types"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	carrier "etop.vn/api/main/shipnow/carrier/types"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/pb/etop"
	"etop.vn/api/pb/etop/order"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/model"
)

func Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in *shipnow.ShipnowFulfillment) *order.ShipnowFulfillment {
	if in == nil {
		return nil
	}
	return &order.ShipnowFulfillment{
		Id:                         in.Id,
		ShopId:                     in.ShopId,
		PartnerId:                  in.PartnerId,
		PickupAddress:              Convert_core_OrderAddress_To_api_OrderAddress(in.PickupAddress),
		DeliveryPoints:             Convert_core_DeliveryPoints_To_api_DeliveryPoints(in.DeliveryPoints),
		Carrier:                    in.Carrier.String(),
		ShippingServiceCode:        in.ShippingServiceCode,
		ShippingServiceFee:         in.ShippingServiceFee,
		ShippingServiceName:        in.ShippingServiceName,
		ShippingServiceDescription: in.ShippingServiceDescription,
		WeightInfo:                 Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:                  Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		ShippingNote:               in.ShippingNote,
		RequestPickupAt:            cmapi.PbTime(in.RequestPickupAt),
		CreatedAt:                  cmapi.PbTime(in.CreatedAt),
		UpdatedAt:                  cmapi.PbTime(in.UpdatedAt),
		Status:                     Pb5(model.Status5(in.Status)),
		ShippingStatus:             Pb5(model.Status5(in.ShippingStatus)),
		ShippingState:              shipnowtypes.StateToString(in.ShippingState),
		ConfirmStatus:              Pb3(model.Status3(in.ConfirmStatus)),
		OrderIds:                   in.OrderIds,
		ShippingCreatedAt:          cmapi.PbTime(in.ShippingCreatedAt),
		ShippingCode:               in.ShippingCode,
		EtopPaymentStatus:          Pb4(model.Status4(in.EtopPaymentStatus)),
		CodEtopTransferedAt:        cmapi.PbTime(in.CodEtopTransferedAt),
		ShippingPickingAt:          cmapi.PbTime(in.ShippingPickingAt),
		ShippingDeliveringAt:       cmapi.PbTime(in.ShippingDeliveringAt),
		ShippingDeliveredAt:        cmapi.PbTime(in.ShippingDeliveredAt),
		ShippingCancelledAt:        cmapi.PbTime(in.ShippingCancelledAt),
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
	}
}

func Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(ins []*shipnow.ShipnowFulfillment) (outs []*order.ShipnowFulfillment) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in))
	}
	return
}

func Convert_api_DeliveryPoint_To_core_DeliveryPoint(in *order.DeliveryPoint) *shipnow.DeliveryPoint {
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

func Conver_api_DeliveryPoints_To_core_DeliveryPoints(ins []*order.DeliveryPoint) (outs []*shipnow.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_api_DeliveryPoint_To_core_DeliveryPoint(in))
	}
	return
}

func Convert_core_DeliveryPoint_To_api_DeliveryPoint(in *shipnow.DeliveryPoint) *order.DeliveryPoint {
	return &order.DeliveryPoint{
		ShippingAddress: Convert_core_OrderAddress_To_api_OrderAddress(in.ShippingAddress),
		Lines:           Convert_core_OrderLines_To_api_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Convert_core_DeliveryPoints_To_api_DeliveryPoints(ins []*shipnow.DeliveryPoint) (outs []*order.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_core_DeliveryPoint_To_api_DeliveryPoint(in))
	}
	return
}

func Convert_api_WeightInfo_To_core_WeightInfo(in order.WeightInfo) shippingtypes.WeightInfo {
	return shippingtypes.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_core_WeightInfo_To_api_WeightInfo(in shippingtypes.WeightInfo) order.WeightInfo {
	return order.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_api_ValueInfo_To_core_ValueInfo(in order.ValueInfo) shippingtypes.ValueInfo {
	return shippingtypes.ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_core_ValueInfo_To_api_ValueInfo(in shippingtypes.ValueInfo) order.ValueInfo {
	return order.ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_api_OrderLine_To_core_OrderLine(in *order.OrderLine) (out *ordertypes.ItemLine) {
	return &ordertypes.ItemLine{
		OrderId:   in.OrderId,
		ProductId: in.ProductId,
		Quantity:  in.Quantity,
		VariantId: in.VariantId,
		IsOutside: in.IsOutsideEtop,
		ProductInfo: ordertypes.ProductInfo{
			ProductName: in.ProductName,
			ImageUrl:    in.ImageUrl,
			Attributes:  Convert_api_Attributes_To_core_Attributes(in.Attributes),
		},
	}
}

func Convert_api_OrderLines_To_core_OrderLines(ins []*order.OrderLine) (outs []*ordertypes.ItemLine) {
	res := make([]*ordertypes.ItemLine, len(ins))
	for i, in := range ins {
		res[i] = Convert_api_OrderLine_To_core_OrderLine(in)
	}
	return res
}

func Convert_core_OrderLine_To_api_OrderLine(in *ordertypes.ItemLine) *order.OrderLine {
	return &order.OrderLine{
		OrderId:       in.OrderId,
		VariantId:     in.VariantId,
		ProductName:   in.ProductInfo.ProductName,
		IsOutsideEtop: in.IsOutside,
		Quantity:      in.Quantity,
		ImageUrl:      in.ProductInfo.ImageUrl,
		Attributes:    Convert_core_Attributes_To_api_Attributes(in.ProductInfo.Attributes),
		ProductId:     in.ProductId,
	}
}

func Convert_core_OrderLines_To_api_OrderLines(ins []*ordertypes.ItemLine) (outs []*order.OrderLine) {
	for _, in := range ins {
		outs = append(outs, Convert_core_OrderLine_To_api_OrderLine(in))
	}
	return
}

func Convert_api_Attributes_To_core_Attributes(ins []*order.Attribute) []*catalogtypes.Attribute {
	res := make([]*catalogtypes.Attribute, len(ins))
	for i, in := range ins {
		res[i] = &catalogtypes.Attribute{
			Name:  in.Name,
			Value: in.Value,
		}
	}
	return res
}

func Convert_core_Attributes_To_api_Attributes(ins []*catalogtypes.Attribute) (outs []*order.Attribute) {
	for _, in := range ins {
		outs = append(outs, &order.Attribute{
			Name:  in.Name,
			Value: in.Value,
		})
	}
	return
}

func Convert_api_OrderAddress_To_core_OrderAddress(in *order.OrderAddress) *ordertypes.Address {
	if in == nil {
		return nil
	}
	return &ordertypes.Address{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Company:  in.Company,
		Address1: in.Address1,
		Address2: in.Address2,
		Location: ordertypes.Location{
			ProvinceCode: in.ProvinceCode,
			DistrictCode: in.DistrictCode,
			WardCode:     in.WardCode,
			Coordinates:  Convert_api_Coordinates_To_core_Coordinates(in.Coordinates),
		},
	}
}

func Convert_core_OrderAddress_To_api_OrderAddress(in *ordertypes.Address) *order.OrderAddress {
	if in == nil {
		return nil
	}
	return &order.OrderAddress{
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

func Convert_api_Coordinates_To_core_Coordinates(in *etop.Coordinates) *ordertypes.Coordinates {
	if in == nil {
		return nil
	}
	return &ordertypes.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func Convert_core_Coordinates_To_api_Coordinates(in *ordertypes.Coordinates) *etop.Coordinates {
	if in == nil {
		return nil
	}
	return &etop.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func Convert_core_ShipnowService_To_api_ShipnowService(in *shipnowtypes.ShipnowService) *order.ShippnowService {
	if in == nil {
		return nil
	}
	return &order.ShippnowService{
		Carrier:     carrier.CarrierToString(in.Carrier),
		Name:        in.Name,
		Code:        in.Code,
		Fee:         in.Fee,
		Description: in.Description,
	}
}

func Convert_core_ShipnowServices_To_api_ShipnowServices(ins []*shipnowtypes.ShipnowService) (outs []*order.ShippnowService) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowService_To_api_ShipnowService(in))
	}
	return
}