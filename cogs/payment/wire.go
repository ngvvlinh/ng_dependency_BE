package _payment

import (
	"github.com/google/wire"

	loggingpayment "o.o/backend/com/etc/logging/payment"
	"o.o/backend/com/external/payment/payment"
)

var WireSet = wire.NewSet(
	payment.WireSet,
	loggingpayment.WireSet,
)
