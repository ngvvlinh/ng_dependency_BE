package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
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
	"id":         "",
	"created_at": "",
	"updated_at": "",
}

var FilterCustomerGroup = sqlstore.FilterWhitelist{}

var SortCustomer = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "updated_at",
	"name":       "",
	"code":       "",
}

var FilterCustomer = sqlstore.FilterWhitelist{
	Contains: []string{"full_name", "phone"},
	Equals:   []string{"type", "code", "email"},
	Numbers:  []string{"id"},
}

var SortShopCustomerGroupCustomer = map[string]string{
	"group_id":    "",
	"customer_id": "",
	"created_at":  "",
	"updated_at":  "",
}

var FilterCustomerGroupCustomer = sqlstore.FilterWhitelist{}

var SortShopTraderAddress = map[string]string{
	"id":         "",
	"updated_at": "",
	"created_at": "",
	"trader_id":  "",
}
