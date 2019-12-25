package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

func NewOrderCommissionSettingStoreFactory(db *cmsql.Database) OrderCommissionSettingStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *OrderCommissionSettingStore {
		return &OrderCommissionSettingStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type OrderCommissionSettingStoreFactory func(ctx context.Context) *OrderCommissionSettingStore

type OrderCommissionSettingStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft OrderCommissionSettingFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *OrderCommissionSettingStore) Pred(pred interface{}) *OrderCommissionSettingStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *OrderCommissionSettingStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.OrderCommissionSetting)(nil))
}

func (s *OrderCommissionSettingStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *OrderCommissionSettingStore) Paging(paging meta.Paging) *OrderCommissionSettingStore {
	s.paging = paging
	return s
}

func (s *OrderCommissionSettingStore) SupplyID(id dot.ID) *OrderCommissionSettingStore {
	s.preds = append(s.preds, s.ft.BySupplyID(id))
	return s
}

func (s *OrderCommissionSettingStore) OrderID(id dot.ID) *OrderCommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *OrderCommissionSettingStore) ProductID(id dot.ID) *OrderCommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *OrderCommissionSettingStore) GetOrderCommissionSettingDB() (*model.OrderCommissionSetting, error) {
	var orderCommissionSetting model.OrderCommissionSetting
	err := s.query().Where(s.preds).ShouldGet(&orderCommissionSetting)
	return &orderCommissionSetting, err
}

func (s *OrderCommissionSettingStore) GetOrderCommissionSettingsDB() ([]*model.OrderCommissionSetting, error) {
	var results model.OrderCommissionSettings
	err := s.query().Where(s.preds).Find(&results)
	return results, err
}

func (s *OrderCommissionSettingStore) CreateOrderCommissionSetting(orderCommissionSetting *model.OrderCommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(orderCommissionSetting)
	return err
}

func (s *OrderCommissionSettingStore) UpdateOrderCommissionSetting(orderCommissionSetting *model.OrderCommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.SupplyID(orderCommissionSetting.SupplyID).OrderID(orderCommissionSetting.OrderID).query().Where(s.preds).Update(orderCommissionSetting)
	return err
}
