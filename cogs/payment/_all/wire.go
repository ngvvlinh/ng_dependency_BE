// +build wireinject

package payment_all

import (
	"github.com/google/wire"

	_payment "o.o/backend/cogs/payment"
	_kpay "o.o/backend/cogs/payment/kpay"
	_vtpay "o.o/backend/cogs/payment/vtpay"
)

var WireSet = wire.NewSet(
	_payment.WireSet,
	_vtpay.WireSet,
	_kpay.WireSet,
	AllSupportedPaymentProviders,
)
