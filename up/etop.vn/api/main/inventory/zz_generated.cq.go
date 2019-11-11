// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package inventory

import (
	context "context"

	meta "etop.vn/api/meta"
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

type AdjustInventoryQuantityCommand struct {
	Overstock bool
	ShopID    int64
	Lines     []*InventoryVariant
	Title     string
	UserID    int64
	Note      string

	Result *AdjustInventoryQuantityRespone `json:"-"`
}

func (h AggregateHandler) HandleAdjustInventoryQuantity(ctx context.Context, msg *AdjustInventoryQuantityCommand) (err error) {
	msg.Result, err = h.inner.AdjustInventoryQuantity(msg.GetArgs(ctx))
	return err
}

type CancelInventoryVoucherCommand struct {
	ShopID    int64
	ID        int64
	UpdatedBy int64
	Reason    string

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleCancelInventoryVoucher(ctx context.Context, msg *CancelInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.CancelInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type ConfirmInventoryVoucherCommand struct {
	ShopID    int64
	ID        int64
	UpdatedBy int64

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleConfirmInventoryVoucher(ctx context.Context, msg *ConfirmInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.ConfirmInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type CreateInventoryVariantCommand struct {
	ShopID    int64
	VariantID int64

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleCreateInventoryVariant(ctx context.Context, msg *CreateInventoryVariantCommand) (err error) {
	return h.inner.CreateInventoryVariant(msg.GetArgs(ctx))
}

type CreateInventoryVoucherCommand struct {
	Overstock   bool
	ShopID      int64
	CreatedBy   int64
	Title       string
	RefID       int64
	RefType     InventoryRefType
	RefName     InventoryVoucherRefName
	TraderID    int64
	TotalAmount int32
	Type        InventoryVoucherType
	Note        string
	Lines       []*InventoryVoucherItem

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleCreateInventoryVoucher(ctx context.Context, msg *CreateInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.CreateInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type UpdateInventoryVoucherCommand struct {
	ID          int64
	ShopID      int64
	Title       dot.NullString
	UpdatedBy   int64
	TraderID    dot.NullInt64
	TotalAmount int32
	Note        dot.NullString
	Lines       []*InventoryVoucherItem

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleUpdateInventoryVoucher(ctx context.Context, msg *UpdateInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.UpdateInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type GetInventoriesQuery struct {
	ShopID int64
	Paging *meta.Paging

	Result *GetInventoriesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventories(ctx context.Context, msg *GetInventoriesQuery) (err error) {
	msg.Result, err = h.inner.GetInventories(msg.GetArgs(ctx))
	return err
}

type GetInventoriesByVariantIDsQuery struct {
	ShopID     int64
	Paging     *meta.Paging
	VariantIDs []int64

	Result *GetInventoriesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoriesByVariantIDs(ctx context.Context, msg *GetInventoriesByVariantIDsQuery) (err error) {
	msg.Result, err = h.inner.GetInventoriesByVariantIDs(msg.GetArgs(ctx))
	return err
}

type GetInventoryQuery struct {
	ShopID    int64
	VariantID int64

	Result *InventoryVariant `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventory(ctx context.Context, msg *GetInventoryQuery) (err error) {
	msg.Result, err = h.inner.GetInventory(msg.GetArgs(ctx))
	return err
}

type GetInventoryVoucherQuery struct {
	ShopID int64
	ID     int64

	Result *InventoryVoucher `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVoucher(ctx context.Context, msg *GetInventoryVoucherQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type GetInventoryVouchersQuery struct {
	ShopID int64
	Paging *meta.Paging

	Result *GetInventoryVouchersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVouchers(ctx context.Context, msg *GetInventoryVouchersQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVouchers(msg.GetArgs(ctx))
	return err
}

type GetInventoryVouchersByIDsQuery struct {
	ShopID int64
	Paging *meta.Paging
	IDs    []int64

	Result *GetInventoryVouchersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVouchersByIDs(ctx context.Context, msg *GetInventoryVouchersByIDsQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVouchersByIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AdjustInventoryQuantityCommand) command() {}
func (q *CancelInventoryVoucherCommand) command()  {}
func (q *ConfirmInventoryVoucherCommand) command() {}
func (q *CreateInventoryVariantCommand) command()  {}
func (q *CreateInventoryVoucherCommand) command()  {}
func (q *UpdateInventoryVoucherCommand) command()  {}
func (q *GetInventoriesQuery) query()              {}
func (q *GetInventoriesByVariantIDsQuery) query()  {}
func (q *GetInventoryQuery) query()                {}
func (q *GetInventoryVoucherQuery) query()         {}
func (q *GetInventoryVouchersQuery) query()        {}
func (q *GetInventoryVouchersByIDsQuery) query()   {}

// implement conversion

func (q *AdjustInventoryQuantityCommand) GetArgs(ctx context.Context) (_ context.Context, Overstock bool, _ *AdjustInventoryQuantityArgs) {
	return ctx,
		q.Overstock,
		&AdjustInventoryQuantityArgs{
			ShopID: q.ShopID,
			Lines:  q.Lines,
			Title:  q.Title,
			UserID: q.UserID,
			Note:   q.Note,
		}
}

func (q *AdjustInventoryQuantityCommand) SetAdjustInventoryQuantityArgs(args *AdjustInventoryQuantityArgs) {
	q.ShopID = args.ShopID
	q.Lines = args.Lines
	q.Title = args.Title
	q.UserID = args.UserID
	q.Note = args.Note
}

func (q *CancelInventoryVoucherCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelInventoryVoucherArgs) {
	return ctx,
		&CancelInventoryVoucherArgs{
			ShopID:    q.ShopID,
			ID:        q.ID,
			UpdatedBy: q.UpdatedBy,
			Reason:    q.Reason,
		}
}

func (q *CancelInventoryVoucherCommand) SetCancelInventoryVoucherArgs(args *CancelInventoryVoucherArgs) {
	q.ShopID = args.ShopID
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
	q.Reason = args.Reason
}

func (q *ConfirmInventoryVoucherCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmInventoryVoucherArgs) {
	return ctx,
		&ConfirmInventoryVoucherArgs{
			ShopID:    q.ShopID,
			ID:        q.ID,
			UpdatedBy: q.UpdatedBy,
		}
}

func (q *ConfirmInventoryVoucherCommand) SetConfirmInventoryVoucherArgs(args *ConfirmInventoryVoucherArgs) {
	q.ShopID = args.ShopID
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
}

func (q *CreateInventoryVariantCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateInventoryVariantArgs) {
	return ctx,
		&CreateInventoryVariantArgs{
			ShopID:    q.ShopID,
			VariantID: q.VariantID,
		}
}

func (q *CreateInventoryVariantCommand) SetCreateInventoryVariantArgs(args *CreateInventoryVariantArgs) {
	q.ShopID = args.ShopID
	q.VariantID = args.VariantID
}

func (q *CreateInventoryVoucherCommand) GetArgs(ctx context.Context) (_ context.Context, Overstock bool, _ *CreateInventoryVoucherArgs) {
	return ctx,
		q.Overstock,
		&CreateInventoryVoucherArgs{
			ShopID:      q.ShopID,
			CreatedBy:   q.CreatedBy,
			Title:       q.Title,
			RefID:       q.RefID,
			RefType:     q.RefType,
			RefName:     q.RefName,
			TraderID:    q.TraderID,
			TotalAmount: q.TotalAmount,
			Type:        q.Type,
			Note:        q.Note,
			Lines:       q.Lines,
		}
}

func (q *CreateInventoryVoucherCommand) SetCreateInventoryVoucherArgs(args *CreateInventoryVoucherArgs) {
	q.ShopID = args.ShopID
	q.CreatedBy = args.CreatedBy
	q.Title = args.Title
	q.RefID = args.RefID
	q.RefType = args.RefType
	q.RefName = args.RefName
	q.TraderID = args.TraderID
	q.TotalAmount = args.TotalAmount
	q.Type = args.Type
	q.Note = args.Note
	q.Lines = args.Lines
}

func (q *UpdateInventoryVoucherCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateInventoryVoucherArgs) {
	return ctx,
		&UpdateInventoryVoucherArgs{
			ID:          q.ID,
			ShopID:      q.ShopID,
			Title:       q.Title,
			UpdatedBy:   q.UpdatedBy,
			TraderID:    q.TraderID,
			TotalAmount: q.TotalAmount,
			Note:        q.Note,
			Lines:       q.Lines,
		}
}

func (q *UpdateInventoryVoucherCommand) SetUpdateInventoryVoucherArgs(args *UpdateInventoryVoucherArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.Title = args.Title
	q.UpdatedBy = args.UpdatedBy
	q.TraderID = args.TraderID
	q.TotalAmount = args.TotalAmount
	q.Note = args.Note
	q.Lines = args.Lines
}

func (q *GetInventoriesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetInventoryRequest) {
	return ctx,
		&GetInventoryRequest{
			ShopID: q.ShopID,
			Paging: q.Paging,
		}
}

func (q *GetInventoriesQuery) SetGetInventoryRequest(args *GetInventoryRequest) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
}

func (q *GetInventoriesByVariantIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetInventoriesByVariantIDsArgs) {
	return ctx,
		&GetInventoriesByVariantIDsArgs{
			ShopID:     q.ShopID,
			Paging:     q.Paging,
			VariantIDs: q.VariantIDs,
		}
}

func (q *GetInventoriesByVariantIDsQuery) SetGetInventoriesByVariantIDsArgs(args *GetInventoriesByVariantIDsArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.VariantIDs = args.VariantIDs
}

func (q *GetInventoryQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID int64, VariantID int64) {
	return ctx,
		q.ShopID,
		q.VariantID
}

func (q *GetInventoryVoucherQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID int64, ID int64) {
	return ctx,
		q.ShopID,
		q.ID
}

func (q *GetInventoryVouchersQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID int64, Paging *meta.Paging) {
	return ctx,
		q.ShopID,
		q.Paging
}

func (q *GetInventoryVouchersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetInventoryVouchersByIDArgs) {
	return ctx,
		&GetInventoryVouchersByIDArgs{
			ShopID: q.ShopID,
			Paging: q.Paging,
			IDs:    q.IDs,
		}
}

func (q *GetInventoryVouchersByIDsQuery) SetGetInventoryVouchersByIDArgs(args *GetInventoryVouchersByIDArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.IDs = args.IDs
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
	b.AddHandler(h.HandleAdjustInventoryQuantity)
	b.AddHandler(h.HandleCancelInventoryVoucher)
	b.AddHandler(h.HandleConfirmInventoryVoucher)
	b.AddHandler(h.HandleCreateInventoryVariant)
	b.AddHandler(h.HandleCreateInventoryVoucher)
	b.AddHandler(h.HandleUpdateInventoryVoucher)
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
	b.AddHandler(h.HandleGetInventories)
	b.AddHandler(h.HandleGetInventoriesByVariantIDs)
	b.AddHandler(h.HandleGetInventory)
	b.AddHandler(h.HandleGetInventoryVoucher)
	b.AddHandler(h.HandleGetInventoryVouchers)
	b.AddHandler(h.HandleGetInventoryVouchersByIDs)
	return QueryBus{b}
}
