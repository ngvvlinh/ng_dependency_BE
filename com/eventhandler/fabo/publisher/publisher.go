package publisher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"

	"o.o/backend/com/eventhandler/fabo/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/common/jsonx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var ll = l.New()

const tenYears = 365 * 10 * 24 * time.Hour

type Publisher struct {
	publisher eventstream.Publisher
}

func New(
	publisherEvent eventstream.Publisher,
) *Publisher {
	h := &Publisher{
		publisher: publisherEvent,
	}
	return h
}

func (h *Publisher) TopicsAndHandlers() map[string]mq.EventHandler {
	return map[string]mq.EventHandler{
		"fb_external_conversation_fabo": wrapHandler(h.HandleFbConversationFaboEvent),
		"fb_external_comment_fabo":      wrapHandler(h.HandleFbCommentFaboEvent),
		"fb_external_message_fabo":      wrapHandler(h.HandleFbMessageFaboEvent),
		"fb_customer_conversation_fabo": wrapHandler(h.HandleFbCustomerConversationFaboEvent),
		"fb_external_post_fabo":         wrapHandler(h.HandleFbPostFaboEvent),
	}
}

func wrapHandler(fn HandlerFunc) mq.EventHandler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
		var event types.FaboEvent
		var err error
		switch {
		case strings.Contains(msg.Topic, "comment") == true:
			err = jsonx.Unmarshal(msg.Value, &event.PgEventComment)
		case strings.Contains(msg.Topic, "message") == true:
			err = jsonx.Unmarshal(msg.Value, &event.PgEventMessage)
		case strings.Contains(msg.Topic, "customer") == true:
			err = jsonx.Unmarshal(msg.Value, &event.PgEventCustomerConversation)
		case strings.Contains(msg.Topic, "post") == true:
			err = jsonx.Unmarshal(msg.Value, &event.PgEventPost)
		case strings.Contains(msg.Topic, "conversation") == true:
			err = jsonx.Unmarshal(msg.Value, &event.PgEventConversation)
		default:
			return mq.CodeStop, wrapError(nil, msg, nil)
		}

		if err != nil {
			return mq.CodeStop, wrapError(err, msg, nil)
		}

		// Skip event too far in the past
		t := time.Unix(event.Timestamp, 0)
		if delta := time.Since(t); delta > 24*time.Hour && delta < tenYears {
			ll.Warn("Skip event",
				l.String("topic", msg.Topic),
				l.Int32("p", msg.Partition),
				l.String("key", string(msg.Key)),
				l.Time("t", msg.Timestamp))
			return mq.CodeOK, nil
		}

		code, err := fn(ctx, &event)
		if err != nil {
			err = wrapError(err, msg, &event)
		}
		return code, err
	}
}

func wrapError(err error, msg *sarama.ConsumerMessage, event *types.FaboEvent) *xerrors.APIError {
	return cm.MapError(err).Throw().
		WithMeta("topic", fmt.Sprintf("%v:%v", msg.Topic, msg.Partition)).
		WithMetaJson("event", event)
}
