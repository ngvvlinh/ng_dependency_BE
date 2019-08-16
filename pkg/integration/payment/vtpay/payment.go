package vtpay

import (
	"context"

	paymentmanager "etop.vn/api/external/payment/manager"
	servicepaymentmanager "etop.vn/backend/com/external/payment/manager"
	vtpayclient "etop.vn/backend/pkg/integration/payment/vtpay/client"
)

var _ servicepaymentmanager.PaymentProvider = &Provider{}

type Provider struct {
	client *vtpayclient.Client
}

func New(cfg vtpayclient.Config) *Provider {
	client := vtpayclient.New(cfg)
	return &Provider{
		client: client,
	}
}

func (p *Provider) BuildUrlConnectPaymentGateway(ctx context.Context, args *paymentmanager.ConnectPaymentGatewayArgs) (string, error) {
	req := &vtpayclient.ConnectPaymentGatewayRequest{
		BillCode:    args.OrderID,
		Desc:        args.Desc,
		OrderID:     args.OrderID,
		ReturnURL:   args.ReturnURL,
		CancelURL:   args.CancelURL,
		TransAmount: args.TransactionAmount,
	}
	return p.client.BuildUrlConnectPaymentGateway(ctx, req)
}

func (p *Provider) GetTransaction(ctx context.Context, args *paymentmanager.GetTransactionArgs) (*paymentmanager.GetTransactionResult, error) {
	req := &vtpayclient.GetTransactionRequest{
		OrderID: args.OrderID,
	}
	res, err := p.client.GetTransaction(ctx, req)
	if err != nil {
		return nil, err
	}

	paymentStatus := vtpayclient.PaymentStatus(res.PaymentStatus)
	return &paymentmanager.GetTransactionResult{
		OrderID:                   res.OrderID,
		ExternalTransactionID:     res.VTTransactionID,
		PaymentState:              paymentStatus.ToState(),
		ExternalPaymentStatus:     res.PaymentStatus,
		ExternalPaymentStatusText: vtpayclient.PaymentStatusMap[paymentStatus],
	}, nil
}

func (p *Provider) CancelTransaction(ctx context.Context, args *paymentmanager.CancelTransactionArgs) (*paymentmanager.CancelTransactionResult, error) {
	req := &vtpayclient.CancelTransactionRequest{
		OrderID:           args.OrderID,
		OriginalRequestID: args.ExternalTransactionID,
		TransAmount:       args.TransactionAmount,
		TransContent:      args.Reason,
	}
	_, err := p.client.CancelTransaction(ctx, req)
	if err != nil {
		return nil, err
	}
	return &paymentmanager.CancelTransactionResult{}, nil
}
