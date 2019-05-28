package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
)

type ProductSourceStore struct {
	db cmsql.Database
}

func NewProductSourceStore(db cmsql.Database) *ProductSourceStore {
	return &ProductSourceStore{db}
}

func (s *ProductSourceStore) WithContext(ctx context.Context) *ProductStoreWithContext {
	return &ProductStoreWithContext{query: s.db.WithContext(ctx)}
}

type ProductSourceStoreWithContext struct {
	db cmsql.QueryInterface
}
