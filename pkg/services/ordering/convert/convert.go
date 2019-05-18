package convert

import (
	"time"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/ordering/v1/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	catalogconvert "etop.vn/backend/pkg/services/catalog/convert"
	orderingmodel "etop.vn/backend/pkg/services/ordering/model"
)

func AddressToModel(in *types.Address) (out *orderingmodel.OrderAddress) {
	if in == nil {
		return nil
	}
	out = &orderingmodel.OrderAddress{
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
		Coordinates:  CoordinatesToModel(in.Coordinates),
	}
	return out
}

func Address(in *orderingmodel.OrderAddress) (out *types.Address) {
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

func Order(in *orderingmodel.Order) (out *ordering.Order) {
	if in == nil {
		return nil
	}
	out = &ordering.Order{
		ID:                        in.ID,
		ShopID:                    in.ShopID,
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
	}
	return out
}

func Orders(ins []*orderingmodel.Order) []*ordering.Order {
	result := make([]*ordering.Order, len(ins))
	for i, in := range ins {
		result[i] = Order(in)
	}
	return result
}

func OrderLineToModel(in *types.ItemLine) (out *orderingmodel.OrderLine) {
	if in == nil {
		return nil
	}
	return &orderingmodel.OrderLine{
		OrderID:          in.OrderId,
		VariantID:        in.VariantId,
		ProductName:      in.ProductInfo.ProductName,
		ProductID:        in.ProductId,
		ShopID:           0,
		UpdatedAt:        time.Time{},
		ClosedAt:         time.Time{},
		ConfirmedAt:      time.Time{},
		CancelledAt:      time.Time{},
		CancelReason:     "",
		Status:           0,
		Weight:           0,
		Quantity:         0,
		WholesalePrice0:  0,
		WholesalePrice:   0,
		ListPrice:        0,
		RetailPrice:      0,
		PaymentPrice:     0,
		LineAmount:       0,
		TotalDiscount:    0,
		TotalLineAmount:  0,
		RequiresShipping: false,
		ImageURL:         in.ProductInfo.ImageUrl,
		Attributes:       catalogconvert.AttributesToModel(in.ProductInfo.Attributes),
		IsOutsideEtop:    in.IsOutside,
		Code:             "",
	}
}

func OrderLinesToModel(ins []*types.ItemLine) (outs []*orderingmodel.OrderLine) {
	for _, in := range ins {
		outs = append(outs, OrderLineToModel(in))
	}
	return outs
}

func OrderLine(in *orderingmodel.OrderLine) (out *types.ItemLine) {
	if in == nil {
		return nil
	}
	return &types.ItemLine{
		OrderId:   in.OrderID,
		ProductId: in.ProductID,
		VariantId: in.VariantID,
		IsOutside: in.IsOutsideEtop,
		ProductInfo: types.ProductInfo{
			ProductName: in.ProductName,
			ImageUrl:    in.ImageURL,
			Attributes:  catalogconvert.Attributes(in.Attributes),
		},
	}
}

func OrderLines(ins []*orderingmodel.OrderLine) (outs []*types.ItemLine) {
	for _, in := range ins {
		outs = append(outs, OrderLine(in))
	}
	return outs
}

func FeeLines(ins []orderingmodel.OrderFeeLine) (outs []ordering.OrderFeeLine) {
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

func OrderToShippingInfo(in *orderingmodel.Order) (out *shippingtypes.ShippingInfo) {
	if in == nil || in.ShopShipping == nil {
		return nil
	}
	shopShipping := in.ShopShipping
	tryOn, _ := shippingtypes.TryOnFromString(string(in.TryOn))

	return &shippingtypes.ShippingInfo{
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

func Coordinates(in *orderingmodel.Coordinates) (out *types.Coordinates) {
	if in == nil {
		return nil
	}
	return &types.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func CoordinatesToModel(in *types.Coordinates) (out *orderingmodel.Coordinates) {
	if in == nil {
		return nil
	}
	return &orderingmodel.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}
