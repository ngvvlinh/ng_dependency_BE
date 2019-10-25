// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package ledgering

import (
	context "context"

	identity "etop.vn/api/main/identity"
	meta "etop.vn/api/meta"
	shopping "etop.vn/api/shopping"
	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
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

type CreateLedgerCommand struct {
	ShopID      int64
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        string
	CreatedBy   int64

	Result *ShopLedger `json:"-"`
}

func (h AggregateHandler) HandleCreateLedger(ctx context.Context, msg *CreateLedgerCommand) (err error) {
	msg.Result, err = h.inner.CreateLedger(msg.GetArgs(ctx))
	return err
}

type DeleteLedgerCommand struct {
	ID     int64
	ShopID int64

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteLedger(ctx context.Context, msg *DeleteLedgerCommand) (err error) {
	msg.Result, err = h.inner.DeleteLedger(msg.GetArgs(ctx))
	return err
}

type UpdateLedgerCommand struct {
	ID          int64
	ShopID      int64
	Name        dot.NullString
	BankAccount *identity.BankAccount
	Note        dot.NullString

	Result *ShopLedger `json:"-"`
}

func (h AggregateHandler) HandleUpdateLedger(ctx context.Context, msg *UpdateLedgerCommand) (err error) {
	msg.Result, err = h.inner.UpdateLedger(msg.GetArgs(ctx))
	return err
}

type GetLedgerByIDQuery struct {
	ID     int64
	ShopID int64

	Result *ShopLedger `json:"-"`
}

func (h QueryServiceHandler) HandleGetLedgerByID(ctx context.Context, msg *GetLedgerByIDQuery) (err error) {
	msg.Result, err = h.inner.GetLedgerByID(msg.GetArgs(ctx))
	return err
}

type ListLedgersQuery struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters

	Result *ShopLedgersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListLedgers(ctx context.Context, msg *ListLedgersQuery) (err error) {
	msg.Result, err = h.inner.ListLedgers(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateLedgerCommand) command() {}
func (q *DeleteLedgerCommand) command() {}
func (q *UpdateLedgerCommand) command() {}
func (q *GetLedgerByIDQuery) query()    {}
func (q *ListLedgersQuery) query()      {}

// implement conversion

func (q *CreateLedgerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateLedgerArgs) {
	return ctx,
		&CreateLedgerArgs{
			ShopID:      q.ShopID,
			Name:        q.Name,
			BankAccount: q.BankAccount,
			Note:        q.Note,
			Type:        q.Type,
			CreatedBy:   q.CreatedBy,
		}
}

func (q *CreateLedgerCommand) SetCreateLedgerArgs(args *CreateLedgerArgs) {
	q.ShopID = args.ShopID
	q.Name = args.Name
	q.BankAccount = args.BankAccount
	q.Note = args.Note
	q.Type = args.Type
	q.CreatedBy = args.CreatedBy
}

func (q *DeleteLedgerCommand) GetArgs(ctx context.Context) (_ context.Context, ID int64, ShopID int64) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *UpdateLedgerCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateLedgerArgs) {
	return ctx,
		&UpdateLedgerArgs{
			ID:          q.ID,
			ShopID:      q.ShopID,
			Name:        q.Name,
			BankAccount: q.BankAccount,
			Note:        q.Note,
		}
}

func (q *UpdateLedgerCommand) SetUpdateLedgerArgs(args *UpdateLedgerArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.Name = args.Name
	q.BankAccount = args.BankAccount
	q.Note = args.Note
}

func (q *GetLedgerByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *GetLedgerByIDQuery) SetIDQueryShopArg(args *shopping.IDQueryShopArg) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *ListLedgersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListLedgersQuery) SetListQueryShopArgs(args *shopping.ListQueryShopArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
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
	b.AddHandler(h.HandleCreateLedger)
	b.AddHandler(h.HandleDeleteLedger)
	b.AddHandler(h.HandleUpdateLedger)
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
	b.AddHandler(h.HandleGetLedgerByID)
	b.AddHandler(h.HandleListLedgers)
	return QueryBus{b}
}