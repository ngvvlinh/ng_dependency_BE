package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

func (ft *FbPageFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var FilterFbPage = sqlstore.FilterWhitelist{
	Dates:   []string{"created_at", "updated_at"},
	Numbers: []string{"id"},
	Status:  []string{"status", "connection_status"},
}

var SortFbPage = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "",
}

var FilterFbPageInternal = sqlstore.FilterWhitelist{
	Dates:   []string{"updated_at"},
	Equals:  []string{"id", "token"},
	Numbers: []string{"expires_in"},
}

var SortFbPageInternal = map[string]string{
	"id":         "id",
	"updated_at": "",
}
