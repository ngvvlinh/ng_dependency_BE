package handler

import (
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/backend/com/eventhandler/pgevent"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/mq"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/common/l"
)

var ll = l.New()

type Handler struct {
	historyStore historysqlstore.HistoryStoreFactory
	producer     *mq.KafkaProducer
	prefix       string

	fbuserQuery      fbusering.QueryBus
	fbMessagingQuery fbmessaging.QueryBus
	fbPagingQuery    fbpaging.QueryBus
	indentityQuery   identity.QueryBus
}

func New(
	db com.MainDB,
	producer *mq.KafkaProducer,
	prefix string,
	fbuserQ fbusering.QueryBus,
	fbMessagingQ fbmessaging.QueryBus,
	fbPageQ fbpaging.QueryBus,
	indentityQuerybus identity.QueryBus,
) *Handler {
	h := &Handler{
		historyStore:     historysqlstore.NewHistoryStore(db),
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
