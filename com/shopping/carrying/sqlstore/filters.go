package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *ShopCarrierFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortCarrier = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"full_name":  "",
}

var FilterCarrier = sqlstore.FilterWhitelist{}
