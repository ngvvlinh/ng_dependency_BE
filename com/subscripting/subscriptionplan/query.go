package subscriptionplan

import (
	"context"

	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/subscripting/subscriptionproduct"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/backend/com/subscripting/subscriptionplan/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscriptionplan.QueryService = &SubrPlanQuery{}

type SubrPlanQuery struct {
	subrPlanStore    sqlstore.SubrPlanStoreFactory
	subrProductQuery subscriptionproduct.QueryBus
}

func NewSubrPlanQuery(db *cmsql.Database, subrProductQuery subscriptionproduct.QueryBus) *SubrPlanQuery {
	return &SubrPlanQuery{
		subrPlanStore:    sqlstore.NewSubrPlanStore(db),
		subrProductQuery: subrProductQuery,
	}
}

func SubrPlanQueryMessageBus(q *SubrPlanQuery) subscriptionplan.QueryBus {
	b := bus.New()
	return subscriptionplan.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *SubrPlanQuery) GetSubrPlanByID(ctx context.Context, id dot.ID) (*subscriptionplan.SubscriptionPlan, error) {
	return q.subrPlanStore(ctx).ID(id).GetSubrPlan()
}

func (q *SubrPlanQuery) ListSubrPlans(ctx context.Context, args *subscriptionplan.ListSubrPlansArgs) ([]*subscriptionplan.SubscriptionPlan, error) {
	query := q.subrPlanStore(ctx)
	if len(args.ProductIDs) > 0 {
		query = query.ProductIDs(args.ProductIDs...)
	}
	return query.ListSubscriptions()
}

func (q *SubrPlanQuery) GetFreeSubrPlanByProductType(ctx context.Context, productType subscription_product_type.ProductSubscriptionType) (*subscriptionplan.SubscriptionPlan, error) {
	queryProduct := &subscriptionproduct.ListSubrProductsQuery{
		Type: productType,
	}
	if err := q.subrProductQuery.Dispatch(ctx, queryProduct); err != nil {
		return nil, err
	}
	productIDs := make([]dot.ID, len(queryProduct.Result))
	for i, p := range queryProduct.Result {
		productIDs[i] = p.ID
	}

	return q.subrPlanStore(ctx).ProductIDs(productIDs...).FreePlan().GetSubrPlan()
}
