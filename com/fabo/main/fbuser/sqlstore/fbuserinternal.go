package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbuser/model"
	"o.o/backend/pkg/common/sql/cmsql"
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

func (s *FbExternalUserInternalStore) CreateFbExternalUserInternal(fbExternalUserInternal *fbusering.FbExternalUserInternal) error {
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
