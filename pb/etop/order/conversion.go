package order

import (
	catalogtypes "etop.vn/api/main/catalog/v1/types"
	orderingv1types "etop.vn/api/main/ordering/v1/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/v1/types"
	"etop.vn/api/meta"
	pbcm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbs5 "etop.vn/backend/pb/etop/etc/status5"
	"etop.vn/backend/pkg/etop/model"
)

func Convert_api_ShipnowFulfillment_To_core_ShipnowFulfillment(in *ShipnowFulfillment) *shipnow.ShipnowFulfillment {
	return &shipnow.ShipnowFulfillment{
		Id:                         in.Id,
		ShopId:                     in.ShopId,
		PartnerId:                  in.PartnerId,
		PickupAddress:              Convert_api_OrderAddress_To_core_OrderAddress(in.PickupAddress),
		DeliveryPoints:             Conver_api_DeliveryPoints_To_core_DeliveryPoints(in.DeliveryPoints),
		Carrier:                    carrier.CarrierFromString(in.Carrier),
		ShippingServiceCode:        in.ShippingServiceCode,
		ShippingServiceFee:         in.ShippingServiceFee,
		ShippingServiceName:        in.ShippingServiceName,
		ShippingServiceDescription: in.ShippingServiceDescription,
		WeightInfo:                 Convert_api_WeightInfo_To_core_WeightInfo(in.WeightInfo),
		ValueInfo:                  Convert_api_ValueInfo_To_core_ValueInfo(in.ValueInfo),
		ShippingNote:               in.ShippingNote,
		RequestPickupAt:            meta.PbTime(pbcm.PbTimeToModel(in.RequestPickupAt)),
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
	}
}

func Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in *shipnow.ShipnowFulfillment) *ShipnowFulfillment {
	if in == nil {
		return nil
	}
	return &ShipnowFulfillment{
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
		RequestPickupAt:            pbcm.PbTime(in.RequestPickupAt.ToTime()),
		CreatedAt:                  pbcm.PbTime(in.CreatedAt.ToTime()),
		UpdatedAt:                  pbcm.PbTime(in.UpdatedAt.ToTime()),
		Status:                     pbs5.Pb(model.Status5(in.Status)),
		ShippingStatus:             pbs5.Pb(model.Status5(in.ShippingStatus)),
		ShippingState:              shipnowtypes.StateToString(in.ShippingState),
		ConfirmStatus:              pbs3.Pb(model.Status3(in.ConfirmStatus)),
		OrderIds:                   in.OrderIds,
		ShippingCreatedAt:          pbcm.PbTime(in.ShippingCreatedAt.ToTime()),
		ShippingCode:               in.ShippingCode,
		EtopPaymentStatus:          pbs4.Pb(model.Status4(in.EtopPaymentStatus)),
		CodEtopTransferedAt:        pbcm.PbTime(in.CodEtopTransferedAt.ToTime()),
		ShippingPickingAt:          pbcm.PbTime(in.ShippingPickingAt.ToTime()),
		ShippingDeliveringAt:       pbcm.PbTime(in.ShippingDeliveringAt.ToTime()),
		ShippingDeliveredAt:        pbcm.PbTime(in.ShippingDeliveredAt.ToTime()),
		ShippingCancelledAt:        pbcm.PbTime(in.ShippingCancelledAt.ToTime()),
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
	}
}

func Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(ins []*shipnow.ShipnowFulfillment) (outs []*ShipnowFulfillment) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(in))
	}
	return
}

func Convert_api_DeliveryPoint_To_core_DeliveryPoint(in *DeliveryPoint) *shipnow.DeliveryPoint {
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

func Conver_api_DeliveryPoints_To_core_DeliveryPoints(ins []*DeliveryPoint) (outs []*shipnow.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_api_DeliveryPoint_To_core_DeliveryPoint(in))
	}
	return
}

func Convert_core_DeliveryPoint_To_api_DeliveryPoint(in *shipnow.DeliveryPoint) *DeliveryPoint {
	return &DeliveryPoint{
		ShippingAddress: Convert_core_OrderAddress_To_api_OrderAddress(in.ShippingAddress),
		Lines:           Convert_core_OrderLines_To_api_OrderLines(in.Lines),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderId,
		WeightInfo:      Convert_core_WeightInfo_To_api_WeightInfo(in.WeightInfo),
		ValueInfo:       Convert_core_ValueInfo_To_api_ValueInfo(in.ValueInfo),
		TryOn:           0,
	}
}

func Convert_core_DeliveryPoints_To_api_DeliveryPoints(ins []*shipnow.DeliveryPoint) (outs []*DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, Convert_core_DeliveryPoint_To_api_DeliveryPoint(in))
	}
	return
}

func Convert_api_WeightInfo_To_core_WeightInfo(in WeightInfo) shippingtypes.WeightInfo {
	return shippingtypes.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_core_WeightInfo_To_api_WeightInfo(in shippingtypes.WeightInfo) WeightInfo {
	return WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
	}
}

func Convert_api_ValueInfo_To_core_ValueInfo(in ValueInfo) shippingtypes.ValueInfo {
	return shippingtypes.ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_core_ValueInfo_To_api_ValueInfo(in shippingtypes.ValueInfo) ValueInfo {
	return ValueInfo{
		BasketValue:      in.BasketValue,
		CodAmount:        in.CodAmount,
		IncludeInsurance: in.IncludeInsurance,
	}
}

func Convert_api_OrderLine_To_core_OrderLine(in *OrderLine) (out *orderingv1types.ItemLine) {
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

func Convert_api_OrderLines_To_core_OrderLines(ins []*OrderLine) (outs []*orderingv1types.ItemLine) {
	res := make([]*orderingv1types.ItemLine, len(ins))
	for i, in := range ins {
		res[i] = Convert_api_OrderLine_To_core_OrderLine(in)
	}
	return res
}

func Convert_core_OrderLine_To_api_OrderLine(in *orderingv1types.ItemLine) *OrderLine {
	return &OrderLine{
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

func Convert_core_OrderLines_To_api_OrderLines(ins []*orderingv1types.ItemLine) (outs []*OrderLine) {
	for _, in := range ins {
		outs = append(outs, Convert_core_OrderLine_To_api_OrderLine(in))
	}
	return
}

func Convert_api_Attributes_To_core_Atributes(ins []*Attribute) []*catalogtypes.Attribute {
	res := make([]*catalogtypes.Attribute, len(ins))
	for i, in := range ins {
		res[i] = &catalogtypes.Attribute{
			Name:  in.Name,
			Value: in.Value,
		}
	}
	return res
}

func Convert_core_Attribute_To_api_Atribures(ins []*catalogtypes.Attribute) (outs []*Attribute) {
	for _, in := range ins {
		outs = append(outs, &Attribute{
			Name:  in.Name,
			Value: in.Value,
		})
	}
	return
}

func Convert_api_OrderAddress_To_core_OrderAddress(in *OrderAddress) *orderingv1types.Address {
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

func Convert_core_OrderAddress_To_api_OrderAddress(in *orderingv1types.Address) *OrderAddress {
	if in == nil {
		return nil
	}
	return &OrderAddress{
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

func Convert_api_Coordinates_To_core_Coordinates(in *etop.Coordinates) *orderingv1types.Coordinates {
	if in == nil {
		return nil
	}
	return &orderingv1types.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func Convert_core_Coordinates_To_api_Coordinates(in *orderingv1types.Coordinates) *etop.Coordinates {
	if in == nil {
		return nil
	}
	return &etop.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func Convert_core_ShipnowService_To_api_ShipnowService(in *shipnowtypes.ShipnowService) *ShippnowService {
	if in == nil {
		return nil
	}
	return &ShippnowService{
		Carrier:     string(carrier.CarrierToString(in.Carrier)),
		Name:        in.Name,
		Code:        in.Code,
		Fee:         in.Fee,
		Description: in.Description,
	}
}

func Convert_core_ShipnowServices_To_api_ShipnowServices(ins []*shipnowtypes.ShipnowService) (outs []*ShippnowService) {
	for _, in := range ins {
		outs = append(outs, Convert_core_ShipnowService_To_api_ShipnowService(in))
	}
	return
}
