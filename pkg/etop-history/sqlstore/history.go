package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq/core"
)

type HistoryStoreFactory func(context.Context) *HistoryStore

type HistoryStore struct {
	query cmsql.QueryFactory
}

func NewHistoryStore(db *cmsql.Database) HistoryStoreFactory {
	return func(ctx context.Context) *HistoryStore {
		return &HistoryStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

func (s *HistoryStore) GetHistory(_model core.IGet, rid int64) (bool, error) {
	return s.query().Where("rid = ?", rid).Get(_model)
}
