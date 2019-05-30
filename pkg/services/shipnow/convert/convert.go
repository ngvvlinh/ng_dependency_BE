package convert

import (
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/etop/model"
	orderconvert "etop.vn/backend/pkg/services/ordering/convert"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
)

func ShipnowToModel(in *shipnow.ShipnowFulfillment) (out *shipnowmodel.ShipnowFulfillment) {
	out = &shipnowmodel.ShipnowFulfillment{
		ID:                  in.Id,
		ShopID:              in.ShopId,
		PartnerID:           in.PartnerId,
		OrderIDs:            in.OrderIds,
		PickupAddress:       orderconvert.AddressToModel(in.PickupAddress),
		Carrier:             shipnowmodel.Carrier(in.Carrier),
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		ChargeableWeight:    in.ChargeableWeight,
		BasketValue:         in.ValueInfo.BasketValue,
		CODAmount:           in.ValueInfo.CodAmount,
		ShippingNote:        in.ShippingNote,
		RequestPickupAt:     in.RequestPickupAt.ToTime(),
		DeliveryPoints:      DeliveryPointsToModel(in.DeliveryPoints),
		ConfirmStatus:       model.Status3(in.ConfirmStatus),
		Status:              model.Status5(in.Status),
		ShippingState:       in.ShippingState.String(),
		ShippingCode:        in.ShippingCode,
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
		Carrier:             string(in.Carrier),
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
		ShippingNote:    in.ShippingNote,
		RequestPickupAt: meta.PbTime(in.RequestPickupAt),
		ConfirmStatus:   etoptypes.Status3FromInt(int(in.ConfirmStatus)),
		Status:          etoptypes.Status5FromInt(int(in.Status)),
		ShippingState:   shipnowtypes.StateFromString(string(in.ShippingState)),
		ShippingCode:    in.ShippingCode,
		OrderIds:        in.OrderIDs,
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
	}
}

func DeliveryPointsToModel(ins []*shipnowtypes.DeliveryPoint) (outs []*shipnowmodel.DeliveryPoint) {
	for _, in := range ins {
		outs = append(outs, DeliveryPointToModel(in))
	}
	return outs
}

func OrderToDeliveryPoint(in *ordering.Order) *shipnowtypes.DeliveryPoint {
	return &shipnowtypes.DeliveryPoint{
		ShippingAddress: in.ShippingAddress,
		Lines:           in.Lines,
		ShippingNote:    in.OrderNote,
		OrderId:         in.ID,
		WeightInfo: shippingtypes.WeightInfo{
			GrossWeight:      int32(in.Shipping.GrossWeight),
			ChargeableWeight: int32(in.Shipping.ChargeableWeight),
			Length:           int32(in.Shipping.Length),
			Width:            0,
			Height:           int32(in.Shipping.Height),
		},
		ValueInfo: shippingtypes.ValueInfo{
			BasketValue:      int32(in.BasketValue),
			CodAmount:        int32(in.Shipping.CODAmount),
			IncludeInsurance: in.Shipping.IncludeInsurance,
		},
		TryOn: in.Shipping.TryOn,
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
