package sqlstore

import "o.o/backend/pkg/common/sql/sqlstore"

var SortSubscriptionBill = map[string]string{
	"created_at": "created_at",
}

var FilterSubscriptionBill = sqlstore.FilterWhitelist{
	Arrays:         nil,
	Bools:          nil,
	Contains:       nil,
	Dates:          []string{"created_at", "updated_at"},
	Equals:         []string{"account_id", "subscription_id"},
	Nullable:       nil,
	Numbers:        nil,
	Status:         []string{"status"},
	Unaccent:       nil,
	PrefixOrRename: nil,
}