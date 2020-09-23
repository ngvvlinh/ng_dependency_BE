package vnp

import (
	"context"

	"o.o/api/top/external/shop"
	"o.o/api/top/external/types"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/mc/vnp
// +gen:swagger:title: MoveCrop - VNPost
// +gen:swagger:version=v1
// +gen:swagger:description=description.md

var _ shop.ShipnowService = ShipnowService(nil)

// +apix:path=/vnposts
// +swagger:tag: Giao tức thì
type ShipnowService interface {

	// Ping
	//
	// Kiểm tra token là hợp lệ và server đang hoạt động.
	//
	// +apix:path=ping
	Ping(context.Context, *cm.Empty) (*cm.Empty, error)

	// Get Services
	//
	// Lấy danh sách các dịch vụ tức thì.
	//
	// +apix:path=getservicesvnpost
	GetShipnowServices(context.Context, *types.GetShipnowServicesRequest) (*types.GetShipnowServicesResponse, error)

	// Create Order
	//
	// Tạo đơn giao tức thì.
	//
	// +apix:path=createordervnpost
	CreateShipnowFulfillment(context.Context, *types.CreateShipnowFulfillmentRequest) (*types.ShipnowFulfillment, error)

	// Cancel Order
	//
	// Huỷ đơn giao tức thì.
	//
	// +apix:path=cancelordervnpost
	CancelShipnowFulfillment(context.Context, *types.CancelShipnowFulfillmentRequest) (*cm.UpdatedResponse, error)

	// Get Order
	//
	// Lấy thông tin đơn giao tức thì.
	//
	// +apix:path=getordervnpost
	GetShipnowFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.ShipnowFulfillment, error)
}

// +apix:path=/vnposts/webhook
type WebhookService interface {
	// Create Webhook
	//
	// Tạo webhook.
	//
	// +apix:path=createwebhook
	CreateWebhook(context.Context, *CreateWebhookRequest) (*Webhook, error)

	// Get Webhooks
	//
	// Lấy danh sách webhooks đã đăng ký.
	//
	// +apix:path=getwebhooks
	GetWebhooks(context.Context, *cm.Empty) (*WebhooksResponse, error)

	// Delete Webhook
	//
	// Xóa webhook
	//
	// +apix:path=deletewebhook
	DeleteWebhook(context.Context, *types.DeleteWebhookRequest) (*WebhooksResponse, error)

	// This API provides an example for webhook data. It's not a real API.
	GetChanges(context.Context, *cm.Empty) (*DataCallback, error)
}
