package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbuser/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

type FbExternalUserInternalFactory func(ctx context.Context) *FbExternalUserInternalStore

func NewFbExternalUserInternalStore(db *cmsql.Database) FbExternalUserInternalFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalUserInternalStore {
		return &FbExternalUserInternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalUserInternalStore struct {
	ft FbExternalUserInternalFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalUserInternalStore) ExternalIDs(externalIDs []string) *FbExternalUserInternalStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalUserInternalStore) ExternalID(externalID string) *FbExternalUserInternalStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalUserInternalStore) CreateOrUpdateFbExternalUserInternal(fbExternalUserInternal *fbusering.FbExternalUserInternal) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalUserInternalDB := new(model.FbExternalUserInternal)
	if err := scheme.Convert(fbExternalUserInternal, fbExternalUserInternalDB); err != nil {
		return err
	}

	if _, err := s.query().Upsert(fbExternalUserInternalDB); err != nil {
		return err
	}

	var tempFbUserInternal model.FbExternalUserInternal
	if err := s.query().Where(s.ft.ByExternalID(fbExternalUserInternal.ExternalID)).ShouldGet(&tempFbUserInternal); err != nil {
		return err
	}
	fbExternalUserInternal.UpdatedAt = tempFbUserInternal.UpdatedAt
	return nil
}

func (s *FbExternalUserInternalStore) GetFbExternalUserInternalDB() (*model.FbExternalUserInternal, error) {
	query := s.query().Where(s.preds)

	var fbExternalUser model.FbExternalUserInternal
	err := query.ShouldGet(&fbExternalUser)
	return &fbExternalUser, err
}

func (s *FbExternalUserInternalStore) GetFbExternalUserInternal() (*fbusering.FbExternalUserInternal, error) {
	fbExternalUser, err := s.GetFbExternalUserInternalDB()
	if err != nil {
		return nil, err
	}
	result := &fbusering.FbExternalUserInternal{}
	err = scheme.Convert(fbExternalUser, result)
	if err != nil {
		return nil, err
	}
	return result, err
}
