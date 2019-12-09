// +build !generator

// Code generated by generator api. DO NOT EDIT.

package suppliering

import (
	context "context"

	meta "etop.vn/api/meta"
	shopping "etop.vn/api/shopping"
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

type CreateSupplierCommand struct {
	ShopID            dot.ID
	FullName          string
	Phone             string
	Email             string
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string

	Result *ShopSupplier `json:"-"`
}

func (h AggregateHandler) HandleCreateSupplier(ctx context.Context, msg *CreateSupplierCommand) (err error) {
	msg.Result, err = h.inner.CreateSupplier(msg.GetArgs(ctx))
	return err
}

type DeleteSupplierCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteSupplier(ctx context.Context, msg *DeleteSupplierCommand) (err error) {
	msg.Result, err = h.inner.DeleteSupplier(msg.GetArgs(ctx))
	return err
}

type UpdateSupplierCommand struct {
	ID                dot.ID
	ShopID            dot.ID
	FullName          dot.NullString
	Note              dot.NullString
	Phone             dot.NullString
	Email             dot.NullString
	CompanyName       dot.NullString
	TaxNumber         dot.NullString
	HeadquaterAddress dot.NullString

	Result *ShopSupplier `json:"-"`
}

func (h AggregateHandler) HandleUpdateSupplier(ctx context.Context, msg *UpdateSupplierCommand) (err error) {
	msg.Result, err = h.inner.UpdateSupplier(msg.GetArgs(ctx))
	return err
}

type GetSupplierByIDQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *ShopSupplier `json:"-"`
}

func (h QueryServiceHandler) HandleGetSupplierByID(ctx context.Context, msg *GetSupplierByIDQuery) (err error) {
	msg.Result, err = h.inner.GetSupplierByID(msg.GetArgs(ctx))
	return err
}

type ListSuppliersQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *SuppliersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListSuppliers(ctx context.Context, msg *ListSuppliersQuery) (err error) {
	msg.Result, err = h.inner.ListSuppliers(msg.GetArgs(ctx))
	return err
}

type ListSuppliersByIDsQuery struct {
	IDs    []dot.ID
	ShopID dot.ID
	Paging meta.Paging

	Result *SuppliersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListSuppliersByIDs(ctx context.Context, msg *ListSuppliersByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListSuppliersByIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateSupplierCommand) command() {}
func (q *DeleteSupplierCommand) command() {}
func (q *UpdateSupplierCommand) command() {}

func (q *GetSupplierByIDQuery) query()    {}
func (q *ListSuppliersQuery) query()      {}
func (q *ListSuppliersByIDsQuery) query() {}

// implement conversion

func (q *CreateSupplierCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateSupplierArgs) {
	return ctx,
		&CreateSupplierArgs{
			ShopID:            q.ShopID,
			FullName:          q.FullName,
			Phone:             q.Phone,
			Email:             q.Email,
			CompanyName:       q.CompanyName,
			TaxNumber:         q.TaxNumber,
			HeadquaterAddress: q.HeadquaterAddress,
			Note:              q.Note,
		}
}

func (q *CreateSupplierCommand) SetCreateSupplierArgs(args *CreateSupplierArgs) {
	q.ShopID = args.ShopID
	q.FullName = args.FullName
	q.Phone = args.Phone
	q.Email = args.Email
	q.CompanyName = args.CompanyName
	q.TaxNumber = args.TaxNumber
	q.HeadquaterAddress = args.HeadquaterAddress
	q.Note = args.Note
}

func (q *DeleteSupplierCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, shopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *UpdateSupplierCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateSupplierArgs) {
	return ctx,
		&UpdateSupplierArgs{
			ID:                q.ID,
			ShopID:            q.ShopID,
			FullName:          q.FullName,
			Note:              q.Note,
			Phone:             q.Phone,
			Email:             q.Email,
			CompanyName:       q.CompanyName,
			TaxNumber:         q.TaxNumber,
			HeadquaterAddress: q.HeadquaterAddress,
		}
}

func (q *UpdateSupplierCommand) SetUpdateSupplierArgs(args *UpdateSupplierArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.FullName = args.FullName
	q.Note = args.Note
	q.Phone = args.Phone
	q.Email = args.Email
	q.CompanyName = args.CompanyName
	q.TaxNumber = args.TaxNumber
	q.HeadquaterAddress = args.HeadquaterAddress
}

func (q *GetSupplierByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDQueryShopArg) {
	return ctx,
		&shopping.IDQueryShopArg{
			ID:     q.ID,
			ShopID: q.ShopID,
		}
}

func (q *GetSupplierByIDQuery) SetIDQueryShopArg(args *shopping.IDQueryShopArg) {
	q.ID = args.ID
	q.ShopID = args.ShopID
}

func (q *ListSuppliersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListSuppliersQuery) SetListQueryShopArgs(args *shopping.ListQueryShopArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *ListSuppliersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
			Paging: q.Paging,
		}
}

func (q *ListSuppliersByIDsQuery) SetIDsQueryShopArgs(args *shopping.IDsQueryShopArgs) {
	q.IDs = args.IDs
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
	b.AddHandler(h.HandleCreateSupplier)
	b.AddHandler(h.HandleDeleteSupplier)
	b.AddHandler(h.HandleUpdateSupplier)
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
	b.AddHandler(h.HandleGetSupplierByID)
	b.AddHandler(h.HandleListSuppliers)
	b.AddHandler(h.HandleListSuppliersByIDs)
	return QueryBus{b}
}
