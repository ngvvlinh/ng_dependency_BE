package sqlstore

import (
	"o.o/api/main/location"
	notisqlstore "o.o/backend/com/eventhandler/notifier/sqlstore"
	com "o.o/backend/com/main"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	servicelocation "o.o/backend/com/main/location"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/common/l"
)

var (
	x                 *cmsql.Database
	xNotifier         *cmsql.Database
	ll                = l.New()
	deviceStore       *notisqlstore.DeviceStore
	notificationStore *notisqlstore.NotificationStore
	eventBus          capi.EventBus
	locationBus       = servicelocation.QueryMessageBus(servicelocation.New(nil))
)

type (
	M  map[string]interface{}
	Ms map[string]string

	Query = cmsql.Query
	Qx    = cmsql.QueryInterface
)

type Store struct{}

func New(db com.MainDB, notiDB com.NotifierDB, _locationBus location.QueryBus, _eventBus capi.EventBus) *Store {
	if x != nil {
		if (*x).DB() != nil {
			ll.Panic("Already initialized")
		}
	}
	x = db
	shopProductStore = catalogsqlstore.NewShopProductStore(db)
	locationBus = _locationBus
	eventBus = _eventBus
	if notiDB != nil {
		initDBNotifier(notiDB) // TODO(qv): remove this
	}
	return nil
}

func initDBNotifier(db *cmsql.Database) {
	if xNotifier != nil && (*xNotifier).DB() != nil {
		ll.Panic("Database Notifier already initialized")
	}
	xNotifier = db
	deviceStore = notisqlstore.NewDeviceStore(xNotifier)
	notificationStore = notisqlstore.NewNotificationStore(xNotifier)
}

func inTransaction(callback func(cmsql.QueryInterface) error) (err error) {
	return x.InTransaction(bus.Ctx(), callback)
}

func IDs(items []int64) []interface{} {
	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item
	}
	return res
}
