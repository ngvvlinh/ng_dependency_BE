package aggregate

import (
	"context"

	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/tradering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
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

func (a *TraderAgg) DeleteTrader(ctx context.Context, id int64, shopID int64,
) (deleted int, _ error) {
	_, err := a.store(ctx).ShopID(shopID).ID(id).GetTraderDB()
	if err != nil {
		return 0, err
	}
	deleted, err = a.store(ctx).ID(id).SoftDelete()
	return deleted, err
}
