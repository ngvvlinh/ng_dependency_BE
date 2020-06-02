package publisher

import (
	"context"

	"o.o/backend/com/eventhandler"
	"o.o/backend/com/eventhandler/fabo/types"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/model"
)

type HandlerFunc func(context.Context, *types.FaboEvent) (mq.Code, error)

var Topics = []eventhandler.TopicDef{
	{
		Name:       "fb_external_conversation_fabo",
		Partitions: 64,
		DBName:     model.DBMain,
	},
	{
		Name:       "fb_external_message_fabo",
		Partitions: 64,
		DBName:     model.DBMain,
	},
	{
		Name:       "fb_external_comment_fabo",
		Partitions: 64,
		DBName:     model.DBMain,
	},
	{
		Name:       "fb_customer_conversation_fabo",
		Partitions: 64,
		DBName:     model.DBMain,
	},
}
