package eventstream

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

type Publisher interface {
	Publish(event Event)
}

type Event struct {
	Type string

	Global    bool // send to all users
	AccountID dot.ID
	UserID    dot.ID

	Payload interface{}

	retryInSecond int
}

type Subscriber struct {
	ID int64

	AccountID dot.ID
	UserID    dot.ID

	AllEvents bool
	Events    []string

	ch chan *Event
}

type EventStream struct {
	subscribers  map[int64]*Subscriber
	eventChannel chan *Event
	ctx          context.Context

	m sync.RWMutex
}

func New(ctx context.Context) *EventStream {
	es := &EventStream{
		subscribers:  make(map[int64]*Subscriber),
		eventChannel: make(chan *Event, 256),
		ctx:          ctx,
	}
	go es.RunForwarder()
	return es
}

func (s *EventStream) Publish(event Event) {
	s.eventChannel <- &event
}

func (s *EventStream) RunForwarder() {
	for event := range s.eventChannel {
		s.forward(event)
	}
}

func (s *EventStream) forward(event *Event) {
	s.m.RLock()
	defer s.m.RUnlock()

	ll.Debug("eventstream: received event", l.Any("event", event))
	for _, subscriber := range s.subscribers {
		if ShouldSendEvent(event, subscriber) {
			select {
			case subscriber.ch <- event:
				ll.Debug("send event to", l.ID("Name", subscriber.AccountID), l.ID("UserID", subscriber.AccountID), l.Any("event", event))

			default:
				ll.Info("out of channel buffer, drop event")
			}
		}
	}
}

func ShouldSendEvent(event *Event, subscriber *Subscriber) bool {
	return event.Global ||
		(event.AccountID != 0 && event.AccountID == subscriber.AccountID) ||
		(event.UserID != 0 && event.UserID == subscriber.UserID)
}

func (s *EventStream) SubscribeShop(userID dot.ID) (id int64, ch chan *Event) {
	subscriber := &Subscriber{
		ID:        cm.RandomInt64(),
		AllEvents: true,
		UserID:    userID,

		ch: make(chan *Event, 16),
	}

	s.m.Lock()
	s.subscribers[id] = subscriber
	defer s.m.Unlock()
	return id, subscriber.ch
}

func (s *EventStream) Unsubscribe(id int64) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.subscribers, id)
}

func (s *EventStream) HandleEventStream(c *httpx.Context) error {
	userID := c.Session.GetUserID()
	ctx := c.Context()

	subscriberID, eventChannel := s.SubscribeShop(userID)
	defer s.Unsubscribe(subscriberID)

	w := c.SetResultRaw()
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	w.WriteHeader(200)

	writeEvent(w, &Event{Type: "ping", Payload: "{}", retryInSecond: 3})
	w.(http.Flusher).Flush()

	// flushTimer is not init, as nil channel will be blocked
	var flushTimer <-chan time.Time
	pingTimer := time.NewTicker(10 * time.Second)
	defer pingTimer.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-s.ctx.Done():
			return nil

		case <-flushTimer:
			w.(http.Flusher).Flush()
			flushTimer = nil

		case event := <-eventChannel:
			writeEvent(w, event)
			if flushTimer == nil {
				t := time.NewTimer(100 * time.Millisecond)
				flushTimer = t.C
			}

		case <-pingTimer.C:
			writeEvent(w, &Event{Type: "ping", Payload: "{}"})
			if flushTimer == nil {
				t := time.NewTimer(100 * time.Millisecond)
				flushTimer = t.C
			}
		}
	}
}

func writeEvent(w http.ResponseWriter, event *Event) {
	if event.retryInSecond != 0 {
		_, _ = fmt.Fprintf(w, "retry: %d000\n", event.retryInSecond)
	}
	if event.Type != "" {
		_, _ = fmt.Fprintf(w, "event: %s\n", event.Type)
	}
	switch payload := event.Payload.(type) {
	case []byte:
		_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)

	case string:
		_, _ = fmt.Fprintf(w, "data: %s\n\n", payload)

	default:
		_, _ = fmt.Fprint(w, "data: ")
		// TODO: must implement capi.Message
		if err := jsonx.MarshalTo(w, payload); err != nil {
			panic(err)
		}
		_, _ = fmt.Fprint(w, "\n")
	}
}
