package payment

import (
	"github.com/google/wire"

	"o.o/backend/com/etc/logging/payment/aggregate"
)

var WireSet = wire.NewSet(
	aggregate.New,
)
