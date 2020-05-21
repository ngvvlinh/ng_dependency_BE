package sqlstore

import (
	"o.o/backend/pkg/common/sql/sqlstore"
)

var FilterFbExternalUser = sqlstore.FilterWhitelist{
	Dates:  []string{"created_at", "updated_at"},
	Equals: []string{"type", "external_id"},
}

var SortFbExternalUser = map[string]string{
	"created_at": "",
	"updated_at": "",
}
