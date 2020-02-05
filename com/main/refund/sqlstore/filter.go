package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

var FilterRefund = sqlstore.FilterWhitelist{
	Contains: []string{},
	Dates:    []string{"created_at", "updated_at", "cancelled_at", "confirmed_at"},
	Equals:   []string{"code", "created_by", "id", "order_id"},
	Numbers:  []string{},
	Status:   []string{"status"},
	Arrays:   []string{},
}

var SortRefund = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}
