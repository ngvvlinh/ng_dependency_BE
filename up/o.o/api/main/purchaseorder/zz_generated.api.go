// +build !generator

// Code generated by generator api. DO NOT EDIT.

package purchaseorder

import (
	context "context"

	meta "o.o/api/meta"
	shopping "o.o/api/shopping"
	inttypes "o.o/api/top/int/types"
	inventory_auto "o.o/api/top/types/etc/inventory_auto"
	status3 "o.o/api/top/types/etc/status3"
	capi "o.o/capi"
	dot "o.o/capi/dot"
	filter "o.o/capi/filter"
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

type CancelPurchaseOrderCommand struct {
	ID                   dot.ID
	ShopID               dot.ID
	CancelReason         string
	UpdatedBy            dot.ID
	InventoryOverStock   bool
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher

	Result int `json:"-"`
}

func (h AggregateHandler) HandleCancelPurchaseOrder(ctx context.Context, msg *CancelPurchaseOrderCommand) (err error) {
	msg.Result, err = h.inner.CancelPurchaseOrder(msg.GetArgs(ctx))
	return err
}

type ConfirmPurchaseOrderCommand struct {
	ID                   dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	ShopID               dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleConfirmPurchaseOrder(ctx context.Context, msg *ConfirmPurchaseOrderCommand) (err error) {
	msg.Result, err = h.inner.ConfirmPurchaseOrder(msg.GetArgs(ctx))
	return err
}

type CreatePurchaseOrderCommand struct {
	ShopID        dot.ID
	SupplierID    dot.ID
	BasketValue   int
	DiscountLines []*inttypes.DiscountLine
	TotalDiscount int
	FeeLines      []*inttypes.FeeLine
	TotalFee      int
	TotalAmount   int
	Note          string
	Lines         []*PurchaseOrderLine
	CreatedBy     dot.ID

	Result *PurchaseOrder `json:"-"`
}

func (h AggregateHandler) HandleCreatePurchaseOrder(ctx context.Context, msg *CreatePurchaseOrderCommand) (err error) {
	msg.Result, err = h.inner.CreatePurchaseOrder(msg.GetArgs(ctx))
	return err
}

type DeletePurchaseOrderCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeletePurchaseOrder(ctx context.Context, msg *DeletePurchaseOrderCommand) (err error) {
	msg.Result, err = h.inner.DeletePurchaseOrder(msg.GetArgs(ctx))
	return err
}

type UpdatePurchaseOrderCommand struct {
	ID            dot.ID
	ShopID        dot.ID
	BasketValue   dot.NullInt
	DiscountLines []*inttypes.DiscountLine
	TotalDiscount dot.NullInt
	FeeLines      []*inttypes.FeeLine
	TotalFee      dot.NullInt
	TotalAmount   dot.NullInt
	Note          dot.NullString
	Lines         []*PurchaseOrderLine

	Result *PurchaseOrder `json:"-"`
}

func (h AggregateHandler) HandleUpdatePurchaseOrder(ctx context.Context, msg *UpdatePurchaseOrderCommand) (err error) {
	msg.Result, err = h.inner.UpdatePurchaseOrder(msg.GetArgs(ctx))
	return err
}

type GetPurchaseOrderByIDQuery struct {
	ID             dot.ID
	ShopID         dot.ID
	IncludeDeleted bool

	Result *PurchaseOrder `json:"-"`
}

func (h QueryServiceHandler) HandleGetPurchaseOrderByID(ctx context.Context, msg *GetPurchaseOrderByIDQuery) (err error) {
	msg.Result, err = h.inner.GetPurchaseOrderByID(msg.GetArgs(ctx))
	return err
}

type GetPurchaseOrdersByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID

	Result *PurchaseOrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetPurchaseOrdersByIDs(ctx context.Context, msg *GetPurchaseOrdersByIDsQuery) (err error) {
	msg.Result, err = h.inner.GetPurchaseOrdersByIDs(msg.GetArgs(ctx))
	return err
}

type ListPurchaseOrdersQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
	Name    filter.FullTextSearch

	Result *PurchaseOrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListPurchaseOrders(ctx context.Context, msg *ListPurchaseOrdersQuery) (err error) {
	msg.Result, err = h.inner.ListPurchaseOrders(msg.GetArgs(ctx))
	return err
}

type ListPurchaseOrdersByReceiptIDQuery struct {
	ReceiptID dot.ID
	ShopID    dot.ID

	Result *PurchaseOrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListPurchaseOrdersByReceiptID(ctx context.Context, msg *ListPurchaseOrdersByReceiptIDQuery) (err error) {
	msg.Result, err = h.inner.ListPurchaseOrdersByReceiptID(msg.GetArgs(ctx))
	return err
}

type ListPurchaseOrdersBySupplierIDsAndStatusesQuery struct {
	ShopID      dot.ID
	SupplierIDs []dot.ID
	Statuses    []status3.Status

	Result *PurchaseOrdersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListPurchaseOrdersBySupplierIDsAndStatuses(ctx context.Context, msg *ListPurchaseOrdersBySupplierIDsAndStatusesQuery) (err error) {
	msg.Result, err = h.inner.ListPurchaseOrdersBySupplierIDsAndStatuses(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CancelPurchaseOrderCommand) command()  {}
func (q *ConfirmPurchaseOrderCommand) command() {}
func (q *CreatePurchaseOrderCommand) command()  {}
func (q *DeletePurchaseOrderCommand) command()  {}
func (q *UpdatePurchaseOrderCommand) command()  {}

func (q *GetPurchaseOrderByIDQuery) query()                       {}
func (q *GetPurchaseOrdersByIDsQuery) query()                     {}
func (q *ListPurchaseOrdersQuery) query()                         {}
func (q *ListPurchaseOrdersByReceiptIDQuery) query()              {}
func (q *ListPurchaseOrdersBySupplierIDsAndStatusesQuery) query() {}

// implement conversion

func (q *CancelPurchaseOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CancelPurchaseOrderArgs) {
	return ctx,
		&CancelPurchaseOrderArgs{
			ID:                   q.ID,
			ShopID:               q.ShopID,
			CancelReason:         q.CancelReason,
			UpdatedBy:            q.UpdatedBy,
			InventoryOverStock:   q.InventoryOverStock,
			AutoInventoryVoucher: q.AutoInventoryVoucher,
		}
}

func (q *CancelPurchaseOrderCommand) SetCancelPurchaseOrderArgs(args *CancelPurchaseOrderArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.CancelReason = args.CancelReason
	q.UpdatedBy = args.UpdatedBy
	q.InventoryOverStock = args.InventoryOverStock
	q.AutoInventoryVoucher = args.AutoInventoryVoucher
}

func (q *ConfirmPurchaseOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmPurchaseOrderArgs) {
	return ctx,
		&ConfirmPurchaseOrderArgs{
			ID:                   q.ID,
			AutoInventoryVoucher: q.AutoInventoryVoucher,
			ShopID:               q.ShopID,
		}
}

func (q *ConfirmPurchaseOrderCommand) SetConfirmPurchaseOrderArgs(args *ConfirmPurchaseOrderArgs) {
	q.ID = args.ID
	q.AutoInventoryVoucher = args.AutoInventoryVoucher
	q.ShopID = args.ShopID
}

func (q *CreatePurchaseOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreatePurchaseOrderArgs) {
	return ctx,
		&CreatePurchaseOrderArgs{
			ShopID:        q.ShopID,
			SupplierID:    q.SupplierID,
			BasketValue:   q.BasketValue,
			DiscountLines: q.DiscountLines,
			TotalDiscount: q.TotalDiscount,
			FeeLines:      q.FeeLines,
			TotalFee:      q.TotalFee,
			TotalAmount:   q.TotalAmount,
			Note:          q.Note,
			Lines:         q.Lines,
			CreatedBy:     q.CreatedBy,
		}
}

func (q *CreatePurchaseOrderCommand) SetCreatePurchaseOrderArgs(args *CreatePurchaseOrderArgs) {
	q.ShopID = args.ShopID
	q.SupplierID = args.SupplierID
	q.BasketValue = args.BasketValue
	q.DiscountLines = args.DiscountLines
	q.TotalDiscount = args.TotalDiscount
	q.FeeLines = args.FeeLines
	q.TotalFee = args.TotalFee
	q.TotalAmount = args.TotalAmount
	q.Note = args.Note
	q.Lines = args.Lines
	q.CreatedBy = args.CreatedBy
}

func (q *DeletePurchaseOrderCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, shopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *UpdatePurchaseOrderCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdatePurchaseOrderArgs) {
	return ctx,
		&UpdatePurchaseOrderArgs{
			ID:            q.ID,
			ShopID:        q.ShopID,
			BasketValue:   q.BasketValue,
			DiscountLines: q.DiscountLines,
			TotalDiscount: q.TotalDiscount,
			FeeLines:      q.FeeLines,
			TotalFee:      q.TotalFee,
			TotalAmount:   q.TotalAmount,
			Note:          q.Note,
			Lines:         q.Lines,
		}
}

func (q *UpdatePurchaseOrderCommand) SetUpdatePurchaseOrderArgs(args *UpdatePurchaseOrderArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.BasketValue = args.BasketValue
	q.DiscountLines = args.DiscountLines
	q.TotalDiscount = args.TotalDiscount
	q.FeeLines = args.FeeLines
	q.TotalFee = args.TotalFee
	q.TotalAmount = args.TotalAmount
	q.Note = args.Note
	q.Lines = args.Lines
}

func (q *GetPurchaseOrderByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:             q.ID,
			ShopID:         q.ShopID,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *GetPurchaseOrderByIDQuery) SetIDQueryShopArg(args *shopping.IDQueryShopArg) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.IncludeDeleted = args.IncludeDeleted
}

func (q *GetPurchaseOrdersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID, ShopID dot.ID) {
	return ctx,
		q.IDs,
		q.ShopID
}

func (q *ListPurchaseOrdersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
			Name:    q.Name,
		}
}

func (q *ListPurchaseOrdersQuery) SetListQueryShopArgs(args *shopping.ListQueryShopArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
	q.Name = args.Name
}

func (q *ListPurchaseOrdersByReceiptIDQuery) GetArgs(ctx context.Context) (_ context.Context, receiptID dot.ID, shopID dot.ID) {
	return ctx,
		q.ReceiptID,
		q.ShopID
}

func (q *ListPurchaseOrdersBySupplierIDsAndStatusesQuery) GetArgs(ctx context.Context) (_ context.Context, shopID dot.ID, supplierIDs []dot.ID, statuses []status3.Status) {
	return ctx,
		q.ShopID,
		q.SupplierIDs,
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
	b.AddHandler(h.HandleCancelPurchaseOrder)
	b.AddHandler(h.HandleConfirmPurchaseOrder)
	b.AddHandler(h.HandleCreatePurchaseOrder)
	b.AddHandler(h.HandleDeletePurchaseOrder)
	b.AddHandler(h.HandleUpdatePurchaseOrder)
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
	b.AddHandler(h.HandleGetPurchaseOrderByID)
	b.AddHandler(h.HandleGetPurchaseOrdersByIDs)
	b.AddHandler(h.HandleListPurchaseOrders)
	b.AddHandler(h.HandleListPurchaseOrdersByReceiptID)
	b.AddHandler(h.HandleListPurchaseOrdersBySupplierIDsAndStatuses)
	return QueryBus{b}
}
