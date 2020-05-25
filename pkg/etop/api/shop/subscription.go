package shop

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

type SubscriptionService struct {
	SubscriptionQuery subscription.QueryBus
}

func (s *SubscriptionService) Clone() *SubscriptionService {
	res := *s
	return &res
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, r *GetSubscriptionEndpoint) error {
	query := &subscription.GetSubscriptionByIDQuery{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.SubscriptionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbSubscription(query.Result)
	return nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, r *GetSubscriptionsEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &subscription.ListSubscriptionsQuery{
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.SubscriptionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetSubscriptionsResponse{
		Subscriptions: convertpb.PbSubscriptions(query.Result.Subscriptions),
		Paging:        cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return nil
}
