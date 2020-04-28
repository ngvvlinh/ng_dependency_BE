package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbuser/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
)

type FbUserInternalFactory func(ctx context.Context) *FbUserInternalStore

func NewFbUserInternalStore(db *cmsql.Database) FbUserInternalFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbUserInternalStore {
		return &FbUserInternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbUserInternalStore struct {
	ft FbUserInternalFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbUserInternalStore) CreateFbUserInternal(fbUserInternal *fbusering.FbUserInternal) error {
	sqlstore.MustNoPreds(s.preds)
	fbUserInternalDB := new(model.FbUserInternal)
	if err := scheme.Convert(fbUserInternal, fbUserInternalDB); err != nil {
		return err
	}
	_, err := s.query().Upsert(fbUserInternalDB)
	if err != nil {
		return err
	}

	var tempFbUserInternal model.FbUserInternal
	if err := s.query().Where(s.ft.ByID(fbUserInternal.ID)).ShouldGet(&tempFbUserInternal); err != nil {
		return err
	}
	fbUserInternal.UpdatedAt = tempFbUserInternal.UpdatedAt

	return nil
}
