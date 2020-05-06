package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *FbExternalPageFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var FilterFbExternalPage = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at"},
	Numbers: []string{"id"},
	Status:  []string{"status", "connection_status"},
}

var SortFbExternalPage = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "",
}

var FilterFbExternalPageInternal = sqlstore.FilterWhitelist{
	Dates:   []string{"updated_at"},
	Equals:  []string{"id", "token"},
	Numbers: []string{"expires_in"},
}

var SortFbExternalPageInternal = map[string]string{
	"id":         "id",
	"updated_at": "",
}
