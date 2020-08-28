package convertpb

import (
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	etop "o.o/api/top/int/etop"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in *shipnow.ShipnowFulfillment) *types.ShipnowFulfillment {
	if in == nil {
		return nil
	}
	return &types.ShipnowFulfillment{
		Id:                         in.ID,
		ShopId:                     in.ShopID,
		PartnerId:                  in.PartnerID,
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
		Status:                     in.Status,
		ShippingStatus:             in.ShippingStatus,
		ShippingState:              in.ShippingState,
		ConfirmStatus:              in.ConfirmStatus,
		OrderIds:                   in.OrderIDs,
		ShippingCreatedAt:          cmapi.PbTime(in.ShippingCreatedAt),
		ShippingCode:               in.ShippingCode,
		EtopPaymentStatus:          in.EtopPaymentStatus,
		CodEtopTransferedAt:        cmapi.PbTime(in.CODEtopTransferedAt),
		ShippingPickingAt:          cmapi.PbTime(in.ShippingPickingAt),
		ShippingDeliveringAt:       cmapi.PbTime(in.ShippingDeliveringAt),
		ShippingDeliveredAt:        cmapi.PbTime(in.ShippingDeliveredAt),
		ShippingCancelledAt:        cmapi.PbTime(in.ShippingCancelledAt),
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
		ConnectionID:               in.ConnectionID,
		Coupon:                     in.Coupon,
	}
}

func Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(ins []*shipnow.ShipnowFulfillment) (outs []*types.ShipnowFulfillment) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in))
	}
	return
}

func Convert_api_DeliveryPoint_To_core_DeliveryPoint(in *types.DeliveryPoint) *shipnow.DeliveryPoint {
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

func Conver_api_DeliveryPoints_To_core_DeliveryPoints(ins []*types.DeliveryPoint) (outs []*shipnow.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_api_DeliveryPoint_To_core_DeliveryPoint(in))
	}
	return
}

func Convert_core_DeliveryPoint_To_api_DeliveryPoint(in *shipnow.DeliveryPoint) *types.DeliveryPoint {
	return &types.DeliveryPoint{
		ShippingAddress: Convert_core_OrderAddress_To_api_OrderAddress(in.ShippingAddress),
		Lines:           Convert_core_OrderLines_To_api_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Convert_core_DeliveryPoints_To_api_DeliveryPoints(ins []*shipnow.DeliveryPoint) (outs []*types.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_core_DeliveryPoint_To_api_DeliveryPoint(in))
	}
	return
}

func Convert_api_WeightInfo_To_core_WeightInfo(in types.WeightInfo) shippingtypes.WeightInfo {
	return shippingtypes.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_core_WeightInfo_To_api_WeightInfo(in shippingtypes.WeightInfo) types.WeightInfo {
	return types.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_api_ValueInfo_To_core_ValueInfo(in types.ValueInfo) shippingtypes.ValueInfo {
	return shippingtypes.ValueInfo{
		BasketValue:      in.BasketValue,
		CODAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_core_ValueInfo_To_api_ValueInfo(in shippingtypes.ValueInfo) types.ValueInfo {
	return types.ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CODAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_api_OrderLine_To_core_OrderLine(in *types.OrderLine) (out *ordertypes.ItemLine) {
	return &ordertypes.ItemLine{
		OrderID:   in.OrderId,
		ProductID: in.ProductId,
		Quantity:  in.Quantity,
		VariantID: in.VariantId,
		IsOutSide: in.IsOutsideEtop,
		ProductInfo: ordertypes.ProductInfo{
			ProductName:  in.ProductName,
			ImageURL:     in.ImageUrl,
			Attributes:   in.Attributes,
			ListPrice:    in.ListPrice,
			RetailPrice:  in.RetailPrice,
			PaymentPrice: in.PaymentPrice,
		},
	}
}

func Convert_api_OrderLines_To_core_OrderLines(ins []*types.OrderLine) (outs []*ordertypes.ItemLine) {
	res := make([]*ordertypes.ItemLine, len(ins))
	for i, in := range ins {
		res[i] = Convert_api_OrderLine_To_core_OrderLine(in)
	}
	return res
}

func Convert_core_OrderLine_To_api_OrderLine(in *ordertypes.ItemLine) *types.OrderLine {
	return &types.OrderLine{
		OrderId:       in.OrderID,
		VariantId:     in.VariantID,
		ProductName:   in.ProductInfo.ProductName,
		IsOutsideEtop: in.IsOutSide,
		Quantity:      in.Quantity,
		ListPrice:     in.ProductInfo.ListPrice,
		RetailPrice:   in.ProductInfo.RetailPrice,
		PaymentPrice:  in.ProductInfo.PaymentPrice,
		ImageUrl:      in.ProductInfo.ImageURL,
		Attributes:    in.ProductInfo.Attributes,
		ProductId:     in.ProductID,
		TotalDiscount: 0,
		MetaFields:    nil,
		Code:          "",
	}
}

func Convert_core_OrderLines_To_api_OrderLines(ins []*ordertypes.ItemLine) (outs []*types.OrderLine) {
	for _, in := range ins {
		outs = append(outs, Convert_core_OrderLine_To_api_OrderLine(in))
	}
	return
}

func Convert_api_OrderAddress_To_core_OrderAddress(in *types.OrderAddress) *ordertypes.Address {
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
			Province:     in.Province,
			DistrictCode: in.DistrictCode,
			District:     in.District,
			WardCode:     in.WardCode,
			Ward:         in.Ward,
			Coordinates:  Convert_api_Coordinates_To_core_Coordinates(in.Coordinates),
		},
	}
}

func Convert_core_OrderAddress_To_api_OrderAddress(in *ordertypes.Address) *types.OrderAddress {
	if in == nil {
		return nil
	}
	return &types.OrderAddress{
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Address1:     in.Address1,
		Address2:     in.Address2,
		ProvinceCode: in.Location.ProvinceCode,
		Province:     in.Province,
		DistrictCode: in.Location.DistrictCode,
		District:     in.District,
		WardCode:     in.Location.WardCode,
		Ward:         in.Ward,
		Coordinates:  Convert_core_Coordinates_To_api_Coordinates(in.Coordinates),
	}
}

func Convert_core_OrderAddress_To_api_Address(in *ordertypes.Address) *etop.Address {
	if in == nil {
		return nil
	}
	return &etop.Address{
		Province:     in.Province,
		ProvinceCode: in.ProvinceCode,
		District:     in.District,
		DistrictCode: in.DistrictCode,
		Ward:         in.Ward,
		WardCode:     in.WardCode,
		Address1:     in.Address1,
		Address2:     in.Address2,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
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

func Convert_core_ShipnowService_To_api_ShipnowService(in *shipnowtypes.ShipnowService) *types.ShipnowService {
	if in == nil {
		return nil
	}
	return &types.ShipnowService{
		Carrier:        in.Carrier.String(),
		Name:           in.Name,
		Code:           in.Code,
		Fee:            in.Fee,
		Description:    in.Description,
		ConnectionInfo: Convert_core_ConnectionInfo_To_api_ConnectionInfo(in.ConnectionInfo),
	}
}

func Convert_core_ShipnowServices_To_api_ShipnowServices(ins []*shipnowtypes.ShipnowService) (outs []*types.ShipnowService) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowService_To_api_ShipnowService(in))
	}
	return
}
