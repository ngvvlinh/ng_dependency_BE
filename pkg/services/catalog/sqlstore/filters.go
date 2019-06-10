package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

func (ft ProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft VariantFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopProductFilters) NotDeleted() sq.WriterTo {
	// shop_product does not use deleted_at
	return nil
}

func (ft ShopVariantFilters) NotDeleted() sq.WriterTo {
	// shop_variant does not use deleted_at
	return nil
}

var (
	SortProduct = map[string]string{
		"id":         "",
		"created_at": "",
		"updated_at": "",
		"name":       "",
	}

	SortVariant = map[string]string{
		"id": "",
	}

	SortShopProduct = map[string]string{
		"product_id": "sp.product_id",
		"created_at": "sp.created_at",
		"updated_at": "sp.updated_at",
	}

	FilterProductWhitelist = sqlstore.FilterWhitelist{
		Arrays:   []string{},
		Contains: []string{"name"},
		Equals:   []string{"etop_category_id", "name"},
		Status:   []string{"ed_status", "status", "etop_status"},
		Numbers:  []string{"wholesale_price", "list_price", "retail_price_min", "retail_price_max", "ed_wholesale_price", "ed_list_price", "ed_retail_price_max"},
		Dates:    []string{"created_at", "updated_at"},
		Unaccent: []string{"name"},
		PrefixOrRename: map[string]string{
			"name":       "p",
			"status":     "p",
			"created_at": "p",
			"updated_at": "p",

			"wholesale_price":     "v",
			"list_price":          "v",
			"retail_price_min":    "v",
			"retail_price_max":    "v",
			"ed_wholesale_price":  "v",
			"ed_list_price":       "v",
			"ed_retail_price_max": "v",
		},
	}

	FilterVariantWhitelist = sqlstore.FilterWhitelist{
		Arrays:   []string{},
		Contains: []string{"name"},
		Equals:   []string{"name"},
		Status:   []string{"ed_status", "status", "etop_status"},
		Numbers:  []string{"wholesale_price", "list_price", "retail_price_min", "retail_price_max", "ed_wholesale_price", "ed_list_price", "ed_retail_price_max"},
	}

	FilterShopProductWhitelist = sqlstore.FilterWhitelist{
		Arrays:   []string{"tags"},
		Contains: []string{"external_name", "name"},
		Equals:   []string{"external_code", "external_base_id", "external_id", "collection_id"},
		Status:   []string{"external_status", "ed_status", "status", "etop_status"},
		Numbers:  []string{"retail_price"},
		Dates:    []string{"created_at", "updated_at"},
		Unaccent: []string{"product.name"},

		PrefixOrRename: map[string]string{
			"name":       "sp",
			"status":     "sp",
			"created_at": "sp",
			"updated_at": "sp",

			"product.name": "p.name_norm_ua",
		},
	}
)
