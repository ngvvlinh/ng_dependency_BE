package pm

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionbill"
	"o.o/api/subscripting/subscriptionplan"
	subscriptingtypes "o.o/api/subscripting/types"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/api/webserver"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	subrBillQuery subscriptionbill.QueryBus
	subrBillAggr  subscriptionbill.CommandBus
	subrQuery     subscription.QueryBus
	subrAggr      subscription.CommandBus
	subrPlanQuery subscriptionplan.QueryBus
	identityQuery identity.QueryBus
}

func New(
	subrBillQuery subscriptionbill.QueryBus,
	subrBillAggr subscriptionbill.CommandBus,
	subrQuery subscription.QueryBus,
	subrAggr subscription.CommandBus,
	subrPlanQuery subscriptionplan.QueryBus,
	identityQuery identity.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		subrBillQuery: subrBillQuery,
		subrBillAggr:  subrBillAggr,
		subrQuery:     subrQuery,
		subrAggr:      subrAggr,
		subrPlanQuery: subrPlanQuery,
		identityQuery: identityQuery,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.SubscriptionBillPaid)
	eventBus.AddEventListener(m.WsWebsiteCreated)
}

func (m *ProcessManager) SubscriptionBillPaid(ctx context.Context, event *subscriptionbill.SubscriptionBillPaidEvent) error {
	queryBill := &subscriptionbill.GetSubscriptionBillByIDQuery{
		ID:        event.ID,
		AccountID: event.AccountID,
	}
	if err := m.subrBillQuery.Dispatch(ctx, queryBill); err != nil {
		return err
	}
	subrBill := queryBill.Result
	if subrBill.Status != status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Subscription bill was not paid").WithMetap("event", "SubscriptionBillPaidEvent")
	}

	querySubr := &subscription.GetSubscriptionByIDQuery{
		ID:        subrBill.SubscriptionID,
		AccountID: subrBill.AccountID,
	}
	if err := m.subrQuery.Dispatch(ctx, querySubr); err != nil {
		return err
	}
	subr := querySubr.Result

	// assume that only has 1 line in subscription_line
	var periodStartAt, periodEndAt time.Time
	for _, line := range subrBill.Lines {
		if !periodStartAt.IsZero() && periodStartAt != line.PeriodStartAt {
			return cm.Errorf(cm.InvalidArgument, nil, "Thời gian bắt đầu trong subscription bill line không được khác nhau")
		}
		if !periodEndAt.IsZero() && periodEndAt != line.PeriodStartAt {
			return cm.Errorf(cm.InvalidArgument, nil, "Thời gian kết thúc trong subscription bill line không được khác nhau")
		}
		periodStartAt = line.PeriodStartAt
		periodEndAt = line.PeriodEndAt
	}
	update := &subscription.UpdateSubscriptionPeriodCommand{
		ID:                   subr.ID,
		AccountID:            subr.AccountID,
		CurrentPeriodStartAt: periodStartAt,
		CurrentPeriodEndAt:   periodEndAt,
	}
	if subr.CurrentPeriodStartAt.IsZero() {
		update.StartedAt = periodStartAt
	}
	return m.subrAggr.Dispatch(ctx, update)
}

/*
	WsWebsiteCreatedEvent

	Create subscription with 1 month free plan
*/

func (m *ProcessManager) WsWebsiteCreated(ctx context.Context, event *webserver.WsWebsiteCreatedEvent) error {
	query := &identity.GetShopByIDQuery{
		ID: event.ShopID,
	}
	if err := m.identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	shop := query.Result

	userQuery := &identity.GetUserByIDQuery{
		UserID: shop.OwnerID,
	}
	if err := m.identityQuery.Dispatch(ctx, userQuery); err != nil {
		return err
	}
	user := userQuery.Result

	planQuery := &subscriptionplan.GetFreeSubrPlanByProductTypeQuery{
		ProductType: subscription_product_type.Ecomify,
	}
	if err := m.subrPlanQuery.Dispatch(ctx, planQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return cm.Errorf(cm.NotFound, nil, "Vui lòng tạo gói dùng thử cho ecomify")
		}
		return err
	}
	plan := planQuery.Result

	cmd := &subscription.CreateSubscriptionCommand{
		AccountID: event.ShopID,
		Lines: []*subscription.SubscriptionLine{
			{
				PlanID:   plan.ID,
				Quantity: 1,
			},
		},
		Customer: &subscriptingtypes.CustomerInfo{
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}
	if err := m.subrAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	subr := cmd.Result

	// create subr bill & confirm to active trial product
	billCmd := &subscriptionbill.CreateSubscriptionBillBySubrIDCommand{
		SubscriptionID: subr.ID,
		AccountID:      shop.ID,
		TotalAmount:    0,
		Description:    "Trial",
	}
	if err := m.subrBillAggr.Dispatch(ctx, billCmd); err != nil {
		return err
	}

	billConfirmCmd := &subscriptionbill.ManualPaymentSubscriptionBillCommand{
		ID:          billCmd.Result.ID,
		AccountID:   shop.ID,
		TotalAmount: 0,
	}
	if err := m.subrBillAggr.Dispatch(ctx, billConfirmCmd); err != nil {
		return err
	}

	return nil
}
