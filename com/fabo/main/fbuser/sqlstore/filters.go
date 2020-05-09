package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *FbExternalUserFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var FilterFbExternalUser = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at"},
	Numbers: []string{"id"},
	Equals:  []string{"type", "external_id"},
}

var SortFbExternalUser = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "",
}
