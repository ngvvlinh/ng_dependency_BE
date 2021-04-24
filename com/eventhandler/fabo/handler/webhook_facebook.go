package handler

import (
	"fmt"

	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
	"o.o/backend/com/eventhandler"
	"o.o/backend/com/fabo/pkg/webhook"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/common/xerrors"
)

type HandlerFunc func(context.Context, []byte) (mq.Code, error)

func WrapHandlerFunc(fn HandlerFunc) mq.EventHandler {
	if fn == nil {
		return nil
	}
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
		code, err := fn(ctx, msg.Value)
		if err != nil {
			err = wrapError(err, msg, msg.Value)
		}
		return code, err
	}
}

func wrapError(err error, msg *sarama.ConsumerMessage, event []byte) *xerrors.APIError {
	return cm.MapError(err).Throw().
		WithMeta("topic", fmt.Sprintf("%v:%v", msg.Topic, msg.Partition)).
		WithMeta("event", string(event))
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
		"facebook_webhook_feed":            h.HandleWebhookFbFeedEvent,
		"facebook_webhook_message":         h.HandleWebhookFbMessageEvent,
		"facebook_webhook_user_live_video": h.HandleWebhookUserFbLiveVideo,
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
		{
			Name:       "facebook_webhook_user_live_video",
			Partitions: 64,
		},
	}
}
