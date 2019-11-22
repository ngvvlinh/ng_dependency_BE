package partner_proto

import (
	"context"

	cm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
	"etop.vn/backend/pb/external"
	partner "etop.vn/backend/pb/external/partner"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:apix:doc-path=external/partner

// +apix:path=/partner.Misc
type MiscAPI interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*external.Partner, error)
	GetLocationList(context.Context, *cm.Empty) (*external.LocationResponse, error)
}

// +apix:path=/partner.Shop
type ShopAPI interface {
	AuthorizeShop(context.Context, *partner.AuthorizeShopRequest) (*partner.AuthorizeShopResponse, error)
	CurrentShop(context.Context, *cm.Empty) (*etop.PublicAccountInfo, error)
}

// +apix:path=/partner.Webhook
type WebhookAPI interface {
	CreateWebhook(context.Context, *external.CreateWebhookRequest) (*external.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*external.WebhooksResponse, error)
	DeleteWebhook(context.Context, *external.DeleteWebhookRequest) (*external.WebhooksResponse, error)
}

// +apix:path=/partner.History
type HistoryAPI interface {
	GetChanges(context.Context, *external.GetChangesRequest) (*external.Callback, error)
}

// +apix:path=/partner.Shipping
type ShippingAPI interface {
	GetShippingServices(context.Context, *external.GetShippingServicesRequest) (*external.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *external.CreateOrderRequest) (*external.OrderAndFulfillments, error)
	CancelOrder(context.Context, *external.CancelOrderRequest) (*external.OrderAndFulfillments, error)
	GetOrder(context.Context, *external.OrderIDRequest) (*external.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *external.FulfillmentIDRequest) (*external.Fulfillment, error)
}
