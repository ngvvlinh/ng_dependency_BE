package catalog

import (
	"context"

	"etop.vn/api/meta"
)

type Bus struct{ meta.Bus }

type Aggregate interface {
	UpdateProduct(context.Context, *UpdateProductArgs) (*Product, error)
}

type QueryService interface {
	GetProductByID(context.Context, *GetProductByIDQueryArgs) (*ProductWithVariants, error)
	ListProducts(context.Context, *GetProductsQueryArgs) (*ProductsResonse, error)
	ListProductsWithVariants(context.Context, *GetProductsQueryArgs) (*ProductsWithVariantsResponse, error)

	GetVariantByID(context.Context, *GetVariantByIDQueryArgs) (*VariantWithProduct, error)
	ListVariants(context.Context, *GetVariantsQueryArgs) (*VariantsResponse, error)
	ListVariantsWithProduct(context.Context, *GetVariantsQueryArgs) (*VariantsWithProductResponse, error)

	GetShopProductByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProductExtended, error)
	ListShopProducts(context.Context, *GetShopProductsQueryArgs) (*ShopProductsResponse, error)
	ListShopProductsWithVariants(context.Context, *GetShopProductsQueryArgs) (*ShopProductsWithVariantsResponse, error)

	GetShopVariantByID(context.Context, *GetShopVariantsByIDQueryArgs) (*ShopVariantExtended, error)
	ListShopVariants(context.Context, *GetShopVariantsQueryArgs) (*ShopVariantsResponse, error)
}

//-- query --//

type GetProductByIDQueryArgs struct {
	ProductID int64
}

type GetVariantByIDQueryArgs struct {
	VariantID int64
}

type GetShopProductByIDQueryArgs struct {
	ProductID int64
	ShopID    int64
}

type GetShopVariantsByIDQueryArgs struct {
}

type GetProductsQueryArgs struct {
	IDs             []int64
	ProductSourceID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type GetVariantsQueryArgs struct {
	IDs             int64
	ProductSourceID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type GetShopProductsQueryArgs struct {
	IDs    int64
	ShopID int64
}

type GetShopVariantsQueryArgs struct {
	IDs    int64
	ShopID int64
}

type ProductsResonse struct {
	Products []*Product
	Count    int32
}

type VariantsResponse struct {
	Variants []*Variant
}

type ProductsWithVariantsResponse struct {
	Products []*ProductWithVariants
	Count    int32
}

type VariantsWithProductResponse struct {
	Variants []*VariantWithProduct
}

type ShopProductsResponse struct {
	Products []*ShopProductExtended
}

type ShopProductsWithVariantsResponse struct {
	Products []*ShopProductWithVariants
}

type ShopVariantsResponse struct {
	Variants []*ShopVariantExtended
}

//-- command --//

type UpdateProductArgs struct {
}
