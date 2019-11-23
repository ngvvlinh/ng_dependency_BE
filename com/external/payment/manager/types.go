package manager

import (
	"context"

	paymentmanager "etop.vn/api/external/payment/manager"
	"etop.vn/api/main/etop"
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
	PaymentState              paymentmanager.PaymentState
	PaymentStatus             etop.Status4
	ExternalPaymentStatusText string
}
