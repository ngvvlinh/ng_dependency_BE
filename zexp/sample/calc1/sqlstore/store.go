package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/zexp/sample/calc1/api"
	"o.o/backend/zexp/sample/calc1/model"
	"o.o/capi/dot"
)

type EquationStoreFactory func(ctx context.Context) *EquationStore

func NewEquationStore(db *cmsql.Database) EquationStoreFactory {
	return func(ctx context.Context) *EquationStore {
		return &EquationStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type EquationStore struct {
	ft      EquationFilters
	query   cmsql.QueryFactory
	preds   []interface{}
	ctx     context.Context
	filters meta.Filters
	sqlstore.Paging
}

func (s *EquationStore) WithPaging(paging meta.Paging) *EquationStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *EquationStore) ID(id dot.ID) *EquationStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *EquationStore) CreateEquationDB(eq *model.Equation) error {
	return s.query().ShouldInsert(eq)
}

func (s *EquationStore) CreateEquation(eq *api.Equation) error {
	eqDB := &model.Equation{
		ID:        eq.ID,
		Equation:  eq.Equation,
		Result:    eq.Result,
		CreatedAt: eq.CreatedAt,
		UpdatedAt: eq.UpdatedAt,
	}
	return s.CreateEquationDB(eqDB)
}

func (s *EquationStore) GetEquationDB() (*model.Equation, error) {
	eq := &model.Equation{}
	err := s.query().Where(s.preds).ShouldGet(eq)
	return eq, err
}

func (s *EquationStore) GetEquation() (*api.Equation, error) {
	eqDB, err := s.GetEquationDB()
	if err != nil {
		return nil, err
	}
	eq := &api.Equation{
		ID:        eqDB.ID,
		Equation:  eqDB.Equation,
		Result:    eqDB.Result,
		CreatedAt: eqDB.CreatedAt,
		UpdatedAt: eqDB.UpdatedAt,
	}
	return eq, nil
}
