package manager

import (
	"context"

	"etop.vn/common/bus"

	"etop.vn/api/external/payment"
	paymentmanager "etop.vn/api/external/payment/manager"
	cm "etop.vn/backend/pkg/common"
)

var _ paymentmanager.Aggregate = &PaymentManager{}

type PaymentManager struct {
	vtpay PaymentProvider
}

func NewManager(viettelPay PaymentProvider) *PaymentManager {
	return &PaymentManager{
		vtpay: viettelPay,
	}
}

func (ctrl *PaymentManager) MesssageBus() paymentmanager.CommandBus {
	b := bus.New()
	return paymentmanager.NewAggregateHandler(ctrl).RegisterHandlers(b)
}

func (ctrl *PaymentManager) GetPaymentProviderDriver(provider payment.PaymentProvider) (PaymentProvider, error) {
	switch provider {
	case payment.PaymentProviderVTPay:
		return ctrl.vtpay, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phương thức thanh toán không hợp lệ (%v)", provider)
	}
}

func (ctrl *PaymentManager) BuildUrlConnectPaymentGateway(ctx context.Context, args *paymentmanager.ConnectPaymentGatewayArgs) (string, error) {
	provider, err := ctrl.GetPaymentProviderDriver(args.Provider)
	if err != nil {
		return "", err
	}
	return provider.BuildUrlConnectPaymentGateway(ctx, args)
}

func (ctrl *PaymentManager) GetTransaction(ctx context.Context, args *paymentmanager.GetTransactionArgs) (*paymentmanager.GetTransactionResult, error) {
	provider, err := ctrl.GetPaymentProviderDriver(args.Provider)
	if err != nil {
		return nil, err
	}
	return provider.GetTransaction(ctx, args)
}

func (ctrl *PaymentManager) CancelTransaction(ctx context.Context, args *paymentmanager.CancelTransactionArgs) (*paymentmanager.CancelTransactionResult, error) {
	provider, err := ctrl.GetPaymentProviderDriver(args.Provider)
	if err != nil {
		return nil, err
	}
	return provider.CancelTransaction(ctx, args)
}
