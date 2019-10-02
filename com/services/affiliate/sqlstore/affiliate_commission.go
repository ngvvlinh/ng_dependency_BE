package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sqlstore"

	"etop.vn/api/services/affiliate"

	"etop.vn/backend/com/services/affiliate/convert"
	"etop.vn/backend/com/services/affiliate/model"

	"etop.vn/api/meta"

	"etop.vn/backend/pkg/common/cmsql"
)

type AffiliateCommissionStoreFactory func(ctx context.Context) *AffiliateCommissionStore

func NewAffiliateCommissionSettingStore(db cmsql.Database) AffiliateCommissionStoreFactory {
	return func(ctx context.Context) *AffiliateCommissionStore {
		return &AffiliateCommissionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type AffiliateCommissionStore struct {
	query func() cmsql.QueryInterface
	preds []interface{}

	ft AffiliateCommissionFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *AffiliateCommissionStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.AffiliateCommission)(nil))
}

func (s *AffiliateCommissionStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *AffiliateCommissionStore) Paging(paging meta.Paging) *AffiliateCommissionStore {
	s.paging = paging
	return s
}

func (s *AffiliateCommissionStore) Filters(filters meta.Filters) *AffiliateCommissionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *AffiliateCommissionStore) ID(id int64) *AffiliateCommissionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *AffiliateCommissionStore) AffiliateID(id int64) *AffiliateCommissionStore {
	s.preds = append(s.preds, s.ft.ByAffiliateID(id))
	return s
}

func (s *AffiliateCommissionStore) GetAffiliateCommissionDB() (*model.AffiliateCommission, error) {
	var affiliateCommission model.AffiliateCommission
	err := s.query().Where(s.preds).ShouldGet(&affiliateCommission)
	return &affiliateCommission, err
}

func (s *AffiliateCommissionStore) GetAffiliateCommission() (*affiliate.AffiliateCommission, error) {
	affiliateCommission, err := s.GetAffiliateCommissionDB()
	if err != nil {
		return nil, err
	}
	return convert.AffiliateCommission(affiliateCommission), nil
}

func (s *AffiliateCommissionStore) GetAffiliateCommissions() ([]*affiliate.AffiliateCommission, error) {
	var results model.AffiliateCommissions
	err := s.query().Where(s.preds).Find(&results)
	return convert.AffiliateCommissions(results), err
}

func (s *AffiliateCommissionStore) CreateAffiliateCommission(affiliateCommission *model.AffiliateCommission) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(affiliateCommission)
	return err
}

func (s *AffiliateCommissionStore) updateAffiliateCommission(affiliateCommission *model.AffiliateCommission) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ID(affiliateCommission.ID).query().Where(s.preds).Update(affiliateCommission)
	return err
}
