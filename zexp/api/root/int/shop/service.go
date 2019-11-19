package shop

import (
	"context"

	cm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	"etop.vn/backend/pb/etop/order"
	shop "etop.vn/backend/pb/etop/shop"
)

// +gen:apix

// +apix:path=/shop.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/shop.Brand
type BrandAPI interface {
	CreateBrand(context.Context, *shop.CreateBrandRequest) (*shop.Brand, error)
	UpdateBrandInfo(context.Context, *shop.UpdateBrandRequest) (*shop.Brand, error)
	DeleteBrand(context.Context, *cm.IDsRequest) (*shop.DeleteBrandResponse, error)

	GetBrandByID(context.Context, *cm.IDRequest) (*shop.Brand, error)
	GetBrandsByIDs(context.Context, *cm.IDsRequest) (*shop.GetBrandsByIDsResponse, error)
	GetBrands(context.Context, *shop.GetBrandsRequest) (*shop.GetBrandsResponse, error)
}

// +apix:path=/shop.Inventory
type InventoryAPI interface {
	CreateInventoryVoucher(context.Context, *shop.CreateInventoryVoucherRequest) (*shop.CreateInventoryVoucherResponse, error)
	ConfirmInventoryVoucher(context.Context, *shop.ConfirmInventoryVoucherRequest) (*shop.ConfirmInventoryVoucherResponse, error)
	CancelInventoryVoucher(context.Context, *shop.CancelInventoryVoucherRequest) (*shop.CancelInventoryVoucherResponse, error)
	UpdateInventoryVoucher(context.Context, *shop.UpdateInventoryVoucherRequest) (*shop.UpdateInventoryVoucherResponse, error)
	AdjustInventoryQuantity(context.Context, *shop.AdjustInventoryQuantityRequest) (*shop.AdjustInventoryQuantityResponse, error)

	GetInventoryVariant(context.Context, *shop.GetInventoryVariantRequest) (*shop.InventoryVariant, error)
	GetInventoryVariants(context.Context, *shop.GetInventoryVariantsRequest) (*shop.GetInventoryVariantsResponse, error)
	GetInventoryVariantsByVariantIDs(context.Context, *shop.GetInventoryVariantsByVariantIDsRequest) (*shop.GetInventoryVariantsResponse, error)
	GetInventoryVouchersByReference(context.Context, *shop.GetInventoryVouchersByReferenceRequest) (*shop.GetInventoryVouchersByReferenceResponse, error)
	GetInventoryVoucher(context.Context, *cm.IDRequest) (*shop.InventoryVoucher, error)
	GetInventoryVouchers(context.Context, *shop.GetInventoryVouchersRequest) (*shop.GetInventoryVouchersResponse, error)
	GetInventoryVouchersByIDs(context.Context, *shop.GetInventoryVouchersByIDsRequest) (*shop.GetInventoryVouchersResponse, error)
}

// +apix:path=/shop.Account
type AccountAPI interface {
	RegisterShop(context.Context, *shop.RegisterShopRequest) (*shop.RegisterShopResponse, error)
	UpdateShop(context.Context, *shop.UpdateShopRequest) (*shop.UpdateShopResponse, error)
	DeleteShop(context.Context, *cm.IDRequest) (*cm.Empty, error)
	SetDefaultAddress(context.Context, *etop.SetDefaultAddressRequest) (*cm.UpdatedResponse, error)
	GetBalanceShop(context.Context, *cm.Empty) (*shop.GetBalanceShopResponse, error)

	CreateExternalAccountAhamove(context.Context, *cm.Empty) (*shop.ExternalAccountAhamove, error)
	GetExternalAccountAhamove(context.Context, *cm.Empty) (*shop.ExternalAccountAhamove, error)
	RequestVerifyExternalAccountAhamove(context.Context, *cm.Empty) (*cm.UpdatedResponse, error)
	UpdateExternalAccountAhamoveVerification(context.Context, *shop.UpdateXAccountAhamoveVerificationRequest) (*cm.UpdatedResponse, error)

	// deprecated: backward-compatible, will be removed later
	UpdateExternalAccountAhamoveVerificationImages(context.Context, *shop.UpdateXAccountAhamoveVerificationRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.ExternalAccount
type ExternalAccountAPI interface {
	GetExternalAccountHaravan(context.Context, *cm.Empty) (*shop.ExternalAccountHaravan, error)
	CreateExternalAccountHaravan(context.Context, *shop.ExternalAccountHaravanRequest) (*shop.ExternalAccountHaravan, error)
	UpdateExternalAccountHaravanToken(context.Context, *shop.ExternalAccountHaravanRequest) (*shop.ExternalAccountHaravan, error)
	ConnectCarrierServiceExternalAccountHaravan(context.Context, *cm.Empty) (*cm.UpdatedResponse, error)
	DeleteConnectedCarrierServiceExternalAccountHaravan(context.Context, *cm.Empty) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.Collection
type CollectionAPI interface {
	CreateCollection(context.Context, *shop.CreateCollectionRequest) (*shop.ShopCollection, error)
	GetCollection(context.Context, *cm.IDRequest) (*shop.ShopCollection, error)
	GetCollections(context.Context, *shop.GetCollectionsRequest) (*shop.ShopCollectionsResponse, error)
	UpdateCollection(context.Context, *shop.UpdateCollectionRequest) (*shop.ShopCollection, error)
	GetCollectionsByProductID(context.Context, *shop.GetShopCollectionsByProductIDRequest) (*shop.CollectionsResponse, error)
}

// +apix:path=/shop.Customer
type CustomerAPI interface {
	CreateCustomer(context.Context, *shop.CreateCustomerRequest) (*shop.Customer, error)
	UpdateCustomer(context.Context, *shop.UpdateCustomerRequest) (*shop.Customer, error)
	DeleteCustomer(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	GetCustomer(context.Context, *cm.IDRequest) (*shop.Customer, error)
	GetCustomerDetails(context.Context, *cm.IDRequest) (*shop.CustomerDetailsResponse, error)
	GetCustomers(context.Context, *shop.GetCustomersRequest) (*shop.CustomersResponse, error)
	GetCustomersByIDs(context.Context, *cm.IDsRequest) (*shop.CustomersResponse, error)
	BatchSetCustomersStatus(context.Context, *shop.SetCustomersStatusRequest) (*cm.UpdatedResponse, error)

	//-- address --//

	GetCustomerAddresses(context.Context, *shop.GetCustomerAddressesRequest) (*shop.CustomerAddressesResponse, error)
	CreateCustomerAddress(context.Context, *shop.CreateCustomerAddressRequest) (*shop.CustomerAddress, error)
	UpdateCustomerAddress(context.Context, *shop.UpdateCustomerAddressRequest) (*shop.CustomerAddress, error)
	SetDefaultCustomerAddress(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	DeleteCustomerAddress(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)

	//-- group --//
	AddCustomersToGroup(context.Context, *shop.AddCustomerToGroupRequest) (*cm.UpdatedResponse, error)
	RemoveCustomersFromGroup(context.Context, *shop.RemoveCustomerOutOfGroupRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/shop.CustomerGroup
type CustomerGroupAPI interface {
	CreateCustomerGroup(context.Context, *shop.CreateCustomerGroupRequest) (*shop.CustomerGroup, error)
	GetCustomerGroup(context.Context, *cm.IDRequest) (*shop.CustomerGroup, error)
	GetCustomerGroups(context.Context, *shop.GetCustomerGroupsRequest) (*shop.CustomerGroupsResponse, error)
	UpdateCustomerGroup(context.Context, *shop.UpdateCustomerGroupRequest) (*shop.CustomerGroup, error)
}

// +apix:path=/shop.Product
type ProductAPI interface {

	//-- product --//

	GetProduct(context.Context, *cm.IDRequest) (*shop.ShopProduct, error)
	GetProducts(context.Context, *shop.GetVariantsRequest) (*shop.ShopProductsResponse, error)
	GetProductsByIDs(context.Context, *cm.IDsRequest) (*shop.ShopProductsResponse, error)

	CreateProduct(context.Context, *shop.CreateProductRequest) (*shop.ShopProduct, error)
	UpdateProduct(context.Context, *shop.UpdateProductRequest) (*shop.ShopProduct, error)
	UpdateProductsStatus(context.Context, *shop.UpdateProductStatusRequest) (*shop.UpdateProductStatusResponse, error)
	UpdateProductsTags(context.Context, *shop.UpdateProductsTagsRequest) (*cm.UpdatedResponse, error)
	UpdateProductImages(context.Context, *shop.UpdateVariantImagesRequest) (*shop.ShopProduct, error)
	UpdateProductMetaFields(context.Context, *shop.UpdateProductMetaFieldsRequest) (*shop.ShopProduct, error)
	RemoveProducts(context.Context, *shop.RemoveVariantsRequest) (*cm.RemovedResponse, error)
	//-- variant --//

	GetVariant(context.Context, *cm.IDRequest) (*shop.ShopVariant, error)
	GetVariantsByIDs(context.Context, *cm.IDsRequest) (*shop.ShopVariantsResponse, error)

	CreateVariant(context.Context, *shop.CreateVariantRequest) (*shop.ShopVariant, error)
	UpdateVariant(context.Context, *shop.UpdateVariantRequest) (*shop.ShopVariant, error)
	UpdateVariantImages(context.Context, *shop.UpdateVariantImagesRequest) (*shop.ShopVariant, error)
	UpdateVariantsStatus(context.Context, *shop.UpdateProductStatusRequest) (*shop.UpdateProductStatusResponse, error)
	UpdateVariantAttributes(context.Context, *shop.UpdateVariantAttributesRequest) (*shop.ShopVariant, error)
	RemoveVariants(context.Context, *shop.RemoveVariantsRequest) (*cm.RemovedResponse, error)

	//-- category --//
	UpdateProductCategory(context.Context, *shop.UpdateProductCategoryRequest) (*shop.ShopProduct, error)
	RemoveProductCategory(context.Context, *cm.IDRequest) (*shop.ShopProduct, error)

	//-- collection --//
	AddProductCollection(context.Context, *shop.AddShopProductCollectionRequest) (*cm.UpdatedResponse, error)
	RemoveProductCollection(context.Context, *shop.RemoveShopProductCollectionRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/shop.Category
type CategoryAPI interface {
	CreateCategory(context.Context, *shop.CreateCategoryRequest) (*shop.ShopCategory, error)
	GetCategory(context.Context, *cm.IDRequest) (*shop.ShopCategory, error)
	GetCategories(context.Context, *shop.GetCategoriesRequest) (*shop.ShopCategoriesResponse, error)
	UpdateCategory(context.Context, *shop.UpdateCategoryRequest) (*shop.ShopCategory, error)
	DeleteCategory(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.ProductSource
// deprecated: 2018.07.31+14
type ProductSourceAPI interface {
	CreateProductSource(context.Context, *shop.CreateProductSourceRequest) (*shop.ProductSource, error)
	GetShopProductSources(context.Context, *cm.Empty) (*shop.ProductSourcesResponse, error)
	// deprecated: use shop.Product/CreateVariant instead
	CreateVariant(context.Context, *shop.DeprecatedCreateVariantRequest) (*shop.ShopProduct, error)
	CreateProductSourceCategory(context.Context, *shop.CreatePSCategoryRequest) (*shop.Category, error)
	UpdateProductsPSCategory(context.Context, *shop.UpdateProductsPSCategoryRequest) (*cm.UpdatedResponse, error)
	GetProductSourceCategory(context.Context, *cm.IDRequest) (*shop.Category, error)
	GetProductSourceCategories(context.Context, *shop.GetProductSourceCategoriesRequest) (*shop.CategoriesResponse, error)
	UpdateProductSourceCategory(context.Context, *shop.UpdateProductSourceCategoryRequest) (*shop.Category, error)
	RemoveProductSourceCategory(context.Context, *cm.IDRequest) (*cm.RemovedResponse, error)
}

// +apix:path=/shop.Order
type OrderAPI interface {
	CreateOrder(context.Context, *order.CreateOrderRequest) (*order.Order, error)
	GetOrder(context.Context, *cm.IDRequest) (*order.Order, error)
	GetOrders(context.Context, *shop.GetOrdersRequest) (*order.OrdersResponse, error)
	GetOrdersByIDs(context.Context, *etop.IDsRequest) (*order.OrdersResponse, error)
	GetOrdersByReceiptID(context.Context, *shop.GetOrdersByReceiptIDRequest) (*order.OrdersResponse, error)
	UpdateOrder(context.Context, *order.UpdateOrderRequest) (*order.Order, error)

	// @deprecated
	UpdateOrdersStatus(context.Context, *shop.UpdateOrdersStatusRequest) (*cm.UpdatedResponse, error)

	ConfirmOrder(context.Context, *shop.ConfirmOrderRequest) (*order.Order, error)
	ConfirmOrderAndCreateFulfillments(context.Context, *shop.OrderIDRequest) (*order.OrderWithErrorsResponse, error)
	CancelOrder(context.Context, *shop.CancelOrderRequest) (*order.OrderWithErrorsResponse, error)
	UpdateOrderPaymentStatus(context.Context, *shop.UpdateOrderPaymentStatusRequest) (*cm.UpdatedResponse, error)
	UpdateOrderShippingInfo(context.Context, *shop.UpdateOrderShippingInfoRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Fulfillment
type FulfillmentAPI interface {
	GetFulfillment(context.Context, *cm.IDRequest) (*order.Fulfillment, error)
	GetFulfillments(context.Context, *shop.GetFulfillmentsRequest) (*order.FulfillmentsResponse, error)

	GetPublicExternalShippingServices(context.Context, *order.GetExternalShippingServicesRequest) (*order.GetExternalShippingServicesResponse, error)
	GetExternalShippingServices(context.Context, *order.GetExternalShippingServicesRequest) (*order.GetExternalShippingServicesResponse, error)
	GetPublicFulfillment(context.Context, *shop.GetPublicFulfillmentRequest) (*order.PublicFulfillment, error)
	UpdateFulfillmentsShippingState(context.Context, *shop.UpdateFulfillmentsShippingStateRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/etop.Shipnow
type ShipnowAPI interface {
	GetShipnowFulfillment(context.Context, *cm.IDRequest) (*order.ShipnowFulfillment, error)
	GetShipnowFulfillments(context.Context, *order.GetShipnowFulfillmentsRequest) (*order.ShipnowFulfillments, error)

	CreateShipnowFulfillment(context.Context, *order.CreateShipnowFulfillmentRequest) (*order.ShipnowFulfillment, error)
	ConfirmShipnowFulfillment(context.Context, *cm.IDRequest) (*order.ShipnowFulfillment, error)
	UpdateShipnowFulfillment(context.Context, *order.UpdateShipnowFulfillmentRequest) (*order.ShipnowFulfillment, error)
	CancelShipnowFulfillment(context.Context, *order.CancelShipnowFulfillmentRequest) (*cm.UpdatedResponse, error)

	GetShipnowServices(context.Context, *order.GetShipnowServicesRequest) (*order.GetShipnowServicesResponse, error)
}

// +apix:path=/shop.History
type HistoryAPI interface {
	GetFulfillmentHistory(context.Context, *shop.GetFulfillmentHistoryRequest) (*etop.HistoryResponse, error)
}

// +apix:path=/shop.MoneyTransaction
type MoneyTransactionAPI interface {
	GetMoneyTransaction(context.Context, *cm.IDRequest) (*order.MoneyTransaction, error)
	GetMoneyTransactions(context.Context, *shop.GetMoneyTransactionsRequest) (*order.MoneyTransactionsResponse, error)
}

// +apix:path=/shop.Summary
type SummaryAPI interface {
	SummarizeFulfillments(context.Context, *shop.SummarizeFulfillmentsRequest) (*shop.SummarizeFulfillmentsResponse, error)
	SummarizePOS(context.Context, *shop.SummarizePOSRequest) (*shop.SummarizePOSResponse, error)
	CalcBalanceShop(context.Context, *cm.Empty) (*shop.CalcBalanceShopResponse, error)
}

// +apix:path=/shop.Export
type ExportAPI interface {
	GetExports(context.Context, *shop.GetExportsRequest) (*shop.GetExportsResponse, error)
	RequestExport(context.Context, *shop.RequestExportRequest) (*shop.RequestExportResponse, error)
}

// +apix:path=/shop.Notification
type NotificationAPI interface {
	CreateDevice(context.Context, *etop.CreateDeviceRequest) (*etop.Device, error)
	DeleteDevice(context.Context, *etop.DeleteDeviceRequest) (*cm.DeletedResponse, error)

	GetNotification(context.Context, *cm.IDRequest) (*etop.Notification, error)
	GetNotifications(context.Context, *etop.GetNotificationsRequest) (*etop.NotificationsResponse, error)
	UpdateNotifications(context.Context, *etop.UpdateNotificationsRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Authorize
type AuthorizeAPI interface {
	GetAuthorizedPartners(context.Context, *cm.Empty) (*shop.GetAuthorizedPartnersResponse, error)
	GetAvailablePartners(context.Context, *cm.Empty) (*shop.GetPartnersResponse, error)
	AuthorizePartner(context.Context, *shop.AuthorizePartnerRequest) (*shop.AuthorizedPartnerResponse, error)
}

// +apix:path=/shop.Trading
type TradingAPI interface {
	TradingGetProduct(context.Context, *cm.IDRequest) (*shop.ShopProduct, error)
	TradingGetProducts(context.Context, *cm.CommonListRequest) (*shop.ShopProductsResponse, error)

	TradingCreateOrder(context.Context, *order.TradingCreateOrderRequest) (*order.Order, error)
	TradingGetOrder(context.Context, *cm.IDRequest) (*order.Order, error)
	TradingGetOrders(context.Context, *shop.GetOrdersRequest) (*order.OrdersResponse, error)
}

// +apix:path=/shop.Payment
type PaymentAPI interface {
	PaymentTradingOrder(context.Context, *shop.PaymentTradingOrderRequest) (*shop.PaymentTradingOrderResponse, error)
	PaymentCheckReturnData(context.Context, *shop.PaymentCheckReturnDataRequest) (*cm.MessageResponse, error)
}

// +apix:path=/shop.Receipt
type ReceiptAPI interface {
	CreateReceipt(context.Context, *shop.CreateReceiptRequest) (*shop.Receipt, error)
	UpdateReceipt(context.Context, *shop.UpdateReceiptRequest) (*shop.Receipt, error)
	GetReceipt(context.Context, *cm.IDRequest) (*shop.Receipt, error)
	GetReceipts(context.Context, *shop.GetReceiptsRequest) (*shop.ReceiptsResponse, error)
	GetReceiptsByLedgerType(context.Context, *shop.GetReceiptsByLedgerTypeRequest) (*shop.ReceiptsResponse, error)
	ConfirmReceipt(context.Context, *cm.IDRequest) (*cm.UpdatedResponse, error)
	CancelReceipt(context.Context, *shop.CancelReceiptRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/shop.Supplier
type SupplierAPI interface {
	GetSupplier(context.Context, *cm.IDRequest) (*shop.Supplier, error)
	GetSuppliers(context.Context, *shop.GetSuppliersRequest) (*shop.SuppliersResponse, error)
	GetSuppliersByIDs(context.Context, *cm.IDsRequest) (*shop.SuppliersResponse, error)
	CreateSupplier(context.Context, *shop.CreateSupplierRequest) (*shop.Supplier, error)
	UpdateSupplier(context.Context, *shop.UpdateSupplierRequest) (*shop.Supplier, error)
	DeleteSupplier(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.Carrier
type CarrierAPI interface {
	GetCarrier(context.Context, *cm.IDRequest) (*shop.Carrier, error)
	GetCarriers(context.Context, *shop.GetCarriersRequest) (*shop.CarriersResponse, error)
	GetCarriersByIDs(context.Context, *cm.IDsRequest) (*shop.CarriersResponse, error)
	CreateCarrier(context.Context, *shop.CreateCarrierRequest) (*shop.Carrier, error)
	UpdateCarrier(context.Context, *shop.UpdateCarrierRequest) (*shop.Carrier, error)
	DeleteCarrier(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.Ledger
type LedgerAPI interface {
	GetLedger(context.Context, *cm.IDRequest) (*shop.Ledger, error)
	GetLedgers(context.Context, *shop.GetLedgersRequest) (*shop.LedgersResponse, error)
	CreateLedger(context.Context, *shop.CreateLedgerRequest) (*shop.Ledger, error)
	UpdateLedger(context.Context, *shop.UpdateLedgerRequest) (*shop.Ledger, error)
	DeleteLedger(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/shop.PurchaseOrder
type PurchaseOrderAPI interface {
	GetPurchaseOrder(context.Context, *cm.IDRequest) (*shop.PurchaseOrder, error)
	GetPurchaseOrders(context.Context, *shop.GetPurchaseOrdersRequest) (*shop.PurchaseOrdersResponse, error)
	GetPurchaseOrdersByIDs(context.Context, *cm.IDsRequest) (*shop.PurchaseOrdersResponse, error)
	GetPurchaseOrdersByReceiptID(context.Context, *cm.IDRequest) (*shop.PurchaseOrdersResponse, error)
	CreatePurchaseOrder(context.Context, *shop.CreatePurchaseOrderRequest) (*shop.PurchaseOrder, error)
	UpdatePurchaseOrder(context.Context, *shop.UpdatePurchaseOrderRequest) (*shop.PurchaseOrder, error)
	DeletePurchaseOrder(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
	ConfirmPurchaseOrder(context.Context, *shop.ConfirmPurchaseOrderRequest) (*cm.UpdatedResponse, error)
	CancelPurchaseOrder(context.Context, *shop.CancelPurchaseOrderRequest) (*cm.UpdatedResponse, error)
}
