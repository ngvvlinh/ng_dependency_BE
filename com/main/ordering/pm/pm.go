package pm

import (
	"context"

	"etop.vn/common/l"

	"etop.vn/api/main/ordering"
	ordertrading "etop.vn/api/main/ordering/trading"
	"etop.vn/api/services/affiliate"
	"etop.vn/backend/pkg/common/bus"
)

type ProcessManager struct {
	order     ordering.CommandBus
	orderQS   ordering.QueryBus
	affiliate affiliate.CommandBus
}

var (
	ll = l.New()
)

func New(
	orderAggr ordering.CommandBus,
	orderQuery ordering.QueryBus,
	affiliateAggr affiliate.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		order:     orderAggr,
		orderQS:   orderQuery,
		affiliate: affiliateAggr,
	}
}

func (p *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(p.CheckTradingOrderValid)
	eventBus.AddEventListener(p.TradingOrderCreated)
}

func (p *ProcessManager) CheckTradingOrderValid(ctx context.Context, event *ordertrading.TradingOrderCreatingEvent) error {
	checkCmd := &affiliate.TradingOrderCreatingCommand{
		ReferralCode: event.ReferralCode,
		UserID:       event.UserID,
	}
	if err := p.affiliate.Dispatch(ctx, checkCmd); err != nil {
		return err
	}
	return nil
}

func (p *ProcessManager) TradingOrderCreated(ctx context.Context, event *ordertrading.TradingOrderCreatedEvent) error {
	orderCreatedNotifyCmd := &affiliate.OnTradingOrderCreatedCommand{
		OrderID:      event.OrderID,
		ReferralCode: event.ReferralCode,
	}
	if err := p.affiliate.Dispatch(ctx, orderCreatedNotifyCmd); err != nil {
		return err
	}
	return nil
}
