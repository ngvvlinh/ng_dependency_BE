package handler

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"o.o/backend/com/eventhandler"
	"o.o/backend/pkg/common/mq"
	"o.o/common/l"
)

var ll = l.New()

type Handler struct {
	consumer mq.KafkaConsumer
	prefix   string
	wg       sync.WaitGroup
}

func New(
	consumer mq.KafkaConsumer,
	prefix string,
) *Handler {
	h := &Handler{
		consumer: consumer,
		prefix:   prefix + "_pgrid_",
	}
	return h
}

func (h *Handler) StartConsuming(ctx context.Context, topics []eventhandler.TopicDef, handlers map[string]mq.EventHandler) {
	for _, topic := range topics {
		handler := handlers[topic.Name]
		if handler == nil {
			ll.Info("No handler for topic", l.String("topic", topic.Name))
			continue
		}
		h.consumeAndHandleTopic(ctx, topic, handler)
	}
}

func (h *Handler) consumeAndHandleTopic(ctx context.Context, d eventhandler.TopicDef, handler mq.EventHandler) {
	count := 0
	var m sync.Mutex
	var wg sync.WaitGroup
	kafkaTopic := h.prefix + d.Name
	wg.Add(d.Partitions)
	h.wg.Add(d.Partitions)

	for i := 0; i < d.Partitions; i++ {
		partition := i
		go func() {
			pc, err := h.consumer.Consume(kafkaTopic, partition)
			if err != nil {
				ll.S.Fatalf("Error while consuming topic: %v:%v", kafkaTopic, partition)
				return
			}
			defer h.wg.Done()
			defer func() { _ = pc.Close() }()

			wg.Done()
			m.Lock()
			count++
			m.Unlock()

			err = pc.ConsumeAndHandle(ctx, handler)
			if err != nil {
				ll.S.Errorf("Handler for topic %v:%v stopped: %+v", kafkaTopic, partition, err)
				buf := make([]byte, 2048)
				runtime.Stack(buf, false)
				msg := fmt.Sprintf(
					"🔥 Handler for topic %v:%v stoppped: %+v\n\n%s",
					kafkaTopic, partition, err, buf)
				ll.SendMessage(msg)
			}
		}()
	}
	wg.Wait()
	m.Lock()
	defer m.Unlock()
}

func (h *Handler) Wait() {
	h.wg.Wait()
}
