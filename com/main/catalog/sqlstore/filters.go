package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

func (ft ShopProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopVariantFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortShopVariant = map[string]string{
	"id": "",
}

var SortShopProduct = map[string]string{
	"product_id": "",
	"created_at": "",
	"updated_at": "",
}

var FilterShopProduct = sqlstore.FilterWhitelist{
	Arrays:   []string{"tags"},
	Contains: []string{"external_name", "name"},
	Equals:   []string{"external_code", "external_base_id", "external_id", "collection_id"},
	Status:   []string{"external_status", "ed_status", "status", "etop_status"},
	Numbers:  []string{"retail_price"},
	Dates:    []string{"created_at", "updated_at"},
	Unaccent: []string{"product.name"},

	PrefixOrRename: map[string]string{
		"product.name": "name_norm_ua",
	},
}