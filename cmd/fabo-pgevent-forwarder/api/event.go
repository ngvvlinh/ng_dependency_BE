package api

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/lib/pq"

	api "o.o/api/top/services/pgevent"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/eventhandler/pgevent"
	cm "o.o/backend/pkg/common"
	"o.o/common/l"
)

type EventService struct {
	PgService *pgevent.Service
}

func (s *EventService) Clone() api.EventService {
	res := *s
	return &res
}

func (s *EventService) GenerateEvents(ctx context.Context, q *api.GenerateEventsRequest) (*pbcm.Empty, error) {
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
			return nil, cm.Errorf(cm.InvalidArgument, nil, "invalid raw_events_pg format")
		}
		rawEventsPg := q.RawEventsPg[1 : len(q.RawEventsPg)-1]
		events = strings.Split(rawEventsPg, ",")
	}
	if conditioner != 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide exactly one argument")
	}

	itemsPerBatch := q.ItemsPerBatch
	if itemsPerBatch <= 0 {
		itemsPerBatch = len(events)
	}
	n := min(itemsPerBatch, len(events))
	if err := s.generateEvents(ctx, events[:n]); err != nil {
		return nil, err
	}
	events = events[n:]
	ll.Info("generated events", l.Int("n", n))

	if len(events) > 0 {
		go func() {
			for len(events) > 0 {
				n := min(itemsPerBatch, len(events))
				err := s.generateEvents(ctx, events[:n])
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

	return &pbcm.Empty{}, nil
}

func (s *EventService) generateEvents(ctx context.Context, events []string) error {
	for _, event := range events {
		fakeEvent := &pq.Notification{
			BePid:   0,
			Channel: pgevent.PgChannel,
			Extra:   event,
		}
		err := s.PgService.HandleNotificationWithError(fakeEvent)
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
