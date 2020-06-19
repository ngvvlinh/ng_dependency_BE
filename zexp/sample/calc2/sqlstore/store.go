package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/zexp/sample/calc2/api"
	"o.o/backend/zexp/sample/calc2/convert"
	"o.o/backend/zexp/sample/calc2/model"
	"o.o/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

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

func (s *EquationStore) Filters(filters meta.Filters) *EquationStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
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
	eqDB := &model.Equation{}
	if err := scheme.Convert(eq, eqDB); err != nil {
		return err
	}
	return s.CreateEquationDB(eqDB)
}

func (s *EquationStore) UpdateEquation(eq *api.Equation) (int, error) {
	var eqDB *model.Equation
	eqDB = convert.Convert_api_Equation_calc2model_Equation(eq, eqDB)
	return s.UpdateEquationDB(eqDB)
}

func (s *EquationStore) UpdateEquationDB(eq *model.Equation) (int, error) {
	return s.ID(eq.ID).query().Where(s.preds).Update(eq)
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

	eq := &api.Equation{}
	if err := scheme.Convert(eqDB, eq); err != nil {
		return nil, err
	}
	return eq, nil
}

func (s *EquationStore) ListEquation() (*api.Equations, error) {
	eqDBs, _ := s.ListEquationDB()

	var eqs []*api.Equation
	eqs = convert.Convert_calc2model_Equations_api_Equations(eqDBs)
	res := &api.Equations{
		Equations: eqs,
	}
	return res, nil
}

func (s *EquationStore) ListEquationDB() ([]*model.Equation, error) {
	query := s.query().Where(s.preds)
	query, _, err := sqlstore.Filters(query, s.filters, FilterEquation)

	if err != nil {
		return nil, err
	}

	var equations model.Equations
	err = query.Find(&equations)
	if err != nil {
		return nil, err
	}

	return equations, nil
}
