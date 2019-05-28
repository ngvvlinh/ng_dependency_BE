package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
)

type ProductStoreFactory func(context.Context) *ProductStoreWithContext

func NewProductStore(db cmsql.Database) ProductStoreFactory {
	return func(ctx context.Context) *ProductStoreWithContext {
		return &ProductStoreWithContext{
			query: db.WithContext(ctx),
		}
	}
}

type ProductStoreWithContext struct {
	query cmsql.QueryInterface
	preds []interface{}

	includeDeleted bool
}
