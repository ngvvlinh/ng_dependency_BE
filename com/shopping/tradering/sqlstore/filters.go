package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *ShopTraderFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortTrader = map[string]string{
	"id": "",
}

var FilterTrader = sqlstore.FilterWhitelist{}
