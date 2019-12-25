package cmsql

import (
	"etop.vn/backend/pkg/common/metrics"
	"etop.vn/backend/pkg/common/sql/sq"
)

// TODO: explain analyze
func monitorQuery(entry *sq.LogEntry) {
	if entry.IsQuery() {
		metrics.DatabaseQuery(entry)
		return
	}

	t := entry.Type()
	if t == sq.TypeCommit || t == sq.TypeRollback {
		metrics.DatabaseTransaction(entry)
	}
}
