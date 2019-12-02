// +build !generator

// Code generated by generator api. DO NOT EDIT.

package address

import (
	context "context"

	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
)

type QueryBus struct{ bus capi.Bus }

func NewQueryBus(bus capi.Bus) QueryBus { return QueryBus{bus} }

func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type GetAddressByIDQuery struct {
	ID dot.ID

	Result *Address `json:"-"`
}

func (h QueryServiceHandler) HandleGetAddressByID(ctx context.Context, msg *GetAddressByIDQuery) (err error) {
	msg.Result, err = h.inner.GetAddressByID(msg.GetArgs(ctx))
	return err
}

// implement interfaces
func (q *GetAddressByIDQuery) query() {}

// implement conversion

func (q *GetAddressByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAddressByIDQueryArgs) {
	return ctx,
		&GetAddressByIDQueryArgs{
			ID: q.ID,
		}
}

func (q *GetAddressByIDQuery) SetGetAddressByIDQueryArgs(args *GetAddressByIDQueryArgs) {
	q.ID = args.ID
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
	b.AddHandler(h.HandleGetAddressByID)
	return QueryBus{b}
}
