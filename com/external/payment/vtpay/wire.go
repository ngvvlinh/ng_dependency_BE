// +build wireinject

package vtpay

import (
	"github.com/google/wire"
	"o.o/backend/com/external/payment/vtpay/gateway/aggregate"
)

var WireSet = wire.NewSet(
	NewAggregate, AggregateMessageBus,
	aggregate.NewAggregate, AggregateMessageBus,
)
