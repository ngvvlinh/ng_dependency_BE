package convertpb

import (
	"context"
	"time"

	"etop.vn/api/main/location"
	etop "etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/etc/gender"
	"etop.vn/api/top/types/etc/shipping_provider"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	servicelocation "etop.vn/backend/com/main/location"
	txmodel "etop.vn/backend/com/main/moneytx/model"
	txmodely "etop.vn/backend/com/main/moneytx/modely"
	orderconvert "etop.vn/backend/com/main/ordering/convert"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	shipmodely "etop.vn/backend/com/main/shipping/modely"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

var locationBus = servicelocation.New().MessageBus()

func PbOrdersWithFulfillments(items []ordermodelx.OrderWithFulfillments, accType int, shops []*model.Shop) []*types.Order {
	res := make([]*types.Order, len(items))
	shopsMap := make(map[dot.ID]*model.Shop)
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
		FulfillmentType:           orderconvert.Fulfill(m.FulfillmentType).String(),
		FulfillmentIds:            m.FulfillmentIDs,
		CustomerId:                m.CustomerID,
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
		FulfillmentType:           orderconvert.Fulfill(m.FulfillmentType).String(),
		FulfillmentIds:            m.FulfillmentIDs,
		CustomerId:                m.CustomerID,
		CreatedBy:                 m.CreatedBy,
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
		res.Coordinates = &ordermodel.Coordinates{
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

func PbOrderAddressFromAddress(m *model.Address) *types.OrderAddress {
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

func OrderShippingToModel(m *types.OrderShipping, mo *ordermodel.Order) error {
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

	// check ETOP service
	shippingServiceName, ok := model.GetShippingServiceRegistry().GetName(shipping_provider.Etop, shippingServiceCode)
	if !ok {
		shippingServiceName, ok = model.GetShippingServiceRegistry().GetName(carrier, shippingServiceCode)
	}
	if carrier != 0 && !ok {
		return cm.Errorf(cm.InvalidArgument, err, "Mã dịch vụ không hợp lệ. Vui lòng F5 thử lại hoặc liên hệ hotro@etop.vn")
	}

	orderShipping := &ordermodel.OrderShipping{
		ShopAddress:         modelPickupAddress,
		ReturnAddress:       modelReturnAddress,
		ExternalServiceID:   shippingServiceCode,
		ExternalShippingFee: cm.CoalesceInt(m.ShippingServiceFee, m.XShippingFee),
		ExternalServiceName: shippingServiceName,
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
		mo.GhnNoteCode = model.GHNNoteCodeFromTryOn(m.TryOn)
	} else if mo.GhnNoteCode != 0 {
		mo.TryOn = model.TryOnFromGHNNoteCode(mo.GhnNoteCode)
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
			Amount: item.Amount,
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
			Amount: item.Amount,
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
		Attributes:     PbAttributesFromModel(m.Attributes),
		ProductId:      m.ProductID,
		TotalDiscount:  m.TotalDiscount,
		MetaFields:     metaFields,
	}
}

func PbAttributesToModel(items []*types.Attribute) []*catalogmodel.ProductAttribute {
	if len(items) == 0 {
		return nil
	}
	res := make([]*catalogmodel.ProductAttribute, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		res = append(res, AttributeToModel(item))
	}
	return res
}

func AttributeToModel(m *types.Attribute) *catalogmodel.ProductAttribute {
	return &catalogmodel.ProductAttribute{
		Name:  m.Name,
		Value: m.Value,
	}
}

func PbAttributesFromModel(as []*catalogmodel.ProductAttribute) []*types.Attribute {
	attrs := make([]*types.Attribute, len(as))
	for i, a := range as {
		attrs[i] = &types.Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
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

func PbFulfillment(m *shipmodel.Fulfillment, accType int, shop *model.Shop, mo *ordermodel.Order) *types.Fulfillment {
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
		ChargeableWeight:                   m.TotalWeight,
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
		IncludeInsurance:                   m.IncludeInsurance,
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
		XShippingNote:                      m.ExternalShippingNote,
		XShippingSubState:                  m.ExternalShippingSubState,
		ActualCompensationAmount:           m.ActualCompensationAmount,
	}
	if shop != nil {
		ff.Shop = PbShop(shop)
	}
	if mo != nil {
		ff.Order = PbOrder(mo, nil, accType)
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

func XPbFulfillment(m *ordermodelx.Fulfillment, accType int, shop *model.Shop, mo *ordermodel.Order) *types.XFulfillment {
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
			TotalWeight:                        shipment.TotalWeight,
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

func PbAvailableShippingServices(items []*model.AvailableShippingService) []*types.ExternalShippingService {
	res := make([]*types.ExternalShippingService, len(items))
	for i, item := range items {
		res[i] = PbAvailableShippingService(item)
	}
	return res
}

var exportedShippingService = cm.SortStrings([]string{
	"name", "code", "fee", "carrier", "estimated_pickup_at", "estimated_delivery_at",
})

func PbAvailableShippingService(s *model.AvailableShippingService) *types.ExternalShippingService {
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
	}
}

func PbShippingFeeLines(items []*model.ShippingFeeLine) []*types.ShippingFeeLine {
	result := make([]*types.ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = PbShippingFeeLine(item)
	}
	return result
}

func PbShippingFeeLine(line *model.ShippingFeeLine) *types.ShippingFeeLine {
	return &types.ShippingFeeLine{
		ShippingFeeType:          line.ShippingFeeType,
		Cost:                     line.Cost,
		ExternalServiceId:        line.ExternalServiceID,
		ExternalServiceName:      line.ExternalServiceName,
		ExternalServiceType:      line.ExternalServiceType,
		ExternalShippingOrderId:  line.ExternalShippingOrderID,
		ExternalPaymentChannelId: line.ExternalPaymentChannelID,
	}
}

func PbFulfillmentSyncStates(m *model.FulfillmentSyncStates) *types.FulfillmentSyncStates {
	if m == nil {
		return nil
	}
	return &types.FulfillmentSyncStates{
		SyncAt:            cmapi.PbTime(m.SyncAt),
		NextShippingState: m.NextShippingState,
		Error:             cmapi.PbError(m.Error),
	}
}

func PbMoneyTransactionExtended(m *txmodely.MoneyTransactionExtended) *types.MoneyTransaction {
	if m == nil {
		return nil
	}
	return &types.MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             m.Status,
		TotalCod:                           m.TotalCOD,
		TotalOrders:                        m.TotalOrders,
		TotalAmount:                        m.TotalAmount,
		Code:                               m.Code,
		Provider:                           m.Provider,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		MoneyTransactionShippingEtopId:     m.MoneyTransactionShippingEtopID,
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:                        cmapi.PbTime(m.ConfirmedAt),
		EtopTransferedAt:                   cmapi.PbTime(m.EtopTransferedAt),
		Note:                               m.Note,
	}
}

func PbMoneyTransaction(m *txmodel.MoneyTransactionShipping) *types.MoneyTransaction {
	if m == nil {
		return nil
	}
	return &types.MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             m.Status,
		TotalCod:                           m.TotalCOD,
		TotalOrders:                        m.TotalOrders,
		Code:                               m.Code,
		Provider:                           m.Provider,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:                        cmapi.PbTime(m.ConfirmedAt),
		EtopTransferedAt:                   cmapi.PbTime(m.EtopTransferedAt),
	}
}

func PbMoneyTransactionExtendeds(items []*txmodely.MoneyTransactionExtended) []*types.MoneyTransaction {
	result := make([]*types.MoneyTransaction, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalExtended(m *txmodel.MoneyTransactionShippingExternalExtended) *types.MoneyTransactionShippingExternal {
	if m == nil {
		return nil
	}
	res := &types.MoneyTransactionShippingExternal{
		Id:             m.ID,
		Code:           m.Code,
		TotalCod:       m.TotalCOD,
		TotalOrders:    m.TotalOrders,
		Status:         m.Status,
		Provider:       m.Provider,
		CreatedAt:      cmapi.PbTime(m.CreatedAt),
		UpdatedAt:      cmapi.PbTime(m.UpdatedAt),
		ExternalPaidAt: cmapi.PbTime(m.ExternalPaidAt),
		Note:           m.Note,
		InvoiceNumber:  m.InvoiceNumber,
		BankAccount:    PbBankAccount(m.BankAccount),
		Lines:          PbMoneyTransactionShippingExternalLineExtendeds(m.Lines),
	}

	return res
}

func PbMoneyTransactionShippingExternalExtendeds(items []*txmodel.MoneyTransactionShippingExternalExtended) []*types.MoneyTransactionShippingExternal {
	result := make([]*types.MoneyTransactionShippingExternal, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalLineExtended(m *txmodel.MoneyTransactionShippingExternalLineExtended) *types.MoneyTransactionShippingExternalLine {
	if m == nil {
		return nil
	}
	res := &types.MoneyTransactionShippingExternalLine{
		Id:                                 m.ID,
		ExternalCode:                       m.ExternalCode,
		ExternalCustomer:                   m.ExternalCustomer,
		ExternalAddress:                    m.ExternalAddress,
		ExternalTotalCod:                   m.ExternalTotalCOD,
		ExternalTotalShippingFee:           m.ExternalTotalShippingFee,
		EtopFulfillmentId:                  m.EtopFulfillmentID,
		EtopFulfillmentIdRaw:               m.EtopFulfillmentIdRaw,
		Note:                               m.Note,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		ImportError:                        cmapi.PbCustomError(m.ImportError),
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ExternalCreatedAt:                  cmapi.PbTime(m.ExternalCreatedAt),
		ExternalClosedAt:                   cmapi.PbTime(m.ExternalClosedAt),
	}
	if m.Fulfillment != nil && m.Fulfillment.ID != 0 {
		res.Fulfillment = PbFulfillment(m.Fulfillment, model.TagEtop, m.Shop, m.Order)
	}
	return res
}

func PbMoneyTransactionShippingExternalLineExtendeds(items []*txmodel.MoneyTransactionShippingExternalLineExtended) []*types.MoneyTransactionShippingExternalLine {
	result := make([]*types.MoneyTransactionShippingExternalLine, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalLineExtended(item)
	}
	return result
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

func PbExternalShippingLogs(items []*model.ExternalShippingLog) []*types.ExternalShippingLog {
	result := make([]*types.ExternalShippingLog, len(items))
	for i, item := range items {
		result[i] = PbExternalShippingLog(item)
	}
	return result
}

func PbExternalShippingLog(l *model.ExternalShippingLog) *types.ExternalShippingLog {
	return &types.ExternalShippingLog{
		StateText: l.StateText,
		Time:      l.Time,
		Message:   l.Message,
	}
}

func PbMoneyTransactionShippingEtopExtended(m *txmodely.MoneyTransactionShippingEtopExtended) *types.MoneyTransactionShippingEtop {
	if m == nil {
		return nil
	}
	return &types.MoneyTransactionShippingEtop{
		Id:                m.ID,
		Code:              m.Code,
		TotalCod:          m.TotalCOD,
		TotalOrders:       m.TotalOrders,
		TotalAmount:       m.TotalAmount,
		TotalFee:          m.TotalFee,
		Status:            m.Status,
		MoneyTransactions: PbMoneyTransactionExtendeds(m.MoneyTransactions),
		CreatedAt:         cmapi.PbTime(m.CreatedAt),
		UpdatedAt:         cmapi.PbTime(m.UpdatedAt),
		ConfirmedAt:       cmapi.PbTime(m.ConfirmedAt),
		Note:              m.Note,
		InvoiceNumber:     m.InvoiceNumber,
		BankAccount:       PbBankAccount(m.BankAccount),
	}
}

func PbMoneyTransactionShippingEtopExtendeds(items []*txmodely.MoneyTransactionShippingEtopExtended) []*types.MoneyTransactionShippingEtop {
	result := make([]*types.MoneyTransactionShippingEtop, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingEtopExtended(item)
	}
	return result
}
