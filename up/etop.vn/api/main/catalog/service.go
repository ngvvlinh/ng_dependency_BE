package catalog

import (
	"context"

	"etop.vn/api/meta"
)

type Aggregate interface {

	//-- shop_product --//

	CreateShopProduct(context.Context, *CreateShopProductArgs) (*ShopProduct, error)

	UpdateShopProductInfo(context.Context, *UpdateShopProductInfoArgs) (*ShopProduct, error)

	UpdateShopProductStatus(context.Context, *UpdateStatusArgs) (*ShopProduct, error)

	UpdateShopProductImages(context.Context, *UpdateImagesArgs) (*ShopProduct, error)

	DeleteShopProducts(context.Context, *IDsShopArgs) (*meta.Empty, error)

	//-- shop_variant --//

	CreateShopVariant(context.Context, *CreateShopVariantArgs) (*ShopVariant, error)

	UpdateShopVariantInfo(context.Context, *UpdateShopVariantInfoArgs) (*ShopVariant, error)

	DeleteShopVariants(context.Context, *IDsShopArgs) (*meta.Empty, error)

	UpdateShopVariantStatus(context.Context, *UpdateStatusArgs) (*ShopVariant, error)

	UpdateShopVariantImages(context.Context, *UpdateImagesArgs) (*ShopVariant, error)

	//-- category --//

	//-- collection --//

	//-- tag --//
}

type QueryService interface {
	GetShopProductByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProduct, error)
	GetShopProductWithVariantsByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProductWithVariants, error)
	ListShopProducts(context.Context, *ListShopProductsQueryArgs) (*ShopProductsResponse, error)
	ListShopProductsByIDs(context.Context, *IDsShopArgs) (*ShopProductsResponse, error)
	ListShopProductsWithVariants(context.Context, *ListShopProductsQueryArgs) (*ShopProductsWithVariantsResponse, error)
	ListShopProductsWithVariantsByIDs(context.Context, *IDsShopArgs) (*ShopProductsWithVariantsResponse, error)

	GetShopVariantByID(context.Context, *GetShopVariantByIDQueryArgs) (*ShopVariant, error)
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
	ShopID int64

	Paging  meta.Paging
	Filters meta.Filters
}

type ListVariantsQueryArgs struct {
	ShopID int64

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

type ShopProductsResponse struct {
	Products []*ShopProduct
	Count    int32
	Paging   meta.PageInfo
}

type ShopProductsWithVariantsResponse struct {
	Products []*ShopProductWithVariants
	Count    int32
	Paging   meta.PageInfo
}

type ShopVariantsResponse struct {
	Variants []*ShopVariant
	Count    int32
	Paging   meta.PageInfo
}

type ShopVariantsWithProductResponse struct {
	Variants []*ShopVariantWithProduct
	Count    int32
	Paging   meta.PageInfo
}

//-- command --//

type CreateShopProductArgs struct {
	ShopID int64

	Code      string
	Name      string
	Unit      string
	ImageURLs []string
	Note      string
	DescriptionInfo
	PriceInfo
}

type UpdateShopProductInfoArgs struct {
	ShopID    int64
	ProductID int64

	Code *string
	Name *string
	Unit *string
	Note *string
	*DescriptionInfo
}

type CreateShopVariantArgs struct {
	ShopID    int64
	ProductID int64

	Code      string
	Name      string
	ImageURLs []string
	Note      string
	DescriptionInfo
	PriceInfo
}

type UpdateShopVariantInfoArgs struct {
	ShopID    int64
	VariantID int64

	Code *string
	Name *string
	Unit *string
	Note *string
	*DescriptionInfo
}

type UpdateStatusArgs struct {
	IDs    []int64
	ShopID int64
	Status int16
}

type UpdateImagesArgs struct {
	ID      int64
	ShopID  int64
	Updates []*meta.UpdateSet
}
