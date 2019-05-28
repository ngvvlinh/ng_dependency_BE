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
	GetProducts(context.Context, *GetProductsQueryArgs) (*ProductsResonse, error)
	GetProductsWithVariants(context.Context, *GetProductsQueryArgs) (*ProductsWithVariantsResponse, error)

	GetVariantByID(context.Context, *GetVariantByIDQueryArgs) (*VariantWithProduct, error)
	GetVariants(context.Context, *GetVariantsQueryArgs) (*VariantsResponse, error)
	GetVariantsWithProduct(context.Context, *GetVariantsQueryArgs) (*VariantsWithProductResponse, error)

	GetShopProductByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProductExtended, error)
	GetShopProducts(context.Context, *GetShopProductsQueryArgs) (*ShopProductsResponse, error)
	GetShopProductsWithVariants(context.Context, *GetShopProductsQueryArgs) (*ShopProductsWithVariantsResponse, error)

	GetShopVariantByID(context.Context, *GetShopVariantsByIDQueryArgs) (*ShopVariantExtended, error)
	GetShopVariants(context.Context, *GetShopVariantsQueryArgs) (*ShopVariantsResponse, error)
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
	IDs             int64
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
}

type VariantsResponse struct {
	Variants []*Variant
}

type ProductsWithVariantsResponse struct {
	Products []*ProductWithVariants
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
