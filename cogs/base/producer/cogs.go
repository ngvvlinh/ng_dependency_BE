package _producer

import (
	"context"

	"o.o/backend/com/eventhandler/handler/intctl"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/webhook"
)

func SupportedProducers(ctx context.Context, cfg cc.Kafka) (_ webhook.Producer, _ error) {
	var producer *mq.KafkaProducer
	if !cfg.Enabled {
		return nil, nil
	}
	producer, err := mq.NewKafkaProducer(ctx, cfg.Brokers)
	if err != nil {
		return nil, err
	}
	ctlProducer := producer.WithTopic(intctl.Topic(cfg.TopicPrefix))
	return ctlProducer, err
}
