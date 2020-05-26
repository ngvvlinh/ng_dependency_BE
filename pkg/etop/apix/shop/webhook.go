package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/webhook"
)

type WebhookService struct{}

func (s *WebhookService) Clone() *WebhookService { res := *s; return &res }

func (s *WebhookService) CreateWebhook(ctx context.Context, r *CreateWebhookEndpoint) error {
	resp, err := webhook.CreateWebhook(ctx, r.Context.Shop.ID, r.CreateWebhookRequest)
	r.Result = resp
	return err
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, r *DeleteWebhookEndpoint) error {
	resp, err := webhook.DeleteWebhook(ctx, r.Context.Shop.ID, r.DeleteWebhookRequest)
	r.Result = resp
	return err
}

func (s *WebhookService) GetWebhooks(ctx context.Context, r *GetWebhooksEndpoint) error {
	resp, err := webhook.GetWebhooks(ctx, r.Context.Shop.ID)
	r.Result = resp
	return err
}
