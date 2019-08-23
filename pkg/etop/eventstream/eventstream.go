package eventstream

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	cmService "etop.vn/backend/pkg/common/service"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/common/l"
)

var ll = l.New()

type Publisher interface {
	Publish(event Event)
}

type Event struct {
	Type string

	Global    bool // send to all users
	AccountID int64
	UserID    int64
	Payload   interface{}

	retryInSecond int
}

type Subscriber struct {
	ID int64

	AccountID int64
	UserID    int64

	AllEvents bool
	Events    []string

	ch chan *Event
}

type EventStreamer struct {
	subscribers  map[int64]*Subscriber
	eventChannel chan *Event
	shutdowner   cmService.Shutdowner

	m sync.RWMutex
}

func NewEventStreamer(sd cmService.Shutdowner) *EventStreamer {
	return &EventStreamer{
		subscribers:  make(map[int64]*Subscriber),
		eventChannel: make(chan *Event, 256),
		shutdowner:   sd,
	}
}

func (s *EventStreamer) Publish(event Event) {
	s.eventChannel <- &event
}

func (s *EventStreamer) RunForwarder() {
	for event := range s.eventChannel {
		s.forward(event)
	}
}

func (s *EventStreamer) forward(event *Event) {
	s.m.RLock()
	defer s.m.RUnlock()

	ll.Debug("eventstream: received event", l.Any("event", event))
	for _, subscriber := range s.subscribers {
		if ShouldSendEvent(event, subscriber) {
			select {
			case subscriber.ch <- event:
				ll.Debug("send event to", l.Int64("Name", subscriber.AccountID), l.Int64("UserID", subscriber.UserID), l.Any("event", event))

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

func (s *EventStreamer) Subscribe(accountID int64, userID int64) (id int64, ch <-chan *Event) {
	subscriber := &Subscriber{
		ID:        cm.NewID(),
		AllEvents: true,
		AccountID: accountID,
		UserID:    userID,

		ch: make(chan *Event, 16),
	}

	s.m.Lock()
	defer s.m.Unlock()
	s.subscribers[id] = subscriber
	return subscriber.ID, subscriber.ch
}

func (s *EventStreamer) Unsubscribe(id int64) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.subscribers, id)
}

func (s *EventStreamer) HandleEventStream(c *httpx.Context) error {
	claim := c.Claim.(*claims.ShopClaim)
	userID := c.Session.GetUserID()
	shop := claim.Shop
	ctx := c.Context()
	// TODO(qv): Limit connections per user

	subscriberID, eventChannel := s.Subscribe(shop.ID, userID)
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

		case <-s.shutdowner.Done():
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

var marshaler = jsonpb.Marshaler{OrigName: true, EmitDefaults: true}

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

	case proto.Message:
		_, _ = fmt.Fprint(w, "data: ")
		_ = marshaler.Marshal(w, payload)
		_, _ = fmt.Fprint(w, "\n\n")

	default:
		panic("unsupported payload type")
	}
}
