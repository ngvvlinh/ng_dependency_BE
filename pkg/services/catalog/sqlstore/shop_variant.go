package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
)

type ShopVariantStoreFactory func(context.Context) *ShopVariantStore

func NewShopVariantStore(db cmsql.Database) ShopVariantStoreFactory {
	return func(ctx context.Context) *ShopVariantStore {
		return &ShopVariantStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShopVariantStore struct {
	FtProduct     ProductFilters
	FtShopProduct ShopProductFilters

	// unexported
	ftVariant     VariantFilters
	ftShopVariant ShopVariantFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopVariantStore) extend() *ShopVariantStore {
	s.FtProduct.prefix = "p"
	s.ftVariant.prefix = "v"
	s.FtShopProduct.prefix = "sp"
	s.ftShopVariant.prefix = "sv"
	return s
}
