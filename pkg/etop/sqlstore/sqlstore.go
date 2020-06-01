package sqlstore

import (
	"time"

	"o.o/api/main/location"
	notisqlstore "o.o/backend/com/handler/notifier/sqlstore"
	com "o.o/backend/com/main"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	servicelocation "o.o/backend/com/main/location"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
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

type filterDeletable interface {
	Prefix() string
	ByDeletedAt(time.Time) *sq.ColumnFilter
}

type includeDeleted bool

func (d includeDeleted) filterDeleted(f filterDeletable) sq.WriterTo {
	if d {
		return nil
	}
	s := "deleted_at IS NULL"
	p := f.Prefix()
	if p != "" {
		s = p + "." + s
	}
	return sq.NewExpr(s)
}

type multiplelity bool

func (m multiplelity) ensureMultiplelity(countable interface{ Count() (int, error) }) error {
	n, err := countable.Count()
	if err != nil {
		return err
	}
	if !m && (n > 1) {
		return cm.Errorf(cm.Internal, nil, "unexpected number of changes")
	}
	return nil
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

func ignoreError(err error) {}
