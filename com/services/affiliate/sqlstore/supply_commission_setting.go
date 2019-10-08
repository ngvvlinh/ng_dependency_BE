package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sq"

	"etop.vn/backend/com/services/affiliate/convert"

	"etop.vn/api/services/affiliate"

	"etop.vn/backend/com/services/affiliate/model"

	"etop.vn/backend/pkg/common/sqlstore"

	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
)

type SupplyCommissionSettingStoreFactory func(ctx context.Context) *SupplyCommissionSettingStore

func NewSupplyCommissionSettingStore(db cmsql.Database) SupplyCommissionSettingStoreFactory {
	return func(ctx context.Context) *SupplyCommissionSettingStore {
		return &SupplyCommissionSettingStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type SupplyCommissionSettingStore struct {
	query func() cmsql.QueryInterface
	preds []interface{}

	ft SupplyCommissionSettingFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *SupplyCommissionSettingStore) ShopID(id int64) *SupplyCommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *SupplyCommissionSettingStore) ProductID(id int64) *SupplyCommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *SupplyCommissionSettingStore) ProductIDs(ids ...int64) *SupplyCommissionSettingStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *SupplyCommissionSettingStore) GetSupplyCommissionSettingDB() (*model.SupplyCommissionSetting, error) {
	var supplyCommissionSetting model.SupplyCommissionSetting
	err := s.query().Where(s.preds).ShouldGet(&supplyCommissionSetting)
	return &supplyCommissionSetting, err
}

func (s *SupplyCommissionSettingStore) GetSupplyCommissionSetting() (*affiliate.SupplyCommissionSetting, error) {
	supplyCommissionSetting, err := s.GetSupplyCommissionSettingDB()
	if err != nil {
		return nil, err
	}
	return convert.SupplyCommissionSetting(supplyCommissionSetting), err
}

func (s *SupplyCommissionSettingStore) GetSupplyCommissionSettings() ([]*affiliate.SupplyCommissionSetting, error) {
	var results model.SupplyCommissionSettings
	err := s.query().Where(s.preds).Find(&results)
	return convert.SupplyCommissionSettings(results), err
}

func (s *SupplyCommissionSettingStore) CreateSupplyCommissionSetting(supplyCommissionSetting *model.SupplyCommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(supplyCommissionSetting)
	return err
}

func (s *SupplyCommissionSettingStore) UpdateSupplyCommissionSetting(supplyCommissionSetting *model.SupplyCommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ShopID(supplyCommissionSetting.ShopID).ProductID(supplyCommissionSetting.ProductID).query().Where(s.preds).Update(supplyCommissionSetting)
	return err
}
