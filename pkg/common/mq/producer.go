package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"

	"etop.vn/common/l"
)

var ll = l.New()

// Producer ...
type Producer interface {
	Send(partition int32, id string, data []byte)
	SendJSON(partition int32, id string, v interface{})
}

// KafkaProducer ...
type KafkaProducer struct {
	producer sarama.AsyncProducer
}

// NewKafkaProducer ...
func NewKafkaProducer(ctx context.Context, brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy    // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond  // Flush batches every 500ms
	config.Producer.Partitioner = sarama.NewManualPartitioner // Manually assign partition

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	mq := &KafkaProducer{producer}
	ll.Info("Connected to Kafka. Start goroutine for handling producer error.")
	go mq.handleKafkaErrors()

	go func() {
		select {
		case <-ctx.Done():
			ll.Info("Closing broadcast producer")
			producer.AsyncClose()
		}
	}()

	return mq, nil
}

func (mq *KafkaProducer) handleKafkaErrors() {
	for err := range mq.producer.Errors() {
		ll.Error("Send message failed", l.Error(err.Err))
	}
}

// WithTopic ...
func (mq *KafkaProducer) WithTopic(topic string) Producer {
	return kafkaProducer{mq, topic}
}

func (mq *KafkaProducer) SendJSON(topic string, partition int32, id string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ll.Panic("error marshalling", l.Error(err))
	}
	mq.Send(topic, partition, id, data)
}

func (mq *KafkaProducer) Send(topic string, partition int32, id string, data []byte) {
	if mq == nil {
		ll.Warn("Send event to Kafka (skipped)", l.String("topic", topic), l.Int32("p", partition), l.String("key", id))
		return
	}

	pmsg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(id),
		Value:     sarama.ByteEncoder(data),
		Partition: partition,
	}
	ll.Debug("Send event to Kafka", l.String("topic", topic), l.Int32("p", partition), l.String("key", id))
	mq.producer.Input() <- pmsg
}

type kafkaProducer struct {
	*KafkaProducer
	topic string
}

func (mq kafkaProducer) SendJSON(partition int32, id string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		ll.Panic("error marshalling", l.Error(err))
	}
	mq.Send(partition, id, data)
}

// Send ...
func (mq kafkaProducer) Send(partition int32, id string, data []byte) {
	pmsg := &sarama.ProducerMessage{
		Topic:     mq.topic,
		Key:       sarama.StringEncoder(id),
		Value:     sarama.ByteEncoder(data),
		Partition: partition,
	}
	ll.Debug("Send event to Kafka", l.String("topic", mq.topic), l.String("key", id))
	mq.producer.Input() <- pmsg
}
