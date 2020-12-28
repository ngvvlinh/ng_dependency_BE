package sqlstore

import "o.o/backend/pkg/common/sql/sq"

func (ft FulfillmentFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortFulfillment = map[string]string{
	"id":         "",
	"updated_at": "",
	"created_at": "",
}

var SortFulfillmentExtended = map[string]string{
	"id":         "f.id",
	"updated_at": "f.updated_at",
	"created_at": "f.created_at",
}
