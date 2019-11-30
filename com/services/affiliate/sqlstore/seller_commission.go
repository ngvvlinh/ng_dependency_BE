package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/convert"
	"etop.vn/backend/com/services/affiliate/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type SellerCommissionStoreFactory func(ctx context.Context) *SellerCommissionStore

func NewSellerCommissionSettingStore(db *cmsql.Database) SellerCommissionStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *SellerCommissionStore {
		return &SellerCommissionStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type SellerCommissionStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft SellerCommissionFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *SellerCommissionStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.SellerCommission)(nil))
}

func (s *SellerCommissionStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *SellerCommissionStore) Paging(paging meta.Paging) *SellerCommissionStore {
	s.paging = paging
	return s
}

func (s *SellerCommissionStore) Filters(filters meta.Filters) *SellerCommissionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *SellerCommissionStore) ID(id dot.ID) *SellerCommissionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SellerCommissionStore) SellerID(id dot.ID) *SellerCommissionStore {
	s.preds = append(s.preds, s.ft.BySellerID(id))
	return s
}

func (s *SellerCommissionStore) OrderID(id dot.ID) *SellerCommissionStore {
	s.preds = append(s.preds, s.ft.ByOrderId(id))
	return s
}

func (s *SellerCommissionStore) GetAffiliateCommissionDB() (*model.SellerCommission, error) {
	var affiliateCommission model.SellerCommission
	err := s.query().Where(s.preds).ShouldGet(&affiliateCommission)
	return &affiliateCommission, err
}

func (s *SellerCommissionStore) GetAffiliateCommissionsDB() ([]*model.SellerCommission, error) {
	var affiliateCommission model.SellerCommissions
	err := s.query().Where(s.preds).Find(&affiliateCommission)
	return affiliateCommission, err
}

func (s *SellerCommissionStore) GetAffiliateCommission() (*affiliate.SellerCommission, error) {
	affiliateCommission, err := s.GetAffiliateCommissionDB()
	if err != nil {
		return nil, err
	}
	return convert.SellerCommission(affiliateCommission), nil
}

func (s *SellerCommissionStore) GetAffiliateCommissions() ([]*affiliate.SellerCommission, error) {
	query := s.query().Where(s.preds)
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.paging, SortSellerCommission)
	if err != nil {
		return nil, err
	}

	var results model.SellerCommissions
	err = query.Find(&results)
	return convert.AffiliateCommissions(results), err
}

func (s *SellerCommissionStore) CreateAffiliateCommission(sellerCommission *model.SellerCommission) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(sellerCommission)
	return err
}

func (s *SellerCommissionStore) UpdateAffiliateCommission(sellerCommission *model.SellerCommission) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ID(sellerCommission.ID).query().Where(s.preds).Update(sellerCommission)
	return err
}
