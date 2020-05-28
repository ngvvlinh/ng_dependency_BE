package vtpay

import (
	"github.com/google/wire"

	"o.o/backend/com/external/payment/vtpay/gateway/aggregate"
	"o.o/backend/com/external/payment/vtpay/gateway/server"
)

var WireSet = wire.NewSet(
	server.New,
	NewAggregate, AggregateMessageBus,
	aggregate.NewAggregate, aggregate.AggregateMessageBus,
)
