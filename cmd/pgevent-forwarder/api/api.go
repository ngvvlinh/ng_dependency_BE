package api

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/lib/pq"

	"etop.vn/backend/com/handler/pgevent"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/common/l"

	pbcm "etop.vn/backend/pb/common"
	wrappgevent "etop.vn/backend/wrapper/services/pgevent"
)

func init() {
	bus.AddHandlers("pgevent/api",
		VersionInfo,
		GenerateEvents,
	)
}

var service *pgevent.Service
var ll = l.New()

func Init(s *pgevent.Service) {
	service = s
}

func VersionInfo(ctx context.Context, q *wrappgevent.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "pgevent-forwarder",
		Version: "0.1",
	}
	return nil
}

func GenerateEvents(ctx context.Context, q *wrappgevent.GenerateEventsEndpoint) error {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println(e)
			debug.PrintStack()
		}
	}()

	var events []string
	conditioner := 0
	if q.RawEvents != nil {
		conditioner++
		events = q.RawEvents
	}
	if q.RawEventsPg != "" {
		conditioner++
		if !strings.HasPrefix(q.RawEventsPg, "{") ||
			!strings.HasSuffix(q.RawEventsPg, "}") {
			return cm.Errorf(cm.InvalidArgument, nil, "invalid raw_events_pg format")
		}
		rawEventsPg := q.RawEventsPg[1 : len(q.RawEventsPg)-1]
		events = strings.Split(rawEventsPg, ",")
	}
	if conditioner != 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "must provide exactly one argument")
	}

	itemsPerBatch := int(q.ItemsPerBatch)
	if itemsPerBatch <= 0 {
		itemsPerBatch = len(events)
	}
	n := min(itemsPerBatch, len(events))
	if err := generateEvents(ctx, events[:n]); err != nil {
		return err
	}
	events = events[n:]
	ll.Info("generated events", l.Int("n", n))

	if len(events) > 0 {
		go func() {
			for len(events) > 0 {
				n := min(itemsPerBatch, len(events))
				err := generateEvents(ctx, events[:n])
				if err != nil {
					ll.Error("error sending events", l.Error(err))
					break
				}
				events = events[n:]
				ll.Info("generated events", l.Int("n", n))

				time.Sleep(30 * time.Second)
			}
		}()
	}

	q.Result = &pbcm.Empty{}
	return nil
}

func generateEvents(ctx context.Context, events []string) error {
	for _, event := range events {
		fakeEvent := &pq.Notification{
			BePid:   0,
			Channel: pgevent.PgChannel,
			Extra:   event,
		}
		err := service.HandleNotificationWithError(fakeEvent)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "invalid event").
				WithMeta("event", event)
		}
	}
	return nil
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
