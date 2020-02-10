package sqlstore

import "etop.vn/backend/pkg/common/sql/sq"

func (ft *UserFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}
