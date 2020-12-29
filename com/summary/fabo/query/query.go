package query

import (
	"context"

	"o.o/api/fabo/summary"
	commonsummary "o.o/api/summary"
	com "o.o/backend/com/main"
	"o.o/backend/com/summary/fabo/sqlstore"
	"o.o/backend/com/summary/util"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
)

const (
	redisCurrentVersion = "v1.1"
	redisTTL            = 1 * 60
)

var _ summary.QueryService = &DashboardQuery{}

func DashboardQueryMessageBus(q *DashboardQuery) summary.QueryBus {
	b := bus.New()
	return summary.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func NewDashboardQuery(
	db com.MainDB, redis redis.Store,
) *DashboardQuery {
	return &DashboardQuery{
		store:      sqlstore.NewSummaryStore(db),
		redisStore: redis,
	}
}

type DashboardQuery struct {
	redisStore redis.Store

	store sqlstore.SummaryStore
}

// result of today wasn't cached
// just cached total and perday
func (d *DashboardQuery) SummaryShop(
	ctx context.Context, req *summary.SummaryShopArgs,
) (*summary.SummaryShopResponse, error) {
	var tablesCached []*commonsummary.SummaryTable

	// always summarize today
	tablesToday, err := d.store.SummarizeShopToday(ctx, req)
	if err != nil {
		return nil, err
	}

	// check dashboard total and perday were cached
	keyRedis := util.BuildKey(req.ShopID, req.DateFrom, req.DateTo, "fabo-fulfillment", redisCurrentVersion)
	isCached, err := d.checkRedis(keyRedis, &tablesCached)
	if err != nil {
		return nil, err
	}
	if !isCached {
		tablesShop, err := d.store.SummarizeShop(ctx, req)
		if err != nil {
			return nil, err
		}

		// cache
		if err = d.redisStore.SetWithTTL(keyRedis, tablesShop, redisTTL); err != nil {
			return nil, err
		}

		result := &summary.SummaryShopResponse{
			ListTable: append(tablesToday, tablesShop...),
		}

		return result, nil
	}

	result := &summary.SummaryShopResponse{
		ListTable: append(tablesToday, tablesCached...),
	}

	return result, nil
}

func (d *DashboardQuery) checkRedis(
	key string, obj interface{},
) (bool, error) {
	err := d.redisStore.Get(key, obj)
	switch err {
	case nil:
		return true, nil
	case redis.ErrNil:
		return false, nil
	default:
		return true, err
	}
}
