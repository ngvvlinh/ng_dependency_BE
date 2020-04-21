package pm

import (
	"context"
	"time"

	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionbill"
	"o.o/api/top/types/etc/status4"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	subrBillQuery subscriptionbill.QueryBus
	subrQuery     subscription.QueryBus
	subrAggr      subscription.CommandBus
}

func New(subrBillQuery subscriptionbill.QueryBus, subrQuery subscription.QueryBus, subrAggr subscription.CommandBus) *ProcessManager {
	return &ProcessManager{
		subrBillQuery: subrBillQuery,
		subrQuery:     subrQuery,
		subrAggr:      subrAggr,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.SubscriptionBillPaid)
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
