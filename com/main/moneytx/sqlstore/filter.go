package sqlstore

import "o.o/backend/pkg/common/sql/sqlstore"

var filterMoneyTxShippingExtendedWhitelist = sqlstore.FilterWhitelist{
	Arrays:   nil,
	Contains: []string{"shop.name"},
	Dates:    []string{"created_at", "updated_at", "confirmed_at", "etop_transfered_at"},
	Equals: []string{
		"code", "shop.money_transaction_rrule", "shop.name",
		"shop.phone", "money_transaction_shipping_etop_id",
	},
	Numbers:  []string{"total_orders", "total_cod", "total_amount", "total_fee"},
	Status:   []string{"status"},
	Nullable: []string{"shop.bank_account"},
	PrefixOrRename: map[string]string{
		"shop.name":                    "s.name",
		"shop.phone":                   "s.phone",
		"shop.money_transaction_rrule": "s.money_transaction_rrule",
		"shop.bank_account":            "s.bank_account",
		"code":                         "m",
		"created_at":                   "m",
		"updated_at":                   "m",
		"confirmed_at":                 "m",
		"etop_transfered_at":           "m",
		"status":                       "m",
	},
}

var filterMoneyTxShippingWhitelist = sqlstore.FilterWhitelist{
	Arrays:         nil,
	Contains:       []string{},
	Dates:          []string{"created_at", "updated_at", "confirmed_at", "etop_transfered_at"},
	Equals:         []string{"code", "money_transaction_shipping_etop_id"},
	Numbers:        []string{"total_orders", "total_cod", "total_amount", "total_fee"},
	Status:         []string{"status"},
	PrefixOrRename: map[string]string{},
}

var filterMoneyTxShippingExternalWhitelist = sqlstore.FilterWhitelist{
	Arrays:         nil,
	Contains:       []string{},
	Dates:          []string{"created_at", "updated_at", "external_paid_at"},
	Equals:         []string{"code", "provider"},
	Numbers:        []string{"total_orders", "total_cod"},
	Status:         []string{"status"},
	PrefixOrRename: map[string]string{},
}

var SortMoneyTx = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
}
