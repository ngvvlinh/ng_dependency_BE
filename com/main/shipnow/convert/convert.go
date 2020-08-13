package convert

import (
	"o.o/api/main/ordering"
	"o.o/api/main/shipnow"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/try_on"
	etopconvert "o.o/backend/com/main/etop/convert"
	orderconvert "o.o/backend/com/main/ordering/convert"
	shipnowmodel "o.o/backend/com/main/shipnow/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/capi/dot"
)

// +gen:convert: o.o/backend/com/main/shipnow/model -> o.o/api/main/shipnow
// +gen:convert: o.o/api/main/shipnow

func ShipnowToModel(in *shipnow.ShipnowFulfillment) (out *shipnowmodel.ShipnowFulfillment) {
	if out == nil {
		out = &shipnowmodel.ShipnowFulfillment{}
	}
	convert_shipnow_ShipnowFulfillment_shipnowmodel_ShipnowFulfillment(in, out)
	out.GrossWeight = in.GrossWeight
	out.ChargeableWeight = in.ChargeableWeight
	out.BasketValue = in.ValueInfo.BasketValue
	out.CODAmount = in.ValueInfo.CODAmount
	if len(in.DeliveryPoints) > 0 {
		out.AddressToDistrictCode = in.DeliveryPoints[0].ShippingAddress.DistrictCode
		out.AddressToProvinceCode = in.DeliveryPoints[0].ShippingAddress.ProvinceCode
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
	if out == nil {
		out = &shipnow.ShipnowFulfillment{}
	}
	convert_shipnowmodel_ShipnowFulfillment_shipnow_ShipnowFulfillment(in, out)
	out.WeightInfo = shippingtypes.WeightInfo{
		GrossWeight:      in.GrossWeight,
		ChargeableWeight: in.ChargeableWeight,
		Length:           0,
		Width:            0,
		Height:           0,
	}
	out.ValueInfo = shippingtypes.ValueInfo{
		BasketValue:      in.BasketValue,
		CODAmount:        in.CODAmount,
		IncludeInsurance: false,
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
		TryOn:         0,
		ShippingState: in.ShippingState,
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
		ShippingState:    in.ShippingState,
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

func SyncStateToModel(in *shipnow.SyncStates) *shippingsharemodel.FulfillmentSyncStates {
	if in == nil {
		return nil
	}
	return &shippingsharemodel.FulfillmentSyncStates{
		SyncAt:    in.SyncAt,
		TrySyncAt: in.TrySyncAt,
		Error:     etopconvert.ErrorToModel(in.Error),
	}
}
