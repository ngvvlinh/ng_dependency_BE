package sqlstore

import (
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

func (ft *InvitationFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortInvitation = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterInvitation = sqlstore.FilterWhitelist{
	Arrays: []string{"roles"},
	Dates:  []string{"created_at", "updated_at", "accepted_at", "rejected_at", "expires_at"},
	Equals: []string{"id", "token"},
	Status: []string{"status"},
}
