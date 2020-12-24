// +build !generator

// Code generated by generator api. DO NOT EDIT.

package credit

import (
	context "context"
	time "time"

	meta "o.o/api/meta"
	credit_type "o.o/api/top/types/etc/credit_type"
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

type ConfirmCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result *CreditExtended `json:"-"`
}

func (h AggregateHandler) HandleConfirmCredit(ctx context.Context, msg *ConfirmCreditCommand) (err error) {
	msg.Result, err = h.inner.ConfirmCredit(msg.GetArgs(ctx))
	return err
}

type CreateCreditCommand struct {
	Amount   int
	ShopID   dot.ID
	Type     credit_type.CreditType
	PaidAt   time.Time
	Classify credit_type.CreditClassify

	Result *CreditExtended `json:"-"`
}

func (h AggregateHandler) HandleCreateCredit(ctx context.Context, msg *CreateCreditCommand) (err error) {
	msg.Result, err = h.inner.CreateCredit(msg.GetArgs(ctx))
	return err
}

type DeleteCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteCredit(ctx context.Context, msg *DeleteCreditCommand) (err error) {
	msg.Result, err = h.inner.DeleteCredit(msg.GetArgs(ctx))
	return err
}

type GetCreditQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *CreditExtended `json:"-"`
}

func (h QueryServiceHandler) HandleGetCredit(ctx context.Context, msg *GetCreditQuery) (err error) {
	msg.Result, err = h.inner.GetCredit(msg.GetArgs(ctx))
	return err
}

type GetShippingUserBalanceQuery struct {
	UserID dot.ID

	Result *GetShippingUserBalanceResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetShippingUserBalance(ctx context.Context, msg *GetShippingUserBalanceQuery) (err error) {
	msg.Result, err = h.inner.GetShippingUserBalance(msg.GetArgs(ctx))
	return err
}

type GetTelecomUserBalanceQuery struct {
	UserID dot.ID

	Result int `json:"-"`
}

func (h QueryServiceHandler) HandleGetTelecomUserBalance(ctx context.Context, msg *GetTelecomUserBalanceQuery) (err error) {
	msg.Result, err = h.inner.GetTelecomUserBalance(msg.GetArgs(ctx))
	return err
}

type ListCreditsQuery struct {
	ShopID dot.ID
	Paging *meta.Paging

	Result *ListCreditsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListCredits(ctx context.Context, msg *ListCreditsQuery) (err error) {
	msg.Result, err = h.inner.ListCredits(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *ConfirmCreditCommand) command() {}
func (q *CreateCreditCommand) command()  {}
func (q *DeleteCreditCommand) command()  {}

func (q *GetCreditQuery) query()              {}
func (q *GetShippingUserBalanceQuery) query() {}
func (q *GetTelecomUserBalanceQuery) query()  {}
func (q *ListCreditsQuery) query()            {}

// implement conversion

func (q *ConfirmCreditCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmCreditArgs) {
	return ctx,
		&ConfirmCreditArgs{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *ConfirmCreditCommand) SetConfirmCreditArgs(args *ConfirmCreditArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *CreateCreditCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateCreditArgs) {
	return ctx,
		&CreateCreditArgs{
			Amount:   q.Amount,
			ShopID:   q.ShopID,
			Type:     q.Type,
			PaidAt:   q.PaidAt,
			Classify: q.Classify,
		}
}

func (q *CreateCreditCommand) SetCreateCreditArgs(args *CreateCreditArgs) {
	q.Amount = args.Amount
	q.ShopID = args.ShopID
	q.Type = args.Type
	q.PaidAt = args.PaidAt
	q.Classify = args.Classify
}

func (q *DeleteCreditCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteCreditArgs) {
	return ctx,
		&DeleteCreditArgs{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *DeleteCreditCommand) SetDeleteCreditArgs(args *DeleteCreditArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *GetCreditQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetCreditArgs) {
	return ctx,
		&GetCreditArgs{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *GetCreditQuery) SetGetCreditArgs(args *GetCreditArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *GetShippingUserBalanceQuery) GetArgs(ctx context.Context) (_ context.Context, UserID dot.ID) {
	return ctx,
		q.UserID
}

func (q *GetTelecomUserBalanceQuery) GetArgs(ctx context.Context) (_ context.Context, UserID dot.ID) {
	return ctx,
		q.UserID
}

func (q *ListCreditsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListCreditsArgs) {
	return ctx,
		&ListCreditsArgs{
			ShopID: q.ShopID,
			Paging: q.Paging,
		}
}

func (q *ListCreditsQuery) SetListCreditsArgs(args *ListCreditsArgs) {
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
	b.AddHandler(h.HandleConfirmCredit)
	b.AddHandler(h.HandleCreateCredit)
	b.AddHandler(h.HandleDeleteCredit)
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
	b.AddHandler(h.HandleGetCredit)
	b.AddHandler(h.HandleGetShippingUserBalance)
	b.AddHandler(h.HandleGetTelecomUserBalance)
	b.AddHandler(h.HandleListCredits)
	return QueryBus{b}
}
