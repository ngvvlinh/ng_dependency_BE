// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package transaction

import (
	context "context"

	etop "etop.vn/api/main/etop"
	meta "etop.vn/api/meta"
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

type CancelTransactionCommand struct {
	TrxnID    int64
	AccountID int64

	Result *Transaction `json:"-"`
}

func (h AggregateHandler) HandleCancelTransaction(ctx context.Context, msg *CancelTransactionCommand) (err error) {
	msg.Result, err = h.inner.CancelTransaction(msg.GetArgs(ctx))
	return err
}

type ConfirmTransactionCommand struct {
	TrxnID    int64
	AccountID int64

	Result *Transaction `json:"-"`
}

func (h AggregateHandler) HandleConfirmTransaction(ctx context.Context, msg *ConfirmTransactionCommand) (err error) {
	msg.Result, err = h.inner.ConfirmTransaction(msg.GetArgs(ctx))
	return err
}

type CreateTransactionCommand struct {
	ID        int64
	Amount    int
	AccountID int64
	Status    etop.Status3
	Type      TransactionType
	Note      string
	Metadata  *TransactionMetadata

	Result *Transaction `json:"-"`
}

func (h AggregateHandler) HandleCreateTransaction(ctx context.Context, msg *CreateTransactionCommand) (err error) {
	msg.Result, err = h.inner.CreateTransaction(msg.GetArgs(ctx))
	return err
}

type GetBalanceQuery struct {
	AccountID       int64
	TransactionType TransactionType

	Result int `json:"-"`
}

func (h QueryServiceHandler) HandleGetBalance(ctx context.Context, msg *GetBalanceQuery) (err error) {
	msg.Result, err = h.inner.GetBalance(msg.GetArgs(ctx))
	return err
}

type GetTransactionByIDQuery struct {
	TrxnID    int64
	AccountID int64

	Result *Transaction `json:"-"`
}

func (h QueryServiceHandler) HandleGetTransactionByID(ctx context.Context, msg *GetTransactionByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTransactionByID(msg.GetArgs(ctx))
	return err
}

type ListTransactionsQuery struct {
	AccountID int64
	Paging    meta.Paging

	Result *TransactionResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListTransactions(ctx context.Context, msg *ListTransactionsQuery) (err error) {
	msg.Result, err = h.inner.ListTransactions(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CancelTransactionCommand) command()  {}
func (q *ConfirmTransactionCommand) command() {}
func (q *CreateTransactionCommand) command()  {}
func (q *GetBalanceQuery) query()             {}
func (q *GetTransactionByIDQuery) query()     {}
func (q *ListTransactionsQuery) query()       {}

// implement conversion

func (q *CancelTransactionCommand) GetArgs(ctx context.Context) (_ context.Context, trxnID int64, accountID int64) {
	return ctx,
		q.TrxnID,
		q.AccountID
}

func (q *ConfirmTransactionCommand) GetArgs(ctx context.Context) (_ context.Context, trxnID int64, accountID int64) {
	return ctx,
		q.TrxnID,
		q.AccountID
}

func (q *CreateTransactionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateTransactionArgs) {
	return ctx,
		&CreateTransactionArgs{
			ID:        q.ID,
			Amount:    q.Amount,
			AccountID: q.AccountID,
			Status:    q.Status,
			Type:      q.Type,
			Note:      q.Note,
			Metadata:  q.Metadata,
		}
}

func (q *GetBalanceQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetBalanceArgs) {
	return ctx,
		&GetBalanceArgs{
			AccountID:       q.AccountID,
			TransactionType: q.TransactionType,
		}
}

func (q *GetTransactionByIDQuery) GetArgs(ctx context.Context) (_ context.Context, trxnID int64, accountID int64) {
	return ctx,
		q.TrxnID,
		q.AccountID
}

func (q *ListTransactionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTransactionsArgs) {
	return ctx,
		&GetTransactionsArgs{
			AccountID: q.AccountID,
			Paging:    q.Paging,
		}
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
	b.AddHandler(h.HandleCancelTransaction)
	b.AddHandler(h.HandleConfirmTransaction)
	b.AddHandler(h.HandleCreateTransaction)
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
	b.AddHandler(h.HandleGetBalance)
	b.AddHandler(h.HandleGetTransactionByID)
	b.AddHandler(h.HandleListTransactions)
	return QueryBus{b}
}
