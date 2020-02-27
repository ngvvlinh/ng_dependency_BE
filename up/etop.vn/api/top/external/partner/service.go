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

// +apix:path=/carrier.ShipmentConnection
type ShipmentConnectionService interface {
	GetConnections(context.Context, *cm.Empty) (*GetConnectionsResponse, error)
	CreateConnection(context.Context, *CreateConnectionRequest) (*ShipmentConnection, error)
	UpdateConnection(context.Context, *UpdateConnectionRequest) (*ShipmentConnection, error)
	DeleteConnection(context.Context, *cm.IDRequest) (*cm.DeletedResponse, error)
}

// +apix:path=/carrier.Shipment
type ShipmentService interface {
	UpdateFulfillment(context.Context, *UpdateFulfillmentRequest) (*cm.UpdatedResponse, error)
}

// +apix:path=/partner.Webhook
type WebhookService interface {
	CreateWebhook(context.Context, *types.CreateWebhookRequest) (*types.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*types.WebhooksResponse, error)
	DeleteWebhook(context.Context, *types.DeleteWebhookRequest) (*types.WebhooksResponse, error)
}

// +apix:path=/partner.History
type HistoryService interface {
	// This API provides an example for webhook request content. It's not a real API.
	GetChanges(context.Context, *cm.Empty) (*types.Callback, error)
}

// +apix:path=/partner.Shipping
type ShippingService interface {
	GetShippingServices(context.Context, *types.GetShippingServicesRequest) (*types.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *types.CreateAndConfirmOrderRequest) (*types.OrderAndFulfillments, error)
	CancelOrder(context.Context, *types.CancelOrderRequest) (*types.OrderAndFulfillments, error)
	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)
}

// +apix:path=/partner.Order
// +wrapper:endpoint-prefix=Order
type OrderService interface {
	CreateOrder(context.Context, *types.CreateOrderRequest) (*types.OrderWithoutShipping, error)

	ConfirmOrder(context.Context, *types.ConfirmOrderRequest) (*cm.Empty, error)

	CancelOrder(context.Context, *types.CancelOrderRequest) (*cm.Empty, error)

	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)

	ListOrders(context.Context, *types.ListOrdersRequest) (*types.OrdersResponse, error)
}

// +apix:path=/partner.Fulfillment
// +wrapper:endpoint-prefix=Fulfillment
type FulfillmentService interface {
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)

	ListFulfillments(context.Context, *types.ListFulfillmentsRequest) (*types.FulfillmentsResponse, error)

	CreateFulfillment(context.Context, *types.CreateFulfillmentRequest) (*types.Fulfillment, error)

	CancelFulfillment(context.Context, *types.CancelFulfillmentRequest) (*cm.Empty, error)
}

// +apix:path=/partner.Customer
type CustomerService interface {
	GetCustomer(context.Context, *types.GetCustomerRequest) (*types.Customer, error)

	ListCustomers(context.Context, *types.ListCustomersRequest) (*types.CustomersResponse, error)

	CreateCustomer(context.Context, *types.CreateCustomerRequest) (*types.Customer, error)

	UpdateCustomer(context.Context, *types.UpdateCustomerRequest) (*types.Customer, error)

	DeleteCustomer(context.Context, *types.DeleteCustomerRequest) (*cm.Empty, error)
}

// +apix:path=/partner.CustomerAddress
type CustomerAddressService interface {
	GetAddress(context.Context, *types.OrderIDRequest) (*types.CustomerAddress, error)

	ListAddresses(context.Context, *types.ListCustomerAddressesRequest) (*types.CustomerAddressesResponse, error)

	CreateAddress(context.Context, *types.CreateCustomerAddressRequest) (*types.CustomerAddress, error)

	UpdateAddress(context.Context, *types.UpdateCustomerAddressRequest) (*types.CustomerAddress, error)

	DeleteAddress(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/partner.CustomerGroup
type CustomerGroupService interface {
	GetGroup(context.Context, *cm.IDRequest) (*types.CustomerGroup, error)

	ListGroups(context.Context, *types.ListCustomerGroupsRequest) (*types.CustomerGroupsResponse, error)

	CreateGroup(context.Context, *types.CreateCustomerGroupRequest) (*types.CustomerGroup, error)

	UpdateGroup(context.Context, *types.UpdateCustomerGroupRequest) (*types.CustomerGroup, error)

	DeleteGroup(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/partner.CustomerGroupRelationship
// +wrapper:endpoint-prefix=CustomerGroup
type CustomerGroupRelationshipService interface {
	ListRelationships(context.Context, *types.ListCustomerGroupRelationshipsRequest) (*types.CustomerGroupRelationshipsResponse, error)

	CreateRelationship(context.Context, *types.AddCustomerRequest) (*cm.Empty, error)

	DeleteRelationship(context.Context, *types.RemoveCustomerRequest) (*cm.Empty, error)
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

	DeleteProduct(context.Context, *types.GetProductRequest) (*cm.Empty, error)
}

// +apix:path=/partner.ProductCollection
type ProductCollectionService interface {
	GetCollection(context.Context, *types.GetCollectionRequest) (*types.ProductCollection, error)

	ListCollections(context.Context, *types.ListCollectionsRequest) (*types.ProductCollectionsResponse, error)

	CreateCollection(context.Context, *types.CreateCollectionRequest) (*types.ProductCollection, error)

	UpdateCollection(context.Context, *types.UpdateCollectionRequest) (*types.ProductCollection, error)

	DeleteCollection(context.Context, *types.GetCollectionRequest) (*cm.Empty, error)
}

// +apix:path=/partner.ProductCollectionRelationship
// +wrapper:endpoint-prefix=ProductCollection
type ProductCollectionRelationshipService interface {
	ListRelationships(context.Context, *types.ListProductCollectionRelationshipsRequest) (*types.ProductCollectionRelationshipsResponse, error)

	CreateRelationship(context.Context, *types.CreateProductCollectionRelationshipRequest) (*cm.Empty, error)

	DeleteRelationship(context.Context, *types.RemoveProductCollectionRequest) (*cm.Empty, error)
}

// +apix:path=/partner.Variant
type VariantService interface {
	GetVariant(context.Context, *types.GetVariantRequest) (*types.ShopVariant, error)

	ListVariants(context.Context, *types.ListVariantsRequest) (*types.ShopVariantsResponse, error)

	CreateVariant(context.Context, *types.CreateVariantRequest) (*types.ShopVariant, error)

	UpdateVariant(context.Context, *types.UpdateVariantRequest) (*types.ShopVariant, error)

	DeleteVariant(context.Context, *types.GetVariantRequest) (*cm.Empty, error)
}
