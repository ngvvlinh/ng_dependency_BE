package sqlstore

import "o.o/backend/pkg/common/sql/sq"

func (ft *HotlineFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ExtensionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *TenantFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortCallLog = map[string]string{
	"id":         "id",
	"created_at": "",
	"updated_at": "",
	"started_at": "started_at",
}

var SortTenant = map[string]string{
	"id":         "id",
	"created_at": "created_at",
}
