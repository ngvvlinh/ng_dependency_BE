package shop

import (
	"context"

	cm "etop.vn/api/pb/common"
	"etop.vn/api/pb/etop"
	"etop.vn/api/pb/external"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/shop

// +apix:path=/shop.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*etop.PublicAccountInfo, error)
	GetLocationList(context.Context, *cm.Empty) (*external.LocationResponse, error)
}

// +apix:path=/shop.Webhook
type WebhookService interface {
	CreateWebhook(context.Context, *external.CreateWebhookRequest) (*external.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*external.WebhooksResponse, error)
	DeleteWebhook(context.Context, *external.DeleteWebhookRequest) (*external.WebhooksResponse, error)
}

// +apix:path=/shop.History
type HistoryService interface {
	GetChanges(context.Context, *external.GetChangesRequest) (*external.Callback, error)
}

// +apix:path=/shop.Shipping
type ShippingService interface {
	GetShippingServices(context.Context, *external.GetShippingServicesRequest) (*external.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *external.CreateOrderRequest) (*external.OrderAndFulfillments, error)
	CancelOrder(context.Context, *external.CancelOrderRequest) (*external.OrderAndFulfillments, error)
	GetOrder(context.Context, *external.OrderIDRequest) (*external.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *external.FulfillmentIDRequest) (*external.Fulfillment, error)
}
