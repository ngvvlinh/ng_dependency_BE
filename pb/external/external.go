package external

import (
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop-handler/webhook/sender"
	"etop.vn/backend/pkg/etop/model"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbgender "etop.vn/backend/pb/etop/etc/gender"
	pbshipping "etop.vn/backend/pb/etop/etc/shipping"
	pbsp "etop.vn/backend/pb/etop/etc/shipping_provider"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbs5 "etop.vn/backend/pb/etop/etc/status5"
	pbtryon "etop.vn/backend/pb/etop/etc/try_on"
	pborder "etop.vn/backend/pb/etop/order"
)

func PbPartner(m *model.Partner) *Partner {
	return &Partner{
		Id:              m.ID,
		Name:            m.Name,
		PublicName:      m.PublicName,
		Type:            pbetop.PbAccountType(model.TypePartner),
		Phone:           m.Phone,
		Website:         m.Website(),
		WebsiteUrl:      m.WebsiteURL,
		ImageUrl:        m.ImageURL,
		Email:           m.Email,
		RecognizedHosts: m.RecognizedHosts,
		RedirectUrls:    m.RedirectURLs,
	}
}

func (m *CreateWebhookRequest) ToModel(accountID int64) *model.Webhook {
	if m == nil {
		return nil
	}
	return &model.Webhook{
		AccountID: accountID,
		Entities:  m.Entities,
		Fields:    m.Fields,
		URL:       m.Url,
		Metadata:  m.Metadata,
	}
}

func PbWebhooks(items []*model.Webhook, states []sender.WebhookStates) []*Webhook {
	res := make([]*Webhook, len(items))
	for i, item := range items {
		res[i] = PbWebhook(item, states[i])
	}
	return res
}

func PbWebhook(m *model.Webhook, s sender.WebhookStates) *Webhook {
	if m == nil {
		return nil
	}
	return &Webhook{
		Id:        m.ID,
		Entities:  m.Entities,
		Fields:    m.Fields,
		Url:       m.URL,
		Metadata:  m.Metadata,
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		States: &WebhookStates{
			State:      string(s.State),
			LastSentAt: pbcm.PbTime(s.LastSentAt),
			LastError:  PbWebhookError(s.LastError),
		},
	}
}

func PbWebhookError(m *sender.WebhookStatesError) *WebhookError {
	if m == nil {
		return nil
	}
	return &WebhookError{
		Error:      m.ErrorMsg,
		RespStatus: int32(m.Status),
		RespBody:   m.Response,
		Retried:    int32(m.Retried),
		SentAt:     pbcm.PbTime(m.SentAt),
	}
}

func PbOrders(items []*ordermodel.Order) []*Order {
	res := make([]*Order, len(items))
	for i, item := range items {
		res[i] = PbOrder(item)
	}
	return res
}

func PbOrder(m *ordermodel.Order) *Order {
	res := &Order{
		Id:                        m.ID,
		ShopId:                    m.ShopID,
		Code:                      &m.Code,
		ExternalId:                &m.ExternalOrderID,
		ExternalCode:              &m.EdCode,
		ExternalUrl:               &m.ExternalURL,
		SelfUrl:                   cm.PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 pbcm.PbTime(m.CreatedAt),
		ProcessedAt:               pbcm.PbTime(m.ProcessedAt),
		UpdatedAt:                 pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                  pbcm.PbTime(m.ClosedAt),
		ConfirmedAt:               pbcm.PbTime(m.ConfirmedAt),
		CancelledAt:               pbcm.PbTime(m.CancelledAt),
		CancelReason:              &m.CancelReason,
		ConfirmStatus:             pbs3.PbPtrStatus(m.ConfirmStatus),
		Status:                    pbs5.PbPtrStatus(m.Status),
		FulfillmentShippingStatus: pbs5.PbPtrStatus(m.FulfillmentShippingStatus),
		EtopPaymentStatus:         pbs4.PbPtrStatus(m.EtopPaymentStatus),
		Lines:                     PbOrderLines(m.Lines),
		TotalItems:                pbcm.PbPtrInt32(m.TotalItems),
		BasketValue:               pbcm.PbPtrInt32(m.BasketValue),
		OrderDiscount:             pbcm.PbPtrInt32(m.OrderDiscount),
		TotalDiscount:             pbcm.PbPtrInt32(m.TotalDiscount),
		TotalFee:                  pbcm.PbPtrInt32(m.GetTotalFee()),
		FeeLines:                  pborder.PbOrderFeeLines(m.GetFeeLines()),
		TotalAmount:               pbcm.PbPtrInt32(m.TotalAmount),
		OrderNote:                 &m.OrderNote,
		Shipping:                  PbOrderShipping(m),
	}
	return res
}

func PbOrderShipping(m *ordermodel.Order) *OrderShipping {
	shipping := m.ShopShipping
	if shipping == nil {
		shipping = &ordermodel.OrderShipping{}
	}
	return &OrderShipping{
		PickupAddress:       PbOrderAddress(shipping.GetPickupAddress()),
		ReturnAddress:       PbOrderAddress(shipping.ReturnAddress),
		ShippingServiceName: &shipping.ExternalServiceName,
		ShippingServiceCode: shipping.GetPtrShippingServiceCode(),
		ShippingServiceFee:  pbcm.PbPtrInt32(shipping.ExternalShippingFee),
		Carrier:             pbsp.PbPtrShippingProvider(shipping.GetShippingProvider()),
		IncludeInsurance:    &shipping.IncludeInsurance,
		TryOn:               pbtryon.PbPtrTryOn(m.GetTryOn()),
		ShippingNote:        &m.ShippingNote,
		CodAmount:           pbcm.PbPtrInt32(m.ShopCOD),
		GrossWeight:         pbcm.PbPtrInt32(m.TotalWeight),
		Length:              pbcm.PbPtrInt32(shipping.Length),
		Width:               pbcm.PbPtrInt32(shipping.Width),
		Height:              pbcm.PbPtrInt32(shipping.Height),
		ChargeableWeight:    pbcm.PbPtrInt32(m.TotalWeight),
	}
}

func PbOrderHistories(items []ordermodel.OrderHistory) []*Order {
	res := make([]*Order, len(items))
	for i, item := range items {
		res[i] = PbOrderHistory(item)
	}
	return res
}

func PbOrderHistory(order ordermodel.OrderHistory) *Order {
	var customer *ordermodel.OrderCustomer
	_ = order.Customer().Unmarshal(&customer)
	var customerAddress, shippingAddress *ordermodel.OrderAddress
	_ = order.CustomerAddress().Unmarshal(&customerAddress)
	_ = order.ShippingAddress().Unmarshal(&shippingAddress)
	var lines []*ordermodel.OrderLine
	_ = order.Lines().Unmarshal(&lines)
	var shopShipping *ordermodel.OrderShipping
	_ = order.ShopShipping().Unmarshal(&shopShipping)
	var feeLines []ordermodel.OrderFeeLine
	_ = order.FeeLines().Unmarshal(&feeLines)

	res := &Order{
		Id:                        *order.ID().Int64(),
		ShopId:                    *order.ID().Int64(),
		Code:                      order.Code().String(),
		ExternalId:                order.ExternalOrderID().String(),
		ExternalCode:              order.EdCode().String(),
		ExternalUrl:               order.ExternalURL().String(),
		CustomerAddress:           PbOrderAddress(customerAddress),
		ShippingAddress:           PbOrderAddress(shippingAddress),
		CreatedAt:                 pbcm.PbTime(order.CreatedAt().Time()),
		ProcessedAt:               pbcm.PbTime(order.ProcessedAt().Time()),
		UpdatedAt:                 pbcm.PbTime(order.CreatedAt().Time()),
		ClosedAt:                  pbcm.PbTime(order.ClosedAt().Time()),
		ConfirmedAt:               pbcm.PbTime(order.ConfirmedAt().Time()),
		CancelledAt:               pbcm.PbTime(order.CancelledAt().Time()),
		CancelReason:              order.CancelReason().String(),
		ConfirmStatus:             pbs3.PbPtr(order.ConfirmStatus().Int()),
		Status:                    pbs5.PbPtr(order.Status().Int()),
		FulfillmentShippingStatus: pbs5.PbPtr(order.FulfillmentShippingStatus().Int()),
		EtopPaymentStatus:         pbs4.PbPtr(order.EtopPaymentStatus().Int()),
		Lines:                     PbOrderLines(lines),
		TotalItems:                order.TotalItems().Int32(),
		BasketValue:               order.BasketValue().Int32(),
		OrderDiscount:             order.OrderDiscount().Int32(),
		TotalDiscount:             order.TotalDiscount().Int32(),
		TotalFee:                  order.TotalFee().Int32(),
		FeeLines:                  nil,
		TotalAmount:               order.TotalAmount().Int32(),
		OrderNote:                 order.OrderNote().String(),
		Shipping:                  PbOrderShippingHistory(order, shopShipping),
	}
	if shopShipping != nil {
		res.ShippingAddress = PbOrderAddress(shopShipping.ShopAddress)
	}
	res.FeeLines = pborder.PbOrderFeeLines(ordermodel.GetFeeLinesWithFallback(feeLines, res.TotalDiscount, order.ShopShippingFee().Int32()))
	return res
}

func (m *Order) HasChanged() bool {
	return m.Status != nil ||
		m.ConfirmStatus != nil ||
		m.FulfillmentShippingStatus != nil ||
		m.EtopPaymentStatus != nil ||
		m.BasketValue != nil ||
		m.TotalAmount != nil ||
		m.Shipping != nil ||
		m.CustomerAddress != nil || m.ShippingAddress != nil
}

func PbOrderShippingHistory(order ordermodel.OrderHistory, shipping *ordermodel.OrderShipping) *OrderShipping {
	if shipping == nil {
		shipping = &ordermodel.OrderShipping{}
	}
	res := &OrderShipping{
		PickupAddress:       nil,
		ShippingServiceName: nil,
		ShippingServiceCode: nil,
		ShippingServiceFee:  nil,
		Carrier:             nil,
		IncludeInsurance:    nil,
		TryOn:               nil,
		ShippingNote:        nil,
		CodAmount:           nil,
		GrossWeight:         order.TotalWeight().Int32(),
		Length:              pbcm.PbPtrInt32(shipping.Length),
		Width:               pbcm.PbPtrInt32(shipping.Width),
		Height:              pbcm.PbPtrInt32(shipping.Height),
		ChargeableWeight:    order.TotalWeight().Int32(),
	}
	return res
}

func PbOrderLines(items []*ordermodel.OrderLine) []*OrderLine {
	// send changes as empty instead of "[]"
	if len(items) == 0 {
		return nil
	}
	res := make([]*OrderLine, len(items))
	for i, item := range items {
		res[i] = PbOrderLine(item)
	}
	return res
}

func PbOrderLine(m *ordermodel.OrderLine) *OrderLine {
	if m == nil {
		return nil
	}
	return &OrderLine{
		VariantId:    m.VariantID,
		ProductId:    m.ProductID,
		ProductName:  m.ProductName,
		Quantity:     int32(m.Quantity),
		ListPrice:    int32(m.ListPrice),
		RetailPrice:  int32(m.RetailPrice),
		PaymentPrice: pbcm.PbPtrInt32(m.PaymentPrice),
		ImageUrl:     m.ImageURL,
		Attributes:   pborder.PbAttributes(m.Attributes),
	}
}

func PbOrderAddress(m *ordermodel.OrderAddress) *OrderAddress {
	if m == nil {
		return nil
	}
	return &OrderAddress{
		FullName: m.GetFullName(),
		Phone:    m.Phone,
		Province: m.Province,
		District: m.District,
		Ward:     m.Ward,
		Company:  m.Company,
		Address1: m.Address1,
		Address2: m.Address2,
	}
}

func (m *OrderAddress) ToPbCustomer() *pborder.OrderCustomer {
	if m == nil {
		return nil
	}
	return &pborder.OrderCustomer{
		FirstName: "",
		LastName:  "",
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    0,
	}
}

func (m *OrderAddress) ToPbOrder() *pborder.OrderAddress {
	if m == nil {
		return nil
	}
	return &pborder.OrderAddress{
		FullName: m.FullName,
		Phone:    m.Phone,
		Province: m.Province,
		District: m.District,
		Ward:     m.Ward,
		Address1: m.Address1,
		Address2: m.Address2,
	}
}

func PbOrderAddressFromAddress(m *model.Address) *OrderAddress {
	if m == nil {
		return nil
	}
	return &OrderAddress{
		FullName: m.GetFullName(),
		Phone:    m.Phone,
		Province: m.Province,
		District: m.District,
		Ward:     m.Ward,
		Company:  m.Company,
		Address1: m.Address1,
		Address2: m.Address2,
	}
}

func PbOrderCustomer(m *ordermodel.OrderCustomer) *OrderCustomer {
	if m == nil {
		return nil
	}
	return &OrderCustomer{
		FullName: m.GetFullName(),
		Email:    m.Email,
		Phone:    m.Phone,
		Gender:   pbgender.PbGender(m.Gender),
	}
}

func (m *OrderCustomer) ToPbOrder() *pborder.OrderCustomer {
	if m == nil {
		return nil
	}
	return &pborder.OrderCustomer{
		FirstName: "",
		LastName:  "",
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    m.Gender,
	}
}

func PbFulfillments(items []*shipmodel.Fulfillment) []*Fulfillment {
	res := make([]*Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item)
	}
	return res
}

func PbFulfillment(m *shipmodel.Fulfillment) *Fulfillment {
	return &Fulfillment{
		Id:                       m.ID,
		OrderId:                  m.OrderID,
		ShopId:                   m.ShopID,
		SelfUrl:                  cm.PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		TotalItems:               pbcm.PbPtrInt32(m.TotalItems),
		BasketValue:              pbcm.PbPtrInt32(m.BasketValue),
		CreatedAt:                pbcm.PbTime(m.CreatedAt),
		UpdatedAt:                pbcm.PbTime(m.UpdatedAt),
		ClosedAt:                 pbcm.PbTime(m.ClosedAt),
		CancelledAt:              pbcm.PbTime(m.ShippingCancelledAt),
		CancelReason:             &m.CancelReason,
		Carrier:                  pbsp.PbPtrShippingProvider(m.ShippingProvider),
		ShippingServiceName:      &m.ExternalShippingName,
		ShippingServiceFee:       pbcm.PbPtrInt32(m.ShippingServiceFee),
		ActualShippingServiceFee: pbcm.PbPtrInt32(m.ShippingFeeShop),
		ShippingServiceCode:      &m.ProviderServiceID,
		ShippingCode:             &m.ShippingCode,
		ShippingNote:             &m.ShippingNote,
		TryOn:                    pbtryon.PbPtrTryOn(m.TryOn),
		IncludeInsurance:         &m.IncludeInsurance,
		ConfirmStatus:            pbs3.PbPtrStatus(m.ConfirmStatus),
		ShippingState:            pbshipping.PbPtrShippingState(m.ShippingState),
		ShippingStatus:           pbs5.PbPtrStatus(m.ShippingStatus),
		Status:                   pbs5.PbPtrStatus(m.Status),
		CodAmount:                pbcm.PbPtrInt32(m.OriginalCODAmount),
		ActualCodAmount:          pbcm.PbPtrInt32(m.TotalCODAmount),
		ChargeableWeight:         pbcm.PbPtrInt32(m.TotalWeight),
		PickupAddress:            PbOrderAddressFromAddress(m.AddressFrom),
		ReturnAddress:            PbOrderAddressFromAddress(m.AddressReturn),
		ShippingAddress:          PbOrderAddressFromAddress(m.AddressTo),
		EtopPaymentStatus:        pbs4.PbPtrStatus(m.EtopPaymentStatus),
		EstimatedDeliveryAt:      pbcm.PbTime(m.ExpectedDeliveryAt),
		EstimatedPickupAt:        pbcm.PbTime(m.ExpectedPickAt),
	}
}

func PbFulfillmentHistories(items []shipmodel.FulfillmentHistory) []*Fulfillment {
	res := make([]*Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillmentHistory(item)
	}
	return res
}

func PbFulfillmentHistory(m shipmodel.FulfillmentHistory) *Fulfillment {
	var addressTo, addressFrom, addressReturn *ordermodel.OrderAddress
	_ = m.AddressTo().Unmarshal(&addressTo)
	_ = m.AddressFrom().Unmarshal(&addressFrom)
	_ = m.AddressReturn().Unmarshal(&addressReturn)

	return &Fulfillment{
		Id:                       *m.ID().Int64(),
		OrderId:                  *m.OrderID().Int64(),
		ShopId:                   *m.ShopID().Int64(),
		TotalItems:               m.TotalItems().Int32(),
		BasketValue:              m.BasketValue().Int32(),
		CreatedAt:                pbcm.PbTime(m.CreatedAt().Time()),
		UpdatedAt:                pbcm.PbTime(m.UpdatedAt().Time()),
		ClosedAt:                 pbcm.PbTime(m.ClosedAt().Time()),
		CancelledAt:              pbcm.PbTime(m.CancelReason().Time()),
		CancelReason:             m.CancelReason().String(),
		Carrier:                  pbsp.PbShippingProviderPtr(m.ShippingProvider().String()),
		ShippingServiceName:      m.ExternalShippingName().String(),
		ShippingServiceFee:       m.ShippingServiceFee().Int32(),
		ActualShippingServiceFee: m.ShippingFeeShop().Int32(),
		ShippingServiceCode:      m.ProviderServiceID().String(),
		ShippingCode:             m.ShippingCode().String(),
		ShippingNote:             m.ShippingNote().String(),
		TryOn:                    pbtryon.PbTryOnPtr(m.TryOn().String()),
		IncludeInsurance:         m.IncludeInsurance().Bool(),
		ConfirmStatus:            pbs3.PbPtr(m.ConfirmStatus().Int()),
		ShippingState:            pbshipping.PbPtr(m.ShippingState().String()),
		ShippingStatus:           pbs5.PbPtr(m.ShippingStatus().Int()),
		Status:                   pbs5.PbPtr(m.Status().Int()),
		CodAmount:                m.OriginalCODAmount().Int32(),
		ActualCodAmount:          m.TotalCODAmount().Int32(),
		ChargeableWeight:         m.TotalWeight().Int32(),
		PickupAddress:            PbOrderAddress(addressFrom),
		ReturnAddress:            PbOrderAddress(addressReturn),
		ShippingAddress:          PbOrderAddress(addressTo),
		EtopPaymentStatus:        pbs4.PbPtr(m.EtopPaymentStatus().Int()),
		EstimatedDeliveryAt:      pbcm.PbTime(m.ExpectedDeliveryAt().Time()),
		EstimatedPickupAt:        pbcm.PbTime(m.ExpectedPickAt().Time()),
	}
}

func (m *Fulfillment) HasChanged() bool {
	return m.Status != nil ||
		m.ShippingState != nil ||
		m.EtopPaymentStatus != nil ||
		m.ActualShippingServiceFee != nil ||
		m.CodAmount != nil ||
		m.ActualCodAmount != nil ||
		m.ShippingNote != nil ||
		m.ChargeableWeight != nil
}

func PbShippingServices(items []*model.AvailableShippingService) []*ShippingService {
	res := make([]*ShippingService, len(items))
	for i, item := range items {
		res[i] = PbShippingService(item)
	}
	return res
}

func PbShippingService(m *model.AvailableShippingService) *ShippingService {
	return &ShippingService{
		Code:                m.ProviderServiceID,
		Name:                m.Name,
		Fee:                 int32(m.ServiceFee),
		Carrier:             pbsp.PbShippingProviderType(m.Provider),
		EstimatedPickupAt:   pbcm.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt: pbcm.PbTime(m.ExpectedDeliveryAt),
	}
}

func OrderLinesToCreateOrderLines(items []*OrderLine) (_ []*pborder.CreateOrderLine, err error) {
	res := make([]*pborder.CreateOrderLine, len(items))
	for i, item := range items {
		res[i], err = OrderLineToCreateOrderLine(item)
		if err != nil {
			return
		}
	}
	return res, nil
}

func OrderLineToCreateOrderLine(m *OrderLine) (*pborder.CreateOrderLine, error) {
	if m == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "order_line must not be null")
	}
	if m.PaymentPrice == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp payment_price")
	}

	return &pborder.CreateOrderLine{
		VariantId:    m.VariantId,
		ProductName:  m.ProductName,
		Quantity:     m.Quantity,
		ListPrice:    m.ListPrice,
		RetailPrice:  m.RetailPrice,
		PaymentPrice: *m.PaymentPrice,
		ImageUrl:     m.ImageUrl,
		Attributes:   m.Attributes,
	}, nil
}

func PbOrderAndFulfillments(order *ordermodel.Order, fulfillments []*shipmodel.Fulfillment) *OrderAndFulfillments {
	return &OrderAndFulfillments{
		Order:        PbOrder(order),
		Fulfillments: PbFulfillments(fulfillments),
	}
}
