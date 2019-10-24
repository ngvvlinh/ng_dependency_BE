package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sqlstore"

	"etop.vn/api/meta"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
)

type OrderCreatedNotifyStoreFactory func(ctx context.Context) *OrderCreatedNotifyStore

func NewOrderCreatedNotifyStore(db *cmsql.Database) OrderCreatedNotifyStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *OrderCreatedNotifyStore {
		return &OrderCreatedNotifyStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type OrderCreatedNotifyStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft OrderCreatedNotifyFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *OrderCreatedNotifyStore) Pred(pred interface{}) *OrderCreatedNotifyStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *OrderCreatedNotifyStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ProductPromotion)(nil))
}

func (s *OrderCreatedNotifyStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *OrderCreatedNotifyStore) Paging(paging meta.Paging) *OrderCreatedNotifyStore {
	s.paging = paging
	return s
}

func (s *OrderCreatedNotifyStore) ID(id int64) *OrderCreatedNotifyStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *OrderCreatedNotifyStore) OrderID(id int64) *OrderCreatedNotifyStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *OrderCreatedNotifyStore) ReferralCode(code string) *OrderCreatedNotifyStore {
	s.preds = append(s.preds, s.ft.ByReferralCode(code))
	return s
}

func (s *OrderCreatedNotifyStore) GetOrderCreatedNotifyDB() (*model.OrderCreatedNotify, error) {
	var orderCreatedNotify model.OrderCreatedNotify
	err := s.query().Where(s.preds).ShouldGet(&orderCreatedNotify)
	return &orderCreatedNotify, err
}

func (s *OrderCreatedNotifyStore) GetOrderCreatedNotifiesDB() ([]*model.OrderCreatedNotify, error) {
	var results model.OrderCreatedNotifies
	err := s.query().Where(s.preds).Find(&results)
	return results, err
}

func (s *OrderCreatedNotifyStore) CreateOrderCreatedNotify(orderCreatedNotify *model.OrderCreatedNotify) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(orderCreatedNotify)
	return err
}

func (s *OrderCreatedNotifyStore) UpdateOrderCreatedNotify(orderCreatedNotify *model.OrderCreatedNotify) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Where(`id = ?`, orderCreatedNotify.ID).Update(orderCreatedNotify)
	return err
}
