package sqlstore

import "etop.vn/backend/pkg/common/sql/sq"

func (ft *ConnectionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShopConnectionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}
