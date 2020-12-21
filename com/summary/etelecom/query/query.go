package query

import (
	"context"
	"time"

	summary "o.o/api/etelecom/summary"
	com "o.o/backend/com/main"
	"o.o/backend/com/summary/etelecom/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
)

const (
	currentVersion = "v1.1"
	dateLayout     = "2006-01-02"
)

const redisTTL = 2 * 60

var _ summary.QueryService = &SummaryQuery{}

type SummaryQuery struct {
	store      sqlstore.SummaryStore
	redisStore redis.Store
}

func NewSummaryQuery(db com.EtelecomDB, redis redis.Store) *SummaryQuery {
	return &SummaryQuery{
		store:      sqlstore.NewSummaryStore(db),
		redisStore: redis,
	}
}

func SummaryQueryMessageBus(q *SummaryQuery) summary.QueryBus {
	b := bus.New()
	return summary.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (s *SummaryQuery) Summary(
	ctx context.Context, args *summary.SummaryArgs,
) (*summary.SummaryResponse, error) {
	var summaryTableRedis summary.SummaryResponse

	keyRedis := buildKey(args.ShopID, args.DateFrom, args.DateTo, "etelecom")
	isReturn, err := s.checkRedis(keyRedis, &summaryTableRedis)
	if isReturn {
		return &summaryTableRedis, err
	}

	tables, err := s.store.Summary(ctx, args)
	if err != nil {
		return nil, err
	}

	result := &summary.SummaryResponse{
		ListTable: tables,
	}
	err = s.redisStore.SetWithTTL(keyRedis, result, redisTTL)
	return result, err
}

func buildKey(shopID dot.ID, dateFrom, dateTo time.Time, keycode string) string {
	key := "summary/" + keycode + ":version=" + currentVersion +
		",shop=" + shopID.String() +
		",from=" + dateFrom.Format(dateLayout) +
		",to=" + dateTo.Format(dateLayout)
	return key
}

func (q *SummaryQuery) checkRedis(key string, obj interface{}) (bool, error) {
	err := q.redisStore.Get(key, obj)
	switch err {
	case nil:
		return true, nil
	case redis.ErrNil:
		return false, nil
	default:
		return true, err
	}
}
