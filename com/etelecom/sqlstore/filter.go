package sqlstore

import "o.o/backend/pkg/common/sql/sq"

func (ft *HotlineFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ExtensionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}
