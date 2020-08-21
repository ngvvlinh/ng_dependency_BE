// +build wireinject

package payment_all

import (
	"github.com/google/wire"

	_payment "o.o/backend/cogs/payment"
	_vtpay "o.o/backend/cogs/payment/vtpay"
)

var WireSet = wire.NewSet(
	_payment.WireSet,
	_vtpay.WireSet,
	AllSupportedPaymentProviders,
)
