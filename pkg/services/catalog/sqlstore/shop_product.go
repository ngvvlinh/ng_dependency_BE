package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sql"
)

type ShopProductStoreFactory func(context.Context) *ShopProductStore

func NewShopProductStore(db cmsql.Database) ShopProductStoreFactory {
	return func(ctx context.Context) *ShopProductStore {
		return &ShopProductStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShopProductStore struct {
	FtProduct     ProductFilters
	FtShopProduct ShopProductFilters

	ftVariant     VariantFilters
	ftShopVariant ShopVariantFilters

	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted bool
}

func (s *ShopProductStore) Where(pred sql.FilterQuery) *ShopProductStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ShopProductStore) ID(id int64) *ShopProductStore {
	s.preds = append(s.preds, s.FtProduct.ByID(id))
	return s
}

func (s *ShopProductStore) ShopID(id int64) *ShopProductStore {
	s.preds = append(s.preds, s.FtShopProduct.ByShopID(id))
	return s
}
