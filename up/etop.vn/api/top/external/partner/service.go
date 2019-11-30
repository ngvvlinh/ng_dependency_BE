package partner_proto

import (
	"context"

	cm "etop.vn/api/pb/common"
	"etop.vn/api/pb/etop"
	"etop.vn/api/pb/external"
	partner "etop.vn/api/pb/external/partner"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/partner

// +apix:path=/partner.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
	CurrentAccount(context.Context, *cm.Empty) (*external.Partner, error)
	GetLocationList(context.Context, *cm.Empty) (*external.LocationResponse, error)
}

// +apix:path=/partner.Shop
type ShopService interface {
	AuthorizeShop(context.Context, *partner.AuthorizeShopRequest) (*partner.AuthorizeShopResponse, error)
	CurrentShop(context.Context, *cm.Empty) (*etop.PublicAccountInfo, error)
}

// +apix:path=/partner.Webhook
type WebhookService interface {
	CreateWebhook(context.Context, *external.CreateWebhookRequest) (*external.Webhook, error)
	GetWebhooks(context.Context, *cm.Empty) (*external.WebhooksResponse, error)
	DeleteWebhook(context.Context, *external.DeleteWebhookRequest) (*external.WebhooksResponse, error)
}

// +apix:path=/partner.History
type HistoryService interface {
	GetChanges(context.Context, *external.GetChangesRequest) (*external.Callback, error)
}

// +apix:path=/partner.Shipping
type ShippingService interface {
	GetShippingServices(context.Context, *external.GetShippingServicesRequest) (*external.GetShippingServicesResponse, error)
	CreateAndConfirmOrder(context.Context, *external.CreateOrderRequest) (*external.OrderAndFulfillments, error)
	CancelOrder(context.Context, *external.CancelOrderRequest) (*external.OrderAndFulfillments, error)
	GetOrder(context.Context, *external.OrderIDRequest) (*external.OrderAndFulfillments, error)
	GetFulfillment(context.Context, *external.FulfillmentIDRequest) (*external.Fulfillment, error)
}
