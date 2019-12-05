package shop

import (
	"context"

	"etop.vn/api/top/external/types"
	etop "etop.vn/api/top/int/etop"
	cm "etop.vn/api/top/types/common"
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
	GetChanges(context.Context, *types.GetChangesRequest) (*types.Callback, error)
}

// +apix:path=/shop.Shipping
type ShippingService interface {
	GetShippingServices(context.Context, *types.GetShippingServicesRequest) (*types.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *types.CreateOrderRequest) (*types.OrderAndFulfillments, error)
	CancelOrder(context.Context, *types.CancelOrderRequest) (*types.OrderAndFulfillments, error)
	GetOrder(context.Context, *types.OrderIDRequest) (*types.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.Fulfillment, error)
}
