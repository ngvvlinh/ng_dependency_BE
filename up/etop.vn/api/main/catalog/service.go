package catalog

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
	dot "etop.vn/capi/dot"
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

	//-- brand --//

	CreateBrand(context.Context, *CreateBrandArgs) (*ShopBrand, error)

	UpdateBrandInfo(context.Context, *UpdateBrandArgs) (*ShopBrand, error)

	DeleteShopBrand(ctx context.Context, ids []dot.ID, shopId dot.ID) (int32, error)

	// -- variant_supplier -- //

	CreateVariantSupplier(context.Context, *CreateVariantSupplier) (*ShopVariantSupplier, error)

	CreateVariantsSupplier(context.Context, *CreateVariantsSupplier) (int, error)

	DeleteVariantSupplier(ctx context.Context, variantID dot.ID, supplierID dot.ID, shopID dot.ID) error

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
	ValidateVariantIDs(ctx context.Context, shopId dot.ID, shopVariantIds []dot.ID) error
	//--Category --//
	GetShopCategory(context.Context, *GetShopCategoryArgs) (*ShopCategory, error)
	ListShopCategories(context.Context, *shopping.ListQueryShopArgs) (*ShopCategoriesResponse, error)

	//-- Collection--//
	GetShopCollection(context.Context, *GetShopCollectionArgs) (*ShopCollection, error)
	ListShopCollections(context.Context, *shopping.ListQueryShopArgs) (*ShopCollectionsResponse, error)
	ListShopCollectionsByProductID(context.Context, *ListShopCollectionsByProductIDArgs) ([]*ShopCollection, error)

	//-- Brand --//
	GetBrandByID(ctx context.Context, id dot.ID, shopID dot.ID) (*ShopBrand, error)
	GetBrandsByIDs(ctx context.Context, ids []dot.ID, shopID dot.ID) ([]*ShopBrand, error)
	ListBrands(ctx context.Context, paging meta.Paging, shopId dot.ID) (*ListBrandsResult, error)

	// -- variant_supplier -- //

	GetSupplierIDsByVariantID(ctx context.Context, variantID dot.ID, shopID dot.ID) ([]dot.ID, error)

	GetVariantsBySupplierID(ctx context.Context, supplierID dot.ID, shopID dot.ID) (*ShopVariantsResponse, error)

	//-- query --//
}

type IDsArgs struct {
	IDs []dot.ID
}

type GetShopCategoryArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetShopProductByIDQueryArgs struct {
	ProductID dot.ID
	ShopID    dot.ID
}

type CreateVariantsSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantIDs []dot.ID
}

type GetShopVariantByIDQueryArgs struct {
	VariantID dot.ID
	ShopID    dot.ID
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
	ProductID dot.ID
	ShopID    dot.ID
}

type GetShopCollectionArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type ShopCollectionsResponse struct {
	Collections []*ShopCollection
	Count       int32
	Paging      meta.PageInfo
}

//-- command --//

type CreateShopProductArgs struct {
	ShopID dot.ID

	Code      string
	Name      string
	Unit      string
	ImageURLs []string
	Note      string
	DescriptionInfo
	PriceInfo
	ProductType ProductType
	MetaFields  []*MetaField
	BrandID     dot.ID
}

type CreateShopCategoryArgs struct {
	ID       dot.ID
	ShopID   dot.ID
	ParentID dot.ID
	Name     string
	Status   int
}

type UpdateShopProductInfoArgs struct {
	ShopID    dot.ID
	ProductID dot.ID

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
	BrandID     NullID
	ProductType ProductType
	CategoryID  dot.ID
}

type UpdateShopProductCategoryArgs struct {
	CategoryID dot.ID
	ShopID     dot.ID
	ProductID  dot.ID
}

type UpdateShopCategoryArgs struct {
	ID       dot.ID
	Name     NullString
	ShopID   dot.ID
	ParentID dot.ID
}

type UpdateShopCollectionArgs struct {
	ID     dot.ID
	ShopID dot.ID

	Name        NullString
	Description NullString
	DescHTML    NullString
	ShortDesc   NullString
}

type CreateShopCollectionArgs struct {
	ID     dot.ID
	ShopID dot.ID

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string
}

type DeleteShopCategoryArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type CreateShopVariantArgs struct {
	ShopID    dot.ID
	ProductID dot.ID

	Code       string
	Name       string
	ImageURLs  []string
	Note       string
	Attributes Attributes
	DescriptionInfo
	PriceInfo
}

type UpdateShopVariantInfoArgs struct {
	ShopID    dot.ID
	VariantID dot.ID

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
	IDs    []dot.ID
	ShopID dot.ID
	Status int16
}

type UpdateImagesArgs struct {
	ID      dot.ID
	ShopID  dot.ID
	Updates []*meta.UpdateSet
}

type UpdateShopProductMetaFieldsArgs struct {
	ID         dot.ID
	ShopID     dot.ID
	MetaFields []*MetaField
}

type UpdateShopVariantAttributes struct {
	ShopID     dot.ID
	VariantID  dot.ID
	Attributes Attributes
}

type AddShopProductCollectionArgs struct {
	ProductID     dot.ID
	ShopID        dot.ID
	CollectionIDs []dot.ID
}

type RemoveShopProductColelctionArgs struct {
	ProductID     dot.ID
	ShopID        dot.ID
	CollectionIDs []dot.ID
}

type ValidSupplierIDEvent struct {
	SupplierID dot.ID
	ShopID     dot.ID
}

type RemoveShopProductCategoryArgs struct {
	ProductID dot.ID
	ShopID    dot.ID
}

// +convert:create=ShopBrand
type CreateBrandArgs struct {
	ShopID      dot.ID
	BrandName   string
	Description string
}

// +convert:update=ShopBrand
type UpdateBrandArgs struct {
	ShopID      dot.ID
	ID          dot.ID
	BrandName   string
	Description string
}

type ShopBrand struct {
	ID     dot.ID
	ShopID dot.ID

	BrandName   string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type ListBrandsResult struct {
	ShopBrands []*ShopBrand
	PageInfo   meta.PageInfo
	Total      int32
}

type GetVariantsBySupplierIDResponse struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantIDs []dot.ID
}

type GetSuppliersByVariantIDResponse struct {
	ShopID      dot.ID
	VariantID   dot.ID
	SupplierIDs []dot.ID
}

// +convert:create=ShopVariantSupplier
type CreateVariantSupplier struct {
	ShopID     dot.ID
	SupplierID dot.ID
	VariantID  dot.ID
}
