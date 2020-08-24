package handler

import (
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/api/top/external/types"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/com/eventhandler/webhook/sender"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/apifw/cmapi"
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
	sender           *sender.WebhookSender
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
		db:               db,
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

func pbChange(event *pgevent.PgEvent) *types.Change {
	return &types.Change{
		Time:       cmapi.PbTime(time.Unix(event.Timestamp, 0)),
		ChangeType: pbChangeType(event.Op),
		Entity:     event.Table,
	}
}

func pbChangeType(op pgevent.TGOP) string {
	switch op {
	case pgevent.OpInsert:
		return "create"
	case pgevent.OpUpdate:
		return "update"
	case pgevent.OpDelete:
		return "delete"
	default:
		return ""
	}
}
