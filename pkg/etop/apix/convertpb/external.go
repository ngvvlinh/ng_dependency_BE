package convertpb

import (
	"etop.vn/api/main/catalog"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/meta"
	"etop.vn/api/shopping/customering"
	exttypes "etop.vn/api/top/external/types"
	"etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/account_type"
	"etop.vn/backend/com/handler/etop-handler/webhook/sender"
	"etop.vn/backend/com/main/catalog/convert"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/capi/util"
)

func PNonZeroString(s string) dot.NullString {
	if s == "" {
		return dot.NullString{}
	}
	return dot.String(s)
}

func PbPartner(m *model.Partner) *exttypes.Partner {
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
		ListPrice:    m.ListPrice,
		RetailPrice:  m.RetailPrice,
		PaymentPrice: cmapi.PbPtrInt(m.PaymentPrice),
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

func PbOrderAddressFromAddress(m *model.Address) *exttypes.OrderAddress {
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
		IncludeInsurance:         dot.Bool(m.IncludeInsurance),
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
		Id:                       m.ID().ID().Apply(0),
		OrderId:                  m.OrderID().ID().Apply(0),
		ShopId:                   m.ShopID().ID().Apply(0),
		TotalItems:               m.TotalItems().Int(),
		BasketValue:              m.BasketValue().Int(),
		CreatedAt:                cmapi.PbTime(m.CreatedAt().Time()),
		UpdatedAt:                cmapi.PbTime(m.UpdatedAt().Time()),
		ClosedAt:                 cmapi.PbTime(m.ClosedAt().Time()),
		CancelledAt:              cmapi.PbTime(m.CancelReason().Time()),
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

func PbShippingServices(items []*model.AvailableShippingService) []*exttypes.ShippingService {
	res := make([]*exttypes.ShippingService, len(items))
	for i, item := range items {
		res[i] = PbShippingService(item)
	}
	return res
}

func PbShippingService(m *model.AvailableShippingService) *exttypes.ShippingService {
	return &exttypes.ShippingService{
		Code:                m.ProviderServiceID,
		Name:                m.Name,
		Fee:                 m.ServiceFee,
		Carrier:             m.Provider,
		EstimatedPickupAt:   cmapi.PbTime(m.ExpectedPickAt),
		EstimatedDeliveryAt: cmapi.PbTime(m.ExpectedDeliveryAt),
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
		ListPrice:    m.ListPrice,
		RetailPrice:  m.RetailPrice,
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

func PbPageInfo(arg *common.CursorPaging, metaPageInfo *meta.PageInfo) *common.CursorPageInfo {
	if metaPageInfo == nil {
		return nil
	}
	if arg == nil {
		return nil
	}
	return &common.CursorPageInfo{
		First:  arg.First,
		Last:   arg.Last,
		Before: arg.Before,
		After:  arg.After,

		Limit: metaPageInfo.Limit,
		Sort:  arg.Sort,

		Prev: metaPageInfo.Prev,
		Next: metaPageInfo.Next,
	}
}

func PbShopCustomer(customer *customering.ShopCustomer) *exttypes.Customer {
	if customer == nil {
		return nil
	}
	return &exttypes.Customer{
		Id:           customer.ID,
		ShopId:       customer.ShopID,
		ExternalId:   customer.ExternalID,
		ExternalCode: customer.ExternalCode,
		FullName:     customer.FullName,
		Code:         customer.Code,
		Note:         customer.Note,
		Phone:        customer.Phone,
		Email:        customer.Email,
		Gender:       customer.Gender.String(),
		Type:         customer.Type.String(),
		Birthday:     customer.Birthday,
		CreatedAt:    dot.Time(customer.CreatedAt),
		UpdatedAt:    dot.Time(customer.UpdatedAt),
		Status:       customer.Status,
		GroupIds:     customer.GroupIDs,
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

func PbShopProduct(arg *catalog.ShopProduct) *exttypes.ShopProduct {
	if arg == nil {
		return nil
	}
	return &exttypes.ShopProduct{
		ExternalId:    arg.ExternalID,
		ExternalCode:  arg.ExternalCode,
		Id:            arg.ProductID,
		Name:          arg.Name,
		Description:   arg.Description,
		ShortDesc:     arg.ShortDesc,
		DescHtml:      arg.DescHTML,
		ImageUrls:     arg.ImageURLs,
		CategoryId:    arg.CategoryID,
		Tags:          arg.Tags,
		Note:          arg.Note,
		Status:        arg.Status,
		ListPrice:     arg.ListPrice,
		RetailPrice:   arg.RetailPrice,
		CollectionIds: arg.CollectionIDs,
		CreatedAt:     cmapi.PbTime(arg.CreatedAt),
		UpdatedAt:     cmapi.PbTime(arg.UpdatedAt),
		BrandId:       arg.BrandID,
	}
}

func PbShopProducts(args []*catalog.ShopProduct) []*exttypes.ShopProduct {
	outs := make([]*exttypes.ShopProduct, len(args))
	for i, arg := range args {
		outs[i] = PbShopProduct(arg)
	}
	return outs
}

func ConvertProductWithVariantsToPbProduct(arg *catalog.ShopProductWithVariants) *exttypes.ShopProduct {
	if arg == nil {
		return nil
	}
	return &exttypes.ShopProduct{
		ExternalId:    arg.ExternalID,
		ExternalCode:  arg.ExternalCode,
		Id:            arg.ProductID,
		Name:          arg.Name,
		Description:   arg.Description,
		ShortDesc:     arg.ShortDesc,
		DescHtml:      arg.DescHTML,
		ImageUrls:     arg.ImageURLs,
		CategoryId:    arg.CategoryID,
		Tags:          arg.Tags,
		Note:          arg.Note,
		Status:        arg.Status,
		ListPrice:     arg.ListPrice,
		RetailPrice:   arg.RetailPrice,
		CollectionIds: arg.CollectionIDs,
		CreatedAt:     cmapi.PbTime(arg.CreatedAt),
		UpdatedAt:     cmapi.PbTime(arg.UpdatedAt),
		BrandId:       arg.BrandID,
	}
}

func PbShopVariant(arg *catalog.ShopVariant) *exttypes.ShopVariant {
	if arg == nil {
		return nil
	}
	return &exttypes.ShopVariant{
		ExternalId:   arg.ExternalID,
		ExternalCode: arg.ExternalCode,
		Id:           arg.VariantID,
		Code:         arg.Code,
		Name:         arg.Name,
		Description:  arg.Description,
		ShortDesc:    arg.ShortDesc,
		DescHtml:     arg.DescHTML,
		ImageUrls:    arg.ImageURLs,
		ListPrice:    arg.ListPrice,
		RetailPrice:  util.CoalesceInt(arg.RetailPrice, arg.ListPrice),
		Note:         arg.Note,
		Status:       arg.Status,
		CostPrice:    arg.CostPrice,
		Attributes:   arg.Attributes,
	}
}

func PbShopVariants(args []*catalog.ShopVariant) []*exttypes.ShopVariant {
	outs := make([]*exttypes.ShopVariant, len(args))
	for i, arg := range args {
		outs[i] = PbShopVariant(arg)
	}
	return outs
}
