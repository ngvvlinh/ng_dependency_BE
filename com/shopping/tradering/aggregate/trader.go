package aggregate

import (
	"context"

	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/tradering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ tradering.Aggregate = &TraderAgg{}

type TraderAgg struct {
	store  sqlstore.TraderStoreFactory
	trader tradering.QueryBus
}

func NewTraderAgg(
	db *cmsql.Database,
) *TraderAgg {
	return &TraderAgg{
		store: sqlstore.NewTraderStore(db),
	}
}

func (q *TraderAgg) MessageBus() tradering.CommandBus {
	b := bus.New()
	return tradering.NewAggregateHandler(q).RegisterHandlers(b)
}

func (q *TraderAgg) DeleteTrader(ctx context.Context, id dot.ID, shopID dot.ID,
) (deleted int, _ error) {
	_, err := q.store(ctx).ShopID(shopID).ID(id).GetTraderDB()
	if err != nil {
		return 0, err
	}
	deleted, err = q.store(ctx).ID(id).SoftDelete()
	return deleted, err
}
