// +build !generator

// Code generated by generator api. DO NOT EDIT.

package receipting

import (
	context "context"
	time "time"

	etop "etop.vn/api/main/etop"
	meta "etop.vn/api/meta"
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

type CancelReceiptCommand struct {
	ID     dot.ID
	ShopID dot.ID
	Reason string

	Result int `json:"-"`
}

func (h AggregateHandler) HandleCancelReceipt(ctx context.Context, msg *CancelReceiptCommand) (err error) {
	msg.Result, err = h.inner.CancelReceipt(msg.GetArgs(ctx))
	return err
}

type ConfirmReceiptCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleConfirmReceipt(ctx context.Context, msg *ConfirmReceiptCommand) (err error) {
	msg.Result, err = h.inner.ConfirmReceipt(msg.GetArgs(ctx))
	return err
}

type CreateReceiptCommand struct {
	ShopID      dot.ID
	TraderID    dot.ID
	Title       string
	Type        ReceiptType
	Status      int
	Description string
	Amount      int
	LedgerID    dot.ID
	RefIDs      []dot.ID
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	CreatedBy   dot.ID
	CreatedType ReceiptCreatedType
	ConfirmedAt time.Time

	Result *Receipt `json:"-"`
}

func (h AggregateHandler) HandleCreateReceipt(ctx context.Context, msg *CreateReceiptCommand) (err error) {
	msg.Result, err = h.inner.CreateReceipt(msg.GetArgs(ctx))
	return err
}

type DeleteReceiptCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteReceipt(ctx context.Context, msg *DeleteReceiptCommand) (err error) {
	msg.Result, err = h.inner.DeleteReceipt(msg.GetArgs(ctx))
	return err
}

type UpdateReceiptCommand struct {
	ID          dot.ID
	ShopID      dot.ID
	TraderID    dot.NullID
	Title       dot.NullString
	Description dot.NullString
	Amount      dot.NullInt
	LedgerID    dot.NullID
	RefIDs      []dot.ID
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time

	Result *Receipt `json:"-"`
}

func (h AggregateHandler) HandleUpdateReceipt(ctx context.Context, msg *UpdateReceiptCommand) (err error) {
	msg.Result, err = h.inner.UpdateReceipt(msg.GetArgs(ctx))
	return err
}

type GetReceiptByCodeQuery struct {
	Code   string
	ShopID dot.ID

	Result *Receipt `json:"-"`
}

func (h QueryServiceHandler) HandleGetReceiptByCode(ctx context.Context, msg *GetReceiptByCodeQuery) (err error) {
	msg.Result, err = h.inner.GetReceiptByCode(msg.GetArgs(ctx))
	return err
}

type GetReceiptByIDQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *Receipt `json:"-"`
}

func (h QueryServiceHandler) HandleGetReceiptByID(ctx context.Context, msg *GetReceiptByIDQuery) (err error) {
	msg.Result, err = h.inner.GetReceiptByID(msg.GetArgs(ctx))
	return err
}

type ListReceiptsQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceipts(ctx context.Context, msg *ListReceiptsQuery) (err error) {
	msg.Result, err = h.inner.ListReceipts(msg.GetArgs(ctx))
	return err
}

type ListReceiptsByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceiptsByIDs(ctx context.Context, msg *ListReceiptsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListReceiptsByIDs(msg.GetArgs(ctx))
	return err
}

type ListReceiptsByLedgerIDsQuery struct {
	ShopID    dot.ID
	LedgerIDs []dot.ID
	Paging    meta.Paging
	Filters   meta.Filters

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceiptsByLedgerIDs(ctx context.Context, msg *ListReceiptsByLedgerIDsQuery) (err error) {
	msg.Result, err = h.inner.ListReceiptsByLedgerIDs(msg.GetArgs(ctx))
	return err
}

type ListReceiptsByRefsAndStatusQuery struct {
	ShopID  dot.ID
	RefIDs  []dot.ID
	RefType ReceiptRefType
	Status  int

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceiptsByRefsAndStatus(ctx context.Context, msg *ListReceiptsByRefsAndStatusQuery) (err error) {
	msg.Result, err = h.inner.ListReceiptsByRefsAndStatus(msg.GetArgs(ctx))
	return err
}

type ListReceiptsByTraderIDsAndStatusesQuery struct {
	ShopID    dot.ID
	TraderIDs []dot.ID
	Statuses  []etop.Status3

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceiptsByTraderIDsAndStatuses(ctx context.Context, msg *ListReceiptsByTraderIDsAndStatusesQuery) (err error) {
	msg.Result, err = h.inner.ListReceiptsByTraderIDsAndStatuses(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CancelReceiptCommand) command()  {}
func (q *ConfirmReceiptCommand) command() {}
func (q *CreateReceiptCommand) command()  {}
func (q *DeleteReceiptCommand) command()  {}
func (q *UpdateReceiptCommand) command()  {}

func (q *GetReceiptByCodeQuery) query()                   {}
func (q *GetReceiptByIDQuery) query()                     {}
func (q *ListReceiptsQuery) query()                       {}
func (q *ListReceiptsByIDsQuery) query()                  {}
func (q *ListReceiptsByLedgerIDsQuery) query()            {}
func (q *ListReceiptsByRefsAndStatusQuery) query()        {}
func (q *ListReceiptsByTraderIDsAndStatusesQuery) query() {}

// implement conversion

func (q *CancelReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelReceiptArgs) {
	return ctx,
		&CancelReceiptArgs{
			ID:     q.ID,
			ShopID: q.ShopID,
			Reason: q.Reason,
		}
}

func (q *CancelReceiptCommand) SetCancelReceiptArgs(args *CancelReceiptArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.Reason = args.Reason
}

func (q *ConfirmReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmReceiptArgs) {
	return ctx,
		&ConfirmReceiptArgs{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *ConfirmReceiptCommand) SetConfirmReceiptArgs(args *ConfirmReceiptArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *CreateReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateReceiptArgs) {
	return ctx,
		&CreateReceiptArgs{
			ShopID:      q.ShopID,
			TraderID:    q.TraderID,
			Title:       q.Title,
			Type:        q.Type,
			Status:      q.Status,
			Description: q.Description,
			Amount:      q.Amount,
			LedgerID:    q.LedgerID,
			RefIDs:      q.RefIDs,
			RefType:     q.RefType,
			Lines:       q.Lines,
			Trader:      q.Trader,
			PaidAt:      q.PaidAt,
			CreatedBy:   q.CreatedBy,
			CreatedType: q.CreatedType,
			ConfirmedAt: q.ConfirmedAt,
		}
}

func (q *CreateReceiptCommand) SetCreateReceiptArgs(args *CreateReceiptArgs) {
	q.ShopID = args.ShopID
	q.TraderID = args.TraderID
	q.Title = args.Title
	q.Type = args.Type
	q.Status = args.Status
	q.Description = args.Description
	q.Amount = args.Amount
	q.LedgerID = args.LedgerID
	q.RefIDs = args.RefIDs
	q.RefType = args.RefType
	q.Lines = args.Lines
	q.Trader = args.Trader
	q.PaidAt = args.PaidAt
	q.CreatedBy = args.CreatedBy
	q.CreatedType = args.CreatedType
	q.ConfirmedAt = args.ConfirmedAt
}

func (q *DeleteReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, shopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *UpdateReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateReceiptArgs) {
	return ctx,
		&UpdateReceiptArgs{
			ID:          q.ID,
			ShopID:      q.ShopID,
			TraderID:    q.TraderID,
			Title:       q.Title,
			Description: q.Description,
			Amount:      q.Amount,
			LedgerID:    q.LedgerID,
			RefIDs:      q.RefIDs,
			RefType:     q.RefType,
			Lines:       q.Lines,
			Trader:      q.Trader,
			PaidAt:      q.PaidAt,
		}
}

func (q *UpdateReceiptCommand) SetUpdateReceiptArgs(args *UpdateReceiptArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.TraderID = args.TraderID
	q.Title = args.Title
	q.Description = args.Description
	q.Amount = args.Amount
	q.LedgerID = args.LedgerID
	q.RefIDs = args.RefIDs
	q.RefType = args.RefType
	q.Lines = args.Lines
	q.Trader = args.Trader
	q.PaidAt = args.PaidAt
}

func (q *GetReceiptByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, code string, shopID dot.ID) {
	return ctx,
		q.Code,
		q.ShopID
}

func (q *GetReceiptByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetReceiptByIDArg) {
	return ctx,
		&GetReceiptByIDArg{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *GetReceiptByIDQuery) SetGetReceiptByIDArg(args *GetReceiptByIDArg) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *ListReceiptsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListReceiptsArgs) {
	return ctx,
		&ListReceiptsArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListReceiptsQuery) SetListReceiptsArgs(args *ListReceiptsArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListReceiptsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetReceiptbyIDsArgs) {
	return ctx,
		&GetReceiptbyIDsArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *ListReceiptsByIDsQuery) SetGetReceiptbyIDsArgs(args *GetReceiptbyIDsArgs) {
	q.IDs = args.IDs
	q.ShopID = args.ShopID
}

func (q *ListReceiptsByLedgerIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListReceiptsByLedgerIDsArgs) {
	return ctx,
		&ListReceiptsByLedgerIDsArgs{
			ShopID:    q.ShopID,
			LedgerIDs: q.LedgerIDs,
			Paging:    q.Paging,
			Filters:   q.Filters,
		}
}

func (q *ListReceiptsByLedgerIDsQuery) SetListReceiptsByLedgerIDsArgs(args *ListReceiptsByLedgerIDsArgs) {
	q.ShopID = args.ShopID
	q.LedgerIDs = args.LedgerIDs
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListReceiptsByRefsAndStatusQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListReceiptsByRefsAndStatusArgs) {
	return ctx,
		&ListReceiptsByRefsAndStatusArgs{
			ShopID:  q.ShopID,
			RefIDs:  q.RefIDs,
			RefType: q.RefType,
			Status:  q.Status,
		}
}

func (q *ListReceiptsByRefsAndStatusQuery) SetListReceiptsByRefsAndStatusArgs(args *ListReceiptsByRefsAndStatusArgs) {
	q.ShopID = args.ShopID
	q.RefIDs = args.RefIDs
	q.RefType = args.RefType
	q.Status = args.Status
}

func (q *ListReceiptsByTraderIDsAndStatusesQuery) GetArgs(ctx context.Context) (_ context.Context, shopID dot.ID, traderIDs []dot.ID, statuses []etop.Status3) {
	return ctx,
		q.ShopID,
		q.TraderIDs,
		q.Statuses
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
	b.AddHandler(h.HandleCancelReceipt)
	b.AddHandler(h.HandleConfirmReceipt)
	b.AddHandler(h.HandleCreateReceipt)
	b.AddHandler(h.HandleDeleteReceipt)
	b.AddHandler(h.HandleUpdateReceipt)
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
	b.AddHandler(h.HandleGetReceiptByCode)
	b.AddHandler(h.HandleGetReceiptByID)
	b.AddHandler(h.HandleListReceipts)
	b.AddHandler(h.HandleListReceiptsByIDs)
	b.AddHandler(h.HandleListReceiptsByLedgerIDs)
	b.AddHandler(h.HandleListReceiptsByRefsAndStatus)
	b.AddHandler(h.HandleListReceiptsByTraderIDsAndStatuses)
	return QueryBus{b}
}
