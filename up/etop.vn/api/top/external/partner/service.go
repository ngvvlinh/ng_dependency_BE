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

// +apix:path=/partner.Customer
type CustomerService interface {
	GetCustomers(context.Context, *types.GetCustomersRequest) (*types.CustomersResponse, error)
}

// +apix:path=/partner.Product
type ProductService interface {
	GetProducts(context.Context, *types.GetProductsRequest) (*types.ShopProductsResponse, error)
}

// +apix:path=/partner.Variant
type VariantService interface {
	GetVariants(context.Context, *types.GetVariantsRequest) (*types.ShopVariantsResponse, error)
}
