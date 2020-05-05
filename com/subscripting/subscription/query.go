package subscription

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/subscripting/subscriptionproduct"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/backend/com/subscripting/subscription/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscription.QueryService = &SubscriptionQuery{}

type SubscriptionQuery struct {
	subscriptionStore sqlstore.SubscriptionStoreFactory
	subrPlanQuery     subscriptionplan.QueryBus
	subrProductQuery  subscriptionproduct.QueryBus
}

func NewSubscriptionQuery(db *cmsql.Database, subrPlanQuery subscriptionplan.QueryBus, subrProductQuery subscriptionproduct.QueryBus) *SubscriptionQuery {
	return &SubscriptionQuery{
		subscriptionStore: sqlstore.NewSubscriptionStore(db),
		subrPlanQuery:     subrPlanQuery,
		subrProductQuery:  subrProductQuery,
	}
}

func (q *SubscriptionQuery) MessageBus() subscription.QueryBus {
	b := bus.New()
	return subscription.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *SubscriptionQuery) GetSubscriptionByID(ctx context.Context, id dot.ID, accountID dot.ID) (*subscription.SubscriptionFtLine, error) {
	return q.subscriptionStore(ctx).ID(id).OptionalAccountID(accountID).GetSubscriptionFtLine()
}

func (q *SubscriptionQuery) ListSubscriptions(ctx context.Context, args *subscription.ListSubscriptionsArgs) (*subscription.ListSubscriptionsResponse, error) {
	query := q.subscriptionStore(ctx).OptionalAccountID(args.AccountID).WithPaging(args.Paging).Filters(args.Filters)
	res, err := query.ListSubscriptionFtLines()
	if err != nil {
		return nil, err
	}
	return &subscription.ListSubscriptionsResponse{
		Subscriptions: res,
		Paging:        query.GetPaging(),
	}, nil
}

func (q *SubscriptionQuery) GetLastestSubscriptionByProductType(ctx context.Context, accountID dot.ID, productType subscription_product_type.ProductSubscriptionType) (*subscription.SubscriptionFtLine, error) {
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
	if len(productIDs) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "")
	}

	queryPlan := &subscriptionplan.ListSubrPlansQuery{
		ProductIDs: productIDs,
	}
	if err := q.subrPlanQuery.Dispatch(ctx, queryPlan); err != nil {
		return nil, err
	}
	plans := queryPlan.Result
	if len(plans) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "")
	}
	planIDs := make([]dot.ID, len(plans))
	for i, p := range plans {
		planIDs[i] = p.ID
	}

	return q.subscriptionStore(ctx).AccountID(accountID).PlanIDs(planIDs...).Status(status3.P).GetSubscriptionFtLine()
}
