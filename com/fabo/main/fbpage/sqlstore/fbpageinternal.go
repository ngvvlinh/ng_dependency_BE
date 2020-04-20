package sqlstore

import (
	"context"

	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/meta"
	"etop.vn/backend/com/fabo/main/fbpage/convert"
	"etop.vn/backend/com/fabo/main/fbpage/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
)

type FbPageInternalStoreFactory func(ctx context.Context) *FbPageInternalStore

func NewFbPageInternalStore(db *cmsql.Database) FbPageInternalStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbPageInternalStore {
		return &FbPageInternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbPageInternalStore struct {
	ft FbPageInternalFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbPageInternalStore) CreateFbPageInternal(fbPageInternal *fbpaging.FbPageInternal) error {
	sqlstore.MustNoPreds(s.preds)
	fbPageInternalDB := new(model.FbPageInternal)
	if err := scheme.Convert(fbPageInternal, fbPageInternalDB); err != nil {
		return err
	}

	_, err := s.query().Insert(fbPageInternalDB)
	if err != nil {
		return err
	}

	var tempFbPageInternal model.FbPageInternal
	if err := s.query().Where(s.ft.ByID(fbPageInternal.ID)).ShouldGet(&tempFbPageInternal); err != nil {
		return err
	}
	fbPageInternal.UpdatedAt = tempFbPageInternal.UpdatedAt

	return nil
}

func (s *FbPageInternalStore) CreateFbPageInternals(fbPageInternals []*fbpaging.FbPageInternal) error {
	sqlstore.MustNoPreds(s.preds)
	fbPageInternalsDB := model.FbPageInternals(convert.Convert_fbpaging_FbPageInternals_fbpagemodel_FbPageInternals(fbPageInternals))

	_, err := s.query().Upsert(&fbPageInternalsDB)
	if err != nil {
		return err
	}

	return nil
}

func (s *FbPageInternalStore) ListFbPageInternalsDB() ([]*model.FbPageInternal, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbPageInternal, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbPageInternal)
	if err != nil {
		return nil, err
	}

	var fbPageInternals model.FbPageInternals
	err = query.Find(&fbPageInternals)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbPageInternals)
	return fbPageInternals, nil
}

func (s *FbPageInternalStore) ListFbPageInternals() (result []*fbpaging.FbPageInternal, err error) {
	fbPageInternals, err := s.ListFbPageInternalsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbPageInternals, &result); err != nil {
		return nil, err
	}
	return
}
