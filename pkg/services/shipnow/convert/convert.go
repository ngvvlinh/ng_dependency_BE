package convert

import (
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/etop/model"
	etopconvert "etop.vn/backend/pkg/services/etop/convert"
	orderconvert "etop.vn/backend/pkg/services/ordering/convert"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
)

func CarrierToModel(in carrier.Carrier) shipnowmodel.Carrier {
	str := carrier.CarrierToString(in)
	res := shipnowmodel.Carrier(str)
	return res
}

func ShipnowToModel(in *shipnow.ShipnowFulfillment) (out *shipnowmodel.ShipnowFulfillment) {
	out = &shipnowmodel.ShipnowFulfillment{
		ID:                   in.Id,
		ShopID:               in.ShopId,
		PartnerID:            in.PartnerId,
		OrderIDs:             in.OrderIds,
		PickupAddress:        orderconvert.AddressToModel(in.PickupAddress),
		Carrier:              CarrierToModel(in.Carrier),
		ShippingServiceCode:  in.ShippingServiceCode,
		ShippingServiceFee:   in.ShippingServiceFee,
		ChargeableWeight:     in.ChargeableWeight,
		BasketValue:          in.ValueInfo.BasketValue,
		CODAmount:            in.ValueInfo.CodAmount,
		ShippingNote:         in.ShippingNote,
		RequestPickupAt:      in.RequestPickupAt.ToTime(),
		DeliveryPoints:       DeliveryPointsToModel(in.DeliveryPoints),
		ConfirmStatus:        model.Status3(in.ConfirmStatus),
		Status:               model.Status5(in.Status),
		ShippingState:        in.ShippingState.String(),
		ShippingCode:         in.ShippingCode,
		ShippingCreatedAt:    in.ShippingCreatedAt.ToTime(),
		ShippingPickingAt:    in.ShippingPickingAt.ToTime(),
		ShippingDeliveringAt: in.ShippingDeliveringAt.ToTime(),
		ShippingDeliveredAt:  in.ShippingDeliveredAt.ToTime(),
		ShippingCancelledAt:  in.ShippingCancelledAt.ToTime(),
		CreatedAt:            in.CreatedAt.ToTime(),
		UpdatedAt:            in.UpdatedAt.ToTime(),
		CODEtopTransferedAt:  in.CodEtopTransferedAt.ToTime(),
		EtopPaymentStatus:    model.Status4(in.EtopPaymentStatus),
		ShippingServiceName:  in.ShippingServiceName,
		ShippingSharedLink:   in.ShippingSharedLink,
		CancelReason:         in.CancelReason,
	}
	var orderIDs []int64
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
	out = &shipnow.ShipnowFulfillment{
		Id:                  in.ID,
		ShopId:              in.ShopID,
		PartnerId:           in.PartnerID,
		PickupAddress:       orderconvert.Address(in.PickupAddress),
		DeliveryPoints:      DeliveryPoints(in.DeliveryPoints),
		Carrier:             carrier.CarrierFromString(string(in.Carrier)),
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
			CodAmount:        in.CODAmount,
			IncludeInsurance: false,
		},
		ShippingNote:         in.ShippingNote,
		RequestPickupAt:      meta.PbTime(in.RequestPickupAt),
		ConfirmStatus:        etoptypes.Status3FromInt(int(in.ConfirmStatus)),
		Status:               etoptypes.Status5FromInt(int(in.Status)),
		ShippingState:        shipnowtypes.StateFromString(string(in.ShippingState)),
		ShippingCode:         in.ShippingCode,
		OrderIds:             in.OrderIDs,
		ShippingCreatedAt:    meta.PbTime(in.ShippingCreatedAt),
		CreatedAt:            meta.PbTime(in.CreatedAt),
		UpdatedAt:            meta.PbTime(in.UpdatedAt),
		ShippingPickingAt:    meta.PbTime(in.ShippingPickingAt),
		ShippingDeliveringAt: meta.PbTime(in.ShippingDeliveringAt),
		ShippingDeliveredAt:  meta.PbTime(in.ShippingDeliveredAt),
		ShippingCancelledAt:  meta.PbTime(in.ShippingCancelledAt),
		CodEtopTransferedAt:  meta.PbTime(in.CODEtopTransferedAt),
		EtopPaymentStatus:    etoptypes.Status4FromInt(int(in.EtopPaymentStatus)),
		ShippingServiceName:  in.ShippingServiceName,
		ShippingSharedLink:   in.ShippingSharedLink,
		CancelReason:         in.CancelReason,
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
			CodAmount:        in.CODAmount,
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
		ShippingAddress:  orderconvert.AddressToModel(in.ShippingAddress),
		Items:            orderconvert.OrderLinesToModel(in.Lines),
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           in.Length,
		Width:            in.Width,
		Height:           in.Height,
		BasketValue:      in.BasketValue,
		CODAmount:        in.CodAmount,
		TryOn:            "",
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
	tryOn := shippingtypes.TryOnNone
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
			GrossWeight:      int32(grossWeight),
			ChargeableWeight: int32(chargeableWeight),
			Length:           int32(lenght),
			Width:            0,
			Height:           int32(height),
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      int32(in.BasketValue),
			CodAmount:        int32(codAmount),
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
		GrossWeight:      int32(grossWeight),
		ChargeableWeight: int32(chargeableWeight),
	}
}

func GetValueInfo(orders []*ordering.Order) shippingtypes.ValueInfo {
	basketValue, codAmount := 0, 0
	for _, o := range orders {
		basketValue += o.BasketValue
		codAmount += o.Shipping.CODAmount
	}
	return shippingtypes.ValueInfo{
		BasketValue: int32(basketValue),
		CodAmount:   int32(codAmount),
	}
}

func FeelineToModel(in *shippingtypes.FeeLine) (out *model.ShippingFeeLine) {
	if in == nil {
		return nil
	}
	return &model.ShippingFeeLine{
		ShippingFeeType:     model.ShippingFeeLineType(in.ShippingFeeType.String()),
		Cost:                int(in.Cost),
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
		ShippingFeeType:     shippingtypes.FeelineTypeFromString(string(in.ShippingFeeType)),
		Cost:                int32(in.Cost),
		ExternalServiceName: in.ExternalServiceName,
		ExternalServiceType: in.ExternalServiceType,
	}
}

func SyncStateToModel(in *shipnow.SyncStates) *model.FulfillmentSyncStates {
	if in == nil {
		return nil
	}
	return &model.FulfillmentSyncStates{
		SyncAt:    in.SyncAt.ToTime(),
		TrySyncAt: in.TrySyncAt.ToTime(),
		Error:     etopconvert.ErrorToModel(in.Error),
	}
}
