package sqlstore

import (
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

func (ft *ShopCustomerFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShopTraderAddressFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortCustomer = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"name":       "",
}

var FilterCustomer = sqlstore.FilterWhitelist{}
