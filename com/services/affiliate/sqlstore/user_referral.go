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

type UserReferralStoreFactory func(ctx context.Context) *UserReferralStore

func NewUserReferralStore(db *cmsql.Database) UserReferralStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *UserReferralStore {
		return &UserReferralStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type UserReferralStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft UserReferralFilters

	sqlstore.Paging
	filters meta.Filters
}

func (s *UserReferralStore) UserID(id dot.ID) *UserReferralStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *UserReferralStore) ReferralID(id dot.ID) *UserReferralStore {
	s.preds = append(s.preds, s.ft.ByReferralID(id))
	return s
}

func (s *UserReferralStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ProductPromotion)(nil))
}

func (s *UserReferralStore) WithPaging(paging meta.Paging) *UserReferralStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *UserReferralStore) Filters(filters meta.Filters) *UserReferralStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *UserReferralStore) GetUserReferralDB() (*model.UserReferral, error) {
	var userReferral model.UserReferral
	err := s.query().Where(s.preds).ShouldGet(&userReferral)
	return &userReferral, err
}

func (s *UserReferralStore) GetUserReferral() (*affiliate.UserReferral, error) {
	userReferral, err := s.GetUserReferralDB()
	return convert.UserReferral(userReferral), err
}

func (s *UserReferralStore) GetUserReferrals() ([]*affiliate.UserReferral, error) {
	var userReferrals model.UserReferrals
	err := s.query().Where(s.preds).Find(&userReferrals)
	return convert.UserReferrals(userReferrals), err
}

func (s *UserReferralStore) CreateUserReferral(userReferral *model.UserReferral) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(userReferral)
	return err
}

func (s *UserReferralStore) UpdateUserReferral(userReferral *model.UserReferral) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.UserID(userReferral.UserID).query().Where(s.preds).Update(userReferral)
	return err
}
