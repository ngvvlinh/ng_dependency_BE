package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/sqlstore"
)

func (ft *ShopSupplierFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortVendor = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"full_name":  "",
}

var FilterVendor = sqlstore.FilterWhitelist{}
