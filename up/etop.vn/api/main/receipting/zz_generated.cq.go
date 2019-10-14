// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package receipting

import (
	context "context"

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

type CreateReceiptCommand struct {
	ShopID      int64
	TraderID    int64
	Code        string
	Title       string
	Type        string
	Description string
	Amount      int32
	OrderIDs    []int64
	Lines       []*ReceiptLine
	CreatedBy   int64

	Result *Receipt `json:"-"`
}

func (h AggregateHandler) HandleCreateReceipt(ctx context.Context, msg *CreateReceiptCommand) (err error) {
	msg.Result, err = h.inner.CreateReceipt(msg.GetArgs(ctx))
	return err
}

type DeleteReceiptCommand struct {
	ID     int64
	ShopID int64

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteReceipt(ctx context.Context, msg *DeleteReceiptCommand) (err error) {
	msg.Result, err = h.inner.DeleteReceipt(msg.GetArgs(ctx))
	return err
}

type UpdateReceiptCommand struct {
	ID          int64
	ShopID      int64
	TraderID    dot.NullInt64
	Title       dot.NullString
	Code        dot.NullString
	Description dot.NullString
	Amount      dot.NullInt32
	OrderIDs    []int64
	Lines       []*ReceiptLine

	Result *Receipt `json:"-"`
}

func (h AggregateHandler) HandleUpdateReceipt(ctx context.Context, msg *UpdateReceiptCommand) (err error) {
	msg.Result, err = h.inner.UpdateReceipt(msg.GetArgs(ctx))
	return err
}

type GetReceiptByCodeQuery struct {
	Code   string
	ShopID int64

	Result *Receipt `json:"-"`
}

func (h QueryServiceHandler) HandleGetReceiptByCode(ctx context.Context, msg *GetReceiptByCodeQuery) (err error) {
	msg.Result, err = h.inner.GetReceiptByCode(msg.GetArgs(ctx))
	return err
}

type GetReceiptByIDQuery struct {
	ID     int64
	ShopID int64

	Result *Receipt `json:"-"`
}

func (h QueryServiceHandler) HandleGetReceiptByID(ctx context.Context, msg *GetReceiptByIDQuery) (err error) {
	msg.Result, err = h.inner.GetReceiptByID(msg.GetArgs(ctx))
	return err
}

type ListReceiptsQuery struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceipts(ctx context.Context, msg *ListReceiptsQuery) (err error) {
	msg.Result, err = h.inner.ListReceipts(msg.GetArgs(ctx))
	return err
}

type ListReceiptsByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceiptsByIDs(ctx context.Context, msg *ListReceiptsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListReceiptsByIDs(msg.GetArgs(ctx))
	return err
}

type ListReceiptsByOrderIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *ReceiptsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListReceiptsByOrderIDs(ctx context.Context, msg *ListReceiptsByOrderIDsQuery) (err error) {
	msg.Result, err = h.inner.ListReceiptsByOrderIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateReceiptCommand) command()      {}
func (q *DeleteReceiptCommand) command()      {}
func (q *UpdateReceiptCommand) command()      {}
func (q *GetReceiptByCodeQuery) query()       {}
func (q *GetReceiptByIDQuery) query()         {}
func (q *ListReceiptsQuery) query()           {}
func (q *ListReceiptsByIDsQuery) query()      {}
func (q *ListReceiptsByOrderIDsQuery) query() {}

// implement conversion

func (q *CreateReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateReceiptArgs) {
	return ctx,
		&CreateReceiptArgs{
			ShopID:      q.ShopID,
			TraderID:    q.TraderID,
			Code:        q.Code,
			Title:       q.Title,
			Type:        q.Type,
			Description: q.Description,
			Amount:      q.Amount,
			OrderIDs:    q.OrderIDs,
			Lines:       q.Lines,
			CreatedBy:   q.CreatedBy,
		}
}

func (q *DeleteReceiptCommand) GetArgs(ctx context.Context) (_ context.Context, ID int64, shopID int64) {
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
			Code:        q.Code,
			Description: q.Description,
			Amount:      q.Amount,
			OrderIDs:    q.OrderIDs,
			Lines:       q.Lines,
		}
}

func (q *GetReceiptByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, code string, shopID int64) {
	return ctx,
		q.Code,
		q.ShopID
}

func (q *GetReceiptByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *ListReceiptsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListReceiptsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *ListReceiptsByOrderIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
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
	b.AddHandler(h.HandleListReceiptsByOrderIDs)
	return QueryBus{b}
}
