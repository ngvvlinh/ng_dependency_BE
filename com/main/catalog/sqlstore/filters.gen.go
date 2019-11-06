// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
)

type ShopVariantFilters struct{ prefix string }

func NewShopVariantFilters(prefix string) ShopVariantFilters {
	return ShopVariantFilters{prefix}
}

func (ft *ShopVariantFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopVariantFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopVariantFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopVariantFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopVariantFilters) ByVariantID(VariantID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == 0,
	}
}

func (ft *ShopVariantFilters) ByVariantIDPtr(VariantID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == nil,
		IsZero: VariantID != nil && (*VariantID) == 0,
	}
}

func (ft *ShopVariantFilters) ByProductID(ProductID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == 0,
	}
}

func (ft *ShopVariantFilters) ByProductIDPtr(ProductID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == nil,
		IsZero: ProductID != nil && (*ProductID) == 0,
	}
}

func (ft *ShopVariantFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *ShopVariantFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *ShopVariantFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShopVariantFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShopVariantFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *ShopVariantFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *ShopVariantFilters) ByDescHTML(DescHTML string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == "",
	}
}

func (ft *ShopVariantFilters) ByDescHTMLPtr(DescHTML *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == nil,
		IsZero: DescHTML != nil && (*DescHTML) == "",
	}
}

func (ft *ShopVariantFilters) ByShortDesc(ShortDesc string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "short_desc",
		Value:  ShortDesc,
		IsNil:  ShortDesc == "",
	}
}

func (ft *ShopVariantFilters) ByShortDescPtr(ShortDesc *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "short_desc",
		Value:  ShortDesc,
		IsNil:  ShortDesc == nil,
		IsZero: ShortDesc != nil && (*ShortDesc) == "",
	}
}

func (ft *ShopVariantFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *ShopVariantFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *ShopVariantFilters) ByCostPrice(CostPrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cost_price",
		Value:  CostPrice,
		IsNil:  CostPrice == 0,
	}
}

func (ft *ShopVariantFilters) ByCostPricePtr(CostPrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cost_price",
		Value:  CostPrice,
		IsNil:  CostPrice == nil,
		IsZero: CostPrice != nil && (*CostPrice) == 0,
	}
}

func (ft *ShopVariantFilters) ByListPrice(ListPrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "list_price",
		Value:  ListPrice,
		IsNil:  ListPrice == 0,
	}
}

func (ft *ShopVariantFilters) ByListPricePtr(ListPrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "list_price",
		Value:  ListPrice,
		IsNil:  ListPrice == nil,
		IsZero: ListPrice != nil && (*ListPrice) == 0,
	}
}

func (ft *ShopVariantFilters) ByRetailPrice(RetailPrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "retail_price",
		Value:  RetailPrice,
		IsNil:  RetailPrice == 0,
	}
}

func (ft *ShopVariantFilters) ByRetailPricePtr(RetailPrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "retail_price",
		Value:  RetailPrice,
		IsNil:  RetailPrice == nil,
		IsZero: RetailPrice != nil && (*RetailPrice) == 0,
	}
}

func (ft *ShopVariantFilters) ByStatus(Status model.Status3) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ShopVariantFilters) ByStatusPtr(Status *model.Status3) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ShopVariantFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopVariantFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopVariantFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopVariantFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopVariantFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopVariantFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ShopVariantFilters) ByNameNorm(NameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name_norm",
		Value:  NameNorm,
		IsNil:  NameNorm == "",
	}
}

func (ft *ShopVariantFilters) ByNameNormPtr(NameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name_norm",
		Value:  NameNorm,
		IsNil:  NameNorm == nil,
		IsZero: NameNorm != nil && (*NameNorm) == "",
	}
}

func (ft *ShopVariantFilters) ByAttrNormKv(AttrNormKv string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "attr_norm_kv",
		Value:  AttrNormKv,
		IsNil:  AttrNormKv == "",
	}
}

func (ft *ShopVariantFilters) ByAttrNormKvPtr(AttrNormKv *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "attr_norm_kv",
		Value:  AttrNormKv,
		IsNil:  AttrNormKv == nil,
		IsZero: AttrNormKv != nil && (*AttrNormKv) == "",
	}
}

type ShopProductFilters struct{ prefix string }

func NewShopProductFilters(prefix string) ShopProductFilters {
	return ShopProductFilters{prefix}
}

func (ft *ShopProductFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopProductFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopProductFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopProductFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopProductFilters) ByProductID(ProductID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == 0,
	}
}

func (ft *ShopProductFilters) ByProductIDPtr(ProductID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == nil,
		IsZero: ProductID != nil && (*ProductID) == 0,
	}
}

func (ft *ShopProductFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *ShopProductFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *ShopProductFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShopProductFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShopProductFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *ShopProductFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *ShopProductFilters) ByDescHTML(DescHTML string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == "",
	}
}

func (ft *ShopProductFilters) ByDescHTMLPtr(DescHTML *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == nil,
		IsZero: DescHTML != nil && (*DescHTML) == "",
	}
}

func (ft *ShopProductFilters) ByShortDesc(ShortDesc string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "short_desc",
		Value:  ShortDesc,
		IsNil:  ShortDesc == "",
	}
}

func (ft *ShopProductFilters) ByShortDescPtr(ShortDesc *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "short_desc",
		Value:  ShortDesc,
		IsNil:  ShortDesc == nil,
		IsZero: ShortDesc != nil && (*ShortDesc) == "",
	}
}

func (ft *ShopProductFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *ShopProductFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *ShopProductFilters) ByUnit(Unit string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "unit",
		Value:  Unit,
		IsNil:  Unit == "",
	}
}

func (ft *ShopProductFilters) ByUnitPtr(Unit *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "unit",
		Value:  Unit,
		IsNil:  Unit == nil,
		IsZero: Unit != nil && (*Unit) == "",
	}
}

func (ft *ShopProductFilters) ByCategoryID(CategoryID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "category_id",
		Value:  CategoryID,
		IsNil:  CategoryID == 0,
	}
}

func (ft *ShopProductFilters) ByCategoryIDPtr(CategoryID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "category_id",
		Value:  CategoryID,
		IsNil:  CategoryID == nil,
		IsZero: CategoryID != nil && (*CategoryID) == 0,
	}
}

func (ft *ShopProductFilters) ByCostPrice(CostPrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cost_price",
		Value:  CostPrice,
		IsNil:  CostPrice == 0,
	}
}

func (ft *ShopProductFilters) ByCostPricePtr(CostPrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cost_price",
		Value:  CostPrice,
		IsNil:  CostPrice == nil,
		IsZero: CostPrice != nil && (*CostPrice) == 0,
	}
}

func (ft *ShopProductFilters) ByListPrice(ListPrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "list_price",
		Value:  ListPrice,
		IsNil:  ListPrice == 0,
	}
}

func (ft *ShopProductFilters) ByListPricePtr(ListPrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "list_price",
		Value:  ListPrice,
		IsNil:  ListPrice == nil,
		IsZero: ListPrice != nil && (*ListPrice) == 0,
	}
}

func (ft *ShopProductFilters) ByRetailPrice(RetailPrice int32) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "retail_price",
		Value:  RetailPrice,
		IsNil:  RetailPrice == 0,
	}
}

func (ft *ShopProductFilters) ByRetailPricePtr(RetailPrice *int32) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "retail_price",
		Value:  RetailPrice,
		IsNil:  RetailPrice == nil,
		IsZero: RetailPrice != nil && (*RetailPrice) == 0,
	}
}

func (ft *ShopProductFilters) ByBrandID(BrandID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "brand_id",
		Value:  BrandID,
		IsNil:  BrandID == 0,
	}
}

func (ft *ShopProductFilters) ByBrandIDPtr(BrandID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "brand_id",
		Value:  BrandID,
		IsNil:  BrandID == nil,
		IsZero: BrandID != nil && (*BrandID) == 0,
	}
}

func (ft *ShopProductFilters) ByStatus(Status model.Status3) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ShopProductFilters) ByStatusPtr(Status *model.Status3) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ShopProductFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopProductFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopProductFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopProductFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopProductFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopProductFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ShopProductFilters) ByNameNorm(NameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name_norm",
		Value:  NameNorm,
		IsNil:  NameNorm == "",
	}
}

func (ft *ShopProductFilters) ByNameNormPtr(NameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name_norm",
		Value:  NameNorm,
		IsNil:  NameNorm == nil,
		IsZero: NameNorm != nil && (*NameNorm) == "",
	}
}

func (ft *ShopProductFilters) ByNameNormUa(NameNormUa string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name_norm_ua",
		Value:  NameNormUa,
		IsNil:  NameNormUa == "",
	}
}

func (ft *ShopProductFilters) ByNameNormUaPtr(NameNormUa *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name_norm_ua",
		Value:  NameNormUa,
		IsNil:  NameNormUa == nil,
		IsZero: NameNormUa != nil && (*NameNormUa) == "",
	}
}

func (ft *ShopProductFilters) ByProductType(ProductType string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_type",
		Value:  ProductType,
		IsNil:  ProductType == "",
	}
}

func (ft *ShopProductFilters) ByProductTypePtr(ProductType *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_type",
		Value:  ProductType,
		IsNil:  ProductType == nil,
		IsZero: ProductType != nil && (*ProductType) == "",
	}
}

type ProductShopCollectionFilters struct{ prefix string }

func NewProductShopCollectionFilters(prefix string) ProductShopCollectionFilters {
	return ProductShopCollectionFilters{prefix}
}

func (ft *ProductShopCollectionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ProductShopCollectionFilters) Prefix() string {
	return ft.prefix
}

func (ft *ProductShopCollectionFilters) ByCollectionID(CollectionID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "collection_id",
		Value:  CollectionID,
		IsNil:  CollectionID == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByCollectionIDPtr(CollectionID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "collection_id",
		Value:  CollectionID,
		IsNil:  CollectionID == nil,
		IsZero: CollectionID != nil && (*CollectionID) == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByProductID(ProductID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByProductIDPtr(ProductID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == nil,
		IsZero: ProductID != nil && (*ProductID) == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByStatus(Status int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByStatusPtr(Status *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ProductShopCollectionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ProductShopCollectionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ProductShopCollectionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ProductShopCollectionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type ShopCategoryFilters struct{ prefix string }

func NewShopCategoryFilters(prefix string) ShopCategoryFilters {
	return ShopCategoryFilters{prefix}
}

func (ft *ShopCategoryFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopCategoryFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopCategoryFilters) ByID(ID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopCategoryFilters) ByIDPtr(ID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopCategoryFilters) ByParentID(ParentID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "parent_id",
		Value:  ParentID,
		IsNil:  ParentID == 0,
	}
}

func (ft *ShopCategoryFilters) ByParentIDPtr(ParentID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "parent_id",
		Value:  ParentID,
		IsNil:  ParentID == nil,
		IsZero: ParentID != nil && (*ParentID) == 0,
	}
}

func (ft *ShopCategoryFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopCategoryFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopCategoryFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShopCategoryFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShopCategoryFilters) ByStatus(Status int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ShopCategoryFilters) ByStatusPtr(Status *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ShopCategoryFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopCategoryFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopCategoryFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopCategoryFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopCategoryFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopCategoryFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

type ShopCollectionFilters struct{ prefix string }

func NewShopCollectionFilters(prefix string) ShopCollectionFilters {
	return ShopCollectionFilters{prefix}
}

func (ft *ShopCollectionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopCollectionFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopCollectionFilters) ByID(ID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopCollectionFilters) ByIDPtr(ID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopCollectionFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopCollectionFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopCollectionFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShopCollectionFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShopCollectionFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *ShopCollectionFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *ShopCollectionFilters) ByDescHTML(DescHTML string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == "",
	}
}

func (ft *ShopCollectionFilters) ByDescHTMLPtr(DescHTML *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == nil,
		IsZero: DescHTML != nil && (*DescHTML) == "",
	}
}

func (ft *ShopCollectionFilters) ByShortDesc(ShortDesc string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "short_desc",
		Value:  ShortDesc,
		IsNil:  ShortDesc == "",
	}
}

func (ft *ShopCollectionFilters) ByShortDescPtr(ShortDesc *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "short_desc",
		Value:  ShortDesc,
		IsNil:  ShortDesc == nil,
		IsZero: ShortDesc != nil && (*ShortDesc) == "",
	}
}

func (ft *ShopCollectionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopCollectionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopCollectionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopCollectionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type ShopProductCollectionFilters struct{ prefix string }

func NewShopProductCollectionFilters(prefix string) ShopProductCollectionFilters {
	return ShopProductCollectionFilters{prefix}
}

func (ft *ShopProductCollectionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopProductCollectionFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopProductCollectionFilters) ByProductID(ProductID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == 0,
	}
}

func (ft *ShopProductCollectionFilters) ByProductIDPtr(ProductID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == nil,
		IsZero: ProductID != nil && (*ProductID) == 0,
	}
}

func (ft *ShopProductCollectionFilters) ByCollectionID(CollectionID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "collection_id",
		Value:  CollectionID,
		IsNil:  CollectionID == 0,
	}
}

func (ft *ShopProductCollectionFilters) ByCollectionIDPtr(CollectionID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "collection_id",
		Value:  CollectionID,
		IsNil:  CollectionID == nil,
		IsZero: CollectionID != nil && (*CollectionID) == 0,
	}
}

func (ft *ShopProductCollectionFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopProductCollectionFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopProductCollectionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopProductCollectionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopProductCollectionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopProductCollectionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type ShopBrandFilters struct{ prefix string }

func NewShopBrandFilters(prefix string) ShopBrandFilters {
	return ShopBrandFilters{prefix}
}

func (ft *ShopBrandFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopBrandFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopBrandFilters) ByID(ID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShopBrandFilters) ByIDPtr(ID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShopBrandFilters) ByShopID(ShopID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopBrandFilters) ByShopIDPtr(ShopID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopBrandFilters) ByBrandName(BrandName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "brand_name",
		Value:  BrandName,
		IsNil:  BrandName == "",
	}
}

func (ft *ShopBrandFilters) ByBrandNamePtr(BrandName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "brand_name",
		Value:  BrandName,
		IsNil:  BrandName == nil,
		IsZero: BrandName != nil && (*BrandName) == "",
	}
}

func (ft *ShopBrandFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *ShopBrandFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *ShopBrandFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopBrandFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopBrandFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopBrandFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopBrandFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShopBrandFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}
