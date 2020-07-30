package shop

import (
	"context"

	etop "o.o/api/top/int/etop"
	shoptypes "o.o/api/top/int/shop/types"
	"o.o/api/top/int/types"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop/shop

// +apix:path=/shop.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/shop.Brand
type BrandService interface {
	CreateBrand(context.Context, *CreateBrandRequest) (*Brand, error)
	UpdateBrandInfo(context.Context, *UpdateBrandRequest) (*Brand, error)
	DeleteBrand(context.Context, *cm.IDsRequest) (*DeleteBrandResponse, error)

	GetBrandByID(context.Context, *cm.IDRequest) (*Brand, error)
	GetBrandsByIDs(context.Context, *cm.IDsRequest) (*GetBrandsByIDsResponse, error)
	GetBrands(context.Context, *GetBrandsRequest) (*GetBrandsResponse, error)
}

// +apix:path=/shop.Inventory
type InventoryService interface {
	CreateInventoryVoucher(context.Context, *CreateInventoryVoucherRequest) (*CreateInventoryVoucherResponse, error)
	ConfirmInventoryVoucher(context.Context, *ConfirmInventoryVoucherRequest) (*ConfirmInventoryVoucherResponse, error)
	CancelInventoryVoucher(context.Context, *CancelInventoryVoucherRequest) (*CancelInventoryVoucherResponse, error)
	UpdateInventoryVoucher(context.Context, *UpdateInventoryVoucherRequest) (*UpdateInventoryVoucherResponse, error)
	AdjustInventoryQuantity(context.Context, *AdjustInventoryQuantityRequest) (*AdjustInventoryQuantityResponse, error)
	UpdateInventoryVariantCostPrice(context.Context, *UpdateInventoryVariantCostPriceRequest) (*UpdateInventoryVariantCostPriceResponse, error)

	GetInventoryVariant(context.Context, *GetInventoryVariantRequest) (*InventoryVariant, error)
	GetInventoryVariants(context.Context, *GetInventoryVariantsRequest) (*GetInventoryVariantsResponse, error)
	GetInventoryVariantsByVariantIDs(context.Context, *GetInventoryVariantsByVariantIDsRequest) (*GetInventoryVariantsResponse, error)
	GetInventoryVouchersByReference(context.Context, *GetInventoryVouchersByReferenceRequest) (*GetInventoryVouchersByReferenceResponse, error)
	GetInventoryVoucher(context.Context, *cm.IDRequest) (*InventoryVoucher, error)
	GetInventoryVouchers(context.Context, *GetInventoryVouchersRequest) (*GetInventoryVouchersResponse, error)
	GetInventoryVouchersByIDs(context.Context, *GetInventoryVouchersByIDsRequest) (*GetInventoryVouchersResponse, error)
}

// +apix:path=/shop.Account
type AccountService interface {
	RegisterShop(context.Context, *RegisterShopRequest) (*RegisterShopResponse, error)
	UpdateShop(context.Context, *UpdateShopRequest) (*UpdateShopResponse, error)
	DeleteShop(context.Context, *cm.IDRequest) (*cm.Empty, error)
	SetDefaultAddress(context.Context, *etop.SetDefaultAddressRequest) (*cm.UpdatedResponse, error)

	CreateExternalAccountAhamove(context.Context, *cm.Empty) (*ExternalAccountAhamove, error)
	GetExternalAccountAhamove(context.Context, *cm.Empty) (*ExternalAccountAhamove, error)
	RequestVerifyExternalAccountAhamove(context.Context, *cm.Empty) (*cm.UpdatedResponse, error)
	UpdateExternalAccountAhamoveVerification(context.Context, *UpdateXAccountAhamoveVerificationRequest) (*cm.UpdatedResponse, error)

	// deprecated: backward-compatible, will be removed later
	UpdateExternalAccountAhamoveVerificationImages(context.Context, *UpdateXAccountAhamoveVerificationRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Collection
type CollectionService interface {
	CreateCollection(context.Context, *CreateCollectionRequest) (*ShopCollection, error)
	GetCollection(context.Context, *cm.IDRequest) (*ShopCollection, error)
	GetCollections(context.Context, *GetCollectionsRequest) (*ShopCollectionsResponse, error)
	UpdateCollection(context.Context, *UpdateCollectionRequest) (*ShopCollection, error)
	GetCollectionsByProductID(context.Context, *GetShopCollectionsByProductIDRequest) (*CollectionsResponse, error)
}

// +apix:path=/shop.Customer
type CustomerService interface {
	CreateCustomer(context.Context, *CreateCustomerRequest) (*Customer, error)
	UpdateCustomer(context.Context, *UpdateCustomerRequest) (*Customer, error)
	DeleteCustomer(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	GetCustomer(context.Context, *cm.IDRequest) (*Customer, error)
	GetCustomerDetails(context.Context, *cm.IDRequest) (*CustomerDetailsResponse, error)
	GetCustomers(context.Context, *GetCustomersRequest) (*CustomersResponse, error)
	GetCustomersByIDs(context.Context, *cm.IDsRequest) (*CustomersResponse, error)
	BatchSetCustomersStatus(context.Context, *SetCustomersStatusRequest) (*cm.UpdatedResponse, error)

	//-- address --//

	GetCustomerAddresses(context.Context, *GetCustomerAddressesRequest) (*CustomerAddressesResponse, error)
	CreateCustomerAddress(context.Context, *CreateCustomerAddressRequest) (*CustomerAddress, error)
	UpdateCustomerAddress(context.Context, *UpdateCustomerAddressRequest) (*CustomerAddress, error)
	SetDefaultCustomerAddress(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	DeleteCustomerAddress(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)

	//-- group --//
	AddCustomersToGroup(context.Context, *AddCustomerToGroupRequest) (*cm.UpdatedResponse, error)
	RemoveCustomersFromGroup(context.Context, *RemoveCustomerOutOfGroupRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/shop.CustomerGroup
type CustomerGroupService interface {
	CreateCustomerGroup(context.Context, *CreateCustomerGroupRequest) (*CustomerGroup, error)
	GetCustomerGroup(context.Context, *cm.IDRequest) (*CustomerGroup, error)
	GetCustomerGroups(context.Context, *GetCustomerGroupsRequest) (*CustomerGroupsResponse, error)
	UpdateCustomerGroup(context.Context, *UpdateCustomerGroupRequest) (*CustomerGroup, error)
}

// +apix:path=/shop.Product
type ProductService interface {

	//-- product --//

	GetProduct(context.Context, *cm.IDRequest) (*ShopProduct, error)
	GetProducts(context.Context, *GetVariantsRequest) (*ShopProductsResponse, error)
	GetProductsByIDs(context.Context, *cm.IDsRequest) (*ShopProductsResponse, error)

	CreateProduct(context.Context, *CreateProductRequest) (*ShopProduct, error)
	UpdateProduct(context.Context, *UpdateProductRequest) (*ShopProduct, error)
	UpdateProductsStatus(context.Context, *UpdateProductStatusRequest) (*UpdateProductStatusResponse, error)
	UpdateProductsTags(context.Context, *UpdateProductsTagsRequest) (*cm.UpdatedResponse, error)
	UpdateProductImages(context.Context, *UpdateVariantImagesRequest) (*ShopProduct, error)
	UpdateProductMetaFields(context.Context, *UpdateProductMetaFieldsRequest) (*ShopProduct, error)
	RemoveProducts(context.Context, *RemoveVariantsRequest) (*cm.RemovedResponse, error)
	//-- variant --//

	GetVariant(context.Context, *GetVariantRequest) (*ShopVariant, error)
	GetVariantsByIDs(context.Context, *cm.IDsRequest) (*ShopVariantsResponse, error)

	CreateVariant(context.Context, *CreateVariantRequest) (*ShopVariant, error)
	UpdateVariant(context.Context, *UpdateVariantRequest) (*ShopVariant, error)
	UpdateVariantImages(context.Context, *UpdateVariantImagesRequest) (*ShopVariant, error)
	UpdateVariantsStatus(context.Context, *UpdateProductStatusRequest) (*UpdateProductStatusResponse, error)
	UpdateVariantAttributes(context.Context, *UpdateVariantAttributesRequest) (*ShopVariant, error)
	RemoveVariants(context.Context, *RemoveVariantsRequest) (*cm.RemovedResponse, error)
	GetVariantsBySupplierID(context.Context, *GetVariantsBySupplierIDRequest) (*ShopVariantsResponse, error)
	//-- category --//
	UpdateProductCategory(context.Context, *UpdateProductCategoryRequest) (*ShopProduct, error)
	RemoveProductCategory(context.Context, *cm.IDRequest) (*ShopProduct, error)

	//-- collection --//
	AddProductCollection(context.Context, *AddShopProductCollectionRequest) (*cm.UpdatedResponse, error)
	RemoveProductCollection(context.Context, *RemoveShopProductCollectionRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/shop.Category
type CategoryService interface {
	CreateCategory(context.Context, *CreateCategoryRequest) (*ShopCategory, error)
	GetCategory(context.Context, *cm.IDRequest) (*ShopCategory, error)
	GetCategories(context.Context, *GetCategoriesRequest) (*ShopCategoriesResponse, error)
	UpdateCategory(context.Context, *UpdateCategoryRequest) (*ShopCategory, error)
	DeleteCategory(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.ProductSource
// deprecated: 2018.07.31+14
type ProductSourceService interface {
	CreateProductSource(context.Context, *CreateProductSourceRequest) (*ProductSource, error)
	GetShopProductSources(context.Context, *cm.Empty) (*ProductSourcesResponse, error)
	// deprecated: use shop.Product/CreateVariant instead
	CreateVariant(context.Context, *DeprecatedCreateVariantRequest) (*ShopProduct, error)
	CreateProductSourceCategory(context.Context, *CreatePSCategoryRequest) (*Category, error)
	UpdateProductsPSCategory(context.Context, *UpdateProductsPSCategoryRequest) (*cm.UpdatedResponse, error)
	GetProductSourceCategory(context.Context, *cm.IDRequest) (*Category, error)
	GetProductSourceCategories(context.Context, *GetProductSourceCategoriesRequest) (*CategoriesResponse, error)
	UpdateProductSourceCategory(context.Context, *UpdateProductSourceCategoryRequest) (*Category, error)
	RemoveProductSourceCategory(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/shop.Order
type OrderService interface {
	CreateOrder(context.Context, *types.CreateOrderRequest) (*types.Order, error)
	GetOrder(context.Context, *cm.IDRequest) (*types.Order, error)
	GetOrders(context.Context, *GetOrdersRequest) (*types.OrdersResponse, error)
	GetOrdersByIDs(context.Context, *etop.IDsRequest) (*types.OrdersResponse, error)
	GetOrdersByReceiptID(context.Context, *GetOrdersByReceiptIDRequest) (*types.OrdersResponse, error)
	UpdateOrder(context.Context, *types.UpdateOrderRequest) (*types.Order, error)

	// @deprecated
	UpdateOrdersStatus(context.Context, *UpdateOrdersStatusRequest) (*cm.UpdatedResponse, error)

	ConfirmOrder(context.Context, *ConfirmOrderRequest) (*types.Order, error)
	ConfirmOrderAndCreateFulfillments(context.Context, *OrderIDRequest) (*types.OrderWithErrorsResponse, error)
	CancelOrder(context.Context, *CancelOrderRequest) (*types.OrderWithErrorsResponse, error)
	CompleteOrder(context.Context, *OrderIDRequest) (*cm.UpdatedResponse, error)
	UpdateOrderPaymentStatus(context.Context, *UpdateOrderPaymentStatusRequest) (*cm.UpdatedResponse, error)
	UpdateOrderShippingInfo(context.Context, *UpdateOrderShippingInfoRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Fulfillment
type FulfillmentService interface {
	GetFulfillment(context.Context, *cm.IDRequest) (*types.Fulfillment, error)
	GetFulfillments(context.Context, *GetFulfillmentsRequest) (*types.FulfillmentsResponse, error)
	GetFulfillmentsByIDs(context.Context, *GetFulfillmentsByIDsRequest) (*types.FulfillmentsResponse, error)

	GetPublicExternalShippingServices(context.Context, *types.GetExternalShippingServicesRequest) (*types.GetExternalShippingServicesResponse, error)
	GetExternalShippingServices(context.Context, *types.GetExternalShippingServicesRequest) (*types.GetExternalShippingServicesResponse, error)
	GetPublicFulfillment(context.Context, *GetPublicFulfillmentRequest) (*types.PublicFulfillment, error)
	UpdateFulfillmentsShippingState(context.Context, *UpdateFulfillmentsShippingStateRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Shipment
type ShipmentService interface {
	GetShippingServices(context.Context, *types.GetShippingServicesRequest) (*types.GetShippingServicesResponse, error)
	CreateFulfillments(context.Context, *CreateFulfillmentsRequest) (*CreateFulfillmentsResponse, error)
	UpdateFulfillmentCOD(context.Context, *UpdateFulfillmentCODRequest) (*cm.UpdatedResponse, error)
	UpdateFulfillmentInfo(context.Context, *UpdateFulfillmentInfoRequest) (*cm.UpdatedResponse, error)
	CancelFulfillment(context.Context, *CancelFulfillmentRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Shipnow
type ShipnowService interface {
	GetShipnowFulfillment(context.Context, *cm.IDRequest) (*types.ShipnowFulfillment, error)
	GetShipnowFulfillments(context.Context, *types.GetShipnowFulfillmentsRequest) (*types.ShipnowFulfillments, error)

	CreateShipnowFulfillmentV2(context.Context, *types.CreateShipnowFulfillmentV2Request) (*types.ShipnowFulfillment, error)

	ConfirmShipnowFulfillment(context.Context, *cm.IDRequest) (*types.ShipnowFulfillment, error)
	UpdateShipnowFulfillment(context.Context, *types.UpdateShipnowFulfillmentRequest) (*types.ShipnowFulfillment, error)
	CancelShipnowFulfillment(context.Context, *types.CancelShipnowFulfillmentRequest) (*cm.UpdatedResponse, error)

	GetShipnowServices(context.Context, *types.GetShipnowServicesRequest) (*types.GetShipnowServicesResponse, error)
}

// +apix:path=/shop.History
type HistoryService interface {
	GetFulfillmentHistory(context.Context, *GetFulfillmentHistoryRequest) (*etop.HistoryResponse, error)
}

// +apix:path=/shop.MoneyTransaction
type MoneyTransactionService interface {
	GetMoneyTransaction(context.Context, *cm.IDRequest) (*types.MoneyTransaction, error)
	GetMoneyTransactions(context.Context, *GetMoneyTransactionsRequest) (*types.MoneyTransactionsResponse, error)
}

// +apix:path=/shop.Summary
type SummaryService interface {
	SummarizeFulfillments(context.Context, *SummarizeFulfillmentsRequest) (*SummarizeFulfillmentsResponse, error)
	SummarizeTopShip(context.Context, *SummarizeTopShipRequest) (*SummarizeTopShipResponse, error)
	SummarizePOS(context.Context, *SummarizePOSRequest) (*SummarizePOSResponse, error)
	CalcBalanceUser(context.Context, *cm.Empty) (*CalcBalanceUserResponse, error)
}

// +apix:path=/shop.Export
type ExportService interface {
	GetExports(context.Context, *GetExportsRequest) (*GetExportsResponse, error)
	RequestExport(context.Context, *RequestExportRequest) (*RequestExportResponse, error)
}

// +apix:path=/shop.Notification
type NotificationService interface {
	CreateDevice(context.Context, *etop.CreateDeviceRequest) (*etop.Device, error)
	DeleteDevice(context.Context, *etop.DeleteDeviceRequest) (*cm.DeletedResponse, error)

	GetNotification(context.Context, *cm.IDRequest) (*etop.Notification, error)
	GetNotifications(context.Context, *etop.GetNotificationsRequest) (*etop.NotificationsResponse, error)
	UpdateNotifications(context.Context, *etop.UpdateNotificationsRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Authorize
type AuthorizeService interface {
	GetAuthorizedPartners(context.Context, *cm.Empty) (*GetAuthorizedPartnersResponse, error)
	GetAvailablePartners(context.Context, *cm.Empty) (*GetPartnersResponse, error)
	AuthorizePartner(context.Context, *AuthorizePartnerRequest) (*AuthorizedPartnerResponse, error)
}

// +apix:path=/shop.Trading
type TradingService interface {
	TradingGetProduct(context.Context, *cm.IDRequest) (*ShopProduct, error)
	TradingGetProducts(context.Context, *cm.CommonListRequest) (*ShopProductsResponse, error)

	TradingCreateOrder(context.Context, *types.TradingCreateOrderRequest) (*types.Order, error)
	TradingGetOrder(context.Context, *cm.IDRequest) (*types.Order, error)
	TradingGetOrders(context.Context, *GetOrdersRequest) (*types.OrdersResponse, error)
}

// +apix:path=/shop.Payment
type PaymentService interface {
	PaymentTradingOrder(context.Context, *PaymentTradingOrderRequest) (*PaymentTradingOrderResponse, error)
	PaymentCheckReturnData(context.Context, *PaymentCheckReturnDataRequest) (*cm.MessageResponse, error)
}

// +apix:path=/shop.Receipt
type ReceiptService interface {
	CreateReceipt(context.Context, *CreateReceiptRequest) (*Receipt, error)
	UpdateReceipt(context.Context, *UpdateReceiptRequest) (*Receipt, error)
	GetReceipt(context.Context, *cm.IDRequest) (*Receipt, error)
	GetReceipts(context.Context, *GetReceiptsRequest) (*ReceiptsResponse, error)
	GetReceiptsByLedgerType(context.Context, *GetReceiptsByLedgerTypeRequest) (*ReceiptsResponse, error)
	ConfirmReceipt(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	CancelReceipt(context.Context, *CancelReceiptRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Supplier
type SupplierService interface {
	GetSupplier(context.Context, *cm.IDRequest) (*Supplier, error)
	GetSuppliers(context.Context, *GetSuppliersRequest) (*SuppliersResponse, error)
	GetSuppliersByIDs(context.Context, *cm.IDsRequest) (*SuppliersResponse, error)
	CreateSupplier(context.Context, *CreateSupplierRequest) (*Supplier, error)
	UpdateSupplier(context.Context, *UpdateSupplierRequest) (*Supplier, error)
	DeleteSupplier(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	GetSuppliersByVariantID(context.Context, *GetSuppliersByVariantIDRequest) (*SuppliersResponse, error)
}

// +apix:path=/shop.Carrier
type CarrierService interface {
	GetCarrier(context.Context, *cm.IDRequest) (*Carrier, error)
	GetCarriers(context.Context, *GetCarriersRequest) (*CarriersResponse, error)
	GetCarriersByIDs(context.Context, *cm.IDsRequest) (*CarriersResponse, error)
	CreateCarrier(context.Context, *CreateCarrierRequest) (*Carrier, error)
	UpdateCarrier(context.Context, *UpdateCarrierRequest) (*Carrier, error)
	DeleteCarrier(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.Ledger
type LedgerService interface {
	GetLedger(context.Context, *cm.IDRequest) (*Ledger, error)
	GetLedgers(context.Context, *GetLedgersRequest) (*LedgersResponse, error)
	CreateLedger(context.Context, *CreateLedgerRequest) (*Ledger, error)
	UpdateLedger(context.Context, *UpdateLedgerRequest) (*Ledger, error)
	DeleteLedger(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.PurchaseOrder
type PurchaseOrderService interface {
	GetPurchaseOrder(context.Context, *cm.IDRequest) (*PurchaseOrder, error)
	GetPurchaseOrders(context.Context, *GetPurchaseOrdersRequest) (*PurchaseOrdersResponse, error)
	GetPurchaseOrdersByIDs(context.Context, *cm.IDsRequest) (*PurchaseOrdersResponse, error)
	GetPurchaseOrdersByReceiptID(context.Context, *cm.IDRequest) (*PurchaseOrdersResponse, error)
	CreatePurchaseOrder(context.Context, *CreatePurchaseOrderRequest) (*PurchaseOrder, error)
	UpdatePurchaseOrder(context.Context, *UpdatePurchaseOrderRequest) (*PurchaseOrder, error)
	DeletePurchaseOrder(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	ConfirmPurchaseOrder(context.Context, *ConfirmPurchaseOrderRequest) (*cm.UpdatedResponse, error)
	CancelPurchaseOrder(context.Context, *CancelPurchaseOrderRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Stocktake
type StocktakeService interface {
	CreateStocktake(context.Context, *CreateStocktakeRequest) (*Stocktake, error)
	UpdateStocktake(context.Context, *UpdateStocktakeRequest) (*Stocktake, error)
	ConfirmStocktake(context.Context, *ConfirmStocktakeRequest) (*Stocktake, error)
	CancelStocktake(context.Context, *CancelStocktakeRequest) (*Stocktake, error)

	GetStocktake(context.Context, *cm.IDRequest) (*Stocktake, error)
	GetStocktakesByIDs(context.Context, *cm.IDsRequest) (*GetStocktakesByIDsResponse, error)
	GetStocktakes(context.Context, *GetStocktakesRequest) (*GetStocktakesResponse, error)
}

// +apix:path=/shop.Connection
type ConnectionService interface {
	GetConnections(context.Context, *cm.Empty) (*types.GetConnectionsResponse, error)
	GetAvailableConnections(context.Context, *cm.Empty) (*types.GetConnectionsResponse, error)
	GetShopConnections(context.Context, *cm.Empty) (*types.GetShopConnectionsResponse, error)
	RegisterShopConnection(context.Context, *types.RegisterShopConnectionRequest) (*types.ShopConnection, error)
	LoginShopConnection(context.Context, *types.LoginShopConnectionRequest) (*types.LoginShopConnectionResponse, error)
	LoginShopConnectionWithOTP(context.Context, *types.LoginShopConnectionWithOTPRequest) (*types.LoginShopConnectionWithOTPResponse, error)
	DeleteShopConnection(context.Context, *types.DeleteShopConnectionRequest) (*cm.DeletedResponse, error)
	UpdateShopConnection(context.Context, *types.UpdateShopConnectionRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Refund
type RefundService interface {
	CreateRefund(context.Context, *CreateRefundRequest) (*Refund, error)
	UpdateRefund(context.Context, *UpdateRefundRequest) (*Refund, error)
	ConfirmRefund(context.Context, *ConfirmRefundRequest) (*Refund, error)
	CancelRefund(context.Context, *CancelRefundRequest) (*Refund, error)

	GetRefund(context.Context, *cm.IDRequest) (*Refund, error)
	GetRefundsByIDs(context.Context, *cm.IDsRequest) (*GetRefundsByIDsResponse, error)
	GetRefunds(context.Context, *GetRefundsRequest) (*GetRefundsResponse, error)
}

// +apix:path=/shop.PurchaseRefund
type PurchaseRefundService interface {
	CreatePurchaseRefund(context.Context, *CreatePurchaseRefundRequest) (*PurchaseRefund, error)
	UpdatePurchaseRefund(context.Context, *UpdatePurchaseRefundRequest) (*PurchaseRefund, error)
	ConfirmPurchaseRefund(context.Context, *ConfirmPurchaseRefundRequest) (*PurchaseRefund, error)
	CancelPurchaseRefund(context.Context, *CancelPurchaseRefundRequest) (*PurchaseRefund, error)

	GetPurchaseRefund(context.Context, *cm.IDRequest) (*PurchaseRefund, error)
	GetPurchaseRefundsByIDs(context.Context, *cm.IDsRequest) (*GetPurchaseRefundsByIDsResponse, error)
	GetPurchaseRefunds(context.Context, *GetPurchaseRefundsRequest) (*GetPurchaseRefundsResponse, error)
}

// +apix:path=/shop.WebServer
type WebServerService interface {

	// ws_website

	CreateWsWebsite(context.Context, *CreateWsWebsiteRequest) (*WsWebsite, error)
	UpdateWsWebsite(context.Context, *UpdateWsWebsiteRequest) (*WsWebsite, error)

	GetWsWebsite(context.Context, *GetWsWebsiteRequest) (*WsWebsite, error)
	GetWsWebsites(context.Context, *GetWsWebsitesRequest) (*GetWsWebsitesResponse, error)
	GetWsWebsitesByIDs(context.Context, *GetWsWebsitesByIDsRequest) (*GetWsWebsitesByIDsResponse, error)

	// ws_product

	CreateOrUpdateWsProduct(context.Context, *CreateOrUpdateWsProductRequest) (*WsProduct, error)

	GetWsProduct(context.Context, *GetWsProductRequest) (*WsProduct, error)
	GetWsProducts(context.Context, *GetWsProductsRequest) (*GetWsProductsResponse, error)
	GetWsProductsByIDs(context.Context, *GetWsProductsByIDsRequest) (*GetWsProductsByIDsResponse, error)

	// ws_category

	CreateOrUpdateWsCategory(context.Context, *CreateOrUpdateWsCategoryRequest) (*WsCategory, error)

	GetWsCategory(context.Context, *GetWsCategoryRequest) (*WsCategory, error)
	GetWsCategories(context.Context, *GetWsCategoriesRequest) (*GetWsCategoriesResponse, error)
	GetWsCategoriesByIDs(context.Context, *GetWsCategoriesByIDsRequest) (*GetWsCategoriesByIDsResponse, error)

	// ws_page

	CreateWsPage(context.Context, *CreateWsPageRequest) (*WsPage, error)
	UpdateWsPage(context.Context, *UpdateWsPageRequest) (*WsPage, error)
	DeleteWsPage(context.Context, *DeteleWsPageRequest) (*DeteleWsPageResponse, error)

	GetWsPage(context.Context, *GetWsPageRequest) (*WsPage, error)
	GetWsPages(context.Context, *GetWsPagesRequest) (*GetWsPagesResponse, error)
	GetWsPagesByIDs(context.Context, *GetWsPagesByIDsRequest) (*GetWsPagesByIDsResponse, error)
}

// +apix:path=/shop.Subscription
type SubscriptionService interface {
	GetSubscription(context.Context, *types.SubscriptionIDRequest) (*types.Subscription, error)
	GetSubscriptions(context.Context, *types.GetSubscriptionsRequest) (*types.GetSubscriptionsResponse, error)
}

// +apix:path=/shop.Ticket
type TicketService interface {
	CreateTicket(context.Context, *CreateTicketRequest) (*shoptypes.Ticket, error)
	GetTickets(context.Context, *GetTicketsRequest) (*GetTicketsResponse, error)
	GetTicket(context.Context, *GetTicketRequest) (*shoptypes.Ticket, error)

	CreateTicketComment(context.Context, *CreateTicketCommentRequest) (*shoptypes.TicketComment, error)
	UpdateTicketComment(context.Context, *UpdateTicketCommentRequest) (*shoptypes.TicketComment, error)
	GetTicketComments(context.Context, *GetTicketCommentsRequest) (*GetTicketCommentsResponse, error)
}
