package eventhandler

import (
	"o.o/backend/pkg/etc/dbdecl"
)

type TopicDef struct {
	Name       string
	Partitions int
	DBName     dbdecl.DBName
}

func MapTopics(topics []TopicDef) map[string]TopicDef {
	m := make(map[string]TopicDef)
	for _, topic := range topics {
		m[topic.Name] = topic
	}
	return m
}
