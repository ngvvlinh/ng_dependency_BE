package vtpay

import (
	"context"

	"etop.vn/api/external/payment"
	"etop.vn/api/external/payment/vtpay"
	"etop.vn/api/main/ordering"
	"etop.vn/api/top/types/etc/payment_provider"
	"etop.vn/api/top/types/etc/payment_source"
	paymentutil "etop.vn/backend/com/external/payment"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	vtpayclient "etop.vn/backend/pkg/integration/payment/vtpay/client"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type Aggregate struct {
	db          cmsql.Transactioner
	orderQS     ordering.QueryBus
	orderAggr   ordering.CommandBus
	vtpayClient *vtpayclient.Client
	paymentAggr payment.CommandBus
}

func NewAggregate(db *cmsql.Database, orderQuery ordering.QueryBus, orderA ordering.CommandBus, paymentA payment.CommandBus, vtpayClient *vtpayclient.Client) *Aggregate {
	return &Aggregate{
		db:          db,
		orderQS:     orderQuery,
		orderAggr:   orderA,
		paymentAggr: paymentA,
		vtpayClient: vtpayClient,
	}
}

func (a *Aggregate) MessageBus() vtpay.CommandBus {
	b := bus.New()
	return vtpay.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) HandleExternalDataResponse(ctx context.Context, args *vtpay.HandleExternalDataResponseArgs) error {
	paymentSource, id, err := paymentutil.ParsePaymentCode(args.OrderID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Mã giao dịch không hợp lệ (order_id = %v)", args.OrderID)
	}
	switch paymentSource {
	case payment_source.PaymentSourceOrder:
		return a.HandleExternalDataOrderResponse(ctx, id, args)
	default:
		return cm.Errorf(cm.InvalidArgument, err, "Mã giao dịch không hợp lệ (order_id = %v)", args.OrderID)
	}
}

func (a *Aggregate) HandleExternalDataOrderResponse(ctx context.Context, orderID dot.ID, args *vtpay.HandleExternalDataResponseArgs) error {
	queryOrder := &ordering.GetOrderByIDQuery{ID: orderID}
	if err := a.orderQS.Dispatch(ctx, queryOrder); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Mã đơn hàng không tồn tại (order_id = %v)", orderID)
	}
	order := queryOrder.Result

	paymentStatus := vtpayclient.PaymentStatus(args.PaymentStatus)
	data, _ := jsonx.Marshal(args)

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		cmd := &payment.CreateOrUpdatePaymentCommand{
			Amount:          args.TransAmount,
			Status:          paymentStatus.ToStatus(),
			State:           paymentStatus.ToState(),
			PaymentProvider: payment_provider.VTPay,
			ExternalTransID: args.VtTransactionID,
			ExternalData:    data,
		}
		if err := a.paymentAggr.Dispatch(ctx, cmd); err != nil {
			return cm.Errorf(cm.Internal, err, "Không thể tạo payment")
		}
		paymentData := cmd.Result

		updateOrder := &ordering.UpdateOrderPaymentInfoCommand{
			ID:            order.ID,
			PaymentStatus: paymentData.Status,
			PaymentID:     paymentData.ID,
		}
		if err := a.orderAggr.Dispatch(ctx, updateOrder); err != nil {
			return cm.Errorf(cm.Internal, err, "Cập nhật trạng thái thanh toán đơn hàng thất bại")
		}
		return nil
	})
}
