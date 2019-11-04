package external_shop

import (
	"context"

	cm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	"etop.vn/backend/pb/external"
)

// +gen:apix

// +apix:path=/external_shop.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*etop.PublicAccountInfo, error)
	GetLocationList(context.Context, *cm.Empty) (*external.LocationResponse, error)
}

// +apix:path=/external_shop.Webhook
type WebhookAPI interface {
	CreateWebhook(context.Context, *external.CreateWebhookRequest) (*external.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*external.WebhooksResponse, error)
	DeleteWebhook(context.Context, *external.DeleteWebhookRequest) (*external.WebhooksResponse, error)
}

// +apix:path=/external_shop.History
type HistoryAPI interface {
	GetChanges(context.Context, *external.GetChangesRequest) (*external.Callback, error)
}

// +apix:path=/external_shop.Shipping
type ShippingAPI interface {
	GetShippingServices(context.Context, *external.GetShippingServicesRequest) (*external.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *external.CreateOrderRequest) (*external.OrderAndFulfillments, error)
	CancelOrder(context.Context, *external.CancelOrderRequest) (*external.OrderAndFulfillments, error)
	GetOrder(context.Context, *external.OrderIDRequest) (*external.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *external.FulfillmentIDRequest) (*external.Fulfillment, error)
}
