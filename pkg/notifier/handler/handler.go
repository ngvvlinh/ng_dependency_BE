package handler

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/etop-handler/pgrid"
	handlerwebhook "etop.vn/backend/pkg/etop-handler/webhook"
	"etop.vn/backend/pkg/notifier/model"
	"etop.vn/backend/pkg/notifier/sqlstore"
	"etop.vn/common/l"
)

var (
	x           cmsql.Database
	xNotifier   cmsql.Database
	ll          = l.New()
	notiStore   *sqlstore.NotificationStore
	deviceStore *sqlstore.DeviceStore
)

const ConsumerGroup = "handler/notifier"

func New(dbMain cmsql.Database, dbNotifier cmsql.Database, bot *telebot.Channel, consumer mq.KafkaConsumer, prefix string) (handlerMain *handlerwebhook.Handler, handlerNotifier *handlerwebhook.Handler) {
	x = dbMain
	xNotifier = dbNotifier
	notiStore = sqlstore.NewNotificationStore(dbNotifier)
	deviceStore = sqlstore.NewDeviceStore(dbNotifier)

	handlersMain := TopicsAndHandlersEtop()
	handlerMain = handlerwebhook.NewWithHandlers(dbMain, nil, bot, consumer, prefix, handlersMain)

	handlersNotifier := TopicsAndHandlerNotifier()
	handlerNotifier = handlerwebhook.NewWithHandlers(dbNotifier, nil, bot, consumer, prefix, handlersNotifier)

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

func CreateNotifications(ctx context.Context, cmds []*model.CreateNotificationArgs) error {
	if len(cmds) == 0 {
		return nil
	}
	chErr := make(chan error, len(cmds))
	for _, cmd := range cmds {
		go func(_cmd *model.CreateNotificationArgs) (_err error) {
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
	for i, l := 0, len(cmds); i < l; i++ {
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
