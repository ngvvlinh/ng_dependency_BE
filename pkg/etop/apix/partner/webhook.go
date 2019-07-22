package partner

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/apix/webhook"
	"etop.vn/common/bus"

	partnerW "etop.vn/backend/wrapper/external/partner"
)

func init() {
	bus.AddHandlers("apix",
		CreateWebhook,
		DeleteWebhook,
		GetWebhooks,
		GetChanges,
	)
}

func CreateWebhook(ctx context.Context, r *partnerW.CreateWebhookEndpoint) error {
	resp, err := webhook.CreateWebhook(ctx, r.Context.Partner.ID, r.CreateWebhookRequest)
	r.Result = resp
	return err
}

func DeleteWebhook(ctx context.Context, r *partnerW.DeleteWebhookEndpoint) error {
	resp, err := webhook.DeleteWebhook(ctx, r.Context.Partner.ID, r.DeleteWebhookRequest)
	r.Result = resp
	return err
}

func GetWebhooks(ctx context.Context, r *partnerW.GetWebhooksEndpoint) error {
	resp, err := webhook.GetWebhooks(ctx, r.Context.Partner.ID)
	r.Result = resp
	return err
}

func GetChanges(ctx context.Context, r *partnerW.GetChangesEndpoint) error {
	return cm.ErrTODO
}
