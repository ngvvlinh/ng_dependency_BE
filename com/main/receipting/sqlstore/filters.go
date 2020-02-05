package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
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
	Contains: []string{"trader_full_name", "trader_phone"},
	Dates:    []string{"created_at", "updated_at", "paid_at", "cancelled_at", "confirmed_at"},
	Equals:   []string{"type", "ledger_id", "code", "created_by", "id", "trader_id", "amount", "trader_type"},
	Numbers:  []string{"amount"},
	Status:   []string{"status"},
	Arrays:   []string{"ref_ids"},
}
