package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/sqlstore"
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