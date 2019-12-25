package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/convert"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type AffiliateReferralCodeStoreFactory func(ctx context.Context) *AffiliateReferralCodeStore

func NewAffiliateReferralCodeStore(db *cmsql.Database) AffiliateReferralCodeStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AffiliateReferralCodeStore {
		return &AffiliateReferralCodeStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AffiliateReferralCodeStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft AffiliateReferralCodeFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *AffiliateReferralCodeStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.SellerCommission)(nil))
}

func (s *AffiliateReferralCodeStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *AffiliateReferralCodeStore) Paging(paging meta.Paging) *AffiliateReferralCodeStore {
	s.paging = paging
	return s
}

func (s *AffiliateReferralCodeStore) Filters(filters meta.Filters) *AffiliateReferralCodeStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *AffiliateReferralCodeStore) ID(id dot.ID) *AffiliateReferralCodeStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *AffiliateReferralCodeStore) Code(code string) *AffiliateReferralCodeStore {
	s.preds = append(s.preds, s.ft.ByCode(code))
	return s
}

func (s *AffiliateReferralCodeStore) AffiliateID(id dot.ID) *AffiliateReferralCodeStore {
	s.preds = append(s.preds, s.ft.ByAffiliateID(id))
	return s
}

func (s *AffiliateReferralCodeStore) GetAffiliateReferralCodeDB() (*model.AffiliateReferralCode, error) {
	var affiliateReferralCode model.AffiliateReferralCode
	err := s.query().Where(s.preds).ShouldGet(&affiliateReferralCode)
	return &affiliateReferralCode, err
}

func (s *AffiliateReferralCodeStore) GetAffiliateReferralCode() (*affiliate.AffiliateReferralCode, error) {
	affiliateReferralCode, err := s.GetAffiliateReferralCodeDB()
	if err != nil {
		return nil, err
	}
	return convert.AffiliateReferralCode(affiliateReferralCode), nil
}

func (s *AffiliateReferralCodeStore) GetAffiliateReferralCodes() ([]*affiliate.AffiliateReferralCode, error) {
	var results model.AffiliateReferralCodes
	err := s.query().Where(s.preds).Find(&results)
	return convert.AffiliateReferralCodes(results), err
}

func (s *AffiliateReferralCodeStore) CreateAffiliateReferralCode(affiliateReferralCode *model.AffiliateReferralCode) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(affiliateReferralCode)
	return err
}

func (s *AffiliateReferralCodeStore) UpdateAffiliateReferralCode(affiliateReferralCode *model.AffiliateReferralCode) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ID(affiliateReferralCode.ID).query().Where(s.preds).Update(affiliateReferralCode)
	return err
}
