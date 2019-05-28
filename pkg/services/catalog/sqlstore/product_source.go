package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
)

type ProductSourceStoreFactory func(context.Context) *ProductSourceStore

func NewProductSourceStore(db cmsql.Database) ProductSourceStoreFactory {
	return func(ctx context.Context) *ProductSourceStore {
		return &ProductSourceStore{query: db.WithContext(ctx)}
	}
}

type ProductSourceStore struct {
	query cmsql.Query
}
