package subscription

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/subscripting/subscriptionproduct"
	api "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type SubscriptionService struct {
	session.Session

	SubscriptionQuery subscription.QueryBus
	SubrProductQuery  subscriptionproduct.QueryBus
	SubrPlanQuery     subscriptionplan.QueryBus
}

func (s *SubscriptionService) Clone() api.SubscriptionService {
	res := *s
	return &res
}

func (s *SubscriptionService) GetSubscriptionProducts(ctx context.Context, r *types.GetSubrProductsRequest) (*types.GetSubrProductsResponse, error) {
	query := &subscriptionproduct.ListSubrProductsQuery{
		Type: r.Type,
	}
	if err := s.SubrProductQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetSubrProductsResponse{
		SubscriptionProducts: convertpball.PbSubrProducts(query.Result),
	}
	return result, nil
}

func (s *SubscriptionService) GetSubscriptionPlans(ctx context.Context, r *types.GetSubrPlansRequest) (*types.GetSubrPlansResponse, error) {
	query := &subscriptionplan.ListSubrPlansQuery{}
	if r.ProductID != 0 {
		query.ProductIDs = []dot.ID{r.ProductID}
	}
	if err := s.SubrPlanQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetSubrPlansResponse{
		SubscriptionPlans: convertpball.PbSubrPlans(query.Result),
	}
	return result, nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, r *types.SubscriptionIDRequest) (*types.Subscription, error) {
	query := &subscription.GetSubscriptionByIDQuery{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.SubscriptionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpball.PbSubscription(query.Result)
	return result, nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, r *types.GetSubscriptionsRequest) (*types.GetSubscriptionsResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &subscription.ListSubscriptionsQuery{
		AccountID: s.SS.Shop().ID,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(r.Filters),
	}
	if err := s.SubscriptionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetSubscriptionsResponse{
		Subscriptions: convertpball.PbSubscriptions(query.Result.Subscriptions),
		Paging:        cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return result, nil
}
