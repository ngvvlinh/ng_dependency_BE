package eventhandler

import (
	"o.o/backend/pkg/etop/model"
)

type TopicDef struct {
	Name       string
	Partitions int
	DBName     model.DBName
}

func MapTopics(topics []TopicDef) map[string]TopicDef {
	m := make(map[string]TopicDef)
	for _, topic := range topics {
		m[topic.Name] = topic
	}
	return m
}
