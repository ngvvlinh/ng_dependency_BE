package manager

import (
	"context"

	paymentmanager "etop.vn/api/external/payment/manager"
	"etop.vn/api/top/types/etc/payment_state"
	"etop.vn/api/top/types/etc/status4"
)

type PaymentProvider interface {
	BuildUrlConnectPaymentGateway(context.Context, *paymentmanager.ConnectPaymentGatewayArgs) (string, error)

	GetTransaction(context.Context, *paymentmanager.GetTransactionArgs) (*paymentmanager.GetTransactionResult, error)

	CancelTransaction(context.Context, *paymentmanager.CancelTransactionArgs) (*paymentmanager.CancelTransactionResult, error)

	CheckReturnData(context.Context, *CheckReturnDataArgs) (*CheckReturnDataResult, error)
}

type CheckReturnDataArgs struct {
	OrderID       string
	Code          string
	PaymentStatus string
}

type CheckReturnDataResult struct {
	Message                   string
	PaymentState              payment_state.PaymentState
	PaymentStatus             status4.Status
	ExternalPaymentStatusText string
}
