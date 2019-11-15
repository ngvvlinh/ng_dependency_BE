package sqlstore

import "etop.vn/backend/pkg/etop/sqlstore"

var SortInventoryVoucher = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}
var FilterInventoryVoucher = sqlstore.FilterWhitelist{
	Equals: []string{"trader_id", "variant_id", "type", "ref_name", "code"},
	Status: []string{"status"},
}
