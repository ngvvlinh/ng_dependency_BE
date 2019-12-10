package server

import (
	"context"

	vtpaygateway "etop.vn/api/external/payment/vtpay/gateway"
	"etop.vn/api/top/types/etc/payment_provider"
	paymentlogaggregate "etop.vn/backend/com/etc/logging/payment/aggregate"
	paymentlogmodel "etop.vn/backend/com/etc/logging/payment/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

type Server struct {
	vtpayAggr      vtpaygateway.CommandBus
	paymentLogAggr *paymentlogaggregate.Aggregate
}

var ll = l.New()

func New(vtpayAggregate vtpaygateway.CommandBus, paymentLogA *paymentlogaggregate.Aggregate) *Server {
	return &Server{
		vtpayAggr:      vtpayAggregate,
		paymentLogAggr: paymentLogA,
	}
}

func (s *Server) SaveLog(ctx context.Context, data interface{}, orderID string, action paymentlogmodel.PaymentAction) error {
	id := cm.NewID()
	_data, err := jsonx.Marshal(data)
	if err != nil {
		return err
	}
	paymentLog := &paymentlogmodel.Payment{
		ID:              id,
		Data:            _data,
		OrderID:         orderID,
		Action:          action,
		PaymentProvider: payment_provider.VTPay,
	}
	if err := s.paymentLogAggr.CreatePaymentLog(ctx, paymentLog); err != nil {
		return err
	}
	return nil
}

func (s *Server) ValidateTransaction(c *httpx.Context) error {
	ctx := c.Req.Context()
	var cmd vtpaygateway.ValidateTransactionCommand
	if err := c.DecodeFormUrlEncoded(&cmd); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "VTPay: Can not decode form data in ValidateTransaction")
	}
	{
		_ = s.SaveLog(ctx, cmd, cmd.OrderID, paymentlogmodel.PaymentActionValidate)
	}

	if err := s.vtpayAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	c.SetResult(cmd.Result)
	return nil
}

func (s *Server) GetResult(c *httpx.Context) error {
	ctx := c.Req.Context()
	var cmd vtpaygateway.GetResultCommand
	if err := c.DecodeFormUrlEncoded(&cmd); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "VTPay: Can not decode form data in GetResult")
	}
	{
		_ = s.SaveLog(ctx, cmd, cmd.OrderID, paymentlogmodel.PaymentActionResult)
	}
	if err := s.vtpayAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	c.SetResult(cmd.Result)
	return nil
}
