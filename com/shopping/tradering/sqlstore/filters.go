package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

func (ft *ShopTraderFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortTrader = map[string]string{
	"id": "",
}

var FilterTrader = sqlstore.FilterWhitelist{}
