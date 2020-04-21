package sqlstore

import "o.o/backend/pkg/common/sql/sqlstore"

var SortSubscription = map[string]string{
	"created_at": "created_at",
}

var FilterSubscription = sqlstore.FilterWhitelist{
	Arrays:         nil,
	Bools:          []string{"cancel_at_period_end"},
	Contains:       nil,
	Dates:          []string{"created_at", "updated_at", "current_period_end_at", "current_period_start_at", "billing_cycle_anchor_at", "start_at"},
	Equals:         []string{"account_id"},
	Nullable:       nil,
	Numbers:        nil,
	Status:         []string{"status"},
	Unaccent:       nil,
	PrefixOrRename: nil,
}
