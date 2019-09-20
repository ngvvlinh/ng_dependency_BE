package pm

import (
	"context"

	"etop.vn/common/l"

	"etop.vn/api/main/ordering"
	ordertrading "etop.vn/api/main/ordering/trading"
	"etop.vn/api/services/affiliate"
	"etop.vn/common/bus"
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

func (p *ProcessManager) CheckTradingOrderValid(ctx context.Context, event *ordertrading.CheckTradingOrderValidEvent) error {
	ll.V(3).Debug("CheckTradingOrderValid", l.Object("event", event))
	return nil
}

func (p *ProcessManager) TradingOrderCreated(ctx context.Context, event *ordertrading.TradingOrderCreatedEvent) error {
	ll.V(3).Debug("TradingOrderCreated", l.Object("event", event))
	return nil
}
