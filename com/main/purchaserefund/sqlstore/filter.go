package sqlstore

import (
	"o.o/backend/pkg/common/sql/sqlstore"
)

var FilterPurchaseRefund = sqlstore.FilterWhitelist{
	Contains: []string{},
	Dates:    []string{"created_at", "updated_at", "cancelled_at", "confirmed_at"},
	Equals:   []string{"code", "created_by", "id", "order_id", "updated_by"},
	Numbers:  []string{},
	Status:   []string{"status"},
	Arrays:   []string{},
}

var SortPurchaseRefund = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}
