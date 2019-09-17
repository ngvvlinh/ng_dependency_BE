package catalog

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {

	//-- shop_product --//

	CreateShopProduct(context.Context, *CreateShopProductArgs) (*ShopProduct, error)

	UpdateShopProductInfo(context.Context, *UpdateShopProductInfoArgs) (*ShopProduct, error)

	UpdateShopProductStatus(context.Context, *UpdateStatusArgs) (int, error)

	UpdateShopProductImages(context.Context, *UpdateImagesArgs) (*ShopProduct, error)

	DeleteShopProducts(context.Context, *shopping.IDsQueryShopArgs) (int, error)

	//-- shop_variant --//

	CreateShopVariant(context.Context, *CreateShopVariantArgs) (*ShopVariant, error)

	UpdateShopVariantInfo(context.Context, *UpdateShopVariantInfoArgs) (*ShopVariant, error)

	DeleteShopVariants(context.Context, *shopping.IDsQueryShopArgs) (int, error)

	UpdateShopVariantStatus(context.Context, *UpdateStatusArgs) (int, error)

	UpdateShopVariantImages(context.Context, *UpdateImagesArgs) (*ShopVariant, error)

	UpdateShopVariantAttributes(context.Context, *UpdateShopVariantAttributes) (*ShopVariant, error)

	//-- category --//

	//-- collection --//

	//-- tag --//
}

type QueryService interface {
	GetShopProductByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProduct, error)
	GetShopProductWithVariantsByID(context.Context, *GetShopProductByIDQueryArgs) (*ShopProductWithVariants, error)
	ListShopProducts(context.Context, *shopping.ListQueryShopArgs) (*ShopProductsResponse, error)
	ListShopProductsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*ShopProductsResponse, error)
	ListShopProductsWithVariants(context.Context, *shopping.ListQueryShopArgs) (*ShopProductsWithVariantsResponse, error)
	ListShopProductsWithVariantsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*ShopProductsWithVariantsResponse, error)

	GetShopVariantByID(context.Context, *GetShopVariantByIDQueryArgs) (*ShopVariant, error)
	GetShopVariantWithProductByID(context.Context, *GetShopVariantByIDQueryArgs) (*ShopVariantWithProduct, error)
	ListShopVariants(context.Context, *shopping.ListQueryShopArgs) (*ShopVariantsResponse, error)
	ListShopVariantsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*ShopVariantsResponse, error)
	ListShopVariantsWithProductByIDs(context.Context, *shopping.IDsQueryShopArgs) (*ShopVariantsWithProductResponse, error)
}

//-- query --//

type IDsArgs struct {
	IDs []int64
}

type GetShopProductByIDQueryArgs struct {
	ProductID int64
	ShopID    int64
}

type GetShopVariantByIDQueryArgs struct {
	VariantID int64
	ShopID    int64
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
	ProductType ProductType
}

type UpdateShopProductInfoArgs struct {
	ShopID    int64
	ProductID int64

	Code        NullString
	Name        NullString
	Unit        NullString
	Note        NullString
	ShortDesc   NullString
	Description NullString
	DescHTML    NullString
	CostPrice   NullInt32
	ListPrice   NullInt32
	RetailPrice NullInt32
	ProductType ProductType
}

type CreateShopVariantArgs struct {
	ShopID    int64
	ProductID int64

	Code       string
	Name       string
	ImageURLs  []string
	Note       string
	Attributes Attributes
	DescriptionInfo
	PriceInfo
}

type UpdateShopVariantInfoArgs struct {
	ShopID    int64
	VariantID int64

	Code         NullString
	Name         NullString
	Note         NullString
	ShortDesc    NullString
	Descripttion NullString
	DescHTML     NullString
	CostPrice    NullInt32
	ListPrice    NullInt32
	RetailPrice  NullInt32
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

type UpdateShopVariantAttributes struct {
	ShopID     int64
	VariantID  int64
	Attributes Attributes
}
