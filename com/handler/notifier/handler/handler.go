package handler

import (
	"context"

	handler "o.o/backend/com/handler/etop-handler"
	"o.o/backend/com/handler/etop-handler/pgrid"
	notifiermodel "o.o/backend/com/handler/notifier/model"
	"o.o/backend/com/handler/notifier/sqlstore"
	com "o.o/backend/com/main"
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

func New(dbMain com.MainDB, dbNotifier com.NotifierDB, consumer mq.KafkaConsumer, prefix string) (handlerMain *handler.Handler, handlerNotifier *handler.Handler) {
	x = dbMain
	xNotifier = dbNotifier
	notiStore = sqlstore.NewNotificationStore(dbNotifier)
	deviceStore = sqlstore.NewDeviceStore(dbNotifier)
	historyStore = historysqlstore.NewHistoryStore(dbMain)

	handlersMain := TopicsAndHandlersEtop()
	handlerMain = handler.NewWithHandlers(dbMain, nil, consumer, prefix, handlersMain)

	handlersNotifier := TopicsAndHandlerNotifier()
	handlerNotifier = handler.NewWithHandlers(dbNotifier, nil, consumer, prefix, handlersNotifier)

	return handlerMain, handlerNotifier
}

func TopicsAndHandlersEtop() map[string]pgrid.HandlerFunc {
	return map[string]pgrid.HandlerFunc{
		"fulfillment":                HandleFulfillmentEvent,
		"money_transaction_shipping": HandleMoneyTransactionShippingEvent,
	}
}

func TopicsAndHandlerNotifier() map[string]pgrid.HandlerFunc {
	return map[string]pgrid.HandlerFunc{
		"notification": HandleNotificationEvent,
	}
}

func CreateNotifications(ctx context.Context, cmds []*notifiermodel.CreateNotificationArgs) error {
	if len(cmds) == 0 {
		return nil
	}
	chErr := make(chan error, len(cmds))
	for _, cmd := range cmds {
		go ignoreError(func(_cmd *notifiermodel.CreateNotificationArgs) (_err error) {
			defer func() {
				chErr <- _err
			}()
			_, _err = notiStore.CreateNotification(cmd)
			if _err != nil {
				ll.Debug("err", l.Error(_err))
			}
			return
		}(cmd))
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

func ignoreError(err error) {}
