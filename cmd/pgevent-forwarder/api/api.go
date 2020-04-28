package api

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/lib/pq"

	pbcm "o.o/api/top/types/common"
	"o.o/backend/com/handler/pgevent"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/common/l"
)

func init() {
	bus.AddHandlers("pgevent/api",
		miscService.VersionInfo,
		eventService.GenerateEvents,
	)
}

var pgservice *pgevent.Service
var ll = l.New()

type MiscService struct{}
type EventService struct{}

var miscService = &MiscService{}
var eventService = &EventService{}

func Init(s *pgevent.Service) {
	pgservice = s
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "pgevent-forwarder",
		Version: "0.1",
	}
	return nil
}

func (s *EventService) GenerateEvents(ctx context.Context, q *GenerateEventsEndpoint) error {
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

	itemsPerBatch := q.ItemsPerBatch
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
		err := pgservice.HandleNotificationWithError(fakeEvent)
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
