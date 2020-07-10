package handler

import (
	"context"

	handler "o.o/backend/com/eventhandler/handler"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/notifier/sqlstore"
	"o.o/backend/com/eventhandler/pgevent"
	com "o.o/backend/com/main"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/common/l"
)

var (
	x            *cmsql.Database
	xNotifier    *cmsql.Database
	ll           = l.New()
	notiStore    *sqlstore.NotificationStore
	deviceStore  *sqlstore.DeviceStore
	historyStore historysqlstore.HistoryStoreFactory
)

const ConsumerGroup = "handler/notifier"

func New(dbMain com.MainDB, dbNotifier com.NotifierDB, consumer mq.KafkaConsumer, cfg cc.Kafka) (handlerMain *handler.Handler, handlerNotifier *handler.Handler) {
	x = dbMain
	xNotifier = dbNotifier
	notiStore = sqlstore.NewNotificationStore(dbNotifier)
	deviceStore = sqlstore.NewDeviceStore(dbNotifier)
	historyStore = historysqlstore.NewHistoryStore(dbMain)

	handlerMain = handler.New(consumer, cfg)
	handlerNotifier = handler.New(consumer, cfg)
	return handlerMain, handlerNotifier
}

func TopicsAndHandlersEtop() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		"fulfillment":                HandleFulfillmentEvent,
		"money_transaction_shipping": HandleMoneyTransactionShippingEvent,
	})
}

func TopicsAndHandlerNotifier() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		"notification": HandleNotificationEvent,
	})
}

func CreateNotifications(ctx context.Context, cmds []*notifiermodel.CreateNotificationArgs) error {
	if len(cmds) == 0 {
		return nil
	}
	chErr := make(chan error, len(cmds))
	for _, cmd := range cmds {
		go func(_cmd *notifiermodel.CreateNotificationArgs) (_err error) {
			defer func() {
				chErr <- _err
			}()
			_, _err = notiStore.CreateNotification(cmd)
			if _err != nil {
				ll.Debug("err", l.Error(_err))
			}
			return
		}(cmd)
	}
	var created, errors int
	for i, n := 0, len(cmds); i < n; i++ {
		err := <-chErr
		if err == nil {
			created++
		} else {
			errors++
		}
	}
	ll.S.Infof("Create notifications: success %v/%v, errors %v/%v", created, len(cmds), errors, len(cmds))
	return nil
}
