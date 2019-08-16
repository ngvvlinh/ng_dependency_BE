package manager

import (
	"context"

	"etop.vn/api/external/payment"
)

type PaymentState = payment.PaymentState

type Aggregate interface {
	BuildUrlConnectPaymentGateway(context.Context, *ConnectPaymentGatewayArgs) (string, error)

	GetTransaction(context.Context, *GetTransactionArgs) (*GetTransactionResult, error)

	CancelTransaction(context.Context, *CancelTransactionArgs) (*CancelTransactionResult, error)
}

type ConnectPaymentGatewayArgs struct {
	// Mã giao dịch ETOP
	OrderID           string
	Desc              string
	ReturnURL         string
	CancelURL         string
	TransactionAmount int
	Provider          payment.PaymentProvider
}

type GetTransactionArgs struct {
	OrderID  string
	Provider payment.PaymentProvider
}

type GetTransactionResult struct {
	OrderID                   string
	ExternalTransactionID     string
	PaymentState              PaymentState
	ExternalPaymentStatus     string
	ExternalPaymentStatusText string
}

type CancelTransactionArgs struct {
	OrderID               string
	ExternalTransactionID string
	TransactionAmount     int
	Reason                string
	Provider              payment.PaymentProvider
}

type CancelTransactionResult struct {
}
