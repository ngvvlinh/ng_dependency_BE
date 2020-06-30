package admin

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionbill"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/subscripting/subscriptionproduct"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type SubscriptionService struct {
	session.Session

	SubrProductAggr   subscriptionproduct.CommandBus
	SubrProductQuery  subscriptionproduct.QueryBus
	SubrPlanAggr      subscriptionplan.CommandBus
	SubrPlanQuery     subscriptionplan.QueryBus
	SubscriptionQuery subscription.QueryBus
	SubscriptionAggr  subscription.CommandBus
	SubrBillAggr      subscriptionbill.CommandBus
	SubrBillQuery     subscriptionbill.QueryBus
}

func (s *SubscriptionService) Clone() admin.SubscriptionService {
	res := *s
	return &res
}

func (s *SubscriptionService) CreateSubscriptionProduct(ctx context.Context, r *types.CreateSubrProductRequest) (*types.SubscriptionProduct, error) {
	cmd := &subscriptionproduct.CreateSubrProductCommand{
		Name:        r.Name,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Type:        r.Type,
	}
	if err := s.SubrProductAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbSubrProduct(cmd.Result)
	return result, nil
}

func (s *SubscriptionService) GetSubscriptionProducts(ctx context.Context, r *pbcm.Empty) (*types.GetSubrProductsResponse, error) {
	query := &subscriptionproduct.ListSubrProductsQuery{}
	if err := s.SubrProductQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := convertpb.PbSubrProducts(query.Result)
	result := &types.GetSubrProductsResponse{
		SubscriptionProducts: res,
	}
	return result, nil
}

func (s *SubscriptionService) DeleteSubscriptionProduct(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &subscriptionproduct.DeleteSubrProductCommand{
		ID: r.Id,
	}
	if err := s.SubrProductAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

func (s *SubscriptionService) CreateSubscriptionPlan(ctx context.Context, r *types.CreateSubrPlanRequest) (*types.SubscriptionPlan, error) {
	cmd := &subscriptionplan.CreateSubrPlanCommand{
		Name:          r.Name,
		Price:         r.Price,
		Description:   r.Description,
		ProductID:     r.ProductID,
		Interval:      r.Interval,
		IntervalCount: r.IntervalCount,
	}
	if err := s.SubrPlanAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbSubrPlan(cmd.Result)
	return result, nil
}

func (s *SubscriptionService) UpdateSubscriptionPlan(ctx context.Context, r *types.UpdateSubrPlanRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &subscriptionplan.UpdateSubrPlanCommand{
		ID:            r.ID,
		Name:          r.Name,
		Price:         r.Price,
		Description:   r.Description,
		ProductID:     r.ProductID,
		Interval:      r.Interval,
		IntervalCount: r.IntervalCount,
	}
	if err := s.SubrPlanAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *SubscriptionService) GetSubscriptionPlans(ctx context.Context, r *pbcm.Empty) (*types.GetSubrPlansResponse, error) {
	query := &subscriptionplan.ListSubrPlansQuery{}
	if err := s.SubrPlanQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetSubrPlansResponse{
		SubscriptionPlans: convertpb.PbSubrPlans(query.Result),
	}
	return result, nil
}

func (s *SubscriptionService) DeleteSubscriptionPlan(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &subscriptionplan.DeleteSubrPlanCommand{
		ID: r.Id,
	}
	if err := s.SubrPlanAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
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
	result := convertpb.PbSubscription(query.Result)
	return result, nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, r *types.GetSubscriptionsRequest) (*types.GetSubscriptionsResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &subscription.ListSubscriptionsQuery{
		AccountID: r.AccountID,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(r.Filters),
	}
	if err := s.SubscriptionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetSubscriptionsResponse{
		Subscriptions: convertpb.PbSubscriptions(query.Result.Subscriptions),
		Paging:        cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return result, nil
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, r *types.CreateSubscriptionRequest) (*types.Subscription, error) {
	cmd := &subscription.CreateSubscriptionCommand{
		AccountID:            r.AccountID,
		CancelAtPeriodEnd:    r.CancelAtPeriodEnd,
		Lines:                convertpb.Convert_api_SubscriptionLines_To_core_SubscriptionLines(r.Lines),
		BillingCycleAnchorAt: r.BillingCycleAnchorAt.ToTime(),
		Customer:             convertpb.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
	}
	if err := s.SubscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbSubscription(cmd.Result)
	return result, nil
}

func (s *SubscriptionService) UpdateSubscriptionInfo(ctx context.Context, r *types.UpdateSubscriptionInfoRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &subscription.UpdateSubscriptionInfoCommand{
		ID:                   r.ID,
		AccountID:            r.AccountID,
		CancelAtPeriodEnd:    r.CancelAtPeriodEnd,
		BillingCycleAnchorAt: r.BillingCycleAnchorAt.ToTime(),
		Customer:             convertpb.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
		Lines:                convertpb.Convert_api_SubscriptionLines_To_core_SubscriptionLines(r.Lines),
	}
	if err := s.SubscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *SubscriptionService) CancelSubscription(ctx context.Context, r *types.SubscriptionIDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &subscription.CancelSubscriptionCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.SubscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *SubscriptionService) ActivateSubscription(ctx context.Context, r *types.SubscriptionIDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &subscription.ActivateSubscriptionCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.SubscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, r *types.SubscriptionIDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &subscription.DeleteSubscriptionCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.SubscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}

func (s *SubscriptionService) GetSubscriptionBills(ctx context.Context, r *types.GetSubscriptionBillsRequest) (*types.GetSubscriptionBillsResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &subscriptionbill.ListSubscriptionBillsQuery{
		AccountID: r.AccountID,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(r.Filters),
	}
	if err := s.SubrBillQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := convertpb.PbSubrBills(query.Result.SubscriptionBills)
	result := &types.GetSubscriptionBillsResponse{
		SubscriptionBills: res,
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return result, nil
}

func (s *SubscriptionService) CreateSubscriptionBill(ctx context.Context, r *types.CreateSubscriptionBillRequest) (*types.SubscriptionBill, error) {
	cmd := &subscriptionbill.CreateSubscriptionBillBySubrIDCommand{
		SubscriptionID: r.SubscriptionID,
		AccountID:      r.AccountID,
		TotalAmount:    r.TotalAmount,
		Customer:       convertpb.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
		Description:    r.Description,
	}
	if err := s.SubrBillAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbSubrBill(cmd.Result)
	return result, nil
}

func (s *SubscriptionService) ManualPaymentSubscriptionBill(ctx context.Context, r *types.ManualPaymentSubscriptionBillRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &subscriptionbill.ManualPaymentSubscriptionBillCommand{
		ID:          r.SubscriptionBillID,
		AccountID:   r.AccountID,
		TotalAmount: r.TotalAmount,
	}
	if err := s.SubrBillAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *SubscriptionService) DeleteSubscriptionBill(ctx context.Context, r *types.SubscriptionIDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &subscriptionbill.DeleteSubsciptionBillCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := s.SubrBillAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: 1}
	return result, nil
}
