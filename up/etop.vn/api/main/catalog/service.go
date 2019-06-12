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
	GetProductByID(context.Context, *GetProductByIDQueryArgs) (*Product, error)
	GetProductWithVariantsByID(context.Context, *GetProductByIDQueryArgs) (*ProductWithVariants, error)
	ListProducts(context.Context, *ListProductsQueryArgs) (*ProductsResonse, error)
	ListProductsByIDs(context.Context, *IDsArgs) (*ProductsResonse, error)
	ListProductsWithVariants(context.Context, *ListProductsQueryArgs) (*ProductsWithVariantsResponse, error)
	ListProductsWithVariantsByIDs(context.Context, *IDsArgs) (*ProductsWithVariantsResponse, error)

	GetVariantByID(context.Context, *GetVariantByIDQueryArgs) (*Variant, error)
	GetVariantWithProductByID(context.Context, *GetVariantByIDQueryArgs) (*VariantWithProduct, error)
	ListVariants(context.Context, *ListVariantsQueryArgs) (*VariantsResponse, error)
	ListVariantsByIDs(context.Context, *IDsArgs) (*VariantsResponse, error)
	ListVariantsWithProduct(context.Context, *ListVariantsQueryArgs) (*VariantsWithProductResponse, error)
	ListVariantsWithProductByIDs(context.Context, *IDsArgs) (*VariantsWithProductResponse, error)

	GetShopProductByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProductExtended, error)
	GetShopProductWithVariantsByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProductWithVariants, error)
	ListShopProducts(context.Context, *ListShopProductsQueryArgs) (*ShopProductsResponse, error)
	ListShopProductsByIDs(context.Context, *IDsShopArgs) (*ShopProductsResponse, error)
	ListShopProductsWithVariants(context.Context, *ListShopProductsQueryArgs) (*ShopProductsWithVariantsResponse, error)
	ListShopProductsWithVariantsByIDs(context.Context, *IDsShopArgs) (*ShopProductsWithVariantsResponse, error)

	GetShopVariantByID(context.Context, *GetShopVariantByIDQueryArgs) (*ShopVariantExtended, error)
	GetShopVariantWithProductByID(context.Context, *GetShopVariantByIDQueryArgs) (*ShopVariantWithProduct, error)
	ListShopVariants(context.Context, *ListShopVariantsQueryArgs) (*ShopVariantsResponse, error)
	ListShopVariantsByIDs(context.Context, *IDsShopArgs) (*ShopVariantsResponse, error)
	ListShopVariantsWithProductByIDs(context.Context, *IDsShopArgs) (*ShopVariantsWithProductResponse, error)
}

//-- query --//

type IDsArgs struct {
	IDs []int64
}

type IDsShopArgs struct {
	IDs    []int64
	ShopID int64
}

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

type GetShopVariantByIDQueryArgs struct {
	VariantID int64
	ShopID    int64
}

type ListProductsQueryArgs struct {
	ProductSourceID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type ListVariantsQueryArgs struct {
	ProductSourceID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type ListShopProductsQueryArgs struct {
	ShopID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type ListShopVariantsQueryArgs struct {
	ShopID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type ProductsResonse struct {
	Products []*Product
	Count    int32
	Paging   meta.PageInfo
}

type VariantsResponse struct {
	Variants []*Variant
	Count    int32
	Paging   meta.PageInfo
}

type ProductsWithVariantsResponse struct {
	Products []*ProductWithVariants
	Count    int32
	Paging   meta.PageInfo
}

type VariantsWithProductResponse struct {
	Variants []*VariantWithProduct
	Count    int32
	Paging   meta.PageInfo
}

type ShopProductsResponse struct {
	Products []*ShopProductExtended
	Count    int32
	Paging   meta.PageInfo
}

type ShopProductsWithVariantsResponse struct {
	Products []*ShopProductWithVariants
	Count    int32
	Paging   meta.PageInfo
}

type ShopVariantsResponse struct {
	Variants []*ShopVariantExtended
	Count    int32
	Paging   meta.PageInfo
}

type ShopVariantsWithProductResponse struct {
	Variants []*ShopVariantWithProduct
	Count    int32
	Paging   meta.PageInfo
}

//-- command --//

type UpdateProductArgs struct {
}
