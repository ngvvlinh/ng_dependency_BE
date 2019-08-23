package intctl

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"

	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/common/l"
)

var ll = l.New()

const ConsumerGroup = "handler/intctl"

type Handler struct {
	prefix   string
	consumer mq.KafkaConsumer
	bot      *telebot.Channel
	init     bool
	wg       sync.WaitGroup

	handlers  map[string]mq.EventHandler
	listeners map[string]mq.EventHandler
}

func New(bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string) *Handler {
	h := &Handler{
		prefix:    prefix + "_",
		consumer:  consumer,
		bot:       bot,
		listeners: make(map[string]mq.EventHandler),
	}
	handlers := h.TopicsAndHandlers()
	h.handlers = handlers
	return h
}

func (h *Handler) Subscribe(channel string, listener mq.EventHandler) {
	if h.init {
		ll.Fatal("can not subscribe: already initialized")
	}
	if h.listeners[channel] != nil {
		ll.Fatal("can not subscribe: duplicated listener", l.String("channel", channel))
	}
	h.listeners[channel] = listener
}

func (h *Handler) TopicsAndHandlers() map[string]mq.EventHandler {
	return map[string]mq.EventHandler{
		"intctl": h.HandleInternalControl,
	}
}

func (h *Handler) ConsumeAndHandle(ctx context.Context) {
	h.init = true // prevent calling Subscribe()

	var wg sync.WaitGroup
	var topics []string
	for topic, handler := range h.handlers {
		topics = append(topics, topic)
		partition := 0 // single partition
		if handler == nil {
			ll.S.Fatalf("no handler for topic: %v", topic)
		}

		kafkaTopic := h.prefix + topic
		wg.Add(1)
		h.wg.Add(partition + 1) // single partition
		go func() {
			pc, err := h.consumer.Consume(kafkaTopic, int32(partition))
			if err != nil {
				ll.S.Fatalf("Error while consuming topic: %v:%v", kafkaTopic, partition)
				panic(err)
			}
			defer h.wg.Done()
			defer ignoreError(pc.Close())
			wg.Done()

			err = pc.ConsumeAndHandle(ctx, handler)
			if err != nil {
				ll.S.Errorf("Handler for topic %v:%v stopped: %+v", kafkaTopic, partition, err)
				if h.bot != nil {
					msg := fmt.Sprintf(
						"🔥 Handler for topic %v:%v stoppped: %+v",
						kafkaTopic, partition, err)
					h.bot.SendMessage(msg)
				}
			}
		}()
	}
	wg.Wait()
	ll.S.Infof("Initialized consumers for topics: %v", strings.Join(topics, ", "))
}

func (h *Handler) Wait() {
	h.wg.Wait()
}

func (h *Handler) HandleInternalControl(ctx context.Context, event *sarama.ConsumerMessage) (mq.Code, error) {
	channel := ParseKey(event.Key)
	listener := h.listeners[channel]
	if listener == nil {
		return mq.CodeIgnore, nil
	}
	ll.Debug("intctl: dispatched listener", l.String("channel", channel))
	return listener(ctx, event)
}

func ignoreError(err error) {}
