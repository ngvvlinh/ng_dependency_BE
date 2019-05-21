package sqlstore

import sq "etop.vn/backend/pkg/common/sql"

func (ft ProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}
