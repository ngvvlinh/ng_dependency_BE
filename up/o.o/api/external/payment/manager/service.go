package manager

import (
	"context"

	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_source"
	"o.o/api/top/types/etc/payment_state"
)

// +gen:api

type Aggregate interface {
	BuildUrlConnectPaymentGateway(context.Context, *ConnectPaymentGatewayArgs) (string, error)

	GenerateCode(context.Context, *GenerateCodeArgs) (code string, err error)

	GetTransaction(context.Context, *GetTransactionArgs) (*GetTransactionResult, error)

	CancelTransaction(context.Context, *CancelTransactionArgs) (*CancelTransactionResult, error)

	// Kiểm trả kết trả trả về dựa vào params trên redirect_urls của bên thứ 3 (vtpay)
	CheckReturnData(context.Context, *CheckReturnDataArgs) (*CheckReturnDataResult, error)
}

type ConnectPaymentGatewayArgs struct {
	// Mã giao dịch ETOP
	OrderID           string
	Desc              string
	ReturnURL         string
	CancelURL         string
	TransactionAmount int
	Provider          payment_provider.PaymentProvider
}

type GetTransactionArgs struct {
	OrderID  string
	Provider payment_provider.PaymentProvider
}

type GetTransactionResult struct {
	OrderID                   string
	ExternalTransactionID     string
	PaymentState              payment_state.PaymentState
	ExternalPaymentStatus     string
	ExternalPaymentStatusText string
}

type CancelTransactionArgs struct {
	OrderID               string
	ExternalTransactionID string
	TransactionAmount     int
	Reason                string
	Provider              payment_provider.PaymentProvider
}

type CancelTransactionResult struct {
}

type GenerateCodeArgs struct {
	PaymentSource payment_source.PaymentSource
	ID            string
}

type CheckReturnDataArgs struct {
	ID                    string
	Code                  string
	PaymentStatus         string
	Amount                int
	ExternalTransactionID string
	Provider              payment_provider.PaymentProvider
}

type CheckReturnDataResult struct {
	Code string
	Msg  string
}
