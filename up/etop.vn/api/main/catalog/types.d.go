// Code generated by gen-cmd-query. DO NOT EDIT.

package catalog

import (
	context "context"

	meta "etop.vn/api/meta"
	metav1 "etop.vn/api/meta/v1"
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

type UpdateProductCommand struct {
	Result *Product `json:"-"`
}

type GetProductByIDQuery struct {
	ProductID int64

	Result *Product `json:"-"`
}

type GetProductWithVariantsByIDQuery struct {
	ProductID int64

	Result *ProductWithVariants `json:"-"`
}

type GetShopProductByIDQuery struct {
	ProductID       int64
	ProductSourceID int64

	Result *ShopProductExtended `json:"-"`
}

type GetShopProductWithVariantsByIDQuery struct {
	ProductID       int64
	ProductSourceID int64

	Result *ShopProductWithVariants `json:"-"`
}

type GetShopVariantByIDQuery struct {
	VariantID       int64
	ProductSourceID int64

	Result *ShopVariantExtended `json:"-"`
}

type GetShopVariantWithProductByIDQuery struct {
	VariantID       int64
	ProductSourceID int64

	Result *ShopVariantWithProduct `json:"-"`
}

type GetVariantByIDQuery struct {
	VariantID int64

	Result *Variant `json:"-"`
}

type GetVariantWithProductByIDQuery struct {
	VariantID int64

	Result *VariantWithProduct `json:"-"`
}

type ListProductsQuery struct {
	ProductSourceID int64
	Paging          metav1.Paging
	Filters         meta.Filters

	Result *ProductsResonse `json:"-"`
}

type ListProductsByIDsQuery struct {
	IDs []int64

	Result *ProductsResonse `json:"-"`
}

type ListProductsWithVariantsQuery struct {
	ProductSourceID int64
	Paging          metav1.Paging
	Filters         meta.Filters

	Result *ProductsWithVariantsResponse `json:"-"`
}

type ListProductsWithVariantsByIDsQuery struct {
	IDs []int64

	Result *ProductsWithVariantsResponse `json:"-"`
}

type ListShopProductsQuery struct {
	ProductSourceID int64
	Paging          metav1.Paging
	Filters         meta.Filters

	Result *ShopProductsResponse `json:"-"`
}

type ListShopProductsByIDsQuery struct {
	IDs             []int64
	ProductSourceID int64

	Result *ShopProductsResponse `json:"-"`
}

type ListShopProductsWithVariantsQuery struct {
	ProductSourceID int64
	Paging          metav1.Paging
	Filters         meta.Filters

	Result *ShopProductsWithVariantsResponse `json:"-"`
}

type ListShopProductsWithVariantsByIDsQuery struct {
	IDs             []int64
	ProductSourceID int64

	Result *ShopProductsWithVariantsResponse `json:"-"`
}

type ListShopVariantsQuery struct {
	ShopID  int64
	Paging  metav1.Paging
	Filters meta.Filters

	Result *ShopVariantsResponse `json:"-"`
}

type ListShopVariantsByIDsQuery struct {
	IDs             []int64
	ProductSourceID int64

	Result *ShopVariantsResponse `json:"-"`
}

type ListShopVariantsWithProductByIDsQuery struct {
	IDs             []int64
	ProductSourceID int64

	Result *ShopVariantsWithProductResponse `json:"-"`
}

type ListVariantsQuery struct {
	ProductSourceID int64
	Paging          metav1.Paging
	Filters         meta.Filters

	Result *VariantsResponse `json:"-"`
}

type ListVariantsByIDsQuery struct {
	IDs []int64

	Result *VariantsResponse `json:"-"`
}

type ListVariantsWithProductQuery struct {
	ProductSourceID int64
	Paging          metav1.Paging
	Filters         meta.Filters

	Result *VariantsWithProductResponse `json:"-"`
}

type ListVariantsWithProductByIDsQuery struct {
	IDs []int64

	Result *VariantsWithProductResponse `json:"-"`
}

// implement interfaces

func (q *UpdateProductCommand) command()                 {}
func (q *GetProductByIDQuery) query()                    {}
func (q *GetProductWithVariantsByIDQuery) query()        {}
func (q *GetShopProductByIDQuery) query()                {}
func (q *GetShopProductWithVariantsByIDQuery) query()    {}
func (q *GetShopVariantByIDQuery) query()                {}
func (q *GetShopVariantWithProductByIDQuery) query()     {}
func (q *GetVariantByIDQuery) query()                    {}
func (q *GetVariantWithProductByIDQuery) query()         {}
func (q *ListProductsQuery) query()                      {}
func (q *ListProductsByIDsQuery) query()                 {}
func (q *ListProductsWithVariantsQuery) query()          {}
func (q *ListProductsWithVariantsByIDsQuery) query()     {}
func (q *ListShopProductsQuery) query()                  {}
func (q *ListShopProductsByIDsQuery) query()             {}
func (q *ListShopProductsWithVariantsQuery) query()      {}
func (q *ListShopProductsWithVariantsByIDsQuery) query() {}
func (q *ListShopVariantsQuery) query()                  {}
func (q *ListShopVariantsByIDsQuery) query()             {}
func (q *ListShopVariantsWithProductByIDsQuery) query()  {}
func (q *ListVariantsQuery) query()                      {}
func (q *ListVariantsByIDsQuery) query()                 {}
func (q *ListVariantsWithProductQuery) query()           {}
func (q *ListVariantsWithProductByIDsQuery) query()      {}

// implement conversion

func (q *UpdateProductCommand) GetArgs() *UpdateProductArgs {
	return &UpdateProductArgs{}
}
func (q *GetProductByIDQuery) GetArgs() *GetProductByIDQueryArgs {
	return &GetProductByIDQueryArgs{
		ProductID: q.ProductID,
	}
}
func (q *GetProductWithVariantsByIDQuery) GetArgs() *GetProductByIDQueryArgs {
	return &GetProductByIDQueryArgs{
		ProductID: q.ProductID,
	}
}
func (q *GetShopProductByIDQuery) GetArgs() *GetShopProductByIDQueryArgs {
	return &GetShopProductByIDQueryArgs{
		ProductID:       q.ProductID,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *GetShopProductWithVariantsByIDQuery) GetArgs() *GetShopProductByIDQueryArgs {
	return &GetShopProductByIDQueryArgs{
		ProductID:       q.ProductID,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *GetShopVariantByIDQuery) GetArgs() *GetShopVariantByIDQueryArgs {
	return &GetShopVariantByIDQueryArgs{
		VariantID:       q.VariantID,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *GetShopVariantWithProductByIDQuery) GetArgs() *GetShopVariantByIDQueryArgs {
	return &GetShopVariantByIDQueryArgs{
		VariantID:       q.VariantID,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *GetVariantByIDQuery) GetArgs() *GetVariantByIDQueryArgs {
	return &GetVariantByIDQueryArgs{
		VariantID: q.VariantID,
	}
}
func (q *GetVariantWithProductByIDQuery) GetArgs() *GetVariantByIDQueryArgs {
	return &GetVariantByIDQueryArgs{
		VariantID: q.VariantID,
	}
}
func (q *ListProductsQuery) GetArgs() *ListProductsQueryArgs {
	return &ListProductsQueryArgs{
		ProductSourceID: q.ProductSourceID,
		Paging:          q.Paging,
		Filters:         q.Filters,
	}
}
func (q *ListProductsByIDsQuery) GetArgs() *IDsArgs {
	return &IDsArgs{
		IDs: q.IDs,
	}
}
func (q *ListProductsWithVariantsQuery) GetArgs() *ListProductsQueryArgs {
	return &ListProductsQueryArgs{
		ProductSourceID: q.ProductSourceID,
		Paging:          q.Paging,
		Filters:         q.Filters,
	}
}
func (q *ListProductsWithVariantsByIDsQuery) GetArgs() *IDsArgs {
	return &IDsArgs{
		IDs: q.IDs,
	}
}
func (q *ListShopProductsQuery) GetArgs() *ListShopProductsQueryArgs {
	return &ListShopProductsQueryArgs{
		ProductSourceID: q.ProductSourceID,
		Paging:          q.Paging,
		Filters:         q.Filters,
	}
}
func (q *ListShopProductsByIDsQuery) GetArgs() *IDsShopArgs {
	return &IDsShopArgs{
		IDs:             q.IDs,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *ListShopProductsWithVariantsQuery) GetArgs() *ListShopProductsQueryArgs {
	return &ListShopProductsQueryArgs{
		ProductSourceID: q.ProductSourceID,
		Paging:          q.Paging,
		Filters:         q.Filters,
	}
}
func (q *ListShopProductsWithVariantsByIDsQuery) GetArgs() *IDsShopArgs {
	return &IDsShopArgs{
		IDs:             q.IDs,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *ListShopVariantsQuery) GetArgs() *ListShopVariantsQueryArgs {
	return &ListShopVariantsQueryArgs{
		ShopID:  q.ShopID,
		Paging:  q.Paging,
		Filters: q.Filters,
	}
}
func (q *ListShopVariantsByIDsQuery) GetArgs() *IDsShopArgs {
	return &IDsShopArgs{
		IDs:             q.IDs,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *ListShopVariantsWithProductByIDsQuery) GetArgs() *IDsShopArgs {
	return &IDsShopArgs{
		IDs:             q.IDs,
		ProductSourceID: q.ProductSourceID,
	}
}
func (q *ListVariantsQuery) GetArgs() *ListVariantsQueryArgs {
	return &ListVariantsQueryArgs{
		ProductSourceID: q.ProductSourceID,
		Paging:          q.Paging,
		Filters:         q.Filters,
	}
}
func (q *ListVariantsByIDsQuery) GetArgs() *IDsArgs {
	return &IDsArgs{
		IDs: q.IDs,
	}
}
func (q *ListVariantsWithProductQuery) GetArgs() *ListVariantsQueryArgs {
	return &ListVariantsQueryArgs{
		ProductSourceID: q.ProductSourceID,
		Paging:          q.Paging,
		Filters:         q.Filters,
	}
}
func (q *ListVariantsWithProductByIDsQuery) GetArgs() *IDsArgs {
	return &IDsArgs{
		IDs: q.IDs,
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
	b.AddHandler(h.HandleUpdateProduct)
	return CommandBus{b}
}

func (h AggregateHandler) HandleUpdateProduct(ctx context.Context, cmd *UpdateProductCommand) error {
	result, err := h.inner.UpdateProduct(ctx, cmd.GetArgs())
	cmd.Result = result
	return err
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
	b.AddHandler(h.HandleGetProductByID)
	b.AddHandler(h.HandleGetProductWithVariantsByID)
	b.AddHandler(h.HandleGetShopProductByID)
	b.AddHandler(h.HandleGetShopProductWithVariantsByID)
	b.AddHandler(h.HandleGetShopVariantByID)
	b.AddHandler(h.HandleGetShopVariantWithProductByID)
	b.AddHandler(h.HandleGetVariantByID)
	b.AddHandler(h.HandleGetVariantWithProductByID)
	b.AddHandler(h.HandleListProducts)
	b.AddHandler(h.HandleListProductsByIDs)
	b.AddHandler(h.HandleListProductsWithVariants)
	b.AddHandler(h.HandleListProductsWithVariantsByIDs)
	b.AddHandler(h.HandleListShopProducts)
	b.AddHandler(h.HandleListShopProductsByIDs)
	b.AddHandler(h.HandleListShopProductsWithVariants)
	b.AddHandler(h.HandleListShopProductsWithVariantsByIDs)
	b.AddHandler(h.HandleListShopVariants)
	b.AddHandler(h.HandleListShopVariantsByIDs)
	b.AddHandler(h.HandleListShopVariantsWithProductByIDs)
	b.AddHandler(h.HandleListVariants)
	b.AddHandler(h.HandleListVariantsByIDs)
	b.AddHandler(h.HandleListVariantsWithProduct)
	b.AddHandler(h.HandleListVariantsWithProductByIDs)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleGetProductByID(ctx context.Context, query *GetProductByIDQuery) error {
	result, err := h.inner.GetProductByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetProductWithVariantsByID(ctx context.Context, query *GetProductWithVariantsByIDQuery) error {
	result, err := h.inner.GetProductWithVariantsByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetShopProductByID(ctx context.Context, query *GetShopProductByIDQuery) error {
	result, err := h.inner.GetShopProductByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetShopProductWithVariantsByID(ctx context.Context, query *GetShopProductWithVariantsByIDQuery) error {
	result, err := h.inner.GetShopProductWithVariantsByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetShopVariantByID(ctx context.Context, query *GetShopVariantByIDQuery) error {
	result, err := h.inner.GetShopVariantByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetShopVariantWithProductByID(ctx context.Context, query *GetShopVariantWithProductByIDQuery) error {
	result, err := h.inner.GetShopVariantWithProductByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetVariantByID(ctx context.Context, query *GetVariantByIDQuery) error {
	result, err := h.inner.GetVariantByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetVariantWithProductByID(ctx context.Context, query *GetVariantWithProductByIDQuery) error {
	result, err := h.inner.GetVariantWithProductByID(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListProducts(ctx context.Context, query *ListProductsQuery) error {
	result, err := h.inner.ListProducts(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListProductsByIDs(ctx context.Context, query *ListProductsByIDsQuery) error {
	result, err := h.inner.ListProductsByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListProductsWithVariants(ctx context.Context, query *ListProductsWithVariantsQuery) error {
	result, err := h.inner.ListProductsWithVariants(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListProductsWithVariantsByIDs(ctx context.Context, query *ListProductsWithVariantsByIDsQuery) error {
	result, err := h.inner.ListProductsWithVariantsByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopProducts(ctx context.Context, query *ListShopProductsQuery) error {
	result, err := h.inner.ListShopProducts(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopProductsByIDs(ctx context.Context, query *ListShopProductsByIDsQuery) error {
	result, err := h.inner.ListShopProductsByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopProductsWithVariants(ctx context.Context, query *ListShopProductsWithVariantsQuery) error {
	result, err := h.inner.ListShopProductsWithVariants(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopProductsWithVariantsByIDs(ctx context.Context, query *ListShopProductsWithVariantsByIDsQuery) error {
	result, err := h.inner.ListShopProductsWithVariantsByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopVariants(ctx context.Context, query *ListShopVariantsQuery) error {
	result, err := h.inner.ListShopVariants(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopVariantsByIDs(ctx context.Context, query *ListShopVariantsByIDsQuery) error {
	result, err := h.inner.ListShopVariantsByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListShopVariantsWithProductByIDs(ctx context.Context, query *ListShopVariantsWithProductByIDsQuery) error {
	result, err := h.inner.ListShopVariantsWithProductByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListVariants(ctx context.Context, query *ListVariantsQuery) error {
	result, err := h.inner.ListVariants(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListVariantsByIDs(ctx context.Context, query *ListVariantsByIDsQuery) error {
	result, err := h.inner.ListVariantsByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListVariantsWithProduct(ctx context.Context, query *ListVariantsWithProductQuery) error {
	result, err := h.inner.ListVariantsWithProduct(ctx, query.GetArgs())
	query.Result = result
	return err
}

func (h QueryServiceHandler) HandleListVariantsWithProductByIDs(ctx context.Context, query *ListVariantsWithProductByIDsQuery) error {
	result, err := h.inner.ListVariantsWithProductByIDs(ctx, query.GetArgs())
	query.Result = result
	return err
}
