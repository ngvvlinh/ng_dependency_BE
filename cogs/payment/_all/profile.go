package payment_all

import (
	"github.com/google/wire"

	_payment "o.o/backend/cogs/payment"
	_vtpay "o.o/backend/cogs/payment/vtpay"
	"o.o/backend/com/external/payment/manager"
	"o.o/backend/pkg/integration/payment/vtpay"
)

var WireSet = wire.NewSet(
	_payment.WireSet,
	_vtpay.WireSet,
	AllSupportedPaymentProviders,
)

func AllSupportedPaymentProviders(
	vtpayProvider *vtpay.Provider,
) []manager.PaymentProvider {
	return []manager.PaymentProvider{
		vtpayProvider,
	}
}
