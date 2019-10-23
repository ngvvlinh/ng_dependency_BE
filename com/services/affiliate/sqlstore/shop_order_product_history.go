package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sqlstore"

	"etop.vn/backend/com/services/affiliate/model"

	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
)

func NewShopOrderProductHistoryStore(db *cmsql.Database) ShopOrderProductHistoryStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopOrderProductHistoryStore {
		return &ShopOrderProductHistoryStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type ShopOrderProductHistoryStoreFactory func(ctx context.Context) *ShopOrderProductHistoryStore

type ShopOrderProductHistoryStore struct {
	query func() cmsql.QueryInterface
	preds []interface{}

	ft ShopOrderProductHistoryFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *ShopOrderProductHistoryStore) UserID(id int64) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) ShopID(id int64) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) ProductID(id int64) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) OrderID(id int64) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) GetShopOrderProductHistoryDB() (*model.ShopOrderProductHistory, error) {
	var shopOrderProductHistory model.ShopOrderProductHistory
	err := s.query().Where(s.preds).ShouldGet(&shopOrderProductHistory)
	return &shopOrderProductHistory, err
}

func (s *ShopOrderProductHistoryStore) GetShopOrderProductHistoriesDB() ([]*model.ShopOrderProductHistory, error) {
	var results model.ShopOrderProductHistories
	err := s.query().Where(s.preds).Find(&results)
	return results, err
}

func (s *ShopOrderProductHistoryStore) CreateShopOrderProductHistory(shopOrderProductHistory *model.ShopOrderProductHistory) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(shopOrderProductHistory)
	return err
}

func (s *ShopOrderProductHistoryStore) UpdateShopOrderProductHistory(shopOrderProductHistory *model.ShopOrderProductHistory) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.UserID(shopOrderProductHistory.UserID).OrderID(shopOrderProductHistory.OrderID).ProductID(shopOrderProductHistory.ProductID).query().Where(s.preds).Update(shopOrderProductHistory)
	return err
}

func (s *ShopOrderProductHistoryStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ShopOrderProductHistory)(nil))
}
