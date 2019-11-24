package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type OrderPromotionStoreFactory func(ctx context.Context) *OrderPromotionStore

func NewOrderPromotionStore(db *cmsql.Database) OrderPromotionStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *OrderPromotionStore {
		return &OrderPromotionStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type OrderPromotionStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft OrderPromotionFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *OrderPromotionStore) ID(id dot.ID) *OrderPromotionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *OrderPromotionStore) OrderCreatedNotifyID(id dot.ID) *OrderPromotionStore {
	s.preds = append(s.preds, s.ft.ByOrderCreatedNotifyID(id))
	return s
}

func (s *OrderPromotionStore) GetOrderPromotionDB() (*model.OrderPromotion, error) {
	var orderPromotion model.OrderPromotion
	err := s.query().Where(s.preds).ShouldGet(&orderPromotion)
	return &orderPromotion, err
}

func (s *OrderPromotionStore) GetOrderPromotionsDB() ([]*model.OrderPromotion, error) {
	var orderPromotions model.OrderPromotions
	err := s.query().Where(s.preds).Find(&orderPromotions)
	return orderPromotions, err
}

func (s *OrderPromotionStore) CreateOrderPromotion(orderPromotion *model.OrderPromotion) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(orderPromotion)
	return err
}

func (s *OrderPromotionStore) UpdateOrderPromotion(orderPromotion *model.OrderPromotion) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ID(orderPromotion.ID).query().Where(s.preds).Update(orderPromotion)
	return err
}
