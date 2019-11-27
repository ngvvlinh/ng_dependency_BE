package convertpb

import (
	"etop.vn/api/pb/etop/etc/gender"
	"etop.vn/api/pb/etop/order"
	"etop.vn/api/pb/external"
	"etop.vn/backend/com/handler/etop-handler/webhook/sender"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func PbPartner(m *model.Partner) *external.Partner {
	return &external.Partner{
		Id:              m.ID,
		Name:            m.Name,
		PublicName:      m.PublicName,
		Type:            convertpb.PbAccountType(model.TypePartner),
		Phone:           m.Phone,
		Website:         m.Website(),
		WebsiteUrl:      m.WebsiteURL,
		ImageUrl:        m.ImageURL,
		Email:           m.Email,
		RecognizedHosts: m.RecognizedHosts,
		RedirectUrls:    m.RedirectURLs,
	}
}

func CreateWebhookRequestToModel(m *external.CreateWebhookRequest, accountID dot.ID) *model.Webhook {
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

func PbWebhooks(items []*model.Webhook, states []sender.WebhookStates) []*external.Webhook {
	res := make([]*external.Webhook, len(items))
	for i, item := range items {
		res[i] = PbWebhook(item, states[i])
	}
	return res
}

func PbWebhook(m *model.Webhook, s sender.WebhookStates) *external.Webhook {
	if m == nil {
		return nil
	}
	return &external.Webhook{
		Id:        m.ID,
		Entities:  m.Entities,
		Fields:    m.Fields,
		Url:       m.URL,
		Metadata:  m.Metadata,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		States: &external.WebhookStates{
			State:      string(s.State),
			LastSentAt: cmapi.PbTime(s.LastSentAt),
			LastError:  PbWebhookError(s.LastError),
		},
	}
}

func PbWebhookError(m *sender.WebhookStatesError) *external.WebhookError {
	if m == nil {
		return nil
	}
	return &external.WebhookError{
		Error:      m.ErrorMsg,
		RespStatus: int(m.Status),
		RespBody:   m.Response,
		Retried:    int(m.Retried),
		SentAt:     cmapi.PbTime(m.SentAt),
	}
}

func PbOrders(items []*ordermodel.Order) []*external.Order {
	res := make([]*external.Order, len(items))
	for i, item := range items {
		res[i] = PbOrder(item)
	}
	return res
}

func PbOrder(m *ordermodel.Order) *external.Order {
	res := &external.Order{
		Id:                        m.ID,
		ShopId:                    m.ShopID,
		Code:                      &m.Code,
		ExternalId:                &m.ExternalOrderID,
		ExternalCode:              &m.EdCode,
		ExternalUrl:               &m.ExternalURL,
		SelfUrl:                   cm.PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 cmapi.PbTime(m.CreatedAt),
		ProcessedAt:               cmapi.PbTime(m.ProcessedAt),
		UpdatedAt:                 cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                  cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:               cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:               cmapi.PbTime(m.CancelledAt),
		CancelReason:              &m.CancelReason,
		ConfirmStatus:             convertpb.Pb3PtrStatus(m.ConfirmStatus),
		Status:                    convertpb.Pb5PtrStatus(m.Status),
		FulfillmentShippingStatus: convertpb.Pb5PtrStatus(m.FulfillmentShippingStatus),
		EtopPaymentStatus:         convertpb.Pb4PtrStatus(m.EtopPaymentStatus),
		Lines:                     PbOrderLines(m.Lines),
		TotalItems:                cmapi.PbPtrInt(m.TotalItems),
		BasketValue:               cmapi.PbPtrInt(m.BasketValue),
		OrderDiscount:             cmapi.PbPtrInt(m.OrderDiscount),
		TotalDiscount:             cmapi.PbPtrInt(m.TotalDiscount),
		TotalFee:                  cmapi.PbPtrInt(m.GetTotalFee()),
		FeeLines:                  convertpb.PbOrderFeeLines(m.GetFeeLines()),
		TotalAmount:               cmapi.PbPtrInt(m.TotalAmount),
		OrderNote:                 &m.OrderNote,
		Shipping:                  PbOrderShipping(m),
	}
	return res
}

func PbOrderShipping(m *ordermodel.Order) *external.OrderShipping {
	shipping := m.ShopShipping
	if shipping == nil {
		shipping = &ordermodel.OrderShipping{}
	}
	return &external.OrderShipping{
		PickupAddress:       PbOrderAddress(shipping.GetPickupAddress()),
		ReturnAddress:       PbOrderAddress(shipping.ReturnAddress),
		ShippingServiceName: &shipping.ExternalServiceName,
		ShippingServiceCode: shipping.GetPtrShippingServiceCode(),
		ShippingServiceFee:  cmapi.PbPtrInt(shipping.ExternalShippingFee),
		Carrier:             convertpb.PbPtrShippingProvider(shipping.GetShippingProvider()),
		IncludeInsurance:    &shipping.IncludeInsurance,
		TryOn:               convertpb.PbPtrTryOn(m.GetTryOn()),
		ShippingNote:        &m.ShippingNote,
		CodAmount:           cmapi.PbPtrInt(m.ShopCOD),
		GrossWeight:         cmapi.PbPtrInt(m.TotalWeight),
		Length:              cmapi.PbPtrInt(shipping.Length),
		Width:               cmapi.PbPtrInt(shipping.Width),
		Height:              cmapi.PbPtrInt(shipping.Height),
		ChargeableWeight:    cmapi.PbPtrInt(m.TotalWeight),
	}
}

func PbOrderHistories(items []ordermodel.OrderHistory) []*external.Order {
	res := make([]*external.Order, len(items))
	for i, item := range items {
		res[i] = PbOrderHistory(item)
	}
	return res
}

func PbOrderHistory(order ordermodel.OrderHistory) *external.Order {
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

	res := &external.Order{
		Id:                        *order.ID().ID(),
		ShopId:                    *order.ID().ID(),
		Code:                      order.Code().String(),
		ExternalId:                order.ExternalOrderID().String(),
		ExternalCode:              order.EdCode().String(),
		ExternalUrl:               order.ExternalURL().String(),
		CustomerAddress:           PbOrderAddress(customerAddress),
		ShippingAddress:           PbOrderAddress(shippingAddress),
		CreatedAt:                 cmapi.PbTime(order.CreatedAt().Time()),
		ProcessedAt:               cmapi.PbTime(order.ProcessedAt().Time()),
		UpdatedAt:                 cmapi.PbTime(order.CreatedAt().Time()),
		ClosedAt:                  cmapi.PbTime(order.ClosedAt().Time()),
		ConfirmedAt:               cmapi.PbTime(order.ConfirmedAt().Time()),
		CancelledAt:               cmapi.PbTime(order.CancelledAt().Time()),
		CancelReason:              order.CancelReason().String(),
		ConfirmStatus:             convertpb.Pb3Ptr(order.ConfirmStatus().Int()),
		Status:                    convertpb.Pb5Ptr(order.Status().Int()),
		FulfillmentShippingStatus: convertpb.Pb5Ptr(order.FulfillmentShippingStatus().Int()),
		EtopPaymentStatus:         convertpb.Pb4Ptr(order.EtopPaymentStatus().Int()),
		Lines:                     PbOrderLines(lines),
		TotalItems:                order.TotalItems().Int(),
		BasketValue:               order.BasketValue().Int(),
		OrderDiscount:             order.OrderDiscount().Int(),
		TotalDiscount:             order.TotalDiscount().Int(),
		TotalFee:                  order.TotalFee().Int(),
		FeeLines:                  nil,
		TotalAmount:               order.TotalAmount().Int(),
		OrderNote:                 order.OrderNote().String(),
		Shipping:                  PbOrderShippingHistory(order, shopShipping),
	}
	if shopShipping != nil {
		res.ShippingAddress = PbOrderAddress(shopShipping.ShopAddress)
	}
	res.FeeLines = convertpb.PbOrderFeeLines(ordermodel.GetFeeLinesWithFallback(feeLines, res.TotalDiscount, order.ShopShippingFee().Int()))
	return res
}

func PbOrderShippingHistory(order ordermodel.OrderHistory, shipping *ordermodel.OrderShipping) *external.OrderShipping {
	if shipping == nil {
		shipping = &ordermodel.OrderShipping{}
	}
	res := &external.OrderShipping{
		PickupAddress:       nil,
		ShippingServiceName: nil,
		ShippingServiceCode: nil,
		ShippingServiceFee:  nil,
		Carrier:             nil,
		IncludeInsurance:    nil,
		TryOn:               nil,
		ShippingNote:        nil,
		CodAmount:           nil,
		GrossWeight:         order.TotalWeight().Int(),
		Length:              cmapi.PbPtrInt(shipping.Length),
		Width:               cmapi.PbPtrInt(shipping.Width),
		Height:              cmapi.PbPtrInt(shipping.Height),
		ChargeableWeight:    order.TotalWeight().Int(),
	}
	return res
}

func PbOrderLines(items []*ordermodel.OrderLine) []*external.OrderLine {
	// send changes as empty instead of "[]"
	if len(items) == 0 {
		return nil
	}
	res := make([]*external.OrderLine, len(items))
	for i, item := range items {
		res[i] = PbOrderLine(item)
	}
	return res
}

func PbOrderLine(m *ordermodel.OrderLine) *external.OrderLine {
	if m == nil {
		return nil
	}
	return &external.OrderLine{
		VariantId:    m.VariantID,
		ProductId:    m.ProductID,
		ProductName:  m.ProductName,
		Quantity:     int(m.Quantity),
		ListPrice:    int(m.ListPrice),
		RetailPrice:  int(m.RetailPrice),
		PaymentPrice: cmapi.PbPtrInt(m.PaymentPrice),
		ImageUrl:     m.ImageURL,
		Attributes:   convertpb.PbAttributesFromModel(m.Attributes),
	}
}

func PbOrderAddress(m *ordermodel.OrderAddress) *external.OrderAddress {
	if m == nil {
		return nil
	}
	return &external.OrderAddress{
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

func OrderAddressToPbCustomer(m *external.OrderAddress) *order.OrderCustomer {
	if m == nil {
		return nil
	}
	return &order.OrderCustomer{
		FirstName: "",
		LastName:  "",
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    0,
	}
}

func OrderAddressToPbOrder(m *external.OrderAddress) *order.OrderAddress {
	if m == nil {
		return nil
	}
	return &order.OrderAddress{
		FullName: m.FullName,
		Phone:    m.Phone,
		Province: m.Province,
		District: m.District,
		Ward:     m.Ward,
		Address1: m.Address1,
		Address2: m.Address2,
	}
}

func PbOrderAddressFromAddress(m *model.Address) *external.OrderAddress {
	if m == nil {
		return nil
	}
	return &external.OrderAddress{
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

func PbOrderCustomer(m *ordermodel.OrderCustomer) *external.OrderCustomer {
	if m == nil {
		return nil
	}
	return &external.OrderCustomer{
		FullName: m.GetFullName(),
		Email:    m.Email,
		Phone:    m.Phone,
		Gender:   gender.PbGender(m.Gender),
	}
}

func OrderCustomerToPbOrder(m *external.OrderCustomer) *order.OrderCustomer {
	if m == nil {
		return nil
	}
	return &order.OrderCustomer{
		FirstName: "",
		LastName:  "",
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    m.Gender,
	}
}

func PbFulfillments(items []*shipmodel.Fulfillment) []*external.Fulfillment {
	res := make([]*external.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item)
	}
	return res
}

func PbFulfillment(m *shipmodel.Fulfillment) *external.Fulfillment {
	return &external.Fulfillment{
		Id:                       m.ID,
		OrderId:                  m.OrderID,
		ShopId:                   m.ShopID,
		SelfUrl:                  cm.PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		TotalItems:               cmapi.PbPtrInt(m.TotalItems),
		BasketValue:              cmapi.PbPtrInt(m.BasketValue),
		CreatedAt:                cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                 cmapi.PbTime(m.ClosedAt),
		CancelledAt:              cmapi.PbTime(m.ShippingCancelledAt),
		CancelReason:             &m.CancelReason,
		Carrier:                  convertpb.PbPtrShippingProvider(m.ShippingProvider),
		ShippingServiceName:      &m.ExternalShippingName,
		ShippingServiceFee:       cmapi.PbPtrInt(m.ShippingServiceFee),
		ActualShippingServiceFee: cmapi.PbPtrInt(m.ShippingFeeShop),
		ShippingServiceCode:      &m.ProviderServiceID,
		ShippingCode:             &m.ShippingCode,
		ShippingNote:             &m.ShippingNote,
		TryOn:                    convertpb.PbPtrTryOn(m.TryOn),
		IncludeInsurance:         &m.IncludeInsurance,
		ConfirmStatus:            convertpb.Pb3PtrStatus(m.ConfirmStatus),
		ShippingState:            convertpb.PbPtrShippingState(m.ShippingState),
		ShippingStatus:           convertpb.Pb5PtrStatus(m.ShippingStatus),
		Status:                   convertpb.Pb5PtrStatus(m.Status),
		CodAmount:                cmapi.PbPtrInt(m.OriginalCODAmount),
		ActualCodAmount:          cmapi.PbPtrInt(m.TotalCODAmount),
		ChargeableWeight:         cmapi.PbPtrInt(m.TotalWeight),
		PickupAddress:            PbOrderAddressFromAddress(m.AddressFrom),
		ReturnAddress:            PbOrderAddressFromAddress(m.AddressReturn),
		ShippingAddress:          PbOrderAddressFromAddress(m.AddressTo),
		EtopPaymentStatus:        convertpb.Pb4PtrStatus(m.EtopPaymentStatus),
		EstimatedDeliveryAt:      cmapi.PbTime(m.ExpectedDeliveryAt),
		EstimatedPickupAt:        cmapi.PbTime(m.ExpectedPickAt),
	}
}

func PbFulfillmentHistories(items []shipmodel.FulfillmentHistory) []*external.Fulfillment {
	res := make([]*external.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillmentHistory(item)
	}
	return res
}

func PbFulfillmentHistory(m shipmodel.FulfillmentHistory) *external.Fulfillment {
	var addressTo, addressFrom, addressReturn *ordermodel.OrderAddress
	_ = m.AddressTo().Unmarshal(&addressTo)
	_ = m.AddressFrom().Unmarshal(&addressFrom)
	_ = m.AddressReturn().Unmarshal(&addressReturn)

	return &external.Fulfillment{
		Id:                       *m.ID().ID(),
		OrderId:                  *m.OrderID().ID(),
		ShopId:                   *m.ShopID().ID(),
		TotalItems:               m.TotalItems().Int(),
		BasketValue:              m.BasketValue().Int(),
		CreatedAt:                cmapi.PbTime(m.CreatedAt().Time()),
		UpdatedAt:                cmapi.PbTime(m.UpdatedAt().Time()),
		ClosedAt:                 cmapi.PbTime(m.ClosedAt().Time()),
		CancelledAt:              cmapi.PbTime(m.CancelReason().Time()),
		CancelReason:             m.CancelReason().String(),
		Carrier:                  convertpb.PbShippingProviderPtr(m.ShippingProvider().String()),
		ShippingServiceName:      m.ExternalShippingName().String(),
		ShippingServiceFee:       m.ShippingServiceFee().Int(),
		ActualShippingServiceFee: m.ShippingFeeShop().Int(),
		ShippingServiceCode:      m.ProviderServiceID().String(),
		ShippingCode:             m.ShippingCode().String(),
		ShippingNote:             m.ShippingNote().String(),
		TryOn:                    convertpb.PbTryOnPtr(m.TryOn().String()),
		IncludeInsurance:         m.IncludeInsurance().Bool(),
		ConfirmStatus:            convertpb.Pb3Ptr(m.ConfirmStatus().Int()),
		ShippingState:            convertpb.PbPtr(m.ShippingState().String()),
		ShippingStatus:           convertpb.Pb5Ptr(m.ShippingStatus().Int()),
		Status:                   convertpb.Pb5Ptr(m.Status().Int()),
		CodAmount:                m.OriginalCODAmount().Int(),
		ActualCodAmount:          m.TotalCODAmount().Int(),
		ChargeableWeight:         m.TotalWeight().Int(),
		PickupAddress:            PbOrderAddress(addressFrom),
		ReturnAddress:            PbOrderAddress(addressReturn),
		ShippingAddress:          PbOrderAddress(addressTo),
		EtopPaymentStatus:        convertpb.Pb4Ptr(m.EtopPaymentStatus().Int()),
		EstimatedDeliveryAt:      cmapi.PbTime(m.ExpectedDeliveryAt().Time()),
		EstimatedPickupAt:        cmapi.PbTime(m.ExpectedPickAt().Time()),
	}
}

func PbShippingServices(items []*model.AvailableShippingService) []*external.ShippingService {
	res := make([]*external.ShippingService, len(items))
	for i, item := range items {
		res[i] = PbShippingService(item)
	}
	return res
}

func PbShippingService(m *model.AvailableShippingService) *external.ShippingService {
	return &external.ShippingService{
		Code:                m.ProviderServiceID,
		Name:                m.Name,
		Fee:                 int(m.ServiceFee),
		Carrier:             convertpb.PbShippingProviderType(m.Provider),
		EstimatedPickupAt:   cmapi.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt: cmapi.PbTime(m.ExpectedDeliveryAt),
	}
}

func OrderLinesToCreateOrderLines(items []*external.OrderLine) (_ []*order.CreateOrderLine, err error) {
	res := make([]*order.CreateOrderLine, len(items))
	for i, item := range items {
		res[i], err = OrderLineToCreateOrderLine(item)
		if err != nil {
			return
		}
	}
	return res, nil
}

func OrderLineToCreateOrderLine(m *external.OrderLine) (*order.CreateOrderLine, error) {
	if m == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "order_line must not be null")
	}
	if m.PaymentPrice == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp payment_price")
	}

	return &order.CreateOrderLine{
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

func PbOrderAndFulfillments(order *ordermodel.Order, fulfillments []*shipmodel.Fulfillment) *external.OrderAndFulfillments {
	return &external.OrderAndFulfillments{
		Order:        PbOrder(order),
		Fulfillments: PbFulfillments(fulfillments),
	}
}
