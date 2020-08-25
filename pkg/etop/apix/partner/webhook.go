package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/webhook"
	"o.o/backend/pkg/etop/authorize/session"
)

type WebhookService struct {
	session.Session

	WebhookInner *webhook.Service
}

func (s *WebhookService) Clone() api.WebhookService { res := *s; return &res }

func (s *WebhookService) CreateWebhook(ctx context.Context, r *externaltypes.CreateWebhookRequest) (*externaltypes.Webhook, error) {
	resp, err := s.WebhookInner.CreateWebhook(ctx, s.SS.Partner().ID, r)
	return resp, err
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, r *externaltypes.DeleteWebhookRequest) (*externaltypes.WebhooksResponse, error) {
	resp, err := s.WebhookInner.DeleteWebhook(ctx, s.SS.Partner().ID, r)
	return resp, err
}

func (s *WebhookService) GetWebhooks(ctx context.Context, r *pbcm.Empty) (*externaltypes.WebhooksResponse, error) {
	resp, err := s.WebhookInner.GetWebhooks(ctx, s.SS.Partner().ID)
	return resp, err
}
