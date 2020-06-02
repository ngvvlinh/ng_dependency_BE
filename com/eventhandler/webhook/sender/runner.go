package sender

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"o.o/backend/com/eventhandler/webhook/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

const PrefixRedisWebhook = "webhook_states:"
const MaxChangesItems = 1000
const MaxPayloadSize = 512 * 1024 // 512KB
const MinInterval = 5 * time.Second
const MaxInterval = 24 * time.Hour
const DefaultTTL = 25 * 60 * 60 // 25 hours

// if a webhook is unavailable during 24 hours, we'll disable it
var retryIntervals = []time.Duration{
	30 * time.Second,
	time.Minute,
	2 * time.Minute,
	4 * time.Minute,
	8 * time.Minute,
	15 * time.Minute,
	30 * time.Minute,
	time.Hour,
	2 * time.Hour,
	4 * time.Hour,
	8 * time.Hour,
	12 * time.Hour,
	16 * time.Hour,
	MaxInterval,
}

type State string

func (s State) String() string { return string(s) }

const (
	StateOK    State = "ok"
	StateRetry State = "retry"
	StateStop  State = "stop"
)

type WebhookStatesError struct {
	SentAt   time.Time `json:"sent_at"`
	Status   int       `json:"status,omitempty"`
	Response string    `json:"response,omitempty"`
	Retried  int       `json:"retried"`
	ErrorMsg string    `json:"error_msg"`
}

type WebhookStates struct {
	IntervalDuration time.Duration `json:"-"`
	Interval         int           `json:"interval"`

	LastSentAt time.Time           `json:"last_sent_at"`
	LastError  *WebhookStatesError `json:"last_error"`
	State      State               `json:"state"`
}

type SingleSender struct {
	webhook *model.Webhook

	lastItems *types.MessageCollector
	items     *types.MessageCollector

	m     sync.Mutex
	stop  chan struct{}
	reset chan struct{}
}

func NewSingleSender(wh *model.Webhook) *SingleSender {
	return &SingleSender{
		webhook: wh,
		stop:    make(chan struct{}),
		reset:   make(chan struct{}, 1), // make the channel not block
	}
}

func LoadWebhookStates(redisStore redis.Store, id dot.ID) WebhookStates {
	var current WebhookStates
	err := redisStore.Get(redisKey(id), &current)
	if err != nil && err != redis.ErrNil {
		ll.Error("Can not load from redis", l.Error(err))
	}

	if current.State == "" {
		current = WebhookStates{
			IntervalDuration: MinInterval,
			State:            StateOK,
		}
	} else {
		current.IntervalDuration = time.Duration(current.Interval) * time.Second
	}
	return current
}

func NewSingleSenders(webhooks []*model.Webhook) []*SingleSender {
	ssenders := make([]*SingleSender, len(webhooks))
	for i, item := range webhooks {
		ss := NewSingleSender(item)
		ssenders[i] = ss
	}
	return ssenders
}

func (s *SingleSender) Collect(entity string, entityID dot.ID, msg []byte) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.items == nil {
		s.items = &types.MessageCollector{}
	}
	s.items.Collect(msg)

	// TODO: flush data
	if len(s.items.Messages) >= 2*MaxChangesItems {
		s.items.Messages = truncate(s.items.Messages, MaxChangesItems)
		// 	flushItems := s.items
		// 	s.items = &whtypes.MessageCollector{}
		// 	s.storeToDatabase(flushItems)
	}
}

func (s *SingleSender) Shutdown() {
	// TODO
}

func (s *SingleSender) storeToDatabase(callbackID dot.ID, mc *types.MessageCollector, states *WebhookStatesError) error {
	// TODO: refactor
	changesData := buildJSON(callbackID, mc.Messages)
	changesData = changesData[len(jsonOpen)-1 : len(changesData)-1]
	statesData, _ := jsonx.Marshal(states)

	data := &model.Callback{
		ID:        callbackID,
		WebhookID: s.webhook.ID,
		AccountID: s.webhook.AccountID,
		CreatedAt: time.Now(),
		Changes:   changesData,
		Result:    statesData,
	}
	ll.Debug("store to database", l.ID("id", data.ID), l.ID("account_id", s.webhook.AccountID), l.ID("webhook_id", s.webhook.ID))
	return changesStore.Insert(context.Background(), data)
}

func redisKey(id dot.ID) string {
	return PrefixRedisWebhook + id.String()
}

func truncate(items [][]byte, max int) [][]byte {
	if len(items) <= max {
		return items
	}
	for i, item := range items[len(items)-max:] {
		items[i] = item
	}
	return items[max:]
}

func (s *SingleSender) Run(ctx context.Context, startAfter time.Duration) {
	ll.Debug("Sender start", l.ID("webhook_id", s.webhook.ID), l.ID("account_id", s.webhook.AccountID))
	defer ll.Debug("Sender stopped", l.ID("webhook_id", s.webhook.ID), l.ID("account_id", s.webhook.AccountID))

	current := LoadWebhookStates(redisStore, s.webhook.ID)
	if current.State == StateStop {
		ll.Warn("Webhook is stopped", l.ID("webhook_id", s.webhook.ID))
		current.IntervalDuration = 24 * time.Hour
		// TODO: refactor this
	}

	interval := startAfter + current.IntervalDuration
	t := time.NewTimer(interval)
	for {
		ll.Debug("Sender is running", l.ID("webhook_id", s.webhook.ID), l.ID("account_id", s.webhook.AccountID), l.Any("states", current))

		select {
		case <-ctx.Done():
			return

		case <-s.stop:
			return

		case <-s.reset:
			t.Stop()
			current.IntervalDuration = 5 * time.Second
			t.Reset(current.IntervalDuration)
			ll.Debug("reset webhook state", l.ID("webhook_id", s.webhook.ID), l.Any("states", current))

		case <-t.C:
			states, err := s.Send()
			if err == nil && states == nil {
				t = time.NewTimer(current.IntervalDuration)
				continue
			}

			current = calcNextStates(s.webhook, current, states, err)
			current.Interval = int(current.IntervalDuration / time.Second)
			t = time.NewTimer(current.IntervalDuration)

			ll.Debug("Store to redis", l.ID("webhook_id", s.webhook.ID), l.Any("states", current))
			if err := redisStore.SetWithTTL(redisKey(s.webhook.ID), current, DefaultTTL); err != nil {
				ll.Error("Can not store to redis", l.Error(err))
			}
		}
	}
}

func (s *SingleSender) ResetState() {
	select {
	case s.reset <- struct{}{}:
	default:
	}
}

func calcNextStates(webhook *model.Webhook, current WebhookStates, states *WebhookStatesError, err error) WebhookStates {
	if err == nil {
		// nothing was sent
		if states == nil {
			return current
		}

		next := WebhookStates{
			IntervalDuration: MinInterval,
			LastSentAt:       states.SentAt,
			LastError:        current.LastError,
			State:            StateOK,
		}
		return next
	}

	nextDuration := MinInterval
	states.ErrorMsg = err.Error()
	if current.State == StateRetry && current.LastError != nil {
		states.Retried = current.LastError.Retried + 1
		nextDuration = current.IntervalDuration
	}
	for _, step := range retryIntervals {
		if step > nextDuration {
			nextDuration = step
			break
		}
	}
	if nextDuration >= MaxInterval {
		next := WebhookStates{
			IntervalDuration: nextDuration,
			LastSentAt:       states.SentAt,
			LastError:        states,
			State:            StateStop,
		}
		return next
	}

	next := WebhookStates{
		IntervalDuration: nextDuration,
		LastSentAt:       states.SentAt,
		LastError:        states,
		State:            StateRetry,
	}
	return next
}

func (s *SingleSender) Send() (*WebhookStatesError, error) {
	// retry the last items if it's not success
	items := s.lastItems
	if items == nil || len(items.Messages) == 0 {
		s.m.Lock()

		// early return
		if s.items == nil || len(s.items.Messages) == 0 {
			s.m.Unlock()
			return nil, nil
		}

		// copy items for sending
		items = s.items
		s.lastItems = s.items
		s.items = nil
		s.m.Unlock()
	}

	callbackID := cm.NewID()
	data := buildJSON(callbackID, truncate(items.Messages, MaxChangesItems))
	wh := s.webhook
	status, respData, err := sendWebhookSingleRequest(context.Background(), wh, data)

	lastItems := s.lastItems
	if err == nil {
		// clean the last items
		s.lastItems = nil
		ll.Debug("Sent webhook", l.ID("account_id", wh.AccountID), l.ID("webhook_id", wh.ID), l.Int("status_code", status), l.String("url", wh.URL))

	} else {
		ll.Warn("Sent webhook (unexpected status)", l.ID("account_id", wh.AccountID), l.ID("webhook_id", wh.ID), l.Int("status_code", status), l.Error(err), l.String("url", wh.URL))
	}

	states := &WebhookStatesError{
		SentAt:   time.Now(),
		Status:   status,
		Response: cm.UnsafeBytesToString(respData),
	}
	if err != nil {
		states.ErrorMsg = err.Error()
	}
	if err := s.storeToDatabase(callbackID, lastItems, states); err != nil {
		ll.Error("can not store to database", l.Error(err))
	}
	return states, err
}

func sendWebhookSingleRequest(ctx context.Context, wh *model.Webhook, data []byte) (status int, respData []byte, _ error) {
	res, err := client.Post(wh.URL, "application/json", bytes.NewReader(data))
	if err != nil {
		return 0, nil, err
	}

	status = res.StatusCode
	if status != 200 {
		return status, nil, fmt.Errorf("unexpected status: %v", status)
	}

	var b [256]byte
	n, err := res.Body.Read(b[:])
	if err != nil && err != io.EOF {
		return status, nil, err
	}
	respData = b[:n]
	if b[0] == 'o' && b[1] == 'k' &&
		(n == 2 || n == 3 && b[2] == '\n') {
		return status, respData, nil
	}
	return status, respData, fmt.Errorf("unexpected response")
}

const jsonOpen = `"changes":[`
const jsonClose = `]}`

func buildJSON(callbackID dot.ID, msgs [][]byte) []byte {
	size := 1 + len(jsonOpen) + len(jsonClose) + 30 // id:"",
	for _, msg := range msgs {
		size += len(msg) + 1
	}

	data := make([]byte, 0, size)
	data = append(data, `{"id":"`...)
	data = strconv.AppendInt(data, callbackID.Int64(), 10)
	data = append(data, `",`...)
	data = append(data, jsonOpen...)
	for _, msg := range msgs {
		data = append(data, msg...)
		data = append(data, ',')
	}
	data = data[:len(data)-1]
	data = append(data, jsonClose...)
	return data
}
