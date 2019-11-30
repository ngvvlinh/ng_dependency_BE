// +build !generator

// Code generated by generator api. DO NOT EDIT.

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
	ShopID      dot.ID
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        LedgerType
	CreatedBy   dot.ID

	Result *ShopLedger `json:"-"`
}

func (h AggregateHandler) HandleCreateLedger(ctx context.Context, msg *CreateLedgerCommand) (err error) {
	msg.Result, err = h.inner.CreateLedger(msg.GetArgs(ctx))
	return err
}

type DeleteLedgerCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteLedger(ctx context.Context, msg *DeleteLedgerCommand) (err error) {
	msg.Result, err = h.inner.DeleteLedger(msg.GetArgs(ctx))
	return err
}

type UpdateLedgerCommand struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        dot.NullString
	BankAccount *identity.BankAccount
	Note        dot.NullString

	Result *ShopLedger `json:"-"`
}

func (h AggregateHandler) HandleUpdateLedger(ctx context.Context, msg *UpdateLedgerCommand) (err error) {
	msg.Result, err = h.inner.UpdateLedger(msg.GetArgs(ctx))
	return err
}

type GetLedgerByAccountNumberQuery struct {
	AccountNumber string
	ShopID        dot.ID

	Result *ShopLedger `json:"-"`
}

func (h QueryServiceHandler) HandleGetLedgerByAccountNumber(ctx context.Context, msg *GetLedgerByAccountNumberQuery) (err error) {
	msg.Result, err = h.inner.GetLedgerByAccountNumber(msg.GetArgs(ctx))
	return err
}

type GetLedgerByIDQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *ShopLedger `json:"-"`
}

func (h QueryServiceHandler) HandleGetLedgerByID(ctx context.Context, msg *GetLedgerByIDQuery) (err error) {
	msg.Result, err = h.inner.GetLedgerByID(msg.GetArgs(ctx))
	return err
}

type ListLedgersQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *ShopLedgersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListLedgers(ctx context.Context, msg *ListLedgersQuery) (err error) {
	msg.Result, err = h.inner.ListLedgers(msg.GetArgs(ctx))
	return err
}

type ListLedgersByIDsQuery struct {
	ShopID dot.ID
	IDs    []dot.ID

	Result *ShopLedgersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListLedgersByIDs(ctx context.Context, msg *ListLedgersByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListLedgersByIDs(msg.GetArgs(ctx))
	return err
}

type ListLedgersByTypeQuery struct {
	LedgerType LedgerType
	ShopID     dot.ID

	Result *ShopLedgersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListLedgersByType(ctx context.Context, msg *ListLedgersByTypeQuery) (err error) {
	msg.Result, err = h.inner.ListLedgersByType(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateLedgerCommand) command()         {}
func (q *DeleteLedgerCommand) command()         {}
func (q *UpdateLedgerCommand) command()         {}
func (q *GetLedgerByAccountNumberQuery) query() {}
func (q *GetLedgerByIDQuery) query()            {}
func (q *ListLedgersQuery) query()              {}
func (q *ListLedgersByIDsQuery) query()         {}
func (q *ListLedgersByTypeQuery) query()        {}

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

func (q *DeleteLedgerCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, ShopID dot.ID) {
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

func (q *GetLedgerByAccountNumberQuery) GetArgs(ctx context.Context) (_ context.Context, accountNumber string, shopID dot.ID) {
	return ctx,
		q.AccountNumber,
		q.ShopID
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

func (q *ListLedgersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, shopID dot.ID, IDs []dot.ID) {
	return ctx,
		q.ShopID,
		q.IDs
}

func (q *ListLedgersByTypeQuery) GetArgs(ctx context.Context) (_ context.Context, ledgerType LedgerType, shopID dot.ID) {
	return ctx,
		q.LedgerType,
		q.ShopID
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
	b.AddHandler(h.HandleGetLedgerByAccountNumber)
	b.AddHandler(h.HandleGetLedgerByID)
	b.AddHandler(h.HandleListLedgers)
	b.AddHandler(h.HandleListLedgersByIDs)
	b.AddHandler(h.HandleListLedgersByType)
	return QueryBus{b}
}