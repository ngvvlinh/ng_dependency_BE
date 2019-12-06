// +build !generator

// Code generated by generator api. DO NOT EDIT.

package vtpay

import (
	context "context"

	capi "etop.vn/capi"
)

type CommandBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type HandleExternalDataResponseCommand struct {
	BillCode        string `json:"bill_code"`
	CustMsisdn      string `json:"cust_msisdn"`
	ErrorCode       string `json:"error_code"`
	MerchantCode    string `json:"merchant_code"`
	OrderID         string `json:"order_id"`
	PaymentStatus   string `json:"payment_status"`
	TransAmount     int    `json:"trans_amount"`
	VtTransactionID string `json:"vt_transaction_id"`
	CheckSum        string `json:"check_sum"`

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleHandleExternalDataResponse(ctx context.Context, msg *HandleExternalDataResponseCommand) (err error) {
	return h.inner.HandleExternalDataResponse(msg.GetArgs(ctx))
}

// implement interfaces

func (q *HandleExternalDataResponseCommand) command() {}

// implement conversion

func (q *HandleExternalDataResponseCommand) GetArgs(ctx context.Context) (_ context.Context, _ *HandleExternalDataResponseArgs) {
	return ctx,
		&HandleExternalDataResponseArgs{
			BillCode:        q.BillCode,
			CustMsisdn:      q.CustMsisdn,
			ErrorCode:       q.ErrorCode,
			MerchantCode:    q.MerchantCode,
			OrderID:         q.OrderID,
			PaymentStatus:   q.PaymentStatus,
			TransAmount:     q.TransAmount,
			VtTransactionID: q.VtTransactionID,
			CheckSum:        q.CheckSum,
		}
}

func (q *HandleExternalDataResponseCommand) SetHandleExternalDataResponseArgs(args *HandleExternalDataResponseArgs) {
	q.BillCode = args.BillCode
	q.CustMsisdn = args.CustMsisdn
	q.ErrorCode = args.ErrorCode
	q.MerchantCode = args.MerchantCode
	q.OrderID = args.OrderID
	q.PaymentStatus = args.PaymentStatus
	q.TransAmount = args.TransAmount
	q.VtTransactionID = args.VtTransactionID
	q.CheckSum = args.CheckSum
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleHandleExternalDataResponse)
	return CommandBus{b}
}
