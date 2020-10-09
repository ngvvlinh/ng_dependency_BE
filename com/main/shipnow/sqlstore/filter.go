package sqlstore

import (
	"o.o/backend/pkg/common/sql/sqlstore"
)

var (
	SortShipnow = map[string]string{
		"id":         "",
		"created_at": "",
		"updated_at": "",
		"name":       "",
	}

	FilterShipnowWhitelist = sqlstore.FilterWhitelist{
		Equals: []string{
			"shipping_code", "carrier", "shipping_state",
			"address_to.province_code", "address_to.district_code",
			"connection_id",
		},
		Dates:    []string{"created_at", "updated_at"},
		Status:   []string{"status", "shipping_status"},
		Contains: []string{"address_to.full_name", "address_to.phone"},
		PrefixOrRename: map[string]string{
			"address_to.province_code": "shipnow_fulfillment.address_to_province_code",
			"address_to.district_code": "shipnow_fulfillment.address_to_district_code",
			"address_to.full_name":     "shipnow_fulfillment.address_to_full_name_norm",
			"address_to.phone":         "shipnow_fulfillment.address_to_phone",
		},
	}
)
