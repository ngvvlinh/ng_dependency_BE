// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package tradering

import (
	context "context"

	shopping "etop.vn/api/shopping"
	capi "etop.vn/capi"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus                          { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus                              { return QueryBus{bus} }
func (c CommandBus) Dispatch(ctx context.Context, msg Command) error { return c.bus.Dispatch(ctx, msg) }
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error     { return c.bus.Dispatch(ctx, msg) }
func (c CommandBus) DispatchAll(ctx context.Context, msgs ...Command) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c QueryBus) DispatchAll(ctx context.Context, msgs ...Query) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

type GetTraderByIDQuery struct {
	ID     int64
	ShopID int64

	Result *ShopTrader `json:"-"`
}

func (h QueryServiceHandler) HandleGetTraderByID(ctx context.Context, msg *GetTraderByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTraderByID(msg.GetArgs(ctx))
	return err
}

type ListTradersByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *TradersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListTradersByIDs(ctx context.Context, msg *ListTradersByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListTradersByIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *GetTraderByIDQuery) query()    {}
func (q *ListTradersByIDsQuery) query() {}

// implement conversion

func (q *GetTraderByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *ListTradersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
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
	b.AddHandler(h.HandleGetTraderByID)
	b.AddHandler(h.HandleListTradersByIDs)
	return QueryBus{b}
}
