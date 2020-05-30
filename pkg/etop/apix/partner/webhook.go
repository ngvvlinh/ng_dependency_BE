package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/webhook"
)

type WebhookService struct {
	WebhookInner *webhook.Service
}

func (s *WebhookService) Clone() *WebhookService { res := *s; return &res }

func (s *WebhookService) CreateWebhook(ctx context.Context, r *CreateWebhookEndpoint) error {
	resp, err := s.WebhookInner.CreateWebhook(ctx, r.Context.Partner.ID, r.CreateWebhookRequest)
	r.Result = resp
	return err
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, r *DeleteWebhookEndpoint) error {
	resp, err := s.WebhookInner.DeleteWebhook(ctx, r.Context.Partner.ID, r.DeleteWebhookRequest)
	r.Result = resp
	return err
}

func (s *WebhookService) GetWebhooks(ctx context.Context, r *GetWebhooksEndpoint) error {
	resp, err := s.WebhookInner.GetWebhooks(ctx, r.Context.Partner.ID)
	r.Result = resp
	return err
}
