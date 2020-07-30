package sqlstore

import "o.o/backend/pkg/common/sql/sq"

var SortTicket = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (ft *TicketCommentFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortTicketComment = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}
