package handler

import (
	"context"

	"etop.vn/api/external/haravan"
	handler "etop.vn/backend/com/handler/etop-handler"
	"etop.vn/backend/com/handler/etop-handler/pgrid"
	"etop.vn/backend/com/handler/pgevent"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/telebot"
	historysqlstore "etop.vn/backend/pkg/etop-history/sqlstore"
	serviceHaravan "etop.vn/backend/pkg/external/haravan"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
)

var (
	historyStore historysqlstore.HistoryStoreFactory
	haravanAggr  haravan.CommandBus
)

const ConsumerGroup = "handler/haravan"

func New(db cmsql.Database, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string) *handler.Handler {
	historyStore = historysqlstore.NewHistoryStore(db)
	haravanAggr = serviceHaravan.NewAggregate(db, haravanclient.Config{}).MessageBus()
	handlers := TopicsAndHandlersHaravan()
	h := handler.NewWithHandlers(db, nil, bot, consumer, prefix, handlers)
	return h
}

func TopicsAndHandlersHaravan() map[string]pgrid.HandlerFunc {
	return map[string]pgrid.HandlerFunc{
		"fulfillment": HandleFulfillmentEvent,
	}
}

func HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	// Update shipping state or payment status if this ffm comes from Haravan
	var ffmHistory shipmodel.FulfillmentHistory
	if ok, err := historyStore(ctx).GetHistory(&ffmHistory, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		return mq.CodeIgnore, nil
	}
	id := *ffmHistory.ID().Int64()
	if ffmHistory.ShippingState().String() != nil {
		cmd := &haravan.SendUpdateExternalFulfillmentStateCommand{
			FulfillmentID: id,
		}
		// Ignore err
		_ = haravanAggr.Dispatch(ctx, cmd)
	}
	if ffmHistory.EtopPaymentStatus().Int() != nil {
		cmd := &haravan.SendUpdateExternalPaymentStatusCommand{
			FulfillmentID: id,
		}
		// Ignore err
		_ = haravanAggr.Dispatch(ctx, cmd)
	}

	return mq.CodeOK, nil
}
