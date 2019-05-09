package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/etop-handler/intctl"
	"etop.vn/backend/pkg/etop-handler/pgrid"
	"etop.vn/backend/pkg/etop-handler/webhook/sender"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/pgevent"
	"github.com/Shopify/sarama"

	pbcm "etop.vn/backend/pb/common"
	pbexternal "etop.vn/backend/pb/external"
)

const ConsumerGroup = "handler/webhook"

var ll = l.New()

type Handler struct {
	db  cmsql.Database
	bot *telebot.Channel

	consumer mq.KafkaConsumer
	handlers map[string]pgrid.HandlerFunc
	prefix   string
	wg       sync.WaitGroup

	sender *sender.WebhookSender
}

func New(db cmsql.Database, sender *sender.WebhookSender, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string) *Handler {
	h := &Handler{
		db:       db,
		bot:      bot,
		consumer: consumer,
		prefix:   prefix + "_pgrid_",
		sender:   sender,
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

func NewWithHandlers(db cmsql.Database, sender *sender.WebhookSender, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string, handlers map[string]pgrid.HandlerFunc) *Handler {
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
		"account":                          nil,
		"address":                          nil,
		"etop_category":                    nil,
		"fulfillment":                      h.HandleFulfillmentEvent,
		"order_external":                   nil,
		"order_source":                     nil,
		"order":                            h.HandleOrderEvent,
		"product_brand":                    nil,
		"product_external":                 nil,
		"product_source_category_external": nil,
		"product_source_category":          nil,
		"product_source":                   nil,
		"product":                          nil,
		"shop_collection":                  nil,
		"shop_product":                     nil,
		"shop":                             nil,
		"supplier":                         nil,
		"user":                             nil,
		"variant_external":                 nil,
		"variant":                          nil,
		"money_transaction_shipping":       nil,
		"notification":                     nil,
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
			partition := int32(i)
			go func() {
				pc, err := h.consumer.Consume(kafkaTopic, partition)
				if err != nil {
					ll.S.Fatalf("Error while consuming topic: %v:%v", kafkaTopic, partition)
				}
				defer h.wg.Done()
				defer pc.Close()

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
		fmt.Printf("\t%2v %v\n", partitions, topic)
	}
}

func (h *Handler) Wait() {
	h.wg.Wait()
}

func (h *Handler) handleReloadWebhook(ctx context.Context, msg *sarama.ConsumerMessage) (mq.Code, error) {
	var v intctl.ReloadWebhook
	err := json.Unmarshal(msg.Value, &v)
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
	var history model.OrderHistory
	if ok, err := h.db.Where("rid = ?", event.RID).Get(&history); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := pbexternal.PbOrderHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.Int64("order_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := *history.ID().Int64()
	var order model.Order
	if ok, err := h.db.Where("id = ?", id).Get(&order); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID), l.Int64("id", id))
		return mq.CodeIgnore, nil
	}

	pbexternal.PbOrderHistory(history)
	change := pbChange(event)
	change.Latest = &pbexternal.LatestOneOf{
		Latest: &pbexternal.LatestOneOf_Order{pbexternal.PbOrder(&order)},
	}
	change.Changed = &pbexternal.ChangeOneOf{
		Changed: &pbexternal.ChangeOneOf_Order{changed},
	}
	accountIDs := []int64{order.ShopID, order.PartnerID}
	return h.sender.CollectPb(ctx, event.Table, id, accountIDs, change)
}

func (h *Handler) HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFulfillmentEvent", l.Object("pgevent", event))
	var history model.FulfillmentHistory
	if ok, err := h.db.Where("rid = ?", event.RID).Get(&history); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := pbexternal.PbFulfillmentHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.Int64("fulfillment_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := *history.ID().Int64()
	var ffm model.Fulfillment
	if ok, err := h.db.Where("id = ?", id).Get(&ffm); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID), l.Int64("id", id))
		return mq.CodeIgnore, nil
	}

	pbexternal.PbFulfillmentHistory(history)
	change := pbChange(event)
	change.Latest = &pbexternal.LatestOneOf{
		Latest: &pbexternal.LatestOneOf_Fulfillment{pbexternal.PbFulfillment(&ffm)},
	}
	change.Changed = &pbexternal.ChangeOneOf{
		Changed: &pbexternal.ChangeOneOf_Fulfillment{changed},
	}
	accountIDs := []int64{ffm.ShopID, ffm.PartnerID}
	return h.sender.CollectPb(ctx, event.Table, id, accountIDs, change)
}

func pbChange(event *pgevent.PgEvent) *pbexternal.Change {
	return &pbexternal.Change{
		Time:       pbcm.PbTime(time.Unix(event.Timestamp, 0)),
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
