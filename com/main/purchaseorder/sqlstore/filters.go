package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
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
	Dates:    []string{"created_at", "updated_at", "cancelled_at", "confirmed_at"},
	Equals:   []string{"code", "ref_code", "created_by", "id", "supplier_id"},
	Numbers:  []string{"total_amount"},
	Contains: []string{"supplier_full_name", "supplier_phone"},
	Status:   []string{"status"},
	Arrays:   []string{"variant_ids"},
}
