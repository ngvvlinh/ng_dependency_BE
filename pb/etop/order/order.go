package order

import (
	"context"
	"time"

	"etop.vn/api/main/location"
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
	"etop.vn/backend/pkg/etop/model"

	pbcm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	pbetop "etop.vn/backend/pb/etop"
	pbfee "etop.vn/backend/pb/etop/etc/fee"
	pbgender "etop.vn/backend/pb/etop/etc/gender"
	pbshipping "etop.vn/backend/pb/etop/etc/shipping"
	pbfeetype "etop.vn/backend/pb/etop/etc/shipping_fee_type"
	pbsp "etop.vn/backend/pb/etop/etc/shipping_provider"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbs5 "etop.vn/backend/pb/etop/etc/status5"
	pbtryon "etop.vn/backend/pb/etop/etc/try_on"
	pbsource "etop.vn/backend/pb/etop/order/source"
)

var locationBus = servicelocation.New().MessageBus()

func PbOrdersWithFulfillments(items []ordermodelx.OrderWithFulfillments, accType int, shops []*model.Shop) []*Order {
	res := make([]*Order, len(items))
	shopsMap := make(map[int64]*model.Shop)
	for _, shop := range shops {
		shopsMap[shop.ID] = shop
	}
	for i, item := range items {
		res[i] = XPbOrder(item.Order, item.Fulfillments, accType)
		res[i].ShopName = shopsMap[item.ShopID].GetShopName()
	}
	return res
}

func PbOrders(items []*ordermodel.Order, accType int) []*Order {
	res := make([]*Order, len(items))
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

func PbOrder(m *ordermodel.Order, fulfillments []*shipmodel.Fulfillment, accType int) *Order {
	ffms := make([]*ordermodelx.Fulfillment, len(fulfillments))
	for i, ffm := range fulfillments {
		ffms[i] = &ordermodelx.Fulfillment{
			Shipment: ffm,
		}
	}

	order := &Order{
		ExportedFields: exportedOrder,

		Id:                        m.ID,
		ShopId:                    m.ShopID,
		ShopName:                  "",
		Code:                      m.Code,
		EdCode:                    m.EdCode,
		ExternalCode:              m.EdCode,
		Source:                    pbsource.PbSource(m.OrderSourceType),
		PartnerId:                 m.PartnerID,
		ExternalId:                m.ExternalOrderID,
		ExternalUrl:               m.ExternalURL,
		SelfUrl:                   m.SelfURL(cm.MainSiteBaseURL(), accType),
		PaymentMethod:             m.PaymentMethod,
		Customer:                  PbOrderCustomer(m.Customer),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		BillingAddress:            PbOrderAddress(m.BillingAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 pbcm.PbTime(m.CreatedAt),
		ProcessedAt:               pbcm.PbTime(m.ProcessedAt),
		UpdatedAt:                 pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                  pbcm.PbTime(m.ClosedAt),
		ConfirmedAt:               pbcm.PbTime(m.ConfirmedAt),
		CancelledAt:               pbcm.PbTime(m.CancelledAt),
		CancelReason:              m.CancelReason,
		ShConfirm:                 pbs3.Pb(m.ShopConfirm),
		Confirm:                   pbs3.Pb(m.ConfirmStatus),
		ConfirmStatus:             pbs3.Pb(m.ConfirmStatus),
		Status:                    pbs5.Pb(m.Status),
		FulfillmentStatus:         pbs5.Pb(m.FulfillmentShippingStatus),
		FulfillmentShippingStatus: pbs5.Pb(m.FulfillmentShippingStatus),
		EtopPaymentStatus:         pbs4.Pb(m.EtopPaymentStatus),
		Lines:                     PbOrderLines(m.Lines),
		Discounts:                 PbDiscounts(m.Discounts),
		TotalItems:                int32(m.TotalItems),
		BasketValue:               int32(m.BasketValue),
		TotalWeight:               int32(m.TotalWeight),
		OrderDiscount:             int32(m.OrderDiscount),
		TotalDiscount:             int32(m.TotalDiscount),
		TotalAmount:               int32(m.TotalAmount),
		OrderNote:                 m.OrderNote,
		ShippingFee:               int32(m.ShopShippingFee),
		TotalFee:                  int32(m.GetTotalFee()),
		FeeLines:                  PbOrderFeeLines(m.FeeLines),
		ShopShippingFee:           int32(m.ShopShippingFee),
		ShippingNote:              m.ShippingNote,
		ShopCod:                   int32(m.ShopCOD),
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

func XPbOrder(m *ordermodel.Order, fulfillments []*ordermodelx.Fulfillment, accType int) *Order {
	order := &Order{
		ExportedFields: exportedOrder,

		Id:                        m.ID,
		ShopId:                    m.ShopID,
		ShopName:                  "",
		Code:                      m.Code,
		EdCode:                    m.EdCode,
		ExternalCode:              m.EdCode,
		Source:                    pbsource.PbSource(m.OrderSourceType),
		PartnerId:                 m.PartnerID,
		ExternalId:                m.ExternalOrderID,
		ExternalUrl:               m.ExternalURL,
		SelfUrl:                   m.SelfURL(cm.MainSiteBaseURL(), accType),
		PaymentMethod:             m.PaymentMethod,
		Customer:                  PbOrderCustomer(m.Customer),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		BillingAddress:            PbOrderAddress(m.BillingAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 pbcm.PbTime(m.CreatedAt),
		ProcessedAt:               pbcm.PbTime(m.ProcessedAt),
		UpdatedAt:                 pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                  pbcm.PbTime(m.ClosedAt),
		ConfirmedAt:               pbcm.PbTime(m.ConfirmedAt),
		CancelledAt:               pbcm.PbTime(m.CancelledAt),
		CancelReason:              m.CancelReason,
		ShConfirm:                 pbs3.Pb(m.ShopConfirm),
		Confirm:                   pbs3.Pb(m.ConfirmStatus),
		ConfirmStatus:             pbs3.Pb(m.ConfirmStatus),
		Status:                    pbs5.Pb(m.Status),
		FulfillmentStatus:         pbs5.Pb(m.FulfillmentShippingStatus),
		FulfillmentShippingStatus: pbs5.Pb(m.FulfillmentShippingStatus),
		EtopPaymentStatus:         pbs4.Pb(m.EtopPaymentStatus),
		Lines:                     PbOrderLines(m.Lines),
		Discounts:                 PbDiscounts(m.Discounts),
		TotalItems:                int32(m.TotalItems),
		BasketValue:               int32(m.BasketValue),
		TotalWeight:               int32(m.TotalWeight),
		OrderDiscount:             int32(m.OrderDiscount),
		TotalDiscount:             int32(m.TotalDiscount),
		TotalAmount:               int32(m.TotalAmount),
		OrderNote:                 m.OrderNote,
		ShippingFee:               int32(m.ShopShippingFee),
		TotalFee:                  int32(m.GetTotalFee()),
		FeeLines:                  PbOrderFeeLines(m.FeeLines),
		ShopShippingFee:           int32(m.ShopShippingFee),
		ShippingNote:              m.ShippingNote,
		ShopCod:                   int32(m.ShopCOD),
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

func PbOrderShipping(order *ordermodel.Order) *OrderShipping {
	if order == nil {
		return nil
	}
	item := order.ShopShipping
	if item == nil {
		item = &ordermodel.OrderShipping{}
	}
	return &OrderShipping{
		ExportedFields: exportedOrderShipping,
		// @deprecated fields
		ShAddress:    PbOrderAddress(item.ShopAddress),
		XServiceId:   cm.Coalesce(item.ProviderServiceID, item.ExternalServiceID),
		XShippingFee: int32(item.ExternalShippingFee),
		XServiceName: item.ExternalServiceName,

		PickupAddress:       PbOrderAddress(item.GetPickupAddress()),
		ShippingServiceName: item.ExternalServiceName,
		ShippingServiceCode: item.GetShippingServiceCode(),
		ShippingServiceFee:  int32(item.ExternalShippingFee),
		ShippingProvider:    pbsp.PbShippingProviderType(item.ShippingProvider),
		Carrier:             pbsp.PbShippingProviderType(item.ShippingProvider),
		IncludeInsurance:    item.IncludeInsurance,
		TryOn:               pbtryon.PbTryOn(order.GetTryOn()),
		ShippingNote:        order.ShippingNote,
		CodAmount:           pbcm.PbPtrInt32(order.ShopCOD),
		Weight:              pbcm.PbPtrInt32(order.TotalWeight),
		GrossWeight:         pbcm.PbPtrInt32(order.TotalWeight),
		Length:              pbcm.PbPtrInt32(item.Length),
		Width:               pbcm.PbPtrInt32(item.Width),
		Height:              pbcm.PbPtrInt32(item.Height),
		ChargeableWeight:    pbcm.PbPtrInt32(order.TotalWeight),
	}
}

var exportedOrderCustomer = cm.SortStrings([]string{
	"full_name", "email", "phone", "gender",
})

func PbOrderCustomer(m *ordermodel.OrderCustomer) *OrderCustomer {
	if m == nil {
		return nil
	}
	return &OrderCustomer{
		ExportedFields: exportedOrderCustomer,

		FirstName: m.FirstName,
		LastName:  m.LastName,
		FullName:  m.GetFullName(),
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    pbgender.PbGender(m.Gender),
	}
}

func (m *OrderCustomer) ToModel() *ordermodel.OrderCustomer {
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

func PbOrderAddress(m *ordermodel.OrderAddress) *OrderAddress {
	if m == nil {
		return nil
	}
	res := &OrderAddress{
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

func (m *OrderAddress) ToModel() (*ordermodel.OrderAddress, error) {
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

func (m *OrderAddress) Fulfilled() (*OrderAddress, error) {
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

func PbOrderAddressFromAddress(m *model.Address) *OrderAddress {
	if m == nil {
		return nil
	}
	return &OrderAddress{
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

func (m *OrderShipping) ToModel(order *ordermodel.Order) error {
	if m == nil {
		return nil
	}

	var pickupAddress *OrderAddress
	if m.ShAddress != nil {
		pickupAddress = m.ShAddress
	} else {
		pickupAddress = m.PickupAddress
	}
	modelPickupAddress, err := pickupAddress.ToModel()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ lấy hàng không hợp lệ: %v", err)
	}
	modelReturnAddress, err := m.ReturnAddress.ToModel()
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

	carrierName := carrier.ToModel()
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
		Length:              pbcm.PatchInt32(0, m.Length),
		Width:               pbcm.PatchInt32(0, m.Width),
		Height:              pbcm.PatchInt32(0, m.Height),
		GrossWeight:         grossWeight,
		ChargeableWeight:    chargeableWeight,
	}

	// when adding new fields here, remember to also change UpdateOrderCommand
	order.ShopShipping = orderShipping
	order.ShopCOD = pbcm.PatchInt32(order.ShopCOD, m.CodAmount)
	order.TotalWeight = chargeableWeight

	if m.TryOn != 0 {
		order.TryOn = m.TryOn.ToModel()
		order.GhnNoteCode = model.GHNNoteCodeFromTryOn(m.TryOn.ToModel())
	} else if order.GhnNoteCode != "" {
		order.TryOn = model.TryOnFromGHNNoteCode(order.GhnNoteCode)
	}

	// Coalesce takes from left to right while PatchInt takes from right
	order.ShippingNote = cm.Coalesce(m.ShippingNote, order.ShippingNote)
	return nil
}

func PbDiscounts(items []*ordermodel.OrderDiscount) []*OrderDiscount {
	res := make([]*OrderDiscount, len(items))
	for i, item := range items {
		res[i] = PbDiscount(item)
	}
	return res
}

func PbDiscount(m *ordermodel.OrderDiscount) *OrderDiscount {
	return &OrderDiscount{
		Code:   m.Code,
		Type:   m.Type,
		Amount: int32(m.Amount),
	}
}

func (m *OrderDiscount) ToModel() *ordermodel.OrderDiscount {
	return &ordermodel.OrderDiscount{
		Code:   m.Code,
		Type:   m.Type,
		Amount: int(m.Amount),
	}
}

func PbOrderDiscountsToModel(discounts []*OrderDiscount) []*ordermodel.OrderDiscount {
	res := make([]*ordermodel.OrderDiscount, len(discounts))
	for i, d := range discounts {
		res[i] = d.ToModel()
	}
	return res
}

func PbOrderFeeLinesToModel(items []*OrderFeeLine) []ordermodel.OrderFeeLine {
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
			Type:   item.Type.ToModel(),
		})
	}
	return res
}

func PbOrderFeeLines(items []ordermodel.OrderFeeLine) []*OrderFeeLine {
	res := make([]*OrderFeeLine, len(items))
	for i, item := range items {
		res[i] = &OrderFeeLine{
			Type:   pbfee.Pb(item.Type),
			Name:   item.Name,
			Code:   item.Code,
			Desc:   item.Desc,
			Amount: int32(item.Amount),
		}
	}
	return res
}

func PbOrderLines(items []*ordermodel.OrderLine) []*OrderLine {
	res := make([]*OrderLine, 0, len(items))
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

func PbOrderLine(m *ordermodel.OrderLine) *OrderLine {
	if m == nil {
		return nil
	}
	return &OrderLine{
		ExportedFields: exportedOrderLine,

		OrderId:       m.OrderID,
		VariantId:     m.VariantID,
		ProductName:   m.ProductName,
		IsOutsideEtop: m.IsOutsideEtop,
		Quantity:      int32(m.Quantity),
		ListPrice:     int32(m.ListPrice),
		RetailPrice:   int32(m.RetailPrice),
		PaymentPrice:  int32(m.PaymentPrice),
		ImageUrl:      m.ImageURL,
		Attributes:    PbAttributes(m.Attributes),
		ProductId:     m.ProductID,
		TotalDiscount: int32(m.TotalDiscount),
	}
}

func PbAttributesToModel(items []*Attribute) []*catalogmodel.ProductAttribute {
	if len(items) == 0 {
		return nil
	}
	res := make([]*catalogmodel.ProductAttribute, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		res = append(res, item.ToModel())
	}
	return res
}

func (m *Attribute) ToModel() *catalogmodel.ProductAttribute {
	return &catalogmodel.ProductAttribute{
		Name:  m.Name,
		Value: m.Value,
	}
}

func PbAttributes(as []*catalogmodel.ProductAttribute) []*Attribute {
	attrs := make([]*Attribute, len(as))
	for i, a := range as {
		attrs[i] = &Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
}

func PbFulfillments(items []*shipmodel.Fulfillment, accType int) []*Fulfillment {
	if items == nil {
		return nil
	}
	res := make([]*Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item, accType, nil, nil)
	}
	return res
}

func PbFulfillmentExtendeds(items []*shipmodely.FulfillmentExtended, accType int) []*Fulfillment {
	if items == nil {
		return nil
	}
	res := make([]*Fulfillment, len(items))
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

func PbFulfillment(m *shipmodel.Fulfillment, accType int, shop *model.Shop, order *ordermodel.Order) *Fulfillment {
	if m == nil {
		return nil
	}
	ff := &Fulfillment{
		ExportedFields: exportedFulfillment,

		Id:                                 m.ID,
		OrderId:                            m.OrderID,
		ShopId:                             m.ShopID,
		PartnerId:                          m.PartnerID,
		SelfUrl:                            m.SelfURL(cm.MainSiteBaseURL(), accType),
		Lines:                              PbOrderLines(m.Lines),
		TotalItems:                         int32(m.TotalItems),
		TotalWeight:                        int32(m.TotalWeight),
		BasketValue:                        int32(m.BasketValue),
		TotalCodAmount:                     int32(m.TotalCODAmount),
		CodAmount:                          int32(m.TotalCODAmount),
		TotalAmount:                        int32(m.BasketValue), // deprecated
		ChargeableWeight:                   int32(m.TotalWeight),
		CreatedAt:                          pbcm.PbTime(m.CreatedAt),
		UpdatedAt:                          pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                           pbcm.PbTime(m.ClosedAt),
		CancelledAt:                        pbcm.PbTime(m.ShippingCancelledAt),
		CancelReason:                       m.CancelReason,
		ShippingProvider:                   string(m.ShippingProvider),
		Carrier:                            pbsp.PbShippingProviderType(m.ShippingProvider),
		ShippingServiceName:                m.ExternalShippingName,
		ShippingServiceFee:                 int32(m.ExternalShippingFee),
		ShippingServiceCode:                m.ProviderServiceID,
		ShippingCode:                       m.ShippingCode,
		ShippingNote:                       m.ShippingNote,
		TryOn:                              pbtryon.PbTryOn(m.TryOn),
		IncludeInsurance:                   m.IncludeInsurance,
		ShConfirm:                          pbs3.Pb(m.ShopConfirm),
		ShippingState:                      pbshipping.Pb(m.ShippingState),
		Status:                             pbs5.Pb(m.Status),
		ShippingStatus:                     pbs5.Pb(m.ShippingStatus),
		EtopPaymentStatus:                  pbs4.Pb(m.EtopPaymentStatus),
		ShippingFeeCustomer:                int32(m.ShippingFeeCustomer),
		ShippingFeeShop:                    int32(m.ShippingFeeShop),
		XShippingFee:                       int32(m.ExternalShippingFee),
		XShippingId:                        m.ExternalShippingID,
		XShippingCode:                      m.ExternalShippingCode,
		XShippingCreatedAt:                 pbcm.PbTime(m.ExternalShippingCreatedAt),
		XShippingUpdatedAt:                 pbcm.PbTime(m.ExternalShippingUpdatedAt),
		XShippingCancelledAt:               pbcm.PbTime(m.ExternalShippingCancelledAt),
		XShippingDeliveredAt:               pbcm.PbTime(m.ExternalShippingDeliveredAt),
		XShippingReturnedAt:                pbcm.PbTime(m.ExternalShippingReturnedAt),
		ExpectedDeliveryAt:                 pbcm.PbTime(m.ExpectedDeliveryAt),
		ExpectedPickAt:                     pbcm.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt:                pbcm.PbTime(m.ExpectedDeliveryAt),
		EstimatedPickupAt:                  pbcm.PbTime(m.ExpectedPickAt),
		CodEtopTransferedAt:                pbcm.PbTime(m.CODEtopTransferedAt),
		ShippingFeeShopTransferedAt:        pbcm.PbTime(m.ShippingFeeShopTransferedAt),
		XShippingState:                     m.ExternalShippingState,
		XShippingStatus:                    pbs5.Pb(m.ExternalShippingStatus),
		XSyncStatus:                        pbs4.Pb(m.SyncStatus),
		XSyncStates:                        PbFulfillmentSyncStates(m.SyncStates),
		AddressTo:                          pbetop.PbAddress(m.AddressTo),
		AddressFrom:                        pbetop.PbAddress(m.AddressFrom),
		PickupAddress:                      PbOrderAddressFromAddress(m.AddressFrom),
		ReturnAddress:                      PbOrderAddressFromAddress(m.AddressFrom),
		ShippingAddress:                    PbOrderAddressFromAddress(m.AddressTo),
		Shop:                               nil,
		Order:                              nil,
		ProviderShippingFeeLines:           PbShippingFeeLines(m.ProviderShippingFeeLines),
		ShippingFeeShopLines:               PbShippingFeeLines(m.ShippingFeeShopLines),
		EtopDiscount:                       int32(m.EtopDiscount),
		MoneyTransactionShippingId:         m.MoneyTransactionID,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		XShippingLogs:                      PbExternalShippingLogs(m.ExternalShippingLogs),
		XShippingNote:                      m.ExternalShippingNote,
		XShippingSubState:                  m.ExternalShippingSubState,
		ActualCompensationAmount:           int32(m.ActualCompensationAmount),
	}
	if shop != nil {
		ff.Shop = pbetop.PbShop(shop)
	}
	if order != nil {
		ff.Order = PbOrder(order, nil, accType)
	}
	return ff
}

func XPbFulfillments(items []*ordermodelx.Fulfillment, accType int) []*XFulfillment {
	if items == nil {
		return nil
	}
	res := make([]*XFulfillment, len(items))
	for i, item := range items {
		res[i] = XPbFulfillment(item, accType, nil, nil)
	}
	return res
}

func XPbFulfillment(m *ordermodelx.Fulfillment, accType int, shop *model.Shop, order *ordermodel.Order) *XFulfillment {
	res := &XFulfillment{}
	shipment := PbFulfillment(m.Shipment, accType, shop, order)
	if shipment != nil {
		res = &XFulfillment{
			Fulfill: &XFulfillment_Shipment{
				Shipment: shipment,
			},
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
		res.Fulfill = &XFulfillment_Shipnow{
			Shipnow: Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(m.Shipnow),
		}
	}

	return res
}

func PbAvailableShippingServices(items []*model.AvailableShippingService) []*ExternalShippingService {
	res := make([]*ExternalShippingService, len(items))
	for i, item := range items {
		res[i] = PbAvailableShippingService(item)
	}
	return res
}

var exportedShippingService = cm.SortStrings([]string{
	"name", "code", "fee", "carrier", "estimated_pickup_at", "estimated_delivery_at",
})

func PbAvailableShippingService(s *model.AvailableShippingService) *ExternalShippingService {
	return &ExternalShippingService{
		ExportedFields: exportedShippingService,

		ExternalId:          s.ProviderServiceID,
		ServiceFee:          int32(s.ServiceFee),
		Provider:            pbsp.PbShippingProviderType(s.Provider),
		ExpectedPickAt:      pbcm.PbTime(s.ExpectedPickAt),
		ExpectedDeliveryAt:  pbcm.PbTime(s.ExpectedDeliveryAt),
		Name:                s.Name,
		Code:                s.ProviderServiceID,
		Fee:                 int32(s.ServiceFee),
		Carrier:             pbsp.PbShippingProviderType(s.Provider),
		EstimatedPickupAt:   pbcm.PbTime(s.ExpectedPickAt),
		EstimatedDeliveryAt: pbcm.PbTime(s.ExpectedDeliveryAt),
	}
}

func PbShippingFeeLines(items []*model.ShippingFeeLine) []*ShippingFeeLine {
	result := make([]*ShippingFeeLine, len(items))
	for i, item := range items {
		result[i] = PbShippingFeeLine(item)
	}
	return result
}

func PbShippingFeeLine(line *model.ShippingFeeLine) *ShippingFeeLine {
	return &ShippingFeeLine{
		ShippingFeeType:          pbfeetype.Pb(line.ShippingFeeType),
		Cost:                     int32(line.Cost),
		ExternalServiceId:        line.ExternalServiceID,
		ExternalServiceName:      line.ExternalServiceName,
		ExternalServiceType:      line.ExternalServiceType,
		ExternalShippingOrderId:  line.ExternalShippingOrderID,
		ExternalPaymentChannelId: line.ExternalPaymentChannelID,
	}
}

func PbFulfillmentSyncStates(m *model.FulfillmentSyncStates) *FulfillmentSyncStates {
	if m == nil {
		return nil
	}
	return &FulfillmentSyncStates{
		SyncAt:            pbcm.PbTime(m.SyncAt),
		NextShippingState: string(m.NextShippingState),
		Error:             pbcm.PbError(m.Error),
	}
}

func PbMoneyTransactionExtended(m *txmodely.MoneyTransactionExtended) *MoneyTransaction {
	if m == nil {
		return nil
	}
	return &MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             pbs3.Pb(m.Status),
		TotalCod:                           int64(m.TotalCOD),
		TotalOrders:                        int64(m.TotalOrders),
		TotalAmount:                        int64(m.TotalAmount),
		Code:                               m.Code,
		Provider:                           m.Provider,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		MoneyTransactionShippingEtopId:     m.MoneyTransactionShippingEtopID,
		CreatedAt:                          pbcm.PbTime(m.CreatedAt),
		UpdatedAt:                          pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                           pbcm.PbTime(m.ClosedAt),
		ConfirmedAt:                        pbcm.PbTime(m.ConfirmedAt),
		EtopTransferedAt:                   pbcm.PbTime(m.EtopTransferedAt),
		Note:                               m.Note,
	}
}

func PbMoneyTransaction(m *txmodel.MoneyTransactionShipping) *MoneyTransaction {
	if m == nil {
		return nil
	}
	return &MoneyTransaction{
		Id:                                 m.ID,
		ShopId:                             m.ShopID,
		Status:                             pbs3.Pb(m.Status),
		TotalCod:                           int64(m.TotalCOD),
		TotalOrders:                        int64(m.TotalOrders),
		Code:                               m.Code,
		Provider:                           m.Provider,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		CreatedAt:                          pbcm.PbTime(m.CreatedAt),
		UpdatedAt:                          pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                           pbcm.PbTime(m.ClosedAt),
		ConfirmedAt:                        pbcm.PbTime(m.ConfirmedAt),
		EtopTransferedAt:                   pbcm.PbTime(m.EtopTransferedAt),
	}
}

func PbMoneyTransactionExtendeds(items []*txmodely.MoneyTransactionExtended) []*MoneyTransaction {
	result := make([]*MoneyTransaction, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalExtended(m *txmodel.MoneyTransactionShippingExternalExtended) *MoneyTransactionShippingExternal {
	if m == nil {
		return nil
	}
	res := &MoneyTransactionShippingExternal{
		Id:             m.ID,
		Code:           m.Code,
		TotalCod:       int64(m.TotalCOD),
		TotalOrders:    int64(m.TotalOrders),
		Status:         pbs3.Pb(m.Status),
		Provider:       m.Provider,
		CreatedAt:      pbcm.PbTime(m.CreatedAt),
		UpdatedAt:      pbcm.PbTime(m.UpdatedAt),
		ExternalPaidAt: pbcm.PbTime(m.ExternalPaidAt),
		Note:           m.Note,
		InvoiceNumber:  m.InvoiceNumber,
		BankAccount:    etop.PbBankAccount(m.BankAccount),
		Lines:          PbMoneyTransactionShippingExternalLineExtendeds(m.Lines),
	}

	return res
}

func PbMoneyTransactionShippingExternalExtendeds(items []*txmodel.MoneyTransactionShippingExternalExtended) []*MoneyTransactionShippingExternal {
	result := make([]*MoneyTransactionShippingExternal, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalExtended(item)
	}
	return result
}

func PbMoneyTransactionShippingExternalLineExtended(m *txmodel.MoneyTransactionShippingExternalLineExtended) *MoneyTransactionShippingExternalLine {
	if m == nil {
		return nil
	}
	res := &MoneyTransactionShippingExternalLine{
		Id:                                 m.ID,
		ExternalCode:                       m.ExternalCode,
		ExternalCustomer:                   m.ExternalCustomer,
		ExternalAddress:                    m.ExternalAddress,
		ExternalTotalCod:                   int64(m.ExternalTotalCOD),
		ExternalTotalShippingFee:           int64(m.ExternalTotalShippingFee),
		EtopFulfillmentId:                  m.EtopFulfillmentID,
		EtopFulfillmentIdRaw:               m.EtopFulfillmentIdRaw,
		Note:                               m.Note,
		MoneyTransactionShippingExternalId: m.MoneyTransactionShippingExternalID,
		ImportError:                        pbcm.PbCustomError(m.ImportError),
		CreatedAt:                          pbcm.PbTime(m.CreatedAt),
		UpdatedAt:                          pbcm.PbTime(m.UpdatedAt),
		ExternalCreatedAt:                  pbcm.PbTime(m.ExternalCreatedAt),
		ExternalClosedAt:                   pbcm.PbTime(m.ExternalClosedAt),
	}
	if m.Fulfillment != nil && m.Fulfillment.ID != 0 {
		res.Fulfillment = PbFulfillment(m.Fulfillment, model.TagEtop, m.Shop, m.Order)
	}
	return res
}

func PbMoneyTransactionShippingExternalLineExtendeds(items []*txmodel.MoneyTransactionShippingExternalLineExtended) []*MoneyTransactionShippingExternalLine {
	result := make([]*MoneyTransactionShippingExternalLine, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingExternalLineExtended(item)
	}
	return result
}

func (m *OrderWithErrorsResponse) HasErrors() []*pbcm.Error {
	return m.FulfillmentErrors
}

func (m *ImportOrdersResponse) HasErrors() []*pbcm.Error {
	if len(m.CellErrors) > 0 {
		return m.CellErrors
	}
	return m.ImportErrors
}

func PbPublicFulfillment(item *shipmodel.Fulfillment) *PublicFulfillment {
	timeLayout := "15:04 02/01/2006"

	// use for manychat
	var deliveredAtText string
	if !item.ShippingDeliveredAt.IsZero() {
		deliveredAtText = item.ShippingDeliveredAt.In(time.Local).Format(timeLayout)
	} else {
		deliveredAtText = item.ExpectedDeliveryAt.In(time.Local).Format(timeLayout)
		deliveredAtText += " (dự kiến)"
	}

	return &PublicFulfillment{
		Id:                  item.ID,
		ShippingState:       pbshipping.Pb(item.ShippingState),
		Status:              pbs5.Pb(item.Status),
		ExpectedDeliveryAt:  pbcm.PbTime(item.ExpectedDeliveryAt),
		DeliveredAt:         pbcm.PbTime(item.ShippingDeliveredAt),
		EstimatedPickupAt:   nil,
		EstimatedDeliveryAt: nil,
		ShippingCode:        item.ShippingCode,
		OrderId:             item.OrderID,
		DeliveredAtText:     deliveredAtText,
		ShippingStateText:   item.ShippingState.Text(),
	}
}

func PbExternalShippingLogs(items []*model.ExternalShippingLog) []*ExternalShippingLog {
	result := make([]*ExternalShippingLog, len(items))
	for i, item := range items {
		result[i] = PbExternalShippingLog(item)
	}
	return result
}

func PbExternalShippingLog(l *model.ExternalShippingLog) *ExternalShippingLog {
	return &ExternalShippingLog{
		StateText: l.StateText,
		Time:      l.Time,
		Message:   l.Message,
	}
}

func PbMoneyTransactionShippingEtopExtended(m *txmodely.MoneyTransactionShippingEtopExtended) *MoneyTransactionShippingEtop {
	if m == nil {
		return nil
	}
	return &MoneyTransactionShippingEtop{
		Id:                m.ID,
		Code:              m.Code,
		TotalCod:          int64(m.TotalCOD),
		TotalOrders:       int64(m.TotalOrders),
		TotalAmount:       int64(m.TotalAmount),
		TotalFee:          int64(m.TotalFee),
		Status:            pbs3.Pb(m.Status),
		MoneyTransactions: PbMoneyTransactionExtendeds(m.MoneyTransactions),
		CreatedAt:         pbcm.PbTime(m.CreatedAt),
		UpdatedAt:         pbcm.PbTime(m.UpdatedAt),
		ConfirmedAt:       pbcm.PbTime(m.ConfirmedAt),
		Note:              m.Note,
		InvoiceNumber:     m.InvoiceNumber,
		BankAccount:       etop.PbBankAccount(m.BankAccount),
	}
}

func PbMoneyTransactionShippingEtopExtendeds(items []*txmodely.MoneyTransactionShippingEtopExtended) []*MoneyTransactionShippingEtop {
	result := make([]*MoneyTransactionShippingEtop, len(items))
	for i, item := range items {
		result[i] = PbMoneyTransactionShippingEtopExtended(item)
	}
	return result
}
