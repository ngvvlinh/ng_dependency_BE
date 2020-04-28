package manager

import (
	"context"
	"fmt"

	paymentmanager "o.o/api/external/payment/manager"
	"o.o/api/main/ordering"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_source"
	paymentutil "o.o/backend/com/external/payment"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ paymentmanager.Aggregate = &PaymentManager{}

type PaymentManager struct {
	vtpay   PaymentProvider
	orderQS ordering.QueryBus
}

func NewManager(viettelPay PaymentProvider, orderQuery ordering.QueryBus) *PaymentManager {
	return &PaymentManager{
		vtpay:   viettelPay,
		orderQS: orderQuery,
	}
}

func (ctrl *PaymentManager) MesssageBus() paymentmanager.CommandBus {
	b := bus.New()
	return paymentmanager.NewAggregateHandler(ctrl).RegisterHandlers(b)
}

func (ctrl *PaymentManager) GetPaymentProviderDriver(provider payment_provider.PaymentProvider) (PaymentProvider, error) {
	switch provider {
	case payment_provider.VTPay:
		return ctrl.vtpay, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phương thức thanh toán không hợp lệ (%v)", provider)
	}
}

func (ctrl *PaymentManager) GenerateCode(ctx context.Context, args *paymentmanager.GenerateCodeArgs) (string, error) {
	switch args.PaymentSource {
	case payment_source.PaymentSourceOrder:
		// nothing
	default:
		return "", cm.Errorf(cm.InvalidArgument, nil, "PaymentSource không hợp lệ. Vui lòng kiểm tra lại.")
	}
	if args.ID == "" {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	return fmt.Sprintf("%v_%v", args.PaymentSource.String(), args.ID), nil
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

func (ctrl *PaymentManager) CheckReturnData(ctx context.Context, args *paymentmanager.CheckReturnDataArgs) (*paymentmanager.CheckReturnDataResult, error) {
	if args.Code == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "'Code' không được để trống")
	}
	paymentSource, id, err := paymentutil.ParsePaymentCode(args.ID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Mã giao dịch không hợp lệ (order_id = %v)", args.ID)
	}
	switch paymentSource {
	case payment_source.PaymentSourceOrder:
		return ctrl.HandlerCheckReturnOrderData(ctx, id, args)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, err, "Mã giao dịch không hợp lệ (order_id = %v)", args.ID)
	}
}

func (ctrl *PaymentManager) HandlerCheckReturnOrderData(ctx context.Context, orderID dot.ID, args *paymentmanager.CheckReturnDataArgs) (*paymentmanager.CheckReturnDataResult, error) {
	queryOrder := &ordering.GetOrderByIDQuery{ID: orderID}
	if err := ctrl.orderQS.Dispatch(ctx, queryOrder); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Mã đơn hàng không tồn tại (order_id = %v)", orderID)
	}

	provider, err := ctrl.GetPaymentProviderDriver(args.Provider)
	if err != nil {
		return nil, err
	}

	_args := &CheckReturnDataArgs{
		OrderID:       args.ID,
		Code:          args.Code,
		PaymentStatus: args.PaymentStatus,
	}
	result, err := provider.CheckReturnData(ctx, _args)
	if err != nil {
		return nil, err
	}
	return &paymentmanager.CheckReturnDataResult{
		Code: "ok",
		Msg:  result.Message,
	}, nil
}
