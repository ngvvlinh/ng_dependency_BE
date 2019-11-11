package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/sqlstore"
)

func (ft *PurchaseOrderFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortPurchaseOrder = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
}

var FilterPurchaseOrder = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at", "cancelled_at", "confirmed_at"},
	Equals:  []string{"code", "created_by", "id", "supplier_id"},
	Numbers: []string{},
	Status:  []string{"status"},
	Arrays:  []string{"variant_ids"},
}
