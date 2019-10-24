package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
)

type ShopCashbackStoreFactory func(ctx context.Context) *ShopCashbackStore

func NewShopCashbackStore(db *cmsql.Database) ShopCashbackStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopCashbackStore {
		return &ShopCashbackStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopCashbackStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft ShopCashbackFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *ShopCashbackStore) ID(id int64) *ShopCashbackStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShopCashbackStore) OrderID(id int64) *ShopCashbackStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *ShopCashbackStore) GetShopCashbackDB() (*model.ShopCashback, error) {
	var shopCashback model.ShopCashback
	err := s.query().Where(s.preds).ShouldGet(&shopCashback)
	return &shopCashback, err
}

func (s *ShopCashbackStore) GetShopCashbacksDB() ([]*model.ShopCashback, error) {
	var shopCashback model.ShopCashbacks
	err := s.query().Where(s.preds).Find(&shopCashback)
	return shopCashback, err
}

func (s *ShopCashbackStore) CreateShopCashback(shopCashback *model.ShopCashback) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(shopCashback)
	return err
}

func (s *ShopCashbackStore) UpdateShopCashback(shopCashback *model.ShopCashback) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ID(shopCashback.ID).query().Where(s.preds).Update(shopCashback)
	return err
}
