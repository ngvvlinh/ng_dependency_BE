package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Shopify/sarama"

	pbexternal "etop.vn/api/pb/external"
	"etop.vn/backend/com/handler/etop-handler/intctl"
	"etop.vn/backend/com/handler/etop-handler/pgrid"
	"etop.vn/backend/com/handler/etop-handler/webhook/sender"
	"etop.vn/backend/com/handler/pgevent"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/telebot"
	historysqlstore "etop.vn/backend/pkg/etop-history/sqlstore"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

const ConsumerGroup = "handler/webhook"

var ll = l.New()

type Handler struct {
	db           *cmsql.Database
	historyStore historysqlstore.HistoryStoreFactory
	bot          *telebot.Channel

	consumer mq.KafkaConsumer
	handlers map[string]pgrid.HandlerFunc
	prefix   string
	wg       sync.WaitGroup

	sender *sender.WebhookSender
}

func New(db *cmsql.Database, sender *sender.WebhookSender, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string) *Handler {
	h := &Handler{
		db:           db,
		historyStore: historysqlstore.NewHistoryStore(db),
		bot:          bot,
		consumer:     consumer,
		prefix:       prefix + "_pgrid_",
		sender:       sender,
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

func NewWithHandlers(db *cmsql.Database, sender *sender.WebhookSender, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string, handlers map[string]pgrid.HandlerFunc) *Handler {
	if len(handlers) == 0 {
		ll.Panic("Missing handler!")
	}
	h := &Handler{
		db:       db,
		bot:      bot,
		consumer: consumer,
		prefix:   prefix + "_pgrid_",
		sender:   sender,
		handlers: handlers,
	}
	return h
}

func (h *Handler) TopicsAndHandlers() map[string]pgrid.HandlerFunc {
	return map[string]pgrid.HandlerFunc{
		"fulfillment":                h.HandleFulfillmentEvent,
		"order":                      h.HandleOrderEvent,
		"notification":               nil,
		"money_transaction_shipping": nil,
	}
}

func (h *Handler) RegisterTo(intctlHandler *intctl.Handler) {
	intctlHandler.Subscribe(intctl.ChannelReloadWebhook, h.handleReloadWebhook)
}

func (h *Handler) ConsumeAndHandleAllTopics(ctx context.Context) {
	count := 0
	var m sync.Mutex
	topics := make(map[string]int)

	var wg sync.WaitGroup
	for _, d := range pgevent.Topics {
		handler := h.handlers[d.Name]
		if handler == nil {
			ll.Info("No handler for topic", l.String("topic", d.Name))
			continue
		}

		wrappedHandler := pgrid.WrapHandlerFunc(handler)
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
					if h.bot != nil {
						msg := fmt.Sprintf(
							"🔥 Handler for topic %v:%v stoppped: %+v",
							kafkaTopic, partition, err)
						h.bot.SendMessage(msg)
					}
				}
			}()
		}
	}

	wg.Wait()
	m.Lock()
	defer m.Unlock()

	ll.S.Infof("Initialized %v consumers", count)
	for topic, partitions := range topics {
		ll.S.Infof("%2v %v\n", partitions, topic)
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
	change.Latest = &pbexternal.LatestOneOf{
		Order: convertpb.PbOrder(&order),
	}
	change.Changed = &pbexternal.ChangeOneOf{
		Order: changed,
	}
	accountIDs := []dot.ID{order.ShopID, order.PartnerID}
	return h.sender.CollectPb(ctx, event.Table, id, accountIDs, change)
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
	change.Latest = &pbexternal.LatestOneOf{
		Fulfillment: convertpb.PbFulfillment(&ffm),
	}
	change.Changed = &pbexternal.ChangeOneOf{
		Fulfillment: changed,
	}
	accountIDs := []dot.ID{ffm.ShopID, ffm.PartnerID}
	return h.sender.CollectPb(ctx, event.Table, id, accountIDs, change)
}

func pbChange(event *pgevent.PgEvent) *pbexternal.Change {
	return &pbexternal.Change{
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
