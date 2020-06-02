package pgevent

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/common/jsonx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

type HandlerFunc func(context.Context, *PgEvent) (mq.Code, error)

const tenYears = 365 * 10 * 24 * time.Hour

func WrapHandlerFunc(fn HandlerFunc) mq.EventHandler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
		var event PgEvent
		if err := jsonx.Unmarshal(msg.Value, &event); err != nil {
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

func wrapError(err error, msg *sarama.ConsumerMessage, event *PgEvent) *xerrors.APIError {
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
