package handler

import (
	"context"

	handler "etop.vn/backend/com/handler/etop-handler"
	"etop.vn/backend/com/handler/etop-handler/pgrid"
	notifiermodel "etop.vn/backend/com/handler/notifier/model"
	"etop.vn/backend/com/handler/notifier/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/telebot"
	historysqlstore "etop.vn/backend/pkg/etop-history/sqlstore"
	"etop.vn/common/l"
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

func New(dbMain *cmsql.Database, dbNotifier *cmsql.Database, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string) (handlerMain *handler.Handler, handlerNotifier *handler.Handler) {
	x = dbMain
	xNotifier = dbNotifier
	notiStore = sqlstore.NewNotificationStore(dbNotifier)
	deviceStore = sqlstore.NewDeviceStore(dbNotifier)
	historyStore = historysqlstore.NewHistoryStore(dbMain)

	handlersMain := TopicsAndHandlersEtop()
	handlerMain = handler.NewWithHandlers(dbMain, nil, bot, consumer, prefix, handlersMain)

	handlersNotifier := TopicsAndHandlerNotifier()
	handlerNotifier = handler.NewWithHandlers(dbNotifier, nil, bot, consumer, prefix, handlersNotifier)

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
