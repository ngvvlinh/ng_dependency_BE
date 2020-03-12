package sqlstore

import (
	"time"

	notisqlstore "etop.vn/backend/com/handler/notifier/sqlstore"
	catalogsqlstore "etop.vn/backend/com/main/catalog/sqlstore"
	servicelocation "etop.vn/backend/com/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi"
	"etop.vn/common/l"
)

var (
	x                 *cmsql.Database
	xNotifier         *cmsql.Database
	ll                = l.New()
	deviceStore       *notisqlstore.DeviceStore
	notificationStore *notisqlstore.NotificationStore
	eventBus          capi.EventBus
	locationBus       = servicelocation.New(nil).MessageBus()
)

type (
	M  map[string]interface{}
	Ms map[string]string

	Query = cmsql.Query
	Qx    = cmsql.QueryInterface
)

func Init(db *cmsql.Database) {
	if x != nil {
		if (*x).DB() != nil {
			ll.Panic("Already initialized")
		}
	}
	x = db
	shopProductStore = catalogsqlstore.NewShopProductStore(db)
}

func AddEventBus(_eventBus capi.EventBus) {
	eventBus = _eventBus
}

func InitDBNotifier(db *cmsql.Database) {
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
