package handler

import (
	"golang.org/x/net/context"
	"o.o/backend/com/fabo/pkg/webhook"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/mq"
	"o.o/common/jsonx"
)

func (h *WebhookFacebookHandler) HandleWebhookUserFbLiveVideo(
	ctx context.Context, value []byte,
) (mq.Code, error) {
	var webhookUser webhook.WebhookUser
	if err := jsonx.Unmarshal(value, &webhookUser); err != nil {
		return mq.CodeIgnore, nil
	}

	ctx = bus.Ctx()
	return h.wh.HandleUserLiveVideo(ctx, webhookUser)
}
