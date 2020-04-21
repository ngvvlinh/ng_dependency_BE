package admin

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionbill"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/subscripting/subscriptionproduct"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

func (s *SubscriptionService) Clone() *SubscriptionService {
	res := *s
	return &res
}

func (s *SubscriptionService) CreateSubscriptionProduct(ctx context.Context, r *CreateSubscriptionProductEndpoint) error {
	cmd := &subscriptionproduct.CreateSubrProductCommand{
		Name:        r.Name,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Type:        r.Type,
	}
	if err := subrProductAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbSubrProduct(cmd.Result)
	return nil
}

func (s *SubscriptionService) GetSubscriptionProducts(ctx context.Context, r *GetSubscriptionProductsEndpoint) error {
	query := &subscriptionproduct.ListSubrProductsQuery{}
	if err := subrProductQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	res := convertpb.PbSubrProducts(query.Result)
	r.Result = &types.GetSubrProductsResponse{
		SubscriptionProducts: res,
	}
	return nil
}

func (s *SubscriptionService) DeleteSubscriptionProduct(ctx context.Context, r *DeleteSubscriptionProductEndpoint) error {
	cmd := &subscriptionproduct.DeleteSubrProductCommand{
		ID: r.Id,
	}
	if err := subrProductAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *SubscriptionService) CreateSubscriptionPlan(ctx context.Context, r *CreateSubscriptionPlanEndpoint) error {
	cmd := &subscriptionplan.CreateSubrPlanCommand{
		Name:          r.Name,
		Price:         r.Price,
		Description:   r.Description,
		ProductID:     r.ProductID,
		Interval:      r.Interval,
		IntervalCount: r.IntervalCount,
	}
	if err := subrPlanAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbSubrPlan(cmd.Result)
	return nil
}

func (s *SubscriptionService) UpdateSubscriptionPlan(ctx context.Context, r *UpdateSubscriptionPlanEndpoint) error {
	cmd := &subscriptionplan.UpdateSubrPlanCommand{
		ID:            r.ID,
		Name:          r.Name,
		Price:         r.Price,
		Description:   r.Description,
		ProductID:     r.ProductID,
		Interval:      r.Interval,
		IntervalCount: r.IntervalCount,
	}
	if err := subrPlanAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *SubscriptionService) GetSubscriptionPlans(ctx context.Context, r *GetSubscriptionPlansEndpoint) error {
	query := &subscriptionplan.ListSubrPlansQuery{}
	if err := subrPlanQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetSubrPlansResponse{
		SubscriptionPlans: convertpb.PbSubrPlans(query.Result),
	}
	return nil
}

func (s *SubscriptionService) DeleteSubscriptionPlan(ctx context.Context, r *DeleteSubscriptionPlanEndpoint) error {
	cmd := &subscriptionplan.DeleteSubrPlanCommand{
		ID: r.Id,
	}
	if err := subrPlanAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, r *GetSubscriptionEndpoint) error {
	query := &subscription.GetSubscriptionByIDQuery{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := subscriptionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbSubscription(query.Result)
	return nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, r *GetSubscriptionsEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &subscription.ListSubscriptionsQuery{
		AccountID: r.AccountID,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(r.Filters),
	}
	if err := subscriptionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetSubscriptionsResponse{
		Subscriptions: convertpb.PbSubscriptions(query.Result.Subscriptions),
		Paging:        cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return nil
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, r *CreateSubscriptionEndpoint) error {
	cmd := &subscription.CreateSubscriptionCommand{
		AccountID:            r.AccountID,
		CancelAtPeriodEnd:    r.CancelAtPeriodEnd,
		Lines:                convertpb.Convert_api_SubscriptionLines_To_core_SubscriptionLines(r.Lines),
		BillingCycleAnchorAt: r.BillingCycleAnchorAt.ToTime(),
		Customer:             convertpb.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
	}
	if err := subscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbSubscription(cmd.Result)
	return nil
}

func (s *SubscriptionService) UpdateSubscriptionInfo(ctx context.Context, r *UpdateSubscriptionInfoEndpoint) error {
	cmd := &subscription.UpdateSubscriptionInfoCommand{
		ID:                   r.ID,
		AccountID:            r.AccountID,
		CancelAtPeriodEnd:    r.CancelAtPeriodEnd,
		BillingCycleAnchorAt: r.BillingCycleAnchorAt.ToTime(),
		Customer:             convertpb.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
		Lines:                convertpb.Convert_api_SubscriptionLines_To_core_SubscriptionLines(r.Lines),
	}
	if err := subscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *SubscriptionService) CancelSubscription(ctx context.Context, r *CancelSubscriptionEndpoint) error {
	cmd := &subscription.CancelSubscriptionCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := subscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *SubscriptionService) ActivateSubscription(ctx context.Context, r *ActivateSubscriptionEndpoint) error {
	cmd := &subscription.ActivateSubscriptionCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := subscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, r *DeleteSubscriptionEndpoint) error {
	cmd := &subscription.DeleteSubscriptionCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := subscriptionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}

func (s *SubscriptionService) GetSubscriptionBills(ctx context.Context, r *GetSubscriptionBillsEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &subscriptionbill.ListSubscriptionBillsQuery{
		AccountID: r.AccountID,
		Paging:    *paging,
		Filters:   cmapi.ToFilters(r.Filters),
	}
	if err := subrBillQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	res := convertpb.PbSubrBills(query.Result.SubscriptionBills)
	r.Result = &types.GetSubscriptionBillsResponse{
		SubscriptionBills: res,
		Paging:            cmapi.PbMetaPageInfo(query.Result.Paging),
	}
	return nil
}

func (s *SubscriptionService) CreateSubscriptionBill(ctx context.Context, r *CreateSubscriptionBillEndpoint) error {
	cmd := &subscriptionbill.CreateSubscriptionBillBySubrIDCommand{
		SubscriptionID: r.SubscriptionID,
		AccountID:      r.AccountID,
		TotalAmount:    r.TotalAmount,
		Customer:       convertpb.Convert_api_SubrCustomer_To_core_SubrCustomer(r.Customer),
		Description:    r.Description,
	}
	if err := subrBillAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbSubrBill(cmd.Result)
	return nil
}

func (s *SubscriptionService) ManualPaymentSubscriptionBill(ctx context.Context, r *ManualPaymentSubscriptionBillEndpoint) error {
	cmd := &subscriptionbill.ManualPaymentSubscriptionBillCommand{
		ID:          r.SubscriptionBillID,
		AccountID:   r.AccountID,
		TotalAmount: r.TotalAmount,
	}
	if err := subrBillAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *SubscriptionService) DeleteSubscriptionBill(ctx context.Context, r *DeleteSubscriptionBillEndpoint) error {
	cmd := &subscriptionbill.DeleteSubsciptionBillCommand{
		ID:        r.ID,
		AccountID: r.AccountID,
	}
	if err := subrBillAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}
