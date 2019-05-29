package sqlstore

import "etop.vn/backend/pkg/common/sql"

func (ft FulfillmentFilters) NotDeleted() sql.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}
