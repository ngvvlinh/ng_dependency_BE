// +build !generator

// Code generated by generator api. DO NOT EDIT.

package refund

import (
	context "context"

	meta "o.o/api/meta"
	inttypes "o.o/api/top/int/types"
	inventory_auto "o.o/api/top/types/etc/inventory_auto"
	capi "o.o/capi"
	dot "o.o/capi/dot"
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

type CancelRefundCommand struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	CancelReason         string
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher

	Result *Refund `json:"-"`
}

func (h AggregateHandler) HandleCancelRefund(ctx context.Context, msg *CancelRefundCommand) (err error) {
	msg.Result, err = h.inner.CancelRefund(msg.GetArgs(ctx))
	return err
}

type ConfirmRefundCommand struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher

	Result *Refund `json:"-"`
}

func (h AggregateHandler) HandleConfirmRefund(ctx context.Context, msg *ConfirmRefundCommand) (err error) {
	msg.Result, err = h.inner.ConfirmRefund(msg.GetArgs(ctx))
	return err
}

type CreateRefundCommand struct {
	Lines           []*RefundLine
	OrderID         dot.ID
	AdjustmentLines []*inttypes.AdjustmentLine
	TotalAdjustment int
	TotalAmount     int
	BasketValue     int
	ShopID          dot.ID
	CreatedBy       dot.ID
	Note            string

	Result *Refund `json:"-"`
}

func (h AggregateHandler) HandleCreateRefund(ctx context.Context, msg *CreateRefundCommand) (err error) {
	msg.Result, err = h.inner.CreateRefund(msg.GetArgs(ctx))
	return err
}

type UpdateRefundCommand struct {
	Lines           []*RefundLine
	ID              dot.ID
	ShopID          dot.ID
	AdjustmentLines []*inttypes.AdjustmentLine
	TotalAdjustment dot.NullInt
	TotalAmount     dot.NullInt
	BasketValue     dot.NullInt
	UpdateBy        dot.ID
	Note            dot.NullString

	Result *Refund `json:"-"`
}

func (h AggregateHandler) HandleUpdateRefund(ctx context.Context, msg *UpdateRefundCommand) (err error) {
	msg.Result, err = h.inner.UpdateRefund(msg.GetArgs(ctx))
	return err
}

type GetRefundByIDQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *Refund `json:"-"`
}

func (h QueryServiceHandler) HandleGetRefundByID(ctx context.Context, msg *GetRefundByIDQuery) (err error) {
	msg.Result, err = h.inner.GetRefundByID(msg.GetArgs(ctx))
	return err
}

type GetRefundsQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *GetRefundsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetRefunds(ctx context.Context, msg *GetRefundsQuery) (err error) {
	msg.Result, err = h.inner.GetRefunds(msg.GetArgs(ctx))
	return err
}

type GetRefundsByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID

	Result []*Refund `json:"-"`
}

func (h QueryServiceHandler) HandleGetRefundsByIDs(ctx context.Context, msg *GetRefundsByIDsQuery) (err error) {
	msg.Result, err = h.inner.GetRefundsByIDs(msg.GetArgs(ctx))
	return err
}

type GetRefundsByOrderIDQuery struct {
	OrderID dot.ID
	ShopID  dot.ID

	Result []*Refund `json:"-"`
}

func (h QueryServiceHandler) HandleGetRefundsByOrderID(ctx context.Context, msg *GetRefundsByOrderIDQuery) (err error) {
	msg.Result, err = h.inner.GetRefundsByOrderID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CancelRefundCommand) command()  {}
func (q *ConfirmRefundCommand) command() {}
func (q *CreateRefundCommand) command()  {}
func (q *UpdateRefundCommand) command()  {}

func (q *GetRefundByIDQuery) query()       {}
func (q *GetRefundsQuery) query()          {}
func (q *GetRefundsByIDsQuery) query()     {}
func (q *GetRefundsByOrderIDQuery) query() {}

// implement conversion

func (q *CancelRefundCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelRefundArgs) {
	return ctx,
		&CancelRefundArgs{
			ShopID:               q.ShopID,
			ID:                   q.ID,
			UpdatedBy:            q.UpdatedBy,
			CancelReason:         q.CancelReason,
			AutoInventoryVoucher: q.AutoInventoryVoucher,
		}
}

func (q *CancelRefundCommand) SetCancelRefundArgs(args *CancelRefundArgs) {
	q.ShopID = args.ShopID
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
	q.CancelReason = args.CancelReason
	q.AutoInventoryVoucher = args.AutoInventoryVoucher
}

func (q *ConfirmRefundCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmRefundArgs) {
	return ctx,
		&ConfirmRefundArgs{
			ShopID:               q.ShopID,
			ID:                   q.ID,
			UpdatedBy:            q.UpdatedBy,
			AutoInventoryVoucher: q.AutoInventoryVoucher,
		}
}

func (q *ConfirmRefundCommand) SetConfirmRefundArgs(args *ConfirmRefundArgs) {
	q.ShopID = args.ShopID
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
	q.AutoInventoryVoucher = args.AutoInventoryVoucher
}

func (q *CreateRefundCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateRefundArgs) {
	return ctx,
		&CreateRefundArgs{
			Lines:           q.Lines,
			OrderID:         q.OrderID,
			AdjustmentLines: q.AdjustmentLines,
			TotalAdjustment: q.TotalAdjustment,
			TotalAmount:     q.TotalAmount,
			BasketValue:     q.BasketValue,
			ShopID:          q.ShopID,
			CreatedBy:       q.CreatedBy,
			Note:            q.Note,
		}
}

func (q *CreateRefundCommand) SetCreateRefundArgs(args *CreateRefundArgs) {
	q.Lines = args.Lines
	q.OrderID = args.OrderID
	q.AdjustmentLines = args.AdjustmentLines
	q.TotalAdjustment = args.TotalAdjustment
	q.TotalAmount = args.TotalAmount
	q.BasketValue = args.BasketValue
	q.ShopID = args.ShopID
	q.CreatedBy = args.CreatedBy
	q.Note = args.Note
}

func (q *UpdateRefundCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateRefundArgs) {
	return ctx,
		&UpdateRefundArgs{
			Lines:           q.Lines,
			ID:              q.ID,
			ShopID:          q.ShopID,
			AdjustmentLines: q.AdjustmentLines,
			TotalAdjustment: q.TotalAdjustment,
			TotalAmount:     q.TotalAmount,
			BasketValue:     q.BasketValue,
			UpdateBy:        q.UpdateBy,
			Note:            q.Note,
		}
}

func (q *UpdateRefundCommand) SetUpdateRefundArgs(args *UpdateRefundArgs) {
	q.Lines = args.Lines
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.AdjustmentLines = args.AdjustmentLines
	q.TotalAdjustment = args.TotalAdjustment
	q.TotalAmount = args.TotalAmount
	q.BasketValue = args.BasketValue
	q.UpdateBy = args.UpdateBy
	q.Note = args.Note
}

func (q *GetRefundByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetRefundByIDArgs) {
	return ctx,
		&GetRefundByIDArgs{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *GetRefundByIDQuery) SetGetRefundByIDArgs(args *GetRefundByIDArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *GetRefundsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetRefundsArgs) {
	return ctx,
		&GetRefundsArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *GetRefundsQuery) SetGetRefundsArgs(args *GetRefundsArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *GetRefundsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetRefundsByIDsArgs) {
	return ctx,
		&GetRefundsByIDsArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *GetRefundsByIDsQuery) SetGetRefundsByIDsArgs(args *GetRefundsByIDsArgs) {
	q.IDs = args.IDs
	q.ShopID = args.ShopID
}

func (q *GetRefundsByOrderIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetRefundsByOrderID) {
	return ctx,
		&GetRefundsByOrderID{
			OrderID: q.OrderID,
			ShopID:  q.ShopID,
		}
}

func (q *GetRefundsByOrderIDQuery) SetGetRefundsByOrderID(args *GetRefundsByOrderID) {
	q.OrderID = args.OrderID
	q.ShopID = args.ShopID
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
	b.AddHandler(h.HandleCancelRefund)
	b.AddHandler(h.HandleConfirmRefund)
	b.AddHandler(h.HandleCreateRefund)
	b.AddHandler(h.HandleUpdateRefund)
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
	b.AddHandler(h.HandleGetRefundByID)
	b.AddHandler(h.HandleGetRefunds)
	b.AddHandler(h.HandleGetRefundsByIDs)
	b.AddHandler(h.HandleGetRefundsByOrderID)
	return QueryBus{b}
}