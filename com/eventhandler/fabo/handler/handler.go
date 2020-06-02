package handler

import (
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/common/l"
)

var ll = l.New()

type Handler struct {
	db           *cmsql.Database
	historyStore historysqlstore.HistoryStoreFactory
	consumer     mq.KafkaConsumer
	producer     *mq.KafkaProducer
	prefix       string

	fbuserQuery      fbusering.QueryBus
	fbMessagingQuery fbmessaging.QueryBus
	fbPagingQuery    fbpaging.QueryBus
	indentityQuery   identity.QueryBus
}

func New(
	db *cmsql.Database,
	consumer mq.KafkaConsumer,
	producer *mq.KafkaProducer,
	prefix string,
	fbuserQ fbusering.QueryBus,
	fbMessagingQ fbmessaging.QueryBus,
	fbPageQ fbpaging.QueryBus,
	indentityQuerybus identity.QueryBus,
) *Handler {
	h := &Handler{
		db:               db,
		historyStore:     historysqlstore.NewHistoryStore(db),
		consumer:         consumer,
		producer:         producer,
		prefix:           prefix + "_pgrid_",
		fbuserQuery:      fbuserQ,
		fbMessagingQuery: fbMessagingQ,
		fbPagingQuery:    fbPageQ,
		indentityQuery:   indentityQuerybus,
	}
	return h
}

func (h *Handler) TopicsAndHandlers() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		"fb_external_conversation": h.HandleFbConversationEvent,
		"fb_external_comment":      h.HandleFbCommentEvent,
		"fb_external_message":      h.HandleFbMessageEvent,
		"fb_customer_conversation": h.HandleFbCustomerConversationEvent,
	})
}
