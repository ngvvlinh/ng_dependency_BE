package sqlstore

import "etop.vn/backend/pkg/etop/sqlstore"

var SortInventoryVoucher = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterInventoryVoucher = sqlstore.FilterWhitelist{
	Equals:  []string{"trader_id", "type", "ref_name", "code", "ref_code"},
	Status:  []string{"status"},
	Arrays:  []string{"variant_ids","product_ids"},
	Numbers: []string{"total_amount"},
	Dates:   []string{"created_at"},
}
