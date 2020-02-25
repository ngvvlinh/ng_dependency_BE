package convert

import (
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/ordering/types"
	shiptypes "etop.vn/api/main/shipping/types"
	addressconvert "etop.vn/backend/com/main/address/convert"
	catalogconvert "etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/ordering/model"
	"etop.vn/common/jsonx"
)

// +gen:convert: etop.vn/backend/com/main/ordering/model->etop.vn/api/main/ordering,etop.vn/api/main/ordering/types
// +gen:convert: etop.vn/api/main/ordering

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
		Province:     in.Province,
		District:     in.District,
		Ward:         in.Ward,
		Zip:          "",
		DistrictCode: in.DistrictCode,
		ProvinceCode: in.ProvinceCode,
		WardCode:     in.WardCode,
		Company:      "",
		Address1:     in.Address1,
		Address2:     in.Address2,
		Coordinates:  addressconvert.Convert_orderingtypes_Coordinates_addressmodel_Coordinates(in.Coordinates, nil),
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
			Province:     in.Province,
			DistrictCode: in.DistrictCode,
			District:     in.District,
			WardCode:     in.WardCode,
			Ward:         in.Ward,
			Coordinates:  addressconvert.Convert_addressmodel_Coordinates_orderingtypes_Coordinates(in.Coordinates, nil),
		},
	}
	return out
}

func Order(in *model.Order) *ordering.Order {
	if in == nil {
		return nil
	}
	out := &ordering.Order{}
	convert_orderingmodel_Order_ordering_Order(in, out)
	out.FeeLines = FeeLines(in.FeeLines)
	out.Shipping = OrderToShippingInfo(in)

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
		OrderID:         in.OrderID,
		VariantID:       in.VariantID,
		ProductName:     in.ProductInfo.ProductName,
		ListPrice:       in.ProductInfo.ListPrice,
		RetailPrice:     in.ProductInfo.RetailPrice,
		PaymentPrice:    in.ProductInfo.PaymentPrice,
		ProductID:       in.ProductID,
		Quantity:        in.Quantity,
		TotalLineAmount: in.TotalPrice,
		ImageURL:        in.ProductInfo.ImageURL,
		Attributes:      catalogconvert.Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(in.ProductInfo.Attributes),
		IsOutsideEtop:   in.IsOutSide,
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
		OrderID:   in.OrderID,
		Quantity:  in.Quantity,
		ProductID: in.ProductID,
		VariantID: in.VariantID,
		IsOutSide: in.IsOutsideEtop,
		ProductInfo: types.ProductInfo{
			ProductName:  in.ProductName,
			ImageURL:     in.ImageURL,
			Attributes:   catalogconvert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(in.Attributes),
			ListPrice:    in.ListPrice,
			RetailPrice:  in.RetailPrice,
			PaymentPrice: in.PaymentPrice,
		},
		TotalPrice: in.TotalLineAmount,
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
			Type:   in.Type,
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

	return &shiptypes.ShippingInfo{
		PickupAddress:       Address(shopShipping.ShopAddress),
		ReturnAddress:       Address(shopShipping.ReturnAddress),
		ShippingServiceName: shopShipping.ExternalServiceName,
		ShippingServiceCode: shopShipping.ProviderServiceID,
		ShippingServiceFee:  shopShipping.ExternalShippingFee,
		Carrier:             shopShipping.ShippingProvider,
		IncludeInsurance:    shopShipping.IncludeInsurance,
		TryOn:               in.TryOn,
		ShippingNote:        in.ShippingNote,
		CODAmount:           in.ShopCOD,
		GrossWeight:         shopShipping.GrossWeight,
		Length:              shopShipping.Length,
		Height:              shopShipping.Height,
		ChargeableWeight:    shopShipping.ChargeableWeight,
	}
}
