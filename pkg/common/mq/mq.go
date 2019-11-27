package mq

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"

	"etop.vn/common/l"
)

//go:generate stringer -type=Code -output=error_string.gen.go

// Code represents type of error code.
type Code int

// Error code enumaration.
const (
	CodeOK Code = iota
	CodeIgnore
	CodeRetry
	CodeStop
)

// KafkaConsumer represents handler for SingleConsumer
type KafkaConsumer struct {
	consumer *SingleConsumer
	wg       *sync.WaitGroup
}

// NewKafkaConsumer returns a new KafkaConsumer object.
func NewKafkaConsumer(brokers []string, group string, cfgs ...*sarama.Config) (KafkaConsumer, error) {
	var cfg *sarama.Config
	switch {
	case len(cfgs) > 1:
		ll.Panic("Expect single config")
	case len(cfgs) == 1:
		cfg = cfgs[0]
	}

	consumer, err := NewSingleConsumer(brokers, cfg, group)
	if err != nil {
		return KafkaConsumer{}, err
	}

	return KafkaConsumer{
		consumer: consumer,
		wg:       &sync.WaitGroup{},
	}, nil
}

func (h KafkaConsumer) Consume(topic string, partition int) (Consumer, error) {
	return h.consumer.Consume(topic, partition)
}

// ConsumeAndHandle consumes and handles messages from kafka topic and partitions.
func (h KafkaConsumer) ConsumeAndHandle(ctx context.Context, handler EventHandler, topic string, partitions []int) {
	if len(partitions) == 0 {
		return
	}

	for _, p := range partitions {
		pconsumer, err := h.consumer.Consume(topic, p)
		if err != nil {
			ll.Fatal("Unable to consume messages", l.String("topic", topic), l.Int("partition", p))
		}

		ll.Info("Start consuming messages", l.String("topic", topic), l.Int("partition", p))
		h.wg.Add(1)
		go func() {
			defer h.wg.Done()
			err := pconsumer.ConsumeAndHandle(ctx, handler)
			if err != nil {
				ll.Error("Error", l.Error(err))
			}
		}()
	}
}

// Wait waits for all goroutines to be finished.
func (h KafkaConsumer) Wait() {
	h.wg.Wait()
}
