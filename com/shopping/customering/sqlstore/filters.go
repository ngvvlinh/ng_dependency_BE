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

func (ft *ShopCustomerGroupFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortShopCustomerGroup = map[string]string{
	"group_id":   "",
	"created_at": "",
	"updated_at": "",
}

var FilterCustomerGroup = sqlstore.FilterWhitelist{}

var SortCustomer = map[string]string{
	"id":         "",
	"created_at": "",
	"updated_at": "",
	"name":       "",
}

var FilterCustomer = sqlstore.FilterWhitelist{
	Contains: []string{"full_name"},
	Equals:   []string{"phone"},
}

var SortShopCustomerGroupCustomer = map[string]string{
	"group_id":    "",
	"customer_id": "",
	"created_at":  "",
	"updated_at":  "",
}

var FilterCustomerGroupCustomer = sqlstore.FilterWhitelist{}
