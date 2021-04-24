package handler

import (
	"context"

	"o.o/backend/com/fabo/pkg/webhook"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/common/jsonx"
)

func (h *WebhookFacebookHandler) HandleWebhookFbMessageEvent(
	ctx context.Context, value []byte,
) (mq.Code, error) {
	var webhookMessage webhook.WebhookMessages
	if err := jsonx.Unmarshal(value, &webhookMessage); err != nil {
		return mq.CodeStop, err
	}

	ctx = bus.Ctx()
	return h.wh.HandleMessenger(ctx, webhookMessage)
}
