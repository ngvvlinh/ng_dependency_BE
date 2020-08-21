package payment_all

import (
	"o.o/backend/com/external/payment/manager"
	"o.o/backend/pkg/integration/payment/vtpay"
)

func AllSupportedPaymentProviders(
	vtpayProvider *vtpay.Provider,
) []manager.PaymentProvider {
	return []manager.PaymentProvider{
		vtpayProvider,
	}
}
