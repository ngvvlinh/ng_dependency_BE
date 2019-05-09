package cache

import (
	"context"
	"strconv"
	"sync"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/model"

	"github.com/lib/pq"
)

func init() {
	bus.AddHandler("cache", GetSuppliersRules)
}

var (
	redisStore redis.Store
	waiting    bool

	x  cmsql.Database
	ll = l.New()
	wg sync.WaitGroup
)

const (
	DefaultRulesTTL = 60
	PrefixRules     = "srules:"

	EventSupplierRulesUpdate = "supplier_rules_update"
)

type SupplierRulesCache struct {
}

func NewRedisCache(r redis.Store, connStr string) error {
	_, err := newRedisCache(r, connStr, false)
	return err
}

// NewRedisCacheAndWait should be used only in tests
func NewRedisCacheAndWait(r redis.Store, connStr string) (<-chan string, error) {
	return newRedisCache(r, connStr, true)
}

func newRedisCache(r redis.Store, connStr string, notify bool) (<-chan string, error) {
	ch := make(chan string)
	redisStore = r
	return ch, nil
}

func newRedisCache1(r redis.Store, cfg cc.Postgres, notify bool) (<-chan string, error) {
	if redisStore != nil {
		ll.Panic("Already initialized")
	}
	_, connStr := cfg.ConnectionString()

	// TODO: Migrate to cloudsql
	engine, err := cmsql.Connect(cmsql.ConfigPostgres(cfg))
	if err != nil {
		return nil, err
	}

	redisStore = r
	x = engine

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			ll.Error("Postgres event", l.Error(err))
		}
	}

	listener := pq.NewListener(connStr, 500*time.Millisecond, 5*time.Second, reportProblem)
	err = listener.Listen(EventSupplierRulesUpdate)
	if err != nil {
		return nil, err
	}

	var ch chan string
	if !waiting {
		waiting = true

		if notify {
			ch = make(chan string, 256)
		}
		go waitForNotifications(listener, ch)
	}
	return ch, nil
}

type SupplierPriceRulesWithID struct {
	ID    int64
	Rules *model.SupplierPriceRules
}

type GetSuppliersRulesQuery struct {
	SupplierIDs []int64

	// Fallback must return coressponding ids
	Fallback func(ids []int64) ([]*SupplierPriceRulesWithID, error)

	Result struct {
		SupplierRules map[int64]*model.SupplierPriceRules
	}
}

type SetSuppliersRulesCommand struct {
	SupplierIDs   []int64
	SupplierRules []*model.SupplierPriceRules
}

func GetSuppliersRules(ctx context.Context, query *GetSuppliersRulesQuery) error {
	ids := make([]int64, 0, len(query.SupplierIDs))
	res := make(map[int64]*model.SupplierPriceRules)
	query.Result.SupplierRules = res
	for _, id := range query.SupplierIDs {
		var rules model.SupplierPriceRules
		key := GetSupplierRulesKey(id)
		err := redisStore.Get(key, &rules)
		if err != nil {
			ids = append(ids, id)
			continue
		}
		res[id] = &rules
	}

	if len(ids) > 0 && query.Fallback != nil {
		supplierRules, err := query.Fallback(ids)
		if err != nil {
			return err
		}

		m := make(map[int64]*model.SupplierPriceRules)
		for _, rules := range supplierRules {
			if rules == nil {
				continue
			}
			key := GetSupplierRulesKey(rules.ID)
			redisStore.SetWithTTL(key, rules.Rules, DefaultRulesTTL)
			m[rules.ID] = rules.Rules
		}
		for _, id := range query.SupplierIDs {
			if res[id] == nil {
				res[id] = m[id]
			}
		}
	}
	return nil
}

func SetSuppliersRules(ctx context.Context, cmd *SetSuppliersRulesCommand) error {
	if len(cmd.SupplierIDs) != len(cmd.SupplierRules) {
		return cm.Error(cm.InvalidArgument, "Length not match", nil)
	}
	for i, id := range cmd.SupplierIDs {
		rules := cmd.SupplierRules[i]
		if rules == nil {
			continue
		}
		key := GetSupplierRulesKey(id)
		redisStore.SetWithTTL(key, rules, DefaultRulesTTL)
	}
	return nil
}

func GetSupplierRulesKey(id int64) string {
	return PrefixRules + strconv.Itoa(int(id))
}

func waitForNotifications(listener *pq.Listener, notify chan<- string) {
	ll.Info("Start listening to postgres event")
	for {
		select {
		case event := <-listener.Notify:
			id, err := strconv.Atoi(event.Extra)
			if err != nil {
				ll.Error("Error while decoding id from postgres events", l.Error(err))
				continue
			}

			key := GetSupplierRulesKey(int64(id))
			redisStore.Del(key)

			if notify != nil {
				select {
				case notify <- key:
				default:
				}
			}
		}
	}
}
