package subscriptionplan

import (
	"context"

	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/top/types/common"
	"o.o/backend/com/subscripting/subscriptionplan/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscriptionplan.QueryService = &SubrPlanQuery{}

type SubrPlanQuery struct {
	subrPlanStore sqlstore.SubrPlanStoreFactory
}

func NewSubrPlanQuery(db *cmsql.Database) *SubrPlanQuery {
	return &SubrPlanQuery{
		subrPlanStore: sqlstore.NewSubrPlanStore(db),
	}
}

func (q *SubrPlanQuery) MessageBus() subscriptionplan.QueryBus {
	b := bus.New()
	return subscriptionplan.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *SubrPlanQuery) GetSubrPlanByID(ctx context.Context, id dot.ID) (*subscriptionplan.SubscriptionPlan, error) {
	return q.subrPlanStore(ctx).ID(id).GetSubrPlan()
}

func (q *SubrPlanQuery) ListSubrPlans(ctx context.Context, _ *common.Empty) ([]*subscriptionplan.SubscriptionPlan, error) {
	return q.subrPlanStore(ctx).ListSubscriptions()
}
