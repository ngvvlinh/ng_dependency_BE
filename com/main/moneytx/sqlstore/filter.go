package sqlstore

import "etop.vn/backend/pkg/common/sql/sqlstore"

var filterMoneyTransactionShippingWhitelist = sqlstore.FilterWhitelist{
	Arrays:   nil,
	Contains: []string{"shop.name"},
	Dates:    []string{"created_at", "updated_at", "confirmed_at", "etop_transfered_at"},
	Equals: []string{
		"code", "shop.money_transaction_rrule", "shop.name",
		"shop.phone",
	},
	Numbers: []string{"total_orders", "total_cod", "total_amount", "total_fee"},
	Status:  []string{"status"},
	PrefixOrRename: map[string]string{
		"shop.name":                    "s.name",
		"shop.phone":                   "s.phone",
		"shop.money_transaction_rrule": "s.money_transaction_rrule",
		"code":                         "m",
		"created_at":                   "m",
		"updated_at":                   "m",
		"confirmed_at":                 "m",
		"etop_transfered_at":           "m",
		"status":                       "m",
	},
}

var filterMoneyTransactionWhitelist = sqlstore.FilterWhitelist{
	Arrays:         nil,
	Contains:       []string{},
	Dates:          []string{"created_at", "updated_at", "confirmed_at", "etop_transfered_at"},
	Equals:         []string{"code"},
	Numbers:        []string{"total_orders", "total_cod", "total_amount", "total_fee"},
	Status:         []string{"status"},
	PrefixOrRename: map[string]string{},
}
