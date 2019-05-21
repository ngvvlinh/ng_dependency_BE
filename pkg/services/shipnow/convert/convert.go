package convert

import (
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipping/v1/types"
	"etop.vn/api/meta"
	etopmodel "etop.vn/backend/pkg/etop/model"
	orderconvert "etop.vn/backend/pkg/services/ordering/convert"
	"etop.vn/backend/pkg/services/shipnow/model"
)

func ShipnowToModel(in *shipnow.ShipnowFulfillment) (out *model.ShipnowFulfillment) {
	out = &model.ShipnowFulfillment{
		ID:                  in.Id,
		ShopID:              in.ShopId,
		PartnerID:           in.PartnerId,
		OrderIDs:            nil, // TODO
		PickupAddress:       orderconvert.AddressToModel(in.PickupAddress),
		Carrier:             etopmodel.ShippingProvider(in.Carrier),
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		ChargeableWeight:    in.ChargeableWeight,
		BasketValue:         in.ValueInfo.BasketValue,
		CODAmount:           in.ValueInfo.CodAmount,
		ShippingNote:        in.ShippingNote,
		RequestPickupAt:     in.RequestPickupAt.ToTime(),
		DeliveryPoints:      nil, // TODO
	}
	return out
}

func Shipnow(in *model.ShipnowFulfillment) (out *shipnow.ShipnowFulfillment) {
	out = &shipnow.ShipnowFulfillment{
		Id:                  in.ID,
		ShopId:              in.ShopID,
		PartnerId:           in.PartnerID,
		PickupAddress:       orderconvert.Address(in.PickupAddress),
		DeliveryPoints:      nil, // TODO
		Carrier:             string(in.Carrier),
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		WeightInfo: types.WeightInfo{
			GrossWeight:      in.GrossWeight,
			ChargeableWeight: in.ChargeableWeight,
			Length:           0,
			Width:            0,
			Height:           0,
		},
		ValueInfo: types.ValueInfo{
			BasketValue:      in.BasketValue,
			CodAmount:        in.CODAmount,
			IncludeInsurance: false,
		},
		ShippingNote:    in.ShippingNote,
		RequestPickupAt: meta.PbTime(in.RequestPickupAt),
	}
	return out
}

func DeliveryPoints(ins []*model.DeliveryPoint) (outs []*shipnow.DeliveryPoint) {
	// TODO
	return nil
}
