package _vtpay

import (
	"github.com/google/wire"

	"o.o/backend/com/external/payment/vtpay"
	paymentvtpay "o.o/backend/pkg/integration/payment/vtpay"
)

var WireSet = wire.NewSet(
	vtpay.WireSet,
	paymentvtpay.WireSet,
)
