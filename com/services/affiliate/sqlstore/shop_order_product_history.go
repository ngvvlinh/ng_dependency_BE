package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/backend/com/services/affiliate/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

func NewShopOrderProductHistoryStore(db *cmsql.Database) ShopOrderProductHistoryStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopOrderProductHistoryStore {
		return &ShopOrderProductHistoryStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopOrderProductHistoryStoreFactory func(ctx context.Context) *ShopOrderProductHistoryStore

type ShopOrderProductHistoryStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft ShopOrderProductHistoryFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *ShopOrderProductHistoryStore) UserID(id dot.ID) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) CustomerPolicyGroup(id dot.ID) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByCustomerPolicyGroupID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) ShopID(id dot.ID) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) ProductID(id dot.ID) *ShopOrderProductHistoryStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *ShopOrderProductHistoryStore) OrderID(id dot.ID) *ShopOrderProductHistoryStore {
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

func (s *ShopOrderProductHistoryStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ShopOrderProductHistory)(nil))
}
