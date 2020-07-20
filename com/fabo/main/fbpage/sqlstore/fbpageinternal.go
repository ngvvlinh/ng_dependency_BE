package sqlstore

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbpage/convert"
	"o.o/backend/com/fabo/main/fbpage/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalPageInternalStoreFactory func(ctx context.Context) *FbExternalPageInternalStore

func NewFbExternalPageInternalStore(db *cmsql.Database) FbExternalPageInternalStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalPageInternalStore {
		return &FbExternalPageInternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalPageInternalStore struct {
	ft FbExternalPageInternalFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalPageInternalStore) ID(id dot.ID) *FbExternalPageInternalStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbExternalPageInternalStore) ExternalID(externalID string) *FbExternalPageInternalStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalPageInternalStore) CreateFbExternalPageInternal(fbExternalPageInternal *fbpaging.FbExternalPageInternal) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPageInternalDB := new(model.FbExternalPageInternal)
	if err := scheme.Convert(fbExternalPageInternal, fbExternalPageInternalDB); err != nil {
		return err
	}

	_, err := s.query().Insert(fbExternalPageInternalDB)
	if err != nil {
		return err
	}

	var tempFbPageInternal model.FbExternalPageInternal
	if err := s.query().Where(s.ft.ByID(fbExternalPageInternal.ID)).ShouldGet(&tempFbPageInternal); err != nil {
		return err
	}
	fbExternalPageInternal.UpdatedAt = tempFbPageInternal.UpdatedAt

	return nil
}

func (s *FbExternalPageInternalStore) CreateFbExternalPageInternals(fbExternalPageInternals []*fbpaging.FbExternalPageInternal) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPageInternalsDB := model.FbExternalPageInternals(convert.Convert_fbpaging_FbExternalPageInternals_fbpagemodel_FbExternalPageInternals(fbExternalPageInternals))

	_, err := s.query().Upsert(&fbExternalPageInternalsDB)
	if err != nil {
		return err
	}

	return nil
}

func (s *FbExternalPageInternalStore) ListFbPageInternalsDB() ([]*model.FbExternalPageInternal, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalPageInternal, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalPageInternal)
	if err != nil {
		return nil, err
	}

	var fbExternalPageInternals model.FbExternalPageInternals
	err = query.Find(&fbExternalPageInternals)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalPageInternals)
	return fbExternalPageInternals, nil
}

func (s *FbExternalPageInternalStore) ListFbPageInternals() (result []*fbpaging.FbExternalPageInternal, err error) {
	fbExternalPageInternals, err := s.ListFbPageInternalsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalPageInternals, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalPageInternalStore) GetFbExternalPageInternalDB() (*model.FbExternalPageInternal, error) {
	query := s.query().Where(s.preds)

	var fbExternalPageInternal model.FbExternalPageInternal
	err := query.ShouldGet(&fbExternalPageInternal)
	return &fbExternalPageInternal, err
}

func (s *FbExternalPageInternalStore) GetFbExternalPageInternal() (*fbpaging.FbExternalPageInternal, error) {
	fbExternalPageInternal, err := s.GetFbExternalPageInternalDB()
	if err != nil {
		return nil, err
	}
	result := &fbpaging.FbExternalPageInternal{}
	err = scheme.Convert(fbExternalPageInternal, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbExternalPageInternalStore) GetAccessToken() (string, error) {
	query := s.query().Where(s.preds)
	var fbExternalPageInternal model.FbExternalPageInternal
	err := query.ShouldGet(&fbExternalPageInternal)
	if err != nil {
		return "", err
	}
	return fbExternalPageInternal.Token, err
}
