package handler

import (
	"fmt"
	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
	"o.o/backend/com/eventhandler"
	"o.o/backend/com/fabo/pkg/webhook"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/common/jsonx"
	"o.o/common/xerrors"
)

type HandlerFunc func(context.Context, *webhook.WebhookMessages) (mq.Code, error)

func WrapHandlerFunc(fn HandlerFunc) mq.EventHandler {
	if fn == nil {
		return nil
	}
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
		var webhookContent webhook.WebhookMessages
		if err := jsonx.Unmarshal(msg.Value, &webhookContent); err != nil {
			return mq.CodeStop, wrapError(err, msg, nil)
		}

		code, err := fn(ctx, &webhookContent)
		if err != nil {
			err = wrapError(err, msg, &webhookContent)
		}
		return code, err
	}
}

func wrapError(err error, msg *sarama.ConsumerMessage, event *webhook.WebhookMessages) *xerrors.APIError {
	return cm.MapError(err).Throw().
		WithMeta("topic", fmt.Sprintf("%v:%v", msg.Topic, msg.Partition)).
		WithMetaJson("event", event)
}

func WrapMapHandlers(handlers map[string]HandlerFunc) map[string]mq.EventHandler {
	m := make(map[string]mq.EventHandler)
	for name, h := range handlers {
		m[name] = WrapHandlerFunc(h)
	}
	return m
}

type WebhookFacebookHandler struct {
	wh *webhook.WebhookHandler
}

func NewWebhookFacebookHandler(wh *webhook.WebhookHandler) *WebhookFacebookHandler {
	h := &WebhookFacebookHandler{
		wh: wh,
	}
	return h
}

func (h *WebhookFacebookHandler) TopicsAndHandlers() map[string]mq.EventHandler {
	return WrapMapHandlers(map[string]HandlerFunc{
		"facebook_webhook_feed":    h.HandleWebhookFbFeedEvent,
		"facebook_webhook_message": h.HandleWebhookFbMessageEvent,
	})
}

func (h *WebhookFacebookHandler) Topics() []eventhandler.TopicDef {
	return []eventhandler.TopicDef{
		{
			Name:       "facebook_webhook_message",
			Partitions: 64,
		},
		{
			Name:       "facebook_webhook_feed",
			Partitions: 64,
		},
	}
}
