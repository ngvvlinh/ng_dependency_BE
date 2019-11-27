package convertpb

import (
	"context"
	"time"

	"etop.vn/api/main/location"
	"etop.vn/api/pb/etop"
	"etop.vn/api/pb/etop/etc/gender"
	"etop.vn/api/pb/etop/order"
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

func PbOrdersWithFulfillments(items []ordermodelx.OrderWithFulfillments, accType int, shops []*model.Shop) []*order.Order {
	res := make([]*order.Order, len(items))
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

func PbOrders(items []*ordermodel.Order, accType int) []*order.Order {
	res := make([]*order.Order, len(items))
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

func PbOrder(m *ordermodel.Order, fulfillments []*shipmodel.Fulfillment, accType int) *order.Order {
	ffms := make([]*ordermodelx.Fulfillment, len(fulfillments))
	for i, ffm := range fulfillments {
		ffms[i] = &ordermodelx.Fulfillment{
			Shipment: ffm,
		}
	}

	order := &order.Order{
		ExportedFields: exportedOrder,

		Id:                        m.ID,
		ShopId:                    m.ShopID,
		ShopName:                  "",
		Code:                      m.Code,
		EdCode:                    m.EdCode,
		ExternalCode:              m.EdCode,
		Source:                    PbSource(m.OrderSourceType),
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
		ShConfirm:                 Pb3(m.ShopConfirm),
		Confirm:                   Pb3(m.ConfirmStatus),
		ConfirmStatus:             Pb3(m.ConfirmStatus),
		Status:                    Pb5(m.Status),
		FulfillmentStatus:         Pb5(m.FulfillmentShippingStatus),
		FulfillmentShippingStatus: Pb5(m.FulfillmentShippingStatus),
		EtopPaymentStatus:         Pb4(m.EtopPaymentStatus),
		PaymentStatus:             Pb4(m.PaymentStatus),
		Lines:                     PbOrderLines(m.Lines),
		Discounts:                 PbDiscounts(m.Discounts),
		TotalItems:                int(m.TotalItems),
		BasketValue:               int(m.BasketValue),
		TotalWeight:               int(m.TotalWeight),
		OrderDiscount:             int(m.OrderDiscount),
		TotalDiscount:             int(m.TotalDiscount),
		TotalAmount:               int(m.TotalAmount),
		OrderNote:                 m.OrderNote,
		ShippingFee:               int(m.ShopShippingFee),
		TotalFee:                  int(m.GetTotalFee()),
		FeeLines:                  PbOrderFeeLines(m.FeeLines),
		ShopShippingFee:           int(m.ShopShippingFee),
		ShippingNote:              m.ShippingNote,
		ShopCod:                   int(m.ShopCOD),
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

func XPbOrder(m *ordermodel.Order, fulfillments []*ordermodelx.Fulfillment, accType int) *order.Order {
	order := &order.Order{
		ExportedFields: exportedOrder,

		Id:                        m.ID,
		ShopId:                    m.ShopID,
		ShopName:                  "",
		Code:                      m.Code,
		EdCode:                    m.EdCode,
		ExternalCode:              m.EdCode,
		Source:                    PbSource(m.OrderSourceType),
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
		ShConfirm:                 Pb3(m.ShopConfirm),
		Confirm:                   Pb3(m.ConfirmStatus),
		ConfirmStatus:             Pb3(m.ConfirmStatus),
		Status:                    Pb5(m.Status),
		FulfillmentStatus:         Pb5(m.FulfillmentShippingStatus),
		FulfillmentShippingStatus: Pb5(m.FulfillmentShippingStatus),
		EtopPaymentStatus:         Pb4(m.EtopPaymentStatus),
		PaymentStatus:             Pb4(m.PaymentStatus),
		Lines:                     PbOrderLines(m.Lines),
		Discounts:                 PbDiscounts(m.Discounts),
		TotalItems:                int(m.TotalItems),
		BasketValue:               int(m.BasketValue),
		TotalWeight:               int(m.TotalWeight),
		OrderDiscount:             int(m.OrderDiscount),
		TotalDiscount:             int(m.TotalDiscount),
		TotalAmount:               int(m.TotalAmount),
		OrderNote:                 m.OrderNote,
		ShippingFee:               int(m.ShopShippingFee),
		TotalFee:                  int(m.GetTotalFee()),
		FeeLines:                  PbOrderFeeLines(m.FeeLines),
		ShopShippingFee:           int(m.ShopShippingFee),
		ShippingNote:              m.ShippingNote,
		ShopCod:                   int(m.ShopCOD),
		ReferenceUrl:              m.ReferenceURL,
		Fulfillments:              XPbFulfillments(fulfillments, accType),
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

var exportedOrderShipping = cm.SortStrings([]string{
	"pickup_address",
	"shipping_service_name", "shipping_service_code", "shipping_service_fee",
	"carrier", "try_on", "include_insurance", "shipping_node", "cod_amount",
	"weight", "gross_weight", "length", "width", "height", "chargeable_weight",
})

func PbOrderShipping(m *ordermodel.Order) *order.OrderShipping {
	if m == nil {
		return nil
	}
	item := m.ShopShipping
	if item == nil {
		item = &ordermodel.OrderShipping{}
	}
	return &order.OrderShipping{
		ExportedFields: exportedOrderShipping,
		// @deprecated fields
		ShAddress:    PbOrderAddress(item.ShopAddress),
		XServiceId:   cm.Coalesce(item.ProviderServiceID, item.ExternalServiceID),
		XShippingFee: int(item.ExternalShippingFee),
		XServiceName: item.ExternalServiceName,

		PickupAddress:       PbOrderAddress(item.GetPickupAddress()),
		ShippingServiceName: item.ExternalServiceName,
		ShippingServiceCode: item.GetShippingServiceCode(),
		ShippingServiceFee:  int(item.ExternalShippingFee),
		ShippingProvider:    PbShippingProviderType(item.ShippingProvider),
		Carrier:             PbShippingProviderType(item.ShippingProvider),
		IncludeInsurance:    item.IncludeInsurance,
		TryOn:               PbTryOn(m.GetTryOn()),
		ShippingNote:        m.ShippingNote,
		CodAmount:           cm.PInt(m.ShopCOD),
		Weight:              cm.PInt(m.TotalWeight),
		GrossWeight:         cm.PInt(m.TotalWeight),
		Length:              cm.PInt(item.Length),
		Width:               cm.PInt(item.Width),
		Height:              cm.PInt(item.Height),
		ChargeableWeight:    cm.PInt(m.TotalWeight),
	}
}

var exportedOrderCustomer = cm.SortStrings([]string{
	"full_name", "email", "phone", "gender",
})

func PbOrderCustomer(m *ordermodel.OrderCustomer) *order.OrderCustomer {
	if m == nil {
		return nil
	}
	return &order.OrderCustomer{
		ExportedFields: exportedOrderCustomer,

		FirstName: m.FirstName,
		LastName:  m.LastName,
		FullName:  m.GetFullName(),
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    gender.PbGender(m.Gender),
	}
}

func OrderCustomerToModel(m *order.OrderCustomer) *ordermodel.OrderCustomer {
	if m == nil {
		return nil
	}
	return &ordermodel.OrderCustomer{
		FirstName: m.FirstName,
		LastName:  m.LastName,
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    m.Gender.ToModel(),
	}
}

var exportedOrderAddress = cm.SortStrings([]string{
	"full_name", "phone", "province", "district",
	"ward", "company", "address1", "address2",
})

func PbOrderAddress(m *ordermodel.OrderAddress) *order.OrderAddress {
	if m == nil {
		return nil
	}
	res := &order.OrderAddress{
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

func OrderAddressToModel(m *order.OrderAddress) (*ordermodel.OrderAddress, error) {
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

func OrderAddressFulfilled(m *order.OrderAddress) (*order.OrderAddress, error) {
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

func PbOrderAddressFromAddress(m *model.Address) *order.OrderAddress {
	if m == nil {
		return nil
	}
	return &order.OrderAddress{
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

func OrderShippingToModel(m *order.OrderShipping, mo *ordermodel.Order) error {
	if m == nil {
		return nil
	}

	var pickupAddress *order.OrderAddress
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

	grossWeight := 0
	if m.Weight != nil {
		grossWeight = int(*m.Weight)
	}
	if m.GrossWeight != nil {
		grossWeight = int(*m.GrossWeight)
	}

	chargeableWeight := 0
	if m.ChargeableWeight != nil {
		chargeableWeight = int(*m.ChargeableWeight)

	} else {
		switch {
		case m.Length == nil && m.Width == nil && m.Height == nil:
			// continue
		case m.Length != nil && m.Width != nil && m.Height != nil:
			chargeableWeight = model.CalcChargeableWeight(
				grossWeight, int(*m.Length), int(*m.Width), int(*m.Height))

		default:
			return cm.Errorf(cm.InvalidArgument, err, "Cần cung cấp đủ các giá trị length, width, height (hoặc để trống cả 3)", err)
		}
	}

	carrierName := ShippingProviderToModel(&carrier)
	shippingServiceCode := cm.Coalesce(m.ShippingServiceCode, m.XServiceId)

	// check ETOP service
	shippingServiceName, ok := model.GetShippingServiceRegistry().GetName(model.TypeShippingETOP, shippingServiceCode)
	if !ok {
		shippingServiceName, ok = model.GetShippingServiceRegistry().GetName(carrierName, shippingServiceCode)
	}
	if carrierName != "" && !ok {
		return cm.Errorf(cm.InvalidArgument, err, "Mã dịch vụ không hợp lệ. Vui lòng F5 thử lại hoặc liên hệ hotro@etop.vn")
	}

	orderShipping := &ordermodel.OrderShipping{
		ShopAddress:         modelPickupAddress,
		ReturnAddress:       modelReturnAddress,
		ExternalServiceID:   shippingServiceCode,
		ExternalShippingFee: cm.CoalesceInt(int(m.ShippingServiceFee), int(m.XShippingFee)),
		ExternalServiceName: shippingServiceName,
		ShippingProvider:    carrierName,
		ProviderServiceID:   cm.Coalesce(shippingServiceCode, m.XServiceId),
		IncludeInsurance:    m.IncludeInsurance,
		Length:              cmapi.PatchInt(0, m.Length),
		Width:               cmapi.PatchInt(0, m.Width),
		Height:              cmapi.PatchInt(0, m.Height),
		GrossWeight:         grossWeight,
		ChargeableWeight:    chargeableWeight,
	}

	// when adding new fields here, remember to also change UpdateOrderCommand
	mo.ShopShipping = orderShipping
	mo.ShopCOD = cmapi.PatchInt(mo.ShopCOD, m.CodAmount)
	mo.TotalWeight = chargeableWeight

	if m.TryOn != 0 {
		mo.TryOn = TryOnCodeToModel(&m.TryOn)
		mo.GhnNoteCode = model.GHNNoteCodeFromTryOn(TryOnCodeToModel(&m.TryOn))
	} else if mo.GhnNoteCode != "" {
		mo.TryOn = model.TryOnFromGHNNoteCode(mo.GhnNoteCode)
	}

	// Coalesce takes from left to right while PatchInt takes from right
	mo.ShippingNote = cm.Coalesce(m.ShippingNote, mo.ShippingNote)
	return nil
}

func PbDiscounts(items []*ordermodel.OrderDiscount) []*order.OrderDiscount {
	res := make([]*order.OrderDiscount, len(items))
	for i, item := range items {
		res[i] = PbDiscount(item)
	}
	return res
}

func PbDiscount(m *ordermodel.OrderDiscount) *order.OrderDiscount {
	return &order.OrderDiscount{
		Code:   m.Code,
		Type:   m.Type,
		Amount: int(m.Amount),
	}
}

func OrderDiscountToModel(m *order.OrderDiscount) *ordermodel.OrderDiscount {
	return &ordermodel.OrderDiscount{
		Code:   m.Code,
		Type:   m.Type,
		Amount: int(m.Amount),
	}
}

func PbOrderDiscountsToModel(discounts []*order.OrderDiscount) []*ordermodel.OrderDiscount {
	res := make([]*ordermodel.OrderDiscount, len(discounts))
	for i, d := range discounts {
		res[i] = OrderDiscountToModel(d)
	}
	return res
}

func PbOrderFeeLinesToModel(items []*order.OrderFeeLine) []ordermodel.OrderFeeLine {
	res := make([]ordermodel.OrderFeeLine, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		res = append(res, ordermodel.OrderFeeLine{
			Amount: int(item.Amount),
			Desc:   item.Desc,
			Code:   item.Code,
			Name:   item.Name,
			Type:   FeeTypeToModel(&item.Type),
		})
	}
	return res
}

func PbOrderFeeLines(items []ordermodel.OrderFeeLine) []*order.OrderFeeLine {
	res := make([]*order.OrderFeeLine, len(items))
	for i, item := range items {
		res[i] = &order.OrderFeeLine{
			Type:   Pb(item.Type),
			Name:   item.Name,
			Code:   item.Code,
			Desc:   item.Desc,
			Amount: int(item.Amount),
		}
	}
	return res
}

func PbOrderLines(items []*ordermodel.OrderLine) []*order.OrderLine {
	res := make([]*order.OrderLine, 0, len(items))
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

func PbOrderLine(m *ordermodel.OrderLine) *order.OrderLine {
	metaFields := []*order.OrderLineMetaField{}
	for _, metaField := range m.MetaFields {
		metaFields = append(metaFields, &order.OrderLineMetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
			Name:  metaField.Name,
		})
	}
	if m == nil {
		return nil
	}
	return &order.OrderLine{
		ExportedFields: exportedOrderLine,

		OrderId:       m.OrderID,
		VariantId:     m.VariantID,
		ProductName:   m.ProductName,
		IsOutsideEtop: m.IsOutsideEtop,
		Quantity:      int(m.Quantity),
		ListPrice:     int(m.ListPrice),
		RetailPrice:   int(m.RetailPrice),
		PaymentPrice:  int(m.PaymentPrice),
		ImageUrl:      m.ImageURL,
		Attributes:    PbAttributesFromModel(m.Attributes),
		ProductId:     m.ProductID,
		TotalDiscount: int(m.TotalDiscount),
		MetaFields:    metaFields,
	}
}

func PbAttributesToModel(items []*order.Attribute) []*catalogmodel.ProductAttribute {
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

func AttributeToModel(m *order.Attribute) *catalogmodel.ProductAttribute {
	return &catalogmodel.ProductAttribute{
		Name:  m.Name,
		Value: m.Value,
	}
}

func PbAttributesFromModel(as []*catalogmodel.ProductAttribute) []*order.Attribute {
	attrs := make([]*order.Attribute, len(as))
	for i, a := range as {
		attrs[i] = &order.Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
}

func PbFulfillments(items []*shipmodel.Fulfillment, accType int) []*order.Fulfillment {
	if items == nil {
		return nil
	}
	res := make([]*order.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item, accType, nil, nil)
	}
	return res
}

func PbFulfillmentExtendeds(items []*shipmodely.FulfillmentExtended, accType int) []*order.Fulfillment {
	if items == nil {
		return nil
	}
	res := make([]*order.Fulfillment, len(items))
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

func PbFulfillment(m *shipmodel.Fulfillment, accType int, shop *model.Shop, mo *ordermodel.Order) *order.Fulfillment {
	if m == nil {
		return nil
	}
	ff := &order.Fulfillment{
		ExportedFields: exportedFulfillment,

		Id:                                 m.ID,
		OrderId:                            m.OrderID,
		ShopId:                             m.ShopID,
		PartnerId:                          m.PartnerID,
		SelfUrl:                            m.SelfURL(cm.MainSiteBaseURL(), accType),
		Lines:                              PbOrderLines(m.Lines),
		TotalItems:                         int(m.TotalItems),
		TotalWeight:                        int(m.TotalWeight),
		BasketValue:                        int(m.BasketValue),
		TotalCodAmount:                     int(m.TotalCODAmount),
		CodAmount:                          int(m.TotalCODAmount),
		TotalAmount:                        int(m.BasketValue), // deprecated
		ChargeableWeight:                   int(m.TotalWeight),
		CreatedAt:                          cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                          cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                           cmapi.PbTime(m.ClosedAt),
		CancelledAt:                        cmapi.PbTime(m.ShippingCancelledAt),
		CancelReason:                       m.CancelReason,
		ShippingProvider:                   string(m.ShippingProvider),
		Carrier:                            PbShippingProviderType(m.ShippingProvider),
		ShippingServiceName:                m.ExternalShippingName,
		ShippingServiceFee:                 int(m.ExternalShippingFee),
		ShippingServiceCode:                m.ProviderServiceID,
		ShippingCode:                       m.ShippingCode,
		ShippingNote:                       m.ShippingNote,
		TryOn:                              PbTryOn(m.TryOn),
		IncludeInsurance:                   m.IncludeInsurance,
		ShConfirm:                          Pb3(m.ShopConfirm),
		ShippingState:                      PbShippingState(m.ShippingState),
		Status:                             Pb5(m.Status),
		ShippingStatus:                     Pb5(m.ShippingStatus),
		EtopPaymentStatus:                  Pb4(m.EtopPaymentStatus),
		ShippingFeeCustomer:                int(m.ShippingFeeCustomer),
		ShippingFeeShop:                    int(m.ShippingFeeShop),
		XShippingFee:                       int(m.ExternalShippingFee),
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
		XShippingStatus:                    Pb5(m.ExternalShippingStatus),
		XSyncStatus:                        Pb4(m.SyncStatus),
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
		EtopDiscount:                       int(m.EtopDiscount),
		MoneyTransactionShippingId:         m.MoneyTransactionID,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		XShippingLogs:                      PbExternalShippingLogs(m.ExternalShippingLogs),
		XShippingNote:                      m.ExternalShippingNote,
		XShippingSubState:                  m.ExternalShippingSubState,
		ActualCompensationAmount:           int(m.ActualCompensationAmount),
	}
	if shop != nil {
		ff.Shop = PbShop(shop)
	}
	if mo != nil {
		ff.Order = PbOrder(mo, nil, accType)
	}
	return ff
}

func XPbFulfillments(items []*ordermodelx.Fulfillment, accType int) []*order.XFulfillment {
	if items == nil {
		return nil
	}
	res := make([]*order.XFulfillment, len(items))
	for i, item := range items {
		res[i] = XPbFulfillment(item, accType, nil, nil)
	}
	return res
}

func XPbFulfillment(m *ordermodelx.Fulfillment, accType int, shop *model.Shop, mo *ordermodel.Order) *order.XFulfillment {
	res := &order.XFulfillment{}
	shipment := PbFulfillment(m.Shipment, accType, shop, mo)
	if shipment != nil {
		res = &order.XFulfillment{
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

func PbAvailableShippingServices(items []*model.AvailableShippingService) []*order.ExternalShippingService {
	res := make([]*order.ExternalShippingService, len(items))
	for i, item := range items {
		res[i] = PbAvailableShippingService(item)
	}
	return res
}

var exportedShippingService = cm.SortStrings([]string{
	"name", "code", "fee", "carrier", "estimated_pickup_at", "estimated_delivery_at",
})

func PbAvailableShippingService(s *model.AvailableShippingService) *order.ExternalShippingService {
	return &order.ExternalShippingService{
		ExportedFields: exportedShippingService,

		ExternalId:          s.ProviderServiceID,
		ServiceFee:          int(s.ServiceFee),
		Provider:            PbShippingProviderType(s.Provider),
		ExpectedPickAt:      cmapi.PbTime(s.ExpectedPickAt),
		ExpectedDeliveryAt:  cmapi.PbTime(s.ExpectedDeliveryAt),
		Name:                s.Name,
		Code:                s.ProviderServiceID,
		Fee:                 int(s.ServiceFee),
		Carrier:             PbShippingProviderType(s.Provider),
		EstimatedPickupAt:   cmapi.PbTime(s.ExpectedPickAt),
		EstimatedDeliveryAt: cmapi.PbTime(s.ExpectedDeliveryAt),
	}
}

func PbShippingFeeLines(items []*model.ShippingFeeLine) []*order.ShippingFeeLine {
	result := make([]*order.ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = PbShippingFeeLine(item)
	}
	return result
}

func PbShippingFeeLine(line *model.ShippingFeeLine) *order.ShippingFeeLine {
	return &order.ShippingFeeLine{
		ShippingFeeType:          PbShippingFeeType(line.ShippingFeeType),
		Cost:                     int(line.Cost),
		ExternalServiceId:        line.ExternalServiceID,
		ExternalServiceName:      line.ExternalServiceName,
		ExternalServiceType:      line.ExternalServiceType,
		ExternalShippingOrderId:  line.ExternalShippingOrderID,
		ExternalPaymentChannelId: line.ExternalPaymentChannelID,
	}
}

func PbFulfillmentSyncStates(m *model.FulfillmentSyncStates) *order.FulfillmentSyncStates {
	if m == nil {
		return nil
	}
	return &order.FulfillmentSyncStates{
		SyncAt:            cmapi.PbTime(m.SyncAt),
		NextShippingState: string(m.NextShippingState),
		Error:             cmapi.PbError(m.Error),
	}
}

func PbMoneyTransactionExtended(m *txmodely.MoneyTransactionExtended) *order.MoneyTransaction {
	if m == nil {
		return nil
	}
	return &order.MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             Pb3(m.Status),
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

func PbMoneyTransaction(m *txmodel.MoneyTransactionShipping) *order.MoneyTransaction {
	if m == nil {
		return nil
	}
	return &order.MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             Pb3(m.Status),
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

func PbMoneyTransactionExtendeds(items []*txmodely.MoneyTransactionExtended) []*order.MoneyTransaction {
	result := make([]*order.MoneyTransaction, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalExtended(m *txmodel.MoneyTransactionShippingExternalExtended) *order.MoneyTransactionShippingExternal {
	if m == nil {
		return nil
	}
	res := &order.MoneyTransactionShippingExternal{
		Id:             m.ID,
		Code:           m.Code,
		TotalCod:       m.TotalCOD,
		TotalOrders:    m.TotalOrders,
		Status:         Pb3(m.Status),
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

func PbMoneyTransactionShippingExternalExtendeds(items []*txmodel.MoneyTransactionShippingExternalExtended) []*order.MoneyTransactionShippingExternal {
	result := make([]*order.MoneyTransactionShippingExternal, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalLineExtended(m *txmodel.MoneyTransactionShippingExternalLineExtended) *order.MoneyTransactionShippingExternalLine {
	if m == nil {
		return nil
	}
	res := &order.MoneyTransactionShippingExternalLine{
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

func PbMoneyTransactionShippingExternalLineExtendeds(items []*txmodel.MoneyTransactionShippingExternalLineExtended) []*order.MoneyTransactionShippingExternalLine {
	result := make([]*order.MoneyTransactionShippingExternalLine, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalLineExtended(item)
	}
	return result
}

func PbPublicFulfillment(item *shipmodel.Fulfillment) *order.PublicFulfillment {
	timeLayout := "15:04 02/01/2006"

	// use for manychat
	var deliveredAtText string
	if !item.ShippingDeliveredAt.IsZero() {
		deliveredAtText = item.ShippingDeliveredAt.In(time.Local).Format(timeLayout)
	} else {
		deliveredAtText = item.ExpectedDeliveryAt.In(time.Local).Format(timeLayout)
		deliveredAtText += " (dự kiến)"
	}

	return &order.PublicFulfillment{
		Id:                 item.ID,
		ShippingState:      PbShippingState(item.ShippingState),
		Status:             Pb5(item.Status),
		ExpectedDeliveryAt: cmapi.PbTime(item.ExpectedDeliveryAt),
		DeliveredAt:        cmapi.PbTime(item.ShippingDeliveredAt),
		ShippingCode:       item.ShippingCode,
		OrderId:            item.OrderID,
		DeliveredAtText:    deliveredAtText,
		ShippingStateText:  item.ShippingState.Text(),
	}
}

func PbExternalShippingLogs(items []*model.ExternalShippingLog) []*order.ExternalShippingLog {
	result := make([]*order.ExternalShippingLog, len(items))
	for i, item := range items {
		result[i] = PbExternalShippingLog(item)
	}
	return result
}

func PbExternalShippingLog(l *model.ExternalShippingLog) *order.ExternalShippingLog {
	return &order.ExternalShippingLog{
		StateText: l.StateText,
		Time:      l.Time,
		Message:   l.Message,
	}
}

func PbMoneyTransactionShippingEtopExtended(m *txmodely.MoneyTransactionShippingEtopExtended) *order.MoneyTransactionShippingEtop {
	if m == nil {
		return nil
	}
	return &order.MoneyTransactionShippingEtop{
		Id:                m.ID,
		Code:              m.Code,
		TotalCod:          m.TotalCOD,
		TotalOrders:       m.TotalOrders,
		TotalAmount:       m.TotalAmount,
		TotalFee:          m.TotalFee,
		Status:            Pb3(m.Status),
		MoneyTransactions: PbMoneyTransactionExtendeds(m.MoneyTransactions),
		CreatedAt:         cmapi.PbTime(m.CreatedAt),
		UpdatedAt:         cmapi.PbTime(m.UpdatedAt),
		ConfirmedAt:       cmapi.PbTime(m.ConfirmedAt),
		Note:              m.Note,
		InvoiceNumber:     m.InvoiceNumber,
		BankAccount:       PbBankAccount(m.BankAccount),
	}
}

func PbMoneyTransactionShippingEtopExtendeds(items []*txmodely.MoneyTransactionShippingEtopExtended) []*order.MoneyTransactionShippingEtop {
	result := make([]*order.MoneyTransactionShippingEtop, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingEtopExtended(item)
	}
	return result
}
