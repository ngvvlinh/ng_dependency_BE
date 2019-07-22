package mq

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Shopify/sarama"

	"etop.vn/common/l"
)

// EventHandler ...
type EventHandler func(context.Context, *sarama.ConsumerMessage) (Code, error)

// Consumer ...
type Consumer interface {
	Errors() <-chan *sarama.ConsumerError
	Messages() <-chan *sarama.ConsumerMessage
	Ack(*sarama.ConsumerMessage)
	Close() error

	ConsumeAndHandle(context.Context, EventHandler) error
}

// SingleConsumer ...
type SingleConsumer struct {
	client   sarama.Client
	consumer sarama.Consumer
	om       sarama.OffsetManager
}

// NewSingleConsumer ...
func NewSingleConsumer(brokers []string, config *sarama.Config, group string) (*SingleConsumer, error) {
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	om, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		return nil, err
	}

	return &SingleConsumer{
		client:   client,
		consumer: consumer,
		om:       om,
	}, nil
}

// Consume ...
func (c *SingleConsumer) Consume(topic string, partition int32) (Consumer, error) {
	pom, err := c.om.ManagePartition(topic, partition)
	if err != nil {
		ll.Error("Cannot retrieve partition offset manager for this topic", l.String("topic", topic), l.Int32("partition", partition), l.Error(err))
		return nil, err
	}

	nx, _ := pom.NextOffset()

	retry := 0
retry:
	pConsumer, err := c.consumer.ConsumePartition(topic, partition, nx)
	if err != nil {
		ll.Error("Unable to init partition consumer", l.String("topic", topic), l.Int32("partition", partition), l.Error(err))
		if err == sarama.ErrOffsetOutOfRange && retry < 1 {
			ll.Info("Retry: Start consuming from the newest offset", l.String("topic", topic), l.Int32("partition", partition))
			nx = sarama.OffsetNewest
			retry++
			goto retry
		}
		return nil, err
	}
	ll.Info("Start kafka consumer", l.Int64("offset", nx), l.String("topic", topic), l.Int32("partition", partition))

	return &kafkaConsumer{
		PartitionConsumer: pConsumer,

		pom:       pom,
		topic:     topic,
		partition: partition,
	}, nil
}

// Close ...
func (c *SingleConsumer) Close() (err error) {
	err1 := c.om.Close()
	err2 := c.consumer.Close()
	err3 := c.client.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return err3
}

type kafkaConsumer struct {
	sarama.PartitionConsumer
	pom sarama.PartitionOffsetManager

	topic     string
	partition int32
}

func (c *kafkaConsumer) Ack(msg *sarama.ConsumerMessage) {
	c.pom.MarkOffset(msg.Offset+1, "")
}

func (c *kafkaConsumer) Close() error {
	err1 := c.PartitionConsumer.Close()
	err2 := c.pom.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// ConsumeAndHandle ...
func (c *kafkaConsumer) ConsumeAndHandle(ctx context.Context, handler EventHandler) (_err error) {
	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-c.Errors():
			if err != nil {
				ll.Error("Received Kafka error", l.Object("ConsumerError", err), l.String("topic", c.topic), l.Int32("partition", c.partition))
				// TODO: Send to bot
			}

		case msg := <-c.Messages():
			if msg == nil {
				ll.Warn("Received nil message (the channel has been closed)", l.String("topic", c.topic), l.Int32("partition", c.partition))
				time.Sleep(100 * time.Millisecond)
				continue
			}
			ll.Info("Received message from Kafka",
				l.Int64("offset", msg.Offset),
				l.String("topic", msg.Topic),
				l.Int32("partition", msg.Partition),
				l.String("key", string(msg.Key)))

			// Stop listening
			if err := c.handleMessage(ctx, handler, msg); err != nil {
				return err
			}
		}
	}
}

func (c *kafkaConsumer) handleMessage(ctx context.Context, handler EventHandler, msg *sarama.ConsumerMessage) (_err error) {
	defer func() {
		e := recover()
		if e != nil {
			ll.Error("RECOVER", l.Object("error", e))
			if err, ok := e.(error); ok {
				_err = err
				fmt.Printf("%+v", _err)
			} else {
				_err = fmt.Errorf("Panic: %v", e)

			}
			stack := debug.Stack()
			fmt.Printf("%s", stack)

			// TODO: Send to bot
		}
	}()

	count := 0

retry:
	code, err := handler(ctx, msg)
	if err == nil && code == CodeOK {
		c.Ack(msg)
		return nil
	}

	switch code {
	case CodeIgnore:
		ll.Warn("Error while handling message (ignored)", l.Error(err))
		return nil

	case CodeRetry:
		count++
		if count >= 3 {
			ll.S.Errorf("Unable to handle message (retried %v)", count)
			return err
		}
		ll.S.Warn("Error while handling message (retry %v)", l.Error(err))
		goto retry

	default:
		ll.Error("Unable to handle message (stop)", l.Error(err))
		return err
	}
}
