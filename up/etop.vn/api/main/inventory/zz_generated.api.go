// +build !generator

// Code generated by generator api. DO NOT EDIT.

package inventory

import (
	context "context"

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

type AdjustInventoryQuantityCommand struct {
	Overstock bool
	ShopID    dot.ID
	Lines     []*InventoryVariant
	Title     string
	UserID    dot.ID
	Note      string

	Result *AdjustInventoryQuantityRespone `json:"-"`
}

func (h AggregateHandler) HandleAdjustInventoryQuantity(ctx context.Context, msg *AdjustInventoryQuantityCommand) (err error) {
	msg.Result, err = h.inner.AdjustInventoryQuantity(msg.GetArgs(ctx))
	return err
}

type CancelInventoryVoucherCommand struct {
	ShopID    dot.ID
	ID        dot.ID
	UpdatedBy dot.ID
	Reason    string

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleCancelInventoryVoucher(ctx context.Context, msg *CancelInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.CancelInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type CheckInventoryVariantsQuantityCommand struct {
	Lines              []*InventoryVoucherItem
	InventoryOverStock bool
	ShopID             dot.ID
	Type               InventoryVoucherType

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleCheckInventoryVariantsQuantity(ctx context.Context, msg *CheckInventoryVariantsQuantityCommand) (err error) {
	return h.inner.CheckInventoryVariantsQuantity(msg.GetArgs(ctx))
}

type ConfirmInventoryVoucherCommand struct {
	ShopID    dot.ID
	ID        dot.ID
	UpdatedBy dot.ID

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleConfirmInventoryVoucher(ctx context.Context, msg *ConfirmInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.ConfirmInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type CreateInventoryVariantCommand struct {
	ShopID    dot.ID
	VariantID dot.ID

	Result *InventoryVariant `json:"-"`
}

func (h AggregateHandler) HandleCreateInventoryVariant(ctx context.Context, msg *CreateInventoryVariantCommand) (err error) {
	msg.Result, err = h.inner.CreateInventoryVariant(msg.GetArgs(ctx))
	return err
}

type CreateInventoryVoucherCommand struct {
	Overstock   bool
	ShopID      dot.ID
	CreatedBy   dot.ID
	Title       string
	RefID       dot.ID
	RefType     InventoryRefType
	RefName     InventoryVoucherRefName
	RefCode     string
	TraderID    dot.ID
	TotalAmount int
	Type        InventoryVoucherType
	Note        string
	Lines       []*InventoryVoucherItem

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleCreateInventoryVoucher(ctx context.Context, msg *CreateInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.CreateInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type CreateInventoryVoucherByQuantityChangeCommand struct {
	ShopID    dot.ID
	RefID     dot.ID
	RefType   InventoryRefType
	RefName   InventoryVoucherRefName
	RefCode   string
	NoteIn    string
	NoteOut   string
	Title     string
	Overstock bool
	CreatedBy dot.ID
	Lines     []*InventoryVariantQuantityChange

	Result *CreateInventoryVoucherByQuantityChangeResponse `json:"-"`
}

func (h AggregateHandler) HandleCreateInventoryVoucherByQuantityChange(ctx context.Context, msg *CreateInventoryVoucherByQuantityChangeCommand) (err error) {
	msg.Result, err = h.inner.CreateInventoryVoucherByQuantityChange(msg.GetArgs(ctx))
	return err
}

type CreateInventoryVoucherByReferenceCommand struct {
	RefType   InventoryRefType
	RefID     dot.ID
	Type      InventoryVoucherType
	ShopID    dot.ID
	UserID    dot.ID
	OverStock bool

	Result []*InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleCreateInventoryVoucherByReference(ctx context.Context, msg *CreateInventoryVoucherByReferenceCommand) (err error) {
	msg.Result, err = h.inner.CreateInventoryVoucherByReference(msg.GetArgs(ctx))
	return err
}

type UpdateInventoryVariantCostPriceCommand struct {
	ShopID    dot.ID
	VariantID dot.ID
	CostPrice int

	Result *InventoryVariant `json:"-"`
}

func (h AggregateHandler) HandleUpdateInventoryVariantCostPrice(ctx context.Context, msg *UpdateInventoryVariantCostPriceCommand) (err error) {
	msg.Result, err = h.inner.UpdateInventoryVariantCostPrice(msg.GetArgs(ctx))
	return err
}

type UpdateInventoryVoucherCommand struct {
	ID          dot.ID
	ShopID      dot.ID
	Title       dot.NullString
	UpdatedBy   dot.ID
	TraderID    dot.NullID
	TotalAmount int
	Note        dot.NullString
	Lines       []*InventoryVoucherItem

	Result *InventoryVoucher `json:"-"`
}

func (h AggregateHandler) HandleUpdateInventoryVoucher(ctx context.Context, msg *UpdateInventoryVoucherCommand) (err error) {
	msg.Result, err = h.inner.UpdateInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type GetInventoryVariantQuery struct {
	ShopID    dot.ID
	VariantID dot.ID

	Result *InventoryVariant `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVariant(ctx context.Context, msg *GetInventoryVariantQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVariant(msg.GetArgs(ctx))
	return err
}

type GetInventoryVariantsQuery struct {
	ShopID dot.ID
	Paging *meta.Paging

	Result *GetInventoryVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVariants(ctx context.Context, msg *GetInventoryVariantsQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVariants(msg.GetArgs(ctx))
	return err
}

type GetInventoryVariantsByVariantIDsQuery struct {
	ShopID     dot.ID
	Paging     *meta.Paging
	VariantIDs []dot.ID

	Result *GetInventoryVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVariantsByVariantIDs(ctx context.Context, msg *GetInventoryVariantsByVariantIDsQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVariantsByVariantIDs(msg.GetArgs(ctx))
	return err
}

type GetInventoryVoucherQuery struct {
	ShopID dot.ID
	ID     dot.ID

	Result *InventoryVoucher `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVoucher(ctx context.Context, msg *GetInventoryVoucherQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVoucher(msg.GetArgs(ctx))
	return err
}

type GetInventoryVoucherByReferenceQuery struct {
	ShopID  dot.ID
	RefID   dot.ID
	RefType InventoryRefType

	Result *GetInventoryVoucherByReferenceResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVoucherByReference(ctx context.Context, msg *GetInventoryVoucherByReferenceQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVoucherByReference(msg.GetArgs(ctx))
	return err
}

type GetInventoryVouchersQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *GetInventoryVouchersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVouchers(ctx context.Context, msg *GetInventoryVouchersQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVouchers(msg.GetArgs(ctx))
	return err
}

type GetInventoryVouchersByIDsQuery struct {
	ShopID dot.ID
	Paging *meta.Paging
	IDs    []dot.ID

	Result *GetInventoryVouchersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVouchersByIDs(ctx context.Context, msg *GetInventoryVouchersByIDsQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVouchersByIDs(msg.GetArgs(ctx))
	return err
}

type GetInventoryVouchersByRefIDsQuery struct {
	RefIDs []dot.ID
	ShopID dot.ID

	Result *GetInventoryVouchersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetInventoryVouchersByRefIDs(ctx context.Context, msg *GetInventoryVouchersByRefIDsQuery) (err error) {
	msg.Result, err = h.inner.GetInventoryVouchersByRefIDs(msg.GetArgs(ctx))
	return err
}

type ListInventoryVariantsByVariantIDsQuery struct {
	ShopID     dot.ID
	VariantIDs []dot.ID

	Result *GetInventoryVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListInventoryVariantsByVariantIDs(ctx context.Context, msg *ListInventoryVariantsByVariantIDsQuery) (err error) {
	msg.Result, err = h.inner.ListInventoryVariantsByVariantIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AdjustInventoryQuantityCommand) command()                {}
func (q *CancelInventoryVoucherCommand) command()                 {}
func (q *CheckInventoryVariantsQuantityCommand) command()         {}
func (q *ConfirmInventoryVoucherCommand) command()                {}
func (q *CreateInventoryVariantCommand) command()                 {}
func (q *CreateInventoryVoucherCommand) command()                 {}
func (q *CreateInventoryVoucherByQuantityChangeCommand) command() {}
func (q *CreateInventoryVoucherByReferenceCommand) command()      {}
func (q *UpdateInventoryVariantCostPriceCommand) command()        {}
func (q *UpdateInventoryVoucherCommand) command()                 {}

func (q *GetInventoryVariantQuery) query()               {}
func (q *GetInventoryVariantsQuery) query()              {}
func (q *GetInventoryVariantsByVariantIDsQuery) query()  {}
func (q *GetInventoryVoucherQuery) query()               {}
func (q *GetInventoryVoucherByReferenceQuery) query()    {}
func (q *GetInventoryVouchersQuery) query()              {}
func (q *GetInventoryVouchersByIDsQuery) query()         {}
func (q *GetInventoryVouchersByRefIDsQuery) query()      {}
func (q *ListInventoryVariantsByVariantIDsQuery) query() {}

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

func (q *CheckInventoryVariantsQuantityCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CheckInventoryVariantQuantityRequest) {
	return ctx,
		&CheckInventoryVariantQuantityRequest{
			Lines:              q.Lines,
			InventoryOverStock: q.InventoryOverStock,
			ShopID:             q.ShopID,
			Type:               q.Type,
		}
}

func (q *CheckInventoryVariantsQuantityCommand) SetCheckInventoryVariantQuantityRequest(args *CheckInventoryVariantQuantityRequest) {
	q.Lines = args.Lines
	q.InventoryOverStock = args.InventoryOverStock
	q.ShopID = args.ShopID
	q.Type = args.Type
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
			RefCode:     q.RefCode,
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
	q.RefCode = args.RefCode
	q.TraderID = args.TraderID
	q.TotalAmount = args.TotalAmount
	q.Type = args.Type
	q.Note = args.Note
	q.Lines = args.Lines
}

func (q *CreateInventoryVoucherByQuantityChangeCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateInventoryVoucherByQuantityChangeRequest) {
	return ctx,
		&CreateInventoryVoucherByQuantityChangeRequest{
			ShopID:    q.ShopID,
			RefID:     q.RefID,
			RefType:   q.RefType,
			RefName:   q.RefName,
			RefCode:   q.RefCode,
			NoteIn:    q.NoteIn,
			NoteOut:   q.NoteOut,
			Title:     q.Title,
			Overstock: q.Overstock,
			CreatedBy: q.CreatedBy,
			Lines:     q.Lines,
		}
}

func (q *CreateInventoryVoucherByQuantityChangeCommand) SetCreateInventoryVoucherByQuantityChangeRequest(args *CreateInventoryVoucherByQuantityChangeRequest) {
	q.ShopID = args.ShopID
	q.RefID = args.RefID
	q.RefType = args.RefType
	q.RefName = args.RefName
	q.RefCode = args.RefCode
	q.NoteIn = args.NoteIn
	q.NoteOut = args.NoteOut
	q.Title = args.Title
	q.Overstock = args.Overstock
	q.CreatedBy = args.CreatedBy
	q.Lines = args.Lines
}

func (q *CreateInventoryVoucherByReferenceCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateInventoryVoucherByReferenceArgs) {
	return ctx,
		&CreateInventoryVoucherByReferenceArgs{
			RefType:   q.RefType,
			RefID:     q.RefID,
			Type:      q.Type,
			ShopID:    q.ShopID,
			UserID:    q.UserID,
			OverStock: q.OverStock,
		}
}

func (q *CreateInventoryVoucherByReferenceCommand) SetCreateInventoryVoucherByReferenceArgs(args *CreateInventoryVoucherByReferenceArgs) {
	q.RefType = args.RefType
	q.RefID = args.RefID
	q.Type = args.Type
	q.ShopID = args.ShopID
	q.UserID = args.UserID
	q.OverStock = args.OverStock
}

func (q *UpdateInventoryVariantCostPriceCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateInventoryVariantCostPriceRequest) {
	return ctx,
		&UpdateInventoryVariantCostPriceRequest{
			ShopID:    q.ShopID,
			VariantID: q.VariantID,
			CostPrice: q.CostPrice,
		}
}

func (q *UpdateInventoryVariantCostPriceCommand) SetUpdateInventoryVariantCostPriceRequest(args *UpdateInventoryVariantCostPriceRequest) {
	q.ShopID = args.ShopID
	q.VariantID = args.VariantID
	q.CostPrice = args.CostPrice
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

func (q *GetInventoryVariantQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID, VariantID dot.ID) {
	return ctx,
		q.ShopID,
		q.VariantID
}

func (q *GetInventoryVariantsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetInventoryRequest) {
	return ctx,
		&GetInventoryRequest{
			ShopID: q.ShopID,
			Paging: q.Paging,
		}
}

func (q *GetInventoryVariantsQuery) SetGetInventoryRequest(args *GetInventoryRequest) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
}

func (q *GetInventoryVariantsByVariantIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetInventoryVariantsByVariantIDsArgs) {
	return ctx,
		&GetInventoryVariantsByVariantIDsArgs{
			ShopID:     q.ShopID,
			Paging:     q.Paging,
			VariantIDs: q.VariantIDs,
		}
}

func (q *GetInventoryVariantsByVariantIDsQuery) SetGetInventoryVariantsByVariantIDsArgs(args *GetInventoryVariantsByVariantIDsArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.VariantIDs = args.VariantIDs
}

func (q *GetInventoryVoucherQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID, ID dot.ID) {
	return ctx,
		q.ShopID,
		q.ID
}

func (q *GetInventoryVoucherByReferenceQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID, refID dot.ID, refType InventoryRefType) {
	return ctx,
		q.ShopID,
		q.RefID,
		q.RefType
}

func (q *GetInventoryVouchersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListInventoryVouchersArgs) {
	return ctx,
		&ListInventoryVouchersArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *GetInventoryVouchersQuery) SetListInventoryVouchersArgs(args *ListInventoryVouchersArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
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

func (q *GetInventoryVouchersByRefIDsQuery) GetArgs(ctx context.Context) (_ context.Context, RefIDs []dot.ID, ShopID dot.ID) {
	return ctx,
		q.RefIDs,
		q.ShopID
}

func (q *ListInventoryVariantsByVariantIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListInventoryVariantsByVariantIDsArgs) {
	return ctx,
		&ListInventoryVariantsByVariantIDsArgs{
			ShopID:     q.ShopID,
			VariantIDs: q.VariantIDs,
		}
}

func (q *ListInventoryVariantsByVariantIDsQuery) SetListInventoryVariantsByVariantIDsArgs(args *ListInventoryVariantsByVariantIDsArgs) {
	q.ShopID = args.ShopID
	q.VariantIDs = args.VariantIDs
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
	b.AddHandler(h.HandleCheckInventoryVariantsQuantity)
	b.AddHandler(h.HandleConfirmInventoryVoucher)
	b.AddHandler(h.HandleCreateInventoryVariant)
	b.AddHandler(h.HandleCreateInventoryVoucher)
	b.AddHandler(h.HandleCreateInventoryVoucherByQuantityChange)
	b.AddHandler(h.HandleCreateInventoryVoucherByReference)
	b.AddHandler(h.HandleUpdateInventoryVariantCostPrice)
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
	b.AddHandler(h.HandleGetInventoryVariant)
	b.AddHandler(h.HandleGetInventoryVariants)
	b.AddHandler(h.HandleGetInventoryVariantsByVariantIDs)
	b.AddHandler(h.HandleGetInventoryVoucher)
	b.AddHandler(h.HandleGetInventoryVoucherByReference)
	b.AddHandler(h.HandleGetInventoryVouchers)
	b.AddHandler(h.HandleGetInventoryVouchersByIDs)
	b.AddHandler(h.HandleGetInventoryVouchersByRefIDs)
	b.AddHandler(h.HandleListInventoryVariantsByVariantIDs)
	return QueryBus{b}
}
