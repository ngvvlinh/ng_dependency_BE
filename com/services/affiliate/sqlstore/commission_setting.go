package sqlstore

import (
	"context"

	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/convert"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type CommissionSettingStoreFactory func(ctx context.Context) *AffiliateCommissonStore

func NewCommissionSettingStore(db *cmsql.Database) CommissionSettingStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AffiliateCommissonStore {
		return &AffiliateCommissonStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AffiliateCommissonStore struct {
	query cmsql.QueryFactory
	ft    CommissionSettingFilters
	preds []interface{}
}

func (s *AffiliateCommissonStore) ProductID(id dot.ID) *AffiliateCommissonStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *AffiliateCommissonStore) ProductIDs(ids []dot.ID) *AffiliateCommissonStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *AffiliateCommissonStore) AccountID(id dot.ID) *AffiliateCommissonStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AffiliateCommissonStore) GetCommissionSettingDB() (*model.CommissionSetting, error) {
	var commissionSetting model.CommissionSetting
	err := s.query().Where(s.preds...).ShouldGet(&commissionSetting)
	return &commissionSetting, err
}

func (s *AffiliateCommissonStore) GetCommissionSetting() (*affiliate.CommissionSetting, error) {
	commissionSetting, err := s.GetCommissionSettingDB()
	if err != nil {
		return nil, err
	}
	return convert.CommissionSetting(commissionSetting), nil
}

func (s *AffiliateCommissonStore) GetCommissionSettings() ([]*affiliate.CommissionSetting, error) {
	var results model.CommissionSettings
	err := s.query().Where(s.preds).Find(&results)
	return convert.CommissionSettings(results), err
}

func (s *AffiliateCommissonStore) CreateCommissionSetting(commissionSetting *model.CommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(commissionSetting)
	return err
}

func (s *AffiliateCommissonStore) UpdateCommissionSetting(commissionSetting *model.CommissionSetting) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.AccountID(commissionSetting.AccountID).ProductID(commissionSetting.ProductID).query().Where(s.preds).UpdateAll().Update(commissionSetting)
	return err
}
