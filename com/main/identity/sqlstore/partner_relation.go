package sqlstore

import (
	"context"

	"o.o/api/main/identity"
	"o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type PartnerRelationStoreFactory func(context.Context) *PartnerRelationStore

func NewPartnerRelationStore(db *cmsql.Database) PartnerRelationStoreFactory {
	return func(ctx context.Context) *PartnerRelationStore {
		return &PartnerRelationStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type PartnerRelationStore struct {
	query cmsql.QueryFactory
	preds []interface{}
	ft    PartnerRelationFilters

	includeDeleted sqlstore.IncludeDeleted
}

func (s *PartnerRelationStore) BySubjectType(_type identity.SubjectType) *PartnerRelationStore {
	s.preds = append(s.preds, s.ft.BySubjectType(_type))
	return s
}

func (s *PartnerRelationStore) BySubjectIDs(ids ...dot.ID) *PartnerRelationStore {
	s.preds = append(s.preds, sq.In("subject_id", ids))
	return s
}

func (s *PartnerRelationStore) ListPartnerRelationsDB() ([]*model.PartnerRelation, error) {
	var partnerRelations model.PartnerRelations
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	err := query.Find(&partnerRelations)
	return partnerRelations, err
}

func (s *PartnerRelationStore) ListPartnerRelations() (partners []*identity.PartnerRelation, err error) {
	partnerRelationsDB, err := s.ListPartnerRelationsDB()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(partnerRelationsDB, &partners)
	return partners, err
}
