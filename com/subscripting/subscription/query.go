package subscription

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/backend/com/subscripting/subscription/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ subscription.QueryService = &SubscriptionQuery{}

type SubscriptionQuery struct {
	subscriptionStore sqlstore.SubscriptionStoreFactory
}

func NewSubscriptionQuery(db *cmsql.Database) *SubscriptionQuery {
	return &SubscriptionQuery{
		subscriptionStore: sqlstore.NewSubscriptionStore(db),
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
