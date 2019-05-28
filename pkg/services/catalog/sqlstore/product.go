package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
)

type ProductStore struct {
	db cmsql.Database
}

func NewProductStore(db cmsql.Database) *ProductStore {
	return &ProductStore{db}
}

func (s *ProductStore) WithContext(ctx context.Context) *ProductStoreWithContext {
	return &ProductStoreWithContext{query: s.db.WithContext(ctx)}
}

type ProductStoreWithContext struct {
	query cmsql.QueryInterface
	preds []interface{}

	includeDeleted bool
}
