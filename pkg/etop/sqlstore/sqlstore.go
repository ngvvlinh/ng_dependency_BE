package sqlstore

import (
	"o.o/api/main/location"
	com "o.o/backend/com/main"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	servicelocation "o.o/backend/com/main/location"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/common/l"
)

var (
	ll          = l.New()
	eventBus    capi.EventBus
	locationBus = servicelocation.QueryMessageBus(servicelocation.New(nil))
)

type (
	M  map[string]interface{}
	Ms map[string]string

	Query = cmsql.Query
	Qx    = cmsql.QueryInterface
)

type Store struct{}

func New(db com.MainDB, _locationBus location.QueryBus, _eventBus capi.EventBus) *Store {
	shopProductStore = catalogsqlstore.NewShopProductStore(db)
	locationBus = _locationBus
	eventBus = _eventBus
	return nil
}

func inTransaction(db *cmsql.Database, callback func(cmsql.QueryInterface) error) (err error) {
	return db.InTransaction(bus.Ctx(), callback)
}

func IDs(items []int64) []interface{} {
	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item
	}
	return res
}
