package sqlstore

import (
	"context"

	"o.o/api/main/identity"
	"o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/capi/dot"
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

func (s *PartnerStore) ByID(id dot.ID) *PartnerStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PartnerStore) GetPartnerDB() (*model.Partner, error) {
	var partner model.Partner
	query := s.query().Where(s.preds)
	err := query.ShouldGet(&partner)
	return &partner, err
}

func (s *PartnerStore) GetPartner() (partner *identity.Partner, _ error) {
	partnerDB, err := s.GetPartnerDB()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(partnerDB, partner)
	return partner, err
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
