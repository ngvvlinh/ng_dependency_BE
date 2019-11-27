package convert

import (
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
	carrier "etop.vn/api/main/shipnow/carrier/types"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/top/types/etc/try_on"
	etopconvert "etop.vn/backend/com/main/etop/convert"
	orderconvert "etop.vn/backend/com/main/ordering/convert"
	shipnowmodel "etop.vn/backend/com/main/shipnow/model"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func CarrierToModel(in carrier.Carrier) shipnowmodel.Carrier {
	res := shipnowmodel.Carrier(in.String())
	return res
}

func ShipnowToModel(in *shipnow.ShipnowFulfillment) (out *shipnowmodel.ShipnowFulfillment) {
	out = &shipnowmodel.ShipnowFulfillment{
		ID:                         in.Id,
		ShopID:                     in.ShopId,
		PartnerID:                  in.PartnerId,
		OrderIDs:                   in.OrderIds,
		PickupAddress:              orderconvert.AddressDB(in.PickupAddress),
		Carrier:                    CarrierToModel(in.Carrier),
		ShippingServiceCode:        in.ShippingServiceCode,
		ShippingServiceFee:         in.ShippingServiceFee,
		ChargeableWeight:           in.ChargeableWeight,
		BasketValue:                in.ValueInfo.BasketValue,
		CODAmount:                  in.ValueInfo.CODAmount,
		ShippingNote:               in.ShippingNote,
		RequestPickupAt:            in.RequestPickupAt,
		DeliveryPoints:             DeliveryPointsToModel(in.DeliveryPoints),
		ConfirmStatus:              in.ConfirmStatus,
		Status:                     in.Status,
		ShippingState:              in.ShippingState,
		ShippingCode:               in.ShippingCode,
		ShippingCreatedAt:          in.ShippingCreatedAt,
		ShippingPickingAt:          in.ShippingPickingAt,
		ShippingDeliveringAt:       in.ShippingDeliveringAt,
		ShippingDeliveredAt:        in.ShippingDeliveredAt,
		ShippingCancelledAt:        in.ShippingCancelledAt,
		CreatedAt:                  in.CreatedAt,
		UpdatedAt:                  in.UpdatedAt,
		CODEtopTransferedAt:        in.CodEtopTransferedAt,
		EtopPaymentStatus:          in.EtopPaymentStatus,
		ShippingServiceName:        in.ShippingServiceName,
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
		ShippingServiceDescription: in.ShippingServiceDescription,
		AddressToDistrictCode:      in.DeliveryPoints[0].ShippingAddress.DistrictCode,
		AddressToProvinceCode:      in.DeliveryPoints[0].ShippingAddress.ProvinceCode,
	}
	var orderIDs []dot.ID
	for _, point := range in.DeliveryPoints {
		if point.OrderId == 0 {
			continue
		}
		orderIDs = append(orderIDs, point.OrderId)
	}
	out.OrderIDs = orderIDs
	return out
}

func Shipnow(in *shipnowmodel.ShipnowFulfillment) (out *shipnow.ShipnowFulfillment) {
	_carrier, _ := carrier.ParseCarrier(in.Carrier.String())
	out = &shipnow.ShipnowFulfillment{
		Id:                  in.ID,
		ShopId:              in.ShopID,
		PartnerId:           in.PartnerID,
		PickupAddress:       orderconvert.Address(in.PickupAddress),
		DeliveryPoints:      DeliveryPoints(in.DeliveryPoints),
		Carrier:             _carrier,
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      in.GrossWeight,
			ChargeableWeight: in.ChargeableWeight,
			Length:           0,
			Width:            0,
			Height:           0,
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      in.BasketValue,
			CODAmount:        in.CODAmount,
			IncludeInsurance: false,
		},
		ShippingNote:               in.ShippingNote,
		RequestPickupAt:            in.RequestPickupAt,
		ConfirmStatus:              in.ConfirmStatus,
		Status:                     in.Status,
		ShippingState:              in.ShippingState,
		ShippingCode:               in.ShippingCode,
		OrderIds:                   in.OrderIDs,
		ShippingCreatedAt:          in.ShippingCreatedAt,
		CreatedAt:                  in.CreatedAt,
		UpdatedAt:                  in.UpdatedAt,
		ShippingPickingAt:          in.ShippingPickingAt,
		ShippingDeliveringAt:       in.ShippingDeliveringAt,
		ShippingDeliveredAt:        in.ShippingDeliveredAt,
		ShippingCancelledAt:        in.ShippingCancelledAt,
		CodEtopTransferedAt:        in.CODEtopTransferedAt,
		EtopPaymentStatus:          in.EtopPaymentStatus,
		ShippingServiceName:        in.ShippingServiceName,
		ShippingSharedLink:         in.ShippingSharedLink,
		CancelReason:               in.CancelReason,
		ShippingServiceDescription: in.ShippingServiceDescription,
	}
	return out
}

func Shipnows(ins []*shipnowmodel.ShipnowFulfillment) (outs []*shipnow.ShipnowFulfillment) {
	for _, in := range ins {
		outs = append(outs, Shipnow(in))
	}
	return
}

func DeliveryPoint(in *shipnowmodel.DeliveryPoint) (outs *shipnowtypes.DeliveryPoint) {
	if in == nil {
		return nil
	}
	return &shipnowtypes.DeliveryPoint{
		ShippingAddress: orderconvert.Address(in.ShippingAddress),
		Lines:           orderconvert.OrderLines(in.Items),
		ShippingNote:    in.ShippingNote,
		OrderId:         in.OrderID,
		OrderCode:       in.OrderCode,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      in.GrossWeight,
			ChargeableWeight: in.ChargeableWeight,
			Length:           in.Length,
			Width:            in.Width,
			Height:           in.Height,
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      in.BasketValue,
			CODAmount:        in.CODAmount,
			IncludeInsurance: false,
		},
		TryOn: 0,
	}
}

func DeliveryPoints(ins []*shipnowmodel.DeliveryPoint) (outs []*shipnowtypes.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, DeliveryPoint(in))
	}
	return outs
}

func DeliveryPointToModel(in *shipnowtypes.DeliveryPoint) (out *shipnowmodel.DeliveryPoint) {
	if in == nil {
		return nil
	}
	return &shipnowmodel.DeliveryPoint{
		ShippingAddress:  orderconvert.AddressDB(in.ShippingAddress),
		Items:            orderconvert.OrderLinesToModel(in.Lines),
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
		BasketValue:      in.BasketValue,
		CODAmount:        in.CODAmount,
		TryOn:            0,
		ShippingNote:     in.ShippingNote,
		OrderID:          in.OrderId,
		OrderCode:        in.OrderCode,
	}
}

func DeliveryPointsToModel(ins []*shipnowtypes.DeliveryPoint) (outs []*shipnowmodel.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, DeliveryPointToModel(in))
	}
	return outs
}

func OrderToDeliveryPoint(in *ordering.Order) *shipnowtypes.DeliveryPoint {
	grossWeight, chargeableWeight, lenght, height, codAmount := 0, 0, 0, 0, 0
	tryOn := try_on.None
	shippingNote := ""
	if in.Shipping != nil {
		grossWeight = in.Shipping.GrossWeight
		chargeableWeight = in.Shipping.ChargeableWeight
		lenght = in.Shipping.Length
		height = in.Shipping.Height
		codAmount = in.Shipping.CODAmount
		tryOn = in.Shipping.TryOn
		shippingNote = in.Shipping.ShippingNote
	}
	return &shipnowtypes.DeliveryPoint{
		ShippingAddress: in.ShippingAddress,
		Lines:           in.Lines,
		ShippingNote:    shippingNote,
		OrderId:         in.ID,
		OrderCode:       in.Code,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      grossWeight,
			ChargeableWeight: chargeableWeight,
			Length:           lenght,
			Width:            0,
			Height:           height,
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      in.BasketValue,
			CODAmount:        codAmount,
			IncludeInsurance: in.Shipping.IncludeInsurance,
		},
		TryOn: tryOn,
	}
}

func GetWeightInfo(orders []*ordering.Order) shippingtypes.WeightInfo {
	grossWeight, chargeableWeight := 0, 0
	for _, o := range orders {
		grossWeight += o.Shipping.GrossWeight
		chargeableWeight += o.Shipping.ChargeableWeight
	}
	return shippingtypes.WeightInfo{
		GrossWeight:      grossWeight,
		ChargeableWeight: chargeableWeight,
	}
}

func GetValueInfo(orders []*ordering.Order) shippingtypes.ValueInfo {
	basketValue, codAmount := 0, 0
	for _, o := range orders {
		basketValue += o.BasketValue
		codAmount += o.Shipping.CODAmount
	}
	return shippingtypes.ValueInfo{
		BasketValue: basketValue,
		CODAmount:   codAmount,
	}
}

func FeelineToModel(in *shippingtypes.FeeLine) (out *model.ShippingFeeLine) {
	if in == nil {
		return nil
	}
	return &model.ShippingFeeLine{
		ShippingFeeType:     in.ShippingFeeType,
		Cost:                in.Cost,
		ExternalServiceName: in.ExternalServiceName,
		ExternalServiceType: in.ExternalServiceType,
	}
}

func FeelinesToModel(ins []*shippingtypes.FeeLine) (outs []*model.ShippingFeeLine) {
	for _, in := range ins {
		outs = append(outs, FeelineToModel(in))
	}
	return
}

func Feeline(in *model.ShippingFeeLine) (out *shippingtypes.FeeLine) {
	if in == nil {
		return nil
	}
	return &shippingtypes.FeeLine{
		ShippingFeeType:     in.ShippingFeeType,
		Cost:                in.Cost,
		ExternalServiceName: in.ExternalServiceName,
		ExternalServiceType: in.ExternalServiceType,
	}
}

func SyncStateToModel(in *shipnow.SyncStates) *model.FulfillmentSyncStates {
	if in == nil {
		return nil
	}
	return &model.FulfillmentSyncStates{
		SyncAt:    in.SyncAt,
		TrySyncAt: in.TrySyncAt,
		Error:     etopconvert.ErrorToModel(in.Error),
	}
}
