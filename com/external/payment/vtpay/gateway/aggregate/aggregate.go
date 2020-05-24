package aggregate

import (
	"context"
	"strconv"

	"o.o/api/external/payment/vtpay"
	vtpaygateway "o.o/api/external/payment/vtpay/gateway"
	"o.o/api/main/ordering"
	"o.o/api/top/types/etc/payment_source"
	"o.o/api/top/types/etc/status4"
	paymentutil "o.o/backend/com/external/payment"
	"o.o/backend/pkg/common/bus"
	vtpayclient "o.o/backend/pkg/integration/payment/vtpay/client"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	_  vtpaygateway.Aggregate = &Aggregate{}
	ll                        = l.New()
)

const (
	PathValidateTransaction = "ValidateTransaction"
	PathGetResult           = "GetResult"
)

func BuildGatewayRoute(path string) string {
	return "/payment/vtpay/gateway/" + path
}

type Aggregate struct {
	orderQS     ordering.QueryBus
	orderAggr   ordering.CommandBus
	vtpayAggr   vtpay.CommandBus
	vtpayClient *vtpayclient.Client
}

func NewAggregate(orderQuery ordering.QueryBus, orderAggregate ordering.CommandBus, vtpayA vtpay.CommandBus, vtpayClient *vtpayclient.Client) *Aggregate {
	return &Aggregate{
		orderQS:     orderQuery,
		orderAggr:   orderAggregate,
		vtpayAggr:   vtpayA,
		vtpayClient: vtpayClient,
	}
}

func AggregateMessageBus(a *Aggregate) vtpaygateway.CommandBus {
	b := bus.New()
	return vtpaygateway.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) ValidateTransaction(ctx context.Context, args *vtpaygateway.ValidateTransactionArgs) (*vtpaygateway.ValidateTransactionResult, error) {
	paymentSource, id, err := paymentutil.ParsePaymentCode(args.OrderID)
	if err != nil {
		return &vtpaygateway.ValidateTransactionResult{
			ErrorCode: vtpaygateway.ErrorCodeInvalidData,
		}, nil
	}
	switch paymentSource {
	case payment_source.PaymentSourceOrder:
		return a.HandleValiDateTransactionOrder(ctx, id, args)
	default:
		return &vtpaygateway.ValidateTransactionResult{
			ErrorCode: vtpaygateway.ErrorCodeInvalidData,
		}, nil
	}
}

func (a *Aggregate) HandleValiDateTransactionOrder(ctx context.Context, orderID dot.ID, args *vtpaygateway.ValidateTransactionArgs) (*vtpaygateway.ValidateTransactionResult, error) {
	queryOrder := &ordering.GetOrderByIDQuery{ID: orderID}
	if err := a.orderQS.Dispatch(ctx, queryOrder); err != nil {
		return nil, err
	}
	order := queryOrder.Result
	if isOrderPaid(order) {
		return &vtpaygateway.ValidateTransactionResult{
			ErrorCode: vtpaygateway.ErrorCodeInternal,
		}, nil
	}
	dataCheckSum := args.DataCheckSum(order.TotalAmount)
	checkSum := a.vtpayClient.CheckSum(dataCheckSum)
	ll.Info("checksum", l.String("data checksum", dataCheckSum), l.String("checkSum", checkSum))

	if checkSum != args.CheckSum {
		return &vtpaygateway.ValidateTransactionResult{
			ErrorCode: vtpaygateway.ErrorCodeCheckSum,
		}, nil
	}
	res := &vtpaygateway.ValidateTransactionResult{
		BillCode:     args.BillCode,
		MerchantCode: a.vtpayClient.MerchantCode,
		OrderID:      args.OrderID,
		TransAmount:  strconv.Itoa(order.TotalAmount),
		ErrorCode:    vtpaygateway.ErrorCodeSuccess,
	}
	checkSumRes := a.vtpayClient.CheckSum(res.DataCheckSum())
	res.CheckSum = checkSumRes
	ll.Info("res ValidateTransaction", l.Any("res", res))
	return res, nil
}

func isOrderPaid(order *ordering.Order) bool {
	return order.PaymentStatus == status4.P
}

func (a *Aggregate) GetResult(ctx context.Context, args *vtpaygateway.GetResultArgs) (*vtpaygateway.GetResultResult, error) {
	cmd := &vtpay.HandleExternalDataResponseCommand{
		BillCode:        args.BillCode,
		CustMsisdn:      args.CustMsisdn,
		ErrorCode:       args.ErrorCode,
		MerchantCode:    args.MerchantCode,
		OrderID:         args.OrderID,
		PaymentStatus:   args.PaymentStatus,
		TransAmount:     args.TransAmount,
		VtTransactionID: args.VtTransactionID,
		CheckSum:        args.CheckSum,
	}

	err := a.vtpayAggr.Dispatch(ctx, cmd)
	if err != nil {
		return &vtpaygateway.GetResultResult{
			ErrorCode: vtpaygateway.ErrorCodeInvalidData,
		}, nil
	}

	res := &vtpaygateway.GetResultResult{
		ErrorCode:    vtpaygateway.ErrorCodeSuccess,
		MerchantCode: a.vtpayClient.MerchantCode,
		OrderID:      args.OrderID,
		CheckSum:     "",
	}
	checkSumRes := a.vtpayClient.CheckSum(res.DataCheckSum())
	res.CheckSum = checkSumRes
	ll.Info("res GetResult", l.Any("res", res))
	return res, nil
}
