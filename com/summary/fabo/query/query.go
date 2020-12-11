package query

import (
	"context"

	"o.o/api/fabo/summary"
	com "o.o/backend/com/main"
	"o.o/backend/com/summary/fabo/sqlstore"
	"o.o/backend/com/summary/util"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
)

const (
	redisCurrentVersion = "v1.0"
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

func (d *DashboardQuery) SummaryShop(
	ctx context.Context, req *summary.SummaryShopArgs,
) (*summary.SummaryShopResponse, error) {
	var summaryTableRedis summary.SummaryShopResponse

	// check dashboard was cached
	keyRedis := util.BuildKey(req.ShopID, req.DateFrom, req.DateTo, "fabo-fulfillment", redisCurrentVersion)
	isReturn, err := d.checkRedis(keyRedis, &summaryTableRedis)
	if err != nil {
		return nil, err
	}
	if isReturn {
		return &summaryTableRedis, nil
	}

	tables, err := d.store.SummarizeShop(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &summary.SummaryShopResponse{
		ListTable: tables,
	}
	err = d.redisStore.SetWithTTL(keyRedis, result, redisTTL)
	return result, err
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
