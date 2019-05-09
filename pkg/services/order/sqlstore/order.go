package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/etop/model"
)

type OrderStore struct {
	ctx   context.Context
	db    cmsql.Database
	preds []interface{}
}

func New(db cmsql.Database) *OrderStore {
	return &OrderStore{db: db, ctx: context.Background()}
}

func (s *OrderStore) WithContext(ctx context.Context) *OrderStore {
	return &OrderStore{db: s.db, ctx: ctx}
}

func (s *OrderStore) ID(id int64) *OrderStore {
	s.preds = append(s.preds, sq.NewExpr("id = ?"))
	return s
}

func (s *OrderStore) Get() (*model.Order, error) {
	var item model.Order
	err := s.db.WithContext(s.ctx).ShouldGet(&item)
	return &item, err
}
