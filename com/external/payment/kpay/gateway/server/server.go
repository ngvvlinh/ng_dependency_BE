package server

import (
	"golang.org/x/net/context"
	"o.o/api/external/payment"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/status4"
	paymentlogaggregate "o.o/backend/com/etc/logging/payment/aggregate"
	paymentlogmodel "o.o/backend/com/etc/logging/payment/model"
	"o.o/backend/com/external/payment/kpay/gateway"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

type Server struct {
	paymentLogAggr *paymentlogaggregate.Aggregate
	paymentAggr    payment.CommandBus
	paymentQuery   payment.QueryBus
}

func New(
	paymentLogA *paymentlogaggregate.Aggregate,
	paymentA payment.CommandBus,
	paymentQ payment.QueryBus,
) *Server {
	return &Server{
		paymentLogAggr: paymentLogA,
		paymentAggr:    paymentA,
		paymentQuery:   paymentQ,
	}
}

func (s *Server) SaveLog(
	ctx context.Context, data interface{},
	orderID string, action paymentlogmodel.PaymentAction,
) error {
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

func (s *Server) Callback(c *httpx.Context) (_err error) {
	ctx := c.Req.Context()
	var msg gateway.CallbackTransaction
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "KPay: can not decode JSON callback")
	}
	_data, err := jsonx.Marshal(msg)
	if err != nil {
		return err
	}

	defer func() {
		_ = s.SaveLog(ctx, msg, msg.PartnerTransaction.String(), paymentlogmodel.PaymentActionUnknown)

		if _err == nil {
			c.SetResult(map[string]string{"code": "ok"})
		}
	}()

	getPaymentQuery := &payment.GetPaymentByExternalTransIDQuery{
		TransactionID: msg.Transaction.String(),
	}
	if err := s.paymentQuery.Dispatch(ctx, getPaymentQuery); err != nil {
		return err
	}
	oldPayment := getPaymentQuery.Result

	// không update payment khi payment đã về trạng thái cuối
	if oldPayment.Status == status4.N || oldPayment.Status == status4.P {
		return nil
	}

	updatePaymentCmd := &payment.UpdateExternalPaymentInfoCommand{
		ID:              oldPayment.ID,
		Status:          msg.State.ToStatus(),
		State:           msg.State.ToState(),
		ExternalData:    _data,
		ExternalTransID: msg.Transaction.String(),
	}
	if err := s.paymentAggr.Dispatch(ctx, updatePaymentCmd); err != nil {
		return err
	}

	return nil
}
