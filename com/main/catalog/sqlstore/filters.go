package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

func (ft ShopProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopVariantFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopVariantSupplierFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopCategoryFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopCollectionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopBrandFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortShopVariant = map[string]string{
	"id":         "variant_id",
	"updated_at": "updated_at",
	"created_at": "",
}

var SortShopProduct = map[string]string{
	"id":         "product_id",
	"product_id": "",
	"created_at": "",
	"updated_at": "updated_at",
}

var SortShopCategory = map[string]string{
	"category_id": "",
	"created_at":  "",
	"updated_at":  "",
}

var SortShopCollection = map[string]string{
	"id":            "id",
	"collection_id": "",
	"created_at":    "",
	"updated_at":    "updated_at",
	"deleted_at":    "",
}

var SortShopProductCollection = map[string]string{
	"product_id":    "",
	"collection_id": "",
	"created_at":    "",
	"updated_at":    "updated_at",
}

var SortShopBrand = map[string]string{
	"id":         "",
	"created_at": "",
}

var FilterShopProduct = sqlstore.FilterWhitelist{
	Arrays:   []string{"tags"},
	Contains: []string{"external_name", "name"},
	Equals:   []string{"external_code", "external_base_id", "external_id", "collection_id", "code"},
	Status:   []string{"external_status", "ed_status", "status", "etop_status"},
	Numbers:  []string{"retail_price", "list_price"},
	Dates:    []string{"created_at", "updated_at"},
	Unaccent: []string{"product.name"},

	PrefixOrRename: map[string]string{
		"product.name": "name_norm_ua",
	},
}
