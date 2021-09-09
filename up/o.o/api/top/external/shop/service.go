package shop

import (
	"context"

	"o.o/api/top/external/types"
	etop "o.o/api/top/int/etop"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/shop

// +apix:path=/shop.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*etop.PublicAccountInfo, error)
	GetLocationList(context.Context, *cm.Empty) (*types.LocationResponse, error)
}

// +apix:path=/shop.Webhook
type WebhookService interface {
	CreateWebhook(context.Context, *types.CreateWebhookRequest) (*types.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*types.WebhooksResponse, error)
	DeleteWebhook(context.Context, *types.DeleteWebhookRequest) (*types.WebhooksResponse, error)
}

// +apix:path=/shop.History
type HistoryService interface {
	// This API provides an example for webhook request content. It's not a real API.
	GetChanges(context.Context, *cm.Empty) (*types.Callback, error)
}

// +apix:path=/shop.Shipping
type ShippingService interface {
	GetShippingServices(context.Context, *types.GetShippingServicesRequest) (*types.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *types.CreateAndConfirmOrderRequest) (*types.OrderAndFulfillments, error)
	CancelOrder(context.Context, *types.CancelOrderRequest) (*types.OrderAndFulfillments, error)
	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)
}

// +apix:path=/shop.Shipnow
type ShipnowService interface {
	GetShipnowServices(context.Context, *types.GetShipnowServicesRequest) (*types.GetShipnowServicesResponse, error)

	CreateShipnowFulfillment(context.Context, *types.CreateShipnowFulfillmentRequest) (*types.ShipnowFulfillment, error)

	CancelShipnowFulfillment(context.Context, *types.CancelShipnowFulfillmentRequest) (*cm.UpdatedResponse, error)

	GetShipnowFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.ShipnowFulfillment, error)
}

// +apix:path=/shop.Order
// +wrapper:endpoint-prefix=Order
type OrderService interface {
	CreateOrder(context.Context, *types.CreateOrderRequest) (*types.Order, error)

	ConfirmOrder(context.Context, *types.ConfirmOrderRequest) (*cm.Empty, error)

	CancelOrder(context.Context, *types.CancelOrderRequest) (*types.OrderAndFulfillments, error)

	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)

	ListOrders(context.Context, *types.ListOrdersRequest) (*types.OrdersResponse, error)
}

// +apix:path=/shop.Fulfillment
// +wrapper:endpoint-prefix=Fulfillment
type FulfillmentService interface {
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)

	ListFulfillments(context.Context, *types.ListFulfillmentsRequest) (*types.FulfillmentsResponse, error)
}

// +apix:path=/shop.Customer
type CustomerService interface {
	GetCustomer(context.Context, *types.GetCustomerRequest) (*types.Customer, error)

	ListCustomers(context.Context, *types.ListCustomersRequest) (*types.CustomersResponse, error)

	CreateCustomer(context.Context, *types.CreateCustomerRequest) (*types.Customer, error)

	UpdateCustomer(context.Context, *types.UpdateCustomerRequest) (*types.Customer, error)

	DeleteCustomer(context.Context, *types.DeleteCustomerRequest) (*cm.Empty, error)
}

// +apix:path=/shop.CustomerAddress
type CustomerAddressService interface {
	GetAddress(context.Context, *types.OrderIDRequest) (*types.CustomerAddress, error)

	ListAddresses(context.Context, *types.ListCustomerAddressesRequest) (*types.CustomerAddressesResponse, error)

	CreateAddress(context.Context, *types.CreateCustomerAddressRequest) (*types.CustomerAddress, error)

	UpdateAddress(context.Context, *types.UpdateCustomerAddressRequest) (*types.CustomerAddress, error)

	DeleteAddress(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/shop.CustomerGroup
type CustomerGroupService interface {
	GetGroup(context.Context, *cm.IDRequest) (*types.CustomerGroup, error)

	ListGroups(context.Context, *types.ListCustomerGroupsRequest) (*types.CustomerGroupsResponse, error)

	CreateGroup(context.Context, *types.CreateCustomerGroupRequest) (*types.CustomerGroup, error)

	UpdateGroup(context.Context, *types.UpdateCustomerGroupRequest) (*types.CustomerGroup, error)

	DeleteGroup(context.Context, *cm.IDRequest) (*cm.Empty, error)
}

// +apix:path=/shop.CustomerGroupRelationship
// +wrapper:endpoint-prefix=CustomerGroup
type CustomerGroupRelationshipService interface {
	ListRelationships(context.Context, *types.ListCustomerGroupRelationshipsRequest) (*types.CustomerGroupRelationshipsResponse, error)

	CreateRelationship(context.Context, *types.AddCustomerRequest) (*cm.Empty, error)

	DeleteRelationship(context.Context, *types.RemoveCustomerRequest) (*cm.Empty, error)
}

// +apix:path=/shop.Inventory
type InventoryService interface {
	ListInventoryLevels(context.Context, *types.ListInventoryLevelsRequest) (*types.InventoryLevelsResponse, error)
}

// +apix:path=/shop.Product
type ProductService interface {
	GetProduct(context.Context, *types.GetProductRequest) (*types.ShopProduct, error)

	ListProducts(context.Context, *types.ListProductsRequest) (*types.ShopProductsResponse, error)

	CreateProduct(context.Context, *types.CreateProductRequest) (*types.ShopProduct, error)

	UpdateProduct(context.Context, *types.UpdateProductRequest) (*types.ShopProduct, error)

	DeleteProduct(context.Context, *types.GetProductRequest) (*cm.Empty, error)
}

// +apix:path=/shop.ProductCollection
type ProductCollectionService interface {
	GetCollection(context.Context, *types.GetCollectionRequest) (*types.ProductCollection, error)

	ListCollections(context.Context, *types.ListCollectionsRequest) (*types.ProductCollectionsResponse, error)

	CreateCollection(context.Context, *types.CreateCollectionRequest) (*types.ProductCollection, error)

	UpdateCollection(context.Context, *types.UpdateCollectionRequest) (*types.ProductCollection, error)

	DeleteCollection(context.Context, *types.GetCollectionRequest) (*cm.Empty, error)
}

// +apix:path=/shop.ProductCollectionRelationship
// +wrapper:endpoint-prefix=ProductCollection
type ProductCollectionRelationshipService interface {
	ListRelationships(context.Context, *types.ListProductCollectionRelationshipsRequest) (*types.ProductCollectionRelationshipsResponse, error)

	CreateRelationship(context.Context, *types.CreateProductCollectionRelationshipRequest) (*cm.Empty, error)

	DeleteRelationship(context.Context, *types.RemoveProductCollectionRequest) (*cm.Empty, error)
}

// +apix:path=/shop.Variant
type VariantService interface {
	GetVariant(context.Context, *types.GetVariantRequest) (*types.ShopVariant, error)

	ListVariants(context.Context, *types.ListVariantsRequest) (*types.ShopVariantsResponse, error)

	CreateVariant(context.Context, *types.CreateVariantRequest) (*types.ShopVariant, error)

	UpdateVariant(context.Context, *types.UpdateVariantRequest) (*types.ShopVariant, error)

	DeleteVariant(context.Context, *types.GetVariantRequest) (*cm.Empty, error)
}

// +apix:path=/shop.Etelecom
type EtelecomService interface {
	GetExtensionInfo(context.Context, *types.GetExtensionInfoRequest) (*types.ExtensionInfo, error)
	ListCallLogs(context.Context, *types.ListCallLogsRequest) (*types.CallLogsResponse, error)
}

// +apix:path=/shop.Contact
type ContactService interface {
	ListContacts(context.Context, *types.ListContactsRequest) (*types.ContactsResponse, error)
	CreateContact(context.Context, *types.CreateContactRequest) (*types.Contact, error)
}
