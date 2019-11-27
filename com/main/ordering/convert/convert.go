package convert

import (
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/ordering/types"
	ordertypes "etop.vn/api/main/ordering/types"
	shiptypes "etop.vn/api/main/shipping/types"
	catalogconvert "etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/ordering/model"
	"etop.vn/common/jsonx"
)

func AddressDB(in *types.Address) (out *model.OrderAddress) {
	if in == nil {
		return nil
	}
	out = &model.OrderAddress{
		FullName:     in.FullName,
		FirstName:    "",
		LastName:     "",
		Phone:        in.Phone,
		Email:        in.Email,
		Country:      "",
		City:         "",
		Province:     "",
		District:     "",
		Ward:         "",
		Zip:          "",
		DistrictCode: in.DistrictCode,
		ProvinceCode: in.ProvinceCode,
		WardCode:     in.WardCode,
		Company:      "",
		Address1:     in.Address1,
		Address2:     in.Address2,
		Coordinates:  CoordinatesDB(in.Coordinates),
	}
	return out
}

func Address(in *model.OrderAddress) (out *types.Address) {
	if in == nil {
		return nil
	}
	out = &types.Address{
		FullName: in.FullName,
		Phone:    in.Phone,
		Email:    in.Email,
		Company:  in.Company,
		Address1: in.Address1,
		Address2: in.Address2,
		Location: types.Location{
			ProvinceCode: in.ProvinceCode,
			DistrictCode: in.DistrictCode,
			WardCode:     in.WardCode,
			Coordinates:  Coordinates(in.Coordinates),
		},
	}
	return out
}

func Order(in *model.Order) (out *ordering.Order) {
	if in == nil {
		return nil
	}
	out = &ordering.Order{
		ID:                        in.ID,
		ShopID:                    in.ShopID,
		Code:                      in.Code,
		CustomerAddress:           Address(in.CustomerAddress),
		ShippingAddress:           Address(in.ShippingAddress),
		CancelReason:              in.CancelReason,
		ConfirmStatus:             etoptypes.Status3FromInt(int(in.ConfirmStatus)),
		Status:                    etoptypes.Status5FromInt(int(in.Status)),
		FulfillmentShippingStatus: etoptypes.Status5FromInt(int(in.FulfillmentShippingStatus)),
		EtopPaymentStatus:         etoptypes.Status4FromInt(int(in.EtopPaymentStatus)),
		Lines:                     OrderLines(in.Lines),
		TotalItems:                in.TotalItems,
		BasketValue:               in.BasketValue,
		OrderDiscount:             in.OrderDiscount,
		TotalDiscount:             in.TotalDiscount,
		TotalFee:                  in.TotalFee,
		TotalAmount:               in.TotalAmount,
		OrderNote:                 in.OrderNote,
		FeeLines:                  FeeLines(in.FeeLines),
		Shipping:                  OrderToShippingInfo(in),
		CreatedAt:                 in.CreatedAt,
		UpdatedAt:                 in.UpdatedAt,
		ProcessedAt:               in.ProcessedAt,
		ClosedAt:                  in.ClosedAt,
		ConfirmedAt:               in.ConfirmedAt,
		CancelledAt:               in.CancelledAt,
		FulfillmentIDs:            in.FulfillmentIDs,
		FulfillmentType:           ordertypes.Fulfill(in.FulfillmentType),
		PaymentStatus:             etoptypes.Status4FromInt(int(in.PaymentStatus)),
		PaymentID:                 in.PaymentID,
		CustomerID:                in.CustomerID,
		TradingShopID:             in.TradingShopID,
	}
	var referralMeta ordering.ReferralMeta
	if err := jsonx.Unmarshal(in.ReferralMeta, &referralMeta); err == nil {
		out.ReferralMeta = &referralMeta
	}
	return out
}

func Orders(ins []*model.Order) []*ordering.Order {
	result := make([]*ordering.Order, len(ins))
	for i, in := range ins {
		result[i] = Order(in)
	}
	return result
}

func OrderLineToModel(in *types.ItemLine) (out *model.OrderLine) {
	if in == nil {
		return nil
	}
	return &model.OrderLine{
		OrderID:         in.OrderId,
		VariantID:       in.VariantId,
		ProductName:     in.ProductInfo.ProductName,
		ProductID:       in.ProductId,
		Quantity:        int(in.Quantity),
		TotalLineAmount: int(in.TotalPrice),
		ImageURL:        in.ProductInfo.ImageUrl,
		Attributes:      catalogconvert.AttributesDB(in.ProductInfo.Attributes),
		IsOutsideEtop:   in.IsOutside,
		Code:            "",
	}
}

func OrderLinesToModel(ins []*types.ItemLine) (outs []*model.OrderLine) {
	for _, in := range ins {
		outs = append(outs, OrderLineToModel(in))
	}
	return outs
}

func OrderLine(in *model.OrderLine) (out *types.ItemLine) {
	if in == nil {
		return nil
	}
	return &types.ItemLine{
		OrderId:   in.OrderID,
		Quantity:  int(in.Quantity),
		ProductId: in.ProductID,
		VariantId: in.VariantID,
		IsOutside: in.IsOutsideEtop,
		ProductInfo: types.ProductInfo{
			ProductName: in.ProductName,
			ImageUrl:    in.ImageURL,
			Attributes:  catalogconvert.Attributes(in.Attributes),
		},
		TotalPrice: int(in.TotalLineAmount),
	}
}

func OrderLines(ins []*model.OrderLine) (outs []*types.ItemLine) {
	for _, in := range ins {
		outs = append(outs, OrderLine(in))
	}
	return outs
}

func FeeLines(ins []model.OrderFeeLine) (outs []ordering.OrderFeeLine) {
	for _, in := range ins {
		outs = append(outs, ordering.OrderFeeLine{
			Type:   string(in.Type),
			Name:   in.Name,
			Code:   in.Code,
			Desc:   in.Desc,
			Amount: in.Amount,
		})
	}
	return
}

func OrderToShippingInfo(in *model.Order) (out *shiptypes.ShippingInfo) {
	if in == nil || in.ShopShipping == nil {
		return nil
	}
	shopShipping := in.ShopShipping
	tryOn, _ := shiptypes.TryOnFromString(string(in.TryOn))

	return &shiptypes.ShippingInfo{
		PickupAddress:       Address(shopShipping.ShopAddress),
		ReturnAddress:       Address(shopShipping.ReturnAddress),
		ShippingServiceName: shopShipping.ExternalServiceName,
		ShippingServiceCode: shopShipping.ProviderServiceID,
		ShippingServiceFee:  shopShipping.ExternalShippingFee,
		Carrier:             string(shopShipping.ShippingProvider),
		IncludeInsurance:    shopShipping.IncludeInsurance,
		TryOn:               tryOn,
		ShippingNote:        in.ShippingNote,
		CODAmount:           in.ShopCOD,
		GrossWeight:         shopShipping.GrossWeight,
		Length:              shopShipping.Length,
		Height:              shopShipping.Height,
		ChargeableWeight:    shopShipping.ChargeableWeight,
	}
}

func Coordinates(in *model.Coordinates) (out *types.Coordinates) {
	if in == nil {
		return nil
	}
	return &types.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func CoordinatesDB(in *types.Coordinates) (out *model.Coordinates) {
	if in == nil {
		return nil
	}
	return &model.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func Fulfill(in model.FulfillType) ordertypes.Fulfill {
	res, _ := ordertypes.FulfillFromInt(int(in))
	return res
}
