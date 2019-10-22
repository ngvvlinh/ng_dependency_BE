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

	CreateShopProduct(context.Context, *CreateShopProductArgs) (*ShopProductWithVariants, error)

	UpdateShopProductInfo(context.Context, *UpdateShopProductInfoArgs) (*ShopProductWithVariants, error)

	UpdateShopProductStatus(context.Context, *UpdateStatusArgs) (int, error)

	UpdateShopProductImages(context.Context, *UpdateImagesArgs) (*ShopProductWithVariants, error)

	UpdateShopProductMetaFields(context.Context, *UpdateShopProductMetaFieldsArgs) (*ShopProductWithVariants, error)

	DeleteShopProducts(context.Context, *shopping.IDsQueryShopArgs) (int, error)

	UpdateShopProductCategory(context.Context, *UpdateShopProductCategoryArgs) (*ShopProductWithVariants, error)

	RemoveShopProductCategory(context.Context, *RemoveShopProductCategoryArgs) (*ShopProductWithVariants, error)

	AddShopProductCollection(context.Context, *AddShopProductCollectionArgs) (int, error)

	RemoveShopProductCollection(context.Context, *RemoveShopProductColelctionArgs) (int, error)

	//-- shop_variant --//

	CreateShopVariant(context.Context, *CreateShopVariantArgs) (*ShopVariant, error)

	UpdateShopVariantInfo(context.Context, *UpdateShopVariantInfoArgs) (*ShopVariant, error)

	DeleteShopVariants(context.Context, *shopping.IDsQueryShopArgs) (int, error)

	UpdateShopVariantStatus(context.Context, *UpdateStatusArgs) (int, error)

	UpdateShopVariantImages(context.Context, *UpdateImagesArgs) (*ShopVariant, error)

	UpdateShopVariantAttributes(context.Context, *UpdateShopVariantAttributes) (*ShopVariant, error)

	//-- category --//

	CreateShopCategory(context.Context, *CreateShopCategoryArgs) (*ShopCategory, error)

	UpdateShopCategory(context.Context, *UpdateShopCategoryArgs) (*ShopCategory, error)

	DeleteShopCategory(context.Context, *DeleteShopCategoryArgs) (int, error)

	//-- collection --//

	CreateShopCollection(context.Context, *CreateShopCollectionArgs) (*ShopCollection, error)

	UpdateShopCollection(context.Context, *UpdateShopCollectionArgs) (*ShopCollection, error)

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
	ValidateVariantIDs(ctx context.Context, shopId int64, shopVariantIds []int64) error

	//--Category --//
	GetShopCategory(context.Context, *GetShopCategoryArgs) (*ShopCategory, error)
	ListShopCategories(context.Context, *shopping.ListQueryShopArgs) (*ShopCategoriesResponse, error)

	//-- Collection--//
	GetShopCollection(context.Context, *GetShopCollectionArgs) (*ShopCollection, error)
	ListShopCollections(context.Context, *shopping.ListQueryShopArgs) (*ShopCollectionsResponse, error)
	ListShopCollectionsByProductID(context.Context, *ListShopCollectionsByProductIDArgs) ([]*ShopCollection, error)
}

//-- query --//

type IDsArgs struct {
	IDs []int64
}

type GetShopCategoryArgs struct {
	ID     int64
	ShopID int64
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

type ShopCategoriesResponse struct {
	Categories []*ShopCategory
	Count      int32
	Paging     meta.PageInfo
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

type ListShopCollectionsByProductIDArgs struct {
	ProductID int64
	ShopID    int64
}

type GetShopCollectionArgs struct {
	ID     int64
	ShopID int64
}

type ShopCollectionsResponse struct {
	Collections []*ShopCollection
	Count       int32
	Paging      meta.PageInfo
}

//-- command --//

type CreateShopProductArgs struct {
	ShopID int64

	VendorID  int64
	Code      string
	Name      string
	Unit      string
	ImageURLs []string
	Note      string
	DescriptionInfo
	PriceInfo
	ProductType ProductType
	MetaFields  []*MetaField
}

type CreateShopCategoryArgs struct {
	ID       int64
	ShopID   int64
	ParentID int64
	Name     string
	Status   int
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
	CategoryID  int64
	VendorID    int64
}

type UpdateShopProductCategoryArgs struct {
	CategoryID int64
	ShopID     int64
	ProductID  int64
}

type UpdateShopCategoryArgs struct {
	ID       int64
	Name     NullString
	ShopID   int64
	ParentID int64
}

type UpdateShopCollectionArgs struct {
	ID     int64
	ShopID int64

	Name        NullString
	Description NullString
	DescHTML    NullString
	ShortDesc   NullString
}

type CreateShopCollectionArgs struct {
	ID     int64
	ShopID int64

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
}

type DeleteShopCategoryArgs struct {
	ID     int64
	ShopID int64
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

type UpdateShopProductMetaFieldsArgs struct {
	ID         int64
	ShopID     int64
	MetaFields []*MetaField
}

type UpdateShopVariantAttributes struct {
	ShopID     int64
	VariantID  int64
	Attributes Attributes
}

type AddShopProductCollectionArgs struct {
	ProductID     int64
	ShopID        int64
	CollectionIDs []int64
}

type RemoveShopProductColelctionArgs struct {
	ProductID     int64
	ShopID        int64
	CollectionIDs []int64
}

type ValidVendorIDEvent struct {
	VendorID int64
	ShopID   int64
}

type RemoveShopProductCategoryArgs struct {
	ProductID int64
	ShopID    int64
}
