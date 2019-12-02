// +build !generator

// Code generated by generator api. DO NOT EDIT.

package summary

import (
	context "context"
	time "time"

	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
)

type QueryBus struct{ bus capi.Bus }

func NewQueryBus(bus capi.Bus) QueryBus { return QueryBus{bus} }

func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type SummaryPOSQuery struct {
	DateFrom time.Time
	DateTo   time.Time
	ShopID   dot.ID

	Result *SummaryPOSResponse `json:"-"`
}

func (h QueryServiceHandler) HandleSummaryPOS(ctx context.Context, msg *SummaryPOSQuery) (err error) {
	msg.Result, err = h.inner.SummaryPOS(msg.GetArgs(ctx))
	return err
}

type SummaryTopShipQuery struct {
	DateFrom time.Time
	DateTo   time.Time
	ShopID   dot.ID

	Result *SummaryTopShipResponse `json:"-"`
}

func (h QueryServiceHandler) HandleSummaryTopShip(ctx context.Context, msg *SummaryTopShipQuery) (err error) {
	msg.Result, err = h.inner.SummaryTopShip(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *SummaryPOSQuery) query()     {}
func (q *SummaryTopShipQuery) query() {}

// implement conversion

func (q *SummaryPOSQuery) GetArgs(ctx context.Context) (_ context.Context, _ *SummaryPOSRequest) {
	return ctx,
		&SummaryPOSRequest{
			DateFrom: q.DateFrom,
			DateTo:   q.DateTo,
			ShopID:   q.ShopID,
		}
}

func (q *SummaryPOSQuery) SetSummaryPOSRequest(args *SummaryPOSRequest) {
	q.DateFrom = args.DateFrom
	q.DateTo = args.DateTo
	q.ShopID = args.ShopID
}

func (q *SummaryTopShipQuery) GetArgs(ctx context.Context) (_ context.Context, _ *SummaryTopShipRequest) {
	return ctx,
		&SummaryTopShipRequest{
			DateFrom: q.DateFrom,
			DateTo:   q.DateTo,
			ShopID:   q.ShopID,
		}
}

func (q *SummaryTopShipQuery) SetSummaryTopShipRequest(args *SummaryTopShipRequest) {
	q.DateFrom = args.DateFrom
	q.DateTo = args.DateTo
	q.ShopID = args.ShopID
}

// implement dispatching

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleSummaryPOS)
	b.AddHandler(h.HandleSummaryTopShip)
	return QueryBus{b}
}
