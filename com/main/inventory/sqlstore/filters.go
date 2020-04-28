package sqlstore

import "o.o/backend/pkg/common/sql/sqlstore"

var SortInventoryVoucher = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterInventoryVoucher = sqlstore.FilterWhitelist{
	Equals:  []string{"trader_id", "type", "ref_name", "code", "ref_code", "rollback", "created_by", "updated_by"},
	Status:  []string{"status"},
	Arrays:  []string{"variant_ids", "product_ids"},
	Numbers: []string{"total_amount"},
	Dates:   []string{"created_at"},
}
