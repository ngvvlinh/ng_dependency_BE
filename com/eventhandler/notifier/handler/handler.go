package handler

import (
	"o.o/api/main/authorization"
	handler "o.o/backend/com/eventhandler/handler"
	"o.o/backend/com/eventhandler/notifier/sqlstore"
	"o.o/backend/com/eventhandler/pgevent"
	fabosqlstore "o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	com "o.o/backend/com/main"
	connectioning "o.o/backend/com/main/connectioning/sqlstore"
	identitystore "o.o/backend/com/main/identity/sqlstore"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/common/l"
)

var (
	x                         *cmsql.Database
	xNotifier                 *cmsql.Database
	ll                        = l.New()
	notifyStore               *sqlstore.NotificationStore
	deviceStore               *sqlstore.DeviceStore
	historyStore              historysqlstore.HistoryStoreFactory
	userNotifySettingStore    sqlstore.UserNotiSettingStoreFactory
	accountUserStore          identitystore.AccountUserStoreFactory
	customerConversationStore fabosqlstore.FbCustomerConversationStoreFactory
	connectionStore           connectioning.ConnectionStoreFactory
)

const (
	TopicOrder                    = "order"
	TopicSystem                   = "system"
	TopicFulfillment              = "fulfilment"
	TopicFBComment                = "fb_external_comment"
	TopicFBMessage                = "fb_external_message"
	TopicMoneyTransactionShipping = "money_transaction_shipping"
)

var notifyTopicRolesMap = map[string][]authorization.Role{
	TopicOrder:                    {authorization.RoleShopOwner, authorization.RoleSalesMan},
	TopicSystem:                   {authorization.RoleShopOwner, authorization.RoleSalesMan, authorization.RoleStaffManagement},
	TopicFulfillment:              {authorization.RoleShopOwner, authorization.RoleSalesMan},
	TopicFBComment:                {authorization.RoleShopOwner, authorization.RoleSalesMan},
	TopicFBMessage:                {authorization.RoleShopOwner, authorization.RoleSalesMan},
	TopicMoneyTransactionShipping: {authorization.RoleShopOwner},
}

const ConsumerGroup = "handler/notifier"

func New(dbMain com.MainDB, dbNotifier com.NotifierDB, consumer mq.KafkaConsumer, cfg cc.Kafka) (handlerMain *handler.Handler, handlerNotifier *handler.Handler) {
	x = dbMain
	xNotifier = dbNotifier
	notifyStore = sqlstore.NewNotificationStore(dbNotifier)
	deviceStore = sqlstore.NewDeviceStore(dbNotifier)
	historyStore = historysqlstore.NewHistoryStore(dbMain)
	userNotifySettingStore = sqlstore.NewUserNotiSettingStore(dbMain)
	accountUserStore = identitystore.NewAccountUserStore(dbMain)
	customerConversationStore = fabosqlstore.NewFbCustomerConversationStore(dbMain)
	connectionStore = connectioning.NewConnectionStore(dbMain)

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

func TopicsAndHandlersFabo() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		"fulfillment":         HandleFulfillmentEvent,
		"fb_external_comment": HandleCommentEvent,
		"fb_external_message": HandleMessageEvent,
	})
}

func TopicsAndHandlerNotifier() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		"notification": HandleNotificationEvent,
	})
}
