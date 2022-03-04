// +build !generator

// Code generated by generator api. DO NOT EDIT.

package bankstatement

import (
	context "context"
	time "time"

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

type CreateBankStatementCommand struct {
	ID                    dot.ID
	Amount                int
	Description           string
	AccountID             dot.ID
	TransferedAt          time.Time
	ExternalTransactionID string
	SenderName            string
	SenderBankAccount     string
	OtherInfo             map[string]string
	CreatedAt             time.Time
	UpdatedAt             time.Time

	Result *BankStatement `json:"-"`
}

func (h AggregateHandler) HandleCreateBankStatement(ctx context.Context, msg *CreateBankStatementCommand) (err error) {
	msg.Result, err = h.inner.CreateBankStatement(msg.GetArgs(ctx))
	return err
}

type GetBankStatementQuery struct {
	ID                    dot.ID
	ExternalTransactionID string

	Result *BankStatement `json:"-"`
}

func (h QueryServiceHandler) HandleGetBankStatement(ctx context.Context, msg *GetBankStatementQuery) (err error) {
	msg.Result, err = h.inner.GetBankStatement(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateBankStatementCommand) command() {}

func (q *GetBankStatementQuery) query() {}

// implement conversion

func (q *CreateBankStatementCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateBankStatementArgs) {
	return ctx,
		&CreateBankStatementArgs{
			ID:                    q.ID,
			Amount:                q.Amount,
			Description:           q.Description,
			AccountID:             q.AccountID,
			TransferedAt:          q.TransferedAt,
			ExternalTransactionID: q.ExternalTransactionID,
			SenderName:            q.SenderName,
			SenderBankAccount:     q.SenderBankAccount,
			OtherInfo:             q.OtherInfo,
			CreatedAt:             q.CreatedAt,
			UpdatedAt:             q.UpdatedAt,
		}
}

func (q *CreateBankStatementCommand) SetCreateBankStatementArgs(args *CreateBankStatementArgs) {
	q.ID = args.ID
	q.Amount = args.Amount
	q.Description = args.Description
	q.AccountID = args.AccountID
	q.TransferedAt = args.TransferedAt
	q.ExternalTransactionID = args.ExternalTransactionID
	q.SenderName = args.SenderName
	q.SenderBankAccount = args.SenderBankAccount
	q.OtherInfo = args.OtherInfo
	q.CreatedAt = args.CreatedAt
	q.UpdatedAt = args.UpdatedAt
}

func (q *GetBankStatementQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetBankStatementArgs) {
	return ctx,
		&GetBankStatementArgs{
			ID:                    q.ID,
			ExternalTransactionID: q.ExternalTransactionID,
		}
}

func (q *GetBankStatementQuery) SetGetBankStatementArgs(args *GetBankStatementArgs) {
	q.ID = args.ID
	q.ExternalTransactionID = args.ExternalTransactionID
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
	b.AddHandler(h.HandleCreateBankStatement)
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
	b.AddHandler(h.HandleGetBankStatement)
	return QueryBus{b}
}
