// +build !generator

// Code generated by generator api. DO NOT EDIT.

package tradering

import (
	context "context"

	meta "etop.vn/api/meta"
	shopping "etop.vn/api/shopping"
	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type DeleteTraderCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteTrader(ctx context.Context, msg *DeleteTraderCommand) (err error) {
	msg.Result, err = h.inner.DeleteTrader(msg.GetArgs(ctx))
	return err
}

type GetTraderByIDQuery struct {
	ID             dot.ID
	ShopID         dot.ID
	IncludeDeleted bool

	Result *ShopTrader `json:"-"`
}

func (h QueryServiceHandler) HandleGetTraderByID(ctx context.Context, msg *GetTraderByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTraderByID(msg.GetArgs(ctx))
	return err
}

type GetTraderInfoByIDQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *ShopTrader `json:"-"`
}

func (h QueryServiceHandler) HandleGetTraderInfoByID(ctx context.Context, msg *GetTraderInfoByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTraderInfoByID(msg.GetArgs(ctx))
	return err
}

type ListTradersByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID
	Paging meta.Paging

	Result *TradersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListTradersByIDs(ctx context.Context, msg *ListTradersByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListTradersByIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *DeleteTraderCommand) command() {}

func (q *GetTraderByIDQuery) query()     {}
func (q *GetTraderInfoByIDQuery) query() {}
func (q *ListTradersByIDsQuery) query()  {}

// implement conversion

func (q *DeleteTraderCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, shopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *GetTraderByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:             q.ID,
			ShopID:         q.ShopID,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *GetTraderByIDQuery) SetIDQueryShopArg(args *shopping.IDQueryShopArg) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.IncludeDeleted = args.IncludeDeleted
}

func (q *GetTraderInfoByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, ShopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *ListTradersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
			Paging: q.Paging,
		}
}

func (q *ListTradersByIDsQuery) SetIDsQueryShopArgs(args *shopping.IDsQueryShopArgs) {
	q.IDs = args.IDs
	q.ShopID = args.ShopID
	q.Paging = args.Paging
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
	b.AddHandler(h.HandleDeleteTrader)
	return CommandBus{b}
}

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
	b.AddHandler(h.HandleGetTraderInfoByID)
	b.AddHandler(h.HandleListTradersByIDs)
	return QueryBus{b}
}
