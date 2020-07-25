package convertpb

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/meta"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	exttypes "o.o/api/top/external/types"
	"o.o/api/top/int/etop"
	"o.o/api/top/int/types"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/backend/com/eventhandler/webhook/sender"
	addressmodel "o.o/backend/com/main/address/model"
	"o.o/backend/com/main/catalog/convert"
	catalogmodel "o.o/backend/com/main/catalog/model"
	identitymodel "o.o/backend/com/main/identity/model"
	inventorymodel "o.o/backend/com/main/inventory/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	customermodel "o.o/backend/com/shopping/customering/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

func PNonZeroString(s string) dot.NullString {
	if s == "" {
		return dot.NullString{}
	}
	return dot.String(s)
}

func PbPartner(m *identitymodel.Partner) *exttypes.Partner {
	return &exttypes.Partner{
		Id:              m.ID,
		Name:            m.Name,
		PublicName:      m.PublicName,
		Type:            convertpb.PbAccountType(account_type.Partner),
		Phone:           m.Phone,
		Website:         m.Website(),
		WebsiteUrl:      m.WebsiteURL,
		ImageUrl:        m.ImageURL,
		Email:           m.Email,
		RecognizedHosts: m.RecognizedHosts,
		RedirectUrls:    m.RedirectURLs,
	}
}

func CreateWebhookRequestToModel(m *exttypes.CreateWebhookRequest, accountID dot.ID) *model.Webhook {
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

func PbWebhooks(items []*model.Webhook, states []sender.WebhookStates) []*exttypes.Webhook {
	res := make([]*exttypes.Webhook, len(items))
	for i, item := range items {
		res[i] = PbWebhook(item, states[i])
	}
	return res
}

func PbWebhook(m *model.Webhook, s sender.WebhookStates) *exttypes.Webhook {
	if m == nil {
		return nil
	}
	return &exttypes.Webhook{
		Id:        m.ID,
		Entities:  m.Entities,
		Fields:    m.Fields,
		Url:       m.URL,
		Metadata:  m.Metadata,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		States: &exttypes.WebhookStates{
			State:      s.State.String(),
			LastSentAt: cmapi.PbTime(s.LastSentAt),
			LastError:  PbWebhookError(s.LastError),
		},
	}
}

func PbWebhookError(m *sender.WebhookStatesError) *exttypes.WebhookError {
	if m == nil {
		return nil
	}
	return &exttypes.WebhookError{
		Error:      m.ErrorMsg,
		RespStatus: m.Status,
		RespBody:   m.Response,
		Retried:    m.Retried,
		SentAt:     cmapi.PbTime(m.SentAt),
	}
}

func PbOrders(items []*ordermodel.Order) []*exttypes.Order {
	res := make([]*exttypes.Order, len(items))
	for i, item := range items {
		res[i] = PbOrder(item)
	}
	return res
}

func PbOrder(m *ordermodel.Order) *exttypes.Order {
	if m == nil {
		return nil
	}
	res := &exttypes.Order{
		Id:                        m.ID,
		ShopId:                    m.ShopID,
		Code:                      dot.String(m.Code),
		ExternalId:                dot.String(m.ExternalOrderID),
		ExternalCode:              dot.String(m.EdCode),
		ExternalUrl:               dot.String(m.ExternalURL),
		SelfUrl:                   PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		CustomerAddress:           PbOrderAddress(m.CustomerAddress),
		ShippingAddress:           PbOrderAddress(m.ShippingAddress),
		CreatedAt:                 cmapi.PbTime(m.CreatedAt),
		ProcessedAt:               cmapi.PbTime(m.ProcessedAt),
		UpdatedAt:                 cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                  cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:               cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:               cmapi.PbTime(m.CancelledAt),
		CancelReason:              dot.String(m.CancelReason),
		ConfirmStatus:             m.ConfirmStatus.Wrap(),
		Status:                    m.Status.Wrap(),
		FulfillmentShippingStatus: m.FulfillmentShippingStatus.Wrap(),
		EtopPaymentStatus:         m.EtopPaymentStatus.Wrap(),
		Lines:                     PbOrderLines(m.Lines),
		TotalItems:                cmapi.PbPtrInt(m.TotalItems),
		BasketValue:               cmapi.PbPtrInt(m.BasketValue),
		OrderDiscount:             cmapi.PbPtrInt(m.OrderDiscount),
		TotalDiscount:             cmapi.PbPtrInt(m.TotalDiscount),
		TotalFee:                  cmapi.PbPtrInt(m.GetTotalFee()),
		FeeLines:                  convertpb.PbOrderFeeLines(m.GetFeeLines()),
		TotalAmount:               cmapi.PbPtrInt(m.TotalAmount),
		OrderNote:                 dot.String(m.OrderNote),
		PreOrder:                  m.PreOrder,
		Shipping:                  PbOrderShipping(m),
	}
	return res
}

func PbOrderShipping(m *ordermodel.Order) *exttypes.OrderShipping {
	shipping := m.ShopShipping
	if shipping == nil {
		shipping = &ordermodel.OrderShipping{}
	}
	return &exttypes.OrderShipping{
		PickupAddress:       PbOrderAddress(shipping.GetPickupAddress()),
		ReturnAddress:       PbOrderAddress(shipping.ReturnAddress),
		ShippingServiceName: dot.String(shipping.ExternalServiceName),
		ShippingServiceCode: shipping.GetPtrShippingServiceCode(),
		ShippingServiceFee:  cmapi.PbPtrInt(shipping.ExternalShippingFee),
		Carrier:             shipping.GetShippingProvider(),
		IncludeInsurance:    dot.Bool(shipping.IncludeInsurance),
		TryOn:               m.GetTryOn(),
		ShippingNote:        dot.String(m.ShippingNote),
		CodAmount:           cmapi.PbPtrInt(m.ShopCOD),
		GrossWeight:         cmapi.PbPtrInt(m.TotalWeight),
		Length:              cmapi.PbPtrInt(shipping.Length),
		Width:               cmapi.PbPtrInt(shipping.Width),
		Height:              cmapi.PbPtrInt(shipping.Height),
		ChargeableWeight:    cmapi.PbPtrInt(m.TotalWeight),
	}
}

func PbOrderWithoutShipping(m *ordermodel.Order) *exttypes.OrderWithoutShipping {
	res := &exttypes.OrderWithoutShipping{
		Id:              m.ID,
		ShopId:          m.ShopID,
		Code:            dot.String(m.Code),
		ExternalId:      dot.String(m.ExternalOrderID),
		ExternalCode:    dot.String(m.EdCode),
		ExternalUrl:     dot.String(m.ExternalURL),
		SelfUrl:         PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		CustomerAddress: PbOrderAddress(m.CustomerAddress),
		ShippingAddress: PbOrderAddress(m.ShippingAddress),
		CreatedAt:       cmapi.PbTime(m.CreatedAt),
		ProcessedAt:     cmapi.PbTime(m.ProcessedAt),
		UpdatedAt:       cmapi.PbTime(m.UpdatedAt),
		ClosedAt:        cmapi.PbTime(m.ClosedAt),
		ConfirmedAt:     cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:     cmapi.PbTime(m.CancelledAt),
		CancelReason:    dot.String(m.CancelReason),
		ConfirmStatus:   m.ConfirmStatus.Wrap(),
		Status:          m.Status.Wrap(),
		Lines:           PbOrderLines(m.Lines),
		TotalItems:      cmapi.PbPtrInt(m.TotalItems),
		BasketValue:     cmapi.PbPtrInt(m.BasketValue),
		OrderDiscount:   cmapi.PbPtrInt(m.OrderDiscount),
		TotalDiscount:   cmapi.PbPtrInt(m.TotalDiscount),
		TotalFee:        cmapi.PbPtrInt(m.GetTotalFee()),
		FeeLines:        convertpb.PbOrderFeeLines(m.GetFeeLines()),
		TotalAmount:     cmapi.PbPtrInt(m.TotalAmount),
		OrderNote:       dot.String(m.OrderNote),
	}
	return res
}

func PbOrderToOrderWithoutShipping(m *exttypes.OrderWithoutShipping) *exttypes.Order {
	if m == nil {
		return nil
	}
	res := &exttypes.Order{
		Id:              m.Id,
		ShopId:          m.ShopId,
		Code:            m.Code,
		ExternalId:      m.ExternalId,
		ExternalCode:    m.ExternalCode,
		ExternalUrl:     m.ExternalUrl,
		SelfUrl:         m.SelfUrl,
		CustomerAddress: m.CustomerAddress,
		ShippingAddress: m.ShippingAddress,
		CreatedAt:       m.CreatedAt,
		ProcessedAt:     m.ProcessedAt,
		UpdatedAt:       m.UpdatedAt,
		ClosedAt:        m.ClosedAt,
		ConfirmedAt:     m.ConfirmedAt,
		CancelledAt:     m.CancelledAt,
		CancelReason:    m.CancelReason,
		ConfirmStatus:   m.ConfirmStatus,
		Status:          m.Status,
		Lines:           m.Lines,
		TotalItems:      m.TotalItems,
		BasketValue:     m.BasketValue,
		OrderDiscount:   m.OrderDiscount,
		TotalDiscount:   m.TotalDiscount,
		TotalFee:        m.TotalFee,
		FeeLines:        m.FeeLines,
		TotalAmount:     m.TotalAmount,
		OrderNote:       m.OrderNote,
	}
	return res
}

func PbOrderHistories(items []ordermodel.OrderHistory) []*exttypes.Order {
	res := make([]*exttypes.Order, len(items))
	for i, item := range items {
		res[i] = PbOrderHistory(item)
	}
	return res
}

func PbOrderHistory(order ordermodel.OrderHistory) *exttypes.Order {
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

	res := &exttypes.Order{
		Id:                        order.ID().ID().Apply(0),
		ShopId:                    order.ID().ID().Apply(0),
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

func PbOrderShippingHistory(order ordermodel.OrderHistory, shipping *ordermodel.OrderShipping) *exttypes.OrderShipping {
	if shipping == nil {
		shipping = &ordermodel.OrderShipping{}
	}
	res := &exttypes.OrderShipping{
		GrossWeight:      order.TotalWeight().Int(),
		Length:           cmapi.PbPtrInt(shipping.Length),
		Width:            cmapi.PbPtrInt(shipping.Width),
		Height:           cmapi.PbPtrInt(shipping.Height),
		ChargeableWeight: order.TotalWeight().Int(),
	}
	return res
}

func PbOrderLines(items []*ordermodel.OrderLine) []*exttypes.OrderLine {
	// send changes as empty instead of "[]"
	if len(items) == 0 {
		return nil
	}
	res := make([]*exttypes.OrderLine, len(items))
	for i, item := range items {
		res[i] = PbOrderLine(item)
	}
	return res
}

func PbOrderLine(m *ordermodel.OrderLine) *exttypes.OrderLine {
	if m == nil {
		return nil
	}
	return &exttypes.OrderLine{
		VariantId:    m.VariantID,
		ProductId:    m.ProductID,
		ProductName:  m.ProductName,
		Quantity:     m.Quantity,
		ListPrice:    dot.Int(m.ListPrice),
		RetailPrice:  dot.Int(m.RetailPrice),
		PaymentPrice: dot.Int(m.PaymentPrice),
		ImageUrl:     m.ImageURL,
		Attributes:   convert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(m.Attributes),
	}
}

func PbOrderAddress(m *ordermodel.OrderAddress) *exttypes.OrderAddress {
	if m == nil {
		return nil
	}
	return &exttypes.OrderAddress{
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

func OrderAddressToPbCustomer(m *exttypes.OrderAddress) *types.OrderCustomer {
	if m == nil {
		return nil
	}
	return &types.OrderCustomer{
		FirstName: "",
		LastName:  "",
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Gender:    0,
	}
}

func OrderAddressToPbOrder(m *exttypes.OrderAddress) *types.OrderAddress {
	if m == nil {
		return nil
	}
	return &types.OrderAddress{
		FullName: m.FullName,
		Phone:    m.Phone,
		Province: m.Province,
		District: m.District,
		Ward:     m.Ward,
		Address1: m.Address1,
		Address2: m.Address2,
	}
}

func PbOrderAddressFromAddress(m *addressmodel.Address) *exttypes.OrderAddress {
	if m == nil {
		return nil
	}
	return &exttypes.OrderAddress{
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

func PbFulfillments(items []*shipmodel.Fulfillment) []*exttypes.Fulfillment {
	res := make([]*exttypes.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillment(item)
	}
	return res
}

func PbFulfillment(m *shipmodel.Fulfillment) *exttypes.Fulfillment {
	return &exttypes.Fulfillment{
		Id:                       m.ID,
		OrderId:                  m.OrderID,
		ShopId:                   m.ShopID,
		SelfUrl:                  PNonZeroString(m.SelfURL(cm.MainSiteBaseURL(), model.TagShop)),
		TotalItems:               cmapi.PbPtrInt(m.TotalItems),
		BasketValue:              cmapi.PbPtrInt(m.BasketValue),
		CreatedAt:                cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                cmapi.PbTime(m.UpdatedAt),
		ClosedAt:                 cmapi.PbTime(m.ClosedAt),
		CancelledAt:              cmapi.PbTime(m.ShippingCancelledAt),
		CancelReason:             dot.String(m.CancelReason),
		Carrier:                  m.ShippingProvider,
		ShippingServiceName:      dot.String(m.ExternalShippingName),
		ShippingServiceFee:       cmapi.PbPtrInt(m.ShippingServiceFee),
		ActualShippingServiceFee: cmapi.PbPtrInt(m.ShippingFeeShop),
		ShippingServiceCode:      dot.String(m.ProviderServiceID),
		ShippingCode:             dot.String(m.ShippingCode),
		ShippingNote:             dot.String(m.ShippingNote),
		TryOn:                    m.TryOn,
		IncludeInsurance:         m.IncludeInsurance,
		ConfirmStatus:            m.ConfirmStatus.Wrap(),
		ShippingState:            m.ShippingState.Wrap(),
		ShippingStatus:           m.ShippingStatus.Wrap(),
		Status:                   m.Status.Wrap(),
		CodAmount:                cmapi.PbPtrInt(m.OriginalCODAmount),
		ActualCodAmount:          cmapi.PbPtrInt(m.TotalCODAmount),
		ChargeableWeight:         cmapi.PbPtrInt(m.TotalWeight),
		PickupAddress:            PbOrderAddressFromAddress(m.AddressFrom),
		ReturnAddress:            PbOrderAddressFromAddress(m.AddressReturn),
		ShippingAddress:          PbOrderAddressFromAddress(m.AddressTo),
		EtopPaymentStatus:        m.EtopPaymentStatus.Wrap(),
		EstimatedDeliveryAt:      cmapi.PbTime(m.ExpectedDeliveryAt),
		EstimatedPickupAt:        cmapi.PbTime(m.ExpectedPickAt),
	}
}

func PbFulfillmentHistories(items []shipmodel.FulfillmentHistory) []*exttypes.Fulfillment {
	res := make([]*exttypes.Fulfillment, len(items))
	for i, item := range items {
		res[i] = PbFulfillmentHistory(item)
	}
	return res
}

func PbFulfillmentHistory(m shipmodel.FulfillmentHistory) *exttypes.Fulfillment {
	var addressTo, addressFrom, addressReturn *ordermodel.OrderAddress
	_ = m.AddressTo().Unmarshal(&addressTo)
	_ = m.AddressFrom().Unmarshal(&addressFrom)
	_ = m.AddressReturn().Unmarshal(&addressReturn)

	return &exttypes.Fulfillment{
		Id:          m.ID().ID().Apply(0),
		OrderId:     m.OrderID().ID().Apply(0),
		ShopId:      m.ShopID().ID().Apply(0),
		TotalItems:  m.TotalItems().Int(),
		BasketValue: m.BasketValue().Int(),
		CreatedAt:   cmapi.PbTime(m.CreatedAt().Time()),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt().Time()),
		ClosedAt:    cmapi.PbTime(m.ClosedAt().Time()),
		// CancelledAt: nothing (TODO: fix it)
		CancelReason:             m.CancelReason().String(),
		Carrier:                  convertpb.PbShippingProviderPtr(m.ShippingProvider().String()).Apply(0),
		ShippingServiceName:      m.ExternalShippingName().String(),
		ShippingServiceFee:       m.ShippingServiceFee().Int(),
		ActualShippingServiceFee: m.ShippingFeeShop().Int(),
		ShippingServiceCode:      m.ProviderServiceID().String(),
		ShippingCode:             m.ShippingCode().String(),
		ShippingNote:             m.ShippingNote().String(),
		TryOn:                    convertpb.PbTryOnPtr(m.TryOn().String()).Apply(0),
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

func PbShippingServices(items []*shippingsharemodel.AvailableShippingService) []*exttypes.ShippingService {
	res := make([]*exttypes.ShippingService, len(items))
	for i, item := range items {
		res[i] = PbShippingService(item)
	}
	return res
}

func PbShippingService(m *shippingsharemodel.AvailableShippingService) *exttypes.ShippingService {
	return &exttypes.ShippingService{
		Code:                m.ProviderServiceID,
		Name:                m.Name,
		Fee:                 m.ServiceFee,
		Carrier:             m.Provider,
		EstimatedPickupAt:   cmapi.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt: cmapi.PbTime(m.ExpectedDeliveryAt),
		CarrierInfo:         PbCarrierInfo(m.ConnectionInfo),
	}
}

func PbCarrierInfo(m *shippingsharemodel.ConnectionInfo) *exttypes.CarrierInfo {
	if m == nil {
		return nil
	}
	return &exttypes.CarrierInfo{
		Name:     m.Name,
		ImageURL: m.ImageURL,
	}
}

func OrderLinesToCreateOrderLines(items []*exttypes.OrderLine) (_ []*types.CreateOrderLine, err error) {
	res := make([]*types.CreateOrderLine, len(items))
	for i, item := range items {
		res[i], err = OrderLineToCreateOrderLine(item)
		if err != nil {
			return
		}
	}
	return res, nil
}

func OrderLineToCreateOrderLine(m *exttypes.OrderLine) (*types.CreateOrderLine, error) {
	if m == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "order_line must not be null")
	}
	if !m.PaymentPrice.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp payment_price")
	}

	return &types.CreateOrderLine{
		VariantId:    m.VariantId,
		ProductName:  m.ProductName,
		Quantity:     m.Quantity,
		ListPrice:    m.ListPrice.Apply(0),
		RetailPrice:  m.RetailPrice.Apply(0),
		PaymentPrice: m.PaymentPrice.Apply(0),
		ImageUrl:     m.ImageUrl,
		Attributes:   m.Attributes,
	}, nil
}

func PbOrderAndFulfillments(order *ordermodel.Order, fulfillments []*shipmodel.Fulfillment) *exttypes.OrderAndFulfillments {
	return &exttypes.OrderAndFulfillments{
		Order:        PbOrder(order),
		Fulfillments: PbFulfillments(fulfillments),
	}
}

func PbPageInfo(arg *cm.Paging, metaPageInfo *meta.PageInfo) *common.CursorPageInfo {
	if metaPageInfo == nil {
		return nil
	}
	if arg == nil {
		return nil
	}
	sort := ""
	if len(arg.Sort) > 0 {
		sort = arg.Sort[0]
	}
	return &common.CursorPageInfo{
		Before: arg.Before,
		After:  arg.After,

		Limit: metaPageInfo.Limit,
		Sort:  sort,

		Prev: metaPageInfo.Prev,
		Next: metaPageInfo.Next,
	}
}

func PbShopCustomer(customer *customering.ShopCustomer) *exttypes.Customer {
	if customer == nil {
		return nil
	}
	if customer.Deleted {
		return &exttypes.Customer{
			Id:      customer.ID,
			Deleted: true,
		}
	}
	return &exttypes.Customer{
		Id:           customer.ID,
		ShopId:       customer.ShopID,
		ExternalId:   dot.String(customer.ExternalID),
		ExternalCode: dot.String(customer.ExternalCode),
		FullName:     dot.String(customer.FullName),
		Code:         dot.String(customer.Code),
		Note:         dot.String(customer.Note),
		Phone:        dot.String(customer.Phone),
		Email:        dot.String(customer.Email),
		Gender:       customer.Gender.Wrap(),
		Type:         customer.Type.Wrap(),
		Birthday:     dot.String(customer.Birthday),
		CreatedAt:    dot.Time(customer.CreatedAt),
		UpdatedAt:    dot.Time(customer.UpdatedAt),
		Status:       customer.Status.Wrap(),
		Deleted:      customer.Deleted,
	}
}

func PbShopCustomerHistory(m customermodel.ShopCustomerHistory) *exttypes.Customer {
	return &exttypes.Customer{
		ExternalId:   m.ExternalID().String(),
		ExternalCode: m.ExternalCode().String(),
		Id:           m.ID().ID().Apply(0),
		ShopId:       m.ShopID().ID().Apply(0),
		FullName:     m.FullName().String(),
		Code:         m.Code().String(),
		Note:         m.Note().String(),
		Phone:        m.Phone().String(),
		Email:        m.Email().String(),
		Gender:       gender.ParseGenderWithNull(m.Gender().String(), gender.Unknown),
		Type:         customer_type.ParseCustomerTypeWithNull(m.Type().String(), customer_type.Unknown),
		Birthday:     dot.String(cmapi.PbTime(m.Birthday().Time()).String()),
		CreatedAt:    cmapi.PbTime(m.CreatedAt().Time()),
		UpdatedAt:    cmapi.PbTime(m.UpdatedAt().Time()),
		Status:       convertpb.Pb3Ptr(m.Status().Int()),
	}
}

func PbShopCustomers(customers []*customering.ShopCustomer) []*exttypes.Customer {
	out := make([]*exttypes.Customer, len(customers))
	for i, customer := range customers {
		out[i] = PbShopCustomer(customer)
	}
	return out
}

func PbCoordinates(in *ordertypes.Coordinates) *etop.Coordinates {
	if in == nil {
		return nil
	}
	return &etop.Coordinates{
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}

func PbShopTraderAddressHistory(ctx context.Context, m customermodel.ShopTraderAddressHistory, locationBus location.QueryBus) *exttypes.CustomerAddress {
	query := &location.GetLocationQuery{
		DistrictCode: m.DistrictCode().String().Apply(""),
		WardCode:     m.WardCode().String().Apply(""),
	}
	province, district, ward := (*location.Province)(nil), (*location.District)(nil), (*location.Ward)(nil)
	if query.DistrictCode != "" && query.WardCode != "" {
		if err := locationBus.Dispatch(ctx, query); err != nil {
			panic("Internal: " + err.Error())
		}
		province, district, ward = query.Result.Province, query.Result.District, query.Result.Ward
	}
	out := &exttypes.CustomerAddress{
		Id:          m.ID().ID().Apply(0),
		CustomerID:  m.TraderID().ID().Apply(0),
		Address1:    m.Address1().String(),
		Address2:    m.Address2().String(),
		FullName:    m.FullName().String(),
		Company:     m.Company().String(),
		Phone:       m.Phone().String(),
		Email:       m.Email().String(),
		Position:    m.Position().String(),
		Coordinates: nil, // TODO
	}
	if ward != nil {
		out.Ward = dot.String(ward.Name)
		out.WardCode = dot.String(ward.Code)
	}
	if district != nil {
		out.District = dot.String(district.Name)
	}
	if province != nil {
		out.Province = dot.String(province.Name)
		out.ProvinceCode = dot.String(province.Code)
	}
	return out
}

func PbShopTraderAddress(ctx context.Context, in *addressing.ShopTraderAddress, locationBus location.QueryBus) *exttypes.CustomerAddress {
	if in.Deleted {
		return &exttypes.CustomerAddress{
			Id:      in.ID,
			Deleted: true,
		}
	}
	query := &location.GetLocationQuery{
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		panic("Internal: " + err.Error())
	}
	province, district, ward := query.Result.Province, query.Result.District, query.Result.Ward
	out := &exttypes.CustomerAddress{
		Id:          in.ID,
		CustomerID:  in.TraderID,
		Company:     dot.String(in.Company),
		Address1:    dot.String(in.Address1),
		Address2:    dot.String(in.Address2),
		FullName:    dot.String(in.FullName),
		Phone:       dot.String(in.Phone),
		Email:       dot.String(in.Email),
		Position:    dot.String(in.Position),
		Coordinates: PbCoordinates(in.Coordinates),
	}
	if ward != nil {
		out.Ward = dot.String(ward.Name)
		out.WardCode = dot.String(ward.Code)
	}
	if district != nil {
		out.District = dot.String(district.Name)
		out.DistrictCode = dot.String(district.Code)
	}
	if province != nil {
		out.Province = dot.String(province.Name)
		out.ProvinceCode = dot.String(province.Code)
	}
	return out
}

func PbShopTraderAddresses(ctx context.Context, ins []*addressing.ShopTraderAddress, locationBus location.QueryBus) []*exttypes.CustomerAddress {
	out := make([]*exttypes.CustomerAddress, len(ins))
	for i, trader := range ins {
		out[i] = PbShopTraderAddress(ctx, trader, locationBus)
	}
	return out
}

func PbCustomerGroupHistory(m customermodel.ShopCustomerGroupHistory) *exttypes.CustomerGroup {
	return &exttypes.CustomerGroup{
		Id:     m.ID().ID().Apply(0),
		ShopID: m.ShopID().ID().Apply(0),
		Name:   m.Name().String(),
	}
}

func PbCustomerGroup(arg *customering.ShopCustomerGroup) *exttypes.CustomerGroup {
	if arg == nil {
		return nil
	}
	if arg.Deleted {
		return &exttypes.CustomerGroup{
			Id:      arg.ID,
			Deleted: true,
		}
	}
	return &exttypes.CustomerGroup{
		Id:     arg.ID,
		ShopID: arg.ShopID,
		Name:   dot.String(arg.Name),
	}
}

func PbCustomerGroups(args []*customering.ShopCustomerGroup) []*exttypes.CustomerGroup {
	out := make([]*exttypes.CustomerGroup, len(args))
	for i, arg := range args {
		out[i] = PbCustomerGroup(arg)
	}
	return out
}

func PbInventoryVariantHistory(m inventorymodel.InventoryVariantHistory) *exttypes.InventoryLevel {
	return &exttypes.InventoryLevel{
		VariantId:         m.VariantID().ID().Apply(0),
		AvailableQuantity: m.QuantityOnHand().Int(),
		PickedQuantity:    m.QuantityPicked().Int(),
		UpdatedAt:         cmapi.PbTime(m.UpdatedAt().Time()),
	}
}

func PbInventoryLevel(arg *inventory.InventoryVariant) *exttypes.InventoryLevel {
	if arg == nil {
		return nil
	}
	return &exttypes.InventoryLevel{
		VariantId:         arg.VariantID,
		AvailableQuantity: dot.Int(arg.QuantityOnHand),
		PickedQuantity:    dot.Int(arg.QuantityPicked),
		UpdatedAt:         dot.Time(arg.UpdatedAt),
	}
}

func PbInventoryLevels(args []*inventory.InventoryVariant) []*exttypes.InventoryLevel {
	out := make([]*exttypes.InventoryLevel, len(args))
	for i, arg := range args {
		out[i] = PbInventoryLevel(arg)
	}
	return out
}

func PbShopProductCollectionHistory(m catalogmodel.ShopCollectionHistory) *exttypes.ProductCollection {
	return &exttypes.ProductCollection{
		ID:          m.ID().ID().Apply(0),
		ShopID:      m.ShopID().ID().Apply(0),
		Name:        m.Name().String(),
		Description: m.Description().String(),
		ShortDesc:   m.ShortDesc().String(),
		CreatedAt:   cmapi.PbTime(m.CreatedAt().Time()),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt().Time()),
	}
}

func PbShopProductCollection(arg *catalog.ShopCollection) *exttypes.ProductCollection {
	if arg == nil {
		return nil
	}
	if arg.Deleted {
		return &exttypes.ProductCollection{
			ID:      arg.ID,
			Deleted: true,
		}
	}
	return &exttypes.ProductCollection{
		ID:          arg.ID,
		ShopID:      arg.ShopID,
		Name:        dot.String(arg.Name),
		Description: dot.String(arg.Description),
		ShortDesc:   dot.String(arg.ShortDesc),
		CreatedAt:   cmapi.PbTime(arg.CreatedAt),
		UpdatedAt:   cmapi.PbTime(arg.UpdatedAt),
	}
}

func PbShopProductCollectionRelationship(arg *catalog.ShopProductCollection) *exttypes.ProductCollectionRelationship {
	if arg == nil {
		return nil
	}
	return &exttypes.ProductCollectionRelationship{
		ProductId:    arg.ProductID,
		CollectionId: arg.CollectionID,
		Deleted:      arg.Deleted,
	}
}

func PbShopProductCollectionRelationships(args []*catalog.ShopProductCollection) []*exttypes.ProductCollectionRelationship {
	outs := make([]*exttypes.ProductCollectionRelationship, len(args))
	for i, arg := range args {
		outs[i] = PbShopProductCollectionRelationship(arg)
	}
	return outs
}

func PbShopProductCollections(args []*catalog.ShopCollection) []*exttypes.ProductCollection {
	outs := make([]*exttypes.ProductCollection, len(args))
	for i, arg := range args {
		outs[i] = PbShopProductCollection(arg)
	}
	return outs
}

func PbCustomerGroupRelationship(args *customering.CustomerGroupCustomer) *exttypes.CustomerGroupRelationship {
	if args == nil {
		return nil
	}
	return &exttypes.CustomerGroupRelationship{
		CustomerID: args.CustomerID,
		GroupID:    args.GroupID,
		Deleted:    args.Deleted,
	}
}

func PbCustomerGroupRelationships(args []*customering.CustomerGroupCustomer) []*exttypes.CustomerGroupRelationship {
	outs := make([]*exttypes.CustomerGroupRelationship, len(args))
	for i, arg := range args {
		outs[i] = PbCustomerGroupRelationship(arg)
	}
	return outs
}

func PbShopCustomerGroupCustomerHistory(m customermodel.ShopCustomerGroupCustomerHistory) *exttypes.CustomerGroupRelationship {
	return &exttypes.CustomerGroupRelationship{
		CustomerID: m.CustomerID().ID().Apply(0),
		GroupID:    m.GroupID().ID().Apply(0),
	}
}

func PbShopProductionCollectionRelationshipHistory(m catalogmodel.ProductShopCollectionHistory) *exttypes.ProductCollectionRelationship {
	return &exttypes.ProductCollectionRelationship{
		ProductId:    m.ProductID().ID().Apply(0),
		CollectionId: m.CollectionID().ID().Apply(0),
	}
}
