package pgrid

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/pgevent"
)

var ll = l.New()

type HandlerFunc func(context.Context, *pgevent.PgEvent) (mq.Code, error)

type M map[string]interface{}

const tenYears = 365 * 10 * 24 * time.Hour

func WrapHandlerFunc(fn HandlerFunc) mq.EventHandler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
		var event pgevent.PgEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
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

func wrapError(err error, msg *sarama.ConsumerMessage, event *pgevent.PgEvent) *cm.APIError {
	return cm.MapError(err).Throw().
		WithMeta("topic", fmt.Sprintf("%v:%v", msg.Topic, msg.Partition)).
		WithMetaJson("event", event)

}
