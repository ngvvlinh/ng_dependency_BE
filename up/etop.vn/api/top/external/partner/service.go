package partner_proto

import (
	"context"

	"etop.vn/api/top/external/types"
	etop "etop.vn/api/top/int/etop"
	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/partner

// +apix:path=/partner.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*types.Partner, error)
	GetLocationList(context.Context, *cm.Empty) (*types.LocationResponse, error)
}

// +apix:path=/partner.Shop
type ShopService interface {
	AuthorizeShop(context.Context, *AuthorizeShopRequest) (*AuthorizeShopResponse, error)
	CurrentShop(context.Context, *cm.Empty) (*etop.PublicAccountInfo, error)
}

// +apix:path=/partner.Webhook
type WebhookService interface {
	CreateWebhook(context.Context, *types.CreateWebhookRequest) (*types.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*types.WebhooksResponse, error)
	DeleteWebhook(context.Context, *types.DeleteWebhookRequest) (*types.WebhooksResponse, error)
}

// +apix:path=/partner.History
type HistoryService interface {
	GetChanges(context.Context, *types.GetChangesRequest) (*types.Callback, error)
}

// +apix:path=/partner.Shipping
type ShippingService interface {
	GetShippingServices(context.Context, *types.GetShippingServicesRequest) (*types.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *types.CreateOrderRequest) (*types.OrderAndFulfillments, error)
	CancelOrder(context.Context, *types.CancelOrderRequest) (*types.OrderAndFulfillments, error)
	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)
}

// +apix:path=/partner.Order
// +wrapper:endpoint-prefix=Order
type OrderService interface {
	CreateOrder(context.Context, *types.CreateOrderRequest) (*types.Order, error)
	// TODO:
	ConfirmOrder(context.Context, *types.ConfirmOrderRequest) (*cm.Empty, error)
	CancelOrder(context.Context, *types.CancelOrderRequest) (*types.OrderAndFulfillments, error)
	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)
}

// +apix:path=/partner.Fulfillment
// +wrapper:endpoint-prefix=Fulfillment
type FulfillmentService interface {
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)
	ListFulfillments(context.Context, *types.ListFulfillmentsRequest) (*types.FulfillmentsResponse, error)
	// + CreateFulfillment
	// + Confirm
}

// +apix:path=/partner.Customer
type CustomerService interface {
	GetCustomer(context.Context, *types.GetCustomerRequest) (*types.Customer, error)
	ListCustomers(context.Context, *types.ListCustomersRequest) (*types.CustomersResponse, error)
	CreateCustomer(context.Context, *types.CreateCustomerRequest) (*types.Customer, error)
	UpdateCustomer(context.Context, *types.UpdateCustomerRequest) (*types.Customer, error)
	DeleteCustomer(context.Context, *types.DeleteCustomerRequest) (*cm.Empty, error)

	//-- address --//
	ListAddresses(context.Context, *types.ListCustomerAddressesRequest) (*types.CustomerAddressesResponse, error)
	CreateAddress(context.Context, *types.CreateCustomerAddressRequest) (*types.CustomerAddress, error)
	UpdateAddress(context.Context, *types.UpdateCustomerAddressRequest) (*types.CustomerAddress, error)
	DeleteAddress(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/partner.CustomerGroup
type CustomerGroupService interface {
	GetCustomerGroup(context.Context, *cm.IDRequest) (*types.CustomerGroup, error)
	ListCustomerGroups(context.Context, *types.ListCustomerGroupsRequest) (*types.CustomerGroupsResponse, error)
	CreateCustomerGroup(context.Context, *types.CreateCustomerGroupRequest) (*types.CustomerGroup, error)
	UpdateCustomerGroup(context.Context, *types.UpdateCustomerGroupRequest) (*types.CustomerGroup, error)

	AddCustomers(context.Context, *types.AddCustomersRequest) (*cm.Empty, error)
	RemoveCustomers(context.Context, *types.RemoveCustomersRequest) (*cm.Empty, error)
}

// +apix:path=/partner.Inventory
type InventoryService interface {
	ListInventoryLevels(context.Context, *types.ListInventoryLevelsRequest) (*types.InventoryLevelsResponse, error)
}

// +apix:path=/partner.Product
type ProductService interface {
	GetProduct(context.Context, *types.GetProductRequest) (*types.ShopProduct, error)
	ListProducts(context.Context, *types.ListProductsRequest) (*types.ShopProductsResponse, error)
	CreateProduct(context.Context, *types.CreateProductRequest) (*types.ShopProduct, error)
	UpdateProduct(context.Context, *types.UpdateProductRequest) (*types.ShopProduct, error)
	// TODO
	BatchUpdateProducts(context.Context, *types.BatchUpdateProductsRequest) (*types.ShopProductsResponse, error)
	AddProductCollection(context.Context, *types.AddProductCollectionRequest) (*cm.Empty, error)
	RemoveProductCollection(context.Context, *types.RemoveProductCollectionRequest) (*cm.Empty, error)
}

// +apix:path=/partner.Variant
type VariantService interface {
	GetVariant(context.Context, *types.GetVariantRequest) (*types.ShopVariant, error)
	ListVariants(context.Context, *types.ListVariantsRequest) (*types.ShopVariantsResponse, error)
	CreateVariant(context.Context, *types.CreateVariantRequest) (*types.ShopVariant, error)
	UpdateVariant(context.Context, *types.UpdateVariantRequest) (*types.ShopVariant, error)
	// TODO
	BatchUpdateVariants(context.Context, *types.BatchUpdateVariantsRequest) (*types.ShopVariantsResponse, error)
}
