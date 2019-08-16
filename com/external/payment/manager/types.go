package manager

import (
	"context"

	paymentmanager "etop.vn/api/external/payment/manager"
)

type PaymentProvider interface {
	BuildUrlConnectPaymentGateway(context.Context, *paymentmanager.ConnectPaymentGatewayArgs) (string, error)

	GetTransaction(context.Context, *paymentmanager.GetTransactionArgs) (*paymentmanager.GetTransactionResult, error)

	CancelTransaction(context.Context, *paymentmanager.CancelTransactionArgs) (*paymentmanager.CancelTransactionResult, error)
}
