package partner

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/apix/webhook"
)

func init() {
	bus.AddHandlers("apix",
		webhookService.CreateWebhook,
		webhookService.DeleteWebhook,
		webhookService.GetWebhooks,
		historyService.GetChanges,
	)
}

func (s *WebhookService) CreateWebhook(ctx context.Context, r *CreateWebhookEndpoint) error {
	resp, err := webhook.CreateWebhook(ctx, r.Context.Partner.ID, r.CreateWebhookRequest)
	r.Result = resp
	return err
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, r *DeleteWebhookEndpoint) error {
	resp, err := webhook.DeleteWebhook(ctx, r.Context.Partner.ID, r.DeleteWebhookRequest)
	r.Result = resp
	return err
}

func (s *WebhookService) GetWebhooks(ctx context.Context, r *GetWebhooksEndpoint) error {
	resp, err := webhook.GetWebhooks(ctx, r.Context.Partner.ID)
	r.Result = resp
	return err
}

func (s *HistoryService) GetChanges(ctx context.Context, r *GetChangesEndpoint) error {
	return cm.ErrTODO
}
