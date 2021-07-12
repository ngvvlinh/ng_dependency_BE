package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

func (ft *ContactFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortContact = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterContact = sqlstore.FilterWhitelist{}

func (ft *ContactFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}
