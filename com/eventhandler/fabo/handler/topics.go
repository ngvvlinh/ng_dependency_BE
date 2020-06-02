package handler

import (
	"o.o/backend/com/eventhandler"
	"o.o/backend/pkg/etop/model"
)

func Topics() []eventhandler.TopicDef {
	return []eventhandler.TopicDef{
		{
			Name:       "fb_external_conversation",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "fb_external_message",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "fb_external_comment",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "fb_customer_conversation",
			Partitions: 64,
			DBName:     model.DBMain,
		},
	}
}

var mapTopics = eventhandler.MapTopics(Topics())
