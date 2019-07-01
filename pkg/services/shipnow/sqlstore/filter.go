package sqlstore

import (
	"etop.vn/backend/pkg/common/sqlstore"
)

var (
	SortShipnow = map[string]string{
		"id":         "",
		"created_at": "",
		"updated_at": "",
		"name":       "",
	}

	FilterShipnowWhitelist = sqlstore.FilterWhitelist{
		Equals: []string{"shipping_code", "carrier", "shipping_state"},
		Dates:  []string{"created_at", "updated_at"},
		Status: []string{"status", "shipping_status"},
	}
)