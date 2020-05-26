package smslog

import (
	"github.com/google/wire"
	"o.o/backend/com/etc/logging/smslog/aggregate"
)

var WireSet = wire.NewSet(
	aggregate.NewSmsLogAggregate, aggregate.SmsLogAggregateMessageBus,
)
