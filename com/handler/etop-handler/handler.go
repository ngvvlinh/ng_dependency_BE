package handler

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/api/top/external/types"
	"o.o/backend/com/handler/etop-handler/intctl"
	"o.o/backend/com/handler/etop-handler/pgrid"
	"o.o/backend/com/handler/etop-handler/webhook/sender"
	"o.o/backend/com/handler/pgevent"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

const ConsumerGroup = "handler/webhook"

var (
	ll        = l.New()
	publisher eventstream.Publisher
)

type Handler struct {
	db           *cmsql.Database
	historyStore historysqlstore.HistoryStoreFactory

	consumer     mq.KafkaConsumer
	producer     *mq.KafkaProducer
	handlers     map[string]pgrid.HandlerFunc
	handlersfabo map[string]pgrid.HandlerFuncFabo
	prefix       string
	wg           sync.WaitGroup

	sender *sender.WebhookSender

	catalogQuery     catalog.QueryBus
	customerQuery    customering.QueryBus
	inventoryQuery   inventory.QueryBus
	addressQuery     addressing.QueryBus
	locationQuery    location.QueryBus
	fbuserQuery      fbusering.QueryBus
	fbMessagingQuery fbmessaging.QueryBus
	fbPagingQuery    fbpaging.QueryBus
	fbQuery          fbpaging.QueryBus
	indentityQuery   identity.QueryBus
}

func New(
	db *cmsql.Database, sender *sender.WebhookSender,
	consumer mq.KafkaConsumer,
	prefix string, catalogQ catalog.QueryBus,
	customerQ customering.QueryBus, inventoryQ inventory.QueryBus,
	addressQ addressing.QueryBus, locationQ location.QueryBus,
	fbuserQ fbusering.QueryBus, producer *mq.KafkaProducer,
	fbMessagingQ fbmessaging.QueryBus, fbPageQ fbpaging.QueryBus,
	indentityQuerybus identity.QueryBus,

) *Handler {
	h := &Handler{
		db:               db,
		historyStore:     historysqlstore.NewHistoryStore(db),
		consumer:         consumer,
		producer:         producer,
		prefix:           prefix + "_pgrid_",
		sender:           sender,
		catalogQuery:     catalogQ,
		customerQuery:    customerQ,
		inventoryQuery:   inventoryQ,
		addressQuery:     addressQ,
		locationQuery:    locationQ,
		fbuserQuery:      fbuserQ,
		fbMessagingQuery: fbMessagingQ,
		fbPagingQuery:    fbPageQ,
		indentityQuery:   indentityQuerybus,
	}
	handlers := h.TopicsAndHandlers()
	h.handlers = handlers
	if len(handlers) != len(pgevent.Topics) {
		ll.Panic("Handler list mismatch")
	}
	for _, d := range pgevent.Topics {
		_, ok := handlers[d.Name]
		if !ok {
			ll.Panic("Handler not found", l.String("name", d.Name))
		}
	}
	return h
}

func NewHandlerFabo(
	db *cmsql.Database, prefix string, fbuserQ fbusering.QueryBus,
	consumer mq.KafkaConsumer, publisherEvent eventstream.Publisher,
	fbMessagingQ fbmessaging.QueryBus, fbPageQ fbpaging.QueryBus,
	indentityQuerybus identity.QueryBus, fbQuery fbpaging.QueryBus,
) *Handler {
	h := &Handler{
		db:               db,
		historyStore:     historysqlstore.NewHistoryStore(db),
		consumer:         consumer,
		prefix:           prefix + "_pgrid_",
		fbuserQuery:      fbuserQ,
		fbMessagingQuery: fbMessagingQ,
		fbPagingQuery:    fbPageQ,
		fbQuery:          fbQuery,
		indentityQuery:   indentityQuerybus,
	}
	publisher = publisherEvent
	handlers := h.TopicsFaboAndHandlers()
	h.handlersfabo = handlers
	return h
}

func NewWithHandlers(db *cmsql.Database, sender *sender.WebhookSender, consumer mq.KafkaConsumer, prefix string, handlers map[string]pgrid.HandlerFunc) *Handler {
	if len(handlers) == 0 {
		ll.Panic("Missing handler!")
	}
	h := &Handler{
		db:       db,
		consumer: consumer,
		prefix:   prefix + "_pgrid_",
		sender:   sender,
		handlers: handlers,
	}
	return h
}

func (h *Handler) TopicsAndHandlers() map[string]pgrid.HandlerFunc {
	return map[string]pgrid.HandlerFunc{
		"fulfillment":                   h.HandleFulfillmentEvent,
		"order":                         h.HandleOrderEvent,
		"notification":                  nil,
		"money_transaction_shipping":    nil,
		"shop_product":                  h.HandleShopProductEvent,
		"shop_variant":                  h.HandleShopVariantEvent,
		"shop_customer":                 h.HandleShopCustomerEvent,
		"shop_customer_group":           h.HandleShopCustomerGroupEvent,
		"shop_customer_group_customer":  h.HandleShopCustomerGroupCustomerEvent,
		"inventory_variant":             h.HandleInventoryVariantEvent,
		"shop_trader_address":           h.HandleShopTraderAddressEvent,
		"shop_collection":               h.HandleShopProductCollectionEvent,
		"shop_product_collection":       h.HandleShopProductionCollectionRelationshipEvent,
		"fb_external_conversation":      h.HandleFbConversationEvent,
		"fb_external_comment":           h.HandleFbCommentEvent,
		"fb_external_message":           h.HandleFbMessageEvent,
		"fb_customer_conversation":      h.HandleFbCustomerConversationEvent,
		"fb_external_conversation_fabo": nil,
		"fb_external_comment_fabo":      nil,
		"fb_external_message_fabo":      nil,
		"fb_customer_conversation_fabo": nil,
	}
}

func (h *Handler) TopicsFaboAndHandlers() map[string]pgrid.HandlerFuncFabo {
	return map[string]pgrid.HandlerFuncFabo{
		"fb_external_conversation_fabo": h.HandleFbConversationFaboEvent,
		"fb_external_comment_fabo":      h.HandleFbCommentFaboEvent,
		"fb_external_message_fabo":      h.HandleFbMessageFaboEvent,
		"fb_customer_conversation_fabo": h.HandleFbCustomerConversationFaboEvent,
	}
}

func (h *Handler) RegisterTo(intctlHandler *intctl.Handler) {
	intctlHandler.Subscribe(intctl.ChannelReloadWebhook, h.handleReloadWebhook)
}

func (h *Handler) ConsumeAndHandleAllTopics(ctx context.Context) {
	for _, d := range pgevent.Topics {
		if strings.Contains(d.Name, "fabo") {
			continue
		}
		h.comsumerAndHandlerTopics(ctx, d)
	}
}

func (h *Handler) getWrapHandlerFunc(topic string) mq.EventHandler {
	var result mq.EventHandler
	if strings.Contains(topic, "fabo") {
		handler := h.handlersfabo[topic]
		if handler == nil {
			ll.Info("No handler for topic", l.String("topic", topic))
			return result
		}
		result = pgrid.WrapHandlerFuncFabo(handler)
	} else {
		handler := h.handlers[topic]
		if handler == nil {
			ll.Info("No handler for topic", l.String("topic", topic))
			return result
		}
		result = pgrid.WrapHandlerFunc(handler)
	}
	return result
}

func (h *Handler) comsumerAndHandlerTopics(ctx context.Context, d pgevent.TopicDef) {
	count := 0
	var m sync.Mutex
	topics := make(map[string]int)
	wrappedHandler := h.getWrapHandlerFunc(d.Name)
	var wg sync.WaitGroup
	kafkaTopic := h.prefix + d.Name
	wg.Add(d.Partitions)
	h.wg.Add(d.Partitions)

	for i := 0; i < d.Partitions; i++ {
		partition := i
		go func() {
			pc, err := h.consumer.Consume(kafkaTopic, partition)
			if err != nil {
				ll.S.Fatalf("Error while consuming topic: %v:%v", kafkaTopic, partition)
				return
			}
			defer h.wg.Done()
			defer func() { _ = pc.Close() }()

			wg.Done()
			m.Lock()
			count++
			topics[kafkaTopic]++
			m.Unlock()

			err = pc.ConsumeAndHandle(ctx, wrappedHandler)
			if err != nil {
				ll.S.Errorf("Handler for topic %v:%v stopped: %+v", kafkaTopic, partition, err)
				buf := make([]byte, 2048)
				runtime.Stack(buf, false)
				msg := fmt.Sprintf(
					"🔥 Handler for topic %v:%v stoppped: %+v\n\n%s",
					kafkaTopic, partition, err, buf)
				ll.SendMessage(msg)
			}
		}()
	}
	wg.Wait()
	m.Lock()
	defer m.Unlock()
}

func (h *Handler) ConsumerAndHandlerFaboTopic(ctx context.Context) {
	for _, d := range pgevent.Topics {
		if !strings.Contains(d.Name, "fabo") {
			continue
		}
		h.comsumerAndHandlerTopics(ctx, d)
	}
}

func (h *Handler) Wait() {
	h.wg.Wait()
}

func (h *Handler) handleReloadWebhook(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
	var v intctl.ReloadWebhook
	err := jsonx.Unmarshal(msg.Value, &v)
	if err != nil {
		return mq.CodeStop, nil
	}
	if v.AccountID == 0 {
		ll.Error("webhook/reload: account_id is empty")
		return mq.CodeStop, nil
	}
	if err := h.sender.Reload(ctx, v.AccountID); err != nil {
		ll.Error("webhook/reload: account_id is empty")
		return mq.CodeStop, nil
	}
	return mq.CodeOK, nil
}

// TODO: handle soft delete

func (h *Handler) HandleOrderEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	var history ordermodel.OrderHistory
	if ok, err := h.db.Where("rid = ?", event.RID).Get(&history); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbOrderHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.ID("order_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	var order ordermodel.Order
	if ok, err := h.db.Where("id = ?", id).Get(&order); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}

	convertpb.PbOrderHistory(history)
	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Order: convertpb.PbOrder(&order),
	}
	change.Changed = &types.ChangeOneOf{
		Order: changed,
	}
	accountIDs := []dot.ID{order.ShopID, order.PartnerID}
	return h.sender.CollectPb(ctx, event.Table, id, order.ShopID, accountIDs, change)
}

func (h *Handler) HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFulfillmentEvent", l.Object("pgevent", event))
	var history shipmodel.FulfillmentHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbFulfillmentHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.ID("fulfillment_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	var ffm shipmodel.Fulfillment
	if ok, err := h.db.Where("id = ?", id).Get(&ffm); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}

	convertpb.PbFulfillmentHistory(history)
	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Fulfillment: convertpb.PbFulfillment(&ffm),
	}
	change.Changed = &types.ChangeOneOf{
		Fulfillment: changed,
	}
	accountIDs := []dot.ID{ffm.ShopID, ffm.PartnerID}
	return h.sender.CollectPb(ctx, event.Table, id, ffm.ShopID, accountIDs, change)
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
