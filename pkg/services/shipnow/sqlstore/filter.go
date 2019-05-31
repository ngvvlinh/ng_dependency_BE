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

	FilterShipnowWhitelist = sqlstore.FilterWhitelist{}
)
