package aggregate

import (
	"context"

	"o.o/api/shopping/tradering"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/tradering/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ tradering.Aggregate = &TraderAgg{}

type TraderAgg struct {
	store  sqlstore.TraderStoreFactory
	trader tradering.QueryBus
}

func NewTraderAgg(
	db com.MainDB,
) *TraderAgg {
	return &TraderAgg{
		store: sqlstore.NewTraderStore(db),
	}
}

func TraderAggMessageBus(q *TraderAgg) tradering.CommandBus {
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
