package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/etop/sqlstore"
)

func (ft *ShopSupplierFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortSupplier = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"full_name":  "",
}

var FilterSupplier = sqlstore.FilterWhitelist{
	Contains: []string{"full_name", "phone", "company_name"},
	Equals:   []string{"email", "code"},
}
