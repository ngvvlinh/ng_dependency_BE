package vtpay

import (
	"context"

	paymentmanager "o.o/api/external/payment/manager"
	"o.o/api/top/types/etc/payment_provider"
	servicepaymentmanager "o.o/backend/com/external/payment/manager"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	vtpayclient "o.o/backend/pkg/integration/payment/vtpay/client"
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

func (p *Provider) Code() payment_provider.PaymentProvider {
	return payment_provider.VTPay
}

func (p *Provider) BuildUrlConnectPaymentGateway(
	ctx context.Context, args *paymentmanager.ConnectPaymentGatewayArgs,
) (string, error) {
	req := &vtpayclient.ConnectPaymentGatewayRequest{
		BillCode:    args.OrderID,
		Desc:        args.Desc,
		OrderID:     args.OrderID,
		ReturnURL:   args.ReturnURL,
		CancelURL:   args.CancelURL,
		TransAmount: args.TransactionAmount,
	}
	url, err := p.client.BuildUrlConnectPaymentGateway(ctx, req)
	if err != nil {
		return "", err
	}
	// TODO(ngoc): get externalTransactionID
	return url, nil
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

func (p *Provider) CheckReturnData(ctx context.Context, args *servicepaymentmanager.CheckReturnDataArgs) (*servicepaymentmanager.CheckReturnDataResult, error) {
	errMsg := vtpayclient.TransactionResultCodeMap[args.Code]
	paymentStatus := vtpayclient.PaymentStatus(args.PaymentStatus)
	if args.Code == vtpayclient.TransactionSuccessCode {
		return &servicepaymentmanager.CheckReturnDataResult{
			Message:                   errMsg,
			PaymentState:              paymentStatus.ToState(),
			PaymentStatus:             paymentStatus.ToStatus(),
			ExternalPaymentStatusText: vtpayclient.PaymentStatusMap[paymentStatus],
		}, nil
	}

	if errMsg == "" {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Mã lỗi '%v' của VTPay không hợp lệ. Vui lòng liên hệ %v để biết thêm chi tiết.", args.Code, wl.X(ctx).CSEmail)
	}
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ VTPay: %v", errMsg)
}
