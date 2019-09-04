package sqlstore

import (
	"context"

	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/convert"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type CommissionSettingStoreFactory func(ctx context.Context) *CommissionSettingStore

func NewShopCommissionSettingStore(db cmsql.Database) CommissionSettingStoreFactory {
	return func(ctx context.Context) *CommissionSettingStore {
		return &CommissionSettingStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type CommissionSettingStore struct {
	query func() cmsql.QueryInterface
	ft    CommissionSettingFilters
	preds []interface{}
}

func (s *CommissionSettingStore) ProductID(id int64) *CommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *CommissionSettingStore) ProductIDs(ids []int64) *CommissionSettingStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *CommissionSettingStore) AccountID(id int64) *CommissionSettingStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *CommissionSettingStore) GetShopCommissionSettingDB() (*model.CommissionSetting, error) {
	var commissionSetting model.CommissionSetting
	err := s.query().Where(s.preds...).ShouldGet(&commissionSetting)
	return &commissionSetting, err
}

func (s *CommissionSettingStore) GetShopCommissionSetting() (*affiliate.CommissionSetting, error) {
	commissionSetting, err := s.GetShopCommissionSettingDB()
	if err != nil {
		return nil, err
	}
	return convert.CommissionSetting(commissionSetting), nil
}

func (s *CommissionSettingStore) GetShopCommissionSettings() ([]*affiliate.CommissionSetting, error) {
	var results model.CommissionSettings
	err := s.query().Where(s.preds).Find(&results)
	return convert.CommissionSettings(results), err
}

func (s *CommissionSettingStore) GetShopCommissionSettingByProductID(ShopID int64, productID int64) (*affiliate.CommissionSetting, error) {
	return s.AccountID(ShopID).ProductID(productID).GetShopCommissionSetting()
}

func (s *CommissionSettingStore) CreateShopCommissionSetting(commissionSetting *model.CommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(commissionSetting)
	return err
}

func (s *CommissionSettingStore) UpdateShopCommissionSetting(commissionSetting *model.CommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.AccountID(commissionSetting.AccountID).ProductID(commissionSetting.ProductID).query().Where(s.preds).Update(commissionSetting)
	return err
}
