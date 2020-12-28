package convertpb

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	shiptypes "o.o/api/main/shipping/types"
	etop "o.o/api/top/int/etop"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/gender"
	addressmodel "o.o/backend/com/main/address/model"
	catalogconvert "o.o/backend/com/main/catalog/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	servicelocation "o.o/backend/com/main/location"
	ordermodel "o.o/backend/com/main/ordering/model"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodely "o.o/backend/com/main/shipping/modely"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etc/typeutil"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

var locationBus = servicelocation.QueryMessageBus(servicelocation.New(nil))

func PbOrdersWithFulfillments(items []ordermodelx.OrderWithFulfillments, accType int, shops []*identitymodel.Shop) []*types.Order {
	res := make([]*types.Order, len(items))
	shopsMap := make(map[dot.ID]*identitymodel.Shop)
	for _, shop := range shops {
		shopsMap[shop.ID] = shop
	}
	for i, item := range items {
		res[i] = XPbOrder(item.Order, item.Fulfillments, accType)
		res[i].ShopName = shopsMap[item.ShopID].GetShopName()
	}
	return res
}

func PbOrders(items []*ordermodel.Order, accType int) []*types.Order {
	res := make([]*types.Order, len(items))
	for i, order := range items {
		res[i] = PbOrder(order, nil, accType)
	}
	return res
}

var exportedOrder = cm.SortStrings([]string{
	"id", "shop_id", "code", "external_id", "external_code", "external_url", "self_url",
	"customer_address", "shipping_address",
	"created_at", "processed_at", "updated_at", "closed_at", "confirmed_at", "cancelled_at",
	"confirm_status", "status", "fulfillment_shipping_status", "etop_payment_status",
	"lines", "total_items", "fee_lines", "order_note", "shipping", "cancel_reason",
	"basket_value", "order_discount", "total_discount", "total_fee",
})

func PbOrder(m *ordermodel.Order, fulfillments []*shipmodel.Fulfillment, accType int) *types.Order {
	ffms := make([]*ordermodelx.Fulfillment, len(fulfillments))
	for i, ffm := range fulfillments {
		ffms[i] = &ordermodelx.Fulfillment{
			Shipment: ffm,
		}
	}
	order := &types.Order{
		ExportedFields:            exportedOrder,
		Id:                        m.ID,
		ShopId:                    m.ShopID,
		ShopName:                  "",
		Code:                      m.Code,
		EdCode:                    m.EdCode,
		ExternalCode:              m.EdCode,
		Source:                    m.OrderSourceType,
		PartnerId:                 m.PartnerID,
		ExternalId:                m.ExternalOrderID,
		ExternalUrl:               m.ExternalURL,
		SelfUrl:                   m.SelfURL(cm.MainSiteBaseURL(), accType),
		PaymentMethod:             m.PaymentMethod,
		Customer:                  PbOrderCustomer(m.Customer),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		BillingAddress:            PbOrderAddress(m.BillingAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 cmapi.PbTime(m.CreatedAt),
		CreatedBy:                 m.CreatedBy,
		ProcessedAt:               cmapi.PbTime(m.ProcessedAt),
		UpdatedAt:                 cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                  cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:               cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:               cmapi.PbTime(m.CancelledAt),
		CancelReason:              m.CancelReason,
		ShConfirm:                 m.ShopConfirm,
		Confirm:                   m.ConfirmStatus,
		ConfirmStatus:             m.ConfirmStatus,
		Status:                    m.Status,
		FulfillmentStatus:         m.FulfillmentShippingStatus,
		FulfillmentShippingStatus: m.FulfillmentShippingStatus,
		EtopPaymentStatus:         m.EtopPaymentStatus,
		PaymentStatus:             m.PaymentStatus,
		Lines:                     PbOrderLines(m.Lines),
		Discounts:                 PbDiscounts(m.Discounts),
		TotalItems:                m.TotalItems,
		BasketValue:               m.BasketValue,
		TotalWeight:               m.TotalWeight,
		OrderDiscount:             m.OrderDiscount,
		TotalDiscount:             m.TotalDiscount,
		TotalAmount:               m.TotalAmount,
		OrderNote:                 m.OrderNote,
		ShippingFee:               m.ShopShippingFee,
		TotalFee:                  m.GetTotalFee(),
		FeeLines:                  PbOrderFeeLines(m.FeeLines),
		ShopShippingFee:           m.ShopShippingFee,
		ShippingNote:              m.ShippingNote,
		ShopCod:                   m.ShopCOD,
		ReferenceUrl:              m.ReferenceURL,
		Fulfillments:              XPbFulfillments(ffms, accType),
		ShopShipping:              nil,
		Shipping:                  nil,
		GhnNoteCode:               m.GhnNoteCode,
		FulfillmentType:           m.FulfillmentType.String(),
		FulfillmentIds:            m.FulfillmentIDs,
		CustomerId:                m.CustomerID,
		PreOrder:                  m.PreOrder,
	}
	shipping := PbOrderShipping(m)
	order.ShopShipping = shipping
	order.Shipping = shipping
	return order
}

func XPbOrder(m *ordermodel.Order, fulfillments []*ordermodelx.Fulfillment, accType int) *types.Order {
	order := &types.Order{
		ExportedFields: exportedOrder,

		Id:                        m.ID,
		ShopId:                    m.ShopID,
		ShopName:                  "",
		Code:                      m.Code,
		EdCode:                    m.EdCode,
		ExternalCode:              m.EdCode,
		Source:                    m.OrderSourceType,
		PartnerId:                 m.PartnerID,
		ExternalId:                m.ExternalOrderID,
		ExternalUrl:               m.ExternalURL,
		SelfUrl:                   m.SelfURL(cm.MainSiteBaseURL(), accType),
		PaymentMethod:             m.PaymentMethod,
		Customer:                  PbOrderCustomer(m.Customer),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		BillingAddress:            PbOrderAddress(m.BillingAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 cmapi.PbTime(m.CreatedAt),
		ProcessedAt:               cmapi.PbTime(m.ProcessedAt),
		UpdatedAt:                 cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                  cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:               cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:               cmapi.PbTime(m.CancelledAt),
		CancelReason:              m.CancelReason,
		ShConfirm:                 m.ShopConfirm,
		Confirm:                   m.ConfirmStatus,
		ConfirmStatus:             m.ConfirmStatus,
		Status:                    m.Status,
		FulfillmentStatus:         m.FulfillmentShippingStatus,
		FulfillmentShippingStatus: m.FulfillmentShippingStatus,
		EtopPaymentStatus:         m.EtopPaymentStatus,
		PaymentStatus:             m.PaymentStatus,
		Lines:                     PbOrderLines(m.Lines),
		Discounts:                 PbDiscounts(m.Discounts),
		TotalItems:                m.TotalItems,
		BasketValue:               m.BasketValue,
		TotalWeight:               m.TotalWeight,
		OrderDiscount:             m.OrderDiscount,
		TotalDiscount:             m.TotalDiscount,
		TotalAmount:               m.TotalAmount,
		OrderNote:                 m.OrderNote,
		ShippingFee:               m.ShopShippingFee,
		TotalFee:                  m.GetTotalFee(),
		FeeLines:                  PbOrderFeeLines(m.FeeLines),
		ShopShippingFee:           m.ShopShippingFee,
		ShippingNote:              m.ShippingNote,
		ShopCod:                   m.ShopCOD,
		ReferenceUrl:              m.ReferenceURL,
		Fulfillments:              XPbFulfillments(fulfillments, accType),
		ShopShipping:              nil,
		Shipping:                  nil,
		GhnNoteCode:               m.GhnNoteCode,
		FulfillmentType:           m.FulfillmentType.String(),
		FulfillmentIds:            m.FulfillmentIDs,
		CustomerId:                m.CustomerID,
		CreatedBy:                 m.CreatedBy,
		PreOrder:                  m.PreOrder,
	}
	shipping := PbOrderShipping(m)
	order.ShopShipping = shipping
	order.Shipping = shipping
	return order
}

var exportedOrderShipping = cm.SortStrings([]string{
	"pickup_address",
	"shipping_service_name", "shipping_service_code", "shipping_service_fee",
	"carrier", "try_on", "include_insurance", "shipping_node", "cod_amount",
	"weight", "gross_weight", "length", "width", "height", "chargeable_weight",
})

func PbOrderShipping(m *ordermodel.Order) *types.OrderShipping {
	if m == nil {
		return nil
	}
	item := m.ShopShipping
	if item == nil {
		item = &ordermodel.OrderShipping{}
	}
	return &types.OrderShipping{
		ExportedFields: exportedOrderShipping,
		// @deprecated fields
		ShAddress:    PbOrderAddress(item.ShopAddress),
		XServiceId:   cm.Coalesce(item.ProviderServiceID, item.ExternalServiceID),
		XShippingFee: item.ExternalShippingFee,
		XServiceName: item.ExternalServiceName,

		PickupAddress:       PbOrderAddress(item.GetPickupAddress()),
		ShippingServiceName: item.ExternalServiceName,
		ShippingServiceCode: item.GetShippingServiceCode(),
		ShippingServiceFee:  item.ExternalShippingFee,
		ShippingProvider:    PbShippingProviderType(item.ShippingProvider),
		Carrier:             PbShippingProviderType(item.ShippingProvider),
		IncludeInsurance:    item.IncludeInsurance,
		TryOn:               m.GetTryOn(),
		ShippingNote:        m.ShippingNote,
		CodAmount:           dot.Int(m.ShopCOD),
		Weight:              dot.Int(m.TotalWeight),
		GrossWeight:         dot.Int(m.TotalWeight),
		Length:              dot.Int(item.Length),
		Width:               dot.Int(item.Width),
		Height:              dot.Int(item.Height),
		ChargeableWeight:    dot.Int(m.TotalWeight),
	}
}

var exportedOrderCustomer = cm.SortStrings([]string{
	"full_name", "email", "phone", "gender",
})

func PbOrderCustomer(m *ordermodel.OrderCustomer) *types.OrderCustomer {
	if m == nil {
		return nil
	}
	_gender, _ := gender.ParseGender(m.Gender)
	return &types.OrderCustomer{
		ExportedFields: exportedOrderCustomer,

		FirstName: m.FirstName,
		LastName:  m.LastName,
		FullName:  m.GetFullName(),
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    _gender,
	}
}

func OrderCustomerToModel(m *types.OrderCustomer) *ordermodel.OrderCustomer {
	if m == nil {
		return nil
	}
	return &ordermodel.OrderCustomer{
		FirstName: m.FirstName,
		LastName:  m.LastName,
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    m.Gender.String(),
	}
}

var exportedOrderAddress = cm.SortStrings([]string{
	"full_name", "phone", "province", "district",
	"ward", "company", "address1", "address2",
})

func PbOrderAddress(m *ordermodel.OrderAddress) *types.OrderAddress {
	if m == nil {
		return nil
	}
	res := &types.OrderAddress{
		ExportedFields: exportedOrderAddress,
		FullName:       m.GetFullName(),
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		Phone:          m.Phone,
		Email:          m.Email,
		Country:        m.Country,
		City:           m.City,
		Province:       m.Province,
		District:       m.District,
		Ward:           m.Ward,
		Zip:            m.Zip,
		Company:        m.Company,
		Address1:       m.Address1,
		Address2:       m.Address2,
		ProvinceCode:   m.ProvinceCode,
		DistrictCode:   m.DistrictCode,
		WardCode:       m.WardCode,
	}

	if m.Coordinates != nil {
		res.Coordinates = &etop.Coordinates{
			Latitude:  m.Coordinates.Latitude,
			Longitude: m.Coordinates.Longitude,
		}
	}
	return res
}

func OrderAddressToModel(m *types.OrderAddress) (*ordermodel.OrderAddress, error) {
	if m == nil {
		return nil, nil
	}
	res := &ordermodel.OrderAddress{
		FullName:     m.FullName,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		Phone:        m.Phone,
		Email:        m.Email,
		Country:      m.Country,
		City:         m.City,
		Province:     m.Province,
		District:     m.District,
		Ward:         m.Ward,
		Zip:          m.Zip,
		DistrictCode: m.DistrictCode,
		ProvinceCode: m.ProvinceCode,
		WardCode:     m.WardCode,
		Company:      m.Company,
		Address1:     m.Address1,
		Address2:     m.Address2,
	}
	locationQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: m.ProvinceCode,
		DistrictCode: m.DistrictCode,
		WardCode:     m.WardCode,
		Province:     m.Province,
		District:     m.District,
		Ward:         m.Ward,
	}
	if err := locationBus.Dispatch(context.TODO(), locationQuery); err != nil {
		return nil, err
	}
	loc := locationQuery.Result
	if loc.Province == nil || loc.District == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cần cung cấp thông tin tỉnh/thành phố và quận/huyện")
	}

	res.Province = loc.Province.Name
	res.ProvinceCode = loc.Province.Code
	res.District = loc.District.Name
	res.DistrictCode = loc.District.Code
	if loc.Ward != nil {
		res.Ward = loc.Ward.Name
		res.WardCode = loc.Ward.Code
	}
	if m.Coordinates != nil {
		res.Coordinates = &addressmodel.Coordinates{
			Latitude:  m.Coordinates.Latitude,
			Longitude: m.Coordinates.Longitude,
		}
	}
	return res, nil
}

func OrderAddressFulfilled(m *types.OrderAddress) (*types.OrderAddress, error) {
	if m == nil {
		return nil, nil
	}
	locationQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: m.ProvinceCode,
		DistrictCode: m.DistrictCode,
		WardCode:     m.WardCode,
		Province:     m.Province,
		District:     m.District,
		Ward:         m.Ward,
	}
	if err := locationBus.Dispatch(context.TODO(), locationQuery); err != nil {
		return nil, err
	}
	loc := locationQuery.Result
	if loc.Province == nil || loc.District == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cần cung cấp thông tin tỉnh/thành phố và quận/huyện")
	}

	m.Province = loc.Province.Name
	m.ProvinceCode = loc.Province.Code
	m.District = loc.District.Name
	m.DistrictCode = loc.District.Code
	if loc.Ward != nil {
		m.Ward = loc.Ward.Name
		m.WardCode = loc.Ward.Code
	}
	return m, nil
}

func PbOrderAddressFromAddress(m *addressmodel.Address) *types.OrderAddress {
	if m == nil {
		return nil
	}
	return &types.OrderAddress{
		ExportedFields: exportedOrderAddress,

		FullName:     m.FullName,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		Phone:        m.Phone,
		Country:      m.Country,
		City:         m.City,
		Province:     m.Province,
		District:     m.District,
		Ward:         m.Ward,
		Zip:          m.Zip,
		Company:      m.Company,
		Address1:     m.Address1,
		Address2:     m.Address2,
		ProvinceCode: m.ProvinceCode,
		DistrictCode: m.DistrictCode,
		WardCode:     m.WardCode,
	}
}

func OrderShippingToModel(ctx context.Context, m *types.OrderShipping, mo *ordermodel.Order) error {
	if m == nil {
		return nil
	}

	var pickupAddress *types.OrderAddress
	if m.ShAddress != nil {
		pickupAddress = m.ShAddress
	} else {
		pickupAddress = m.PickupAddress
	}
	modelPickupAddress, err := OrderAddressToModel(pickupAddress)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}
	modelReturnAddress, err := OrderAddressToModel(m.ReturnAddress)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ trả hàng không hợp lệ: %v", err)
	}

	carrier := m.ShippingProvider
	if m.Carrier != 0 {
		carrier = m.Carrier
	}

	grossWeight := m.Weight.Apply(0)
	grossWeight = m.GrossWeight.Apply(grossWeight)

	chargeableWeight := m.ChargeableWeight.Apply(0)
	if chargeableWeight == 0 {
		switch {
		case !m.Length.Valid && !m.Width.Valid && !m.Height.Valid:
			// continue
		case m.Length.Valid && m.Width.Valid && m.Height.Valid:
			chargeableWeight = model.CalcChargeableWeight(grossWeight, m.Length.Int, m.Width.Int, m.Height.Int)
		default:
			return cm.Errorf(cm.InvalidArgument, err, "Cần cung cấp đủ các giá trị length, width, height (hoặc để trống cả 3)", err)
		}
	}

	shippingServiceCode := cm.Coalesce(m.ShippingServiceCode, m.XServiceId)

	orderShipping := &ordermodel.OrderShipping{
		ShopAddress:         modelPickupAddress,
		ReturnAddress:       modelReturnAddress,
		ExternalServiceID:   shippingServiceCode,
		ExternalShippingFee: cm.CoalesceInt(m.ShippingServiceFee, m.XShippingFee),
		ExternalServiceName: m.ShippingServiceName,
		ShippingProvider:    carrier,
		ProviderServiceID:   cm.Coalesce(shippingServiceCode, m.XServiceId),
		IncludeInsurance:    m.IncludeInsurance,
		Length:              m.Length.Apply(0),
		Width:               m.Width.Apply(0),
		Height:              m.Height.Apply(0),
		GrossWeight:         grossWeight,
		ChargeableWeight:    chargeableWeight,
	}

	// when adding new fields here, remember to also change UpdateOrderCommand
	mo.ShopShipping = orderShipping
	mo.ShopCOD = m.CodAmount.Apply(mo.ShopCOD)
	mo.TotalWeight = chargeableWeight

	if m.TryOn != 0 {
		mo.TryOn = m.TryOn
		mo.GhnNoteCode = typeutil.GHNNoteCodeFromTryOn(m.TryOn)
	} else if mo.GhnNoteCode != 0 {
		mo.TryOn = typeutil.TryOnFromGHNNoteCode(mo.GhnNoteCode)
	}

	// Coalesce takes from left to right while PatchInt takes from right
	mo.ShippingNote = cm.Coalesce(m.ShippingNote, mo.ShippingNote)
	return nil
}

func PbDiscounts(items []*ordermodel.OrderDiscount) []*types.OrderDiscount {
	res := make([]*types.OrderDiscount, len(items))
	for i, item := range items {
		res[i] = PbDiscount(item)
	}
	return res
}

func PbDiscount(m *ordermodel.OrderDiscount) *types.OrderDiscount {
	return &types.OrderDiscount{
		Code:   m.Code,
		Type:   m.Type,
		Amount: m.Amount,
	}
}

func OrderDiscountToModel(m *types.OrderDiscount) *ordermodel.OrderDiscount {
	return &ordermodel.OrderDiscount{
		Code:   m.Code,
		Type:   m.Type,
		Amount: m.Amount,
	}
}

func PbOrderDiscountsToModel(discounts []*types.OrderDiscount) []*ordermodel.OrderDiscount {
	res := make([]*ordermodel.OrderDiscount, len(discounts))
	for i, d := range discounts {
		res[i] = OrderDiscountToModel(d)
	}
	return res
}

func PbOrderFeeLinesToModel(items []*types.OrderFeeLine) []ordermodel.OrderFeeLine {
	res := make([]ordermodel.OrderFeeLine, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		res = append(res, ordermodel.OrderFeeLine{
			Amount: item.Amount.Int(),
			Desc:   item.Desc,
			Code:   item.Code,
			Name:   item.Name,
			Type:   item.Type,
		})
	}
	return res
}

func PbOrderFeeLines(items []ordermodel.OrderFeeLine) []*types.OrderFeeLine {
	res := make([]*types.OrderFeeLine, len(items))
	for i, item := range items {
		res[i] = &types.OrderFeeLine{
			Type:   item.Type,
			Name:   item.Name,
			Code:   item.Code,
			Desc:   item.Desc,
			Amount: types.Int(item.Amount),
		}
	}
	return res
}

func PbOrderLines(items []*ordermodel.OrderLine) []*types.OrderLine {
	res := make([]*types.OrderLine, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		res = append(res, PbOrderLine(item))
	}
	return res
}

var exportedOrderLine = cm.SortStrings([]string{
	"order_id", "product_id", "variant_id",
	"product_name", "image_url", "attributes",
	"quantity", "list_price", "retail_price", "payment_price",
})

func PbOrderLine(m *ordermodel.OrderLine) *types.OrderLine {
	metaFields := []*types.OrderLineMetaField{}
	for _, metaField := range m.MetaFields {
		metaFields = append(metaFields, &types.OrderLineMetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
			Name:  metaField.Name,
		})
	}
	if m == nil {
		return nil
	}
	return &types.OrderLine{
		ExportedFields: exportedOrderLine,
		Code:           m.Code,
		OrderId:        m.OrderID,
		VariantId:      m.VariantID,
		ProductName:    m.ProductName,
		IsOutsideEtop:  m.IsOutsideEtop,
		Quantity:       m.Quantity,
		ListPrice:      m.ListPrice,
		RetailPrice:    m.RetailPrice,
		PaymentPrice:   m.PaymentPrice,
		ImageUrl:       m.ImageURL,
		Attributes:     catalogconvert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(m.Attributes),
		ProductId:      m.ProductID,
		TotalDiscount:  m.TotalDiscount,
		MetaFields:     metaFields,
	}
}

func PbFulfillments(items []*shipmodel.Fulfillment, accType int) []*types.Fulfillment {
	if items == nil {
		return nil
	}
	res := make([]*types.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item, accType, nil, nil)
	}
	return res
}

func PbFulfillmentExtendeds(items []*shipmodely.FulfillmentExtended, accType int) []*types.Fulfillment {
	if items == nil {
		return nil
	}
	res := make([]*types.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item.Fulfillment, accType, item.Shop, item.Order)
	}
	return res
}

var exportedFulfillment = cm.SortStrings([]string{
	"id", "order_id", "shop_id", "self_url",
	"created_at", "updated_at", "closed_at", "shipping_cancelled_at",
	"lines", "total_items", "basket_value", "cod_amount", "chargeable_weight",
	"include_insurance", "try_on", "carrier",
	"shipping_service_name", "shipping_service_fee", "shipping_service_code",
	"cancel_reason", "pickup_address", "return_address", "shipping_address",
	"status", "confirm_status", "shipping_state", "shipping_status", "etop_payment_status",
	"estimated_delivery_at", "estimated_pickup_at",
})

func PbFulfillment(m *shipmodel.Fulfillment, accType int, shop *identitymodel.Shop, mo *ordermodel.Order) *types.Fulfillment {
	if m == nil {
		return nil
	}
	ff := &types.Fulfillment{
		ExportedFields: exportedFulfillment,

		Id:                                 m.ID,
		OrderId:                            m.OrderID,
		ShopId:                             m.ShopID,
		PartnerId:                          m.PartnerID,
		SelfUrl:                            m.SelfURL(cm.MainSiteBaseURL(), accType),
		Lines:                              PbOrderLines(m.Lines),
		TotalItems:                         m.TotalItems,
		TotalWeight:                        m.TotalWeight,
		BasketValue:                        m.BasketValue,
		TotalCodAmount:                     m.TotalCODAmount,
		CodAmount:                          m.TotalCODAmount,
		TotalAmount:                        m.BasketValue, // deprecated
		ChargeableWeight:                   m.ChargeableWeight,
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		CancelledAt:                        cmapi.PbTime(m.ShippingCancelledAt),
		CancelReason:                       m.CancelReason,
		ShippingProvider:                   m.ShippingProvider,
		Carrier:                            PbShippingProviderType(m.ShippingProvider),
		ShippingServiceName:                m.ExternalShippingName,
		ShippingServiceFee:                 m.ExternalShippingFee,
		ShippingServiceCode:                m.ProviderServiceID,
		ShippingCode:                       m.ShippingCode,
		ShippingNote:                       m.ShippingNote,
		TryOn:                              m.TryOn,
		ShippingPaymentType:                m.ShippingPaymentType,
		IncludeInsurance:                   m.IncludeInsurance.Apply(false),
		ShConfirm:                          m.ShopConfirm,
		ShippingState:                      m.ShippingState,
		Status:                             m.Status,
		ShippingStatus:                     m.ShippingStatus,
		EtopPaymentStatus:                  m.EtopPaymentStatus,
		ShippingFeeCustomer:                m.ShippingFeeCustomer,
		ShippingFeeShop:                    m.ShippingFeeShop,
		XShippingFee:                       m.ExternalShippingFee,
		XShippingId:                        m.ExternalShippingID,
		XShippingCode:                      m.ExternalShippingCode,
		XShippingCreatedAt:                 cmapi.PbTime(m.ExternalShippingCreatedAt),
		XShippingUpdatedAt:                 cmapi.PbTime(m.ExternalShippingUpdatedAt),
		XShippingCancelledAt:               cmapi.PbTime(m.ExternalShippingCancelledAt),
		XShippingDeliveredAt:               cmapi.PbTime(m.ExternalShippingDeliveredAt),
		XShippingReturnedAt:                cmapi.PbTime(m.ExternalShippingReturnedAt),
		ExpectedDeliveryAt:                 cmapi.PbTime(m.ExpectedDeliveryAt),
		ExpectedPickAt:                     cmapi.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt:                cmapi.PbTime(m.ExpectedDeliveryAt),
		EstimatedPickupAt:                  cmapi.PbTime(m.ExpectedPickAt),
		CodEtopTransferedAt:                cmapi.PbTime(m.CODEtopTransferedAt),
		ShippingFeeShopTransferedAt:        cmapi.PbTime(m.ShippingFeeShopTransferedAt),
		XShippingState:                     m.ExternalShippingState,
		XShippingStatus:                    m.ExternalShippingStatus,
		XSyncStatus:                        m.SyncStatus,
		XSyncStates:                        PbFulfillmentSyncStates(m.SyncStates),
		AddressTo:                          PbAddress(m.AddressTo),
		AddressFrom:                        PbAddress(m.AddressFrom),
		PickupAddress:                      PbOrderAddressFromAddress(m.AddressFrom),
		ReturnAddress:                      PbOrderAddressFromAddress(m.AddressFrom),
		ShippingAddress:                    PbOrderAddressFromAddress(m.AddressTo),
		Shop:                               nil,
		Order:                              nil,
		ProviderShippingFeeLines:           PbShippingFeeLines(m.ProviderShippingFeeLines),
		ShippingFeeShopLines:               PbShippingFeeLines(m.ShippingFeeShopLines),
		EtopDiscount:                       m.EtopDiscount,
		MoneyTransactionShippingId:         m.MoneyTransactionID,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		XShippingLogs:                      PbExternalShippingLogs(m.ExternalShippingLogs),
		XShippingNote:                      m.ExternalShippingNote.String,
		XShippingSubState:                  m.ExternalShippingSubState.String,
		ActualCompensationAmount:           m.ActualCompensationAmount,
		ConnectionID:                       m.ConnectionID,
		InsuranceValue:                     m.InsuranceValue.Apply(0),
		UpdatedBy:                          m.UpdatedBy,
		ShopCarrierID:                      m.ShopCarrierID,
		EdCode:                             m.EdCode,
		LinesContent:                       m.LinesContent,
		ShippingSubstate:                   m.ShippingSubstate.Enum,
	}
	if shop != nil {
		ff.Shop = PbShop(shop)
	}
	if mo != nil {
		ff.Order = PbOrder(mo, nil, accType)
	}
	if accType == account_tag.TagEtop {
		ff.AdminNote = m.AdminNote
	}
	return ff
}

func XPbFulfillments(items []*ordermodelx.Fulfillment, accType int) []*types.XFulfillment {
	if items == nil {
		return nil
	}
	res := make([]*types.XFulfillment, len(items))
	for i, item := range items {
		res[i] = XPbFulfillment(item, accType, nil, nil)
	}
	return res
}

func XPbFulfillment(m *ordermodelx.Fulfillment, accType int, shop *identitymodel.Shop, mo *ordermodel.Order) *types.XFulfillment {
	res := &types.XFulfillment{}
	shipment := PbFulfillment(m.Shipment, accType, shop, mo)
	if shipment != nil {
		res = &types.XFulfillment{
			Shipment:                           shipment,
			Id:                                 shipment.Id,
			OrderId:                            shipment.OrderId,
			ShopId:                             shipment.ShopId,
			PartnerId:                          shipment.PartnerId,
			SelfUrl:                            shipment.SelfUrl,
			Lines:                              shipment.Lines,
			TotalItems:                         shipment.TotalItems,
			TotalWeight:                        shipment.ChargeableWeight,
			BasketValue:                        shipment.BasketValue,
			TotalCodAmount:                     shipment.TotalCodAmount,
			CodAmount:                          shipment.CodAmount,
			TotalAmount:                        shipment.TotalAmount,
			ChargeableWeight:                   shipment.ChargeableWeight,
			CreatedAt:                          shipment.CreatedAt,
			UpdatedAt:                          shipment.UpdatedAt,
			ClosedAt:                           shipment.ClosedAt,
			CancelledAt:                        shipment.CancelledAt,
			CancelReason:                       shipment.CancelReason,
			ShippingProvider:                   shipment.ShippingProvider,
			Carrier:                            shipment.Carrier,
			ShippingServiceName:                shipment.ShippingServiceName,
			ShippingServiceFee:                 shipment.ShippingServiceFee,
			ShippingServiceCode:                shipment.ShippingServiceCode,
			ShippingCode:                       shipment.ShippingCode,
			ShippingNote:                       shipment.ShippingNote,
			TryOn:                              shipment.TryOn,
			IncludeInsurance:                   shipment.IncludeInsurance,
			ShConfirm:                          shipment.ShConfirm,
			ShippingState:                      shipment.ShippingState,
			Status:                             shipment.Status,
			ShippingStatus:                     shipment.ShippingStatus,
			EtopPaymentStatus:                  shipment.EtopPaymentStatus,
			ShippingFeeCustomer:                shipment.ShippingFeeCustomer,
			ShippingFeeShop:                    shipment.ShippingFeeShop,
			XShippingFee:                       shipment.XShippingFee,
			XShippingId:                        shipment.XShippingId,
			XShippingCode:                      shipment.XShippingCode,
			XShippingCreatedAt:                 shipment.XShippingCreatedAt,
			XShippingUpdatedAt:                 shipment.XShippingUpdatedAt,
			XShippingCancelledAt:               shipment.XShippingCancelledAt,
			XShippingDeliveredAt:               shipment.XShippingDeliveredAt,
			XShippingReturnedAt:                shipment.XShippingReturnedAt,
			ExpectedDeliveryAt:                 shipment.ExpectedDeliveryAt,
			ExpectedPickAt:                     shipment.ExpectedPickAt,
			EstimatedDeliveryAt:                shipment.EstimatedDeliveryAt,
			EstimatedPickupAt:                  shipment.EstimatedPickupAt,
			CodEtopTransferedAt:                shipment.CodEtopTransferedAt,
			ShippingFeeShopTransferedAt:        shipment.ShippingFeeShopTransferedAt,
			XShippingState:                     shipment.XShippingState,
			XShippingStatus:                    shipment.XShippingStatus,
			XSyncStatus:                        shipment.XSyncStatus,
			XSyncStates:                        shipment.XSyncStates,
			AddressTo:                          shipment.AddressTo,
			AddressFrom:                        shipment.AddressFrom,
			PickupAddress:                      shipment.PickupAddress,
			ReturnAddress:                      shipment.ReturnAddress,
			ShippingAddress:                    shipment.ShippingAddress,
			Shop:                               shipment.Shop,
			Order:                              shipment.Order,
			ProviderShippingFeeLines:           shipment.ProviderShippingFeeLines,
			ShippingFeeShopLines:               shipment.ShippingFeeShopLines,
			EtopDiscount:                       shipment.EtopDiscount,
			MoneyTransactionShippingId:         shipment.MoneyTransactionShippingId,
			MoneyTransactionShippingExternalId: shipment.MoneyTransactionShippingExternalId,
			XShippingLogs:                      shipment.XShippingLogs,
			XShippingNote:                      shipment.XShippingNote,
			XShippingSubState:                  shipment.XShippingSubState,
			Code:                               shipment.Code,
			ActualCompensationAmount:           shipment.ActualCompensationAmount,
		}
	}
	if m.Shipnow != nil {
		res.Shipnow = Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(m.Shipnow)
	}

	return res
}

func PbAvailableShippingServices(items []*shippingsharemodel.AvailableShippingService) []*types.ExternalShippingService {
	res := make([]*types.ExternalShippingService, len(items))
	for i, item := range items {
		res[i] = PbAvailableShippingService(item)
	}
	return res
}

var exportedShippingService = cm.SortStrings([]string{
	"name", "code", "fee", "carrier", "estimated_pickup_at", "estimated_delivery_at",
})

func PbAvailableShippingService(s *shippingsharemodel.AvailableShippingService) *types.ExternalShippingService {
	return &types.ExternalShippingService{
		ExportedFields: exportedShippingService,

		ExternalId:          s.ProviderServiceID,
		ServiceFee:          s.ServiceFee,
		Provider:            PbShippingProviderType(s.Provider),
		ExpectedPickAt:      cmapi.PbTime(s.ExpectedPickAt),
		ExpectedDeliveryAt:  cmapi.PbTime(s.ExpectedDeliveryAt),
		Name:                s.Name,
		Code:                s.ProviderServiceID,
		Fee:                 s.ServiceFee,
		Carrier:             PbShippingProviderType(s.Provider),
		EstimatedPickupAt:   cmapi.PbTime(s.ExpectedPickAt),
		EstimatedDeliveryAt: cmapi.PbTime(s.ExpectedDeliveryAt),
		ConnectionInfo:      PbConnectionInfo(s.ConnectionInfo),
		ShipmentServiceInfo: PbShipmentServiceInfo(s.ShipmentServiceInfo),
		ShipmentPriceInfo:   PbShipmentPriceInfo(s.ShipmentPriceInfo),
	}
}

func PbConnectionInfo(item *shippingsharemodel.ConnectionInfo) *types.ConnectionInfo {
	if item == nil {
		return nil
	}
	return &types.ConnectionInfo{
		ID:       item.ID,
		Name:     item.Name,
		ImageURL: item.ImageURL,
	}
}

func PbShipmentServiceInfo(item *shippingsharemodel.ShipmentServiceInfo) *types.ShipmentServiceInfo {
	if item == nil {
		return nil
	}
	return &types.ShipmentServiceInfo{
		ID:           item.ID,
		Code:         item.Code,
		Name:         item.Name,
		IsAvailable:  item.IsAvailable,
		ErrorMessage: item.ErrorMessage,
	}
}

func PbShipmentPriceInfo(item *shippingsharemodel.ShipmentPriceInfo) *types.ShipmentPriceInfo {
	if item == nil {
		return nil
	}
	return &types.ShipmentPriceInfo{
		ShipmentPriceID: item.ShipmentPriceID,
		OriginFee:       item.OriginFee,
		MakeupFee:       item.MakeupFee,
	}
}

func PbShippingFeeLines(items []*shippingsharemodel.ShippingFeeLine) []*types.ShippingFeeLine {
	result := make([]*types.ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = PbShippingFeeLine(item)
	}
	return result
}

func PbShippingFeeLine(line *shippingsharemodel.ShippingFeeLine) *types.ShippingFeeLine {
	return &types.ShippingFeeLine{
		ShippingFeeType: line.ShippingFeeType,
		Cost:            line.Cost,
	}
}

func PbFulfillmentSyncStates(m *shippingsharemodel.FulfillmentSyncStates) *types.FulfillmentSyncStates {
	if m == nil {
		return nil
	}
	return &types.FulfillmentSyncStates{
		SyncAt:            cmapi.PbTime(m.SyncAt),
		NextShippingState: m.NextShippingState,
		Error:             cmapi.PbError(m.Error),
	}
}

func PbPublicFulfillment(item *shipmodel.Fulfillment) *types.PublicFulfillment {
	timeLayout := "15:04 02/01/2006"

	// use for manychat
	var deliveredAtText string
	if !item.ShippingDeliveredAt.IsZero() {
		deliveredAtText = item.ShippingDeliveredAt.In(time.Local).Format(timeLayout)
	} else {
		deliveredAtText = item.ExpectedDeliveryAt.In(time.Local).Format(timeLayout)
		deliveredAtText += " (dự kiến)"
	}

	return &types.PublicFulfillment{
		Id:                 item.ID,
		ShippingState:      item.ShippingState,
		Status:             item.Status,
		ExpectedDeliveryAt: cmapi.PbTime(item.ExpectedDeliveryAt),
		DeliveredAt:        cmapi.PbTime(item.ShippingDeliveredAt),
		ShippingCode:       item.ShippingCode,
		OrderId:            item.OrderID,
		DeliveredAtText:    deliveredAtText,
		ShippingStateText:  item.ShippingState.Text(),
	}
}

func PbExternalShippingLogs(items []*shipmodel.ExternalShippingLog) []*types.ExternalShippingLog {
	result := make([]*types.ExternalShippingLog, len(items))
	for i, item := range items {
		result[i] = PbExternalShippingLog(item)
	}
	return result
}

func PbExternalShippingLog(l *shipmodel.ExternalShippingLog) *types.ExternalShippingLog {
	return &types.ExternalShippingLog{
		StateText: l.StateText,
		Time:      l.Time,
		Message:   l.Message,
	}
}

func Convert_core_Fulfillment_To_api_Fulfillment(m *shipping.Fulfillment, accType int, shop *identity.Shop, mo *ordering.Order) *types.Fulfillment {
	if m == nil {
		return nil
	}
	ff := &types.Fulfillment{
		ExportedFields: exportedFulfillment,

		Id:                                 m.ID,
		OrderId:                            m.OrderID,
		ShopId:                             m.ShopID,
		PartnerId:                          m.PartnerID,
		SelfUrl:                            FulfillmentSelfURL(m, accType),
		Lines:                              Convert_core_OrderLines_To_api_OrderLines(m.Lines),
		TotalItems:                         m.TotalItems,
		TotalWeight:                        m.ChargeableWeight,
		BasketValue:                        m.BasketValue,
		TotalCodAmount:                     m.TotalCODAmount,
		CodAmount:                          m.TotalCODAmount,
		TotalAmount:                        m.BasketValue, // deprecated
		ChargeableWeight:                   m.ChargeableWeight,
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		CancelledAt:                        cmapi.PbTime(m.ShippingCancelledAt),
		CancelReason:                       m.CancelReason,
		ShippingProvider:                   m.ShippingProvider,
		Carrier:                            m.ShippingProvider,
		ShippingServiceName:                m.ExternalShippingName,
		ShippingServiceFee:                 m.ExternalShippingFee,
		ShippingServiceCode:                m.ProviderServiceID,
		ShippingCode:                       m.ShippingCode,
		ShippingNote:                       m.ShippingNote,
		TryOn:                              m.TryOn,
		IncludeInsurance:                   m.IncludeInsurance.Apply(false),
		ShippingPaymentType:                m.ShippingPaymentType,
		ShConfirm:                          m.ShopConfirm,
		ShippingState:                      m.ShippingState,
		Status:                             m.Status,
		ShippingStatus:                     m.ShippingStatus,
		EtopPaymentStatus:                  m.EtopPaymentStatus,
		ShippingFeeCustomer:                m.ShippingFeeCustomer,
		ShippingFeeShop:                    m.ShippingFeeShop,
		XShippingFee:                       m.ExternalShippingFee,
		XShippingId:                        m.ExternalShippingID,
		XShippingCode:                      m.ExternalShippingCode,
		XShippingCreatedAt:                 cmapi.PbTime(m.ExternalShippingCreatedAt),
		XShippingUpdatedAt:                 cmapi.PbTime(m.ExternalShippingUpdatedAt),
		XShippingCancelledAt:               cmapi.PbTime(m.ExternalShippingCancelledAt),
		XShippingDeliveredAt:               cmapi.PbTime(m.ExternalShippingDeliveredAt),
		XShippingReturnedAt:                cmapi.PbTime(m.ExternalShippingReturnedAt),
		ExpectedDeliveryAt:                 cmapi.PbTime(m.ExpectedDeliveryAt),
		ExpectedPickAt:                     cmapi.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt:                cmapi.PbTime(m.ExpectedDeliveryAt),
		EstimatedPickupAt:                  cmapi.PbTime(m.ExpectedPickAt),
		CodEtopTransferedAt:                cmapi.PbTime(m.CODEtopTransferedAt),
		ShippingFeeShopTransferedAt:        cmapi.PbTime(m.ShippingFeeShopTransferedAt),
		XShippingState:                     m.ExternalShippingState,
		XShippingStatus:                    m.ExternalShippingStatus,
		XSyncStatus:                        m.SyncStatus,
		XSyncStates:                        Convert_core_FulfillmentSyncStates_To_api_FulfillmentSyncStates(m.SyncStates),
		AddressTo:                          Convert_core_OrderAddress_To_api_Address(m.AddressTo),
		AddressFrom:                        Convert_core_OrderAddress_To_api_Address(m.AddressFrom),
		PickupAddress:                      Convert_core_OrderAddress_To_api_OrderAddress(m.AddressFrom),
		ReturnAddress:                      Convert_core_OrderAddress_To_api_OrderAddress(m.AddressFrom),
		ShippingAddress:                    Convert_core_OrderAddress_To_api_OrderAddress(m.AddressTo),
		Shop:                               nil,
		Order:                              nil,
		ProviderShippingFeeLines:           Convert_core_ShippingFeeLines_To_api_ShippingFeeLines(m.ProviderShippingFeeLines),
		ShippingFeeShopLines:               Convert_core_ShippingFeeLines_To_api_ShippingFeeLines(m.ShippingFeeShopLines),
		EtopDiscount:                       m.EtopDiscount,
		MoneyTransactionShippingId:         m.MoneyTransactionID,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		XShippingLogs:                      Convert_core_ExternalShippingLogs_To_api_ExternalShippingLogs(m.ExternalShippingLogs),
		XShippingNote:                      m.ExternalShippingNote.String,
		XShippingSubState:                  m.ExternalShippingSubState.String,
		ActualCompensationAmount:           m.ActualCompensationAmount,
	}
	if shop != nil {
		ff.Shop = Convert_core_Shop_To_api_Shop(shop)
	}
	if mo != nil {
		ff.Order = Convert_core_Order_To_api_Order(mo, nil, accType)
	}
	return ff
}

func Convert_core_Fulfillments_To_api_Fulfillments(items []*shipping.Fulfillment, accType int) []*types.Fulfillment {
	result := make([]*types.Fulfillment, len(items))
	for i, item := range items {
		result[i] = Convert_core_Fulfillment_To_api_Fulfillment(item, accType, nil, nil)
	}
	return result
}

func Convert_core_Fulfillment_To_api_XFulfillment(m *shipping.Fulfillment, accType int, shop *identity.Shop, mo *ordering.Order) *types.XFulfillment {
	ffm := Convert_core_Fulfillment_To_api_Fulfillment(m, accType, shop, mo)
	res := &types.XFulfillment{
		Shipment: ffm,
	}
	return res
}

func Convert_core_Fulfillments_To_api_XFulfillments(items []*shipping.Fulfillment, accType int) []*types.XFulfillment {
	result := make([]*types.XFulfillment, len(items))
	for i, item := range items {
		result[i] = Convert_core_Fulfillment_To_api_XFulfillment(item, accType, nil, nil)
	}
	return result
}

func FulfillmentSelfURL(ffm *shipping.Fulfillment, accType int) string {
	baseURL := cm.MainSiteBaseURL()
	switch accType {
	case account_tag.TagEtop:
		return ""

	case account_tag.TagShop:
		if baseURL == "" || ffm.ShopID == 0 || ffm.ID == 0 {
			return ""
		}
		return fmt.Sprintf("%v/s/%v/fulfillments/%v", baseURL, ffm.ShopID, ffm.ID)

	default:
		panic(fmt.Sprintf("unsupported account type: %v", accType))
	}
}

func Convert_core_FulfillmentSyncStates_To_api_FulfillmentSyncStates(s *shipping.FulfillmentSyncStates) *types.FulfillmentSyncStates {
	if s == nil {
		return nil
	}
	return &types.FulfillmentSyncStates{
		SyncAt:            cmapi.PbTime(s.SyncAt),
		NextShippingState: s.NextShippingState,
		Error:             cmapi.PbMetaError(s.Error),
	}
}

func Convert_core_ShippingFeeLine_To_api_ShippingFeeLine(line *shiptypes.ShippingFeeLine) *types.ShippingFeeLine {
	if line == nil {
		return nil
	}
	return &types.ShippingFeeLine{
		ShippingFeeType: line.ShippingFeeType,
		Cost:            line.Cost,
	}
}

func Convert_core_ShippingFeeLines_To_api_ShippingFeeLines(items []*shiptypes.ShippingFeeLine) []*types.ShippingFeeLine {
	result := make([]*types.ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = Convert_core_ShippingFeeLine_To_api_ShippingFeeLine(item)
	}
	return result
}

func Convert_api_ShippingFeeLine_To_core_ShippingFeeLine(in *types.ShippingFeeLine) *shiptypes.ShippingFeeLine {
	if in == nil {
		return nil
	}
	return &shiptypes.ShippingFeeLine{
		ShippingFeeType: in.ShippingFeeType,
		Cost:            in.Cost,
	}
}

func Convert_api_ShippingFeeLines_To_core_ShippingFeeLines(items []*types.ShippingFeeLine) []*shiptypes.ShippingFeeLine {
	if items == nil {
		return nil
	}
	result := make([]*shiptypes.ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = Convert_api_ShippingFeeLine_To_core_ShippingFeeLine(item)
	}
	return result
}

func Convert_core_ExternalShippingLog_To_api_ExternalShippingLog(log *shipping.ExternalShippingLog) *types.ExternalShippingLog {
	if log == nil {
		return nil
	}
	return &types.ExternalShippingLog{
		StateText: log.StateText,
		Time:      log.Time,
		Message:   log.Message,
	}
}

func Convert_core_ExternalShippingLogs_To_api_ExternalShippingLogs(items []*shipping.ExternalShippingLog) []*types.ExternalShippingLog {
	result := make([]*types.ExternalShippingLog, len(items))
	for i, item := range items {
		result[i] = Convert_core_ExternalShippingLog_To_api_ExternalShippingLog(item)
	}
	return result
}

func Convert_core_Order_To_api_Order(in *ordering.Order, ffms []*shipping.Fulfillment, accType int) *types.Order {
	if in == nil {
		return nil
	}

	fulfillments := Convert_core_Fulfillments_To_api_XFulfillments(ffms, accType)
	order := &types.Order{
		ExportedFields:            exportedOrder,
		Id:                        in.ID,
		ShopId:                    in.ShopID,
		ShopName:                  "",
		Code:                      in.Code,
		EdCode:                    in.EdCode,
		ExternalCode:              in.EdCode,
		Source:                    in.OrderSourceType,
		PartnerId:                 in.PartnerID,
		ExternalId:                in.ExternalOrderID,
		SelfUrl:                   OrderSelfURL(in.ID, in.ShopID, accType),
		PaymentMethod:             in.PaymentMethod,
		Customer:                  Convert_core_OrderCustomer_To_api_OrderCustomer(in.Customer),
		CustomerAddress:           Convert_core_OrderAddress_To_api_OrderAddress(in.CustomerAddress),
		BillingAddress:            Convert_core_OrderAddress_To_api_OrderAddress(in.BillingAddress),
		ShippingAddress:           Convert_core_OrderAddress_To_api_OrderAddress(in.ShippingAddress),
		CreatedAt:                 cmapi.PbTime(in.CreatedAt),
		CreatedBy:                 in.CreatedBy,
		ProcessedAt:               cmapi.PbTime(in.ProcessedAt),
		UpdatedAt:                 cmapi.PbTime(in.UpdatedAt),
		ClosedAt:                  cmapi.PbTime(in.ClosedAt),
		ConfirmedAt:               cmapi.PbTime(in.ConfirmedAt),
		CancelledAt:               cmapi.PbTime(in.CancelledAt),
		CancelReason:              in.CancelReason,
		ShConfirm:                 in.ShopConfirm,
		Confirm:                   in.ConfirmStatus,
		ConfirmStatus:             in.ConfirmStatus,
		Status:                    in.Status,
		FulfillmentStatus:         in.FulfillmentShippingStatus,
		FulfillmentShippingStatus: in.FulfillmentShippingStatus,
		EtopPaymentStatus:         in.EtopPaymentStatus,
		PaymentStatus:             in.PaymentStatus,
		Lines:                     Convert_core_OrderLines_To_api_OrderLines(in.Lines),
		Discounts:                 Convert_core_OrderDiscounts_To_api_OrderDiscounts(in.Discounts),
		TotalItems:                in.TotalItems,
		BasketValue:               in.BasketValue,
		TotalWeight:               in.TotalWeight,
		OrderDiscount:             in.OrderDiscount,
		TotalDiscount:             in.TotalDiscount,
		TotalAmount:               in.TotalAmount,
		OrderNote:                 in.OrderNote,
		ShippingFee:               in.ShopShippingFee,
		TotalFee:                  in.GetTotalFee(),
		FeeLines:                  Convert_core_OrderFeeLines_To_api_OrderFeeLines(in.FeeLines),
		ShopShippingFee:           in.ShopShippingFee,
		ShippingNote:              in.ShippingNote,
		ShopCod:                   in.ShopCOD,
		ReferenceUrl:              in.ReferenceURL,
		Fulfillments:              fulfillments,
		ShopShipping:              nil,
		Shipping:                  nil,
		GhnNoteCode:               in.GhnNoteCode,
		FulfillmentType:           in.FulfillmentType.String(),
		FulfillmentIds:            in.FulfillmentIDs,
		CustomerId:                in.CustomerID,
	}
	shipping := Convert_core_ShippingInfo_To_api_OrderShipping(in)
	order.ShopShipping = shipping
	order.Shipping = shipping
	return order
}

func OrderSelfURL(orderID dot.ID, shopID dot.ID, accType int) string {
	baseURL := cm.MainSiteBaseURL()
	switch accType {
	case account_tag.TagEtop:
		return ""

	case account_tag.TagShop:
		if baseURL == "" || shopID == 0 || orderID == 0 {
			return ""
		}
		return fmt.Sprintf("%v/s/%v/orders/%v", baseURL, shopID, orderID)

	default:
		panic(fmt.Sprintf("unsupported account type: %v", accType))
	}
}

func Convert_core_OrderCustomer_To_api_OrderCustomer(in *ordering.OrderCustomer) *types.OrderCustomer {
	if in == nil {
		return nil
	}
	_gender, _ := gender.ParseGender(in.Gender)
	return &types.OrderCustomer{
		ExportedFields: exportedOrderCustomer,

		FirstName: in.FirstName,
		LastName:  in.LastName,
		FullName:  in.GetFullName(),
		Email:     in.Email,
		Phone:     in.Phone,
		Gender:    _gender,
	}
}

func Convert_core_OrderDiscount_To_api_OrderDiscount(in *ordering.OrderDiscount) *types.OrderDiscount {
	if in == nil {
		return nil
	}
	return &types.OrderDiscount{
		Code:   in.Code,
		Type:   in.Type,
		Amount: in.Amount,
	}
}

func Convert_core_OrderDiscounts_To_api_OrderDiscounts(items []*ordering.OrderDiscount) []*types.OrderDiscount {
	result := make([]*types.OrderDiscount, len(items))
	for i, item := range items {
		result[i] = Convert_core_OrderDiscount_To_api_OrderDiscount(item)
	}
	return result
}

func Convert_core_OrderFeeLine_To_api_OrderFeeLine(in ordering.OrderFeeLine) *types.OrderFeeLine {
	return &types.OrderFeeLine{
		Type:   in.Type,
		Name:   in.Name,
		Code:   in.Code,
		Desc:   in.Desc,
		Amount: types.Int(in.Amount),
	}
}

func Convert_core_OrderFeeLines_To_api_OrderFeeLines(items []ordering.OrderFeeLine) []*types.OrderFeeLine {
	result := make([]*types.OrderFeeLine, len(items))
	for i, item := range items {
		result[i] = Convert_core_OrderFeeLine_To_api_OrderFeeLine(item)
	}
	return result
}

func Convert_core_ShippingInfo_To_api_OrderShipping(m *ordering.Order) *types.OrderShipping {
	if m == nil {
		return nil
	}
	in := m.Shipping
	if in == nil {
		in = &shiptypes.ShippingInfo{}
	}
	return &types.OrderShipping{
		ExportedFields: exportedOrderShipping,
		// @deprecated fields
		ShAddress:    Convert_core_OrderAddress_To_api_OrderAddress(in.PickupAddress),
		XServiceId:   in.ShippingServiceCode,
		XShippingFee: in.ShippingServiceFee,
		XServiceName: in.ShippingServiceName,

		PickupAddress:       Convert_core_OrderAddress_To_api_OrderAddress(in.PickupAddress),
		ShippingServiceName: in.ShippingServiceName,
		ShippingServiceCode: in.ShippingServiceCode,
		ShippingServiceFee:  in.ShippingServiceFee,
		ShippingProvider:    in.Carrier,
		Carrier:             in.Carrier,
		IncludeInsurance:    in.IncludeInsurance,
		TryOn:               in.TryOn,
		ShippingNote:        in.ShippingNote,
		CodAmount:           dot.Int(in.CODAmount),
		Weight:              dot.Int(cm.CoalesceInt(in.ChargeableWeight, m.TotalWeight)),
		GrossWeight:         dot.Int(cm.CoalesceInt(in.GrossWeight, m.TotalWeight)),
		Length:              dot.Int(in.Length),
		Width:               dot.Int(in.Width),
		Height:              dot.Int(in.Height),
		ChargeableWeight:    dot.Int(cm.CoalesceInt(in.ChargeableWeight, m.TotalWeight)),
	}
}
