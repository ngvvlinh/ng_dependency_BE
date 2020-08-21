package publisher

import (
	"context"

	"o.o/backend/com/eventhandler"
	"o.o/backend/com/eventhandler/fabo/types"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etc/dbdecl"
)

type HandlerFunc func(context.Context, *types.FaboEvent) (mq.Code, error)

func Topics() []eventhandler.TopicDef {
	return []eventhandler.TopicDef{
		{
			Name:       "fb_external_conversation_fabo",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "fb_external_message_fabo",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "fb_external_comment_fabo",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "fb_customer_conversation_fabo",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
	}
}

var mapTopics = eventhandler.MapTopics(Topics())

func GetTopics(topics map[string]mq.EventHandler) []eventhandler.TopicDef {
	var result []eventhandler.TopicDef
	for name := range topics {
		result = append(result, mapTopics[name])
	}
	return result
}
