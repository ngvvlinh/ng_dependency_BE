package sqlstore

import "etop.vn/backend/pkg/common/sql/sq"

func (ft *ShipmentPriceFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShipmentPriceFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}
