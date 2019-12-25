package sqlstore

import "etop.vn/backend/pkg/common/sql/sq"

func (ft FulfillmentFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}
