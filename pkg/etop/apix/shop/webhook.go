package xshop

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/apix/webhook"
	wrapxshop "etop.vn/backend/wrapper/external/shop"
)

func init() {
	bus.AddHandlers("apix",
		CreateWebhook,
		DeleteWebhook,
		GetWebhooks,
		GetChanges,
	)
}

func CreateWebhook(ctx context.Context, r *wrapxshop.CreateWebhookEndpoint) error {
	resp, err := webhook.CreateWebhook(ctx, r.Context.Shop.ID, r.CreateWebhookRequest)
	r.Result = resp
	return err
}

func DeleteWebhook(ctx context.Context, r *wrapxshop.DeleteWebhookEndpoint) error {
	resp, err := webhook.DeleteWebhook(ctx, r.Context.Shop.ID, r.DeleteWebhookRequest)
	r.Result = resp
	return err
}

func GetWebhooks(ctx context.Context, r *wrapxshop.GetWebhooksEndpoint) error {
	resp, err := webhook.GetWebhooks(ctx, r.Context.Shop.ID)
	r.Result = resp
	return err
}

func GetChanges(ctx context.Context, r *wrapxshop.GetChangesEndpoint) error {
	return cm.ErrTODO
}
