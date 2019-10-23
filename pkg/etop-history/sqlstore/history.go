package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq/core"
)

type HistoryStoreFactory func(context.Context) *HistoryStore

type HistoryStore struct {
	query func() cmsql.QueryInterface
}

func NewHistoryStore(db *cmsql.Database) HistoryStoreFactory {
	return func(ctx context.Context) *HistoryStore {
		return &HistoryStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

func (s *HistoryStore) GetHistory(_model core.IGet, rid int64) (bool, error) {
	return s.query().Where("rid = ?", rid).Get(_model)
}
