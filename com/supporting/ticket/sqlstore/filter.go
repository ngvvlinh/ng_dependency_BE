package sqlstore

import "o.o/backend/pkg/common/sql/sq"

func (ft *TicketFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (ft *TicketLabelFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

var SortTicket = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (ft *TicketFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *TicketCommentFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *TicketLabelFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *TicketLabelExternalFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *TicketLabelTicketLabelExternalFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortTicketComment = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var SortTicketLabelExternal = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}
