package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/sqlstore"
)

func (ft *ReceiptFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortReceipt = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
	"title":      "title",
	"paid_at":    "paid_at",
}

var FilterReceipt = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at", "paid_at", "cancelled_at", "confirmed_at"},
	Equals:  []string{"type", "ledger_id", "code", "created_by", "id", "trader_id", "amount"},
	Numbers: []string{},
	Status:  []string{"status"},
	Arrays:  []string{"ref_ids"},
}
