package kpay

import (
	"golang.org/x/net/context"
	paymentmanager "o.o/api/external/payment/manager"
	"o.o/api/top/types/etc/payment_provider"
	servicepaymentmanager "o.o/backend/com/external/payment/manager"
	cm "o.o/backend/pkg/common"
	kpayclient "o.o/backend/pkg/integration/payment/kpay/client"
)

var _ servicepaymentmanager.PaymentProvider = &Provider{}

type Provider struct {
	client *kpayclient.Client
}

func New(cfg kpayclient.Config) *Provider {
	client := kpayclient.New(cfg)
	return &Provider{
		client: client,
	}
}

func (p Provider) Code() payment_provider.PaymentProvider {
	return payment_provider.KPay
}

func (p Provider) BuildUrlConnectPaymentGateway(
	ctx context.Context, args *paymentmanager.ConnectPaymentGatewayArgs,
) (string, error) {
	createTransactionReq := &kpayclient.CreateTransactionRequest{
		Amount:             args.TransactionAmount,
		Desc:               args.Desc,
		PartnerTransaction: args.OrderID,
		FailedURL:          args.CancelURL,
		RedirectURL:        args.ReturnURL,
	}
	createTransactionResp, err := p.client.CreateTransaction(ctx, createTransactionReq)
	if err != nil {
		return "", err
	}

	return createTransactionResp.Url.String(), nil
}

func (p Provider) GetTransaction(
	ctx context.Context, args *paymentmanager.GetTransactionArgs,
) (*paymentmanager.GetTransactionResult, error) {
	getTransactionReq := &kpayclient.GetTransactionRequest{
		Transaction: args.ExternalTransactionID,
	}
	getTransactionResp, err := p.client.GetTransaction(ctx, getTransactionReq)
	if err != nil {
		return nil, err
	}
	return &paymentmanager.GetTransactionResult{
		OrderID:                   getTransactionResp.PartnerTransaction.String(),
		ExternalTransactionID:     getTransactionResp.Transaction.String(),
		PaymentState:              getTransactionResp.State.ToState(),
		ExternalPaymentStatus:     getTransactionResp.State.ToString(),
		ExternalPaymentStatusText: kpayclient.PaymentStatusMap[getTransactionResp.State],
	}, nil
}

func (p Provider) CancelTransaction(
	ctx context.Context, args *paymentmanager.CancelTransactionArgs,
) (*paymentmanager.CancelTransactionResult, error) {
	cancelTransactionReq := &kpayclient.CancelTransactionRequest{
		Transaction: args.ExternalTransactionID,
	}
	if _, err := p.client.CancelTransaction(ctx, cancelTransactionReq); err != nil {
		return nil, err
	}
	return &paymentmanager.CancelTransactionResult{}, nil
}

func (p Provider) CheckReturnData(
	ctx context.Context, args *servicepaymentmanager.CheckReturnDataArgs,
) (*servicepaymentmanager.CheckReturnDataResult, error) {
	paymentStatus := kpayclient.PaymentStatus(args.PaymentStatus)
	if paymentStatus == kpayclient.PaymentStatusSucceeded {
		return &servicepaymentmanager.CheckReturnDataResult{
			Message:                   kpayclient.PaymentStatusMap[paymentStatus],
			PaymentState:              paymentStatus.ToState(),
			PaymentStatus:             paymentStatus.ToStatus(),
			ExternalPaymentStatusText: kpayclient.PaymentStatusMap[paymentStatus],
		}, nil
	}

	return nil, cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ KPay: %v", kpayclient.PaymentStatusMap[paymentStatus])
}
