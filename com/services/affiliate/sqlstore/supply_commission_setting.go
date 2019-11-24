package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/convert"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type SupplyCommissionSettingStoreFactory func(ctx context.Context) *SupplyCommissionSettingStore

func NewSupplyCommissionSettingStore(db *cmsql.Database) SupplyCommissionSettingStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *SupplyCommissionSettingStore {
		return &SupplyCommissionSettingStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type SupplyCommissionSettingStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft SupplyCommissionSettingFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *SupplyCommissionSettingStore) ShopID(id dot.ID) *SupplyCommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *SupplyCommissionSettingStore) ProductID(id dot.ID) *SupplyCommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *SupplyCommissionSettingStore) ProductIDs(ids ...dot.ID) *SupplyCommissionSettingStore {
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
	return s.ShopID(supplyCommissionSetting.ShopID).ProductID(supplyCommissionSetting.ProductID).query().Where(s.preds).UpdateAll().ShouldUpdate(supplyCommissionSetting)
}
