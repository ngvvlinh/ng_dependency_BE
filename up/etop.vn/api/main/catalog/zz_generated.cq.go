// Code generated by generator cq. DO NOT EDIT.

// +build !generator

package catalog

import (
	context "context"

	types "etop.vn/api/main/catalog/types"
	meta "etop.vn/api/meta"
	shopping "etop.vn/api/shopping"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus meta.Bus }
type QueryBus struct{ bus meta.Bus }

func (c CommandBus) Dispatch(ctx context.Context, msg Command) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error {
	return c.bus.Dispatch(ctx, msg)
}
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

type CreateShopProductCommand struct {
	ShopID          int64
	Code            string
	Name            string
	Unit            string
	ImageURLs       []string
	Note            string
	DescriptionInfo DescriptionInfo
	PriceInfo       PriceInfo

	Result *ShopProduct `json:"-"`
}

func (h AggregateHandler) HandleCreateShopProduct(ctx context.Context, msg *CreateShopProductCommand) (err error) {
	msg.Result, err = h.inner.CreateShopProduct(msg.GetArgs(ctx))
	return err
}

type CreateShopVariantCommand struct {
	ShopID          int64
	ProductID       int64
	Code            string
	Name            string
	ImageURLs       []string
	Note            string
	Attributes      types.Attributes
	DescriptionInfo DescriptionInfo
	PriceInfo       PriceInfo

	Result *ShopVariant `json:"-"`
}

func (h AggregateHandler) HandleCreateShopVariant(ctx context.Context, msg *CreateShopVariantCommand) (err error) {
	msg.Result, err = h.inner.CreateShopVariant(msg.GetArgs(ctx))
	return err
}

type DeleteShopProductsCommand struct {
	IDs    []int64
	ShopID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleDeleteShopProducts(ctx context.Context, msg *DeleteShopProductsCommand) (err error) {
	msg.Result, err = h.inner.DeleteShopProducts(msg.GetArgs(ctx))
	return err
}

type DeleteShopVariantsCommand struct {
	IDs    []int64
	ShopID int64

	Result *meta.Empty `json:"-"`
}

func (h AggregateHandler) HandleDeleteShopVariants(ctx context.Context, msg *DeleteShopVariantsCommand) (err error) {
	msg.Result, err = h.inner.DeleteShopVariants(msg.GetArgs(ctx))
	return err
}

type UpdateShopProductImagesCommand struct {
	ID      int64
	ShopID  int64
	Updates []*meta.UpdateSet

	Result *ShopProduct `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopProductImages(ctx context.Context, msg *UpdateShopProductImagesCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopProductImages(msg.GetArgs(ctx))
	return err
}

type UpdateShopProductInfoCommand struct {
	ShopID          int64
	ProductID       int64
	Code            *string
	Name            *string
	Unit            *string
	Note            *string
	DescriptionInfo *DescriptionInfo

	Result *ShopProduct `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopProductInfo(ctx context.Context, msg *UpdateShopProductInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopProductInfo(msg.GetArgs(ctx))
	return err
}

type UpdateShopProductStatusCommand struct {
	IDs    []int64
	ShopID int64
	Status int16

	Result *ShopProduct `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopProductStatus(ctx context.Context, msg *UpdateShopProductStatusCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopProductStatus(msg.GetArgs(ctx))
	return err
}

type UpdateShopVariantImagesCommand struct {
	ID      int64
	ShopID  int64
	Updates []*meta.UpdateSet

	Result *ShopVariant `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopVariantImages(ctx context.Context, msg *UpdateShopVariantImagesCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopVariantImages(msg.GetArgs(ctx))
	return err
}

type UpdateShopVariantInfoCommand struct {
	ShopID          int64
	VariantID       int64
	Code            *string
	Name            *string
	Unit            *string
	Note            *string
	DescriptionInfo *DescriptionInfo

	Result *ShopVariant `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopVariantInfo(ctx context.Context, msg *UpdateShopVariantInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopVariantInfo(msg.GetArgs(ctx))
	return err
}

type UpdateShopVariantStatusCommand struct {
	IDs    []int64
	ShopID int64
	Status int16

	Result *ShopVariant `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopVariantStatus(ctx context.Context, msg *UpdateShopVariantStatusCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopVariantStatus(msg.GetArgs(ctx))
	return err
}

type GetShopProductByIDQuery struct {
	ProductID int64
	ShopID    int64

	Result *ShopProduct `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopProductByID(ctx context.Context, msg *GetShopProductByIDQuery) (err error) {
	msg.Result, err = h.inner.GetShopProductByID(msg.GetArgs(ctx))
	return err
}

type GetShopProductWithVariantsByIDQuery struct {
	ProductID int64
	ShopID    int64

	Result *ShopProductWithVariants `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopProductWithVariantsByID(ctx context.Context, msg *GetShopProductWithVariantsByIDQuery) (err error) {
	msg.Result, err = h.inner.GetShopProductWithVariantsByID(msg.GetArgs(ctx))
	return err
}

type GetShopVariantByIDQuery struct {
	VariantID int64
	ShopID    int64

	Result *ShopVariant `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopVariantByID(ctx context.Context, msg *GetShopVariantByIDQuery) (err error) {
	msg.Result, err = h.inner.GetShopVariantByID(msg.GetArgs(ctx))
	return err
}

type GetShopVariantWithProductByIDQuery struct {
	VariantID int64
	ShopID    int64

	Result *ShopVariantWithProduct `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopVariantWithProductByID(ctx context.Context, msg *GetShopVariantWithProductByIDQuery) (err error) {
	msg.Result, err = h.inner.GetShopVariantWithProductByID(msg.GetArgs(ctx))
	return err
}

type ListShopProductsQuery struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters

	Result *ShopProductsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopProducts(ctx context.Context, msg *ListShopProductsQuery) (err error) {
	msg.Result, err = h.inner.ListShopProducts(msg.GetArgs(ctx))
	return err
}

type ListShopProductsByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *ShopProductsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopProductsByIDs(ctx context.Context, msg *ListShopProductsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListShopProductsByIDs(msg.GetArgs(ctx))
	return err
}

type ListShopProductsWithVariantsQuery struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters

	Result *ShopProductsWithVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopProductsWithVariants(ctx context.Context, msg *ListShopProductsWithVariantsQuery) (err error) {
	msg.Result, err = h.inner.ListShopProductsWithVariants(msg.GetArgs(ctx))
	return err
}

type ListShopProductsWithVariantsByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *ShopProductsWithVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopProductsWithVariantsByIDs(ctx context.Context, msg *ListShopProductsWithVariantsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListShopProductsWithVariantsByIDs(msg.GetArgs(ctx))
	return err
}

type ListShopVariantsQuery struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters

	Result *ShopVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopVariants(ctx context.Context, msg *ListShopVariantsQuery) (err error) {
	msg.Result, err = h.inner.ListShopVariants(msg.GetArgs(ctx))
	return err
}

type ListShopVariantsByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *ShopVariantsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopVariantsByIDs(ctx context.Context, msg *ListShopVariantsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListShopVariantsByIDs(msg.GetArgs(ctx))
	return err
}

type ListShopVariantsWithProductByIDsQuery struct {
	IDs    []int64
	ShopID int64

	Result *ShopVariantsWithProductResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopVariantsWithProductByIDs(ctx context.Context, msg *ListShopVariantsWithProductByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListShopVariantsWithProductByIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateShopProductCommand) command()             {}
func (q *CreateShopVariantCommand) command()             {}
func (q *DeleteShopProductsCommand) command()            {}
func (q *DeleteShopVariantsCommand) command()            {}
func (q *UpdateShopProductImagesCommand) command()       {}
func (q *UpdateShopProductInfoCommand) command()         {}
func (q *UpdateShopProductStatusCommand) command()       {}
func (q *UpdateShopVariantImagesCommand) command()       {}
func (q *UpdateShopVariantInfoCommand) command()         {}
func (q *UpdateShopVariantStatusCommand) command()       {}
func (q *GetShopProductByIDQuery) query()                {}
func (q *GetShopProductWithVariantsByIDQuery) query()    {}
func (q *GetShopVariantByIDQuery) query()                {}
func (q *GetShopVariantWithProductByIDQuery) query()     {}
func (q *ListShopProductsQuery) query()                  {}
func (q *ListShopProductsByIDsQuery) query()             {}
func (q *ListShopProductsWithVariantsQuery) query()      {}
func (q *ListShopProductsWithVariantsByIDsQuery) query() {}
func (q *ListShopVariantsQuery) query()                  {}
func (q *ListShopVariantsByIDsQuery) query()             {}
func (q *ListShopVariantsWithProductByIDsQuery) query()  {}

// implement conversion

func (q *CreateShopProductCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShopProductArgs) {
	return ctx,
		&CreateShopProductArgs{
			ShopID:          q.ShopID,
			Code:            q.Code,
			Name:            q.Name,
			Unit:            q.Unit,
			ImageURLs:       q.ImageURLs,
			Note:            q.Note,
			DescriptionInfo: q.DescriptionInfo,
			PriceInfo:       q.PriceInfo,
		}
}

func (q *CreateShopVariantCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShopVariantArgs) {
	return ctx,
		&CreateShopVariantArgs{
			ShopID:          q.ShopID,
			ProductID:       q.ProductID,
			Code:            q.Code,
			Name:            q.Name,
			ImageURLs:       q.ImageURLs,
			Note:            q.Note,
			Attributes:      q.Attributes,
			DescriptionInfo: q.DescriptionInfo,
			PriceInfo:       q.PriceInfo,
		}
}

func (q *DeleteShopProductsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *DeleteShopVariantsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *UpdateShopProductImagesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateImagesArgs) {
	return ctx,
		&UpdateImagesArgs{
			ID:      q.ID,
			ShopID:  q.ShopID,
			Updates: q.Updates,
		}
}

func (q *UpdateShopProductInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShopProductInfoArgs) {
	return ctx,
		&UpdateShopProductInfoArgs{
			ShopID:          q.ShopID,
			ProductID:       q.ProductID,
			Code:            q.Code,
			Name:            q.Name,
			Unit:            q.Unit,
			Note:            q.Note,
			DescriptionInfo: q.DescriptionInfo,
		}
}

func (q *UpdateShopProductStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateStatusArgs) {
	return ctx,
		&UpdateStatusArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
			Status: q.Status,
		}
}

func (q *UpdateShopVariantImagesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateImagesArgs) {
	return ctx,
		&UpdateImagesArgs{
			ID:      q.ID,
			ShopID:  q.ShopID,
			Updates: q.Updates,
		}
}

func (q *UpdateShopVariantInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShopVariantInfoArgs) {
	return ctx,
		&UpdateShopVariantInfoArgs{
			ShopID:          q.ShopID,
			VariantID:       q.VariantID,
			Code:            q.Code,
			Name:            q.Name,
			Unit:            q.Unit,
			Note:            q.Note,
			DescriptionInfo: q.DescriptionInfo,
		}
}

func (q *UpdateShopVariantStatusCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateStatusArgs) {
	return ctx,
		&UpdateStatusArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
			Status: q.Status,
		}
}

func (q *GetShopProductByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetShopProductByIDQueryArgs) {
	return ctx,
		&GetShopProductByIDQueryArgs{
			ProductID: q.ProductID,
			ShopID:    q.ShopID,
		}
}

func (q *GetShopProductWithVariantsByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetShopProductByIDQueryArgs) {
	return ctx,
		&GetShopProductByIDQueryArgs{
			ProductID: q.ProductID,
			ShopID:    q.ShopID,
		}
}

func (q *GetShopVariantByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetShopVariantByIDQueryArgs) {
	return ctx,
		&GetShopVariantByIDQueryArgs{
			VariantID: q.VariantID,
			ShopID:    q.ShopID,
		}
}

func (q *GetShopVariantWithProductByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetShopVariantByIDQueryArgs) {
	return ctx,
		&GetShopVariantByIDQueryArgs{
			VariantID: q.VariantID,
			ShopID:    q.ShopID,
		}
}

func (q *ListShopProductsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListShopProductsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *ListShopProductsWithVariantsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListShopProductsWithVariantsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *ListShopVariantsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.ListQueryShopArgs) {
	return ctx,
		&shopping.ListQueryShopArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListShopVariantsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
	return ctx,
		&shopping.IDsQueryShopArgs{
			IDs:    q.IDs,
			ShopID: q.ShopID,
		}
}

func (q *ListShopVariantsWithProductByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *shopping.IDsQueryShopArgs) {
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
	meta.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCreateShopProduct)
	b.AddHandler(h.HandleCreateShopVariant)
	b.AddHandler(h.HandleDeleteShopProducts)
	b.AddHandler(h.HandleDeleteShopVariants)
	b.AddHandler(h.HandleUpdateShopProductImages)
	b.AddHandler(h.HandleUpdateShopProductInfo)
	b.AddHandler(h.HandleUpdateShopProductStatus)
	b.AddHandler(h.HandleUpdateShopVariantImages)
	b.AddHandler(h.HandleUpdateShopVariantInfo)
	b.AddHandler(h.HandleUpdateShopVariantStatus)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetShopProductByID)
	b.AddHandler(h.HandleGetShopProductWithVariantsByID)
	b.AddHandler(h.HandleGetShopVariantByID)
	b.AddHandler(h.HandleGetShopVariantWithProductByID)
	b.AddHandler(h.HandleListShopProducts)
	b.AddHandler(h.HandleListShopProductsByIDs)
	b.AddHandler(h.HandleListShopProductsWithVariants)
	b.AddHandler(h.HandleListShopProductsWithVariantsByIDs)
	b.AddHandler(h.HandleListShopVariants)
	b.AddHandler(h.HandleListShopVariantsByIDs)
	b.AddHandler(h.HandleListShopVariantsWithProductByIDs)
	return QueryBus{b}
}
