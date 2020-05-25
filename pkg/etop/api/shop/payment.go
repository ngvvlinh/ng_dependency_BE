package shop

import (
	"context"

	paymentmanager "o.o/api/external/payment/manager"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_source"
	cm "o.o/backend/pkg/common"
)

type PaymentService struct {
	PaymentAggr paymentmanager.CommandBus
}

func (s *PaymentService) Clone() *PaymentService { res := *s; return &res }

func (s *PaymentService) PaymentTradingOrder(ctx context.Context, q *PaymentTradingOrderEndpoint) error {
	if q.OrderId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}
	if q.ReturnUrl == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ReturnURL")
	}

	argGenCode := &paymentmanager.GenerateCodeCommand{
		PaymentSource: payment_source.PaymentSourceOrder,
		ID:            q.OrderId.String(),
	}
	if err := s.PaymentAggr.Dispatch(ctx, argGenCode); err != nil {
		return err
	}
	args := &paymentmanager.BuildUrlConnectPaymentGatewayCommand{
		OrderID:           argGenCode.Result,
		Desc:              q.Desc,
		ReturnURL:         q.ReturnUrl,
		TransactionAmount: q.Amount,
		Provider:          payment_provider.PaymentProvider(q.PaymentProvider),
	}

	if err := s.PaymentAggr.Dispatch(ctx, args); err != nil {
		return err
	}
	q.Result = &shop.PaymentTradingOrderResponse{
		Url: args.Result,
	}
	return nil
}

func (s *PaymentService) PaymentCheckReturnData(ctx context.Context, q *PaymentCheckReturnDataEndpoint) error {
	if q.Id == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Mã giao dịch không được để trống")
	}
	if q.Code == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Mã 'Code' không được để trống")
	}
	args := &paymentmanager.CheckReturnDataCommand{
		ID:                    q.Id,
		Code:                  q.Code,
		PaymentStatus:         q.PaymentStatus,
		Amount:                q.Amount,
		ExternalTransactionID: q.ExternalTransactionId,
		Provider:              payment_provider.PaymentProvider(q.PaymentProvider),
	}
	if err := s.PaymentAggr.Dispatch(ctx, args); err != nil {
		return err
	}
	q.Result = &pbcm.MessageResponse{
		Code: "ok",
		Msg:  args.Result.Msg,
	}
	return nil
}
