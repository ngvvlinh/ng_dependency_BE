package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

func (ft *ShopLedgerFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortLedger = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
}

var FilterLedger = sqlstore.FilterWhitelist{
	Equals: []string{"type"},
}
