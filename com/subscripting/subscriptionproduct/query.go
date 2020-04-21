package subscriptionproduct

import (
	"context"

	"o.o/api/subscripting/subscriptionproduct"
	"o.o/api/top/types/common"
	"o.o/backend/com/subscripting/subscriptionproduct/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscriptionproduct.QueryService = &SubrProductQuery{}

type SubrProductQuery struct {
	subrProductStore sqlstore.SubrProductStoreFactory
}

func NewSubrProductQuery(db *cmsql.Database) *SubrProductQuery {
	return &SubrProductQuery{
		subrProductStore: sqlstore.NewSubscriptionProductStore(db),
	}
}

func (q *SubrProductQuery) MessageBus() subscriptionproduct.QueryBus {
	b := bus.New()
	return subscriptionproduct.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *SubrProductQuery) GetSubrProductByID(ctx context.Context, id dot.ID) (*subscriptionproduct.SubscriptionProduct, error) {
	return q.subrProductStore(ctx).ID(id).GetSubrProduct()
}

func (q *SubrProductQuery) ListSubrProducts(ctx context.Context, _ *common.Empty) ([]*subscriptionproduct.SubscriptionProduct, error) {
	return q.subrProductStore(ctx).ListSubscriptions()
}
