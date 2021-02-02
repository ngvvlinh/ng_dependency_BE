package handler

import (
	"context"
	"o.o/backend/com/fabo/pkg/webhook"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
)

func (h *WebhookFacebookHandler) HandleWebhookFbFeedEvent(
	ctx context.Context, webhookMessage *webhook.WebhookMessages,
) (mq.Code, error) {
	ctx = bus.Ctx()
	return h.wh.HandleFeed(ctx, *webhookMessage)
}
