package handler

import (
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/com/fabo/pkg/fbclient"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/common/l"
)

var ll = l.New()

type Handler struct {
	db               *cmsql.Database
	historyStore     historysqlstore.HistoryStoreFactory
	producer         *mq.KafkaProducer
	prefix           string
	fbuserQuery      fbusering.QueryBus
	fbMessagingQuery fbmessaging.QueryBus
	fbPagingQuery    fbpaging.QueryBus
	indentityQuery   identity.QueryBus
	shippingQuery    shipping.QueryBus
	orderQuery       ordering.QueryBus
	fbClient         *fbclient.FbClient
}

func New(
	db com.MainDB,
	producer *mq.KafkaProducer,
	prefix string,
	fbuserQ fbusering.QueryBus,
	fbMessagingQ fbmessaging.QueryBus,
	fbPageQ fbpaging.QueryBus,
	indentityQuerybus identity.QueryBus,
	shippingQ shipping.QueryBus,
	orderQ ordering.QueryBus,
	fbClient *fbclient.FbClient,
) *Handler {
	h := &Handler{
		db:               db,
		historyStore:     historysqlstore.NewHistoryStore(db),
		producer:         producer,
		prefix:           prefix + "_pgrid_",
		fbuserQuery:      fbuserQ,
		fbMessagingQuery: fbMessagingQ,
		fbPagingQuery:    fbPageQ,
		indentityQuery:   indentityQuerybus,
		shippingQuery:    shippingQ,
		orderQuery:       orderQ,
		fbClient:         fbClient,
	}
	return h
}

func (h *Handler) TopicsAndHandlers() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		// turn off feature send message shipping_state changed
		//"fulfillment":              h.HandleFulfillmentEvent,
		"fb_external_conversation": h.HandleFbConversationEvent,
		"fb_external_comment":      h.HandleFbCommentEvent,
		"fb_external_message":      h.HandleFbMessageEvent,
		"fb_customer_conversation": h.HandleFbCustomerConversationEvent,
	})
}
