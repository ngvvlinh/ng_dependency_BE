package sqlstore

import (
	"strconv"
	"strings"
	"time"

	notisqlstore "etop.vn/backend/com/handler/notifier/sqlstore"
	catalogsqlstore "etop.vn/backend/com/main/catalog/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/backend/pkg/etop/model"
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

func LimitSort(s Query, p *cm.Paging, sortWhitelist map[string]string) (cmsql.Query, error) {
	query, err := sqlstore.LimitSort(s, p, sortWhitelist)
	if err != nil {
		return cmsql.Query{}, err
	}
	return query, nil
}

func Filters(s cmsql.Query, filters []cm.Filter, whitelist FilterWhitelist) (cmsql.Query, bool, error) {
	query, ok, err := sqlstore.Filters(s, filters, whitelist)
	if err != nil {
		return cmsql.Query{}, ok, err
	}
	return query, ok, nil
}

func IDs(items []int64) []interface{} {
	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item
	}
	return res
}

func FilterStatus(s cmsql.Query, prefix string, query model.StatusQuery) cmsql.Query {
	if query.Status != nil {
		s = s.Where(prefix+"status = ?", query.Status)
	}
	return s
}

func Sort(s Query, sorts []string, whitelist map[string]string) (Query, error) {
	for _, sort := range sorts {
		sort = strings.TrimSpace(sort)
		if sort == "" {
			continue
		}

		field := sort
		desc := ""
		if sort[0] == '-' {
			field = sort[1:]
			desc = " DESC"
		}

		if sortField, ok := whitelist[field]; ok {
			if sortField == "" {
				sortField = field
			}
			s = s.OrderBy(sortField + desc)
		} else {
			return s, cm.Errorf(cm.InvalidArgument, nil, "Sort by %v is not allowed", field)
		}
	}
	return s, nil
}

type FilterWhitelist = sqlstore.FilterWhitelist

func countBool(A ...bool) int {
	c := 0
	for _, a := range A {
		if a {
			c++
		}
	}
	return c
}

type valueConfig struct {
	isNumber   bool
	isDate     bool
	isStatus   bool
	isBool     bool
	isNullable bool
}

func parseValue(v string, cfg valueConfig) (interface{}, error) {
	if cfg.isNullable {
		n, err := strconv.ParseBool(v)
		if err != nil {
			return 0, cm.Error(cm.InvalidArgument, "Invalid bool: "+v, nil)
		}
		// nullable will be handled specially at caller
		return n, nil
	}
	if cfg.isBool {
		n, err := strconv.ParseBool(v)
		if err != nil {
			return 0, cm.Error(cm.InvalidArgument, "Invalid bool: "+v, nil)
		}
		return n, nil
	}
	if cfg.isNumber {
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, cm.Error(cm.InvalidArgument, "Invalid number: "+v, nil)
		}
		return n, nil
	}
	if cfg.isDate {
		t, ok := httpreq.ParseAsISO8601([]byte(v))
		if !ok {
			return 0, cm.Error(cm.InvalidArgument, "Invalid date: "+v, nil)
		}
		return t, nil
	}
	if cfg.isStatus {
		switch v {
		case "P", "1":
			return 1, nil
		case "Z", "0":
			return 0, nil
		case "N", "-1":
			return -1, nil
		case "S", "2":
			return 2, nil
		case "NS", "-2":
			return -2, nil
		}
	}
	return v, nil
}

func ignoreError(err error) {}
