package sqlstore

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/backend/com/main/identity/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
)

type PartnerStoreFactory func(context.Context) *PartnerStore

func NewPartnerStore(db *cmsql.Database) PartnerStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *PartnerStore {
		return &PartnerStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type PartnerStore struct {
	query cmsql.QueryFactory
	preds []interface{}
	ft    PartnerFilters
}

func (s *PartnerStore) WhiteLabel() *PartnerStore {
	s.preds = append(s.preds, sq.NewExpr("white_label_key IS NOT NULL"))
	return s
}

func (s *PartnerStore) ListPartnersDB() ([]*model.Partner, error) {
	var partners model.Partners
	err := s.query().Where(s.preds).Find(&partners)
	return partners, err
}

func (s *PartnerStore) ListPartners() (partners []*identity.Partner, err error) {
	partnersDB, err := s.ListPartnersDB()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(partnersDB, &partners)
	return partners, err
}
